package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

//	Declares a point of interest in a Scene.
type Node struct {
	matModelProj   unum.Mat4
	glMatModelProj ugl.GlMat4

	//	If true, this Node is ignored by the rendering runtime.
	Disabled bool

	//	Allows the Node to recursively define hierarchy.
	ChildNodes Nodes

	//	Encapsulates all parent-relative transformations for this Node.
	Transform NodeTransforms

	mat                                *FxMaterial
	mesh                               *Mesh
	model                              *Model
	curSubNode, parentNode             *Node
	curKey, matID, meshID, modelID, id string
}

func newNode(id, meshID, modelID string, parent *Node) (me *Node) {
	me = &Node{id: id, parentNode: parent}
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

func (me *Node) render() {
	if !me.Disabled {
		curNode, curMesh, curModel = me, me.mesh, me.model
		if me.model != nil {
			curTechnique.onRenderNode()
			if curCam.Perspective.Use {
				me.matModelProj.SetFromMult4(&curCam.matCamProj, &me.Transform.matModelView)
			} else {
				me.matModelProj = me.Transform.matModelView
			}
			me.glMatModelProj.Load(&me.matModelProj)
			gl.UniformMatrix4fv(curProg.UnifLocs["uMatModelProj"], 1, gl.FALSE, &me.glMatModelProj[0])
			me.model.render()
		}
		for me.curKey, me.curSubNode = range me.ChildNodes.M {
			me.curSubNode.render()
		}
	}
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
