package assets

type SourceAccessor struct {
	Count  uint64
	Offset uint64
	Source string
	Stride uint64
	Params []*Param
}

type Source struct {
	BaseDef
	Data struct {
		Bools   *ListBools
		Floats  *ListFloats
		IdRefs  *ListStrings
		Ints    *ListInts
		Names   *ListStrings
		SidRefs *ListStrings
		Tokens  *ListStrings
	}
	TechniqueCommon struct {
		Accessor *SourceAccessor
	}
	Techniques []*Technique
}

type Sources map[string]*Source
