package core

import (
	"runtime"
	"time"

	glfw "github.com/go-gl/glfw"
)

var (
	//	Manages your main-thread's render loop.
	//	Call it's Run() method once after go:ngine initialization (see examples).
	Loop NgLoop
)

//	NgLoop is a singleton type, only used for the Loop variable.
//	It is only aware of that instance and does not support any other NgLoop instances.
type NgLoop struct {
	//	Set to true by Loop.Run(). Set to false to stop looping.
	Running bool

	Delay time.Duration

	MaxIterations float64

	On struct {
		//	While Loop.Run() is running, this callback is invoked (in its own "app thread")
		//	every loop iteration (ie. once per frame).
		//	This callback may run in parallel with On.EverySec(), but never with On.WinThread().
		AppThread func()

		//	While Loop.Run() is running, this callback is invoked (on the main windowing thread)
		//	every loop iteration (ie. once per frame).
		//	This callback is guaranteed to never run in parallel with
		//	(and always after) the On.AppThread() and On.EverySec() callbacks.
		WinThread func()

		//	While Loop.Run() is running, this callback is invoked (on the main windowing thread)
		//	at least and at most once per second, a useful entry point for non-real-time periodically recurring code.
		//	Caution: unlike On.WinThread(), this callback runs in parallel with your On.AppThread() callback.
		EverySec func()
	}

	Tick struct {
		//	The tick-time when the Loop.On.EverySec() callback was last invoked.
		PrevSec int

		//	While Loop.Run() is running, is set to the current "tick-time":
		//	the time in seconds expired ever since Loop.Run() was last called.
		Now float64

		//	While Loop.Run() is running, is set to the previous tick-time.
		Prev float64

		//	The delta between Tick.Prev and Tick.Now.
		Delta float64
	}
}

func (_ *NgLoop) init() {
	Loop.On.EverySec, Loop.On.AppThread, Loop.On.WinThread = func() {}, func() {}, func() {}
}

func (_ *NgLoop) onGC() {
	Stats.Gc.begin()
	runtime.GC()
	Stats.Gc.end()
}

func (_ *NgLoop) onSwap() {
	Stats.FrameRenderGpu.begin()
	glfw.SwapBuffers()
	Stats.FrameRenderGpu.end()
}

func (_ *NgLoop) onThreadApp() {
	Stats.FrameAppThread.begin()
	Loop.On.AppThread()
	Stats.FrameAppThread.end()
	thrApp.Unlock()
}

func (_ *NgLoop) onThreadPrep() {
	Stats.FramePrepThread.begin()
	Core.onPrep()
	Stats.FramePrepThread.end()
	thrPrep.Unlock()
}

func (_ *NgLoop) onThreadWin() {
	Stats.FrameWinThread.begin()
	if glfw.PollEvents(); glfw.WindowParam(glfw.Opened) == 1 {
		Loop.On.WinThread()
	} else {
		Loop.Running = false
	}
	Stats.FrameWinThread.end()
}

func (_ *NgLoop) onWaitForThreads() {
	Stats.FrameThreadSync.begin()

	thrPrep.Lock()
	Core.copyPrepToRend()
	thrPrep.Unlock()

	thrApp.Lock()
	Core.copyAppToPrep()
	thrApp.Unlock()

	Stats.FrameThreadSync.end()
}

//	Initiates a rendering loop. This method returns only when the loop is stopped for whatever reason.
//	
//	(Before entering the loop, this method performs a one-off GC invokation.)
func (_ *NgLoop) Run() {
	var (
		secTick int
		runGc   = Options.Loop.GcEvery.Frame
	)
	if !Loop.Running {
		Loop.Running = true
		glfw.SetTime(0)
		Loop.Tick.Now = glfw.Time()
		Core.copyAppToPrep()
		Core.copyPrepToRend()
		Loop.Tick.Prev = Loop.Tick.Now
		Loop.Tick.Now = glfw.Time()
		Loop.Tick.PrevSec, Loop.Tick.Delta = int(Loop.Tick.Now), Loop.Tick.Now-Loop.Tick.Prev
		Stats.reset()
		runtime.GC()
		Diag.LogMisc("Enter loop...")
		Loop.Running = glfw.WindowParam(glfw.Opened) == 1
		for Loop.Running {
			//	STEP 0. Fire off the prep thread (for next frame) and app thread (for next-after-next frame).
			thrPrep.Lock()
			go Loop.onThreadPrep()
			thrApp.Lock()
			go Loop.onThreadApp()

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
				Stats.FrameOnSec.begin()
				Stats.FpsLastSec, Loop.Tick.PrevSec = Stats.fpsCounter, secTick
				Stats.fpsCounter = 0
				if Diag.LogGLErrorsInLoopOnSec {
					Diag.LogIfGlErr("onSec")
				}
				Loop.On.EverySec()
				runGc = Options.Loop.GcEvery.Sec
				Stats.FrameOnSec.end()
				Stats.enable() // the first few frames were warm-ups that don't count towards the stats
			}
			//	Wait for threads -- waits until both app and prep threads are done and copies stage states around
			Loop.onWaitForThreads()
			//	Call On.WinThread() -- for main-thread user code (mostly input polling) without affecting On.AppThread
			Loop.onThreadWin()

			//	Must do this here so that current-tick won't change half-way through OnAppTread(),
			//	and then we'd also like this frame's On.WinThread() to have the same current-tick.
			Loop.Tick.Prev, Loop.Tick.Now = Loop.Tick.Now, glfw.Time()
			Stats.Frame.measureStartTime, Loop.Tick.Delta = Loop.Tick.Prev, Loop.Tick.Now-Loop.Tick.Prev
			Stats.Frame.end()
			//	GC stops-the-world so do it after go-routines have finished. Now is a good time, as the GPU
			//	is likely still busy processing commands from step 1 and won't be interrupted by Go's GC --
			//	the subsequent buffer-swap step block-waits for the GPU anyway.
			if runGc {
				runGc = Options.Loop.GcEvery.Frame
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
			Stats.FrameRenderBoth.combine(&Stats.FrameRenderCpu, &Stats.FrameRenderGpu)
			if Loop.Delay > 0 {
				time.Sleep(Loop.Delay)
			}
			if Loop.MaxIterations > 0 && Stats.fpsAll >= Loop.MaxIterations {
				Loop.Running = false
			}
		}
		Loop.Running = false
		Diag.LogMisc("Exited loop.")
		Diag.LogIfGlErr("ngine.PostLoop")
		for rbi, rbe := range Core.Render.Canvases[1].Views[0].Technique_Scene().thrRend.batch.all {
			println(strf("%d\t==>\tP:%v\tT:%v\tB:%v\tD:%v", rbi, rbe.prog, rbe.texes, Core.Libs.Meshes[rbe.mesh].meshBuffer.glIbo.GlHandle, rbe.dist))
		}
	}
}

//	Returns the number of seconds expired ever since Loop.Run() was last called.
func (_ *NgLoop) Time() float64 {
	return glfw.Time()
}
