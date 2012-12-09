package assets

type HasAsset struct {
	Asset *Asset
}

type HasExtras struct {
	Extras []*Extra
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

type Asset struct {
	HasExtras
	Created      string
	Modified     string
	Keywords     string
	Revision     string
	Subject      string
	Title        string
	UpAxis       string
	Contributors []*Contributor
	Coverage     *GeographicLocation
}

type Contributor struct {
	Author        string
	AuthorEmail   string
	AuthorWebsite string
	AuthoringTool string
	Comments      string
	Copyright     string
	SourceData    string
}

type Extra struct {
	HasID
	HasName
	HasAsset
	Type       string
	Techniques []*Technique
}

type FxAnnotation struct {
	HasName
	Value interface{}
}

type FxParamDef struct {
	ParamDef
	Annotations        []*FxAnnotation
	Modifier, Semantic string
}

type FxParamDefs map[string]*FxParamDef

type GeographicLocation struct {
	Longitude        float64
	Latitude         float64
	Altitude         float64
	AltitudeAbsolute bool
}

type Input struct {
	Semantic string
	Source   string
}

type InputShared struct {
	Input
	Offset uint64
	Set    uint64
}

type Param struct {
	HasName
	HasSid
	Semantic string
	Type     string
}

type ParamDef struct {
	HasSid
	Value interface{}
}

type ParamDefs map[string]*ParamDef

type ParamInst struct {
	Ref   string
	Value interface{}
}

type Technique struct {
	Profile string
	Data    interface{}
}
