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
	apputil.MaxKeyHint = 4
	apputil.OnSec = onSec
	apputil.Main(setupExample_01_TriQuad, onAppThread, onWinThread)
}

//	called once per second in main thread
func onSec() {
	fxID := apputil.LibIDs.Fx["cat"]
	ng.Core.Libs.Effects[fxID].Ops.ToggleOrangify(0)
	ng.Core.Libs.Effects[fxID].UpdateRoutine()

	fxID = apputil.LibIDs.Fx["dog"]
	ng.Core.Libs.Effects[fxID].Ops.ToggleOrangify(0)
	ng.Core.Libs.Effects[fxID].UpdateRoutine()
}

//	called once per frame in main thread
func onWinThread() {
	apputil.CheckAndHandleToggleKeys()
}

//	called once per frame in app thread
func onAppThread() {
	if !apputil.Paused {
		tri.Transform.Rot.Add3(-0.005, -0.0005, 0)
		tri.Transform.Pos.Set(0.85, 1*math.Sin(ng.Loop.Tick.Now), 4)
		tri.ApplyTransform()
		quad.Transform.Rot.Add3(0, 0.005, 0.0005)
		quad.Transform.Pos.Set(-0.85, 1*math.Cos(ng.Loop.Tick.Now), 4)
		quad.ApplyTransform()
	}
}

func setupExample_01_TriQuad() {
	var (
		err                   error
		scene                 *ng.Scene
		meshTriID, meshQuadID int
		meshBuf               *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cat": "tex/cat.png",
		"dog": "tex/dog.png",
	})
	fx := &ng.Core.Libs.Effects[apputil.LibIDs.Fx["dog"]]
	fx.Ops.EnableOrangify(0)
	fx.UpdateRoutine()

	//	meshes / models

	if meshTriID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_tri", ng.MeshProviderPrefabTri); err != nil {
		panic(err)
	}
	if meshQuadID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_quad", ng.MeshProviderPrefabQuad); err != nil {
		panic(err)
	}

	if meshBuf, err = ng.Core.MeshBuffers.AddNew("meshbuf", ng.MeshBufferParams{9, 9}); err != nil {
		panic(err)
	}
	if err = meshBuf.Add(&ng.Core.Libs.Meshes[meshTriID]); err != nil {
		panic(err)
	}
	if err = meshBuf.Add(&ng.Core.Libs.Meshes[meshQuadID]); err != nil {
		panic(err)
	}

	//	scene
	scene = apputil.AddMainScene()
	tri = scene.RootNode.ChildNodes.AddNew("node_tri", meshTriID)
	quad = scene.RootNode.ChildNodes.AddNew("node_quad", meshQuadID)
	tri.MatID = apputil.LibIDs.Mat["cat"]
	quad.MatID = apputil.LibIDs.Mat["dog"]
}
