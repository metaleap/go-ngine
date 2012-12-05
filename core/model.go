package core

type models map[string]*Model

func (me models) Default() *Model {
	return me[""]
}

type Model struct {
	mat     *Material
	matName string
	mesh    *Mesh
	name    string
}

func newModel(name string, mesh *Mesh) (me *Model) {
	me = &Model{name: name, mesh: mesh}
	return
}

func (me *Model) Clone(modelName string) (clonedModel *Model) {
	if (modelName != me.name) && (me.mesh.Models[modelName] == nil) {
		clonedModel = newModel(modelName, me.mesh)
		me.mesh.Models[modelName] = clonedModel
	}
	return
}

func (me *Model) MatName() string {
	return me.matName
}

func (me *Model) render() {
	curTechnique.onRenderMeshModel()
	me.mesh.render()
}

func (me *Model) SetMatName(newMatName string) {
	if newMatName != me.matName {
		me.mat, me.matName = Core.Materials[newMatName], newMatName
	}
}
