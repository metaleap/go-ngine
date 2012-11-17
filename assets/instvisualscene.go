package assets

type VisualSceneInst struct {
	baseInst
	Def *VisualSceneDef
}

func newVisualSceneInst (def *VisualSceneDef, id string) (me *VisualSceneInst) {
	me = &VisualSceneInst { Def: def }
	me.base.init(id)
	return
}
