package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

//	A camera embodies the eye point of the viewer looking at the visual scene.
type Camera2D struct {
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

	technique renderTechnique
}

func (me *Camera2D) init() {
	rend := &me.Rendering
	rend.Enabled = true
	rend.States.DepthTest, rend.States.FaceCulling, rend.States.ScissorTest, rend.States.StencilTest = false, true, true, false
	rend.States.ClearFlags = gl.COLOR_BUFFER_BIT
	rend.States.ClearColor = Core.Options.Rendering.DefaultClearColor
	rend.ViewPort.init()
	me.SetTechnique(Core.Options.Rendering.DefaultTechnique2D)
}

func (me *Camera2D) dispose() {
}

func (me *Camera2D) render() {
	if me.Rendering.Enabled {
		Core.Rendering.states.Apply(&me.Rendering.States.RenderStatesBag)
		curScene = Core.Libs.Scenes[me.Rendering.SceneID]
		Core.useTechnique(me.technique)
		me.technique.onPreRender()
		gl.Viewport(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		gl.Scissor(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		gl.Clear(me.Rendering.States.ClearFlags)
		curScene.RootNode.render()
	}
}

func (me *Camera2D) SetTechnique(name string) {
	if (me.technique == nil) || (me.technique.name() != name) {
		me.technique = techs[name]
	}
}

//#begin-gt -gen-lib.gt T:Camera2D

//	Initializes and returns a new Camera2D with default parameters.
func NewCamera2D() (me *Camera2D) {
	me = &Camera2D{}
	me.init()
	return
}

//	A hash-table of Camera2Ds associated by IDs. Only for use in Core.Libs.
type LibCamera2Ds map[string]*Camera2D

//	Creates and initializes a new Camera2D with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibCamera2Ds) AddNew(id string) (obj *Camera2D) {
	obj = NewCamera2D()
	me[id] = obj
	return
}

func (me *LibCamera2Ds) ctor() {
	*me = LibCamera2Ds{}
}

func (me *LibCamera2Ds) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibCamera2Ds) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
