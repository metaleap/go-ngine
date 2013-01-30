package core

import (
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type nodeCameraGlMats map[*Camera]*ugl.GlMat4

type nodeCameraMats map[*Camera]*unum.Mat4

//	Declares a point of interest in a Scene.
type Node struct {
	// matModelProj   unum.Mat4
	// glMatModelProj ugl.GlMat4

	//	Defaults to true. If false, this Node is ignored by the rendering runtime.
	Enabled bool

	//	Allows the Node to recursively define hierarchy.
	ChildNodes Nodes

	//	Encapsulates all parent-relative transformations for this Node.
	Transform NodeTransforms

	thrPrep struct {
		curId        string
		curSubNode   *Node
		matModelView unum.Mat4
		model        *Model
	}
	thrRend struct {
		curId          string
		curSubNode     *Node
		matModelView   unum.Mat4
		matModelProj   unum.Mat4
		glMatModelProj ugl.GlMat4
	}

	mat                               *FxMaterial
	mesh                              *Mesh
	model                             *Model
	curSubNode, parentNode            *Node
	curID, matID, meshID, modelID, id string
}

func newNode(id, meshID, modelID string, parent *Node) (me *Node) {
	me = &Node{id: id, parentNode: parent, Enabled: true}
	me.ChildNodes.init(me)
	me.Transform.init(me)
	me.SetMeshModelID(meshID, modelID)
	return
}

func (me *Node) EffectiveMaterial() *FxMaterial {
	if me.mat != nil {
		return me.mat
	}
	return me.model.mat
}

func (me *Node) MatID() string {
	return me.matID
}

func (me *Node) MeshID() string {
	return me.meshID
}

func (me *Node) ModelID() string {
	return me.modelID
}

func (me *Node) SetMatID(newMatID string) {
	if newMatID != me.matID {
		me.mat, me.matID = Core.Libs.Materials[newMatID], newMatID
	}
}

func (me *Node) SetMeshModelID(meshID, modelID string) {
	if meshID != me.meshID {
		me.mesh, me.meshID = Core.Libs.Meshes[meshID], meshID
	}
	if me.mesh == nil {
		me.model, me.modelID = nil, ""
	} else {
		me.model, me.modelID = me.mesh.Models.Default(), ""
		if modelID != me.modelID {
			me.model, me.modelID = me.mesh.Models[modelID], modelID
		}
	}
}

func (me *Node) Walk(onNode func(*Node)) {
	onNode(me)
	for me.curID, me.curSubNode = range me.ChildNodes.M {
		me.curSubNode.Walk(onNode)
	}
}
