package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type Cameras map[string]*Camera

func (me Cameras) AddNew(id string) (cam *Camera) {
	cam = NewCamera()
	me[id] = cam
	return
}

type Camera struct {
	ViewPort   CameraViewPort
	Options    CameraOptions
	Controller Controller
	Disabled   bool

	technique renderTechnique
	matProj   unum.Mat4
	glmatProj ugl.GlMat4
}

func NewCamera() (me *Camera) {
	me = &Camera{}
	me.Options.init()
	me.matProj.Identity()
	me.glmatProj.Load(&me.matProj)
	me.Controller.init()
	me.ViewPort.init()
	me.UpdatePerspective()
	me.SetTechnique(Core.Options.DefaultRenderTechnique)
	return
}

func (me *Camera) dispose() {
}

func (me *Camera) render() {
	curScene = Core.Libs.Scenes[me.Options.SceneName]
	Core.useTechnique(me.technique)
	gl.UniformMatrix4fv(curProg.UnifLocs["uMatCam"], 1, gl.FALSE, &me.Controller.glMat[0])
	gl.UniformMatrix4fv(curProg.UnifLocs["uMatProj"], 1, gl.FALSE, &me.glmatProj[0])
	me.technique.onPreRender()
	gl.Viewport(me.ViewPort.glVpX, me.ViewPort.glVpY, me.ViewPort.glVpW, me.ViewPort.glVpH)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	curScene.RootNode.render()
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

func (me *Camera) UpdatePerspective() {
	me.matProj.Perspective(me.Options.FovY, me.ViewPort.aspect, me.Options.ZNear, me.Options.ZFar)
	me.glmatProj.Load(&me.matProj)
}

type CameraOptions struct {
	FovY, ZFar, ZNear float64
	SceneName         string
}

func (me *CameraOptions) init() {
	me.FovY = 37.8493
	me.ZFar = 30000
	me.ZNear = 0.3
}

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

func (me *CameraViewPort) SetAbs(x, y, width, height int) {
	me.absolute, me.absX, me.absY, me.absW, me.absH = true, x, y, width, height
	me.update()
}

func (me *CameraViewPort) SetRel(x, y, width, height float64) {
	me.absolute, me.relX, me.relY, me.relW, me.relH = false, x, y, width, height
	me.update()
}

func (me *CameraViewPort) update() {
	if !me.absolute {
		me.absW, me.absH = int(me.relW*float64(curCanvas.viewWidth)), int(me.relH*float64(curCanvas.viewHeight))
		me.absX, me.absY = int(me.relX*float64(curCanvas.viewWidth)), int(me.relY*float64(curCanvas.viewHeight))
	}
	me.glVpX, me.glVpY, me.glVpW, me.glVpH = gl.Int(me.absX), gl.Int(me.absY), gl.Sizei(me.absW), gl.Sizei(me.absH)
	me.aspect = float64(me.absW) / float64(me.absH)
}
