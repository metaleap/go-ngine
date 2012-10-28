package core

import (
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
	glLogLastError("ngine.Pre_Dispose")
	if Core != nil { Core.Dispose() }
	glLogLastError("ngine.Core_Dispose")
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
		glLogLastError("ngine.PreLoop")
		log.Printf("Enter loop...")
		for me.IsLooping {
			Stats.FrameRender.begin(); Core.onRender(); Stats.FrameRender.end()
			Stats.fpsCounter++
			Stats.fpsAll++
			me.TickLast = me.TickNow
			me.TickNow = glfw.Time()
			Stats.Frame.measureStartTime = me.TickLast; Stats.Frame.end()
			me.TickDelta = me.TickNow - me.TickLast
			if tickNowFloor = math.Floor(me.TickNow); tickNowFloor != me.SecTickLast {
				Stats.FpsLastSec, me.SecTickLast = Stats.fpsCounter, tickNowFloor
				Stats.fpsCounter = 0
				Core.onSecTick()
				me.OnSecTick()
				Stats.Gc.begin(); runtime.GC(); Stats.Gc.end()
				// if Stats.isWarmup { Stats.reset(); Stats.isWarmup = false }
			}
			Stats.FrameCoreCode.begin(); Core.onLoop(); Stats.FrameCoreCode.end()
			Stats.FrameUserCode.begin(); me.OnLoop(); Stats.FrameUserCode.end()
			Stats.FrameSwap.begin(); UserIO.onLoop(); Stats.FrameSwap.end()
		}
		glLogLastError("ngine.PostLoop")
	}
}

func (me *tEngineLoop) Time () float64 {
	return glfw.Time()
}
