package main

import (
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, cube, dog *ng.Node
)

func main() {
	apputil.MaxKeyHint = 0
	//	You better check out this function, it's part of the "minimal go:ngine setup":
	apputil.Main(setupScene, onAppThread, onWinThread)
}

func onAppThread() {
	apputil.HandleCamCtlKeys()
}

//	called once per frame in main thread
func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()
}

func setupScene() {
	var (
		err                 error
		scene               *ng.Scene
		meshFloor, meshCube *ng.Mesh
		bufRest             *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
	})

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(200, 200)); err != nil {
		panic(err)
	}

	if meshFloor, err = ng.Core.Libs.Meshes.AddLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloor)

	if meshCube, err = ng.Core.Libs.Meshes.AddLoad("mesh_cube", ng.MeshProviderPrefabCube); err != nil {
		panic(err)
	}
	bufRest.Add(meshCube)

	//	scene
	scene = apputil.AddScene("", true, "mesh_cube")
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	floor.SetMatID("mat_cobbles")
	floor.Transform.SetScaleN(100)

	cube = scene.RootNode.ChildNodes.AddNew("node_cube", "mesh_cube", "")
	cube.SetMatID("mat_sky")

	dog = scene.RootNode.ChildNodes.AddNew("node_dog", "mesh_cube", "")
	dog.SetMatID("mat_dog")
	dog.Transform.SetPosX(-2)
	dog.Transform.SetPosZ(2)

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-1, 2, -5)
	camCtl.EndUpdate()
}
