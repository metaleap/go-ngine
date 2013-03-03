package main

import (
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floorNodeID int
)

func main() {
	apputil.Main(setupExample_02_EmptyGround, onAppThread, onWinThread)
}

func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()
}

func onAppThread() {
	apputil.HandleCamCtlKeys()
}

func setupExample_02_EmptyGround() {
	var (
		err         error
		meshPlaneID int
		bufRest     *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
	})

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.AddNew("buf_rest", 100); err != nil {
		panic(err)
	}
	if meshPlaneID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.MeshDescriptorPlane); err != nil {
		panic(err)
	}

	bufRest.Add(meshPlaneID) // bufRest.Add(meshCube);

	//	scene
	scene := apputil.AddMainScene()
	floor := apputil.AddNode(scene, 0, meshPlaneID, apputil.LibIDs.Mat["cobbles"], -1)
	floorNodeID = floor.ID
	floor.Transform.SetPos(0.1, 0, -8)
	floor.Transform.SetScale(100)
	scene.ApplyNodeTransforms(floorNodeID)

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Y = 1.6
	camCtl.EndUpdate()
}
