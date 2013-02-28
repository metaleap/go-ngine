package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, pyr *ng.Node
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
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["pulse"]].Ops.GetColor(1).SetMixWeight(0.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*4)))
}

func setupScene() {
	var (
		err                    error
		scene                  *ng.Scene
		meshFloorID, meshPyrID int
		bufRest                *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
		"cat":     "tex/cat.png",
	})

	fxBlueID := ng.Core.Libs.Effects.AddNew()
	fxBlue := &ng.Core.Libs.Effects[fxBlueID]
	fxBlue.Ops.EnableColor(-1).SetRgb(0, 0.33, 0.66)
	fxBlue.Ops.EnableTex2D(-1).SetImageID(apputil.LibIDs.Img2D["dog"]).SetMixWeight(0.33)
	fxBlue.UpdateRoutine()

	fxGreenID := ng.Core.Libs.Effects.AddNew()
	fxGreen := &ng.Core.Libs.Effects[fxGreenID]
	fxGreen.Ops.EnableColor(-1).SetRgb(0, 0.66, 0.33)
	fxGreen.UpdateRoutine()

	fxCat := &ng.Core.Libs.Effects[apputil.LibIDs.Fx["cat"]]
	fxCat.Ops.EnableOrangify(-1)
	fxCat.UpdateRoutine()

	fxPulseID := ng.Core.Libs.Effects.AddNew()
	apputil.LibIDs.Fx["pulse"] = fxPulseID
	fxPulse := &ng.Core.Libs.Effects[fxPulseID]
	fxPulse.Ops.EnableColor(0).SetRgb(0.6, 0, 0)
	fxPulse.Ops.EnableColor(1).SetRgb(0.9, 0.7, 0).SetMixWeight(0.25)
	fxPulse.UpdateRoutine()

	dogMat := &ng.Core.Libs.Materials[apputil.LibIDs.Mat["dog"]]
	dogMat.DefaultEffectID = fxBlueID
	dogMat.FaceEffects.ByID["t0"] = apputil.LibIDs.Fx["cat"]
	dogMat.FaceEffects.ByID["t1"] = fxPulseID
	dogMat.FaceEffects.ByID["t2"] = fxGreenID

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(200, 200)); err != nil {
		panic(err)
	}
	if meshFloorID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}
	if meshPyrID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_pyr", ng.MeshProviderPrefabPyramid); err != nil {
		panic(err)
	}
	bufRest.Add(&ng.Core.Libs.Meshes[meshFloorID])
	bufRest.Add(&ng.Core.Libs.Meshes[meshPyrID])

	scene = apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshPyrID)
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", meshFloorID)
	floor.MatID = apputil.LibIDs.Mat["cobbles"]
	floor.Transform.SetScale(100)
	floor.ApplyTransform()

	pyr = scene.RootNode.ChildNodes.AddNew("node_pyr", meshPyrID)
	pyr.MatID = apputil.LibIDs.Mat["dog"]
	pyr.Transform.Pos.Y = 2
	pyr.ApplyTransform()

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-1.5, 2, -4)
	camCtl.EndUpdate()
}
