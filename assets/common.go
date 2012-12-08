package assets

type FxParamDefs map[string]*FxParamDef

type FxParamDef struct {
	Modifier, Semantic, Sid string
}
