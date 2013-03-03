package core

import (
	unum "github.com/metaleap/go-util/num"
)

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
		onPrepNode                  func(*Node)
		matCamProj, matProj, matPos unum.Mat4
	}
}

func (me *Camera) dispose() {
	me.clearNodeMats()
}

func (me *Camera) init() {
	me.sceneID = -1
	me.thrPrep.onPrepNode = func(n *Node) { me.onPrepNode(n) }
	me.Perspective = Options.Cameras.PerspectiveDefaults
	unum.Mat4Identities(&me.thrApp.matProj, &me.thrPrep.matProj)
	me.Controller.init()
}

func (me *Camera) clearNodeMats() {
	if scene := me.scene(); scene != nil {
		scene.RootNode.Walk(func(node *Node) {
			delete(node.thrPrep.camProjMats, me)
			delete(node.thrRend.camProjMats, me)
			delete(node.thrPrep.camRender, me)
			delete(node.thrRend.camRender, me)
		})
	}
}

func (me *Camera) scene() *Scene {
	return Core.Libs.Scenes.get(me.sceneID)
}

func (me *Camera) SetScene(sceneID int) {
	if sceneID != me.sceneID {
		me.clearNodeMats()
		me.sceneID = sceneID
		if scene := me.scene(); scene != nil {
			scene.RootNode.Walk(func(node *Node) {
				node.initCamData(me)
			})
		}
	}
}
