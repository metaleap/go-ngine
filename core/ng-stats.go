package core

import (
	glfw "github.com/go-gl/glfw"
)

var (
	//	Tracks various go:ngine performance counters over time.
	Stats EngineStats
)

//	Consider EngineStats a "Singleton" type, only valid use is the core.Stats global variable.
//	Tracks various go:ngine performance indicators over time.
type EngineStats struct {
	//	Gives the total number of frames rendered during the "previous"
	//	(not the current) second. Good enough for just a simple-minded FPS indicator.
	FpsLastSec int

	//	This TimingStats instance combines all the individual FrameFoo fields
	//	to track over time (both average and maximum) total cost per frame.
	Frame TimingStats

	//	"Rendering" consists of a CPU-side and a GPU-side cost.
	//	This TimingStats instance combines both to track over time
	//	(both average and maximum) total rendering cost per frame.
	FrameRenderBoth TimingStats

	//	The CPU-side cost of rendering comprises sending pre-batched
	//	rendering commands (prepared by the "prep" stage) to the GPU.
	//	This TimingStats instance tracks over time (both average
	//	and maximum) CPU-side rendering cost per frame.
	FrameRenderCpu TimingStats

	//	The GPU-side cost of rendering comprises execution of all draw calls
	//	sent by the CPU-side, plus waiting for V-sync if enabled.
	//	This TimingStats instance tracks over time (both average
	//	and maximum) GPU-side rendering cost per frame.
	FrameRenderGpu TimingStats

	//	"Prep code" comprises all go:ngine logic executed every frame in parallel to cull
	//	geometry and prepare a batch of rendering commands for the next (not current) frame.
	//	This TimingStats instance tracks over time (both average and maximum) "prep code" cost per frame.
	FramePrepThread TimingStats

	//	"App code" comprises (mostly user-specific) logic executed every frame in parallel in
	//	your Loop.OnAppThread() callback. Such code may freely modify dynamic Cameras, Nodes etc.
	//	Unlike OnWinThread() code, "app code" always runs in its own thread in parallel to the prep and main threads.
	//	This TimingStats instance tracks over time (both average and maximum) "app code" cost per frame.
	FrameAppThread TimingStats

	//	Tracks over the time (both average and maximum) cost for Loop.On.EverySec() and,
	//	if Diag.LogGLErrorsInLoopOnSec is true, Diag.LogIfGlErr().
	FrameOnSec TimingStats

	//	"Windowing/GPU/IO code" comprises user-specific logic executed every frame via your own
	//	Loop.OnWinThread() callback. This should be kept to a minimum to fully enjoy
	//	the benefits of multi-threading. Main use-cases are calls resulting in GPU state
	//	changes (such as toggling effects in Core.Render.PostFx) and working with UserIO
	//	to poll for user input -- but do consider executing resulting logic in your OnAppThread().
	//	This TimingStats instance tracks over time (both average and maximum) "input code" cost per frame.
	FrameWinThread TimingStats

	//	When CPU-side rendering is completed, Loop waits for the app thread and prep thread
	//	to finish (either before or after GPU-side rendering depending on Loop.SwapLast).
	//	It then moves "prep results" to the render thread and "app results" to the prep thread.
	//	This TimingStats instance tracks over time (both average and maximum) "thread sync" cost per frame.
	FrameThreadSync TimingStats

	//	During the Loop, the Go Garbge Collector is invoked at least and at most once per second.
	//	
	//	Forcing GC "that often" practically guarantees it will almost never have so much work to do as to
	//	noticably block user interaction --- typically well below 10ms, most often around 1ms.
	//	
	//	This TimingStats instance over time tracks the maximum and average time spent on that
	//	1x-per-second-during-Loop GC invokation (but does not track any other GC invokations).
	Gc TimingStats

	enabled    bool
	fpsCounter int
	fpsAll     float64

	Programs struct {
		NumProgsCompiled int
		TotalTimeCost    int64
	}
}

func (_ *EngineStats) addProgCompile(numProgs int, duration int64) {
	Stats.Programs.NumProgsCompiled += numProgs
	Stats.Programs.TotalTimeCost += duration
}

//	Returns the average number of frames-per-second since Loop.Loop() was last called.
func (_ *EngineStats) AverageFps() float64 {
	return Stats.fpsAll / glfw.Time()
}

func (_ *EngineStats) enable() {
	if !Stats.enabled {
		Stats.enabled, Stats.fpsAll = true, 0
	}
}

func (_ *EngineStats) reset() {
	Stats.FpsLastSec, Stats.fpsCounter, Stats.fpsAll, Stats.enabled = 0, 0, 0, false
	Stats.Frame.reset()
	Stats.FrameAppThread.reset()
	Stats.FrameOnSec.reset()
	Stats.FramePrepThread.reset()
	Stats.FrameRenderBoth.reset()
	Stats.FrameRenderCpu.reset()
	Stats.FrameRenderGpu.reset()
	Stats.FrameThreadSync.reset()
	Stats.FrameWinThread.reset()
	Stats.Gc.reset()
}

func (_ *EngineStats) TotalFrames() float64 {
	return Stats.fpsAll
}

//	Helps track average and maximum cost for a variety of performance indicators.
type TimingStats struct {
	max, measuredCounter, measureStartTime, thisTime, totalAccum float64
}

//	Returns the average cost tracked by this performance indicator.
func (me *TimingStats) Average() float64 {
	return me.totalAccum / me.measuredCounter
}

func (me *TimingStats) combine(c1, c2 *TimingStats) {
	if Stats.enabled {
		me.max = c1.max + c2.max
		me.measuredCounter = (c1.measuredCounter + c2.measuredCounter) * 0.5
		me.totalAccum = c1.totalAccum + c2.totalAccum
	}
}

func (me *TimingStats) begin() {
	if Stats.enabled {
		me.measureStartTime = glfw.Time()
	}
}

func (me *TimingStats) end() {
	if Stats.enabled {
		if me.thisTime = glfw.Time() - me.measureStartTime; me.thisTime > me.max {
			me.max = me.thisTime
		}
		me.measuredCounter++
		me.totalAccum += me.thisTime
	}
}

//	Returns the maximum cost tracked by this performance indicator.
func (me *TimingStats) Max() float64 {
	return me.max
}

func (me *TimingStats) reset() {
	*me = TimingStats{}
}
