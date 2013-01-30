package core

import (
	"math"
	"runtime"

	glfw "github.com/go-gl/glfw"
	ugl "github.com/go3d/go-glutil"
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
	//	at least and at most once per second.
	//	Caution: unlike OnWinThread(), this callback most likely runs in parallel with your OnAppThread() callback.
	OnSec func()

	//	If true, Loop() waits for the app and prep threads to finish before swapping buffers.
	//	If false, Loop() allows the app and prep threads to continue running while swapping buffers.
	//	Defaults to false. Setting this to true may prove beneficial if your OnAppThread() callback
	//	isn't doing any computationally intensive work.
	SwapLast bool

	threadBusy struct {
		app, prep bool
	}
}

func (me *EngineLoop) init() {
	me.OnSec, me.OnAppThread, me.OnWinThread = func() {}, func() {}, func() {}
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
		for me.IsLooping && (glfw.WindowParam(glfw.Opened) == 1) {
			//	STEP 0. Fire off the prep thread (for next frame) and app thread (for next-after-next frame).
			go me.loopThreadApp()
			go me.loopThreadPrep()

			//	STEP 1. Send rendering commands (batched together in the previous prep thread) to the GPU / GL pipeline
			Stats.FrameRenderCpu.begin()
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
			}

			//	STEP 3, 4 & 5:
			//	Swap buffers -- also waits for GPU to finish commands sent in Step 1, and for V-sync (if set).
			//	Call OnWinThread() -- for main-thread user code without affecting the OnAppThread thread
			//	Wait for threads -- waits until both app and prep threads are done and copies stage states around.
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

			//	This is best handled right here in between the most-recent swap and before the next frame rendering.
			if (UserIO.lastWinResize > 0) && ((me.TickNow - UserIO.lastWinResize) > UserIO.WinResizeMinDelay) {
				UserIO.lastWinResize = 0
				Core.onResizeWindow(Core.Options.winWidth, Core.Options.winHeight)
			}
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
	me.threadBusy.app = true
	Stats.FrameAppThread.begin()
	me.OnAppThread()
	Stats.FrameAppThread.end()
	me.threadBusy.app = false
}

func (me *EngineLoop) loopThreadPrep() {
	me.threadBusy.prep = true
	Stats.FramePrepCode.begin()
	Core.onPrep()
	Stats.FramePrepCode.end()
	me.threadBusy.prep = false
}

func (me *EngineLoop) loopWaitForThreads() {
	for me.threadBusy.prep && me.IsLooping {
		runtime.Gosched()
	}
	Core.copyPrepToRend()
	for me.threadBusy.app && me.IsLooping {
		runtime.Gosched()
	}
	Core.copyAppToPrep()
}

//	Stops the currently running Loop.Loop().
func (me *EngineLoop) Stop() {
	me.IsLooping = false
}

//	Returns the number of seconds expired ever since Loop.Loop() was last called.
func (_ *EngineLoop) Time() float64 {
	return glfw.Time()
}
