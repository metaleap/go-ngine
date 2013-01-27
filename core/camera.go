package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
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

	Rendering struct {
		Enabled bool

		//	The device-relative or absolute view-port for this Camera.
		ViewPort CameraViewPort

		//	The ID of the Scene (in Core.Libs.Scenes) this camera is looking at.
		SceneID string

		States struct {
			ugl.RenderStatesBag
			ClearFlags gl.Bitfield
		}
	}

	technique           renderTechnique
	matCamProj, matProj unum.Mat4
}

//	Initializes and returns a new Camera without any projection.
func NewCamera2D() (me *Camera) {
	me = &Camera{}
	me.init(false)
	return
}

//	Initializes and returns a new Camera with perspective projection.
func NewCamera3D() (me *Camera) {
	me = &Camera{}
	me.init(true)
	return
}

func (me *Camera) init(persp3d bool) {
	rend := &me.Rendering
	rend.Enabled = true
	rend.States.DepthTest, rend.States.FaceCulling, rend.States.ScissorTest, rend.States.StencilTest = persp3d, true, false, false
	rend.States.ClearColor = Core.Options.Rendering.DefaultClearColor
	if persp3d {
		rend.States.ClearFlags = gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT
	} else {
		rend.States.ClearFlags = gl.COLOR_BUFFER_BIT
	}
	persp := &me.Perspective
	persp.Use = persp3d
	persp.FovY = 37.8493
	persp.ZFar = 30000
	persp.ZNear = 0.3
	me.matProj.Identity()
	me.Controller.init()
	rend.ViewPort.init()
	me.ApplyMatrices()
	me.SetTechnique(Core.Options.Rendering.DefaultTechnique3D)
}

//	Applies changes made to the FovY, ZNear and/or ZFar parameters in me.Perspective.
func (me *Camera) ApplyMatrices() {
	me.matProj.Perspective(me.Perspective.FovY, me.Rendering.ViewPort.aspect, me.Perspective.ZNear, me.Perspective.ZFar)
}

func (me *Camera) dispose() {
}

func (me *Camera) render() {
	if me.Rendering.Enabled {
		Core.Rendering.states.Apply(&me.Rendering.States.RenderStatesBag)
		Core.Rendering.states.ForceEnableScissorTest()
		curScene = Core.Libs.Scenes[me.Rendering.SceneID]
		Core.useTechnique(me.technique)
		me.matCamProj.SetFromMult4(&me.matProj, &me.Controller.mat)
		//me.glMatCamProj.Load(&me.matCamProj)
		//gl.UniformMatrix4fv(curProg.UnifLocs["uMatCamProj"], 1, gl.FALSE, &me.glMatCamProj[0])
		// gl.UniformMatrix4fv(curProg.UnifLocs["uMatProj"], 1, gl.FALSE, &me.glmatProj[0])
		me.technique.onPreRender()
		gl.Scissor(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		gl.Viewport(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		gl.Clear(me.Rendering.States.ClearFlags)
		curScene.RootNode.render()
		Core.Rendering.states.ForceDisableScissorTest()
	}
}

func (me *Camera) SetTechnique(name string) {
	if (me.technique == nil) || (me.technique.name() != name) {
		me.technique = techs[name]
	}
}

type Cameras []*Camera

func (me *Cameras) Add(camera *Camera) *Camera {
	for _, cam := range *me {
		if cam == camera {
			return cam
		}
	}
	*me = append(*me, camera)
	return camera
}

func (me *Cameras) Insert(camera *Camera, index int) {
	if index >= len(*me) {
		me.Add(camera)
	} else {
		for _, cam := range *me {
			if cam == camera {
				return
			}
		}
		pre, post := (*me)[:index], (*me)[index:]
		*me = append(append(pre, camera), post...)
	}
}

func (me *Cameras) Remove(camera *Camera) {
	for i, cam := range *me {
		if cam == camera {
			pre, post := (*me)[:i], (*me)[i+1:]
			*me = append(pre, post...)
		}
	}
}

//	Encapsulates a device-relative or absolute camera view-port.
type CameraViewPort struct {
	absolute               bool
	relX, relY, relW, relH float64
	absX, absY, absW, absH int
	aspect                 float64
	glVpX, glVpY           gl.Int
	glVpW, glVpH           gl.Sizei
}

func (me *CameraViewPort) init() {
	me.SetRel(0, 0, 1, 1)
}

//	Sets the absolute viewport origin and dimensions in pixels.
func (me *CameraViewPort) SetAbs(x, y, width, height int) {
	me.absolute, me.absX, me.absY, me.absW, me.absH = true, x, y, width, height
	me.update()
}

//	Sets the device-relative viewport origin and dimensions, with the value 1.0
//	representing the maximum extent of the viewport on that respective axis.
func (me *CameraViewPort) SetRel(x, y, width, height float64) {
	me.absolute, me.relX, me.relY, me.relW, me.relH = false, x, y, width, height
	me.update()
}

func (me *CameraViewPort) update() {
	if !me.absolute {
		me.absW, me.absH = int(me.relW*float64(curCanvas.absViewWidth)), int(me.relH*float64(curCanvas.absViewHeight))
		me.absX, me.absY = int(me.relX*float64(curCanvas.absViewWidth)), int(me.relY*float64(curCanvas.absViewHeight))
	}
	me.glVpX, me.glVpY, me.glVpW, me.glVpH = gl.Int(me.absX), gl.Int(me.absY), gl.Sizei(me.absW), gl.Sizei(me.absH)
	me.aspect = float64(me.absW) / float64(me.absH)
}
