package assets

import (
	unum "github.com/metaleap/go-util/num"
)

type Bool2 [2]bool

type Bool3 [3]bool

type Bool4 [4]bool

type Float2 [2]float64

type Float2x2 [4]float64

type Float2x3 [6]float64

type Float2x4 [8]float64

type Float3 [3]float64

type Float3x2 [6]float64

type Float3x3 [9]float64

type Float3x4 [12]float64

type Float4 [4]float64

type Float4x2 [8]float64

type Float4x3 [12]float64

type Float4x4 [16]float64

type Float7 [7]float64

type Int2 [2]int64

type Int2x2 [4]int64

type Int3 [3]int64

type Int3x3 [9]int64

type Int4 [4]int64

type Int4x4 [16]int64

/*
type ListBools struct {
	HasID
	HasName
	B []bool
}

type ListFloats struct {
	HasID
	HasName
	F []float64
}

type ListInts struct {
	HasID
	HasName
	I []int64
}

type ListStrings struct {
	HasID
	HasName
	S []string
}
*/

type ParamBool struct {
	B   bool
	Ref string
}

type ParamFloat struct {
	F   float64
	Ref string
}

type ParamFloat2 struct {
	F   Float2
	Ref string
}

type ParamInt struct {
	I   int64
	Ref string
}

type ParamScopedFloat struct {
	F   ScopedFloat
	Ref string
}

type ParamSidRef struct {
	S   string
	Ref string
}

type ParamUint struct {
	U   uint64
	Ref string
}

type ScopedBool struct {
	HasSid
	B bool
}

type ScopedFloat struct {
	HasSid
	F float64
}

type ScopedFloat3 struct {
	HasSid
	F Float3
}

type ScopedString struct {
	HasSid
	S string
}

type ScopedVec3 struct {
	HasSid
	V unum.Vec3
}
