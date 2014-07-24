// Renders a scene of many animated 3D objects surrounded by a sky-pyramid, adding a 2nd smaller (tinted & pulsating) "look-behind" camera.
package main

import (
	"math"

	"github.com/go-utils/unum"
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
	ng "github.com/go3d/go-ngine/core"
)

var (
	gui2d apputil.Gui2D

	tmpScene                                                     *ng.Scene
	tmpNode                                                      *ng.SceneNode
	floorID, pyrID, boxID, meshCubeID, meshPyrID, modelCubeCatID int
	pyrIDs                                                       [4]int

	crateIDs         = make([]int, 0, 9)
	initialCrateAdds = 1
)

func main() {
	apputil.AddKeyHint("F10", "Add three more crates")
	apputil.AddKeyHint("F11", "Remove last three crates")
	apputil.AddKeyHint("F12", "Toggle 'Rear-View-Mirror' Camera")
	// ng.Options.Rendering.DefaultBatcher.Enabled = false
	// ng.Loop.MaxIterations = 2000
	apputil.Main(setupExample_04_PyrsCubes, onAppThread, onWinThread)
}

func onWinThread() {
	apputil.CheckCamCtlKeys()
	apputil.CheckAndHandleToggleKeys()
	if ng.UserIO.KeyToggled(apputil.KeyF12) {
		apputil.RearView.Toggle()
	}
	if ng.UserIO.KeyToggled(apputil.KeyF10) {
		addCrates(apputil.SceneCam.Scene(), 3)
	}
	if ng.UserIO.KeyToggled(apputil.KeyF11) {
		removeCrates(apputil.SceneCam.Scene(), 3)
	}
	apputil.RearView.OnWin()

	//	pulsating fx anims
	ng.Core.Libs.Effects[apputil.LibIDs.Fx["mosaic"]].GetTex2D(1).SetMixWeight(0.5 + (0.5 * math.Cos(ng.Loop.Tick.Now*2)))
	apputil.RearView.View.FxProcs.GetOrangify(0).SetMixWeight(0.75 + (0.25 * math.Sin(ng.Loop.Tick.Now*4)))
}

func onAppThread() {
	if apputil.Paused {
		return
	}

	apputil.HandleCamCtlKeys()
	apputil.RearView.OnApp()

	//	node geometry anims
	if tmpScene = gui2d.Cam.Scene(); tmpScene != nil {
		tmpScene.Node(gui2d.DogNodeID).Transform.Rot.Add3(0, -0.005, 0.001)
		tmpScene.Node(gui2d.CatNodeID).Transform.Rot.X += 0.003
		tmpScene.ApplyNodeTransforms(0)
	}

	if tmpScene = apputil.SceneCam.Scene(); tmpScene != nil {
		if tmpNode = tmpScene.Node(pyrID); tmpNode != nil {
			tmpNode.Transform.Rot.Add3(-0.0005, -0.0005, 0)
			tmpNode.Transform.Pos.Set(-13.75, 2*math.Sin(ng.Loop.Tick.Now), 2)
		}
		if tmpNode = tmpScene.Node(boxID); tmpNode != nil {
			tmpNode.Transform.Rot.Add3(0.0004, 0, 0.0006)
			tmpNode.Transform.Pos.Set(-8.125, 2*math.Cos(ng.Loop.Tick.Now), -2)
		}

		var f float64
		for i := 0; i < len(crateIDs); i++ {
			if tmpNode = tmpScene.Node(crateIDs[i]); tmpNode != nil {
				f = float64(i)
				f = (f + 1) * (f + 1)
				tmpNode.Transform.Rot.Add3(f*0.00001, f*0.0001, f*0.001)
			}
		}

		if tmpNode = tmpScene.Node(pyrIDs[0]); tmpNode != nil {
			tmpNode.Transform.Pos.X = math.Sin(ng.Loop.Tick.Now) * 100
		}
		if tmpNode = tmpScene.Node(pyrIDs[1]); tmpNode != nil {
			tmpNode.Transform.Pos.Z = math.Cos(ng.Loop.Tick.Now) * 1000
		}
		tmpScene.ApplyNodeTransforms(0)
	}
}

func setupExample_04_PyrsCubes() {
	var (
		err               error
		meshPlaneID       int
		bufFloor, bufRest *ng.MeshBuffer
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
	fx.EnableTex2D(1).Tex_SetImageID(apputil.LibIDs.Img2D["gopher"]).SetMixWeight(0.5)
	fx.UpdateRoutine()

	//	meshes / models
	if bufFloor, err = ng.Core.Mesh.Buffers.AddNew("buf_floor", 6); err != nil {
		panic(err)
	}
	if bufRest, err = ng.Core.Mesh.Buffers.AddNew("buf_rest", 36+12); err != nil {
		panic(err)
	}

	if meshPlaneID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_plane", ng.Core.Mesh.Desc.Plane); err != nil {
		panic(err)
	}
	modelPlaneDefaultID := ng.Core.Libs.Models.AddNew()
	ng.Core.Libs.Models[modelPlaneDefaultID].MatID = apputil.LibIDs.Mat["cobbles"]
	ng.Core.Libs.Meshes[meshPlaneID].DefaultModelID = modelPlaneDefaultID

	if meshPyrID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_pyramid", ng.Core.Mesh.Desc.Pyramid); err != nil {
		panic(err)
	}
	modelPyrDefaultID := ng.Core.Libs.Models.AddNew()
	ng.Core.Libs.Models[modelPyrDefaultID].MatID = apputil.LibIDs.Mat["mosaic"]
	ng.Core.Libs.Meshes[meshPyrID].DefaultModelID = modelPyrDefaultID
	modelPyrDogID := ng.Core.Libs.Models.AddNew()
	ng.Core.Libs.Models[modelPyrDogID].MatID = apputil.LibIDs.Mat["dog"]

	if meshCubeID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_cube", ng.Core.Mesh.Desc.Cube); err != nil {
		panic(err)
	}
	modelCubeDefaultID := ng.Core.Libs.Models.AddNew()
	ng.Core.Libs.Models[modelCubeDefaultID].MatID = apputil.LibIDs.Mat["crate"]
	ng.Core.Libs.Meshes[meshCubeID].DefaultModelID = modelCubeDefaultID
	modelCubeCatID = ng.Core.Libs.Models.AddNew()
	ng.Core.Libs.Models[modelCubeCatID].MatID = apputil.LibIDs.Mat["cat"]

	tmpMat := &ng.Core.Libs.Materials[apputil.LibIDs.Mat["crate"]]
	tmpMat.FaceEffects.ByTag["front"] = apputil.LibIDs.Fx["cat"]
	tmpMat.FaceEffects.ByTag["back"] = apputil.LibIDs.Fx["dog"]

	tmpMat = &ng.Core.Libs.Materials[ng.Core.Libs.Materials.AddNew()]
	tmpMat.DefaultEffectID = apputil.LibIDs.Fx["crate"]
	tmpMat.FaceEffects.ByTag["front"], tmpMat.FaceEffects.ByTag["back"] = apputil.LibIDs.Fx["cat"], apputil.LibIDs.Fx["cat"]
	tmpMat.FaceEffects.ByTag["top"], tmpMat.FaceEffects.ByTag["bottom"] = apputil.LibIDs.Fx["dog"], apputil.LibIDs.Fx["dog"]
	tmpMat.FaceEffects.ByTag["left"], tmpMat.FaceEffects.ByTag["right"] = apputil.LibIDs.Fx["gopher"], apputil.LibIDs.Fx["gopher"]
	apputil.LibIDs.Mat["mix"] = tmpMat.ID

	bufFloor.Add(meshPlaneID)
	bufRest.Add(meshCubeID)
	bufRest.Add(meshPyrID)

	//	scene
	scene := apputil.AddMainScene()
	apputil.AddSkyMesh(scene, meshPyrID)
	floor := apputil.AddNode(scene, 0, meshPlaneID, -1, -1)
	floorID = floor.ID
	floor.Transform.SetPos(0.1, 0, -8)
	floor.Transform.SetScale(10000)

	pyrID = apputil.AddNode(scene, 0, meshPyrID, -1, -1).ID
	boxID = apputil.AddNode(scene, 0, meshCubeID, -1, -1).ID
	for i := 0; i < initialCrateAdds; i++ {
		addCrates(scene, 3)
	}
	var f float64
	for i := 0; i < len(pyrIDs); i++ {
		tmpNode = apputil.AddNode(scene, 0, meshPyrID, -1, -1)
		pyrIDs[i] = tmpNode.ID
		if i > 1 {
			tmpNode.Render.ModelID = modelPyrDogID
		}
		f = float64(len(pyrIDs) - i)
		tmpNode.Transform.SetScale((f + 1) * 2)
		tmpNode.Transform.SetPos((f+3)*-4, (f+2)*3, (f+2)*14)
		if i > 1 {
			if i == 2 {
				f = 45
			} else {
				f = 135
			}
			tmpNode.Transform.Rot.Z = unum.DegToRad(f)
		} else {
			if i == 0 {
				f = 180
			} else {
				f = 90
			}
			tmpNode.Transform.Rot.X = unum.DegToRad(f)
		}
		if i == 1 {
			tmpNode.Transform.SetScale(100)
			tmpNode.Transform.Pos.Y += 100
		}
	}

	scene.ApplyNodeTransforms(0)
	apputil.RearView.Setup(scene.ID)
	if err = gui2d.Setup(); err != nil {
		panic(err)
	}

	camCtl := &apputil.SceneCam.Controller
	camCtl.BeginUpdate()
	camCtl.Pos.X, camCtl.Pos.Y, camCtl.Pos.Z = -3.57, 3.63, 19.53
	camCtl.TurnRightBy(-111)
	camCtl.EndUpdate()
}

func addCrates(scene *ng.Scene, num int) {
	var f float64
	oldLen := len(crateIDs)
	nuLen := oldLen + 3
	if nuLen > cap(crateIDs) {
		nu := make([]int, len(crateIDs), len(crateIDs)*2)
		copy(nu, crateIDs)
		crateIDs = nu
	}
	for i := 0; i < num; i++ {
		tmpNode = apputil.AddNode(scene, 0, meshCubeID, -1, -1)
		crateIDs = append(crateIDs, tmpNode.ID)
		f = float64(oldLen + i)
		tmpNode.Transform.SetPos((f+3)*-2, (f+1)*2, (f+2)*3)
		switch i {
		case 0:
			tmpNode.Render.MatID = apputil.LibIDs.Mat["mix"]
		case 1:
			apputil.AddNode(scene, tmpNode.ID, meshPyrID, apputil.LibIDs.Mat["cat"], -1).Transform.Pos.Y = 2.125
		case 2:
			tmpNode.Render.ModelID = modelCubeCatID
		}
	}
	scene.ApplyNodeTransforms(0)
}

func removeCrates(scene *ng.Scene, num int) {
	if len(crateIDs) >= num {
		for i := len(crateIDs) - num; i < len(crateIDs); i++ {
			scene.RemoveNode(crateIDs[i])
		}
		crateIDs = crateIDs[:len(crateIDs)-num]
	}
}
