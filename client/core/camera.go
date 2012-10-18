package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/metaleap/go-util/gl"
	numutil "github.com/metaleap/go-util/num"

	nglcore "github.com/metaleap/go-ngine/client/glcore"
)

type TCamera struct {
	Controller *TController
	Disabled bool
	RenderSceneKey string
	ViewPort *TViewPort

	technique ITechnique
	canvas *TCanvas
	nearPlane, farPlane, fieldOfView float64
	matProj *numutil.TMat4
	glMatProj *glutil.TGlMat4
}

	func NewCamera (parentCanvas *TCanvas, technique string) *TCamera {
		var cam = &TCamera {}
		cam.SetTechnique(technique)
		cam.canvas = parentCanvas
		cam.matProj = &numutil.TMat4 {}
		cam.glMatProj = &glutil.TGlMat4 {}
		cam.ViewPort = newViewPort(cam)
		cam.SetPerspective(0.3, 30000, 45)
		cam.Controller = NewController()
		return cam
	}

	func (me *TCamera) Dispose () {
	}

	func (me *TCamera) FarPlane () float64 {
		return me.farPlane
	}

	func (me *TCamera) FieldOfView () float64 {
		return me.fieldOfView
	}

	func (me *TCamera) NearPlane () float64 {
		return me.nearPlane
	}

	func (me *TCamera) render () {
		Core.useTechnique(me.technique)
		gl.UniformMatrix4fv(curProg.UnifLocs["uMatCam"], 1, gl.FALSE, &me.Controller.glMat[0])
		gl.UniformMatrix4fv(curProg.UnifLocs["uMatProj"], 1, gl.FALSE, &me.glMatProj[0])
		me.technique.OnPreRender()

		gl.Viewport(me.ViewPort.glX, me.ViewPort.glY, me.ViewPort.glW, me.ViewPort.glH)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		curScene.RootNode.render()
	}

	func (me *TCamera) SetPerspective (nearPlane, farPlane, fieldOfView float64) {
		me.nearPlane, me.farPlane, me.fieldOfView = nearPlane, farPlane, fieldOfView
		me.updatePerspective()
	}

	func (me *TCamera) SetTechnique (name string) {
		if (me.technique == nil) || (me.technique.Name() != name) {
			me.technique = GetTechnique(name)
		}
	}

	func (me *TCamera) ToggleTechnique () {
		var allNames, curTech, name = nglcore.ShaderMan.AllNames, curTechnique.Name(), ""
		var curIndex, i int
		var tech ITechnique = nil
		for i, name = range allNames { if name == curTech { curIndex = i; break } }
		if curIndex < (len(allNames) - 1) { for i = curIndex + 1; i < len(allNames); i++ { if tech = GetTechnique(allNames[i]); tech != nil { break } } }
		if tech == nil { for i = 0; i < curIndex; i++ { if tech = GetTechnique(allNames[i]); tech != nil { break } } }
		if tech != nil { me.technique = tech }
	}

	func (me *TCamera) updatePerspective () {
		me.matProj.Perspective(me.fieldOfView, me.ViewPort.aspect, me.nearPlane, me.farPlane)
		me.glMatProj.Load(me.matProj)
	}

type TViewPort struct {
	absolute bool
	camera *TCamera
	glX, glY gl.Int
	glW, glH gl.Sizei
	relX, relY, relW, relH float64
	absX, absY, absW, absH int
	aspect float64
}

	func newViewPort (cam *TCamera) *TViewPort {
		var vp = &TViewPort {}
		vp.camera = cam
		vp.SetRel(0, 0, 1, 1)
		return vp
	}

	func (me *TViewPort) SetAbs (x, y, width, height int) {
		me.absolute, me.absX, me.absY, me.absW, me.absH = true, x, y, width, height
		me.update()
	}

	func (me *TViewPort) SetRel (x, y, width, height float64) {
		me.absolute, me.relX, me.relY, me.relW, me.relH = false, x, y, width, height
		me.update()
	}

	func (me *TViewPort) update () {
		if !me.absolute {
			me.absW, me.absH = int(me.relW * float64(me.camera.canvas.viewWidth)), int(me.relH * float64(me.camera.canvas.viewHeight))
			me.absX, me.absY = int(me.relX * float64(me.camera.canvas.viewWidth)), int(me.relY * float64(me.camera.canvas.viewHeight))
		}
		me.glX, me.glY, me.glW, me.glH = gl.Int(me.absX), gl.Int(me.absY), gl.Sizei(me.absW), gl.Sizei(me.absH)
		me.aspect = float64(me.absW) / float64(me.absH)
	}
