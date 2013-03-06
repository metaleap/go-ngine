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

	FrustumCull bool

	sceneID int

	thrApp struct {
		matProj unum.Mat4
	}
	thrPrep struct {
		frustum                     cameraFrustum
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
	// me.FrustumCull = true
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

func (me *Camera) updateFrustumPlanes(mat *unum.Mat4) {
	f := &me.thrPrep.frustum

	m11, m21, m31, m41 := mat[0], mat[4], mat[8], mat[12]
	m12, m22, m32, m42 := mat[1], mat[5], mat[9], mat[13]
	m13, m23, m33, m43 := mat[2], mat[6], mat[10], mat[14]
	m14, m24, m34, m44 := mat[3], mat[7], mat[11], mat[15]

	f[cfpLeft].X = m41 + m11
	f[cfpLeft].Y = m42 + m12
	f[cfpLeft].Z = m43 + m13
	f[cfpLeft].d = m44 + m14

	f[cfpRight].X = m41 - m11
	f[cfpRight].Y = m42 - m12
	f[cfpRight].Z = m43 - m13
	f[cfpRight].d = m44 - m14

	f[cfpBottom].X = m41 + m21
	f[cfpBottom].Y = m42 + m22
	f[cfpBottom].Z = m43 + m23
	f[cfpBottom].d = m44 + m24

	f[cfpTop].X = m41 - m21
	f[cfpTop].Y = m42 - m22
	f[cfpTop].Z = m43 - m23
	f[cfpTop].d = m44 - m24

	f[cfpNear].X = m41 + m31
	f[cfpNear].Y = m42 + m32
	f[cfpNear].Z = m43 + m33
	f[cfpNear].d = m44 + m34

	f[cfpFar].X = m41 - m31
	f[cfpFar].Y = m42 - m32
	f[cfpFar].Z = m43 - m33
	f[cfpFar].d = m44 - m34

	f.normalize()
}

type cfp int

const (
	cfpNear   cfp = 0
	cfpFar    cfp = 1
	cfpBottom cfp = 2
	cfpTop    cfp = 3
	cfpLeft   cfp = 4
	cfpRight  cfp = 5
)

type cameraFrustum [6]cameraFrustumPlane

func (me cameraFrustum) containsSphere(center *unum.Vec3, radius float64) bool {
	var dist float64
	for i := 0; i < len(me); i++ {
		dist = me[i].Dot(center) + me[i].d
		if dist < -radius {
			//	fully outside this plane
			return false
		}
		if math.Abs(dist) < radius {
			//	intersecting plane
			return true
		}
	}
	//	fully in frustum
	return true
}

func (me cameraFrustum) normalize() {
	for i := 0; i < len(me); i++ {
		me[i].normalize()
	}
}

type cameraFrustumPlane struct {
	unum.Vec3
	d float64
}

func (me *cameraFrustumPlane) normalize() {
	me.Vec3.Normalize()
	me.d = me.d / me.Vec3.Magnitude()
}

func (me *cameraFrustumPlane) signedDistTo(point *unum.Vec3) float64 {
	return me.Dot(point) + me.d
}
