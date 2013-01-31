package main

import (
	"fmt"
	"math"

	glfw "github.com/go-gl/glfw"
	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
	ng "github.com/go3d/go-ngine/core"
	unum "github.com/metaleap/go-util/num"
)

var (
	gui2d      ngsamples.Gui2D
	rearMirror ngsamples.RearMirror

	floor, pyr, box *ng.Node
	crates          [3]*ng.Node
	pyramids        [4]*ng.Node
	i               int
	f               float64
)

func main() {
	ngsamples.KeyHints = append([]string{"[F4]  --  Toggle Rear-Mirror Camera"}, ngsamples.KeyHints...)
	ngsamples.Keys.CheckFor.Toggled = append(ngsamples.Keys.CheckFor.Toggled, glfw.KeyF4)
	ngsamples.SamplesMainFunc(LoadSampleScene_03_PyrsCubes)
}

func onWinThread() {
	ngsamples.CheckCamCtlKeys()
	ngsamples.CheckToggleKeys()
	ngsamples.HandleToggleKeys()
	if ngsamples.Keys.Toggled[glfw.KeyF4] {
		rearMirror.Cam.Enabled = !rearMirror.Cam.Enabled
	}
	rearMirror.Cam.Rendering.States.FaceCulling = ngsamples.Cam.Rendering.States.FaceCulling
}

func onAppThread() {
	ngsamples.HandleCamCtlKeys()
	rearMirror.OnApp()

	//	animate mesh nodes
	pyr.Transform.Rot.X -= 0.0005
	pyr.Transform.Rot.Y -= 0.0005
	pyr.Transform.Pos.Set(-13.75, 2*math.Sin(ng.Loop.TickNow), 2)
	pyr.Transform.OnPosRotChanged()

	gui2d.Dog.Transform.Rot.Y -= 0.005
	gui2d.Dog.Transform.Rot.Z += 0.001
	gui2d.Dog.Transform.OnRotChanged()
	gui2d.Cat.Transform.Rot.X += 0.003
	gui2d.Cat.Transform.OnRotXChanged()

	box.Transform.Rot.Y += 0.0004
	box.Transform.Rot.Z += 0.0006
	box.Transform.Pos.Set(-8.125, 2*math.Cos(ng.Loop.TickNow), -2)
	box.Transform.OnPosRotChanged()

	for i = 0; i < len(crates); i++ {
		f = float64(i)
		f = (f + 1) * (f + 1)
		crates[i].Transform.Rot.X += f * 0.00001
		crates[i].Transform.Rot.Y += f * 0.0001
		crates[i].Transform.Rot.Z += f * 0.001
		crates[i].Transform.OnRotChanged()
	}

	pyramids[0].Transform.SetPosX(math.Sin(ng.Loop.TickNow) * 100)
	pyramids[1].Transform.SetPosZ(math.Cos(ng.Loop.TickNow) * 1000)
}

func LoadSampleScene_03_PyrsCubes() {
	var (
		err                          error
		scene                        *ng.Scene
		meshFloor, meshPyr, meshCube *ng.Mesh
		bufFloor, bufRest            *ng.MeshBuffer
		str                          string
	)

	ng.Loop.OnAppThread, ng.Loop.OnWinThread = onAppThread, onWinThread
	ngsamples.Cam.Rendering.States.FaceCulling = false
	rearMirror.Setup()

	//	textures / materials
	ngsamples.AddTextureMaterials(map[string]string{
		"cobbles": "tex/cobbles.png",
		"crate":   "tex/crate.jpeg",
		"mosaic":  "tex/mosaic.jpeg",
		"cat":     "tex/cat.png",
		"dog":     "tex/dog.png",
	})
	ngsamples.AddColorMaterials(map[string][]float64{
		"yellow": []float64{1, 1, 0},
		"green":  []float64{0, 1, 0},
		"blue":   []float64{0, 0, 1},
	})

	if err = gui2d.Setup(); err != nil {
		panic(err)
	}

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
	meshFloor.Models.Default().SetMatID("mat_cobbles")
	if meshPyr, err = ng.Core.Libs.Meshes.AddLoad("mesh_pyramid", ng.MeshProviderPrefabPyramid); err != nil {
		panic(err)
	}
	meshPyr.Models.Default().SetMatID("mat_mosaic")
	meshPyr.Models.Default().Clone("model_pyramid_dog").SetMatID("mat_dog")

	if meshCube, err = ng.Core.Libs.Meshes.AddLoad("mesh_cube", ng.MeshProviderPrefabCube); err != nil {
		panic(err)
	}
	meshCube.Models.Default().SetMatID("mat_crate")
	meshCube.Models.Default().Clone("model_cube_cat").SetMatID("mat_cat")

	bufFloor.Add(meshFloor)
	bufRest.Add(meshCube)
	bufRest.Add(meshPyr)

	//	scene
	scene = ngsamples.AddScene("")
	floor = scene.RootNode.ChildNodes.AddNew("node_floor", "mesh_plane", "")
	pyr = scene.RootNode.ChildNodes.AddNew("node_pyr", "mesh_pyramid", "")
	box = scene.RootNode.ChildNodes.AddNew("node_box", "mesh_cube", "")

	for i = 0; i < len(crates); i++ {
		if i == 0 {
			str = "model_cube_cat"
		} else {
			str = ""
		}
		crates[i] = scene.RootNode.ChildNodes.AddNew(fmt.Sprintf("node_box_%v", i), "mesh_cube", str)
		f = float64(i)
		crates[i].Transform.SetPosXYZ((f+3)*-2, (f+1)*2, (f+2)*3)
		if i == 2 {
			crates[i].SetMatID("mat_dog")
		}
	}

	for i = 0; i < len(pyramids); i++ {
		if i > 1 {
			str = "model_pyramid_dog"
		} else {
			str = ""
		}
		pyramids[i] = scene.RootNode.ChildNodes.AddNew(fmt.Sprintf("nody_pyr_%v", i), "mesh_pyramid", str)
		f = float64(len(pyramids) - i)
		pyramids[i].Transform.SetScaleN((f + 1) * 2)
		pyramids[i].Transform.SetPosXYZ((f+3)*-4, (f+2)*3, (f+2)*14)
		if i > 1 {
			if i == 2 {
				f = 45
			} else {
				f = 135
			}
			pyramids[i].Transform.SetRotZ(unum.DegToRad(f))
		} else {
			if i == 0 {
				f = 180
			} else {
				f = 90
			}
			pyramids[i].Transform.SetRotX(unum.DegToRad(f))
		}
		if i == 1 {
			pyramids[i].Transform.SetScaleN(100)
			pyramids[i].Transform.Pos.Y += 100
		}
	}

	// pyramids[len(pyramids)-1].SetMatID("mat_yellow")
	floor.Transform.SetPosXYZ(0.1, 0, -8)
	floor.Transform.SetScaleN(10000)

	ngsamples.CamCtl.BeginUpdate()
	ngsamples.CamCtl.Pos.X, ngsamples.CamCtl.Pos.Y, ngsamples.CamCtl.Pos.Z = 35, 1.6, 24
	ngsamples.CamCtl.EndUpdate()
}
