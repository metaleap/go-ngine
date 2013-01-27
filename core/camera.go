package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type Camera interface {
	render()
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
