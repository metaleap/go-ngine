package samplescenes

import (
	ng "github.com/go3d/go-ngine/core"
	unum "github.com/metaleap/go-util/num"
)

type Gui2D struct {
	Cat, Dog *ng.Node
}

func (me *Gui2D) Setup() (err error) {
	var (
		meshBuf  *ng.MeshBuffer
		quadMesh *ng.Mesh
	)
	scene := AddScene("gui2d", false)
	cam := ng.Core.Rendering.Canvases[0].AddNewCamera2D(true)
	cam.Rendering.States.ClearColor.Set(0.75, 0.25, 0.1, 1)
	cam.Rendering.States.FaceCulling = false
	cam.Rendering.Viewport.SetRel(0.02, 0.04, 0.125, 0.222)
	cam.SetScene("gui2d")

	if meshBuf, err = ng.Core.MeshBuffers.Add("buf_quad", ng.Core.MeshBuffers.NewParams(6, 6)); err != nil {
		return
	}
	if quadMesh, err = ng.Core.Libs.Meshes.AddLoad("mesh_quad", ng.MeshProviderPrefabQuad); err != nil {
		return
	}
	if err = meshBuf.Add(quadMesh); err != nil {
		return
	}

	quadMesh.Models.Default().SetMatID("mat_dog")
	me.Dog = scene.RootNode.ChildNodes.AddNew("gui_dog", "mesh_quad", "")
	me.Dog.Transform.SetScaleN(0.85)
	me.Dog.Transform.SetRotZ(unum.DegToRad(90))

	quadMesh.Models.Default().Clone("model_cat").SetMatID("mat_cat")
	me.Cat = scene.RootNode.ChildNodes.AddNew("gui_cat", "mesh_quad", "model_cat")
	me.Cat.Transform.SetScaleN(0.85)
	me.Cat.Transform.SetRotZ(unum.DegToRad(90))

	me.Dog.Transform.SetPosZ(0.1)
	me.Cat.Transform.SetPosZ(0.11)
	return
}
