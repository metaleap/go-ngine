package main

import (
	"math"

	ngine "github.com/go3d/go-ngine/core"
	ngine_samples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	tri, quad *ngine.Node
)

func main () {
	ngine_samples.MaxKeyHint = 3
	ngine_samples.SamplesMainFunc(LoadSampleScene_01_TriQuad)
}

func onLoop () {
	ngine_samples.CheckToggleKeys()
	tri.Transform.Rot.X -= 0.005
	tri.Transform.Rot.Y -= 0.005
	tri.Transform.Pos.Set(-3.75, 1 * math.Sin(ngine.Loop.TickNow), 1)
	tri.Transform.OnPosRotChanged()
	quad.Transform.Rot.Y += 0.0001
	quad.Transform.Rot.Z += 0.0001
	quad.Transform.Pos.Set(-4.125, 1 * math.Cos(ngine.Loop.TickNow), 0)
	quad.Transform.OnPosRotChanged()
}

func LoadSampleScene_01_TriQuad () {
	var meshTri, meshQuad *ngine.Mesh
	var meshBuf *ngine.MeshBuffer
	var err error

	ngine.Loop.OnLoop = onLoop
	ngine.Core.Canvases[0].Cameras[0].Options.BackfaceCulling = false

	//	textures / materials

	ngine.Core.Textures["tex_cat"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "tex/cat.png")
	ngine.Core.Textures["tex_dog"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "tex/dog.png")

	ngine.Core.Materials["mat_cat"] = ngine.Core.Materials.New("tex_cat")
	ngine.Core.Materials["mat_dog"] = ngine.Core.Materials.New("tex_dog")

	//	meshes / models

	if meshTri, err = ngine.Core.Meshes.Load("mesh_tri", ngine.MeshProviders.PrefabTri); err != nil { panic(err) }
	if meshQuad, err = ngine.Core.Meshes.Load("mesh_quad", ngine.MeshProviders.PrefabQuad); err != nil { panic(err) }
	ngine.Core.Meshes.AddRange(meshTri, meshQuad)

	if meshBuf, err = ngine.Core.MeshBuffers.Add("meshbuf", ngine.Core.MeshBuffers.NewParams(9, 9)); err != nil { panic(err) }
	if err = meshBuf.Add(meshTri); err != nil { panic(err) }
	if err = meshBuf.Add(meshQuad); err != nil { panic(err) }

	//	scene
	var scene = ngine.NewScene()
	ngine.Core.Scenes[""] = scene
	scene.RootNode.SubNodes.MakeN("node_tri", "mesh_tri", "", "node_quad", "mesh_quad", "")
	tri, quad = scene.RootNode.SubNodes.M["node_tri"], scene.RootNode.SubNodes.M["node_quad"]
	tri.SetMatName("mat_cat")
	quad.SetMatName("mat_dog")

	//	upload everything
	ngine.Core.SyncUpdates()
}
