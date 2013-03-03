package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	triNodeID, quadNodeID int

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
			if tmpNode = tmpScene.Node(triNodeID); tmpNode != nil {
				tmpNode.Transform.Rot.Add3(-0.005, -0.0005, 0)
				tmpNode.Transform.Pos.Set(0.85, 1*math.Sin(ng.Loop.Tick.Now), 4)
				tmpScene.ApplyNodeTransforms(triNodeID)
			}
			if tmpNode = tmpScene.Node(quadNodeID); tmpNode != nil {
				tmpNode.Transform.Rot.Add3(0, 0.005, 0.0005)
				tmpNode.Transform.Pos.Set(-0.85, 1*math.Cos(ng.Loop.Tick.Now), 4)
				tmpScene.ApplyNodeTransforms(quadNodeID)
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
	fxID := apputil.LibIDs.Fx["cat"]
	ng.Core.Libs.Effects[fxID].ToggleOrangify(0)
	ng.Core.Libs.Effects[fxID].UpdateRoutine()

	fxID = apputil.LibIDs.Fx["dog"]
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
	fx := &ng.Core.Libs.Effects[apputil.LibIDs.Fx["dog"]]
	fx.EnableOrangify(0)
	fx.UpdateRoutine()

	//	meshes / models

	if meshTriID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_tri", ng.MeshDescriptorTri); err != nil {
		panic(err)
	}
	if meshQuadID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_quad", ng.MeshDescriptorQuad); err != nil {
		panic(err)
	}

	if meshBuf, err = ng.Core.MeshBuffers.AddNew("meshbuf", 9); err != nil {
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
	triNodeID = apputil.AddNode(scene, 0, meshTriID, apputil.LibIDs.Mat["cat"], -1).ID
	quadNodeID = apputil.AddNode(scene, 0, meshQuadID, apputil.LibIDs.Mat["dog"], -1).ID
}
