package main

import (
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor *ng.Node
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
		err       error
		scene     *ng.Scene
		meshFloor *ng.Mesh
		bufRest   *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
	})

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(100, 100)); err != nil {
		panic(err)
	}
	if meshFloor, err = ng.Core.Libs.Meshes.AddLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}

	bufRest.Add(meshFloor) // bufRest.Add(meshCube);

	//	scene
	scene = apputil.AddScene("", true, "")
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	floor.SetMatID("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScaleN(100)

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Y = 1.6
	camCtl.EndUpdate()
}
