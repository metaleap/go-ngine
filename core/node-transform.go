package core

import (
	unum "github.com/metaleap/go-util/num"
)

//	Represents one or more transformations of a Node.
//	This is only used by Node objects, which initialize their NodeTransforms with the
//	proper defaults and associate themselves with their NodeTransforms. (Any other
//	NodeTransforms are invalid.)
//	
//	A single NodeTransforms encapsulates an unexported 4x4 matrix that is recalculated
//	from its exported fields via the SetFoo() or ApplyMatrices() methods.
type NodeTransforms struct {
	//	Translation of the from origin, if any.
	Pos unum.Vec3

	//	Rotation for each axis in radians, if any.
	Rot unum.Vec3

	//	Scaling of this node, if any. Defaults to (1, 1, 1) for no scaling.
	Scale unum.Vec3

	//	An arbitrary additional transformation to apply. Defaults to the 4x4 identity matrix for 'none'.
	Other *unum.Mat4

	matModelView unum.Mat4
	owner        *Node
}

func (me *NodeTransforms) init(owner *Node) {
	me.owner = owner
	me.Scale.X, me.Scale.Y, me.Scale.Z = 1, 1, 1
	me.matModelView.Identity()
	me.ApplyMatrices()
}

//	Adds the specified vector to me.Rot and calls ApplyMatrices().
func (me *NodeTransforms) AddRot(rot *unum.Vec3) {
	if rot != nil {
		me.Rot.Add(rot)
		me.ApplyMatrices()
	}
}

//	Adds the specified values to me.Rot and calls ApplyMatrices().
func (me *NodeTransforms) AddRotXYZ(x, y, z float64) {
	me.Rot.Add3(x, y, z)
	me.ApplyMatrices()
}

//	Updates the internal 4x4 transformation matrix for all transformations
//	in me. It is only this matrix that is used by the rendering runtime.
func (me *NodeTransforms) ApplyMatrices() {
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
func (me *NodeTransforms) SetPos(pos *unum.Vec3) {
	if pos != nil {
		me.Pos = *pos
		me.ApplyMatrices()
	} else {
		me.SetPosXYZ(0, 0, 0)
	}
}

//	Sets me.Pos.X and calls ApplyMatrices().
func (me *NodeTransforms) SetPosX(posX float64) {
	me.Pos.X = posX
	me.ApplyMatrices()
}

//	Sets me.Pos and calls ApplyMatrices().
func (me *NodeTransforms) SetPosXYZ(posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
	me.ApplyMatrices()
}

//	Sets me.Pos.Y and calls ApplyMatrices().
func (me *NodeTransforms) SetPosY(posY float64) {
	me.Pos.Y = posY
	me.ApplyMatrices()
}

//	Sets me.Pos.Z and calls ApplyMatrices().
func (me *NodeTransforms) SetPosZ(posZ float64) {
	me.Pos.Z = posZ
	me.ApplyMatrices()
}

//	Sets me.Rot and calls ApplyMatrices().
func (me *NodeTransforms) SetRot(rot *unum.Vec3) {
	if rot != nil {
		me.Rot = *rot
		me.ApplyMatrices()
	} else {
		me.SetRotXYZ(0, 0, 0)
	}
}

//	Sets me.Rot.X and calls ApplyMatrices().
func (me *NodeTransforms) SetRotX(rad float64) {
	me.Rot.X = rad
	me.ApplyMatrices()
}

//	Sets me.Rot and calls ApplyMatrices().
func (me *NodeTransforms) SetRotXYZ(radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
	me.ApplyMatrices()
}

//	Sets me.Rot.Y and calls ApplyMatrices().
func (me *NodeTransforms) SetRotY(rad float64) {
	me.Rot.Y = rad
	me.ApplyMatrices()
}

//	Sets me.Rot.Z and calls ApplyMatrices().
func (me *NodeTransforms) SetRotZ(rad float64) {
	me.Rot.Z = rad
	me.ApplyMatrices()
}

//	Sets me.Scale and calls ApplyMatrices().
func (me *NodeTransforms) SetScale(scale *unum.Vec3) {
	if scale != nil {
		me.Scale = *scale
		me.ApplyMatrices()
	} else {
		me.SetScaleN(1)
	}
}

//	Sets me.Scale and calls ApplyMatrices().
func (me *NodeTransforms) SetScaleN(scale float64) {
	me.SetScaleXYZ(scale, scale, scale)
}

//	Sets me.Scale and calls ApplyMatrices().
func (me *NodeTransforms) SetScaleXYZ(x, y, z float64) {
	me.Scale.X, me.Scale.Y, me.Scale.Z = x, y, z
	me.ApplyMatrices()
}

//	Returns the result of multiplying deltaPerSecond with EngineLoop.TickDelta.
func (me *NodeTransforms) StepDelta(deltaPerSecond float64) float64 {
	return Loop.TickDelta * deltaPerSecond
}
