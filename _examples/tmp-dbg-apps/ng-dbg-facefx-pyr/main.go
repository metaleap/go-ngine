package main

import (
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, pyr *ng.Node
)

func main() {
	apputil.Main(setupScene, onAppThread, onWinThread)
}

func onAppThread() {
	apputil.HandleCamCtlKeys()
}

func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()
}

func setupScene() {
	var (
		err                error
		scene              *ng.Scene
		meshFloor, meshPyr *ng.Mesh
		bufRest            *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
		"cat":     "tex/cat.png",
	})

	ng.Core.Libs.Materials["mat_dog"].FaceEffects.ByID["t1"] = "fx_cat"
	ng.Core.Libs.Materials["mat_dog"].FaceEffects.ByID["t3"] = "fx_cat"

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(200, 200)); err != nil {
		panic(err)
	}

	if meshFloor, err = ng.Core.Libs.Meshes.AddLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloor)

	if meshPyr, err = ng.Core.Libs.Meshes.AddLoad("mesh_pyr", ng.MeshProviderPrefabPyramid); err != nil {
		panic(err)
	}
	bufRest.Add(meshPyr)

	scene = apputil.AddScene("", true, "mesh_pyr")
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	floor.SetMatID("mat_cobbles")
	floor.Transform.SetScaleN(100)

	pyr = scene.RootNode.ChildNodes.AddNew("node_pyr", "mesh_pyr", "")
	pyr.SetMatID("mat_dog")
	pyr.Transform.Pos.Y = 2
	pyr.Transform.ApplyMatrices()

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-2.5, 2, -7)
	camCtl.EndUpdate()
}
