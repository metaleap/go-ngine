package core

import (
	glfw "github.com/go-gl/glfw"
)

type tEngineStats struct {
	FpsLastSec int
	Frame, FrameRender, FrameSwap, FrameCoreCode, FrameUserCode, Gc *TEngineStatsTiming

	fpsCounter int
	fpsAll float64
	isWarmup bool
}

	func newEngineStats () *tEngineStats {
		var stats = &tEngineStats {}
		stats.reset()
		return stats
	}

	func (me *tEngineStats) AverageFps () float64 {
		return me.fpsAll / glfw.Time()
	}

	func (me *tEngineStats) reset () {
		var nt = func () *TEngineStatsTiming { return &TEngineStatsTiming {} }
		me.isWarmup, me.FpsLastSec, me.fpsCounter, me.fpsAll = true, 0, 0, 0
		me.Frame, me.FrameRender, me.FrameSwap, me.FrameCoreCode, me.FrameUserCode, me.Gc = nt(), nt(), nt(), nt(), nt(), nt()
	}

type TEngineStatsTiming struct {
	max, measuredCounter, measureStartTime, thisTime, totalAccum float64
}

	func (me *TEngineStatsTiming) Average () float64 {
		return me.totalAccum / me.measuredCounter
	}

	func (me *TEngineStatsTiming) Max () float64 {
		return me.max
	}

	func (me *TEngineStatsTiming) begin () {
		me.measureStartTime = glfw.Time()
	}

	func (me *TEngineStatsTiming) end () {
		if me.thisTime = glfw.Time() - me.measureStartTime; me.thisTime > me.max { me.max = me.thisTime }
		me.measuredCounter++
		me.totalAccum += me.thisTime
	}
