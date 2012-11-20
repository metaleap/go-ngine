package core

import (
	glfw "github.com/go-gl/glfw"
)

//	Consider EngineStats a "Singleton" type, only valid use is the Stats global variable.
//	Tracks various go:ngine performance counters over time.
type EngineStats struct {
	//	Gives the total number of frames rendered during the "previous" (not the current) second. Good enough for just a simple-minded FPS indicator.
	FpsLastSec int

	//	This TimingStats instance combines all the individual FrameFoo fields to track over time (both average and maximum) total cost per frame.
	Frame *TimingStats

	//	"Rendering" consists of a CPU-side and a GPU-side cost.
	//	This TimingStats instance combines both to track over time (both average and maximum) total rendering cost per frame.
	FrameRenderBoth *TimingStats

	//	The CPU-side cost of rendering comprises geometry culling, and batching draw calls to the GPU.
	//	This TimingStats instance tracks over time (both average and maximum) CPU-side rendering cost per frame.
	FrameRenderCpu *TimingStats

	//	The GPU-side cost of rendering comprises execution of all draw calls sent by the CPU-side, plus waiting for V-sync if enabled.
	//	This TimingStats instance tracks over time (both average and maximum) GPU-side rendering cost per frame.
	FrameRenderGpu *TimingStats

	//	"Core code" comprises non-rendering go:ngine logic executed every frame.
	//	This TimingStats instance tracks over time (both average and maximum) "core code" cost per frame.
	FrameCoreCode *TimingStats

	//	"User code" comprises user-specific logic executed every frame in your own EngineLoop.OnLoop() callback.
	//	This TimingStats instance tracks over time (both average and maximum) "user code" cost per frame.
	FrameUserCode *TimingStats

	//	During the Loop, the Go Garbge Collector is invoked at least and at most once per second.
	//	Forcing GC "that often" practically guarantees it will almost never have so much work to do as to
	//	noticably block user interaction --- 99.9% of the time it will complete in less than 10ms (most often even under 1ms).
	//	This TimingStats instance over time tracks the maximum and average time spent on that 1x-per-second GC.
	Gc *TimingStats

	fpsCounter int
	fpsAll float64
}

	func newEngineStats () *EngineStats {
		var stats = &EngineStats {}
		stats.reset()
		return stats
	}

	func (me *EngineStats) AverageFps () float64 {
		return me.fpsAll / glfw.Time()
	}

	func (me *EngineStats) reset () {
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
