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
	handlers []func()
	OnSecTick func()
}

func newEngineLoop () *tEngineLoop {
	var loop = &tEngineLoop {}
	loop.OnSecTick = func () {
	}
	loop.handlers = [] func() {}
	return loop
}

func (me *tEngineLoop) AddHandler (loopHandler func ()) {
	me.handlers = append(me.handlers, loopHandler)
}

func (me *tEngineLoop) Loop () {
	var onLoopHandler func()
	me.SecTickLast, me.TickNow = me.Time(), me.Time()
	Stats.fps = 0
	if (!me.IsLooping) {
		me.IsLooping = true
		glLogLastError("ngine.PreLoop")
		log.Printf("Enter loop...")
		for me.IsLooping {
			Core.renderLoop()
			Stats.fps++
			Stats.fpsAll++
			me.TickLast = me.TickNow
			me.TickNow = me.Time()
			me.TickDelta = me.TickNow - me.TickLast
			if math.Floor(me.TickNow) != me.SecTickLast {
				runtime.GC()
				if Stats.TrackGC {
					if Stats.GcTime = me.Time() - me.TickNow; Stats.GcTime > Stats.GcMaxTime { Stats.GcMaxTime = Stats.GcTime }
					Stats.gcAll += Stats.GcTime
				}
				Stats.totalSecs++
				Stats.FpsLastSec, Stats.fps, me.SecTickLast = Stats.fps, 0, math.Floor(me.TickNow)
				Core.onSecTick()
				me.OnSecTick()
			}
			for _, onLoopHandler = range me.handlers {
				onLoopHandler()
			}
			UserIO.onLoop()
		}
		glLogLastError("ngine.PostLoop")
	}
}

func (me *tEngineLoop) Time () float64 {
	return glfw.Time()
}

type tEngineStats struct {
	FpsLastSec int
	GcTime, GcMaxTime float64
	TrackGC bool

	fps int
	fpsAll, totalSecs int64
	gcAll float64
}

	func newEngineStats () *tEngineStats {
		return &tEngineStats {}
	}

	func (me *tEngineStats) FpsOverallAverage () int64 {
		return me.fpsAll / me.totalSecs
	}

	func (me *tEngineStats) GcOverallAverage () float64 {
		return me.gcAll / float64(me.totalSecs)
	}
