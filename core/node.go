package core

import (
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type nodeCamProjMats map[*Camera]*unum.Mat4

type nodeCamProjGlMats map[*Camera]*ugl.GlMat4

//	Declares a point of interest in a Scene.
type Node struct {
	//	Defaults to true. If false, this Node is ignored by the rendering runtime.
	Enabled bool

	//	Allows the Node to recursively define hierarchy.
	ChildNodes Nodes

	//	Encapsulates all parent-relative transformations for this Node.
	Transform NodeTransform

	thrPrep struct {
		copyDone, done bool
		matModelView   unum.Mat4
		tmpMat         *unum.Mat4
		model          *Model
		curId          string
		curSubNode     *Node
		tmpCam         *Camera
		matProjs       nodeCamProjMats
	}

	thrRend struct {
		copyDone   bool
		matProjs   nodeCamProjGlMats
		curId      string
		curSubNode *Node
	}

	mat                               *FxMaterial
	mesh                              *Mesh
	model                             *Model
	curSubNode, parentNode            *Node
	rootScene                         *Scene
	curID, matID, meshID, modelID, id string
}

func newNode(id, meshID, modelID string, parent *Node, scene *Scene) (me *Node) {
	me = &Node{id: id, parentNode: parent, Enabled: true, rootScene: scene}
	me.ChildNodes.init(me)
	me.Transform.init(me)
	me.SetMeshModelID(meshID, modelID)
	me.initProjMats()
	return
}

func (me *Node) EffectiveMaterial() *FxMaterial {
	if me.mat != nil {
		return me.mat
	}
	return me.model.mat
}

func (me *Node) initProjMat(cam *Camera) {
	if cam.scene == me.rootScene {
		me.thrPrep.matProjs[cam], me.thrRend.matProjs[cam] = unum.NewMat4Identity(), ugl.NewGlMat4(nil)
	}
}

func (me *Node) initProjMats() {
	me.thrPrep.matProjs, me.thrRend.matProjs = nodeCamProjMats{}, nodeCamProjGlMats{}
	Core.Rendering.Canvases.Walk(func(cam *Camera) {
		me.initProjMat(cam)
	})
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

func (me *Node) Root() (root *Node) {
	if me.parentNode == nil {
		root = me
	} else {
		root = me.parentNode.Root()
	}
	return
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
