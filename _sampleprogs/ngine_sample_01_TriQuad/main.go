package main

import (
	"math"

	ngine "github.com/go3d/go-ngine/core"
	ngine_samples "github.com/go3d/go-ngine/_sampleprogs/_sharedcode"
)

var (
	tri, quad *ngine.TNode
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
	var meshTri, meshQuad *ngine.TMesh
	var bufTri, bufQuad *ngine.TMeshBuffer
	var err error

	ngine.Loop.OnLoop = onLoop
	ngine.Core.Options.SetGlBackfaceCulling(false)

	//	textures / materials
	ngine.Core.Textures["tex_cat"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "misc/cat.png")
	ngine.Core.Textures["tex_dog"] = ngine.Core.Textures.Load(ngine.TextureProviders.LocalFile, "misc/dog.png")
	ngine.Core.Materials["mat_cat"] = ngine.Core.Materials.New("tex_cat")
	ngine.Core.Materials["mat_dog"] = ngine.Core.Materials.New("tex_dog")

	//	meshes / models
	if bufTri, err = ngine.Core.MeshBuffers.Add("buf_tri", ngine.Core.MeshBuffers.NewParams(3, 3)); err != nil { panic(err) }
	if bufQuad, err = ngine.Core.MeshBuffers.Add("buf_quad", ngine.Core.MeshBuffers.NewParams(6, 6)); err != nil { panic(err) }
	if meshTri, err = ngine.Core.Meshes.Load("mesh_tri", ngine.MeshProviders.PrefabTri); err != nil { panic(err) }
	if meshQuad, err = ngine.Core.Meshes.Load("mesh_quad", ngine.MeshProviders.PrefabQuad); err != nil { panic(err) }
	ngine.Core.Meshes.AddRange(meshTri, meshQuad)
	bufTri.Add(meshTri); bufQuad.Add(meshQuad)

	//	scene
	var scene = ngine.NewScene()
	ngine.Core.Scenes[""] = scene
	scene.RootNode.AddSubNodesNamed(map[string]string { "node_tri":  "mesh_tri", "node_quad": "mesh_quad" })
	tri, quad = scene.RootNode.SubNodes["node_tri"], scene.RootNode.SubNodes["node_quad"]
	tri.SetMatKey("mat_cat")
	quad.SetMatKey("mat_dog")

	//	upload everything
	ngine.Core.SyncUpdates()
}
