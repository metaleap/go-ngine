package main

import (
	ngine "github.com/go3d/go-ngine/core"
	ngine_samples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	floor *ngine.Node
)

func main () {
	ngine_samples.SamplesMainFunc(LoadSampleScene_04_ColladaCube)
}

func onLoop () {
	ngine_samples.CheckToggleKeys()
	ngine_samples.CheckCamCtlKeys()
}

func LoadSampleScene_04_ColladaCube () {
	var err error
	var meshFloor *ngine.Mesh
	var bufRest *ngine.MeshBuffer

	ngine.Loop.OnLoop = onLoop
	ngine.Core.Canvases[0].Cameras[0].Options.BackfaceCulling = false

	//	textures / materials
	ngine.Core.Textures["tex_cobbles"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "tex/cobbles.png")
	ngine.Core.Materials["mat_cobbles"] = ngine.Core.Materials.New("tex_cobbles")

	//	meshes / models
	if bufRest, err = ngine.Core.MeshBuffers.Add("buf_rest", ngine.Core.MeshBuffers.NewParams(100, 100)); err != nil { panic(err) }
	if meshFloor, err = ngine.Core.Meshes.Load("mesh_plane", ngine.MeshProviders.PrefabPlane); err != nil { panic(err) }

	ngine.Core.Meshes.AddRange(meshFloor)
	bufRest.Add(meshFloor); // bufRest.Add(meshCube);

	//	scene
	var scene = ngine.NewScene()
	ngine.Core.Scenes[""] = scene
	scene.RootNode.SubNodes.MakeN("node_floor", "mesh_plane", "" /*"node_box", "mesh_cube", ""*/)
	floor = scene.RootNode.SubNodes.M["node_floor"]

	floor.SetMatName("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScalingN(100)

	ngine_samples.CamCtl.BeginUpdate(); ngine_samples.CamCtl.Pos.Y = 1.6; ngine_samples.CamCtl.EndUpdate()

	//	upload everything
	ngine.Core.SyncUpdates()
}
