package core

import (
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type transformable interface {
	transformChildrenUpdateMatrices ()
	transformParent () transformable
	transform () *transform
}

type transform struct {
	Pos, Rot, Scaling *unum.Vec3

	owner transformable
	matModelView, matRotX, matRotY, matRotZ, matScaling, matTranslation *unum.Mat4
	matNormal *unum.Mat3
	glMatModelView *ugl.GlMat4
	glMatNormal *ugl.GlMat3
}

func newTransform (owner transformable) *transform {
	var t = &transform {}
	t.owner = owner
	t.Pos, t.Rot = &unum.Vec3 {}, &unum.Vec3 {}
	t.Scaling = &unum.Vec3 { 1, 1, 1 }
	t.matModelView, t.matNormal = unum.NewMat4Identity(), unum.NewMat3Identity()
	t.matRotX, t.matRotY, t.matRotZ = unum.NewMat4RotationX(t.Rot.X), unum.NewMat4RotationX(t.Rot.Y), unum.NewMat4RotationX(t.Rot.Z)
	t.matScaling, t.matTranslation = unum.NewMat4Scaling(t.Scaling), unum.NewMat4Translation(t.Pos)
	t.glMatModelView, t.glMatNormal = ugl.NewGlMat4(t.matModelView), ugl.NewGlMat3(t.matNormal)
	return t
}

func (me *transform) OnAnyChanged () {
	me.onScalingChanged()
	me.onPosChanged()
	me.onRotChanged()
	me.updateMatrices()
}

func (me *transform) onPosChanged () {
	me.matTranslation.Translation(me.Pos)
}

func (me *transform) OnPosChanged () {
	me.onPosChanged()
	me.updateMatrices()
}

func (me *transform) OnPosRotChanged () {
	me.onPosChanged()
	me.onRotChanged()
	me.updateMatrices()
}

func (me *transform) OnPosScalingChanged () {
	me.onPosChanged()
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *transform) onRotChanged () {
	me.onRotXChanged()
	me.onRotYChanged()
	me.onRotZChanged()
}

func (me *transform) OnRotChanged () {
	me.onRotChanged()
	me.updateMatrices()
}

func (me *transform) OnRotScalingChanged () {
	me.onRotChanged()
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *transform) onRotXChanged () {
	me.matRotX.RotationX(me.Rot.X)
}

func (me *transform) OnRotXChanged () {
	me.onRotXChanged()
	me.updateMatrices()
}

func (me *transform) onRotYChanged () {
	me.matRotY.RotationY(me.Rot.Y)
}

func (me *transform) OnRotYChanged () {
	me.onRotYChanged()
	me.updateMatrices()
}

func (me *transform) onRotZChanged () {
	me.matRotZ.RotationZ(me.Rot.Z)
}

func (me *transform) OnRotZChanged () {
	me.onRotZChanged()
	me.updateMatrices()
}

func (me *transform) onScalingChanged () {
	me.matScaling.Scaling(me.Scaling)
}

func (me *transform) OnScalingChanged () {
	me.onScalingChanged()
	me.updateMatrices()
}

func (me *transform) OnSomeChanged (scaling, pos, rotX, rotY, rotZ bool) {
	if scaling { me.onScalingChanged() }
	if pos { me.onPosChanged() }
	if rotX { me.onRotXChanged() }
	if rotY { me.onRotYChanged() }
	if rotZ { me.onRotZChanged() }
	me.updateMatrices()
}

func (me *transform) SetPos (pos *unum.Vec3) {
	me.Pos = pos
	me.OnPosChanged()
}

func (me *transform) SetPosX (posX float64) {
	me.Pos.X = posX
	me.OnPosChanged()
}

func (me *transform) SetPosXYZ (posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
	me.OnPosChanged()
}

func (me *transform) SetPosY (posY float64) {
	me.Pos.Y = posY
	me.OnPosChanged()
}

func (me *transform) SetPosZ (posZ float64) {
	me.Pos.Z = posZ
	me.OnPosChanged()
}

func (me *transform) SetRot (rot *unum.Vec3) {
	me.Rot = rot
	me.OnRotChanged()
}

func (me *transform) SetRotX (rad float64) {
	me.Rot.X = rad
	me.OnRotXChanged()
}

func (me *transform) SetRotXYZ (radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
	me.OnRotChanged()
}

func (me *transform) SetRotY (rad float64) {
	me.Rot.Y = rad
	me.OnRotYChanged()
}

func (me *transform) SetRotZ (rad float64) {
	me.Rot.Z = rad
	me.OnRotZChanged()
}

func (me *transform) SetScaling (scaling *unum.Vec3) {
	me.Scaling = scaling
	me.OnScalingChanged()
}

func (me *transform) SetScalingN (scaling float64) {
	me.Scaling.X, me.Scaling.Y, me.Scaling.Z = scaling, scaling, scaling
	me.OnScalingChanged()
}

func (me *transform) StepDelta (deltaPerSecond float64) float64 {
	return Loop.TickDelta * deltaPerSecond
}

func (me *transform) updateMatrices () {
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
