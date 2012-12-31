package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
	unum "github.com/metaleap/go-util/num"
)

type cameras map[string]*Camera

func (me cameras) AddNew(id string) (cam *Camera) {
	cam = me.New()
	me[id] = cam
	return
}

func (me cameras) New() *Camera {
	return newCamera()
}

type Camera struct {
	ViewPort   *CameraViewPort
	MatProj    *unum.Mat4
	Options    *CameraOptions
	Controller *Controller
	Disabled   bool
	SceneName  string

	technique renderTechnique
	glMatProj *ugl.GlMat4
}

func newCamera() (me *Camera) {
	me = &Camera{}
	me.Options = newCameraOptions()
	me.MatProj = &unum.Mat4{}
	me.glMatProj = &ugl.GlMat4{}
	me.Controller = newController()
	me.ViewPort = newCameraViewPort(me)
	me.UpdatePerspective()
	me.SetTechnique(Core.Options.DefaultRenderTechnique)
	return
}

func (me *Camera) dispose() {
}

func (me *Camera) render() {
	curScene = Core.Scenes[me.SceneName]
	glSetBackfaceCulling(me.Options.BackfaceCulling)
	Core.useTechnique(me.technique)
	gl.UniformMatrix4fv(curProg.UnifLocs["uMatCam"], 1, gl.FALSE, &me.Controller.glMat[0])
	gl.UniformMatrix4fv(curProg.UnifLocs["uMatProj"], 1, gl.FALSE, &me.glMatProj[0])
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
	allNames, curTech, name := glShaderMan.AllNames, curTechnique.name(), ""
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
	me.MatProj.Perspective(me.Options.FovY, me.ViewPort.aspect, me.Options.ZNear, me.Options.ZFar)
	me.glMatProj.Load(me.MatProj)
}

type CameraViewPort struct {
	Absolute               bool
	RelX, RelY, RelW, RelH float64
	AbsX, AbsY, AbsW, AbsH int
	aspect                 float64
	cam                    *Camera
	glVpX, glVpY           gl.Int
	glVpW, glVpH           gl.Sizei
}

func newCameraViewPort(cam *Camera) (me *CameraViewPort) {
	me = &CameraViewPort{cam: cam}
	me.SetRel(0, 0, 1, 1)
	return
}

func (me *CameraViewPort) SetAbs(x, y, width, height int) {
	me.Absolute, me.AbsX, me.AbsY, me.AbsW, me.AbsH = true, x, y, width, height
	me.Update()
}

func (me *CameraViewPort) SetRel(x, y, width, height float64) {
	me.Absolute, me.RelX, me.RelY, me.RelW, me.RelH = false, x, y, width, height
	me.Update()
}

func (me *CameraViewPort) Update() {
	if !me.Absolute {
		me.AbsW, me.AbsH = int(me.RelW*float64(curCanvas.viewWidth)), int(me.RelH*float64(curCanvas.viewHeight))
		me.AbsX, me.AbsY = int(me.RelX*float64(curCanvas.viewWidth)), int(me.RelY*float64(curCanvas.viewHeight))
	}
	me.glVpX, me.glVpY, me.glVpW, me.glVpH = gl.Int(me.AbsX), gl.Int(me.AbsY), gl.Sizei(me.AbsW), gl.Sizei(me.AbsH)
	me.aspect = float64(me.AbsW) / float64(me.AbsH)
}
