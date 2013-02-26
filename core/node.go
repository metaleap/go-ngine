package core

import (
	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
)

type nodeCamProjMats map[*Camera]*unum.Mat4

type nodeCamProjGlMats map[*Camera]*ugl.GlMat4

type NodeVisitor func(*Node)

//	Declares a point of interest in a Scene.
type Node struct {
	Rendering struct {
		//	Defaults to true. If false, this Node is ignored by the rendering runtime.
		Enabled bool

		skyMode bool
	}

	//	Allows the Node to recursively define hierarchy.
	ChildNodes Nodes

	//	Encapsulates all parent-relative transformations for this Node.
	Transform NodeTransform

	thrPrep struct {
		copyDone, done bool
		matModelView   unum.Mat4
		camProjMats    nodeCamProjMats
		camRender      map[*Camera]bool
	}

	thrRend struct {
		copyDone    bool
		camProjMats nodeCamProjGlMats
	}

	MatID  int
	MeshID int

	model       *Model
	parentNode  *Node
	rootScene   *Scene
	modelID, id string
}

func newNode(id string, meshID int, modelID string, parent *Node, scene *Scene) (me *Node) {
	me = &Node{id: id, parentNode: parent, rootScene: scene, MatID: -1, MeshID: meshID}
	me.Rendering.Enabled = true
	me.Rendering.skyMode = (parent == nil)
	me.ChildNodes.init(me)
	me.Transform.init()
	me.ApplyTransform()
	me.SetModelID(modelID)
	me.initCamDatas()
	return
}

func (me *Node) ApplyTransform() {
	me.Transform.applyMatrices(me)
}

func (me *Node) initCamData(cam *Camera) {
	if cam.scene() == me.rootScene {
		me.thrPrep.camProjMats[cam], me.thrRend.camProjMats[cam] = unum.NewMat4Identity(), ugl.NewGlMat4(nil)
		me.thrPrep.camRender[cam] = me.Rendering.Enabled
	}
}

func (me *Node) initCamDatas() {
	me.thrPrep.camRender = map[*Camera]bool{}
	me.thrPrep.camProjMats, me.thrRend.camProjMats = nodeCamProjMats{}, nodeCamProjGlMats{}
	Core.Rendering.Canvases.Walk(nil, func(cam *Camera) {
		me.initCamData(cam)
	})
}

func (me *Node) Material() (mat *FxMaterial) {
	if mat = Core.Libs.Materials.Get(me.MatID); mat == nil {
		mat = Core.Libs.Materials.Get(me.model.MatID)
	}
	return
}

func (me *Node) ModelID() string {
	return me.modelID
}

func (me *Node) mesh() *Mesh {
	return Core.Libs.Meshes.Get(me.MeshID)
}

func (me *Node) Root() (root *Node) {
	if me.parentNode == nil {
		root = me
	} else {
		root = me.parentNode.Root()
	}
	return
}

func (me *Node) SetModelID(modelID string) {
	if mesh := me.mesh(); mesh == nil {
		me.model, me.modelID = nil, ""
	} else {
		me.model, me.modelID = mesh.Models.Default(), ""
		if modelID != me.modelID {
			me.model, me.modelID = mesh.Models[modelID], modelID
		}
	}
}

func (me *Node) Walk(onNode NodeVisitor) {
	onNode(me)
	for _, subNode := range me.ChildNodes.M {
		subNode.Walk(onNode)
	}
}
