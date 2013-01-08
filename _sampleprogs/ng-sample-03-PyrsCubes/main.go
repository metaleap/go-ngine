package main

import (
	"math"

	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
	// nga "github.com/go3d/go-ngine/assets"
	ng "github.com/go3d/go-ngine/core"
	unum "github.com/metaleap/go-util/num"
)

var (
	floor, pyr, box *ng.Node
	crates          [3]*ng.Node
	pyramids        [4]*ng.Node
	i               int
	f               float64
)

func main() {
	ngsamples.SamplesMainFunc(LoadSampleScene_03_PyrsCubes)
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

	for i = 0; i < len(crates); i++ {
		f = float64(i)
		f = (f + 1) * (f + 1)
		crates[i].Transform.Rot.X += f * 0.00001
		crates[i].Transform.Rot.Y += f * 0.0001
		crates[i].Transform.Rot.Z += f * 0.001
		crates[i].Transform.OnRotChanged()
	}

	pyramids[0].Transform.SetPosX(math.Sin(ng.Loop.TickNow) * 100)
	pyramids[1].Transform.SetPosZ(math.Cos(ng.Loop.TickNow) * 1000)
}

func LoadSampleScene_03_PyrsCubes() {
	var (
		err                          error
		scene                        *ng.Scene
		meshFloor, meshPyr, meshCube *ng.Mesh
		bufFloor, bufRest            *ng.MeshBuffer
		str                          string
	)

	ng.Loop.OnLoop = onLoop
	ng.Core.Rendering.States.DisableFaceCulling()

	//	textures / materials
	ngsamples.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"crate":   "tex/crate.jpeg",
		"mosaic":  "tex/mosaic.jpeg",
		"cat":     "tex/cat.png",
		"dog":     "tex/dog.png",
	})

	//	meshes / models
	if bufFloor, err = ng.Core.MeshBuffers.Add("buf_floor", ng.Core.MeshBuffers.NewParams(6, 6)); err != nil {
		panic(err)
	}
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(36+12, 36+12)); err != nil {
		panic(err)
	}
	if meshFloor, err = ng.Core.Libs.Meshes.Load("mesh_plane", ng.MeshProviders.PrefabPlane); err != nil {
		panic(err)
	}

	if meshPyr, err = ng.Core.Libs.Meshes.Load("mesh_pyramid", ng.MeshProviders.PrefabPyramid); err != nil {
		panic(err)
	}
	meshPyr.Models.Default().SetMatName("mat_mosaic")
	meshPyr.Models.Default().Clone("model_pyramid_dog").SetMatName("mat_dog")

	if meshCube, err = ng.Core.Libs.Meshes.Load("mesh_cube", ng.MeshProviders.PrefabCube); err != nil {
		panic(err)
	}
	meshCube.Models.Default().SetMatName("mat_crate")
	meshCube.Models.Default().Clone("model_cube_cat").SetMatName("mat_cat")

	ng.Core.Libs.Meshes.AddRange(meshFloor, meshPyr, meshCube)
	bufFloor.Add(meshFloor)
	bufRest.Add(meshCube)
	bufRest.Add(meshPyr)

	//	scene
	scene = ngsamples.AddScene("")
	scene.RootNode.SubNodes.MakeN("node_floor", "mesh_plane", "", "node_pyr", "mesh_pyramid", "", "node_box", "mesh_cube", "")
	floor, pyr, box = /* scene.RootNode.SubNodes.Get("node_floor", "node_pyr", "node_box") */ scene.RootNode.SubNodes.M["node_floor"], scene.RootNode.SubNodes.M["node_pyr"], scene.RootNode.SubNodes.M["node_box"]

	for i = 0; i < len(crates); i++ {
		if i == 0 {
			str = "model_cube_cat"
		} else {
			str = ""
		}
		crates[i] = scene.RootNode.SubNodes.Make(ng.Sfmt("node_box_%v", i), "mesh_cube", str)
		f = float64(i)
		crates[i].Transform.SetPosXYZ((f+3)*-2, (f+1)*2, (f+2)*3)
		if i == 2 {
			crates[i].SetMatName("mat_dog")
		}
	}

	for i = 0; i < len(pyramids); i++ {
		if i > 1 {
			str = "model_pyramid_dog"
		} else {
			str = ""
		}
		pyramids[i] = scene.RootNode.SubNodes.Make(ng.Sfmt("nody_pyr_%v", i), "mesh_pyramid", str)
		f = float64(len(pyramids) - i)
		pyramids[i].Transform.SetScaleN((f + 1) * 2)
		pyramids[i].Transform.SetPosXYZ((f+3)*-4, (f+2)*3, (f+2)*14)
		if i > 1 {
			if i == 2 {
				f = 45
			} else {
				f = 135
			}
			pyramids[i].Transform.SetRotZ(unum.DegToRad(f))
		} else {
			if i == 0 {
				f = 180
			} else {
				f = 90
			}
			pyramids[i].Transform.SetRotX(unum.DegToRad(f))
		}
		if i == 1 {
			pyramids[i].Transform.SetScaleN(100)
			pyramids[i].Transform.Pos.Y += 100
		}
	}

	floor.SetMatName("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScaleN(10000)

	ngsamples.CamCtl.BeginUpdate()
	ngsamples.CamCtl.Pos.X, ngsamples.CamCtl.Pos.Y, ngsamples.CamCtl.Pos.Z = 35, 1.6, 24
	ngsamples.CamCtl.EndUpdate()
}
