//	Renders a cube-mapped crate next to a normal texture-mapped crate.
package main

import (
	apputil "github.com/go3d/go-ngine/___old2013/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/___old2013/core"
)

var (
	floorID, cubeID, dogID int
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
		meshFloorID, meshCubeID int
		bufRest                 *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
	})

	//	meshes / models
	if bufRest, err = ng.Core.Mesh.Buffers.AddNew("buf_rest", 200); err != nil {
		panic(err)
	}

	if meshFloorID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.Core.Mesh.Desc.Plane); err != nil {
		panic(err)
	}

	if meshCubeID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_cube", ng.Core.Mesh.Desc.Cube); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloorID)
	bufRest.Add(meshCubeID)

	//	scene
	scene := apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshCubeID)
	floor := apputil.AddNode(scene, 0, meshFloorID, apputil.LibIDs.Mat["cobbles"], -1)
	floorID = floor.ID
	floor.Transform.SetScale(100)

	cubeID = apputil.AddNode(scene, 0, meshCubeID, apputil.LibIDs.Mat["sky"], -1).ID
	dog := apputil.AddNode(scene, 0, meshCubeID, apputil.LibIDs.Mat["dog"], -1)
	dogID = dog.ID
	dog.Transform.Pos.X, dog.Transform.Pos.Z = -2, 2
	scene.ApplyNodeTransforms(0)

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-1, 2, -5)
	camCtl.EndUpdate()
}
