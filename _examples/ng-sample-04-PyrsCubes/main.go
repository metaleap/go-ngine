package main

import (
	"fmt"
	"math"

	glfw "github.com/go-gl/glfw"
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
	unum "github.com/metaleap/go-util/num"
)

var (
	gui2d apputil.Gui2D

	floor, pyr, box *ng.Node
	crates          [3]*ng.Node
	pyramids        [4]*ng.Node
	i               int
	f               float64
)

func main() {
	apputil.AddKeyHint("F12", "Toggle 'Rear-View-Mirror' Camera")
	apputil.Main(setupExample_04_PyrsCubes, onAppThread, onWinThread)
}

func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()
	if ng.UserIO.KeyToggled(glfw.KeyF12) {
		apputil.RearView.Toggle()
	}
	apputil.RearView.OnWin()

	//	pulsating fx anims
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["mosaic"]].Ops.GetTex2D(1).SetMixWeight(0.5 + (0.5 * math.Cos(ng.Loop.Tick.Now*2)))
	apputil.RearView.Cam.Rendering.FxOps.GetOrangify(0).SetMixWeight(0.75 + (0.25 * math.Sin(ng.Loop.Tick.Now*4)))
}

func onAppThread() {
	if apputil.Paused {
		return
	}

	apputil.HandleCamCtlKeys()
	apputil.RearView.OnApp()

	//	node geometry anims
	gui2d.Dog.Transform.Rot.Add3(0, -0.005, 0.001)
	gui2d.Dog.ApplyTransform()
	gui2d.Cat.Transform.Rot.X += 0.003
	gui2d.Cat.ApplyTransform()

	pyr.Transform.Rot.Add3(-0.0005, -0.0005, 0)
	pyr.Transform.Pos.Set(-13.75, 2*math.Sin(ng.Loop.Tick.Now), 2)
	pyr.ApplyTransform()

	box.Transform.Rot.Add3(0.0004, 0, 0.0006)
	box.Transform.Pos.Set(-8.125, 2*math.Cos(ng.Loop.Tick.Now), -2)
	box.ApplyTransform()

	for i = 0; i < len(crates); i++ {
		f = float64(i)
		f = (f + 1) * (f + 1)
		crates[i].Transform.Rot.Add3(f*0.00001, f*0.0001, f*0.001)
		crates[i].ApplyTransform()
	}

	pyramids[0].Transform.Pos.X = math.Sin(ng.Loop.Tick.Now) * 100
	pyramids[0].ApplyTransform()
	pyramids[1].Transform.Pos.Z = math.Cos(ng.Loop.Tick.Now) * 1000
	pyramids[1].ApplyTransform()
}

func setupExample_04_PyrsCubes() {
	var (
		err                          error
		meshPlane, meshPyr, meshCube *ng.Mesh
		bufFloor, bufRest            *ng.MeshBuffer
		scene                        *ng.Scene
	)

	//	textures / materials
	apputil.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"crate":   "tex/crate.jpeg",
		"mosaic":  "tex/mosaic.jpeg",
		"cat":     "tex/cat.png",
		"dog":     "tex/dog.png",
		"gopher":  "tex/gopher.png",
	})

	fx := &ng.Core.Libs.Effects[apputil.LibIDs.Fx["mosaic"]]
	fx.Ops.EnableTex2D(1).SetImageID(apputil.LibIDs.Img["gopher"]).SetMixWeight(0.5)
	fx.UpdateRoutine()

	//	meshes / models
	if bufFloor, err = ng.Core.MeshBuffers.Add("buf_floor", ng.Core.MeshBuffers.NewParams(6, 6)); err != nil {
		panic(err)
	}
	if bufRest, err = ng.Core.MeshBuffers.Add("buf_rest", ng.Core.MeshBuffers.NewParams(36+12, 36+12)); err != nil {
		panic(err)
	}

	if meshPlane, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.MeshProviderPrefabPlane); err != nil {
		panic(err)
	}
	modelPlaneDefault := ng.Core.Libs.Models.AddNew()
	modelPlaneDefault.MatID = apputil.LibIDs.Mat["cobbles"]
	meshPlane.DefaultModelID = modelPlaneDefault.ID

	if meshPyr, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_pyramid", ng.MeshProviderPrefabPyramid); err != nil {
		panic(err)
	}
	modelPyrDefault := ng.Core.Libs.Models.AddNew()
	modelPyrDefault.MatID = apputil.LibIDs.Mat["mosaic"]
	meshPyr.DefaultModelID = modelPyrDefault.ID
	modelPyrDog := ng.Core.Libs.Models.AddNew()
	modelPyrDog.MatID = apputil.LibIDs.Mat["dog"]

	if meshCube, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_cube", ng.MeshProviderPrefabCube); err != nil {
		panic(err)
	}
	modelCubeDefault := ng.Core.Libs.Models.AddNew()
	modelCubeDefault.MatID = apputil.LibIDs.Mat["crate"]
	meshCube.DefaultModelID = modelCubeDefault.ID
	modelCubeCat := ng.Core.Libs.Models.AddNew()
	modelCubeCat.MatID = apputil.LibIDs.Mat["cat"]

	tmpMat := &ng.Core.Libs.Materials[apputil.LibIDs.Mat["crate"]]
	tmpMat.FaceEffects.ByTag["front"] = apputil.LibIDs.Fx["cat"]
	tmpMat.FaceEffects.ByTag["back"] = apputil.LibIDs.Fx["dog"]

	tmpMat = ng.Core.Libs.Materials.AddNew()
	tmpMat.DefaultEffectID = apputil.LibIDs.Fx["crate"]
	tmpMat.FaceEffects.ByTag["front"], tmpMat.FaceEffects.ByTag["back"] = apputil.LibIDs.Fx["cat"], apputil.LibIDs.Fx["cat"]
	tmpMat.FaceEffects.ByTag["top"], tmpMat.FaceEffects.ByTag["bottom"] = apputil.LibIDs.Fx["dog"], apputil.LibIDs.Fx["dog"]
	tmpMat.FaceEffects.ByTag["left"], tmpMat.FaceEffects.ByTag["right"] = apputil.LibIDs.Fx["gopher"], apputil.LibIDs.Fx["gopher"]
	apputil.LibIDs.Mat["mix"] = tmpMat.ID

	bufFloor.Add(meshPlane)
	bufRest.Add(meshCube)
	bufRest.Add(meshPyr)

	//	scene
	scene = apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshPyr.ID)
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", meshPlane.ID)
	pyr = scene.RootNode.ChildNodes.AddNew("node_pyr", meshPyr.ID)
	box = scene.RootNode.ChildNodes.AddNew("node_box", meshCube.ID)

	for i = 0; i < len(crates); i++ {
		crates[i] = scene.RootNode.ChildNodes.AddNew(fmt.Sprintf("node_box_%v", i), meshCube.ID)
		if i == 0 {
			crates[i].ModelID = modelCubeCat.ID
		}
		f = float64(i)
		crates[i].Transform.SetPos((f+3)*-2, (f+1)*2, (f+2)*3)
		crates[i].ApplyTransform()
		if i == 2 {
			crates[i].MatID = apputil.LibIDs.Mat["mix"]
		}
	}

	for i = 0; i < len(pyramids); i++ {
		pyramids[i] = scene.RootNode.ChildNodes.AddNew(fmt.Sprintf("nody_pyr_%v", i), meshPyr.ID)
		if i > 1 {
			pyramids[i].ModelID = modelPyrDog.ID
		}
		f = float64(len(pyramids) - i)
		pyramids[i].Transform.SetScale((f + 1) * 2)
		pyramids[i].Transform.SetPos((f+3)*-4, (f+2)*3, (f+2)*14)
		if i > 1 {
			if i == 2 {
				f = 45
			} else {
				f = 135
			}
			pyramids[i].Transform.Rot.Z = unum.DegToRad(f)
		} else {
			if i == 0 {
				f = 180
			} else {
				f = 90
			}
			pyramids[i].Transform.Rot.X = unum.DegToRad(f)
		}
		if i == 1 {
			pyramids[i].Transform.SetScale(100)
			pyramids[i].Transform.Pos.Y += 100
		}
		pyramids[i].ApplyTransform()
	}

	floor.Transform.SetPos(0.1, 0, -8)
	floor.Transform.SetScale(10000)
	floor.ApplyTransform()

	if err = gui2d.Setup(); err != nil {
		panic(err)
	}
	apputil.RearView.Setup(scene.ID)

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.X, camCtl.Pos.Y, camCtl.Pos.Z = -3.57, 3.63, 19.53
	camCtl.TurnRightBy(155)
	camCtl.EndUpdate()
}
