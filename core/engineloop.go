package core

import (
	"log"
	"math"
	"runtime"
)

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
	me.SecTickLast, me.TickNow = Windowing.Time(), Windowing.Time()
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
			me.TickNow = Windowing.Time()
			me.TickDelta = me.TickNow - me.TickLast
			if math.Floor(me.TickNow) != me.SecTickLast {
				runtime.GC()
				if Stats.TrackGC {
					if Stats.GcTime = Windowing.Time() - me.TickNow; Stats.GcTime > Stats.GcMaxTime { Stats.GcMaxTime = Stats.GcTime }
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
			Windowing.OnLoop()
		}
		glLogLastError("ngine.PostLoop")
	}
}
