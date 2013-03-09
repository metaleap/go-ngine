package core

type SceneNode struct {
	ID        int
	Transform SceneNodeTransform

	Render struct {
		Cull struct {
			Frustum bool
		}
		Enabled bool
		MatID   int
		ModelID int

		meshID  int
		skyMode bool
	}

	parentID     int
	childNodeIDs []int

	thrApp struct {
		bounding geometryBounds
	}
	thrPrep struct {
		bounding geometryBounds
	}
}

func (me *SceneNode) dispose() {
}

func (me *SceneNode) init() {
	me.Render.Enabled, me.Render.Cull.Frustum = true, true
	me.Render.MatID, me.Render.meshID, me.Render.ModelID = -1, -1, -1
	me.Transform.init()
}

func (me *SceneNode) mesh() *Mesh {
	return Core.Libs.Meshes.get(me.Render.meshID)
}

func (me *SceneNode) meshMat() (mesh *Mesh, mat *FxMaterial) {
	if mesh = me.mesh(); mesh != nil {
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
