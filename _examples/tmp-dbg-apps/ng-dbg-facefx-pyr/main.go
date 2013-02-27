package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, pyr *ng.Node
	fxPulse    *ng.FxEffect
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
	fxPulse.Ops.GetColor(1).SetMixWeight(0.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*4)))
}

func setupScene() {
	var (
		err                error
		scene              *ng.Scene
		meshFloor, meshPyr *ng.Mesh
		bufRest            *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
		"cat":     "tex/cat.png",
	})

	fxBlue := ng.Core.Libs.Effects.AddNew()
	fxBlue.Ops.EnableColor(-1).SetRgb(0, 0.33, 0.66)
	fxBlue.Ops.EnableTex2D(-1).SetImageID(apputil.LibIDs.Img["dog"]).SetMixWeight(0.33)
	fxBlue.UpdateRoutine()

	fxGreen := ng.Core.Libs.Effects.AddNew()
	fxGreen.Ops.EnableColor(-1).SetRgb(0, 0.66, 0.33)
	fxGreen.UpdateRoutine()

	fxCat := &ng.Core.Libs.Effects[apputil.LibIDs.Fx["cat"]]
	fxCat.Ops.EnableOrangify(-1)
	fxCat.UpdateRoutine()

	fxPulse = ng.Core.Libs.Effects.AddNew()
	fxPulse.Ops.EnableColor(0).SetRgb(0.6, 0, 0)
	fxPulse.Ops.EnableColor(1).SetRgb(0.9, 0.7, 0).SetMixWeight(0.25)
	fxPulse.UpdateRoutine()

	dogMat := &ng.Core.Libs.Materials[apputil.LibIDs.Mat["dog"]]
	dogMat.DefaultEffectID = fxBlue.ID
	dogMat.FaceEffects.ByID["t0"] = apputil.LibIDs.Fx["cat"]
	dogMat.FaceEffects.ByID["t1"] = fxPulse.ID
	dogMat.FaceEffects.ByID["t2"] = fxGreen.ID

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(200, 200)); err != nil {
		panic(err)
	}

	if meshFloor, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloor)

	if meshPyr, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_pyr", ng.MeshProviderPrefabPyramid); err != nil {
		panic(err)
	}
	bufRest.Add(meshPyr)

	scene = apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshPyr.ID)
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", meshFloor.ID)
	floor.MatID = apputil.LibIDs.Mat["cobbles"]
	floor.Transform.SetScale(100)
	floor.ApplyTransform()

	pyr = scene.RootNode.ChildNodes.AddNew("node_pyr", meshPyr.ID)
	pyr.MatID = apputil.LibIDs.Mat["dog"]
	pyr.Transform.Pos.Y = 2
	pyr.ApplyTransform()

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-1.5, 2, -4)
	camCtl.EndUpdate()
}
