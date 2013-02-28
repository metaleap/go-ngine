package core

import (
	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
)

type refCanvCam struct {
	canvID, camID int
}

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

	MatID   int
	MeshID  int
	ModelID int

	parentNode *Node
	rootScene  *Scene
	id         string

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
}

func newNode(id string, meshID int, parent *Node, scene *Scene) (me *Node) {
	me = &Node{id: id, parentNode: parent, rootScene: scene, MatID: -1, MeshID: meshID, ModelID: -1}
	me.Rendering.Enabled = true
	me.Rendering.skyMode = (parent == nil)
	me.ChildNodes.init(me)
	me.Transform.init()
	me.ApplyTransform()
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
	var cam int
	for canv := 0; canv < len(Core.Render.Canvases); canv++ {
		for cam = 0; cam < len(Core.Render.Canvases[canv].Cams); cam++ {
			me.initCamData(Core.Render.Canvases[canv].Cams[cam])
		}
	}
}

func (me *Node) material() (mat *FxMaterial) {
	if mat = Core.Libs.Materials.get(me.MatID); mat == nil {
		if model := me.model(); model != nil {
			mat = Core.Libs.Materials.get(model.MatID)
		}
	}
	return
}

func (me *Node) mesh() *Mesh {
	return Core.Libs.Meshes.get(me.MeshID)
}

func (me *Node) model() *Model {
	return Core.Libs.Models.get(me.modelID())
}

func (me *Node) modelID() (id int) {
	if id = me.ModelID; id < 0 {
		if mesh := me.mesh(); mesh != nil {
			id = mesh.DefaultModelID
		}
	}
	return
}

func (me *Node) Root() (root *Node) {
	if me.parentNode == nil {
		root = me
	} else {
		root = me.parentNode.Root()
	}
	return
}

func (me *Node) Walk(onNode NodeVisitor) {
	onNode(me)
	for _, subNode := range me.ChildNodes.M {
		subNode.Walk(onNode)
	}
}
