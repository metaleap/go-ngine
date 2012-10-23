package main

import (
	"math"

	gl "github.com/chsc/gogl/gl42"
	glfw "github.com/go-gl/glfw"

	ngine "github.com/go3d/go-ngine/core"
	ngine_samplescenes "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	floor, tri, quad *ngine.TNode
	camCtl *ngine.TController
)

func main () {
	ngine_samplescenes.SamplesMainFunc(LoadSampleScene_02_PyrCube)
}

func onLoop () {
	//	check option-toggle keys
	if ngine.Windowing.KeyToggled(glfw.KeyF2) { ngine.Core.CurCamera.ToggleTechnique() }
	if ngine.Windowing.KeyToggled(glfw.KeyF3) { ngine.Core.Options.ToggleGlBackfaceCulling() }
	if ngine.Windowing.KeyToggled(glfw.KeyF4) { ngine.Core.Options.DefaultTextureParams.ToggleFilter() }
	if ngine.Windowing.KeyToggled(glfw.KeyF5) { ngine.Core.Options.DefaultTextureParams.ToggleFilterAnisotropy() }

	//	check camera-control keys
	camCtl = ngine.Core.CurCamera.Controller
	camCtl.MoveSpeedupFactor = 1
	if ngine.Windowing.KeyPressed(glfw.KeyLshift) {
		camCtl.MoveSpeedupFactor = 10
	} else if ngine.Windowing.KeyPressed(glfw.KeyRshift) {
		camCtl.MoveSpeedupFactor = 100
	} else if ngine.Windowing.KeyPressed(glfw.KeyLalt) {
		camCtl.MoveSpeedupFactor = 0.1
	}
	if ngine.Windowing.KeyPressed(glfw.KeyUp) { camCtl.MoveForward() }
	if ngine.Windowing.KeyPressed(glfw.KeyDown) { camCtl.MoveBackward() }
	if ngine.Windowing.KeyPressed('A') { camCtl.MoveLeft() }
	if ngine.Windowing.KeyPressed('D') { camCtl.MoveRight() }
	if ngine.Windowing.KeyPressed('W') { camCtl.MoveUp() }
	if ngine.Windowing.KeyPressed('S') { camCtl.MoveDown() }
	if ngine.Windowing.KeyPressed(glfw.KeyLeft) { camCtl.TurnLeft() }
	if ngine.Windowing.KeyPressed(glfw.KeyRight) { camCtl.TurnRight() }
	if ngine.Windowing.KeysPressedAny2(glfw.KeyPageup, glfw.KeyKP9) { camCtl.TurnUp() }
	if ngine.Windowing.KeysPressedAny2(glfw.KeyPagedown, glfw.KeyKP3) { camCtl.TurnDown() }

	//	animate mesh nodes
	tri.Transform.Rot.X -= 0.0005
	tri.Transform.Rot.Y -= 0.0005
	tri.Transform.Pos.Set(-13.75, 2 * math.Sin(ngine.Loop.TickNow), 2)
	tri.Transform.OnAnyChanged()

	quad.Transform.Rot.Y += 0.0004
	quad.Transform.Rot.Z += 0.0006
	quad.Transform.Pos.Set(-8.125, 2 * math.Cos(ngine.Loop.TickNow), -2)
	quad.Transform.OnAnyChanged()
}

func LoadSampleScene_02_PyrCube () {
	ngine.Core.Options.SetGlBackfaceCulling(false)

	ngine.Core.Materials["cobbles"] = ngine_samplescenes.NewMaterialFromRemoteTextureImageFile("http://dl.dropbox.com/u/136375/misc/cobbles.png")
	ngine.Core.Materials["crate"] = ngine_samplescenes.NewMaterialFromLocalTextureImageFile("misc/crate.jpeg")
	ngine.Core.Materials["mosaic"] = ngine_samplescenes.NewMaterialFromLocalTextureImageFile("misc/mosaic.jpeg")

	ngine.Core.Meshes["face3"] = ngine.Core.Meshes.NewPyramid()
	ngine.Core.Meshes["face4"] = ngine.Core.Meshes.NewCube()
	ngine.Core.Meshes["plane"] = ngine.Core.Meshes.NewPlane()

	var scene = ngine.NewScene()
	ngine.Core.Scenes[""] = scene
	scene.RootNode.AddSubNodesNamed(map[string]string { "floor": "plane", "tri":   "face3", "quad":  "face4" })
	floor, tri, quad = scene.RootNode.SubNodes["floor"], scene.RootNode.SubNodes["tri"], scene.RootNode.SubNodes["quad"]

	floor.SetMatKey("cobbles", []gl.Float {
		0, 0,
		10, 0,
		0, 10,
		10, 10,
	})
	tri.SetMatKey("mosaic", []gl.Float {
		// Front face
		0, 0,
		1, 0,
		1, 1,
		// Right face
		1, 0,
		1, 1,
		0, 1,
		// Back face
		1, 0,
		1, 1,
		0, 1,
		// Left face
		0, 0,
		1, 0,
		1, 1,
	})
	quad.SetMatKey("crate", []gl.Float {
		// Front face
		0, 0,
		1, 0,
		1, 1,
		0, 1,
		// Back face
		1, 0,
		1, 1,
		0, 1,
		0, 0,
		// Top face
		0, 1,
		0, 0,
		1, 0,
		1, 1,
		// Bottom face
		1, 1,
		0, 1,
		0, 0,
		1, 0,
		// Right face
		1, 0,
		1, 1,
		0, 1,
		0, 0,
		// Left face
		0, 0,
		1, 0,
		1, 1,
		0, 1,
	})

	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScalingN(100)
	ngine.Core.SyncUpdates()
	ngine.Core.CurCamera.Controller.With(func (ctl *ngine.TController) { ctl.Pos.Y = 1.6 })
	ngine.Loop.AddHandler(onLoop)
}
