package core

import (
	glutil "github.com/metaleap/go-util/gl"
	numutil "github.com/metaleap/go-util/num"
)

type ITransformable interface {
	transformChildrenUpdateMatrices ()
	transformParent () ITransformable
	transform () *TTransform
}

type TTransform struct {
	Pos, Rot, Scaling *numutil.TVec3

	owner ITransformable
	matModelView, matRotX, matRotY, matRotZ, matScaling, matTranslation *numutil.TMat4
	matNormal *numutil.TMat3
	glMatModelView *glutil.TGlMat4
	glMatNormal *glutil.TGlMat3
}

func NewTransform (owner ITransformable) *TTransform {
	var t = &TTransform {}
	t.owner = owner
	t.Pos, t.Rot = &numutil.TVec3 {}, &numutil.TVec3 {}
	t.Scaling = &numutil.TVec3 { 1, 1, 1 }
	t.matModelView, t.matNormal = numutil.NewMat4Identity(), numutil.NewMat3Identity()
	t.matRotX, t.matRotY, t.matRotZ = numutil.NewMat4RotationX(t.Rot.X), numutil.NewMat4RotationX(t.Rot.Y), numutil.NewMat4RotationX(t.Rot.Z)
	t.matScaling, t.matTranslation = numutil.NewMat4Scaling(t.Scaling), numutil.NewMat4Translation(t.Pos)
	t.glMatModelView, t.glMatNormal = glutil.NewGlMat4(t.matModelView), glutil.NewGlMat3(t.matNormal)
	return t
}

func (me *TTransform) OnAnyChanged () {
	me.onScalingChanged()
	me.onPosChanged()
	me.onRotChanged()
	me.updateMatrices()
}

func (me *TTransform) onPosChanged () {
	me.matTranslation.Translation(me.Pos)
}

func (me *TTransform) OnPosChanged () {
	me.onPosChanged()
	me.updateMatrices()
}

func (me *TTransform) OnPosRotChanged () {
	me.onPosChanged()
	me.onRotChanged()
	me.updateMatrices()
}

func (me *TTransform) OnPosScalingChanged () {
	me.onPosChanged()
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *TTransform) onRotChanged () {
	me.onRotXChanged()
	me.onRotYChanged()
	me.onRotZChanged()
}

func (me *TTransform) OnRotChanged () {
	me.onRotChanged()
	me.updateMatrices()
}

func (me *TTransform) OnRotScalingChanged () {
	me.onRotChanged()
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *TTransform) onRotXChanged () {
	me.matRotX.RotationX(me.Rot.X)
}

func (me *TTransform) OnRotXChanged () {
	me.onRotXChanged()
	me.updateMatrices()
}

func (me *TTransform) onRotYChanged () {
	me.matRotY.RotationY(me.Rot.Y)
}

func (me *TTransform) OnRotYChanged () {
	me.onRotYChanged()
	me.updateMatrices()
}

func (me *TTransform) onRotZChanged () {
	me.matRotZ.RotationZ(me.Rot.Z)
}

func (me *TTransform) OnRotZChanged () {
	me.onRotZChanged()
	me.updateMatrices()
}

func (me *TTransform) onScalingChanged () {
	me.matScaling.Scaling(me.Scaling)
}

func (me *TTransform) OnScalingChanged () {
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *TTransform) OnSomeChanged (scaling, pos, rotX, rotY, rotZ bool) {
	if scaling { me.onScalingChanged() }
	if pos { me.onPosChanged() }
	if rotX { me.onRotXChanged() }
	if rotY { me.onRotYChanged() }
	if rotZ { me.onRotZChanged() }
	me.updateMatrices()
}

func (me *TTransform) SetPos (pos *numutil.TVec3) {
	me.Pos = pos
	me.OnPosChanged()
}

func (me *TTransform) SetRot (rot *numutil.TVec3) {
	me.Rot = rot
	me.OnRotChanged()
}

func (me *TTransform) SetRotX (rad float64) {
	me.Rot.X = rad
	me.OnRotXChanged()
}

func (me *TTransform) SetRotY (rad float64) {
	me.Rot.Y = rad
	me.OnRotYChanged()
}

func (me *TTransform) SetRotZ (rad float64) {
	me.Rot.Z = rad
	me.OnRotZChanged()
}

func (me *TTransform) SetScaling (scaling *numutil.TVec3) {
	me.Scaling = scaling
	me.OnScalingChanged()
}

func (me *TTransform) SetScalingN (scaling float64) {
	me.Scaling.X, me.Scaling.Y, me.Scaling.Z = scaling, scaling, scaling
	me.OnScalingChanged()
}

func (me *TTransform) StepDelta (deltaPerSecond float64) float64 {
	return Core.Timer.TickDelta * deltaPerSecond
}

func (me *TTransform) updateMatrices () {
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
