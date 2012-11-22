package main

import (
	ng "github.com/go3d/go-ngine/core"
	nga "github.com/go3d/go-ngine/assets"
	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	floor *ng.Node
)

func main () {
	ngsamples.SamplesMainFunc(LoadSampleScene_01_EmptyPlane)
}

func onLoop () {
	ngsamples.CheckToggleKeys()
	ngsamples.CheckCamCtlKeys()
}

func LoadSampleScene_01_EmptyPlane () {
	var err error
	var meshFloor *ng.Mesh
	var bufRest *ng.MeshBuffer

	ng.Loop.OnLoop = onLoop
	ngsamples.Cam.Options.BackfaceCulling = false

	//	textures / materials
	nga.ImageDefs.AddNew("tex_cobbles").InitFrom.RefUrl = "tex/cobbles.png"
	ng.Core.Materials["mat_cobbles"] = ng.Core.Materials.New("tex_cobbles")

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(100, 100)); err != nil { panic(err) }
	if meshFloor, err = ng.Core.Meshes.Load("mesh_plane", ng.MeshProviders.PrefabPlane); err != nil { panic(err) }

	ng.Core.Meshes.AddRange(meshFloor)
	bufRest.Add(meshFloor); // bufRest.Add(meshCube);

	//	scene
	var scene = ng.NewScene()
	ng.Core.Scenes[""] = scene
	scene.RootNode.SubNodes.MakeN("node_floor", "mesh_plane", "" /*"node_box", "mesh_cube", ""*/)
	floor = scene.RootNode.SubNodes.M["node_floor"]

	floor.SetMatName("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScalingN(100)

	ngsamples.CamCtl.BeginUpdate(); ngsamples.CamCtl.Pos.Y = 1.6; ngsamples.CamCtl.EndUpdate()

	//	upload everything
	ng.Core.SyncUpdates()
}
