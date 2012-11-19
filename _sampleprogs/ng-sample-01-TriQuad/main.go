package main

import (
	"math"

	ng "github.com/go3d/go-ngine/core"
	ngsamples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	tri, quad *ng.Node
)

func main () {
	ngsamples.MaxKeyHint = 3
	ngsamples.SamplesMainFunc(LoadSampleScene_01_TriQuad)
}

func onLoop () {
	ngsamples.CheckToggleKeys()
	tri.Transform.Rot.X -= 0.005
	tri.Transform.Rot.Y -= 0.005
	tri.Transform.Pos.Set(-3.75, 1 * math.Sin(ng.Loop.TickNow), 1)
	tri.Transform.OnPosRotChanged()
	quad.Transform.Rot.Y += 0.0001
	quad.Transform.Rot.Z += 0.0001
	quad.Transform.Pos.Set(-4.125, 1 * math.Cos(ng.Loop.TickNow), 0)
	quad.Transform.OnPosRotChanged()
}

func LoadSampleScene_01_TriQuad () {
	var meshTri, meshQuad *ng.Mesh
	var meshBuf *ng.MeshBuffer
	var err error

	ng.Loop.OnLoop = onLoop
	ngsamples.Cam.Options.BackfaceCulling = false

	//	textures / materials

	ng.Core.Textures["tex_cat"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/cat.png")
	ng.Core.Textures["tex_dog"] = ng.Core.Textures.Load(ng.TextureProviders.LocalFile, "tex/dog.png")

	ng.Core.Materials["mat_cat"] = ng.Core.Materials.New("tex_cat")
	ng.Core.Materials["mat_dog"] = ng.Core.Materials.New("tex_dog")

	//	meshes / models

	if meshTri, err = ng.Core.Meshes.Load("mesh_tri", ng.MeshProviders.PrefabTri); err != nil { panic(err) }
	if meshQuad, err = ng.Core.Meshes.Load("mesh_quad", ng.MeshProviders.PrefabQuad); err != nil { panic(err) }
	ng.Core.Meshes.AddRange(meshTri, meshQuad)

	if meshBuf, err = ng.Core.MeshBuffers.Add("meshbuf", ng.Core.MeshBuffers.NewParams(9, 9)); err != nil { panic(err) }
	if err = meshBuf.Add(meshTri); err != nil { panic(err) }
	if err = meshBuf.Add(meshQuad); err != nil { panic(err) }

	//	scene
	var scene = ng.NewScene()
	ng.Core.Scenes[""] = scene
	scene.RootNode.SubNodes.MakeN("node_tri", "mesh_tri", "", "node_quad", "mesh_quad", "")
	tri, quad = scene.RootNode.SubNodes.M["node_tri"], scene.RootNode.SubNodes.M["node_quad"]
	tri.SetMatName("mat_cat")
	quad.SetMatName("mat_dog")

	//	upload everything
	ng.Core.SyncUpdates()
}
