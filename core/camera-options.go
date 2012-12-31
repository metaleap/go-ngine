package core

type CameraOptions struct {
	BackfaceCulling   bool
	FovY, ZFar, ZNear float64
}

func newCameraOptions() (me *CameraOptions) {
	me = &CameraOptions{}
	me.FovY = 37.8493
	me.ZFar = 30000
	me.ZNear = 0.3
	return
}

func (me *CameraOptions) ToggleGlBackfaceCulling() {
	me.BackfaceCulling = !me.BackfaceCulling
}
