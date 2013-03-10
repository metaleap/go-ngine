package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	catTri, dogQuad int

	tmpScene *ng.Scene
	tmpNode  *ng.SceneNode
)

func main() {
	apputil.MaxKeyHint = 4
	apputil.OnSec = onSec
	apputil.Main(setupExample_01_TriQuad, onAppThread, onWinThread)
}

//	called once per frame in app thread
func onAppThread() {
	if !apputil.Paused {
		if tmpScene = apputil.SceneCam.Scene(); tmpScene != nil {
			if tmpNode = tmpScene.Node(catTri); tmpNode != nil {
				tmpNode.Transform.Rot.Add3(-0.005, -0.0005, 0)
				tmpNode.Transform.Pos.Set(0.85, 1*math.Sin(ng.Loop.Tick.Now), 4)
				tmpScene.ApplyNodeTransforms(catTri)
			}
			if tmpNode = tmpScene.Node(dogQuad); tmpNode != nil {
				tmpNode.Transform.Rot.Add3(0, 0.005, 0.0005)
				tmpNode.Transform.Pos.Set(-0.85, 1*math.Cos(ng.Loop.Tick.Now), 4)
				tmpScene.ApplyNodeTransforms(dogQuad)
			}
		}
	}
}

//	called once per frame in main thread
func onWinThread() {
	apputil.CheckAndHandleToggleKeys()
}

//	called once per second in main thread
func onSec() {
	toggleOrange(apputil.LibIDs.Fx["cat"])
	toggleOrange(apputil.LibIDs.Fx["dog"])
}

func toggleOrange(fxID int) {
	ng.Core.Libs.Effects[fxID].ToggleOrangify(0)
	ng.Core.Libs.Effects[fxID].UpdateRoutine()
}

func setupExample_01_TriQuad() {
	var (
		err                   error
		meshTriID, meshQuadID int
		meshBuf               *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cat": "tex/cat.png",
		"dog": "tex/dog.png",
	})
	toggleOrange(apputil.LibIDs.Fx["dog"])

	//	meshes / models

	if meshTriID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_tri", ng.Core.Mesh.Desc.Tri); err != nil {
		panic(err)
	}
	if meshQuadID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_quad", ng.Core.Mesh.Desc.Quad); err != nil {
		panic(err)
	}

	if meshBuf, err = ng.Core.Mesh.Buffers.AddNew("meshbuf", 3+6); err != nil {
		panic(err)
	}
	if err = meshBuf.Add(meshTriID); err != nil {
		panic(err)
	}
	if err = meshBuf.Add(meshQuadID); err != nil {
		panic(err)
	}

	//	scene
	scene := apputil.AddMainScene()
	catTri = apputil.AddNode(scene, 0, meshTriID, apputil.LibIDs.Mat["cat"], -1).ID
	dogQuad = apputil.AddNode(scene, 0, meshQuadID, apputil.LibIDs.Mat["dog"], -1).ID
}
