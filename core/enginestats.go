package core

type tEngineStats struct {
	FpsLastSec int
	GcTime, GcMaxTime float64
	TrackGC bool
	fps, fpsAll, fpsSecs int
}

func (me *tEngineStats) FpsOverallAverage () int {
	return me.fpsAll / me.fpsSecs
}
