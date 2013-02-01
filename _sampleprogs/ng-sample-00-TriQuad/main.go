package main

import (
	"math"

	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
	ng "github.com/go3d/go-ngine/core"
)

var (
	tri, quad *ng.Node
)

func main() {
	ngsamples.MaxKeyHint = 3
	ngsamples.SamplesMainFunc(LoadSampleScene_00_TriQuad)
}

func onWinThread() {
	ngsamples.CheckToggleKeys()
	ngsamples.HandleToggleKeys()
}

func onAppThread() {
	if !ngsamples.Paused {
		tri.Transform.Rot.Add3(-0.005, -0.005, 0)
		tri.Transform.Pos.Set(-3.75, 1*math.Sin(ng.Loop.TickNow), 1)
		tri.Transform.ApplyMatrices()
		quad.Transform.Rot.Add3(0, 0.001, 0.001)
		quad.Transform.Pos.Set(-4.125, 1*math.Cos(ng.Loop.TickNow), 0)
		quad.Transform.ApplyMatrices()
	}
}

func LoadSampleScene_00_TriQuad() {
	var (
		err               error
		scene             *ng.Scene
		meshTri, meshQuad *ng.Mesh
		meshBuf           *ng.MeshBuffer
	)

	ng.Loop.OnAppThread, ng.Loop.OnWinThread = onAppThread, onWinThread
	ngsamples.Cam.Rendering.States.FaceCulling = false

	//	textures / materials
	ngsamples.AddTextureMaterials(map[string]string{
		"cat": "tex/cat.png",
		"dog": "tex/dog.png",
	})

	//	meshes / models

	if meshTri, err = ng.Core.Libs.Meshes.AddLoad("mesh_tri", ng.MeshProviderPrefabTri); err != nil {
		panic(err)
	}
	if meshQuad, err = ng.Core.Libs.Meshes.AddLoad("mesh_quad", ng.MeshProviderPrefabQuad); err != nil {
		panic(err)
	}

	if meshBuf, err = ng.Core.MeshBuffers.Add("meshbuf", ng.Core.MeshBuffers.NewParams(9, 9)); err != nil {
		panic(err)
	}
	if err = meshBuf.Add(meshTri); err != nil {
		panic(err)
	}
	if err = meshBuf.Add(meshQuad); err != nil {
		panic(err)
	}

	//	scene
	scene = ngsamples.AddScene("", true)
	tri = scene.RootNode.ChildNodes.AddNew("node_tri", "mesh_tri", "")
	quad = scene.RootNode.ChildNodes.AddNew("node_quad", "mesh_quad", "")
	tri.SetMatID("mat_cat")
	quad.SetMatID("mat_dog")
}
