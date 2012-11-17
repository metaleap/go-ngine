package core

type CameraOptions struct {
	BackfaceCulling bool
}

func newCameraOptions () (me *CameraOptions) {
	me = &CameraOptions {}
	return
}

func (me *CameraOptions) ToggleGlBackfaceCulling () {
	me.BackfaceCulling = !me.BackfaceCulling
}
