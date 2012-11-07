package core

import (
	glfw "github.com/go-gl/glfw"
)

type engineStats struct {
	FpsLastSec int
	Frame, FrameRenderBoth, FrameRenderCpu, FrameRenderGpu, FrameCoreCode, FrameUserCode, Gc *TimingStats

	fpsCounter int
	fpsAll float64
}

	func newEngineStats () *engineStats {
		var stats = &engineStats {}
		stats.reset()
		return stats
	}

	func (me *engineStats) AverageFps () float64 {
		return me.fpsAll / glfw.Time()
	}

	func (me *engineStats) reset () {
		var nt = func () *TimingStats { return &TimingStats {} }
		me.FpsLastSec, me.fpsCounter, me.fpsAll = 0, 0, 0
		me.Frame, me.FrameRenderBoth, me.FrameRenderCpu, me.FrameRenderGpu, me.FrameCoreCode, me.FrameUserCode, me.Gc = nt(), nt(), nt(), nt(), nt(), nt(), nt()
	}

type TimingStats struct {
	max, measuredCounter, measureStartTime, thisTime, totalAccum float64
	comb1, comb2 *TimingStats
}

	func (me *TimingStats) Average () float64 {
		return me.totalAccum / me.measuredCounter
	}

	func (me *TimingStats) combine () {
		me.max = me.comb1.max + me.comb2.max
		me.measuredCounter = (me.comb1.measuredCounter + me.comb2.measuredCounter) * 0.5
		me.totalAccum = me.comb1.totalAccum + me.comb2.totalAccum
	}

	func (me *TimingStats) begin () {
		me.measureStartTime = glfw.Time()
	}

	func (me *TimingStats) end () {
		if me.thisTime = glfw.Time() - me.measureStartTime; me.thisTime > me.max { me.max = me.thisTime }
		me.measuredCounter++
		me.totalAccum += me.thisTime
	}

	func (me *TimingStats) Max () float64 {
		return me.max
	}
