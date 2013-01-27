package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

//	A camera embodies the eye point of the viewer looking at the visual scene.
type Camera3D struct {
	//	Optical and imager properties for this camera.
	Params struct {
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

func (me *Camera3D) init() {
	rend := &me.Rendering
	rend.Enabled = true
	rend.States.DepthTest, rend.States.FaceCulling, rend.States.ScissorTest, rend.States.StencilTest = true, true, true, false
	rend.States.ClearFlags = gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT
	rend.States.ClearColor = Core.Options.Rendering.DefaultClearColor
	opt := &me.Params
	opt.FovY = 37.8493
	opt.ZFar = 30000
	opt.ZNear = 0.3
	me.matProj.Identity()
	me.Controller.init()
	rend.ViewPort.init()
	me.ApplyMatrices()
	me.SetTechnique(Core.Options.Rendering.DefaultTechnique3D)
}

//	Applies changes made to the FovY, ZNear and/or ZFar parameters in me.Params.
func (me *Camera3D) ApplyMatrices() {
	me.matProj.Perspective(me.Params.FovY, me.Rendering.ViewPort.aspect, me.Params.ZNear, me.Params.ZFar)
}

func (me *Camera3D) dispose() {
}

func (me *Camera3D) render() {
	if me.Rendering.Enabled {
		Core.Rendering.states.Apply(&me.Rendering.States.RenderStatesBag)
		curScene = Core.Libs.Scenes[me.Rendering.SceneID]
		Core.useTechnique(me.technique)
		me.matCamProj.SetFromMult4(&me.matProj, &me.Controller.mat)
		//me.glMatCamProj.Load(&me.matCamProj)
		//gl.UniformMatrix4fv(curProg.UnifLocs["uMatCamProj"], 1, gl.FALSE, &me.glMatCamProj[0])
		// gl.UniformMatrix4fv(curProg.UnifLocs["uMatProj"], 1, gl.FALSE, &me.glmatProj[0])
		me.technique.onPreRender()
		Core.Rendering.states.ForceEnableScissorTest()
		gl.Scissor(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		gl.Viewport(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		gl.Clear(me.Rendering.States.ClearFlags)
		curScene.RootNode.render()
		Core.Rendering.states.ForceDisableScissorTest()
	}
}

func (me *Camera3D) SetTechnique(name string) {
	if (me.technique == nil) || (me.technique.name() != name) {
		me.technique = techs[name]
	}
}

//#begin-gt -gen-lib.gt T:Camera3D

//	Initializes and returns a new Camera3D with default parameters.
func NewCamera3D() (me *Camera3D) {
	me = &Camera3D{}
	me.init()
	return
}

//	A hash-table of Camera3Ds associated by IDs. Only for use in Core.Libs.
type LibCamera3Ds map[string]*Camera3D

//	Creates and initializes a new Camera3D with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibCamera3Ds) AddNew(id string) (obj *Camera3D) {
	obj = NewCamera3D()
	me[id] = obj
	return
}

func (me *LibCamera3Ds) ctor() {
	*me = LibCamera3Ds{}
}

func (me *LibCamera3Ds) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibCamera3Ds) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
