package core

import (
	glutil "github.com/go3d/go-util/gl"
	numutil "github.com/go3d/go-util/num"
)

type iTransformable interface {
	transformChildrenUpdateMatrices ()
	transformParent () iTransformable
	transform () *tTransform
}

type tTransform struct {
	Pos, Rot, Scaling *numutil.TVec3

	owner iTransformable
	matModelView, matRotX, matRotY, matRotZ, matScaling, matTranslation *numutil.TMat4
	matNormal *numutil.TMat3
	glMatModelView *glutil.TGlMat4
	glMatNormal *glutil.TGlMat3
}

func newTransform (owner iTransformable) *tTransform {
	var t = &tTransform {}
	t.owner = owner
	t.Pos, t.Rot = &numutil.TVec3 {}, &numutil.TVec3 {}
	t.Scaling = &numutil.TVec3 { 1, 1, 1 }
	t.matModelView, t.matNormal = numutil.NewMat4Identity(), numutil.NewMat3Identity()
	t.matRotX, t.matRotY, t.matRotZ = numutil.NewMat4RotationX(t.Rot.X), numutil.NewMat4RotationX(t.Rot.Y), numutil.NewMat4RotationX(t.Rot.Z)
	t.matScaling, t.matTranslation = numutil.NewMat4Scaling(t.Scaling), numutil.NewMat4Translation(t.Pos)
	t.glMatModelView, t.glMatNormal = glutil.NewGlMat4(t.matModelView), glutil.NewGlMat3(t.matNormal)
	return t
}

func (me *tTransform) OnAnyChanged () {
	me.onScalingChanged()
	me.onPosChanged()
	me.onRotChanged()
	me.updateMatrices()
}

func (me *tTransform) onPosChanged () {
	me.matTranslation.Translation(me.Pos)
}

func (me *tTransform) OnPosChanged () {
	me.onPosChanged()
	me.updateMatrices()
}

func (me *tTransform) OnPosRotChanged () {
	me.onPosChanged()
	me.onRotChanged()
	me.updateMatrices()
}

func (me *tTransform) OnPosScalingChanged () {
	me.onPosChanged()
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *tTransform) onRotChanged () {
	me.onRotXChanged()
	me.onRotYChanged()
	me.onRotZChanged()
}

func (me *tTransform) OnRotChanged () {
	me.onRotChanged()
	me.updateMatrices()
}

func (me *tTransform) OnRotScalingChanged () {
	me.onRotChanged()
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *tTransform) onRotXChanged () {
	me.matRotX.RotationX(me.Rot.X)
}

func (me *tTransform) OnRotXChanged () {
	me.onRotXChanged()
	me.updateMatrices()
}

func (me *tTransform) onRotYChanged () {
	me.matRotY.RotationY(me.Rot.Y)
}

func (me *tTransform) OnRotYChanged () {
	me.onRotYChanged()
	me.updateMatrices()
}

func (me *tTransform) onRotZChanged () {
	me.matRotZ.RotationZ(me.Rot.Z)
}

func (me *tTransform) OnRotZChanged () {
	me.onRotZChanged()
	me.updateMatrices()
}

func (me *tTransform) onScalingChanged () {
	me.matScaling.Scaling(me.Scaling)
}

func (me *tTransform) OnScalingChanged () {
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *tTransform) OnSomeChanged (scaling, pos, rotX, rotY, rotZ bool) {
	if scaling { me.onScalingChanged() }
	if pos { me.onPosChanged() }
	if rotX { me.onRotXChanged() }
	if rotY { me.onRotYChanged() }
	if rotZ { me.onRotZChanged() }
	me.updateMatrices()
}

func (me *tTransform) SetPos (pos *numutil.TVec3) {
	me.Pos = pos
	me.OnPosChanged()
}

func (me *tTransform) SetRot (rot *numutil.TVec3) {
	me.Rot = rot
	me.OnRotChanged()
}

func (me *tTransform) SetRotX (rad float64) {
	me.Rot.X = rad
	me.OnRotXChanged()
}

func (me *tTransform) SetRotY (rad float64) {
	me.Rot.Y = rad
	me.OnRotYChanged()
}

func (me *tTransform) SetRotZ (rad float64) {
	me.Rot.Z = rad
	me.OnRotZChanged()
}

func (me *tTransform) SetScaling (scaling *numutil.TVec3) {
	me.Scaling = scaling
	me.OnScalingChanged()
}

func (me *tTransform) SetScalingN (scaling float64) {
	me.Scaling.X, me.Scaling.Y, me.Scaling.Z = scaling, scaling, scaling
	me.OnScalingChanged()
}

func (me *tTransform) StepDelta (deltaPerSecond float64) float64 {
	return Loop.TickDelta * deltaPerSecond
}

func (me *tTransform) updateMatrices () {
	var mat *numutil.TMat4
	var parent = me.owner.transformParent()
	if parent != nil { mat = parent.transform().matModelView } else { mat = numutil.Mat4Identity }
	me.matModelView.SetFromMultN(mat, me.matTranslation, me.matScaling, me.matRotX, me.matRotY, me.matRotZ)
	me.glMatModelView.Load(me.matModelView)
	me.matModelView.ToInverseMat3(me.matNormal)
	me.matNormal.Transpose()
	me.glMatNormal.Load(me.matNormal)
	me.owner.transformChildrenUpdateMatrices()
}
