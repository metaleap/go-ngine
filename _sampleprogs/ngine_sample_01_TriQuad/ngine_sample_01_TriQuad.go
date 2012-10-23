package main

import (
	"math"

	gl "github.com/chsc/gogl/gl42"
	glfw "github.com/go-gl/glfw"

	numutil "github.com/go3d/go-util/num"

	ngine "github.com/go3d/go-ngine/core"
	ngine_samplescenes "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

func main() {
	ngine_samplescenes.MaxKeyHint = 1
	ngine_samplescenes.SamplesMainFunc(LoadSampleScene_01_TriQuad)
}

func LoadSampleScene_01_TriQuad() {
	ngine.Core.Materials["cat"] = ngine_samplescenes.NewMaterialFromLocalTextureImageFile("misc/cat.png")
	ngine.Core.Materials["dog"] = ngine_samplescenes.NewMaterialFromLocalTextureImageFile("misc/dog.png")
	ngine.Core.Meshes["face3"] = ngine.Core.Meshes.NewTriangle()
	ngine.Core.Meshes["face4"] = ngine.Core.Meshes.NewQuad()
	var scene = ngine.NewScene()
	ngine.Core.Options.SetGlBackfaceCulling(false)
	ngine.Core.Scenes[""] = scene
	scene.RootNode.AddSubNodesNamed(map[string]string{
		"tri":  "face3",
		"quad": "face4",
	})
	var tri, quad = scene.RootNode.SubNodes["tri"], scene.RootNode.SubNodes["quad"]
	tri.SetMatKey("cat", []gl.Float {
		0, 0,
		3, 0,
		3, 3,
	})
	quad.SetMatKey("dog", []gl.Float {
		// Front face
		1, 1,
		0, 1,
		1, 0,
		0, 0,
	})
	ngine.Core.SyncUpdates()
	ngine.Loop.OnLoopHandlers = append(ngine.Loop.OnLoopHandlers, func() {
		if ngine.Windowing.KeyToggled(glfw.KeyF2) { ngine.Core.CurCamera.ToggleTechnique() }
		if ngine.Windowing.KeyToggled(glfw.KeyF3) { ngine.Core.Options.ToggleGlBackfaceCulling() }
		tri.Transform.Rot.X -= 0.0005
		tri.Transform.Rot.Y -= 0.0005
		tri.Transform.Pos = &numutil.TVec3 { -3.75, 1 * math.Sin(ngine.Loop.TickNow), 1 }
		// tri.Transform.Scaling = &numutil.TVec3 { math.Sin(ngine.Loop.LoopNowTime) + 2, math.Cos(ngine.Loop.LoopNowTime) + 2, 1 }
		tri.Transform.OnAnyChanged()
		quad.Transform.Rot.Y += 0.0004
		quad.Transform.Rot.Z += 0.0006
		quad.Transform.Pos = &numutil.TVec3 { -8.125, 1 * math.Cos(ngine.Loop.TickNow), -2 }
		// quad.Transform.Scaling = &numutil.TVec3 { 1, 1, (math.Sin(ngine.Loop.LoopNowTime * 2) + 2) * 1.5 }
		quad.Transform.OnAnyChanged()
	})
}
