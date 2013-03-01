package main

import (
	"math"

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
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["pulse"]].GetTex2D(1).SetMixWeight(0.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*4)))
}

func setupScene() {
	var (
		err                    error
		scene                  *ng.Scene
		meshFloorID, meshBoxID int
		bufRest                *ng.MeshBuffer
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"dog":     "tex/dog.png",
		"cat":     "tex/cat.png",
		"gopher":  "tex/gopher.png",
		"crate":   "tex/crate.jpeg",
	})
	fxPulseID := ng.Core.Libs.Effects.AddNew()
	apputil.LibIDs.Fx["pulse"] = fxPulseID
	fxPulse := &ng.Core.Libs.Effects[fxPulseID]
	fxPulse.EnableTex2D(0).Tex_SetImageID(apputil.LibIDs.Img2D["crate"])
	fxPulse.EnableTex2D(1).Tex_SetImageID(apputil.LibIDs.Img2D["gopher"]).SetMixWeight(0.5)
	fxPulse.UpdateRoutine()

	dogMat := &ng.Core.Libs.Materials[apputil.LibIDs.Mat["dog"]]
	dogMat.FaceEffects.ByTag["top"] = apputil.LibIDs.Fx["cat"]
	dogMat.FaceEffects.ByTag["front"] = fxPulseID

	//	meshes / models
	if bufRest, err = ng.Core.MeshBuffers.AddNew("buf_rest", 200); err != nil {
		panic(err)
	}
	if meshFloorID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.MeshDescriptorPlane); err != nil {
		panic(err)
	}
	if meshBoxID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_box", ng.MeshDescriptorCube); err != nil {
		panic(err)
	}
	bufRest.Add(meshFloorID)
	bufRest.Add(meshBoxID)

	scene = apputil.AddMainScene()
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", meshFloorID)
	floor.MatID = apputil.LibIDs.Mat["cobbles"]
	floor.Transform.SetScale(100)
	floor.ApplyTransform()

	box = scene.RootNode.ChildNodes.AddNew("node_box", meshBoxID)
	box.MatID = apputil.LibIDs.Mat["dog"]
	box.Transform.Pos.Y = 2
	box.ApplyTransform()

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.Set(-2.5, 2, -7)
	camCtl.EndUpdate()
}
