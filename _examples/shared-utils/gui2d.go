package exampleutils

import (
	ng "github.com/go3d/go-ngine/core"
	unum "github.com/metaleap/go-util/num"
)

//	A fake "2D GUI" concoction. There will be much better ngine-provided support for stuff like this.
//	Has two textured quads, a cat and a dog one, shows them both animated and overlapping
//	inside a red 64x64 px square in the bottom left canvas corner.
type Gui2D struct {
	View     *ng.RenderView
	Cat, Dog *ng.Node
}

func (me *Gui2D) Setup() (err error) {
	var (
		meshBuf    *ng.MeshBuffer
		quadMeshID int
	)
	sceneID := ng.Core.Libs.Scenes.AddNew()
	scene := &ng.Core.Libs.Scenes[sceneID]
	me.View = SceneCanvas.AddNewView("Scene")
	cam := &me.View.Technique_Scene().Camera
	cam.Perspective.Enabled = false
	me.View.RenderStates.ClearColor.Set(0.75, 0.25, 0.1, 1)
	me.View.Port.SetAbsolute(8, 8, 64, 64) //SetRel(0.02, 0.04, 0.125, 0.222)
	cam.SetScene(sceneID)
	if meshBuf, err = ng.Core.MeshBuffers.AddNew("buf_quad", 6); err != nil {
		return
	}
	if quadMeshID, err = ng.Core.Libs.Meshes.AddNewAndLoad("mesh_quad", ng.MeshDescriptorQuad); err != nil {
		return
	}
	if err = meshBuf.Add(quadMeshID); err != nil {
		return
	}

	me.Dog = scene.RootNode.ChildNodes.AddNew("gui_dog", quadMeshID)
	me.Dog.MatID = LibIDs.Mat["dog"]
	me.Dog.Transform.SetScale(0.85)
	me.Dog.Transform.Rot.Z = unum.DegToRad(90)

	me.Cat = scene.RootNode.ChildNodes.AddNew("gui_cat", quadMeshID)
	me.Cat.MatID = LibIDs.Mat["cat"]
	me.Cat.Transform.SetScale(0.85)
	me.Cat.Transform.Rot.Z = unum.DegToRad(90)

	me.Dog.Transform.Pos.Z = 0.1
	me.Cat.Transform.Pos.Z = 0.11
	me.Dog.ApplyTransform()
	me.Cat.ApplyTransform()
	return
}
