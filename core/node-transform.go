package core

import (
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

//	Represents one or more transformations of a Node.
//	This is only used by Node objects, which initialize their NodeTransforms with the
//	proper defaults and associate themselves with their NodeTransforms. (Any other
//	NodeTransforms are invalid.)
//	
//	A single NodeTransforms encapsulates a number of private 4x4 matrices that are
//	set from its exported fields via the SetFoo() or OnFooChanged() methods, and
//	are multiplied together into a final single 4x4 transformation matrix via the
//	ApplyMatrices() (or OnFooChanged()) methods.
type NodeTransforms struct {
	//	Translation of the from origin, if any.
	Pos unum.Vec3

	//	Rotation for each axis in radians, if any.
	Rot unum.Vec3

	//	Scaling of this node, if any. Defaults to (1, 1, 1) for no scaling.
	Scale unum.Vec3

	//	An arbitrary additional transformation to apply. Defaults to the 4x4 identity matrix for 'none'.
	//	Must not be nil.
	Other unum.Mat4

	matModelView, matRotX, matRotY, matRotZ, matScale, matTranslation, matParent unum.Mat4
	glMatModelView                                                               ugl.GlMat4
	owner                                                                        *Node
}

func newNodeTransforms(owner *Node) (me *NodeTransforms) {
	me = &NodeTransforms{}
	me.owner = owner
	me.Scale.X, me.Scale.Y, me.Scale.Z = 1, 1, 1
	me.matScale.Scaling(&me.Scale)
	me.matTranslation.Translation(&me.Pos)
	me.matRotX.RotationX(0)
	me.matRotY.RotationY(0)
	me.matRotZ.RotationZ(0)
	me.Other.Identity()
	me.matModelView.Identity()
	me.glMatModelView.Load(&me.matModelView)
	return
}

//	Creates a final single 4x4 transformation matrix for all transformations
//	in me. It is this matrix that is used by the rendering pipeline.
func (me *NodeTransforms) ApplyMatrices() {
	if me.owner.parentNode != nil {
		me.matParent.CopyFrom(&me.owner.parentNode.Transform.matModelView)
	} else {
		me.matParent.Identity()
	}
	me.matModelView.SetFromMultN(&me.matParent, &me.matTranslation, &me.Other, &me.matScale, &me.matRotX, &me.matRotY, &me.matRotZ)
	me.glMatModelView.Load(&me.matModelView)
	for _, me.owner.curSubNode = range me.owner.SubNodes.M {
		me.owner.curSubNode.Transform.ApplyMatrices()
	}
}

//	Updates the translation matrix from the current me.Pos value.
//	Updates the scaling matrix from the current me.Scale value.
//	Updates the X, Y and Z axis rotation matrices from the current me.Rot value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnAnyChanged() {
	me.onScaleChanged()
	me.onPosChanged()
	me.onRotChanged()
	me.ApplyMatrices()
}

func (me *NodeTransforms) onPosChanged() {
	me.matTranslation.Translation(&me.Pos)
}

//	Updates the translation matrix from the current me.Pos value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnPosChanged() {
	me.onPosChanged()
	me.ApplyMatrices()
}

//	Updates the translation matrix from the current me.Pos value.
//	Updates the X, Y and Z axis rotation matrices from the current me.Rot value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnPosRotChanged() {
	me.onPosChanged()
	me.onRotChanged()
	me.ApplyMatrices()
}

//	Updates the translation matrix from the current me.Pos value.
//	Updates the scaling matrix from the current me.Scale value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnPosScaleChanged() {
	me.onPosChanged()
	me.onScaleChanged()
	me.ApplyMatrices()
}

func (me *NodeTransforms) onRotChanged() {
	me.onRotXChanged()
	me.onRotYChanged()
	me.onRotZChanged()
}

//	Updates the X, Y and Z axis rotation matrices from the current me.Rot value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnRotChanged() {
	me.onRotChanged()
	me.ApplyMatrices()
}

//	Updates the scaling matrix from the current me.Scale value.
//	Updates the X, Y and Z axis rotation matrices from the current me.Rot value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnRotScaleChanged() {
	me.onRotChanged()
	me.onScaleChanged()
	me.ApplyMatrices()
}

func (me *NodeTransforms) onRotXChanged() {
	me.matRotX.RotationX(me.Rot.X)
}

//	Updates the X-axis rotation matrix from the current me.Rot.X value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnRotXChanged() {
	me.onRotXChanged()
	me.ApplyMatrices()
}

func (me *NodeTransforms) onRotYChanged() {
	me.matRotY.RotationY(me.Rot.Y)
}

//	Updates the Y-axis rotation matrix from the current me.Rot.Y value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnRotYChanged() {
	me.onRotYChanged()
	me.ApplyMatrices()
}

func (me *NodeTransforms) onRotZChanged() {
	me.matRotZ.RotationZ(me.Rot.Z)
}

//	Updates the Z-axis rotation matrix from the current me.Rot.Z value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnRotZChanged() {
	me.onRotZChanged()
	me.ApplyMatrices()
}

func (me *NodeTransforms) onScaleChanged() {
	me.matScale.Scaling(&me.Scale)
}

//	Updates the scaling matrix from the current me.Scale value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnScaleChanged() {
	me.onScaleChanged()
	me.ApplyMatrices()
}

//	If pos is true, updates the translation matrix from the current me.Pos value.
//	If scale is true, updates the scaling matrix from the current me.Scale value.
//	If rotX is true, updates the X-axis rotation matrix from the current me.Rot.X value.
//	If rotY is true, updates the Y-axis rotation matrix from the current me.Rot.Y value.
//	If rotZ is true, updates the Z-axis rotation matrix from the current me.Rot.Z value.
//	Then calls ApplyMatrices().
func (me *NodeTransforms) OnSomeChanged(scale, pos, rotX, rotY, rotZ bool) {
	if scale {
		me.onScaleChanged()
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
	me.ApplyMatrices()
}

//	Sets me.Pos and calls OnPosChanged().
func (me *NodeTransforms) SetPos(pos *unum.Vec3) {
	if pos != nil {
		me.Pos = *pos
		me.OnPosChanged()
	} else {
		me.SetPosXYZ(0, 0, 0)
	}
}

//	Sets me.Pos.X and calls OnPosChanged().
func (me *NodeTransforms) SetPosX(posX float64) {
	me.Pos.X = posX
	me.OnPosChanged()
}

//	Sets me.Pos and calls OnPosChanged().
func (me *NodeTransforms) SetPosXYZ(posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
	me.OnPosChanged()
}

//	Sets me.Pos.Y and calls OnPosChanged().
func (me *NodeTransforms) SetPosY(posY float64) {
	me.Pos.Y = posY
	me.OnPosChanged()
}

//	Sets me.Pos.Z and calls OnPosChanged().
func (me *NodeTransforms) SetPosZ(posZ float64) {
	me.Pos.Z = posZ
	me.OnPosChanged()
}

//	Sets me.Rot and calls OnRotChanged().
func (me *NodeTransforms) SetRot(rot *unum.Vec3) {
	if rot != nil {
		me.Rot = *rot
		me.OnRotChanged()
	} else {
		me.SetRotXYZ(0, 0, 0)
	}
}

//	Sets me.Rot.X and calls OnRotChanged().
func (me *NodeTransforms) SetRotX(rad float64) {
	me.Rot.X = rad
	me.OnRotXChanged()
}

//	Sets me.Rot and calls OnRotChanged().
func (me *NodeTransforms) SetRotXYZ(radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
	me.OnRotChanged()
}

//	Sets me.Rot.Y and calls OnRotChanged().
func (me *NodeTransforms) SetRotY(rad float64) {
	me.Rot.Y = rad
	me.OnRotYChanged()
}

//	Sets me.Rot.Z and calls OnRotChanged().
func (me *NodeTransforms) SetRotZ(rad float64) {
	me.Rot.Z = rad
	me.OnRotZChanged()
}

//	Sets me.Scale and calls OnScaleChanged().
func (me *NodeTransforms) SetScale(scale *unum.Vec3) {
	if scale != nil {
		me.Scale = *scale
		me.OnScaleChanged()
	} else {
		me.SetScaleN(1)
	}
}

//	Sets me.Scale and calls OnScaleChanged().
func (me *NodeTransforms) SetScaleN(scale float64) {
	me.SetScaleXYZ(scale, scale, scale)
}

//	Sets me.Scale and calls OnScaleChanged().
func (me *NodeTransforms) SetScaleXYZ(x, y, z float64) {
	me.Scale.X, me.Scale.Y, me.Scale.Z = x, y, z
	me.OnScaleChanged()
}

//	Returns the result of multiplying deltaPerSecond with EngineLoop.TickDelta.
func (me *NodeTransforms) StepDelta(deltaPerSecond float64) float64 {
	return Loop.TickDelta * deltaPerSecond
}
