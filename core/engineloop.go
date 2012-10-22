package core

import (
	"log"
	"math"
	"runtime"

	nglcore "github.com/go3d/go-ngine/core/glcore"
)

type tEngineLoop struct {
	IsLooping bool
	SecTickLast, TickNow, TickLast, TickDelta float64
	OnLoopHandlers []func()
	OnSecTick func()
}

func newEngineLoop (onSecTick func()) *tEngineLoop {
	var loop = &tEngineLoop {}
	if loop.OnSecTick = onSecTick; onSecTick == nil {
		loop.OnSecTick = func () {
		}
	}
	loop.OnLoopHandlers = [] func() {}
	return loop
}

func (me *tEngineLoop) Loop () {
	var onLoopHandler func()
	me.SecTickLast, me.TickNow = Windowing.Time(), Windowing.Time()
	Stats.fps = 0
	if (!me.IsLooping) {
		me.IsLooping = true
		nglcore.LogLastError("ngine.PreLoop")
		log.Printf("Enter loop...")
		for me.IsLooping {
			Core.RenderLoop()
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
			for _, onLoopHandler = range me.OnLoopHandlers {
				onLoopHandler()
			}
			Windowing.OnLoop()
		}
		nglcore.LogLastError("ngine.PostLoop")
	}
}
