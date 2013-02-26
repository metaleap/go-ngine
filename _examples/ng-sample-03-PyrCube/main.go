package main

import (
	"math"

	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	floor, pyr, box *ng.Node
)

func main() {
	apputil.Main(setupExample_03_PyrCube, onAppThread, onWinThread)
}

func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()

	//	pulsating materials
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["cat"]].Ops.GetOrangify(0).SetMixWeight(0.5 + (0.5 * math.Sin(ng.Loop.Tick.Now*4)))
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["dog"]].Ops.GetTex2D(1).SetMixWeight(0.5 + (0.5 * math.Cos(ng.Loop.Tick.Now*2)))
}

func onAppThread() {
	if !apputil.Paused {
		apputil.HandleCamCtlKeys()

		//	animate mesh nodes
		pyr.Transform.Rot.Add3(-0.0005, -0.0005, 0)
		pyr.Transform.Pos.Set(-1.5, 1.5+(2*math.Sin(ng.Loop.Tick.Now*3)), 7)
		pyr.ApplyTransform()

		box.Transform.Rot.Add3(0, 0.0004, 0.0006)
		box.Transform.Pos.Set(1.5, 1.5+(2*math.Cos(ng.Loop.Tick.Now*0.3333)), 7)
		box.ApplyTransform()
	}
}

func setupExample_03_PyrCube() {
	var (
		err                          error
		scene                        *ng.Scene
		meshFloor, meshPyr, meshCube *ng.Mesh
		bufFloor, bufRest            *ng.MeshBuffer
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

	//	meshes / models
	if bufFloor, err = ng.Core.MeshBuffers.Add("buf_floor", ng.Core.MeshBuffers.NewParams(6, 6)); err != nil {
		panic(err)
	}
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(36+12, 36+12)); err != nil {
		panic(err)
	}
	if meshFloor, err = ng.Core.Libs.Meshes.AddLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}
	if meshPyr, err = ng.Core.Libs.Meshes.AddLoad("mesh_pyramid", ng.MeshProviderPrefabPyramid); err != nil {
		panic(err)
	}
	if meshCube, err = ng.Core.Libs.Meshes.AddLoad("mesh_cube", ng.MeshProviderPrefabCube); err != nil {
		panic(err)
	}
	bufFloor.Add(meshFloor)
	bufRest.Add(meshCube)
	bufRest.Add(meshPyr)
	meshPyr.Models.Default().MatID = apputil.LibIDs.Mat["mosaic"]
	meshCube.Models.Default().MatID = apputil.LibIDs.Mat["crate"]

	fx := &ng.Core.Libs.Effects[apputil.LibIDs.Fx["cat"]]
	fx.Ops.EnableOrangify(-1).SetMixWeight(0.5)
	fx.UpdateRoutine()
	fx = &ng.Core.Libs.Effects[apputil.LibIDs.Fx["dog"]]
	fx.Ops.EnableTex2D(1).SetImageID(apputil.LibIDs.Img["gopher"]).SetMixWeight(0.5)
	fx.UpdateRoutine()
	ng.Core.Libs.Materials[apputil.LibIDs.Mat["crate"]].FaceEffects.ByTag["front"] = apputil.LibIDs.Fx["dog"]
	ng.Core.Libs.Materials[apputil.LibIDs.Mat["crate"]].FaceEffects.ByTag["back"] = apputil.LibIDs.Fx["dog"]
	ng.Core.Libs.Materials[apputil.LibIDs.Mat["mosaic"]].FaceEffects.ByID["t3"] = apputil.LibIDs.Fx["cat"]

	//	scene
	scene = apputil.AddMainScene()
	apputil.AddSkyMesh(scene, "mesh_cube")
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	pyr = scene.RootNode.ChildNodes.AddNew("node_pyr", "mesh_pyramid", "")
	box = scene.RootNode.ChildNodes.AddNew("node_box", "mesh_cube", "")

	floor.MatID = apputil.LibIDs.Mat["cobbles"]
	floor.Transform.SetPos(0.1, 0, -8)
	floor.Transform.SetScale(1000)
	floor.ApplyTransform()

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.X, camCtl.Pos.Y, camCtl.Pos.Z = -1, 1.6, -2
	camCtl.EndUpdate()
}
