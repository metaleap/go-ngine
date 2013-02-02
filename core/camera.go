package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
)

//	A camera embodies the eye point of the viewer looking at the visual scene.
type Camera struct {
	//	Optical and imager properties for this camera.
	Perspective struct {
		//	Whether this is a perspective-projection camera. Defaults to true.
		//	If false, no projection transformation is applied.
		Use bool

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

	//	Encapsulates the position and direction of this camera.
	Controller Controller

	Enabled bool

	Rendering struct {
		//	The device-relative or absolute view-port for this Camera.
		Viewport CameraViewport

		States ugl.RenderStatesBag
	}

	scene *Scene

	thrApp struct {
		matProj unum.Mat4
	}
	thrPrep struct {
		matCamProj, matProj unum.Mat4
	}
	thrRend struct {
		states    ugl.RenderStatesBag
		technique renderTechnique
	}
}

func newCamera2D(canv *RenderCanvas, depth bool) (me *Camera) {
	me = &Camera{}
	me.init(canv, false, depth)
	return
}

func newCamera3D(canv *RenderCanvas) (me *Camera) {
	me = &Camera{}
	me.init(canv, true, true)
	return
}

func (me *Camera) init(canv *RenderCanvas, persp3d bool, depth bool) {
	me.Enabled = true
	rend := &me.Rendering
	rend.Viewport.canvas = canv
	rend.States.DepthTest, rend.States.FaceCulling, rend.States.StencilTest = depth, true, false
	rend.States.ClearColor = Core.Options.Rendering.DefaultClearColor
	if depth {
		rend.States.Other.ClearBits = gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT
	} else {
		rend.States.Other.ClearBits = gl.COLOR_BUFFER_BIT
	}
	persp := &me.Perspective
	persp.Use = persp3d
	persp.FovY = 37.8493
	persp.ZFar = 30000
	persp.ZNear = 0.3
	me.thrApp.matProj.Identity()
	me.thrPrep.matProj.Identity()
	me.Controller.init()
	rend.Viewport.init()
	me.ApplyMatrices()
	me.SetTechnique(Core.Options.Rendering.DefaultTechnique3D)
}

//	Applies changes made to the FovY, ZNear and/or ZFar parameters in me.Perspective.
func (me *Camera) ApplyMatrices() {
	me.thrApp.matProj.Perspective(&thrApp.numBag, me.Perspective.FovY, me.Rendering.Viewport.aspect, me.Perspective.ZNear, me.Perspective.ZFar)
}

func (me *Camera) clearNodeMats() {
	if me.scene != nil {
		me.scene.RootNode.Walk(func(node *Node) {
			delete(node.thrPrep.matProjs, me)
			delete(node.thrRend.matProjs, me)
		})
	}
}

func (me *Camera) dispose() {
	me.clearNodeMats()
}

func (me *Camera) Scene() *Scene {
	return me.scene
}

func (me *Camera) setScene(scene *Scene) {
	if scene != me.scene {
		me.clearNodeMats()
		if me.scene = scene; me.scene != nil {
			me.scene.RootNode.Walk(func(node *Node) {
				node.initProjMat(me)
			})
		}
	}
}

func (me *Camera) SetScene(id string) {
	me.setScene(Core.Libs.Scenes[id])
}

func (me *Camera) SetTechnique(name string) {
	if (me.thrRend.technique == nil) || (me.thrRend.technique.name() != name) {
		me.thrRend.technique = techs[name]
	}
}

type Cameras []*Camera

func (me *Cameras) Remove(camera *Camera) {
	for i, cam := range *me {
		if cam == camera {
			pre, post := (*me)[:i], (*me)[i+1:]
			*me = append(pre, post...)
		}
	}
}

//	Encapsulates a device-relative or absolute camera view-port.
type CameraViewport struct {
	absolute               bool
	relX, relY, relW, relH float64
	absX, absY, absW, absH int
	aspect                 float64
	glVpX, glVpY           gl.Int
	glVpW, glVpH           gl.Sizei
	canvas                 *RenderCanvas
}

func (me *CameraViewport) init() {
	me.SetRel(0, 0, 1, 1)
}

//	Sets the absolute viewport origin and dimensions in pixels.
func (me *CameraViewport) SetAbs(x, y, width, height int) {
	me.absolute, me.absX, me.absY, me.absW, me.absH = true, x, y, width, height
	me.update()
}

//	Sets the device-relative viewport origin and dimensions, with the value 1.0
//	representing the maximum extent of the viewport on that respective axis.
func (me *CameraViewport) SetRel(x, y, width, height float64) {
	me.absolute, me.relX, me.relY, me.relW, me.relH = false, x, y, width, height
	me.update()
}

func (me *CameraViewport) update() {
	if !me.absolute {
		me.absW, me.absH = int(me.relW*float64(me.canvas.absViewWidth)), int(me.relH*float64(me.canvas.absViewHeight))
		me.absX, me.absY = int(me.relX*float64(me.canvas.absViewWidth)), int(me.relY*float64(me.canvas.absViewHeight))
	}
	me.glVpX, me.glVpY, me.glVpW, me.glVpH = gl.Int(me.absX), gl.Int(me.absY), gl.Sizei(me.absW), gl.Sizei(me.absH)
	me.aspect = float64(me.absW) / float64(me.absH)
}
