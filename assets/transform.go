package assets

import (
	unum "github.com/metaleap/go-util/num"
)

type Transform interface {
}

type Transforms struct {
	All []Transform
}

	func newTransforms () (me *Transforms) {
		me = &Transforms {}
		return
	}

	func (me *Transforms) Add (transform Transform) {
		me.All = append(me.All, transform)
	}

type TransformLookat struct {
	EyePos, LookTarget, UpAxis unum.Vec3
}

	func NewTransformLookat (eyePos, lookTarget, axisUp *unum.Vec3) (me *TransformLookat) {
		me = &TransformLookat { EyePos: *eyePos, LookTarget: *lookTarget, UpAxis: *axisUp }
		return
	}

type TransformMatrix struct {
	Matrix unum.Mat4
}

	func NewTransformMatrix (matrix *unum.Mat4) (me *TransformMatrix) {
		me = &TransformMatrix { Matrix: *matrix }
		return
	}

type TransformRotate struct {
	Axis unum.Vec3
	Deg float64
}

	func NewTransformRotate (axis *unum.Vec3, deg float64) (me *TransformRotate) {
		me = &TransformRotate { Axis: *axis, Deg: deg }
		return
	}

type TransformScale struct {
	Factor unum.Vec3
}

	func NewTransformScale (factor *unum.Vec3) (me *TransformScale) {
		me = &TransformScale { Factor: *factor }
		return
	}

type TransformSkew struct {
	Deg float64
	Rotation, Translation unum.Vec3
}

	func NewTransformSkew (deg float64, rotation, translation *unum.Vec3) (me *TransformSkew) {
		me = &TransformSkew { Deg: deg, Rotation: *rotation, Translation: *translation }
		return
	}

type TransformTranslate struct {
	Translation unum.Vec3
}

	func NewTransformTranslate (translation *unum.Vec3) (me *TransformTranslate) {
		me = &TransformTranslate { Translation: *translation }
		return
	}
