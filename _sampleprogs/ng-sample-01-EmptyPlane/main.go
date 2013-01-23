package main

import (
	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
	// nga "github.com/go3d/go-ngine/assets"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor *ng.Node
)

func main() {
	ngsamples.SamplesMainFunc(LoadSampleScene_01_EmptyPlane)
}

func onLoop() {
	ngsamples.CheckToggleKeys()
	ngsamples.CheckCamCtlKeys()
}

func LoadSampleScene_01_EmptyPlane() {
	var (
		err       error
		scene     *ng.Scene
		meshFloor *ng.Mesh
		bufRest   *ng.MeshBuffer
	)

	ng.Loop.OnLoop = onLoop

	//	textures / materials
	ngsamples.AddTextureMaterials(map[string]string{
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
	scene = ngsamples.AddScene("")
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	floor.SetMatID("mat_cobbles")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScaleN(100)

	ngsamples.CamCtl.BeginUpdate()
	ngsamples.CamCtl.Pos.Y = 1.6
	ngsamples.CamCtl.EndUpdate()
}
