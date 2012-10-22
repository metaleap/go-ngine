package core

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
