package main

import (
	"math"

	unum "github.com/metaleap/go-util/num"

	ng "github.com/go3d/go-ngine/core"
	nga "github.com/go3d/go-ngine/assets"
	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	floor, pyr, box *ng.Node
	crates [3]*ng.Node
	pyramids [4]*ng.Node
	i int
	f float64
)

func main () {
	ngsamples.SamplesMainFunc(LoadSampleScene_03_PyrsCubes)
}

func onLoop () {
	ngsamples.CheckToggleKeys()
	ngsamples.CheckCamCtlKeys()

	//	animate mesh nodes
	pyr.NodeTransform.Rot.X -= 0.0005
	pyr.NodeTransform.Rot.Y -= 0.0005
	pyr.NodeTransform.Pos.Set(-13.75, 2 * math.Sin(ng.Loop.TickNow), 2)
	pyr.NodeTransform.OnPosRotChanged()

	box.NodeTransform.Rot.Y += 0.0004
	box.NodeTransform.Rot.Z += 0.0006
	box.NodeTransform.Pos.Set(-8.125, 2 * math.Cos(ng.Loop.TickNow), -2)
	box.NodeTransform.OnPosRotChanged()

	for i = 0; i < len(crates); i++ {
		f = float64(i)
		f = (f + 1) * (f + 1)
		crates[i].NodeTransform.Rot.X += f * 0.00001
		crates[i].NodeTransform.Rot.Y += f * 0.0001
		crates[i].NodeTransform.Rot.Z += f * 0.001
		crates[i].NodeTransform.OnRotChanged()
	}

	pyramids[0].NodeTransform.SetPosX(math.Sin(ng.Loop.TickNow) * 100)
	pyramids[1].NodeTransform.SetPosZ(math.Cos(ng.Loop.TickNow) * 1000)
}

func LoadSampleScene_03_PyrsCubes () {
	var err error
	var meshFloor, meshPyr, meshCube *ng.Mesh
	var bufFloor, bufRest *ng.MeshBuffer
	var str string

	ng.Loop.OnLoop = onLoop
	ngsamples.Cam.Options.BackfaceCulling = false

	//	textures / materials
	ng.Core.Textures["tex_cobbles"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/cobbles.png")
	ng.Core.Textures["tex_crate"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/crate.jpeg")
	ng.Core.Textures["tex_mosaic"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/mosaic.jpeg")
	ng.Core.Textures["tex_cat"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/cat.png")
	ng.Core.Textures["tex_dog"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/dog.png")
	nga.Materials["mat_cobbles"] = nga.Materials.New("tex_cobbles")
	nga.Materials["mat_crate"] = nga.Materials.New("tex_crate")
	nga.Materials["mat_mosaic"] = nga.Materials.New("tex_mosaic")
	nga.Materials["mat_cat"] = nga.Materials.New("tex_cat")
	nga.Materials["mat_dog"] = nga.Materials.New("tex_dog")

	//	meshes / models
	if bufFloor, err = ng.Core.MeshBuffers.Add("buf_floor", ng.Core.MeshBuffers.NewParams(6, 6)); err != nil { panic(err) }
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(36 + 12, 36 + 12)); err != nil { panic(err) }
	if meshFloor, err = ng.Core.Meshes.Load("mesh_plane", ng.MeshProviders.PrefabPlane); err != nil { panic(err) }

	if meshPyr, err = ng.Core.Meshes.Load("mesh_pyramid", ng.MeshProviders.PrefabPyramid); err != nil { panic(err) }
	meshPyr.Models.Default().SetMatName("mat_mosaic")
	meshPyr.Models.Default().Clone("model_pyramid_dog").SetMatName("mat_dog")

	if meshCube, err = ng.Core.Meshes.Load("mesh_cube", ng.MeshProviders.PrefabCube); err != nil { panic(err) }
	meshCube.Models.Default().SetMatName("mat_crate")
	meshCube.Models.Default().Clone("model_cube_cat").SetMatName("mat_cat")

	ng.Core.Meshes.AddRange(meshFloor, meshPyr, meshCube)
	bufFloor.Add(meshFloor); bufRest.Add(meshCube); bufRest.Add(meshPyr)

	//	scene
	var scene = ng.NewScene()
	ng.Core.Scenes[""] = scene
	scene.RootNode.SubNodes.MakeN("node_floor", "mesh_plane", "", "node_pyr", "mesh_pyramid", "", "node_box", "mesh_cube", "")
	floor, pyr, box = /* scene.RootNode.SubNodes.Get("node_floor", "node_pyr", "node_box") */ scene.RootNode.SubNodes.M["node_floor"], scene.RootNode.SubNodes.M["node_pyr"], scene.RootNode.SubNodes.M["node_box"]

	for i = 0; i < len(crates); i++ {
		if i == 0 { str = "model_cube_cat" } else { str = "" }
		crates[i] = scene.RootNode.SubNodes.Make(ng.Fmt("node_box_%v", i), "mesh_cube", str)
		f = float64(i)
		crates[i].NodeTransform.SetPosXYZ((f + 3) * -2, (f + 1) * 2, (f + 2) * 3)
		if i == 2 {
			crates[i].SetMatName("mat_dog")
		}
	}

	for i = 0; i < len(pyramids); i++ {
		if i > 1 { str = "model_pyramid_dog" } else { str = "" }
		pyramids[i] = scene.RootNode.SubNodes.Make(ng.Fmt("nody_pyr_%v", i), "mesh_pyramid", str)
		f = float64(len(pyramids) - i)
		pyramids[i].NodeTransform.SetScalingN((f + 1) * 2)
		pyramids[i].NodeTransform.SetPosXYZ((f + 3) * -4, (f + 2) * 3, (f + 2) * 14)
		if i > 1 {
			if i == 2 { f = 45 } else { f = 135 }
			pyramids[i].NodeTransform.SetRotZ(unum.DegToRad(f))
		} else {
			if i == 0 { f = 180 } else { f = 90 }
			pyramids[i].NodeTransform.SetRotX(unum.DegToRad(f))
		}
		if i == 1 {
			pyramids[i].NodeTransform.SetScalingN(100)
			pyramids[i].NodeTransform.Pos.Y += 100
		}
	}

	floor.SetMatName("mat_cobbles")
	floor.NodeTransform.SetPosXYZ(0.1, 0, -8)
	floor.NodeTransform.SetScalingN(10000)

	ngsamples.CamCtl.BeginUpdate(); ngsamples.CamCtl.Pos.Y = 1.6; ngsamples.CamCtl.EndUpdate()

	//	upload everything
	ng.Core.SyncUpdates()
}
