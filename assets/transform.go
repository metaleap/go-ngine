package assets

//	Categorizes the kind of a Transform.
type TransformKind int

const (
	//	A position and orientation transformation suitable for aiming a camera.
	TransformKindLookat TransformKind = iota + 1

	//	A transformation that embodies mathematical changes to points within a coordinate system
	//	or the coordinate system itself.
	TransformKindMatrix

	//	A transformation that specifies how to rotate an object around an axis.
	TransformKindRotate

	//	A transformation that specifies how to deform an object along one axis.
	TransformKindSkew

	//	A transformation that specifies how to change an object's size.
	TransformKindScale

	//	A transformation that changes the position of an object in a local coordinate system.
	TransformKindTranslate
)

//	Represents a single transformation of a specific kind.
type Transform struct {
	//	Sid
	HasSid

	//	The type of this transformation (rotation, skewing, scaling, translation, "look-at", or matrix).
	//	Must be one of the TransformKind* enumerated constants.
	Kind TransformKind

	//	Contains one or more vectors and values representing this transformation.
	//	For TransformKindLookat:
	//		9 values representing three 3D vectors (eye position, interest point, up-axis).
	//	For TransformKindMatrix:
	//		16 values representing one column-major 4x4 matrix.
	//	For TransformKindSkew:
	//		7 values -- one angle in degrees, then two 3D vectors for the axes of rotation and translation.
	//	For TransformKindRotate:
	//		4 values -- one 3D vector specifying the axis of rotation, then one angle in degrees.
	//	For TransformKindTranslate or TransformKindScale:
	//		3 values representing a single 3D vector.
	F []float64
}

func (me *Transform) accessField(fn string) interface{} {
	switch fn {
	case "X":
		return &me.F[0]
	case "Y":
		return &me.F[1]
	case "Z":
		return &me.F[2]
	}
	return nil
}

func (me *Transform) accessIndex(i int) interface{} {
	return &me.F[i]
}
