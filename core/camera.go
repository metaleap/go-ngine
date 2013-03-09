package core

import (
	"math"

	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
	usl "github.com/metaleap/go-util/slice"
)

type CameraPerspective struct {
	//	Whether this is a perspective-projection camera. Defaults to true.
	//	If false, no projection transformation is applied.
	Enabled bool

	//	Vertical field-of-view angle.
	FovY struct {
		//	In degrees. Defaults to 37.8493.
		//	After changing this value, call Camera.ApplyPerspective().
		Deg float64

		//	In radians, times 0.5. This is auto-set (from Deg) by Camera.ApplyPerspective().
		RadHalf float64

		tanRadHalf, tanRadHalfAspect float64
	}

	//	Distance of the far-plane from the camera. Defaults to 30000.
	//	After changing this value, call Camera.ApplyPerspective().
	ZFar float64

	//	Distance of the near-plane from the camera. Defaults to 0.3.
	//	After changing this value, call Camera.ApplyPerspective().
	ZNear float64
}

type Camera struct {
	//	Optical and imager properties for this camera.
	Perspective CameraPerspective

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
		frustum                     struct {
			axes struct {
				x, y, z unum.Vec3
			}
			sphereFactor unum.Vec2
		}
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
		me.Perspective.FovY.tanRadHalf = math.Tan(me.Perspective.FovY.RadHalf)
		me.Perspective.FovY.tanRadHalfAspect = me.Perspective.FovY.tanRadHalf * me.viewportAspectRatio
		me.thrPrep.frustum.sphereFactor.Y = 1 / math.Cos(me.Perspective.FovY.RadHalf)
		me.thrPrep.frustum.sphereFactor.X = 1 / math.Cos(math.Atan(me.Perspective.FovY.tanRadHalfAspect))
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
	usl.BoolEnsureLen(&me.thrPrep.nodeRender, len(all))
	usl.BoolEnsureLen(&me.thrRend.nodeRender, len(all))
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

func (me *Camera) frustumHasPoint(point *unum.Vec3) bool {
	var axisPos float64
	pp := point.Sub(&me.Controller.thrPrep.pos)
	if axisPos = pp.Dot(&me.thrPrep.frustum.axes.z); axisPos > me.Perspective.ZFar || axisPos < me.Perspective.ZNear {
		return false
	}
	halfHeight := axisPos * me.Perspective.FovY.tanRadHalf
	if axisPos = pp.Dot(&me.thrPrep.frustum.axes.y); -halfHeight > axisPos || axisPos > halfHeight {
		return false
	}
	halfWidth := halfHeight * me.viewportAspectRatio
	if axisPos = pp.Dot(&me.thrPrep.frustum.axes.x); -halfWidth > axisPos || axisPos > halfWidth {
		return false
	}
	return true
}

func (me *Camera) frustumHasSphere(center *unum.Vec3, radius float64) bool {
	if radius == 0 {
		return me.frustumHasPoint(center)
	}
	var axPos, z, d float64
	cp := center.Sub(&me.Controller.thrPrep.pos)
	if axPos = cp.Dot(&me.thrPrep.frustum.axes.z); axPos > me.Perspective.ZFar+radius || axPos < me.Perspective.ZNear-radius {
		return false
	}

	z, d = axPos*me.Perspective.FovY.tanRadHalfAspect, me.thrPrep.frustum.sphereFactor.X*radius
	if axPos = cp.Dot(&me.thrPrep.frustum.axes.x); axPos > z+d || axPos < -z-d {
		return false
	}

	z, d = z/me.viewportAspectRatio, me.thrPrep.frustum.sphereFactor.Y*radius
	if axPos = cp.Dot(&me.thrPrep.frustum.axes.y); axPos > z+d || axPos < -z-d {
		return false
	}

	// if pc.Z > me.Perspective.ZFar-radius || pc.Z < me.Perspective.ZNear+radius {
	// 	intersect = true
	// }
	// if pc.Y > zz.Y-d.Y || pc.Y < -zz.Y+d.Y {
	// 	intersect = true
	// }
	// if pc.X > zz.X-d.X || pc.X < -zz.X+d.X {
	// 	intersect = true
	// }

	return true
}
