package core

import (
	"fmt"
	"log"
	"math"
	"runtime"

	glfw "github.com/go-gl/glfw"
)

var (
	Loop *tEngineLoop
	Core *tEngineCore
	Stats *tEngineStats
	UserIO = newUserIO()
)

func Dispose () {
	if Core != nil { Core.Dispose() }
	glDispose()
	UserIO.dispose()
	Core, Loop, Stats = nil, nil, nil
}

func Init (options *tOptions, winTitle string) error {
	var err error
	var isVerErr bool
	var forceContext = false
	tryInit:
	if err = UserIO.init(options, winTitle, forceContext); err == nil {
		if err, isVerErr = glInit(); err == nil {
			Loop, Stats, Core = newEngineLoop(), newEngineStats(), newEngineCore(options)
		} else if isVerErr && !forceContext {
			forceContext = true
			UserIO.isGlfwInit, UserIO.isGlfwWindow = false, false
			goto tryInit
		}
	}
	return err
}

func Key (format string, fmtArgs ... interface {}) string {
	return fmt.Sprintf(format, fmtArgs ...)
}

func LogError (err error) {
	if err != nil { log.Println(err.Error()) }
}

func LogLastGlError (step string, fmtArgs ... interface {}) {
	glLogLastError(step, fmtArgs ...)
}

type tEngineLoop struct {
	IsLooping bool
	SecTickLast, TickNow, TickLast, TickDelta float64
	OnLoop, OnSecTick func()
}

func newEngineLoop () *tEngineLoop {
	var loop = &tEngineLoop {}
	loop.OnSecTick, loop.OnLoop = func () {}, func () {}
	return loop
}

func (me *tEngineLoop) Loop () {
	var tickNowFloor float64
	if (!me.IsLooping) {
		me.IsLooping = true
		glfw.SetTime(0)
		me.SecTickLast, me.TickNow = glfw.Time(), glfw.Time()
		Stats.reset()
		Stats.FrameRenderBoth.comb1, Stats.FrameRenderBoth.comb2 = Stats.FrameRenderCpu, Stats.FrameRenderGpu
		glLogLastError("ngine.PreLoop")
		log.Printf("Enter loop...")
		for me.IsLooping && (glfw.WindowParam(glfw.Opened) == 1) {
			//	STEP 1. Send rendering commands to the GPU / GL pipeline
			Stats.FrameRenderCpu.begin(); Core.onRender(); Stats.FrameRenderCpu.end()
			//	STEP 2. While those are sent (and executed GPU-side), do other non-rendering CPU-side stuff
			Stats.fpsCounter++
			Stats.fpsAll++
			me.TickLast = me.TickNow
			me.TickNow = glfw.Time()
			Stats.Frame.measureStartTime = me.TickLast; Stats.Frame.end()
			me.TickDelta = me.TickNow - me.TickLast
			if tickNowFloor = math.Floor(me.TickNow); tickNowFloor != me.SecTickLast {
				Stats.FpsLastSec, me.SecTickLast = Stats.fpsCounter, tickNowFloor
				Stats.fpsCounter = 0; Core.onSecTick(); me.OnSecTick()
				Stats.Gc.begin(); runtime.GC(); Stats.Gc.end()
			}
			Stats.FrameCoreCode.begin(); Core.onLoop(); Stats.FrameCoreCode.end()
			Stats.FrameUserCode.begin(); me.OnLoop(); Stats.FrameUserCode.end()
			//	STEP 3. Swap buffers -- also waits for GPU to finish commands sent in Step 1, and for V-sync (if set).
			Stats.FrameRenderGpu.begin(); glfw.SwapBuffers(); Stats.FrameRenderGpu.end()
			Stats.FrameRenderBoth.combine()
		}
		me.IsLooping = false
		glLogLastError("ngine.PostLoop")
	}
}

func (me *tEngineLoop) Stop () {
	me.IsLooping = false
}

func (me *tEngineLoop) Time () float64 {
	return glfw.Time()
}
