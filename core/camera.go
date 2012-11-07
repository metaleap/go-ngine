package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type Camera struct {
	Controller *Controller
	Disabled bool
	Options *cameraOptions
	SceneName string
	ViewPort *camViewPort

	technique renderTechnique
	canvas *RenderCanvas
	nearPlane, farPlane, fieldOfView float64
	matProj *unum.Mat4
	glMatProj *ugl.GlMat4
}

	func NewCamera (parentCanvas *RenderCanvas, technique string) *Camera {
		var cam = &Camera {}
		cam.Options = newCameraOptions(cam)
		cam.SetTechnique(technique)
		cam.canvas = parentCanvas
		cam.matProj = &unum.Mat4 {}
		cam.glMatProj = &ugl.GlMat4 {}
		cam.ViewPort = newViewPort(cam)
		cam.SetPerspective(0.3, 30000, 45)
		cam.Controller = newController()
		return cam
	}

	func (me *Camera) Dispose () {
	}

	func (me *Camera) FarPlane () float64 {
		return me.farPlane
	}

	func (me *Camera) FieldOfView () float64 {
		return me.fieldOfView
	}

	func (me *Camera) NearPlane () float64 {
		return me.nearPlane
	}

	func (me *Camera) render () {
		curScene = Core.Scenes[me.SceneName]
		glSetBackfaceCulling(me.Options.BackfaceCulling)
		Core.useTechnique(me.technique)
		gl.UniformMatrix4fv(curProg.UnifLocs["uMatCam"], 1, gl.FALSE, &me.Controller.glMat[0])
		gl.UniformMatrix4fv(curProg.UnifLocs["uMatProj"], 1, gl.FALSE, &me.glMatProj[0])
		me.technique.onPreRender()
		gl.Viewport(me.ViewPort.glX, me.ViewPort.glY, me.ViewPort.glW, me.ViewPort.glH)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		curScene.RootNode.render()
	}

	func (me *Camera) SetPerspective (nearPlane, farPlane, fieldOfView float64) {
		me.nearPlane, me.farPlane, me.fieldOfView = nearPlane, farPlane, fieldOfView
		me.updatePerspective()
	}

	func (me *Camera) SetTechnique (name string) {
		if (me.technique == nil) || (me.technique.name() != name) {
			me.technique = techs[name]
		}
	}

	func (me *Camera) ToggleTechnique () {
		var allNames, curTech, name = glShaderMan.AllNames, curTechnique.name(), ""
		var curIndex, i int
		var tech renderTechnique = nil
		for i, name = range allNames { if name == curTech { curIndex = i; break } }
		if curIndex < (len(allNames) - 1) { for i = curIndex + 1; i < len(allNames); i++ { if tech = techs[allNames[i]]; tech != nil { break } } }
		if tech == nil { for i = 0; i < curIndex; i++ { if tech = techs[allNames[i]]; tech != nil { break } } }
		if tech != nil { me.technique = tech }
	}

	func (me *Camera) updatePerspective () {
		me.matProj.Perspective(me.fieldOfView, me.ViewPort.aspect, me.nearPlane, me.farPlane)
		me.glMatProj.Load(me.matProj)
	}

type camViewPort struct {
	absolute bool
	camera *Camera
	glX, glY gl.Int
	glW, glH gl.Sizei
	relX, relY, relW, relH float64
	absX, absY, absW, absH int
	aspect float64
}

	func newViewPort (cam *Camera) *camViewPort {
		var vp = &camViewPort {}
		vp.camera = cam
		vp.SetRel(0, 0, 1, 1)
		return vp
	}

	func (me *camViewPort) SetAbs (x, y, width, height int) {
		me.absolute, me.absX, me.absY, me.absW, me.absH = true, x, y, width, height
		me.update()
	}

	func (me *camViewPort) SetRel (x, y, width, height float64) {
		me.absolute, me.relX, me.relY, me.relW, me.relH = false, x, y, width, height
		me.update()
	}

	func (me *camViewPort) update () {
		if !me.absolute {
			me.absW, me.absH = int(me.relW * float64(me.camera.canvas.viewWidth)), int(me.relH * float64(me.camera.canvas.viewHeight))
			me.absX, me.absY = int(me.relX * float64(me.camera.canvas.viewWidth)), int(me.relY * float64(me.camera.canvas.viewHeight))
		}
		me.glX, me.glY, me.glW, me.glH = gl.Int(me.absX), gl.Int(me.absY), gl.Sizei(me.absW), gl.Sizei(me.absH)
		me.aspect = float64(me.absW) / float64(me.absH)
	}
