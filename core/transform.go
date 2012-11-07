package core

import (
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type iTransformable interface {
	transformChildrenUpdateMatrices ()
	transformParent () iTransformable
	transform () *tTransform
}

type tTransform struct {
	Pos, Rot, Scaling *unum.Vec3

	owner iTransformable
	matModelView, matRotX, matRotY, matRotZ, matScaling, matTranslation *unum.Mat4
	matNormal *unum.Mat3
	glMatModelView *ugl.GlMat4
	glMatNormal *ugl.GlMat3
}

func newTransform (owner iTransformable) *tTransform {
	var t = &tTransform {}
	t.owner = owner
	t.Pos, t.Rot = &unum.Vec3 {}, &unum.Vec3 {}
	t.Scaling = &unum.Vec3 { 1, 1, 1 }
	t.matModelView, t.matNormal = unum.NewMat4Identity(), unum.NewMat3Identity()
	t.matRotX, t.matRotY, t.matRotZ = unum.NewMat4RotationX(t.Rot.X), unum.NewMat4RotationX(t.Rot.Y), unum.NewMat4RotationX(t.Rot.Z)
	t.matScaling, t.matTranslation = unum.NewMat4Scaling(t.Scaling), unum.NewMat4Translation(t.Pos)
	t.glMatModelView, t.glMatNormal = ugl.NewGlMat4(t.matModelView), ugl.NewGlMat3(t.matNormal)
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

func (me *tTransform) SetPos (pos *unum.Vec3) {
	me.Pos = pos
	me.OnPosChanged()
}

func (me *tTransform) SetPosX (posX float64) {
	me.Pos.X = posX
	me.OnPosChanged()
}

func (me *tTransform) SetPosXYZ (posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
	me.OnPosChanged()
}

func (me *tTransform) SetPosY (posY float64) {
	me.Pos.Y = posY
	me.OnPosChanged()
}

func (me *tTransform) SetPosZ (posZ float64) {
	me.Pos.Z = posZ
	me.OnPosChanged()
}

func (me *tTransform) SetRot (rot *unum.Vec3) {
	me.Rot = rot
	me.OnRotChanged()
}

func (me *tTransform) SetRotX (rad float64) {
	me.Rot.X = rad
	me.OnRotXChanged()
}

func (me *tTransform) SetRotXYZ (radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
	me.OnRotChanged()
}

func (me *tTransform) SetRotY (rad float64) {
	me.Rot.Y = rad
	me.OnRotYChanged()
}

func (me *tTransform) SetRotZ (rad float64) {
	me.Rot.Z = rad
	me.OnRotZChanged()
}

func (me *tTransform) SetScaling (scaling *unum.Vec3) {
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
	var mat *unum.Mat4
	var parent = me.owner.transformParent()
	if parent != nil { mat = parent.transform().matModelView } else { mat = unum.Mat4Identity }
	me.matModelView.SetFromMultN(mat, me.matTranslation, me.matScaling, me.matRotX, me.matRotY, me.matRotZ)
	me.glMatModelView.Load(me.matModelView)
	me.matModelView.ToInverseMat3(me.matNormal)
	me.matNormal.Transpose()
	me.glMatNormal.Load(me.matNormal)
	me.owner.transformChildrenUpdateMatrices()
}
