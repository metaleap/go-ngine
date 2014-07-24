package core

import (
	"github.com/go-utils/unum"
	u3d "github.com/go3d/go-3dutil"
)

type nodeBounds struct {
	full, self u3d.Bounds
}

//	Represents one or more transformations of a Node.
//	This is only used by Node objects, which initialize their SceneNodeTransform with the
//	proper defaults and associate themselves with their SceneNodeTransform. (Any other
//	SceneNodeTransform are invalid.)
//
//	A single SceneNodeTransform encapsulates an unexported 4x4 matrix that is recalculated
//	whenever its parent Node.ApplyTransform() method is called.
type SceneNodeTransform struct {
	//	Translation of the from origin.
	Pos unum.Vec3

	//	Rotation for each axis in radians.
	Rot unum.Vec3

	//	Scaling of this node, if any. Defaults to (1, 1, 1) for no scaling.
	Scale unum.Vec3

	// Other unum.Mat4

	thrApp struct {
		matModelView unum.Mat4
	}
	thrPrep struct {
		matModelView unum.Mat4
	}
}

func (me *SceneNodeTransform) init() {
	me.Scale.X, me.Scale.Y, me.Scale.Z = 1, 1, 1
	me.thrApp.matModelView.Identity()
}

func (me *SceneNodeTransform) AddRot(rot *unum.Vec3) {
	me.Rot.Add(rot)
}

func (me *SceneNodeTransform) AddRotXYZ(x, y, z float64) {
	me.Rot.Add3(x, y, z)
}

func (me *SceneNodeTransform) SetPos(posX, posY, posZ float64) {
	me.Pos.X, me.Pos.Y, me.Pos.Z = posX, posY, posZ
}

func (me *SceneNodeTransform) SetRot(radX, radY, radZ float64) {
	me.Rot.X, me.Rot.Y, me.Rot.Z = radX, radY, radZ
}

func (me *SceneNodeTransform) SetScale(s float64) {
	me.Scale.X, me.Scale.Y, me.Scale.Z = s, s, s
}

func (me *SceneNodeTransform) SetScaleXyz(x, y, z float64) {
	me.Scale.X, me.Scale.Y, me.Scale.Z = x, y, z
}

//	Returns the result of multiplying deltaPerSecond with EngineLoop.TickDelta.
func (me *SceneNodeTransform) StepDelta(deltaPerSecond float64) float64 {
	return Loop.Tick.Delta * deltaPerSecond
}

//	Updates the internal 4x4 transformation matrix for all transformations of the specified
//	node and child-nodes. It is only this matrix that is used by the rendering runtime.
func (me *Scene) ApplyNodeTransforms(nodeID int) {
	if me.allNodes.IsOk(nodeID) {
		//	this node
		var matParent, matTrans, matScale, matRotX, matRotY, matRotZ unum.Mat4
		matScale.Scaling(&me.allNodes[nodeID].Transform.Scale)
		matTrans.Translation(&me.allNodes[nodeID].Transform.Pos)
		matRotX.RotationX(me.allNodes[nodeID].Transform.Rot.X)
		matRotY.RotationY(me.allNodes[nodeID].Transform.Rot.Y)
		matRotZ.RotationZ(me.allNodes[nodeID].Transform.Rot.Z)
		if me.allNodes[nodeID].parentID < 0 {
			matParent.Identity()
		} else {
			matParent.CopyFrom(&me.allNodes[me.allNodes[nodeID].parentID].Transform.thrApp.matModelView)
		}
		me.allNodes[nodeID].Transform.thrApp.matModelView.SetFromMultN(&matParent, &matTrans /*me.Other,*/, &matScale, &matRotX, &matRotY, &matRotZ)
		//	child-nodes
		for i := 0; i < len(me.allNodes[nodeID].childNodeIDs); i++ {
			me.ApplyNodeTransforms(me.allNodes[nodeID].childNodeIDs[i])
		}
		me.allNodes[nodeID].thrApp.bounding.full.Clear()
		me.allNodes[nodeID].thrApp.bounding.self.Clear()
		if Core.Libs.Meshes.IsOk(me.allNodes[nodeID].Render.meshID) {
			me.applyBounds(nodeID, &Core.Libs.Meshes[me.allNodes[nodeID].Render.meshID].raw.bounding)
		} else {
			//	this node has no geometry of its own but its child-nodes might
			me.applyBounds(nodeID, nil)
		}
	}
}

func (me *Scene) applyBounds(n int, src *u3d.Bounds) {
	if src != nil {
		me.allNodes[n].thrApp.bounding.self.AaBox = src.AaBox
		me.allNodes[n].thrApp.bounding.self.AaBox.Transform(&me.allNodes[n].Transform.thrApp.matModelView)
		me.allNodes[n].thrApp.bounding.self.Sphere = me.allNodes[n].thrApp.bounding.self.AaBox.BoundingSphere(&me.allNodes[n].Transform.Pos)
	}
	me.allNodes[n].thrApp.bounding.full = me.allNodes[n].thrApp.bounding.self
	for _, cid := range me.allNodes[n].childNodeIDs {
		me.allNodes[n].thrApp.bounding.full.AaBox.UpdateMinMaxFrom(&me.allNodes[cid].thrApp.bounding.full.AaBox)
	}
	me.allNodes[n].thrApp.bounding.full.Sphere = me.allNodes[n].thrApp.bounding.full.AaBox.BoundingSphere(&me.allNodes[n].Transform.Pos)
}
