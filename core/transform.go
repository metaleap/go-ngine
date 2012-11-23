package core

import (
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type Transformable interface {
	ChildrenUpdateMatrices()
	Parent() Transformable
	Transforms() *Transforms
}

type Transforms struct {
	Pos, Rot, Scaling                                                   *unum.Vec3
	Owner                                                               Transformable
	MatModelView, MatRotX, MatRotY, MatRotZ, MatScaling, MatTranslation *unum.Mat4

	glMatModelView *ugl.GlMat4
	glMatNormal    *ugl.GlMat3
}

func newTransforms(owner Transformable) (me *Transforms) {
	me = &Transforms{}
	// me.Transforms.Init(owner, func () {  })
	me.Owner = owner
	me.Pos, me.Rot = &unum.Vec3{}, &unum.Vec3{}
	me.Scaling = &unum.Vec3{1, 1, 1}
	me.MatModelView = unum.NewMat4Identity()
	me.MatRotX, me.MatRotY, me.MatRotZ = unum.NewMat4RotationX(me.Rot.X), unum.NewMat4RotationX(me.Rot.Y), unum.NewMat4RotationX(me.Rot.Z)
	me.MatScaling, me.MatTranslation = unum.NewMat4Scaling(me.Scaling), unum.NewMat4Translation(me.Pos)
	me.glMatModelView = ugl.NewGlMat4(me.MatModelView)
	return
}

func (me *Transforms) OnAnyChanged() {
	me.onScalingChanged()
	me.onPosChanged()
	me.onRotChanged()
	me.UpdateMatrices()
}

func (me *Transforms) onPosChanged() {
	me.MatTranslation.Translation(me.Pos)
}

func (me *Transforms) OnPosChanged() {
	me.onPosChanged()
	me.UpdateMatrices()
}

func (me *Transforms) OnPosRotChanged() {
	me.onPosChanged()
	me.onRotChanged()
	me.UpdateMatrices()
}

func (me *Transforms) OnPosScalingChanged() {
	me.onPosChanged()
	me.onScalingChanged()
	me.UpdateMatrices()
}

func (me *Transforms) onRotChanged() {
	me.onRotXChanged()
	me.onRotYChanged()
	me.onRotZChanged()
}

func (me *Transforms) OnRotChanged() {
	me.onRotChanged()
	me.UpdateMatrices()
}

func (me *Transforms) OnRotScalingChanged() {
	me.onRotChanged()
	me.onScalingChanged()
	me.UpdateMatrices()
}

func (me *Transforms) onRotXChanged() {
	me.MatRotX.RotationX(me.Rot.X)
}

func (me *Transforms) OnRotXChanged() {
	me.onRotXChanged()
	me.UpdateMatrices()
}

func (me *Transforms) onRotYChanged() {
	me.MatRotY.RotationY(me.Rot.Y)
}

func (me *Transforms) OnRotYChanged() {
	me.onRotYChanged()
	me.UpdateMatrices()
}

func (me *Transforms) onRotZChanged() {
	me.MatRotZ.RotationZ(me.Rot.Z)
}

func (me *Transforms) OnRotZChanged() {
	me.onRotZChanged()
	me.UpdateMatrices()
}

func (me *Transforms) onScalingChanged() {
	me.MatScaling.Scaling(me.Scaling)
}

func (me *Transforms) OnScalingChanged() {
	me.onScalingChanged()
	me.UpdateMatrices()
}

func (me *Transforms) OnSomeChanged(scaling, pos, rotX, rotY, rotZ bool) {
	if scaling {
		me.onScalingChanged()
	}
	if pos {
		me.onPosChanged()
	}
	if rotX {
		me.onRotXChanged()
	}
	if rotY {
		me.onRotYChanged()
	}
	if rotZ {
		me.onRotZChanged()
	}
	me.UpdateMatrices()
}

func (me *Transforms) SetPos(pos *unum.Vec3) {
	me.Pos = pos
	me.OnPosChanged()
}

func (me *Transforms) SetPosX(posX float64) {
	me.Pos.X = posX
	me.OnPosChanged()
}

func (me *Transforms) SetPosXYZ(posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
	me.OnPosChanged()
}

func (me *Transforms) SetPosY(posY float64) {
	me.Pos.Y = posY
	me.OnPosChanged()
}

func (me *Transforms) SetPosZ(posZ float64) {
	me.Pos.Z = posZ
	me.OnPosChanged()
}

func (me *Transforms) SetRot(rot *unum.Vec3) {
	me.Rot = rot
	me.OnRotChanged()
}

func (me *Transforms) SetRotX(rad float64) {
	me.Rot.X = rad
	me.OnRotXChanged()
}

func (me *Transforms) SetRotXYZ(radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
	me.OnRotChanged()
}

func (me *Transforms) SetRotY(rad float64) {
	me.Rot.Y = rad
	me.OnRotYChanged()
}

func (me *Transforms) SetRotZ(rad float64) {
	me.Rot.Z = rad
	me.OnRotZChanged()
}

func (me *Transforms) SetScaling(scaling *unum.Vec3) {
	me.Scaling = scaling
	me.OnScalingChanged()
}

func (me *Transforms) SetScalingN(scaling float64) {
	me.Scaling.X, me.Scaling.Y, me.Scaling.Z = scaling, scaling, scaling
	me.OnScalingChanged()
}

func (me *Transforms) StepDelta(deltaPerSecond float64) float64 {
	return Loop.TickDelta * deltaPerSecond
}

func (me *Transforms) UpdateMatrices() {
	var mat *unum.Mat4
	var parent = me.Owner.Parent()
	if parent != nil {
		mat = parent.Transforms().MatModelView
	} else {
		mat = unum.Mat4Identity
	}
	me.MatModelView.SetFromMultN(mat, me.MatTranslation, me.MatScaling, me.MatRotX, me.MatRotY, me.MatRotZ)
	me.glMatModelView.Load(me.MatModelView)
	me.Owner.ChildrenUpdateMatrices()
}
