package core

import (
	"log"
	"math"
	"runtime"

	glfw "github.com/go-gl/glfw"
	ugl "github.com/go3d/go-glutil"
)

//	Consider EngineLoop a "Singleton" type, only valid use is the core.Loop global variable.
//	Manages your main-thread's "game loop".
type EngineLoop struct {
	//	Set to true by EngineLoop.Loop(). Set to false to stop looping.
	IsLooping bool

	//	The tick-time when the EngineLoop.OnSec() callback was last invoked.
	SecTickLast float64

	//	While EngineLoop.Loop() is running, is set to the current "tick-time":
	//	the time in seconds expired ever since EngineLoop.Loop() was last called.
	TickNow float64

	//	While EngineLoop.Loop() is running, is set to the previous tick-time.
	TickLast float64

	//	The delta between TickLast and TickNow.
	TickDelta float64

	//	While EngineLoop.Loop() is running, this callback is invoked every loop iteration (ie. once per frame).
	OnLoop func()

	//	While EngineLoop.Loop() is running, this callback is invoked at least and at most once per second.
	OnSec func()
}

func (me *EngineLoop) init() {
	me.OnSec, me.OnLoop = func() {}, func() {}
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
		log.Printf("Enter loop...")
		for me.IsLooping && (glfw.WindowParam(glfw.Opened) == 1) {
			//	STEP 1. Send rendering commands to the GPU / GL pipeline
			Stats.FrameRenderCpu.begin()
			Core.onRender()
			Stats.FrameRenderCpu.end()
			//	STEP 2. While those are sent (and executed GPU-side), do other non-rendering CPU-side stuff
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
			Stats.FrameCoreCode.begin()
			Core.onLoop()
			Stats.FrameCoreCode.end()
			Stats.FrameUserCode.begin()
			me.OnLoop()
			Stats.FrameUserCode.end()
			//	STEP 3. Swap buffers -- also waits for GPU to finish commands sent in Step 1, and for V-sync (if set).
			Stats.FrameRenderGpu.begin()
			glfw.SwapBuffers()
			Stats.FrameRenderGpu.end()
			Stats.FrameRenderBoth.combine()
			if (UserIO.lastWinResize > 0) && ((me.TickNow - UserIO.lastWinResize) > UserIO.WinResizeMinDelay) {
				UserIO.lastWinResize = 0
				Core.onResizeWindow(Core.Options.winWidth, Core.Options.winHeight)
			}
		}
		me.IsLooping = false
		log.Printf("Exited loop.")
		ugl.LogLastError("ngine.PostLoop")
	}
}

//	Stops the currently running EngineLoop.Loop().
func (me *EngineLoop) Stop() {
	me.IsLooping = false
}

//	Returns the number of seconds expired ever since EngineLoop.Loop() was last called.
func (me *EngineLoop) Time() float64 {
	return glfw.Time()
}
