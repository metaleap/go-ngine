package core

import (
	"math"
	"runtime"

	glfw "github.com/go-gl/glfw"
	ugl "github.com/go3d/go-opengl/util"
)

var (
	//	Manages your main-thread's "game loop". You'll need to call it's Loop() method once after go:ngine initialization (see samples).
	Loop EngineLoop
)

//	Consider EngineLoop a "Singleton" type, only valid use is the core.Loop global variable.
//	Manages your main-thread's "game loop".
type EngineLoop struct {
	//	Set to true by Loop.Loop(). Set to false to stop looping.
	IsLooping bool

	//	The tick-time when the Loop.OnSec() callback was last invoked.
	SecTickLast float64

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
	SwapLast bool
}

func (me *EngineLoop) init() {
	me.SwapLast, me.OnSec, me.OnAppThread, me.OnWinThread = true, func() {}, func() {}, func() {}
}

//	Initiates a rendering loop. This method returns only when the loop is stopped for whatever reason.
//	
//	(Before entering the loop, this method performs a one-off GC invokation.)
func (me *EngineLoop) Loop() {
	var tickNowFloor float64
	if !me.IsLooping {
		runtime.GC()
		me.IsLooping = true
		glfw.SetTime(0)
		me.SecTickLast, me.TickNow = glfw.Time(), glfw.Time()
		Stats.reset()
		Stats.FrameRenderBoth.comb1, Stats.FrameRenderBoth.comb2 = &Stats.FrameRenderCpu, &Stats.FrameRenderGpu
		ugl.LogLastError("ngine.PreLoop")
		Diag.LogMisc("Enter loop...")
		Core.copyAppToPrep()
		Core.copyPrepToRend()
		Stats.enabled = false // Allow a "warm-up phase" for the first few frames (1sec max or less)
		for me.IsLooping && (glfw.WindowParam(glfw.Opened) == 1) {
			//	STEP 0. Fire off the prep thread (for next frame) and app thread (for next-after-next frame).
			thrApp.Lock()
			go me.loopThreadApp()
			thrPrep.Lock()
			go me.loopThreadPrep()

			//	STEP 1. Send rendering commands (batched together in the previous prep thread) to the GPU / GL pipeline
			Stats.FrameRenderCpu.begin()
			//	Check for resize before render
			if (UserIO.lastWinResize > 0) && ((me.TickNow - UserIO.lastWinResize) > UserIO.WinResizeMinDelay) {
				UserIO.lastWinResize = 0
				Core.onResizeWindow(Core.Options.winWidth, Core.Options.winHeight)
			}
			Core.onRender()
			Stats.FrameRenderCpu.end()

			//	STEP 2. While those are processed GPU-side, do some minor non-rendering CPU-side stuff
			Stats.fpsCounter++
			Stats.fpsAll++
			me.TickLast = me.TickNow
			me.TickNow = glfw.Time()
			Stats.Frame.measureStartTime = me.TickLast
			Stats.Frame.end()
			me.TickDelta = me.TickNow - me.TickLast
			if tickNowFloor = math.Floor(me.TickNow); tickNowFloor != me.SecTickLast {
				Stats.FpsLastSec, me.SecTickLast = Stats.fpsCounter, tickNowFloor
				Stats.fpsCounter = 0
				Core.onSec()
				me.OnSec()
				Stats.Gc.begin()
				runtime.GC()
				Stats.Gc.end()
				Stats.enabled = true
			}

			//	STEP 3, 4 & 5:
			//	Wait for threads -- waits until both app and prep threads are done and copies stage states around
			//	Call OnWinThread() -- for main-thread user code (mostly input) without affecting OnAppThread
			//	Swap buffers -- wait for GPU to (a) finish processing commands sent in Step 1, (b) swap buffers and (c) for V-sync (if any)
			if me.SwapLast {
				me.loopWaitForThreads()
				me.loopOnWinThread()
				me.loopSwap()
			} else {
				me.loopSwap()
				me.loopWaitForThreads()
				me.loopOnWinThread()
			}
			Stats.FrameRenderBoth.combine()
		}
		me.IsLooping = false
		Diag.LogMisc("Exited loop.")
		ugl.LogLastError("ngine.PostLoop")
	}
}

func (me *EngineLoop) loopOnWinThread() {
	Stats.FrameWinThread.begin()
	glfw.PollEvents()
	me.OnWinThread()
	Stats.FrameWinThread.end()
}

func (_ *EngineLoop) loopSwap() {
	Stats.FrameRenderGpu.begin()
	glfw.SwapBuffers()
	Stats.FrameRenderGpu.end()
}

func (me *EngineLoop) loopThreadApp() {
	Stats.FrameAppThread.begin()
	me.OnAppThread()
	Stats.FrameAppThread.end()
	thrApp.Unlock()
}

func (me *EngineLoop) loopThreadPrep() {
	Stats.FramePrepThread.begin()
	Core.onPrep()
	Stats.FramePrepThread.end()
	thrPrep.Unlock()
}

func (me *EngineLoop) loopWaitForThreads() {
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
func (me *EngineLoop) Stop() {
	me.IsLooping = false
}

//	Returns the number of seconds expired ever since Loop.Loop() was last called.
func (_ *EngineLoop) Time() float64 {
	return glfw.Time()
}
