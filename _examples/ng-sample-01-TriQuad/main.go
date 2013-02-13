package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	tri, quad *ng.Node
)

func main() {
	apputil.MaxKeyHint = 3
	apputil.OnSec = onSec
	apputil.Main(setupExample_01_TriQuad, onAppThread, onWinThread)
}

//	called once per second in main thread
func onSec() {
	var fx = ng.Core.Libs.Effects["fx_cat"]
	fx.Ops.ToggleOrangify(0)
	fx.UpdateRoutine()

	fx = ng.Core.Libs.Effects["fx_dog"]
	fx.Ops.ToggleOrangify(0)
	fx.UpdateRoutine()
}

//	called once per frame in main thread
func onWinThread() {
	apputil.CheckAndHandleToggleKeys()
}

//	called once per frame in app thread
func onAppThread() {
	if !apputil.Paused {
		tri.Transform.Rot.Add3(-0.005, -0.005, 0)
		tri.Transform.Pos.Set(-3.75, 1*math.Sin(ng.Loop.Tick.Now), 1)
		tri.Transform.ApplyMatrices()
		quad.Transform.Rot.Add3(0, 0.001, 0.001)
		quad.Transform.Pos.Set(-4.125, 1*math.Cos(ng.Loop.Tick.Now), 0)
		quad.Transform.ApplyMatrices()
	}
}

func setupExample_01_TriQuad() {
	var (
		err               error
		scene             *ng.Scene
		meshTri, meshQuad *ng.Mesh
		meshBuf           *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cat": "tex/cat.png",
		"dog": "tex/dog.png",
	})
	fx := ng.Core.Libs.Effects["fx_dog"]
	fx.Ops.EnableOrangify(0)
	fx.UpdateRoutine()

	//	meshes / models

	if meshTri, err = ng.Core.Libs.Meshes.AddLoad("mesh_tri", ng.MeshProviderPrefabTri); err != nil {
		panic(err)
	}
	if meshQuad, err = ng.Core.Libs.Meshes.AddLoad("mesh_quad", ng.MeshProviderPrefabQuad); err != nil {
		panic(err)
	}

	if meshBuf, err = ng.Core.MeshBuffers.Add("meshbuf", ng.Core.MeshBuffers.NewParams(9, 9)); err != nil {
		panic(err)
	}
	if err = meshBuf.Add(meshTri); err != nil {
		panic(err)
	}
	if err = meshBuf.Add(meshQuad); err != nil {
		panic(err)
	}

	//	scene
	scene = apputil.AddScene("", true)
	tri = scene.RootNode.ChildNodes.AddNew("node_tri", "mesh_tri", "")
	quad = scene.RootNode.ChildNodes.AddNew("node_quad", "mesh_quad", "")
	tri.SetMatID("mat_cat")
	quad.SetMatID("mat_dog")
}
