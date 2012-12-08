package assets

type FxParamDefs map[string]*FxParamDef

type FxParamDef struct {
	HasSid
	Modifier, Semantic string
	Value              interface{}
}

type HasID struct {
	//	The unique identifier of this *Def*, *Inst* or *Lib*.
	ID string
}

type HasName struct {
	//	The optional pretty-print name/title of this *Def*, *Inst* or *Lib*.
	Name string
}

type HasSid struct {
	//	Scoped ID
	Sid string
}
