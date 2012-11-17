package assets

type CameraDef struct {
	baseDef
	FovOrMagX, FovOrMagY, ZfarPlane, ZnearPlane float64
	IsOrtho bool
}

func newCameraDef (id string) (me *CameraDef) {
	me = &CameraDef {}
	me.base.init(id)
	return
}

func (me *CameraDef) NewInst (id string) *CameraInst {
	return newCameraInst(me, id)
}
