package main

import (
	"math"
	"time"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floorNodeID, pyrNodeID, boxNodeID int
	tmpScene                          *ng.Scene
	tmpNode                           *ng.SceneNode
)

func main() {
	apputil.Main(setupExample_03_PyrCube, onAppThread, onWinThread)
}

func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()

	//	pulsating materials
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["cat"]].GetOrangify(0).SetMixWeight(0.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*4)))
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["dog"]].GetTex2D(1).SetMixWeight(0.5 + (0.5 * math.Cos(ng.Loop.Tick.Now*2)))
}

func onAppThread() {
	if !apputil.Paused {
		apputil.HandleCamCtlKeys()
		//	animate mesh nodes
		if tmpScene = apputil.SceneCam.Scene(); tmpScene != nil {
			if tmpNode = tmpScene.Node(pyrNodeID); tmpNode != nil {
				tmpNode.Transform.Rot.Add3(-0.0005, -0.0005, 0)
				tmpNode.Transform.Pos.Set(-1.5, 1.5+(2*math.Sin(ng.Loop.Tick.Now*3)), 7)
				tmpScene.ApplyNodeTransforms(tmpNode.ID)
			}
			if tmpNode = tmpScene.Node(boxNodeID); tmpNode != nil {
				tmpNode.Transform.Rot.Add3(0, 0.0004, 0.0006)
				tmpNode.Transform.Pos.Set(1.5, 1.5+(2*math.Cos(ng.Loop.Tick.Now*0.3333)), 7)
				tmpScene.ApplyNodeTransforms(tmpNode.ID)
			}
		}
	}
}

func setupExample_03_PyrCube() {
	var (
		err                                error
		meshPlaneID, meshPyrID, meshCubeID int
		bufFloor, bufRest                  *ng.MeshBuffer
	)

	urlPrefix := "http://dl.dropbox.com/u/136375/go-ngine/assets/"
	urlPrefix = ""

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": urlPrefix + "tex/cobbles.png",
		"crate":   "tex/crate.jpeg",
		"mosaic":  "tex/mosaic.jpeg",
		"gopher":  "tex/gopher.png",
		"dog":     "tex/dog.png",
		"cat":     "tex/cat.png",
	})

	//	meshes
	if bufFloor, err = ng.Core.MeshBuffers.AddNew("buf_floor", 6); err != nil {
		panic(err)
	}
	if bufRest, err = ng.Core.MeshBuffers.AddNew("buf_rest", 36+12); err != nil {
		panic(err)
	}
	if meshPlaneID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.MeshDescriptorPlane); err != nil {
		panic(err)
	}
	if meshPyrID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_pyramid", ng.MeshDescriptorPyramid); err != nil {
		panic(err)
	}
	if meshCubeID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_cube", ng.MeshDescriptorCube); err != nil {
		panic(err)
	}
	bufFloor.Add(meshPlaneID)
	bufRest.Add(meshCubeID)
	bufRest.Add(meshPyrID)

	fx := &ng.Core.Libs.Effects[apputil.LibIDs.Fx["cat"]]
	fx.EnableOrangify(-1).SetMixWeight(0.5)
	fx.UpdateRoutine()
	fx = &ng.Core.Libs.Effects[apputil.LibIDs.Fx["dog"]]
	fx.EnableTex2D(1).Tex_SetImageID(apputil.LibIDs.Img2D["gopher"]).SetMixWeight(0.5)
	fx.UpdateRoutine()
	ng.Core.Libs.Materials[apputil.LibIDs.Mat["crate"]].FaceEffects.ByTag["front"] = apputil.LibIDs.Fx["dog"]
	ng.Core.Libs.Materials[apputil.LibIDs.Mat["mosaic"]].FaceEffects.ByID["t3"] = apputil.LibIDs.Fx["cat"]

	//	scene / nodes
	scene := apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshCubeID)
	floor := apputil.AddNode(scene, 0, meshPlaneID, apputil.LibIDs.Mat["cobbles"], -1)
	floorNodeID = floor.ID
	floor.Transform.SetPos(0.1, 0, -8)
	floor.Transform.SetScale(1000)
	scene.ApplyNodeTransforms(floorNodeID)

	pyrNodeID = apputil.AddNode(scene, 0, meshPyrID, apputil.LibIDs.Mat["mosaic"], -1).ID
	boxNodeID = apputil.AddNode(scene, 0, meshCubeID, apputil.LibIDs.Mat["crate"], -1).ID

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.X, camCtl.Pos.Y, camCtl.Pos.Z = 0, 1.7, -10
	camCtl.EndUpdate()
	if false {
		ng.Loop.Delay = 100 * time.Millisecond
	}
}
