package core

import (
	unum "github.com/metaleap/go-util/num"
)

//	Represents one or more transformations of a Node.
//	This is only used by Node objects, which initialize their NodeTransform with the
//	proper defaults and associate themselves with their NodeTransform. (Any other
//	NodeTransform are invalid.)
//	
//	A single NodeTransform encapsulates an unexported 4x4 matrix that is recalculated
//	from its exported fields via the AddFoo(), SetFoo() or ApplyMatrices() methods.
type NodeTransform struct {
	//	Translation of the from origin.
	Pos unum.Vec3

	//	Rotation for each axis in radians.
	Rot unum.Vec3

	//	Scaling of this node, if any. Defaults to (1, 1, 1) for no scaling.
	Scale unum.Vec3

	//	An arbitrary additional transformation to apply, if any.
	Other *unum.Mat4

	matModelView unum.Mat4
	owner        *Node
}

func (me *NodeTransform) init(owner *Node) {
	me.owner = owner
	me.Scale.X, me.Scale.Y, me.Scale.Z = 1, 1, 1
	me.matModelView.Identity()
	me.ApplyMatrices()
}

//	Adds the specified vector to me.Rot and calls ApplyMatrices().
func (me *NodeTransform) AddRot(rot *unum.Vec3) {
	me.Rot.Add(rot)
	me.ApplyMatrices()
}

//	Adds the specified values to me.Rot and calls ApplyMatrices().
func (me *NodeTransform) AddRotXYZ(x, y, z float64) {
	me.Rot.Add3(x, y, z)
	me.ApplyMatrices()
}

//	Updates the internal 4x4 transformation matrix for all transformations
//	in me. It is only this matrix that is used by the rendering runtime.
func (me *NodeTransform) ApplyMatrices() {
	thrApp.nodeTrans.matScale.Scaling(&me.Scale)
	thrApp.nodeTrans.matTrans.Translation(&me.Pos)
	thrApp.nodeTrans.matRotX.RotationX(&thrApp.numBag, me.Rot.X)
	thrApp.nodeTrans.matRotY.RotationY(&thrApp.numBag, me.Rot.Y)
	thrApp.nodeTrans.matRotZ.RotationZ(&thrApp.numBag, me.Rot.Z)
	if me.owner.parentNode == nil {
		thrApp.nodeTrans.matParent.Identity()
	} else {
		thrApp.nodeTrans.matParent.CopyFrom(&me.owner.parentNode.Transform.matModelView)
	}
	me.matModelView.SetFromMultN(&thrApp.numBag, &thrApp.nodeTrans.matParent, &thrApp.nodeTrans.matTrans, me.Other, &thrApp.nodeTrans.matScale, &thrApp.nodeTrans.matRotX, &thrApp.nodeTrans.matRotY, &thrApp.nodeTrans.matRotZ)
	for _, me.owner.curSubNode = range me.owner.ChildNodes.M {
		me.owner.curSubNode.Transform.ApplyMatrices()
	}
}

//	Sets me.Pos and calls ApplyMatrices().
func (me *NodeTransform) SetPos(pos *unum.Vec3) {
	me.Pos = *pos
	me.ApplyMatrices()
}

//	Sets me.Pos.X and calls ApplyMatrices().
func (me *NodeTransform) SetPosX(posX float64) {
	me.Pos.X = posX
	me.ApplyMatrices()
}

//	Sets me.Pos and calls ApplyMatrices().
func (me *NodeTransform) SetPosXYZ(posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
	me.ApplyMatrices()
}

//	Sets me.Pos.Y and calls ApplyMatrices().
func (me *NodeTransform) SetPosY(posY float64) {
	me.Pos.Y = posY
	me.ApplyMatrices()
}

//	Sets me.Pos.Z and calls ApplyMatrices().
func (me *NodeTransform) SetPosZ(posZ float64) {
	me.Pos.Z = posZ
	me.ApplyMatrices()
}

//	Sets me.Rot and calls ApplyMatrices().
func (me *NodeTransform) SetRot(rot *unum.Vec3) {
	me.Rot = *rot
	me.ApplyMatrices()
}

//	Sets me.Rot.X and calls ApplyMatrices().
func (me *NodeTransform) SetRotX(rad float64) {
	me.Rot.X = rad
	me.ApplyMatrices()
}

//	Sets me.Rot and calls ApplyMatrices().
func (me *NodeTransform) SetRotXYZ(radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
	me.ApplyMatrices()
}

//	Sets me.Rot.Y and calls ApplyMatrices().
func (me *NodeTransform) SetRotY(rad float64) {
	me.Rot.Y = rad
	me.ApplyMatrices()
}

//	Sets me.Rot.Z and calls ApplyMatrices().
func (me *NodeTransform) SetRotZ(rad float64) {
	me.Rot.Z = rad
	me.ApplyMatrices()
}

//	Sets me.Scale and calls ApplyMatrices().
func (me *NodeTransform) SetScale(scale *unum.Vec3) {
	me.Scale = *scale
	me.ApplyMatrices()
}

//	Sets me.Scale and calls ApplyMatrices().
func (me *NodeTransform) SetScaleN(scale float64) {
	me.SetScaleXYZ(scale, scale, scale)
}

//	Sets me.Scale and calls ApplyMatrices().
func (me *NodeTransform) SetScaleXYZ(x, y, z float64) {
	me.Scale.X, me.Scale.Y, me.Scale.Z = x, y, z
	me.ApplyMatrices()
}

//	Returns the result of multiplying deltaPerSecond with EngineLoop.TickDelta.
func (me *NodeTransform) StepDelta(deltaPerSecond float64) float64 {
	return Loop.TickDelta * deltaPerSecond
}
