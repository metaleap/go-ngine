package core

import (
	"runtime"

	glfw "github.com/go-gl/glfw"
)

var (
	//	Manages your main-thread's render loop. You'll need to call it's Loop() method once after go:ngine initialization (see examples).
	Loop EngineLoop
)

//	EngineLoop is a singleton type, only used for the core.Loop package-global exported variable.
//	It is only aware of that instance and does not support any other EngineLoop instances.
type EngineLoop struct {
	//	Set to true by Loop.Loop(). Set to false or call Loop.Stop() to stop looping.
	Looping bool

	On struct {
		//	While Loop.Loop() is running, this callback is invoked (in its own "app thread")
		//	every loop iteration (ie. once per frame).
		//	This callback may run in parallel with On.EverySec(), but never with On.WinThread().
		AppThread func()

		//	While Loop.Loop() is running, this callback is invoked (on the main windowing thread)
		//	every loop iteration (ie. once per frame).
		//	This callback is guaranteed to never run in parallel with
		//	(and always after) the On.AppThread() and On.EverySec() callbacks.
		WinThread func()

		//	While Loop.Loop() is running, this callback is invoked (on the main windowing thread)
		//	at least and at most once per second, a useful entry point for non-real-time periodically recurring code.
		//	Caution: unlike On.WinThread(), this callback runs in parallel with your On.AppThread() callback.
		EverySec func()
	}

	Tick struct {
		//	The tick-time when the Loop.On.EverySec() callback was last invoked.
		PrevSec int

		//	While Loop.Loop() is running, is set to the current "tick-time":
		//	the time in seconds expired ever since Loop.Loop() was last called.
		Now float64

		//	While Loop.Loop() is running, is set to the previous tick-time.
		Prev float64

		//	The delta between Tick.Prev and Tick.Now.
		Delta float64
	}

	//	If true (default), Loop() waits for the app and prep threads to finish & sync before
	//	first calling On.WinThread() and finally waiting for GPU vsync/buffer-swap.
	//	If false, Loop() allows the app and prep threads to continue running while waiting
	//	for GPU vsync/buffer-swap; then On.WinThread() is called when all of those 3 waits are over.
	// SwapLast bool
}

func (_ *EngineLoop) init() {
	Loop.On.EverySec, Loop.On.AppThread, Loop.On.WinThread = func() {}, func() {}, func() {}
}

//	Initiates a rendering loop. This method returns only when the loop is stopped for whatever reason.
//	
//	(Before entering the loop, this method performs a one-off GC invokation.)
func (_ *EngineLoop) Loop() {
	var (
		secTick int
		runGc   = Core.Options.Loop.GcEveryFrame
	)
	if !Loop.Looping {
		Loop.Looping = true
		glfw.SetTime(0)
		Loop.Tick.Now = glfw.Time()
		runtime.GC()
		Loop.Tick.Prev = Loop.Tick.Now
		Loop.Tick.Now = glfw.Time()
		Loop.Tick.PrevSec, Loop.Tick.Delta = int(Loop.Tick.Now), Loop.Tick.Now-Loop.Tick.Prev
		Stats.reset()
		Stats.FrameRenderBoth.comb1, Stats.FrameRenderBoth.comb2 = &Stats.FrameRenderCpu, &Stats.FrameRenderGpu
		Diag.LogIfGlErr("ngine.PreLoop")
		Diag.LogMisc("Enter loop...")
		Core.copyAppToPrep()
		Core.copyPrepToRend()
		Stats.enabled = false // Allow a "warm-up phase" for the first few frames (1sec max or less)
		for Loop.Looping {
			//	STEP 0. Fire off the prep thread (for next frame) and app thread (for next-after-next frame).
			thrApp.Lock()
			go Loop.onThreadApp()
			thrPrep.Lock()
			go Loop.onThreadPrep()

			//	STEP 1. Fill the GPU command queue with rendering commands (batched together by the previous prep thread)
			Stats.FrameRenderCpu.begin()
			Core.onRender()
			Stats.FrameRenderCpu.end()

			//	STEP 2. While the GL driver processes its command queue and other CPU cores work on
			//	app and prep threads, this CPU core can now perform some other minor duties
			Stats.fpsCounter++
			Stats.fpsAll++
			//	This branch runs at most and at least 1x per second
			if secTick = int(Loop.Tick.Now); secTick != Loop.Tick.PrevSec {
				Stats.FpsLastSec, Loop.Tick.PrevSec = Stats.fpsCounter, secTick
				Stats.fpsCounter = 0
				Core.onSec()
				Loop.On.EverySec()
				runGc, Stats.enabled = true, true
			}

			//	Wait for threads -- waits until both app and prep threads are done and copies stage states around
			Loop.onWaitForThreads()

			//	Call On.WinThread() -- for main-thread user code (mostly input polling) without affecting On.AppThread
			Loop.onThreadWin()

			//	Must do this here so that current-tick won't change half-way through OnAppTread(),
			//	and then we'd also like this frame's On.WinThread() to have the same current-tick.
			Loop.Tick.Prev = Loop.Tick.Now
			Loop.Tick.Now = glfw.Time()
			Loop.Tick.Delta = Loop.Tick.Now - Loop.Tick.Prev
			Stats.Frame.measureStartTime = Loop.Tick.Prev
			Stats.Frame.end()

			//	GC stops-the-world so do it after go-routines have finished. Now is a good time, as the GL driver
			//	is likely still processing its command queue from step 1. and won't be interrupted by Go's GC --
			//	the following buf-swap step block-waits for the GPU anyway.
			if runGc {
				runGc = Core.Options.Loop.GcEveryFrame
				Loop.onGC()
			}

			//	STEP 3. Swap buffers -- this waits for the GPU/GL to finish processing its command
			//	queue filled in Step 1, swap buffers and for V-sync (if any)
			Loop.onSwap()

			//	Check for resize before next render
			if (UserIO.Window.lastResize > 0) && ((Loop.Tick.Now - UserIO.Window.lastResize) > UserIO.Window.ResizeMinDelay) {
				UserIO.Window.lastResize = 0
				Core.onResizeWindow(UserIO.Window.width, UserIO.Window.height)
			}
			Stats.FrameRenderBoth.combine()
		}
		Loop.Looping = false
		Diag.LogMisc("Exited loop.")
		Diag.LogIfGlErr("ngine.PostLoop")
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
	if Core.Options.Loop.ForceThreads.App {
		runtime.LockOSThread()
	}
	Loop.On.AppThread()
	Stats.FrameAppThread.end()
	thrApp.Unlock()
}

func (_ *EngineLoop) onThreadPrep() {
	Stats.FramePrepThread.begin()
	if Core.Options.Loop.ForceThreads.Prep {
		runtime.LockOSThread()
	}
	Core.onPrep()
	Stats.FramePrepThread.end()
	thrPrep.Unlock()
}

func (_ *EngineLoop) onThreadWin() {
	Stats.FrameWinThread.begin()
	if glfw.PollEvents(); glfw.WindowParam(glfw.Opened) == 1 {
		Loop.On.WinThread()
	} else {
		Loop.Looping = false
	}
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
	Loop.Looping = false
}

//	Returns the number of seconds expired ever since Loop.Loop() was last called.
func (_ *EngineLoop) Time() float64 {
	return glfw.Time()
}
