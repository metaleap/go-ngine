package main

import (
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, box *ng.Node
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
		meshFloor, meshBox *ng.Mesh
		bufRest            *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
		"cat":     "tex/cat.png",
	})

	ng.Core.Libs.Materials["mat_dog"].FaceEffects.ByTag["top"] = "fx_cat"

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(200, 200)); err != nil {
		panic(err)
	}

	if meshFloor, err = ng.Core.Libs.Meshes.AddLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloor)

	if meshBox, err = ng.Core.Libs.Meshes.AddLoad("mesh_box", ng.MeshProviderPrefabCube); err != nil {
		panic(err)
	}
	bufRest.Add(meshBox)

	scene = apputil.AddScene("", true, "")
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	floor.SetMatID("mat_cobbles")
	floor.Transform.SetScaleN(100)

	box = scene.RootNode.ChildNodes.AddNew("node_box", "mesh_box", "")
	box.SetMatID("mat_dog")
	box.Transform.Pos.Y = 2
	box.Transform.ApplyMatrices()

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-2.5, 2, -7)
	camCtl.EndUpdate()
}
