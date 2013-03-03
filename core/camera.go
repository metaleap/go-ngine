package core

import (
	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
)

type camNodeProjMats map[*Node]*unum.Mat4

type camNodeProjGlMats map[*Node]*ugl.GlMat4

type CameraPerspective struct {
	//	Whether this is a perspective-projection camera. Defaults to true.
	//	If false, no projection transformation is applied.
	Enabled bool

	//	Vertical field-of-view. Defaults to 37.8493.
	//	After changing this value, you must call the ApplyMatrices() method.
	FovY float64

	//	Distance of the far-plane from the camera. Defaults to 30000.
	//	After changing this value, you must call the ApplyMatrices() method.
	ZFar float64

	//	Distance of the near-plane from the camera. Defaults to 0.3.
	//	After changing this value, you must call the ApplyMatrices() method.
	ZNear float64
}

type Camera struct {
	//	Optical and imager properties for this camera.
	Perspective CameraPerspective

	//	Encapsulates the position and direction of this camera.
	Controller Controller

	sceneID int

	thrApp struct {
		matProj unum.Mat4
	}
	thrPrep struct {
		matCamProj, matProj, matPos unum.Mat4
		onPrepNode                  func(*Node)
		nodeRender                  map[*Node]bool
		nodeProjMats                camNodeProjMats
	}
	thrRend struct {
		nodeProjMats camNodeProjGlMats
		nodeRender   map[*Node]bool
	}
}

func (me *Camera) init() {
	me.SetScene(-1)
	me.thrPrep.onPrepNode = func(n *Node) { me.onPrepNode(n) }
	me.Perspective = Options.Cameras.PerspectiveDefaults
	unum.Mat4Identities(&me.thrApp.matProj, &me.thrPrep.matProj)
	me.Controller.init()
}

func (me *Camera) initNodeCamData(node *Node) {
	me.thrPrep.nodeProjMats[node] = unum.NewMat4Identity()
	me.thrRend.nodeProjMats[node] = ugl.NewGlMat4(nil)
	me.thrPrep.nodeRender[node] = node.Rendering.Enabled
	me.thrRend.nodeRender[node] = node.Rendering.Enabled
}

func (me *Camera) scene() *Scene {
	return Core.Libs.Scenes.get(me.sceneID)
}

func (me *Camera) SetScene(sceneID int) {
	if sceneID != me.sceneID {
		me.thrPrep.nodeRender, me.thrRend.nodeRender = make(map[*Node]bool, 1), make(map[*Node]bool, 1)
		me.thrPrep.nodeProjMats, me.thrRend.nodeProjMats = make(camNodeProjMats, 1), make(camNodeProjGlMats, 1)
		me.sceneID = sceneID
		if scene := me.scene(); scene != nil {
			scene.RootNode.Walk(func(node *Node) {
				me.initNodeCamData(node)
			})
		}
	}
}
