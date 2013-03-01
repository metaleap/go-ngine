package main

import (
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, cube, dog *ng.Node
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
		err                     error
		scene                   *ng.Scene
		meshFloorID, meshCubeID int
		bufRest                 *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
	})

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.AddNew("buf_rest", 200); err != nil {
		panic(err)
	}

	if meshFloorID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}

	if meshCubeID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_cube", ng.MeshProviderPrefabCube); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloorID)
	bufRest.Add(meshCubeID)

	//	scene
	scene = apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshCubeID)
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", meshFloorID)
	floor.MatID = apputil.LibIDs.Mat["cobbles"]
	floor.Transform.SetScale(100)
	floor.ApplyTransform()

	cube = scene.RootNode.ChildNodes.AddNew("node_cube", meshCubeID)
	cube.MatID = apputil.LibIDs.Mat["sky"]

	dog = scene.RootNode.ChildNodes.AddNew("node_dog", meshCubeID)
	dog.MatID = apputil.LibIDs.Mat["dog"]
	dog.Transform.Pos.X, dog.Transform.Pos.Z = -2, 2
	dog.ApplyTransform()

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-1, 2, -5)
	camCtl.EndUpdate()
}
