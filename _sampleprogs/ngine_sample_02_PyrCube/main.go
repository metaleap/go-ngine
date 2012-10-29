package main

import (
	"math"

	ngine "github.com/go3d/go-ngine/core"
	ngine_samples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	floor, pyr, box *ngine.TNode
)

func main () {
	ngine_samples.SamplesMainFunc(LoadSampleScene_02_PyrCube)
}

func onLoop () {
	ngine_samples.CheckToggleKeys()
	ngine_samples.CheckCamCtlKeys()

	//	animate mesh nodes
	pyr.Transform.Rot.X -= 0.0005
	pyr.Transform.Rot.Y -= 0.0005
	pyr.Transform.Pos.Set(-13.75, 2 * math.Sin(ngine.Loop.TickNow), 2)
	pyr.Transform.OnPosRotChanged()

	box.Transform.Rot.Y += 0.0004
	box.Transform.Rot.Z += 0.0006
	box.Transform.Pos.Set(-8.125, 2 * math.Cos(ngine.Loop.TickNow), -2)
	box.Transform.OnPosRotChanged()
}

func LoadSampleScene_02_PyrCube () {
	var err error
	var meshFloor, meshPyr, meshCube *ngine.TMesh
	var bufFloor, bufRest *ngine.TMeshBuffer

	ngine.Loop.OnLoop = onLoop
	ngine.Core.Options.SetGlBackfaceCulling(false)

	//	textures / materials
	ngine.Core.Textures["tex_cobbles"] = ngine.Core.Textures.LoadAsync(ngine.TextureProviders.RemoteFile, "http://dl.dropbox.com/u/136375/misc/cobbles.png")
	ngine.Core.Textures["tex_crate"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "misc/crate.jpeg")
	ngine.Core.Textures["tex_mosaic"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "misc/mosaic.jpeg")
	ngine.Core.Materials["mat_cobbles"] = ngine.Core.Materials.New("tex_cobbles")
	ngine.Core.Materials["mat_crate"] = ngine.Core.Materials.New("tex_crate")
	ngine.Core.Materials["mat_mosaic"] = ngine.Core.Materials.New("tex_mosaic")

	//	meshes / models
	if bufFloor, err = ngine.Core.MeshBuffers.Add("buf_floor", ngine.Core.MeshBuffers.NewParams(6, 6)); err != nil { panic(err) }
	if bufRest, err = ngine.Core.MeshBuffers.Add("buf_rest", ngine.Core.MeshBuffers.NewParams(36 + 12, 36 + 12)); err != nil { panic(err) }
	if meshFloor, err = ngine.Core.Meshes.Load("mesh_plane", ngine.MeshProviders.PrefabPlane); err != nil { panic(err) }
	if meshPyr, err = ngine.Core.Meshes.Load("mesh_pyramid", ngine.MeshProviders.PrefabPyramid); err != nil { panic(err) }
	if meshCube, err = ngine.Core.Meshes.Load("mesh_cube", ngine.MeshProviders.PrefabCube); err != nil { panic(err) }
	ngine.Core.Meshes.AddRange(meshFloor, meshPyr, meshCube)
	bufFloor.Add(meshFloor); bufRest.Add(meshCube); bufRest.Add(meshPyr)
	meshPyr.Models.Default().SetMatName("mat_mosaic")
	meshCube.Models.Default().SetMatName("mat_crate")

	//	scene
	var scene = ngine.NewScene()
	ngine.Core.Scenes[""] = scene
	scene.RootNode.MakeSubNodes("node_floor", "mesh_plane", "", "node_pyr", "mesh_pyramid", "", "node_box", "mesh_cube", "")
	floor, pyr, box = scene.RootNode.SubNodes["node_floor"], scene.RootNode.SubNodes["node_pyr"], scene.RootNode.SubNodes["node_box"]

	floor.SetMatName("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScalingN(100)

	ngine_samples.CamCtl.BeginUpdate(); ngine_samples.CamCtl.Pos.Y = 1.6; ngine_samples.CamCtl.EndUpdate()

	//	upload everything
	ngine.Core.SyncUpdates()
}
