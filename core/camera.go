package core

import (
	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
	usl "github.com/metaleap/go-util/slice"
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
		matCamProj, matProj, matPos unum.Mat4
		nodeRender                  []bool
		nodeProjMats                []unum.Mat4
	}
	thrRend struct {
		nodeProjMats []ugl.GlMat4
		nodeRender   []bool
	}
}

func (me *Camera) init() {
	me.SetScene(-1)
	me.Perspective = Options.Cameras.PerspectiveDefaults
	unum.Mat4Identities(&me.thrApp.matProj, &me.thrPrep.matProj)
	me.Controller.init()
}

func (me *Camera) ensureProjMats(length int) {
	if len(me.thrPrep.nodeProjMats) < length {
		nu := make([]unum.Mat4, length)
		copy(nu, me.thrPrep.nodeProjMats)
		me.thrPrep.nodeProjMats = nu
	}
	if len(me.thrRend.nodeProjMats) < length {
		nu := make([]ugl.GlMat4, length)
		copy(nu, me.thrRend.nodeProjMats)
		me.thrRend.nodeProjMats = nu
	}
}

func (me *Camera) initNodeCamData(all SceneNodeLib, nodeID int) {
	me.ensureProjMats(len(all))
	me.thrPrep.nodeProjMats[nodeID].Identity()
	usl.BoolEnsureLen(&me.thrPrep.nodeRender, len(all))
	usl.BoolEnsureLen(&me.thrRend.nodeRender, len(all))
	me.thrPrep.nodeRender[nodeID] = all[nodeID].Render.Enabled
	me.thrRend.nodeRender[nodeID] = all[nodeID].Render.Enabled
}

func (me *Camera) Scene() *Scene {
	return Core.Libs.Scenes.get(me.sceneID)
}

func (me *Camera) SetScene(sceneID int) {
	if sceneID != me.sceneID {
		me.sceneID = sceneID
		if scene := me.Scene(); scene != nil {
			cap := len(scene.allNodes)
			me.ensureProjMats(cap)
			usl.BoolEnsureLen(&me.thrPrep.nodeRender, cap)
			usl.BoolEnsureLen(&me.thrRend.nodeRender, cap)
			for i := 0; i < len(scene.allNodes); i++ {
				if scene.allNodes.Ok(i) {
					me.initNodeCamData(scene.allNodes, i)
				}
			}
		}
	}
}
