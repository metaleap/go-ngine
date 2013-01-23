package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

//	A camera embodies the eye point of the viewer looking at the visual scene.
type Camera struct {
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

		//	The ID of the Scene (in Core.Libs.Scenes) this camera is looking at.
		SceneID string
	}

	//	Encapsulates the position and direction of this camera.
	Controller Controller

	Rendering struct {
		Enabled bool

		//	The device-relative or absolute view-port for this Camera.
		ViewPort CameraViewPort

		States struct {
			DepthTest   bool
			FaceCulling bool
			ClearColor  ugl.GlVec4
		}
	}

	technique           renderTechnique
	matCamProj, matProj unum.Mat4
}

func (me *Camera) init() {
	me.Rendering.States.DepthTest, me.Rendering.Enabled, me.Rendering.States.FaceCulling = true, true, true
	me.Rendering.States.ClearColor = Core.Options.Rendering.DefaultClearColor
	opt := &me.Params
	opt.FovY = 37.8493
	opt.ZFar = 30000
	opt.ZNear = 0.3
	me.matProj.Identity()
	me.Controller.init()
	me.Rendering.ViewPort.init()
	me.ApplyMatrices()
	me.SetTechnique(Core.Options.Rendering.DefaultTechnique)
}

//	Applies changes made to the FovY, ZNear and/or ZFar parameters in me.Params.
func (me *Camera) ApplyMatrices() {
	me.matProj.Perspective(me.Params.FovY, me.Rendering.ViewPort.aspect, me.Params.ZNear, me.Params.ZFar)
}

func (me *Camera) dispose() {
}

func (me *Camera) render() {
	if me.Rendering.Enabled {
		Core.Rendering.states.SetClearColor(me.Rendering.States.ClearColor)
		Core.Rendering.states.SetDepthTest(me.Rendering.States.DepthTest)
		Core.Rendering.states.SetFaceCulling(me.Rendering.States.FaceCulling)
		curScene = Core.Libs.Scenes[me.Params.SceneID]
		Core.useTechnique(me.technique)
		me.matCamProj.SetFromMult4(&me.matProj, &me.Controller.mat)
		//me.glMatCamProj.Load(&me.matCamProj)
		//gl.UniformMatrix4fv(curProg.UnifLocs["uMatCamProj"], 1, gl.FALSE, &me.glMatCamProj[0])
		// gl.UniformMatrix4fv(curProg.UnifLocs["uMatProj"], 1, gl.FALSE, &me.glmatProj[0])
		me.technique.onPreRender()
		gl.Viewport(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		curScene.RootNode.render()
	}
}

func (me *Camera) SetTechnique(name string) {
	if (me.technique == nil) || (me.technique.name() != name) {
		me.technique = techs[name]
	}
}

func (me *Camera) ToggleTechnique() {
	var (
		curIndex, i int
		tech        renderTechnique = nil
	)
	allNames, curTech, name := glShaderMan.names, curTechnique.name(), ""
	for i, name = range allNames {
		if name == curTech {
			curIndex = i
			break
		}
	}
	if curIndex < (len(allNames) - 1) {
		for i = curIndex + 1; i < len(allNames); i++ {
			if tech = techs[allNames[i]]; tech != nil {
				break
			}
		}
	}
	if tech == nil {
		for i = 0; i < curIndex; i++ {
			if tech = techs[allNames[i]]; tech != nil {
				break
			}
		}
	}
	if tech != nil {
		me.technique = tech
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

//#begin-gt -gen-lib.gt T:Camera

//	Initializes and returns a new Camera with default parameters.
func NewCamera() (me *Camera) {
	me = &Camera{}
	me.init()
	return
}

//	A hash-table of Cameras associated by IDs. Only for use in Core.Libs.
type LibCameras map[string]*Camera

//	Creates and initializes a new Camera with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibCameras) AddNew(id string) (obj *Camera) {
	obj = NewCamera()
	me[id] = obj
	return
}

func (me *LibCameras) ctor() {
	*me = LibCameras{}
}

func (me *LibCameras) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibCameras) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
