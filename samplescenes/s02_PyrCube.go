package samplescenes

import (
	"math"
	"path/filepath"

	gl "github.com/chsc/gogl/gl42"
	glfw "github.com/jteeuwen/glfw"

	numutil "github.com/go-ngine/go-util/num"

	ngine "github.com/go-ngine/go-ngine/client"
	ncore "github.com/go-ngine/go-ngine/client/core"
)

func LoadSampleScene_02_PyrCube() {
	ncore.Core.Materials["cobbles"] = ncore.NewMaterialFromLocalTextureImageFile(filepath.Join(ngine.AssetRootDirPath, "misc", "cobbles.png"))
	ncore.Core.Materials["crate"] = ncore.NewMaterialFromLocalTextureImageFile(filepath.Join(ngine.AssetRootDirPath, "misc", "crate.jpg"))
	ncore.Core.Materials["mosaic"] = ncore.NewMaterialFromLocalTextureImageFile(filepath.Join(ngine.AssetRootDirPath, "misc", "mosaic.jpg"))
	ncore.Core.Meshes["face3"] = NewMeshPyramid()
	ncore.Core.Meshes["face4"] = NewMeshCube()
	ncore.Core.Meshes["plane"] = NewMeshPlane()
	ngine.Core.Options.SetGlBackfaceCulling(false)
	var scene = ncore.NewScene()
	ncore.Core.Scenes[""] = scene
	scene.RootNode.AddSubNodesNamed(map[string]string{
		"floor": "plane",
		"tri":   "face3",
		"quad":  "face4",
	})
	var floor, tri, quad = scene.RootNode.SubNodes["floor"], scene.RootNode.SubNodes["tri"], scene.RootNode.SubNodes["quad"]
	floor.SetMatKey("cobbles", []gl.Float{
		0, 0,
		10, 0,
		0, 10,
		10, 10,
	})
	tri.SetMatKey("mosaic", []gl.Float{
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
	quad.SetMatKey("crate", []gl.Float{
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
	floor.Transform.SetPos(&numutil.TVec3{0.1, -1.75, -8})
	floor.Transform.SetScalingN(100)
	ncore.Core.SyncUpdates()
	ngine.Loop.OnLoopHandlers = append(ngine.Loop.OnLoopHandlers, func() {
		ngine.Core.CurCamera.Controller.MoveSpeedupFactor = 1
		if ngine.Windowing.KeyToggled(glfw.KeyF2) { ngine.Core.CurCamera.ToggleTechnique() }
		if ngine.Windowing.KeyToggled(glfw.KeyF3) { ngine.Core.Options.ToggleGlBackfaceCulling() }
		if ngine.Windowing.KeyPressed(glfw.KeyLshift) {
			ngine.Core.CurCamera.Controller.MoveSpeedupFactor = 10
		} else if ngine.Windowing.KeyPressed(glfw.KeyRshift) {
			ngine.Core.CurCamera.Controller.MoveSpeedupFactor = 100
		} else if ngine.Windowing.KeyPressed(glfw.KeyLalt) {
			ngine.Core.CurCamera.Controller.MoveSpeedupFactor = 0.1
		}
		if ngine.Windowing.KeyPressed(glfw.KeyUp) { ngine.Core.CurCamera.Controller.MoveForward() }
		if ngine.Windowing.KeyPressed(glfw.KeyDown) { ngine.Core.CurCamera.Controller.MoveBackward() }
		if ngine.Windowing.KeyPressed('A') { ngine.Core.CurCamera.Controller.MoveLeft() }
		if ngine.Windowing.KeyPressed('D') { ngine.Core.CurCamera.Controller.MoveRight() }
		if ngine.Windowing.KeyPressed('W') { ngine.Core.CurCamera.Controller.MoveUp() }
		if ngine.Windowing.KeyPressed('S') { ngine.Core.CurCamera.Controller.MoveDown() }
		if ngine.Windowing.KeyPressed(glfw.KeyLeft) { ngine.Core.CurCamera.Controller.TurnLeft() }
		if ngine.Windowing.KeyPressed(glfw.KeyRight) { ngine.Core.CurCamera.Controller.TurnRight() }
		if ngine.Windowing.KeysPressedAny2(glfw.KeyPageup, glfw.KeyKP9) { ngine.Core.CurCamera.Controller.TurnUp() }
		if ngine.Windowing.KeysPressedAny2(glfw.KeyPagedown, glfw.KeyKP3) { ngine.Core.CurCamera.Controller.TurnDown() }
		tri.Transform.Rot.X -= 0.0005
		tri.Transform.Rot.Y -= 0.0005
		tri.Transform.Pos = &numutil.TVec3{-13.75, 2 * math.Sin(ngine.Loop.TickNow), 2}
		tri.Transform.OnAnyChanged()
		quad.Transform.Rot.Y += 0.0004
		quad.Transform.Rot.Z += 0.0006
		quad.Transform.Pos = &numutil.TVec3{-8.125, 2 * math.Cos(ngine.Loop.TickNow), -2}
		quad.Transform.OnAnyChanged()
	})
}
