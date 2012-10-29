package main

import (
	"math"

	numutil "github.com/go3d/go-util/num"

	ngine "github.com/go3d/go-ngine/core"
	ngine_samples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	floor, pyr, box *ngine.TNode
	crates [3]*ngine.TNode
	pyramids [4]*ngine.TNode
	i int
	f float64
)

func main () {
	ngine_samples.SamplesMainFunc(LoadSampleScene_03_PyrsCubes)
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

	for i = 0; i < len(crates); i++ {
		f = float64(i)
		f = (f + 1) * (f + 1)
		crates[i].Transform.Rot.X += f * 0.00001
		crates[i].Transform.Rot.Y += f * 0.0001
		crates[i].Transform.Rot.Z += f * 0.001
		crates[i].Transform.OnRotChanged()
	}
}

func LoadSampleScene_03_PyrsCubes () {
	var err error
	var meshFloor, meshPyr, meshCube *ngine.TMesh
	var bufFloor, bufRest *ngine.TMeshBuffer
	var str string

	ngine.Loop.OnLoop = onLoop
	ngine.Core.Options.SetGlBackfaceCulling(false)

	//	textures / materials
	ngine.Core.Textures["tex_cobbles"] = ngine.Core.Textures.LoadAsync(ngine.TextureProviders.RemoteFile, "http://dl.dropbox.com/u/136375/misc/cobbles.png")
	ngine.Core.Textures["tex_crate"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "misc/crate.jpeg")
	ngine.Core.Textures["tex_mosaic"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "misc/mosaic.jpeg")
	ngine.Core.Textures["tex_cat"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "misc/cat.png")
	ngine.Core.Textures["tex_dog"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "misc/dog.png")
	ngine.Core.Materials["mat_cobbles"] = ngine.Core.Materials.New("tex_cobbles")
	ngine.Core.Materials["mat_crate"] = ngine.Core.Materials.New("tex_crate")
	ngine.Core.Materials["mat_mosaic"] = ngine.Core.Materials.New("tex_mosaic")
	ngine.Core.Materials["mat_cat"] = ngine.Core.Materials.New("tex_cat")
	ngine.Core.Materials["mat_dog"] = ngine.Core.Materials.New("tex_dog")

	//	meshes / models
	if bufFloor, err = ngine.Core.MeshBuffers.Add("buf_floor", ngine.Core.MeshBuffers.NewParams(6, 6)); err != nil { panic(err) }
	if bufRest, err = ngine.Core.MeshBuffers.Add("buf_rest", ngine.Core.MeshBuffers.NewParams(36 + 12, 36 + 12)); err != nil { panic(err) }
	if meshFloor, err = ngine.Core.Meshes.Load("mesh_plane", ngine.MeshProviders.PrefabPlane); err != nil { panic(err) }
	if meshPyr, err = ngine.Core.Meshes.Load("mesh_pyramid", ngine.MeshProviders.PrefabPyramid); err != nil { panic(err) }
	if meshCube, err = ngine.Core.Meshes.Load("mesh_cube", ngine.MeshProviders.PrefabCube); err != nil { panic(err) }
	ngine.Core.Meshes.AddRange(meshFloor, meshPyr, meshCube)
	bufFloor.Add(meshFloor); bufRest.Add(meshCube); bufRest.Add(meshPyr)

	//	scene
	var scene = ngine.NewScene()
	ngine.Core.Scenes[""] = scene
	scene.RootNode.MakeSubNodes("node_floor", "mesh_plane", "", "node_pyr", "mesh_pyramid", "", "node_box", "mesh_cube", "")
	floor, pyr, box = scene.RootNode.SubNodes["node_floor"], scene.RootNode.SubNodes["node_pyr"], scene.RootNode.SubNodes["node_box"]
	meshPyr.Models.Default().SetMatName("mat_mosaic")
	meshPyr.Models.Default().Clone("model_pyramid_dog").SetMatName("mat_dog")
	meshCube.Models.Default().SetMatName("mat_crate")
	meshCube.Models.Default().Clone("model_cube_cat").SetMatName("mat_cat")

	for i = 0; i < len(crates); i++ {
		if i == 0 { str = "model_cube_cat" } else { str = "" }
		crates[i] = scene.RootNode.MakeSubNode(ngine.Key("node_box_%v", i), "mesh_cube", str)
		f = float64(i)
		crates[i].Transform.SetPosXYZ((f + 3) * -2, (f + 1) * 2, (f + 2) * 3)
	}

	for i = 0; i < len(pyramids); i++ {
		if i > 1 { str = "model_pyramid_dog" } else { str = "" }
		pyramids[i] = scene.RootNode.MakeSubNode(ngine.Key("nody_pyr_%v", i), "mesh_pyramid", str)
		f = float64(len(pyramids) - i)
		pyramids[i].Transform.SetScalingN((f + 1) * 2)
		pyramids[i].Transform.SetPosXYZ((f + 3) * -4, (f + 2) * 3, (f + 2) * 14)
		if i > 1 {
			if i == 2 { f = 45 } else { f = 135 }
			pyramids[i].Transform.SetRotZ(numutil.DegToRad(f))
		} else {
			if i == 0 { f = 180 } else { f = 90 }
			pyramids[i].Transform.SetRotX(numutil.DegToRad(f))
		}
	}

	floor.SetMatName("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScalingN(100)

	ngine_samples.CamCtl.BeginUpdate(); ngine_samples.CamCtl.Pos.Y = 1.6; ngine_samples.CamCtl.EndUpdate()

	//	upload everything
	ngine.Core.SyncUpdates()
}
