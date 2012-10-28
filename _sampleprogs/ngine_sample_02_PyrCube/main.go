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
	var bufFloor, bufPyr, bufCube *ngine.TMeshBuffer

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
	if bufPyr, err = ngine.Core.MeshBuffers.Add("buf_pyr", ngine.Core.MeshBuffers.NewParams(24, 24)); err != nil { panic(err) }
	if bufCube, err = ngine.Core.MeshBuffers.Add("buf_cube", ngine.Core.MeshBuffers.NewParams(36, 36)); err != nil { panic(err) }
	if meshFloor, err = ngine.Core.Meshes.Load("mesh_plane", ngine.MeshProviders.PrefabPlane); err != nil { panic(err) }
	if meshPyr, err = ngine.Core.Meshes.Load("mesh_pyramid", ngine.MeshProviders.PrefabPyramid); err != nil { panic(err) }
	if meshCube, err = ngine.Core.Meshes.Load("mesh_cube", ngine.MeshProviders.PrefabCube); err != nil { panic(err) }
	ngine.Core.Meshes.AddRange(meshFloor, meshPyr, meshCube)
	bufFloor.Add(meshFloor); bufCube.Add(meshCube); bufPyr.Add(meshPyr)

	//	scene
	var scene = ngine.NewScene()
	ngine.Core.Scenes[""] = scene
	scene.RootNode.AddSubNodesNamed(map[string]string { "floor": "mesh_plane", "pyr": "mesh_pyramid", "box": "mesh_cube" })
	floor, pyr, box = scene.RootNode.SubNodes["floor"], scene.RootNode.SubNodes["pyr"], scene.RootNode.SubNodes["box"]
	pyr.SetMatKey("mat_mosaic")
	box.SetMatKey("mat_crate")
	floor.SetMatKey("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScalingN(100)

	ngine_samples.CamCtl.BeginUpdate(); ngine_samples.CamCtl.Pos.Y = 1.6; ngine_samples.CamCtl.EndUpdate()

	//	upload everything
	ngine.Core.SyncUpdates()
}
