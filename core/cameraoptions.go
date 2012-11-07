package core

type cameraOptions struct {
	BackfaceCulling bool
	cam *Camera
}

func newCameraOptions (cam *Camera) *cameraOptions {
	var opt = &cameraOptions {}
	opt.cam = cam
	return opt
}

func (me *cameraOptions) ToggleGlBackfaceCulling () {
	me.BackfaceCulling = !me.BackfaceCulling
}
