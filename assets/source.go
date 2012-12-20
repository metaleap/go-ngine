package assets

//	Declares a data repository that provides values
//	according to the semantics of an Input that refers to it.
type Source struct {
	//	Id
	HasId
	//	Name
	HasName
	//	Asset
	HasAsset
	//	Techniques
	HasTechniques
	//	The data array of this Source.
	Array SourceArray
	//	Common-technique profile
	TC struct {
		//	Describes a stream of values from this array data source.
		Accessor *SourceAccessor
	}
}

//	Describes a stream of values from an array data source.
type SourceAccessor struct {
	//	The number of times the array is accessed. Required.
	Count uint64
	//	The index of the first value to be read from the array. The default is 0. Optional.
	Offset uint64
	//	The Id of the array to access. Required.
	Source RefId
	//	The number of values that are to be considered a unit during each access to the array.
	//	The default is 1, indicating that a single value is accessed. Optional.
	Stride uint64
	//	The number and order of Params define the output of the accessor.
	//	Parameters are bound to values in the order in which both are specified. No reordering of the data can occur.
	//	A Param without a Name indicates that the value is not part of the output, so the Param is unbound.
	Params []*Param
}

//	Allocates, initializes and returns a new SourceAccessor.
func NewSourceAccessor() (me *SourceAccessor) {
	me = &SourceAccessor{Stride: 1}
	return
}

//	The data array of a Source.
//	Of all its []slice fields, only ONE should ever be non-nil/non-empty at any time.
type SourceArray struct {
	//	Id
	HasId
	//	Name
	HasName
	//	A slice into the array of bools that this Source represents, if any.
	Bools []bool
	//	A slice into the array of floats that this Source represents, if any.
	Floats []float64
	//	A slice into the array of RefIds that this Source represents, if any.
	IdRefs []string
	//	A slice into the array of ints that this Source represents, if any.
	Ints []int64
	//	A slice into the array of names that this Source represents, if any.
	Names []string
	//	A slice into the array of RefSids that this Source represents, if any.
	SidRefs []string
	//	A slice into the array of tokens that this Source represents, if any.
	Tokens []string
}

//	A hash-table of Sources, each keyed with its Id.
type Sources map[string]*Source
