package assets

import (
	unum "github.com/metaleap/go-util/num"
)

//	Contains two bool values.
type Bool2 [2]bool

//	Contains three bool values.
type Bool3 [3]bool

//	Contains four bool values.
type Bool4 [4]bool

//	Contains two float64 values.
type Float2 [2]float64

//	Contains four float64 values.
type Float2x2 [4]float64

//	Contains six float64 values.
type Float2x3 [6]float64

//	Contains eight float64 values.
type Float2x4 [8]float64

//	Contains three float64 values.
type Float3 [3]float64

//	Contains six float64 values.
type Float3x2 [6]float64

//	Contains nine float64 values.
type Float3x3 [9]float64

//	Contains twelve float64 values.
type Float3x4 [12]float64

//	Contains four float64 values.
type Float4 [4]float64

//	Contains eight float64 values.
type Float4x2 [8]float64

//	Contains twelve float64 values.
type Float4x3 [12]float64

//	Contains sixteen float64 values.
type Float4x4 [16]float64

//	Contains seven float64 values.
type Float7 [7]float64

//	Contains two int64 values.
type Int2 [2]int64

//	Contains four int64 values.
type Int2x2 [4]int64

//	Contains three int64 values.
type Int3 [3]int64

//	Contains nine int64 values.
type Int3x3 [9]int64

//	Contains four int64 values.
type Int4 [4]int64

//	Contains sixteen int64 values.
type Int4x4 [16]int64

//	Provides a bool value.
type ParamBool struct {
	//	The value provided if Ref is empty.
	B bool
	//	If set, refers to a previously defined parameter providing the value.
	Ref RefParam
}

//	Provides a float64 value.
type ParamFloat struct {
	//	The value provided if Ref is empty.
	F float64
	//	If set, refers to a previously defined parameter providing the value.
	Ref RefParam
}

//	Provides two float64 values.
type ParamFloat2 struct {
	//	The values provided if Ref is empty.
	F Float2
	//	If set, refers to a previously defined parameter providing the values.
	Ref RefParam
}

//	Provides a int64 value.
type ParamInt struct {
	//	The value provided if Ref is empty.
	I int64
	//	If set, refers to a previously defined parameter providing the value.
	Ref RefParam
}

//	Provides a scoped float64 value.
type ParamScopedFloat struct {
	//	The value provided if Ref is empty.
	F ScopedFloat
	//	If set, refers to a previously defined parameter providing the value.
	Ref RefParam
}

//	Provides a RefSid value.
type ParamSidRef struct {
	//	The value provided if Ref is empty.
	Sr RefSid
	//	If set, refers to a previously defined parameter providing the value.
	Ref RefParam
}

//	Provides a uint64 value.
type ParamUint struct {
	//	The value provided if Ref is empty.
	U uint64
	//	If set, refers to a previously defined parameter providing the value.
	Ref RefParam
}

//	A bool value that has a scoped identifier.
type ScopedBool struct {
	//	Sid
	HasSid
	//	The value.
	B bool
}

//	A float64 value that has a scoped identifier.
type ScopedFloat struct {
	//	Sid
	HasSid
	//	The value.
	F float64
}

//	Three float64 values that have a scoped identifier.
type ScopedFloat3 struct {
	//	Sid
	HasSid
	//	The values.
	F Float3
}

//	A string value that has a scoped identifier.
type ScopedString struct {
	//	Sid
	HasSid
	//	The value.
	S string
}

//	A 3D vector that has a scoped identifier.
type ScopedVec3 struct {
	//	Sid
	HasSid
	//	The vector.
	V unum.Vec3
}
