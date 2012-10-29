package core

type tModels map[string]*TModel

	func (me tModels) Default () *TModel {
		return me[""]
	}

type TModel struct {
	mat *TMaterial
	matName string
	mesh *TMesh
	name string
}

func newModel (name string, mesh *TMesh) *TModel {
	var model = &TModel {}
	model.name, model.mesh = name, mesh
	return model
}

func (me *TModel) Clone (modelName string) (clonedModel *TModel) {
	if (modelName != me.name) && (me.mesh.Models[modelName] == nil) {
		clonedModel = newModel(modelName, me.mesh)
		me.mesh.Models[modelName] = clonedModel
	}
	return
}

func (me *TModel) MatName () string {
	return me.matName
}

func (me *TModel) render () {
	curTechnique.onRenderMeshModel()
	me.mesh.render()
}

func (me *TModel) SetMatName (newMatName string) {
	if newMatName != me.matName { me.mat, me.matName = Core.Materials[newMatName], newMatName }
}
