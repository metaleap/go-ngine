package core

type tCameraOptions struct {
	BackfaceCulling bool
	cam *TCamera
}

func newCameraOptions (cam *TCamera) *tCameraOptions {
	var opt = &tCameraOptions {}
	opt.cam = cam
	return opt
}

func (me *tCameraOptions) ToggleGlBackfaceCulling () {
	me.BackfaceCulling = !me.BackfaceCulling
}
