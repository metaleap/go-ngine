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

func LoadSampleScene_01_TriQuad() {
	ngine.Core.Materials["cat"] = ncore.NewMaterialFromLocalTextureImageFile(filepath.Join(ngine.AssetRootDirPath, "misc", "cat.png"))
	ngine.Core.Materials["dog"] = ncore.NewMaterialFromLocalTextureImageFile(filepath.Join(ngine.AssetRootDirPath, "misc", "dog.png"))
	ngine.Core.Meshes["face3"] = NewMeshFace3()
	ngine.Core.Meshes["face4"] = NewMeshFace4()
	var scene = ncore.NewScene()
	ngine.Core.Options.SetGlBackfaceCulling(false)
	ngine.Core.Scenes[""] = scene
	scene.RootNode.AddSubNodesNamed(map[string]string{
		"tri":  "face3",
		"quad": "face4",
	})
	var tri, quad = scene.RootNode.SubNodes["tri"], scene.RootNode.SubNodes["quad"]
	tri.SetMatKey("cat", []gl.Float{
		0, 0,
		3, 0,
		3, 3,
	})
	quad.SetMatKey("dog", []gl.Float{
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
		tri.Transform.Pos = &numutil.TVec3{-3.75, 1 * math.Sin(ngine.Loop.TickNow), 1}
		// tri.Transform.Scaling = &numutil.TVec3 { math.Sin(ngine.Loop.LoopNowTime) + 2, math.Cos(ngine.Loop.LoopNowTime) + 2, 1 }
		tri.Transform.OnAnyChanged()
		quad.Transform.Rot.Y += 0.0004
		quad.Transform.Rot.Z += 0.0006
		quad.Transform.Pos = &numutil.TVec3{-8.125, 1 * math.Cos(ngine.Loop.TickNow), -2}
		// quad.Transform.Scaling = &numutil.TVec3 { 1, 1, (math.Sin(ngine.Loop.LoopNowTime * 2) + 2) * 1.5 }
		quad.Transform.OnAnyChanged()
	})
}
