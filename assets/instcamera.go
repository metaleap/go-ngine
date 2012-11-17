package assets

type CameraInst struct {
	baseInst
	Def *CameraDef
}

func newCameraInst (def *CameraDef, id string) (me *CameraInst) {
	me = &CameraInst { Def: def }
	me.base.init(id)
	return
}
