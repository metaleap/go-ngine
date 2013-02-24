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
//	whenever its parent Node.ApplyTransform() method is called.
type NodeTransform struct {
	//	Translation of the from origin.
	Pos unum.Vec3

	//	Rotation for each axis in radians.
	Rot unum.Vec3

	//	Scaling of this node, if any. Defaults to (1, 1, 1) for no scaling.
	Scale unum.Vec3

	// Other unum.Mat4

	matModelView unum.Mat4
}

func (me *NodeTransform) init() {
	me.Scale.X, me.Scale.Y, me.Scale.Z = 1, 1, 1
	me.matModelView.Identity()
}

func (me *NodeTransform) AddRot(rot *unum.Vec3) {
	me.Rot.Add(rot)
}

func (me *NodeTransform) AddRotXYZ(x, y, z float64) {
	me.Rot.Add3(x, y, z)
}

//	Updates the internal 4x4 transformation matrix for all transformations
//	in me. It is only this matrix that is used by the rendering runtime.
func (me *NodeTransform) applyMatrices(owner *Node) {
	var matParent, matTrans, matScale, matRotX, matRotY, matRotZ unum.Mat4
	matScale.Scaling(&me.Scale)
	matTrans.Translation(&me.Pos)
	matRotX.RotationX(me.Rot.X)
	matRotY.RotationY(me.Rot.Y)
	matRotZ.RotationZ(me.Rot.Z)
	if owner.parentNode == nil {
		matParent.Identity()
	} else {
		matParent.CopyFrom(&owner.parentNode.Transform.matModelView)
	}
	me.matModelView.SetFromMultN(&matParent, &matTrans /*me.Other,*/, &matScale, &matRotX, &matRotY, &matRotZ)
	for _, subNode := range owner.ChildNodes.M {
		subNode.Transform.applyMatrices(subNode)
	}
}

func (me *NodeTransform) SetPos(posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
}

func (me *NodeTransform) SetRot(radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
}

func (me *NodeTransform) SetScale(s float64) {
	me.Scale.X, me.Scale.Y, me.Scale.Z = s, s, s
}

func (me *NodeTransform) SetScaleXyz(x, y, z float64) {
	me.Scale.X, me.Scale.Y, me.Scale.Z = x, y, z
}

//	Returns the result of multiplying deltaPerSecond with EngineLoop.TickDelta.
func (me *NodeTransform) StepDelta(deltaPerSecond float64) float64 {
	return Loop.Tick.Delta * deltaPerSecond
}
