package main

import (
	"math"

	gl "github.com/chsc/gogl/gl42"

	ngine "github.com/go3d/go-ngine/core"
	ngine_samples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	floor, tri, quad *ngine.TNode
	cam *ngine.TCamera
	camCtl *ngine.TController
)

func main () {
	ngine_samples.SamplesMainFunc(LoadSampleScene_02_PyrCube)
}

func onLoop () {
	ngine_samples.CheckToggleKeys()
	ngine_samples.CheckCamCtlKeys()

	//	animate mesh nodes
	tri.Transform.Rot.X -= 0.0005
	tri.Transform.Rot.Y -= 0.0005
	tri.Transform.Pos.Set(-13.75, 2 * math.Sin(ngine.Loop.TickNow), 2)
	tri.Transform.OnPosRotChanged()

	quad.Transform.Rot.Y += 0.0004
	quad.Transform.Rot.Z += 0.0006
	quad.Transform.Pos.Set(-8.125, 2 * math.Cos(ngine.Loop.TickNow), -2)
	quad.Transform.OnPosRotChanged()
}

func LoadSampleScene_02_PyrCube () {
	ngine.Core.Options.SetGlBackfaceCulling(false)

	ngine.Core.Materials["cobbles"] = ngine_samples.NewMaterialFromRemoteTextureImageFile("http://dl.dropbox.com/u/136375/misc/cobbles.png")
	ngine.Core.Materials["crate"] = ngine_samples.NewMaterialFromLocalTextureImageFile("misc/crate.jpeg")
	ngine.Core.Materials["mosaic"] = ngine_samples.NewMaterialFromLocalTextureImageFile("misc/mosaic.jpeg")

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

	cam, camCtl = ngine_samples.Cam, ngine_samples.CamCtl
	camCtl.BeginUpdate()
	camCtl.Pos.Y = 1.6
	camCtl.EndUpdate()

	ngine.Core.SyncUpdates()
	ngine.Loop.AddHandler(onLoop)
}
