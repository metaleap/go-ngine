package assets

type ImageInst struct {
	baseInst
	Def *ImageDef
}

func newImageInst (def *ImageDef, id string) (me *ImageInst) {
	me = &ImageInst { Def: def }
	me.base.init(id)
	return
}
