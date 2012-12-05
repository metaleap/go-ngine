package assets

type FxNewParams map[string]*FxNewParam

type FxNewParam struct {
	Annotations        map[string]interface{}
	Modifier, Semantic string
}

func NewFxNewParam(modifier, semantic string) (me *FxNewParam) {
	me = &FxNewParam{Modifier: modifier, Semantic: semantic}
	me.Annotations = map[string]interface{}{}
	return
}
