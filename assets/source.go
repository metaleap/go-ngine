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
	HasTechniques
	Data struct {
		HasID
		HasName
		Bools   []bool
		Floats  []float64
		IdRefs  []string
		Ints    []int64
		Names   []string
		SidRefs []string
		Tokens  []string
	}
	TC struct {
		Accessor *SourceAccessor
	}
}

type Sources map[string]*Source
