package assets

import (
	unum "github.com/metaleap/go-util/num"
)

type Transformable interface {
	ChildrenUpdateMatrices ()
	Parent () Transformable
	Transform () *OldTransforms
}

type OldTransforms struct {
	Pos, Rot, Scaling *unum.Vec3

	OnMatrixUpdate func ()
	Owner Transformable
	MatModelView, MatRotX, MatRotY, MatRotZ, MatScaling, MatTranslation *unum.Mat4
}

func (me *OldTransforms) Init (owner Transformable, onMatrixUpdate func ()) {
	me.Owner, me.OnMatrixUpdate = owner, onMatrixUpdate
	me.Pos, me.Rot = &unum.Vec3 {}, &unum.Vec3 {}
	me.Scaling = &unum.Vec3 { 1, 1, 1 }
	me.MatModelView = unum.NewMat4Identity()
	me.MatRotX, me.MatRotY, me.MatRotZ = unum.NewMat4RotationX(me.Rot.X), unum.NewMat4RotationX(me.Rot.Y), unum.NewMat4RotationX(me.Rot.Z)
	me.MatScaling, me.MatTranslation = unum.NewMat4Scaling(me.Scaling), unum.NewMat4Translation(me.Pos)
}

func (me *OldTransforms) OnAnyChanged () {
	me.onScalingChanged()
	me.onPosChanged()
	me.onRotChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) onPosChanged () {
	me.MatTranslation.Translation(me.Pos)
}

func (me *OldTransforms) OnPosChanged () {
	me.onPosChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) OnPosRotChanged () {
	me.onPosChanged()
	me.onRotChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) OnPosScalingChanged () {
	me.onPosChanged()
	me.onScalingChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) onRotChanged () {
	me.onRotXChanged()
	me.onRotYChanged()
	me.onRotZChanged()
}

func (me *OldTransforms) OnRotChanged () {
	me.onRotChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) OnRotScalingChanged () {
	me.onRotChanged()
	me.onScalingChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) onRotXChanged () {
	me.MatRotX.RotationX(me.Rot.X)
}

func (me *OldTransforms) OnRotXChanged () {
	me.onRotXChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) onRotYChanged () {
	me.MatRotY.RotationY(me.Rot.Y)
}

func (me *OldTransforms) OnRotYChanged () {
	me.onRotYChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) onRotZChanged () {
	me.MatRotZ.RotationZ(me.Rot.Z)
}

func (me *OldTransforms) OnRotZChanged () {
	me.onRotZChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) onScalingChanged () {
	me.MatScaling.Scaling(me.Scaling)
}

func (me *OldTransforms) OnScalingChanged () {
	me.onScalingChanged()
	me.UpdateMatrices()
}

func (me *OldTransforms) OnSomeChanged (scaling, pos, rotX, rotY, rotZ bool) {
	if scaling { me.onScalingChanged() }
	if pos { me.onPosChanged() }
	if rotX { me.onRotXChanged() }
	if rotY { me.onRotYChanged() }
	if rotZ { me.onRotZChanged() }
	me.UpdateMatrices()
}

func (me *OldTransforms) SetPos (pos *unum.Vec3) {
	me.Pos = pos
	me.OnPosChanged()
}

func (me *OldTransforms) SetPosX (posX float64) {
	me.Pos.X = posX
	me.OnPosChanged()
}

func (me *OldTransforms) SetPosXYZ (posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
	me.OnPosChanged()
}

func (me *OldTransforms) SetPosY (posY float64) {
	me.Pos.Y = posY
	me.OnPosChanged()
}

func (me *OldTransforms) SetPosZ (posZ float64) {
	me.Pos.Z = posZ
	me.OnPosChanged()
}

func (me *OldTransforms) SetRot (rot *unum.Vec3) {
	me.Rot = rot
	me.OnRotChanged()
}

func (me *OldTransforms) SetRotX (rad float64) {
	me.Rot.X = rad
	me.OnRotXChanged()
}

func (me *OldTransforms) SetRotXYZ (radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
	me.OnRotChanged()
}

func (me *OldTransforms) SetRotY (rad float64) {
	me.Rot.Y = rad
	me.OnRotYChanged()
}

func (me *OldTransforms) SetRotZ (rad float64) {
	me.Rot.Z = rad
	me.OnRotZChanged()
}

func (me *OldTransforms) SetScaling (scaling *unum.Vec3) {
	me.Scaling = scaling
	me.OnScalingChanged()
}

func (me *OldTransforms) SetScalingN (scaling float64) {
	me.Scaling.X, me.Scaling.Y, me.Scaling.Z = scaling, scaling, scaling
	me.OnScalingChanged()
}

func (me *OldTransforms) UpdateMatrices () {
	var mat *unum.Mat4
	var parent = me.Owner.Parent()
	if parent != nil { mat = parent.Transform().MatModelView } else { mat = unum.Mat4Identity }
	me.MatModelView.SetFromMultN(mat, me.MatTranslation, me.MatScaling, me.MatRotX, me.MatRotY, me.MatRotZ)
	me.OnMatrixUpdate()
	me.Owner.ChildrenUpdateMatrices()
}
