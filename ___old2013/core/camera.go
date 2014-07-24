package core

import (
	"github.com/go-utils/unum"
	"github.com/go-utils/uslice"
	u3d "github.com/go3d/go-3dutil"
	ugl "github.com/go3d/go-opengl/util"
)

type Camera struct {
	//	Optical and imager properties for this camera.
	Perspective u3d.Perspective

	//	Encapsulates the position and direction of this camera.
	Controller Controller

	Cull struct {
		Frustum bool
	}

	viewportAspectRatio float64 // copied over from the parent RenderView.Port
	sceneID             int

	thrApp struct {
		matProj unum.Mat4
	}
	thrPrep struct {
		matCamProj, matProj, matPos unum.Mat4
		nodeRender                  []bool
		nodeProjMats                []unum.Mat4
		frustum                     u3d.Frustum
	}
	thrRend struct {
		nodeProjMats []ugl.GlMat4
		nodeRender   []bool
	}
}

func (me *Camera) init() {
	me.Cull.Frustum = true
	me.SetScene(-1)
	me.Perspective = Options.Cameras.PerspectiveDefaults
	unum.Mat4Identities(&me.thrApp.matProj, &me.thrPrep.matProj)
	me.Controller.init()
}

//	Applies changes made to the Enabled, FovY, ZNear and/or ZFar parameters in me.Camera.Perspective.
func (me *Camera) applyPerspective() {
	if me.Perspective.Enabled {
		me.Perspective.FovY.RadHalf = me.thrApp.matProj.Perspective(me.Perspective.FovY.Deg, me.viewportAspectRatio, me.Perspective.ZNear, me.Perspective.ZFar)
		me.thrPrep.frustum.UpdateRatio(&me.Perspective, me.viewportAspectRatio)
	} else {
		me.thrApp.matProj.Identity()
	}
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
	uslice.BoolEnsureLen(&me.thrPrep.nodeRender, len(all))
	uslice.BoolEnsureLen(&me.thrRend.nodeRender, len(all))
	me.thrPrep.nodeRender[nodeID] = false
	me.thrRend.nodeRender[nodeID] = false
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
			uslice.BoolEnsureLen(&me.thrPrep.nodeRender, cap)
			uslice.BoolEnsureLen(&me.thrRend.nodeRender, cap)
			for i := 0; i < len(scene.allNodes); i++ {
				if scene.allNodes.Ok(i) {
					me.initNodeCamData(scene.allNodes, i)
				}
			}
		}
	}
}

func (me *Camera) frustumHasPoint(point *unum.Vec3) bool {
	return me.thrPrep.frustum.HasPoint(&me.Controller.thrPrep.pos, point, me.Perspective.ZNear, me.Perspective.ZFar)
}

func (me *Camera) frustumHasSphere(center *unum.Vec3, radius float64) (fullyInside, intersect bool) {
	return me.thrPrep.frustum.HasSphere(&me.Controller.thrPrep.pos, center, radius, me.Perspective.ZNear, me.Perspective.ZFar)
}
