package core

import (
	"runtime"

	glfw "github.com/go-gl/glfw"
	ugl "github.com/go3d/go-opengl/util"
)

var (
	//	Manages your main-thread's "game loop". You'll need to call it's Loop() method once after go:ngine initialization (see samples).
	Loop EngineLoop
)

//	EngineLoop is a singleton type, only used for the core.Loop package-global exported variable.
//	It is only aware of that instance and does not support any other EngineLoop instances.
type EngineLoop struct {
	//	Set to true by Loop.Loop(). Set to false to stop looping.
	IsLooping bool

	//	The tick-time when the Loop.OnSec() callback was last invoked.
	SecTickLast int

	//	While Loop.Loop() is running, is set to the current "tick-time":
	//	the time in seconds expired ever since Loop.Loop() was last called.
	TickNow float64

	//	While Loop.Loop() is running, is set to the previous tick-time.
	TickLast float64

	//	The delta between TickLast and TickNow.
	TickDelta float64

	//	While Loop.Loop() is running, this callback is invoked (in its own "app thread")
	//	every loop iteration (ie. once per frame).
	//	This callback may run in parallel with OnSec(), but never with OnWinThread().
	OnAppThread func()

	//	While Loop.Loop() is running, this callback is invoked (on the main windowing thread)
	//	every loop iteration (ie. once per frame).
	//	This callback is guaranteed to never run in parallel with
	//	(and always after) the OnAppThread() and OnSec() callbacks.
	OnWinThread func()

	//	While Loop.Loop() is running, this callback is invoked (on the main windowing thread)
	//	at least and at most once per second, a useful entry point for non-real-time periodically recurring code.
	//	Caution: unlike OnWinThread(), this callback runs in parallel with your OnAppThread() callback.
	OnSec func()

	//	If true (default), Loop() waits for the app and prep threads to finish & sync before
	//	first calling OnWinThread() and finally waiting for GPU vsync/buffer-swap.
	//	If false, Loop() allows the app and prep threads to continue running while waiting
	//	for GPU vsync/buffer-swap; then OnWinThread() is called when all of those 3 waits are over.
	// SwapLast bool
}

func (_ *EngineLoop) init() {
	Loop.OnSec, Loop.OnAppThread, Loop.OnWinThread = func() {}, func() {}, func() {}
}

//	Initiates a rendering loop. This method returns only when the loop is stopped for whatever reason.
//	
//	(Before entering the loop, this method performs a one-off GC invokation.)
func (_ *EngineLoop) Loop() {
	var (
		secTick int
		runGc   bool
	)
	if !Loop.IsLooping {
		Loop.IsLooping = true
		runtime.GC()
		glfw.SetTime(0)
		Loop.SecTickLast, Loop.TickNow = int(glfw.Time()), glfw.Time()
		Stats.reset()
		Stats.FrameRenderBoth.comb1, Stats.FrameRenderBoth.comb2 = &Stats.FrameRenderCpu, &Stats.FrameRenderGpu
		ugl.LogLastError("ngine.PreLoop")
		Diag.LogMisc("Enter loop...")
		Core.copyAppToPrep()
		Core.copyPrepToRend()
		Stats.enabled = false // Allow a "warm-up phase" for the first few frames (1sec max or less)
		for Loop.IsLooping && (glfw.WindowParam(glfw.Opened) == 1) {
			//	STEP 0. Fire off the prep thread (for next frame) and app thread (for next-after-next frame).
			thrApp.Lock()
			go Loop.onThreadApp()
			thrPrep.Lock()
			go Loop.onThreadPrep()

			//	STEP 1. Fill the GPU command queue with rendering commands (batched together by the previous prep thread)
			//	Check for resize before render
			Stats.FrameRenderCpu.begin()
			if (UserIO.lastWinResize > 0) && ((Loop.TickNow - UserIO.lastWinResize) > UserIO.WinResizeMinDelay) {
				UserIO.lastWinResize = 0
				Core.onResizeWindow(Core.Options.winWidth, Core.Options.winHeight)
			}
			Core.onRender()
			Stats.FrameRenderCpu.end()

			//	STEP 2. While the GL driver processes its command queue and other CPU cores work on
			//	app and prep threads, this CPU core can now perform some other minor duties
			Stats.fpsCounter++
			Stats.fpsAll++
			Loop.TickLast = Loop.TickNow
			Loop.TickNow = glfw.Time()
			Stats.Frame.measureStartTime = Loop.TickLast
			Stats.Frame.end()
			Loop.TickDelta = Loop.TickNow - Loop.TickLast
			//	This branch runs at most and at least 1x per second
			if secTick = int(Loop.TickNow); secTick != Loop.SecTickLast {
				Stats.FpsLastSec, Loop.SecTickLast = Stats.fpsCounter, secTick
				Stats.fpsCounter = 0
				Core.onSec()
				Loop.OnSec()
				runGc, Stats.enabled = true, true
			}

			//	Wait for threads -- waits until both app and prep threads are done and copies stage states around
			Loop.onWaitForThreads()

			//	Call OnWinThread() -- for main-thread user code (mostly input polling) without affecting OnAppThread
			Loop.onThreadWin()

			//	GC stops-the-world so do it after go-routines have finished. Now is a good time, as the GL driver
			//	may still be busy processing its command queue from step 1. and won't be interrupted by Go's GC.
			if runGc {
				runGc = false
				Loop.onGC()
			}

			//	STEP 3. Swap buffers -- this waits for the GPU/GL to finish processing its command
			//	queue filled in Step 1, swap buffers and for V-sync (if any)
			Loop.onSwap()
			Stats.FrameRenderBoth.combine()
		}
		Loop.IsLooping = false
		Diag.LogMisc("Exited loop.")
		ugl.LogLastError("ngine.PostLoop")
	}
}

func (_ *EngineLoop) onGC() {
	Stats.Gc.begin()
	runtime.GC()
	Stats.Gc.end()
}

func (_ *EngineLoop) onSwap() {
	Stats.FrameRenderGpu.begin()
	glfw.SwapBuffers()
	Stats.FrameRenderGpu.end()
}

func (_ *EngineLoop) onThreadApp() {
	Stats.FrameAppThread.begin()
	Loop.OnAppThread()
	Stats.FrameAppThread.end()
	thrApp.Unlock()
}

func (_ *EngineLoop) onThreadPrep() {
	Stats.FramePrepThread.begin()
	Core.onPrep()
	Stats.FramePrepThread.end()
	thrPrep.Unlock()
}

func (_ *EngineLoop) onThreadWin() {
	Stats.FrameWinThread.begin()
	glfw.PollEvents()
	Loop.OnWinThread()
	Stats.FrameWinThread.end()
}

func (_ *EngineLoop) onWaitForThreads() {
	Stats.FrameThreadSync.begin()

	thrPrep.Lock()
	Core.copyPrepToRend()
	thrPrep.Unlock()

	thrApp.Lock()
	Core.copyAppToPrep()
	thrApp.Unlock()

	Stats.FrameThreadSync.end()
}

//	Stops the currently running Loop.Loop().
func (_ *EngineLoop) Stop() {
	Loop.IsLooping = false
}

//	Returns the number of seconds expired ever since Loop.Loop() was last called.
func (_ *EngineLoop) Time() float64 {
	return glfw.Time()
}
