package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floorID, pyrID int
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
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["pulse"]].GetColor(1).SetMixWeight(0.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*4)))
}

func setupScene() {
	var (
		err                    error
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
	fxBlue.EnableColor(-1).Color_SetRgb(0, 0.33, 0.66)
	fxBlue.EnableTex2D(-1).Tex_SetImageID(apputil.LibIDs.Img2D["dog"]).SetMixWeight(0.33)
	fxBlue.UpdateRoutine()

	fxGreenID := ng.Core.Libs.Effects.AddNew()
	fxGreen := &ng.Core.Libs.Effects[fxGreenID]
	fxGreen.EnableColor(-1).Color_SetRgb(0, 0.66, 0.33)
	fxGreen.UpdateRoutine()

	fxCat := &ng.Core.Libs.Effects[apputil.LibIDs.Fx["cat"]]
	fxCat.EnableOrangify(-1)
	fxCat.UpdateRoutine()

	fxPulseID := ng.Core.Libs.Effects.AddNew()
	apputil.LibIDs.Fx["pulse"] = fxPulseID
	fxPulse := &ng.Core.Libs.Effects[fxPulseID]
	fxPulse.EnableColor(0).Color_SetRgb(0.6, 0, 0)
	fxPulse.EnableColor(1).Color_SetRgb(0.9, 0.7, 0).SetMixWeight(0.25)
	fxPulse.UpdateRoutine()

	dogMat := &ng.Core.Libs.Materials[apputil.LibIDs.Mat["dog"]]
	dogMat.DefaultEffectID = fxBlueID
	dogMat.FaceEffects.ByID["t0"] = apputil.LibIDs.Fx["cat"]
	dogMat.FaceEffects.ByID["t1"] = fxPulseID
	dogMat.FaceEffects.ByID["t2"] = fxGreenID

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.AddNew("buf_rest", 200); err != nil {
		panic(err)
	}
	if meshFloorID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.MeshDescriptorPlane); err != nil {
		panic(err)
	}
	if meshPyrID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_pyr", ng.MeshDescriptorPyramid); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloorID)
	bufRest.Add(meshPyrID)

	scene := apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshPyrID)

	floor := apputil.AddNode(scene, 0, meshFloorID, apputil.LibIDs.Mat["cobbles"], -1)
	floorID = floor.ID
	floor.Transform.SetScale(100)

	pyr := apputil.AddNode(scene, 0, meshPyrID, apputil.LibIDs.Mat["dog"], -1)
	pyrID = pyr.ID
	pyr.Transform.Pos.Y = 2

	scene.ApplyNodeTransform(0)
	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-1.5, 2, -4)
	camCtl.EndUpdate()
}
