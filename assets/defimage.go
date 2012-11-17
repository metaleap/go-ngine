package assets

type ImageDef struct {
	baseDef
}

func newImageDef (id string) (me *ImageDef) {
	me = &ImageDef {}
	me.base.init(id)
	return
}

func (me *ImageDef) NewInst (id string) *ImageInst {
	return newImageInst(me, id)
}
