package core

type SceneNode struct {
	ID        int
	Transform SceneNodeTransform

	Render struct {
		Enabled bool
		MatID   int
		MeshID  int
		ModelID int

		skyMode bool
	}

	parentID     int
	childNodeIDs []int
}

func (me *SceneNode) dispose() {
}

func (me *SceneNode) init() {
	me.Render.Enabled = true
	me.Render.MatID, me.Render.MeshID, me.Render.ModelID = -1, -1, -1
	me.Transform.init()
}

func (me *SceneNode) meshMat() (mesh *Mesh, mat *FxMaterial) {
	if mesh = Core.Libs.Meshes.get(me.Render.MeshID); mesh != nil {
		if mat = Core.Libs.Materials.get(me.Render.MatID); mat == nil {
			model := Core.Libs.Models.get(me.Render.ModelID)
			if model == nil {
				model = Core.Libs.Models.get(mesh.DefaultModelID)
			}
			if model != nil {
				mat = Core.Libs.Materials.get(model.MatID)
			}
		}
	}
	return
}
