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
	Transform SceneNodeTransform

	MatID   int
	MeshID  int
	ModelID int

	parentNode *Node
	rootScene  *Scene
	id         string

	thrPrep struct {
		matModelView unum.Mat4
		camProjMats  nodeCamProjMats
		camRender    map[*Camera]bool
	}

	thrRend struct {
		camProjMats nodeCamProjGlMats
		camRender   map[*Camera]bool
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

func (me *Node) initCamData(view *Camera) {
	if view.scene() == me.rootScene {
		me.thrPrep.camProjMats[view], me.thrRend.camProjMats[view] = unum.NewMat4Identity(), ugl.NewGlMat4(nil)
		me.thrPrep.camRender[view] = me.Rendering.Enabled
		me.thrRend.camRender[view] = me.Rendering.Enabled
	}
}

func (me *Node) initCamDatas() {
	me.thrPrep.camRender, me.thrRend.camRender = map[*Camera]bool{}, map[*Camera]bool{}
	me.thrPrep.camProjMats, me.thrRend.camProjMats = nodeCamProjMats{}, nodeCamProjGlMats{}
	var view int
	var rts *RenderTechniqueScene
	for canv := 0; canv < len(Core.Render.Canvases); canv++ {
		for view = 0; view < len(Core.Render.Canvases[canv].Views); view++ {
			if rts = Core.Render.Canvases[canv].Views[view].RenderTechniqueScene(); rts != nil {
				me.initCamData(&rts.Camera)
			}
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
