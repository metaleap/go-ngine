package samplescenes

import (
	ng "github.com/go3d/go-ngine/core"
	unum "github.com/metaleap/go-util/num"
)

type RearMirror struct {
	Cam        *ng.Camera
	CamRotAxis unum.Vec3
}

func (me *RearMirror) OnApp() {
	if !Paused {
		me.Cam.Controller.CopyFrom(&Cam.Controller)
		me.Cam.Controller.TurnRightBy(180)
	}
}

func (me *RearMirror) Setup() {
	me.CamRotAxis.Y = 1
	me.Cam = ng.Core.Rendering.Canvases.Main().AddNewCamera3D()
	me.Cam.Rendering.States.ClearColor.Set(0.125, 0.25, 0.45, 1)
	me.Cam.Rendering.States.FaceCulling = false
	me.Cam.Rendering.Viewport.SetRel(0.66, 0.73, 0.33, 0.25)
	me.Cam.Perspective.FovY *= 2
}
