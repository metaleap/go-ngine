package main

import (
	"math"

	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
	nga "github.com/go3d/go-ngine/assets"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, pyr, box *ng.Node
)

func main() {
	ngsamples.SamplesMainFunc(LoadSampleScene_02_PyrCube)
}

func onLoop() {
	ngsamples.CheckToggleKeys()
	ngsamples.CheckCamCtlKeys()

	//	animate mesh nodes
	pyr.Transform.Rot.X -= 0.0005
	pyr.Transform.Rot.Y -= 0.0005
	pyr.Transform.Pos.Set(-13.75, 2*math.Sin(ng.Loop.TickNow), 2)
	pyr.Transform.OnPosRotChanged()

	box.Transform.Rot.Y += 0.0004
	box.Transform.Rot.Z += 0.0006
	box.Transform.Pos.Set(-8.125, 2*math.Cos(ng.Loop.TickNow), -2)
	box.Transform.OnPosRotChanged()
}

func LoadSampleScene_02_PyrCube() {
	var (
		err                          error
		scene                        *ng.Scene
		meshFloor, meshPyr, meshCube *ng.Mesh
		bufFloor, bufRest            *ng.MeshBuffer
	)

	ng.Loop.OnLoop = onLoop
	ngsamples.Cam.Options.BackfaceCulling = false

	//	textures / materials

	nga.ImageDefs.AddNew("tex_cobbles").InitFrom.RefUrl = "http://dl.dropbox.com/u/136375/go-ngine/assets/tex/cobbles.png"
	nga.ImageDefs.AddNew("tex_crate").InitFrom.RefUrl = "tex/crate.jpeg"
	nga.ImageDefs.AddNew("tex_mosaic").InitFrom.RefUrl = "tex/mosaic.jpeg"

	ng.Core.Materials["mat_cobbles"] = ng.Core.Materials.New("tex_cobbles")
	ng.Core.Materials["mat_crate"] = ng.Core.Materials.New("tex_crate")
	ng.Core.Materials["mat_mosaic"] = ng.Core.Materials.New("tex_mosaic")

	//	meshes / models
	if bufFloor, err = ng.Core.MeshBuffers.Add("buf_floor", ng.Core.MeshBuffers.NewParams(6, 6)); err != nil {
		panic(err)
	}
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(36+12, 36+12)); err != nil {
		panic(err)
	}
	if meshFloor, err = ng.Core.Meshes.Load("mesh_plane", ng.MeshProviders.PrefabPlane); err != nil {
		panic(err)
	}
	if meshPyr, err = ng.Core.Meshes.Load("mesh_pyramid", ng.MeshProviders.PrefabPyramid); err != nil {
		panic(err)
	}
	if meshCube, err = ng.Core.Meshes.Load("mesh_cube", ng.MeshProviders.PrefabCube); err != nil {
		panic(err)
	}
	ng.Core.Meshes.AddRange(meshFloor, meshPyr, meshCube)
	bufFloor.Add(meshFloor)
	bufRest.Add(meshCube)
	bufRest.Add(meshPyr)
	meshPyr.Models.Default().SetMatName("mat_mosaic")
	meshCube.Models.Default().SetMatName("mat_crate")

	//	scene
	scene = ng.NewScene()
	ng.Core.Scenes[""] = scene
	scene.RootNode.SubNodes.MakeN("node_floor", "mesh_plane", "", "node_pyr", "mesh_pyramid", "", "node_box", "mesh_cube", "")
	floor, pyr, box = scene.RootNode.SubNodes.M["node_floor"], scene.RootNode.SubNodes.M["node_pyr"], scene.RootNode.SubNodes.M["node_box"]

	floor.SetMatName("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScalingN(1000)

	ngsamples.CamCtl.BeginUpdate()
	ngsamples.CamCtl.Pos.Y = 1.6
	ngsamples.CamCtl.EndUpdate()

	//	upload everything
	ng.Core.SyncUpdates()
}
