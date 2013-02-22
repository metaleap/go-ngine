package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, pyr, box *ng.Node
)

func main() {
	apputil.Main(setupExample_03_PyrCube, onAppThread, onWinThread)
}

func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()

	//	pulsating materials
	ng.Core.Libs.Effects["fx_cat"].Ops.GetOrangify(0).SetMixWeight(0.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*4)))
	ng.Core.Libs.Effects["fx_dog"].Ops.GetTex2D(1).SetMixWeight(0.5 + (0.5 * math.Cos(ng.Loop.Tick.Now*2)))
}

func onAppThread() {
	if !apputil.Paused {
		apputil.HandleCamCtlKeys()

		//	animate mesh nodes
		pyr.Transform.Rot.Add3(-0.0005, -0.0005, 0)
		pyr.Transform.Pos.Set(-1.5, 1.5+(2*math.Sin(ng.Loop.Tick.Now*3)), 7)
		pyr.Transform.ApplyMatrices()

		box.Transform.Rot.Add3(0, 0.0004, 0.0006)
		box.Transform.Pos.Set(1.5, 1.5+(2*math.Cos(ng.Loop.Tick.Now*0.3333)), 7)
		box.Transform.ApplyMatrices()
	}
}

func setupExample_03_PyrCube() {
	var (
		err                          error
		scene                        *ng.Scene
		meshFloor, meshPyr, meshCube *ng.Mesh
		bufFloor, bufRest            *ng.MeshBuffer
	)

	urlPrefix := "http://dl.dropbox.com/u/136375/go-ngine/assets/"
	urlPrefix = ""

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": urlPrefix + "tex/cobbles.png",
		"crate":   "tex/crate.jpeg",
		"mosaic":  "tex/mosaic.jpeg",
		"gopher":  "tex/gopher.png",
		"dog":     "tex/dog.png",
		"cat":     "tex/cat.png",
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

	fx := ng.Core.Libs.Effects["fx_cat"]
	fx.Ops.EnableOrangify(-1).SetMixWeight(0.5)
	fx.UpdateRoutine()
	fx = ng.Core.Libs.Effects["fx_dog"]
	fx.Ops.EnableTex2D(1).SetImageID("img_gopher").SetMixWeight(0.5)
	fx.UpdateRoutine()
	ng.Core.Libs.Materials["mat_crate"].FaceEffects.ByTag["front"] = "fx_dog"
	ng.Core.Libs.Materials["mat_crate"].FaceEffects.ByTag["back"] = "fx_dog"
	ng.Core.Libs.Materials["mat_mosaic"].FaceEffects.ByID["t3"] = "fx_cat"

	//	scene
	scene = apputil.AddScene("", true, "mesh_cube")
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	pyr = scene.RootNode.ChildNodes.AddNew("node_pyr", "mesh_pyramid", "")
	box = scene.RootNode.ChildNodes.AddNew("node_box", "mesh_cube", "")

	floor.SetMatID("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScaleN(1000)

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.X, camCtl.Pos.Y, camCtl.Pos.Z = -1, 1.6, -2
	camCtl.EndUpdate()
}
