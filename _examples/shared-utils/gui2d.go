package exampleutils

import (
	ng "github.com/go3d/go-ngine/core"
	unum "github.com/metaleap/go-util/num"
)

//	A fake "2D GUI" concoction. There will be much better ngine-provided support for stuff like this.
//	Has two textured quads, a cat and a dog one, shows them both animated and overlapping
//	inside a red 64x64 px square in the bottom left canvas corner.
type Gui2D struct {
	View                 *ng.RenderView
	Cam                  *ng.Camera
	CatNodeID, DogNodeID int
}

func (me *Gui2D) Setup() (err error) {
	var (
		meshBuf    *ng.MeshBuffer
		quadMeshID int
	)
	sceneID := ng.Core.Libs.Scenes.AddNew()
	scene := &ng.Core.Libs.Scenes[sceneID]
	me.View = SceneCanvas.AddNewView("Scene")
	rts := me.View.Technique_Scene()
	rts.Batch.Enabled, me.Cam = false, &rts.Camera
	me.Cam.Perspective.Enabled = false
	me.View.RenderStates.ClearColor.Set(0.75, 0.25, 0.1, 1)
	me.View.Port.SetAbsolute(8, 8, 64, 64) //SetRel(0.02, 0.04, 0.125, 0.222)
	me.Cam.SetScene(sceneID)
	if meshBuf, err = ng.Core.MeshBuffers.AddNew("buf_quad", 6); err != nil {
		return
	}
	if quadMeshID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_quad", ng.MeshDescriptorQuad); err != nil {
		return
	}
	if err = meshBuf.Add(quadMeshID); err != nil {
		return
	}

	me.DogNodeID = scene.AddNewChildNode(0)
	dog := scene.Node(me.DogNodeID)
	dog.Render.MeshID = quadMeshID
	dog.Render.MatID = LibIDs.Mat["dog"]
	dog.Transform.SetScale(0.85)
	dog.Transform.Rot.Z = unum.DegToRad(90)

	me.CatNodeID = scene.AddNewChildNode(0)
	cat := scene.Node(me.CatNodeID)
	cat.Render.MeshID = quadMeshID
	cat.Render.MatID = LibIDs.Mat["cat"]
	cat.Transform.SetScale(0.85)
	cat.Transform.Rot.Z = unum.DegToRad(90)

	dog.Transform.Pos.Z = 0.1
	cat.Transform.Pos.Z = 0.11
	scene.ApplyNodeTransforms(me.DogNodeID)
	scene.ApplyNodeTransforms(me.CatNodeID)
	return
}
