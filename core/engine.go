package core

var (
	Loop *tEngineLoop
	Core *tEngineCore
	Stats *tEngineStats
	Windowing = newWindowing()
)

func Dispose () {
	glLogLastError("ngine.Pre_Dispose")
	if Core != nil { Core.Dispose() }
	glLogLastError("ngine.Core_Dispose")
	glDispose()
	Windowing.dispose()
	Core, Loop, Stats = nil, nil, nil
}

func Init (options *tOptions, winTitle string) error {
	var err error
	var isVerErr bool
	var forceContext = false
	tryInit:
	if err = Windowing.init(options, winTitle, forceContext); err == nil {
		if err, isVerErr = glInit(); err == nil {
			Loop, Stats, Core = newEngineLoop(), newEngineStats(), newEngineCore(options)
		} else if isVerErr && !forceContext {
			forceContext = true
			Windowing.isGlfwInit, Windowing.isGlfwWindow = false, false
			goto tryInit
		}
	}
	return err
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
