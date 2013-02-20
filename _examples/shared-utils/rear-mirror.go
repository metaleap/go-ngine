package exampleutils

import (
	ng "github.com/go3d/go-ngine/core"
)

var RearView RearMirror

//	A rather simple "rear-view mirror" camera that can be added to the example-program's main render canvas.
type RearMirror struct {
	Cam *ng.Camera
}

//	Copies the main camera's current position and direction and reverses its by 180 degrees to achieve a rear-view mirror.
func (me *RearMirror) OnApp() {
	if me.Cam.Enabled && !Paused {
		me.Cam.Controller.CopyFrom(&SceneCam.Controller)
		me.Cam.Controller.TurnRightBy(180)
	}
}

//	Syncs the rear-view camera's render states with the main camera's.
func (me *RearMirror) OnWin() {
	me.Cam.Rendering.States.FaceCulling = SceneCam.Rendering.States.FaceCulling
}

//	Adds the rear-view mirror's camera to the main render canvas, at 1/3rd of its width and 1/4th of its height.
func (me *RearMirror) Setup() {
	me.Cam = SceneCanvas.AddNewCamera3D()
	me.Cam.Rendering.States.ClearColor.Set(0.125, 0.25, 0.45, 1)
	me.Cam.Rendering.Viewport.SetRel(0.66, 0.73, 0.33, 0.25)
	me.Cam.Perspective.FovY *= 2
}

func (me *RearMirror) Toggle() {
	me.Cam.Enabled = !me.Cam.Enabled
}
