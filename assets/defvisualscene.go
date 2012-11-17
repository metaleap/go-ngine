package assets

type VisualSceneDef struct {
	baseDef
}

func newVisualSceneDef (id string) (me *VisualSceneDef) {
	me = &VisualSceneDef {}
	me.base.init(id)
	return
}

func (me *VisualSceneDef) NewInst (id string) *VisualSceneInst {
	return newVisualSceneInst(me, id)
}
