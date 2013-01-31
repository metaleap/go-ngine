package main

import (
	"math"

	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, pyr, box *ng.Node
)

func main() {
	ngsamples.SamplesMainFunc(LoadSampleScene_02_PyrCube)
}

func onWinThread() {
	ngsamples.CheckCamCtlKeys()
	ngsamples.CheckToggleKeys()
	ngsamples.HandleToggleKeys()
}

func onAppThread() {
	if !ngsamples.Paused {
		ngsamples.HandleCamCtlKeys()

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
}

func LoadSampleScene_02_PyrCube() {
	var (
		err                          error
		scene                        *ng.Scene
		meshFloor, meshPyr, meshCube *ng.Mesh
		bufFloor, bufRest            *ng.MeshBuffer
	)

	ng.Loop.OnAppThread, ng.Loop.OnWinThread = onAppThread, onWinThread

	//	textures / materials
	ngsamples.AddTextureMaterials(map[string]string{
		"cobbles": "http://dl.dropbox.com/u/136375/go-ngine/assets/tex/cobbles.png",
		"crate":   "tex/crate.jpeg",
		"mosaic":  "tex/mosaic.jpeg",
	})
	ngsamples.AddColorMaterials(map[string][]float64{
		"yellow": []float64{1, 1, 0},
		"blue":   []float64{0, 0, 1},
	})

	//	meshes / models
	if bufFloor, err = ng.Core.MeshBuffers.Add("buf_floor", ng.Core.MeshBuffers.NewParams(6, 6)); err != nil {
		panic(err)
	}
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(36+12, 36+12)); err != nil {
		panic(err)
	}
	if meshFloor, err = ng.Core.Libs.Meshes.AddLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}
	if meshPyr, err = ng.Core.Libs.Meshes.AddLoad("mesh_pyramid", ng.MeshProviderPrefabPyramid); err != nil {
		panic(err)
	}
	if meshCube, err = ng.Core.Libs.Meshes.AddLoad("mesh_cube", ng.MeshProviderPrefabCube); err != nil {
		panic(err)
	}
	bufFloor.Add(meshFloor)
	bufRest.Add(meshCube)
	bufRest.Add(meshPyr)
	meshPyr.Models.Default().SetMatID("mat_mosaic")
	meshCube.Models.Default().SetMatID("mat_crate")
	ng.Core.Libs.Materials["mat_crate"].FaceEffects.ByTag["back"] = "fx_mosaic"

	//	scene
	scene = ngsamples.AddScene("", true)
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	pyr = scene.RootNode.ChildNodes.AddNew("node_pyr", "mesh_pyramid", "")
	box = scene.RootNode.ChildNodes.AddNew("node_box", "mesh_cube", "")

	floor.SetMatID("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScaleN(1000)

	ngsamples.CamCtl.BeginUpdate()
	ngsamples.CamCtl.Pos.Y = 1.6
	ngsamples.CamCtl.EndUpdate()
}
