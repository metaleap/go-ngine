package core

import (
	nga "github.com/go3d/go-ngine/assets"
)

type models map[string]*Model

	func (me models) Default () *Model {
		return me[""]
	}

type Model struct {
	mat *nga.Material
	matName string
	mesh *Mesh
	name string
}

func newModel (name string, mesh *Mesh) *Model {
	var model = &Model {}
	model.name, model.mesh = name, mesh
	return model
}

func (me *Model) Clone (modelName string) (clonedModel *Model) {
	if (modelName != me.name) && (me.mesh.Models[modelName] == nil) {
		clonedModel = newModel(modelName, me.mesh)
		me.mesh.Models[modelName] = clonedModel
	}
	return
}

func (me *Model) MatName () string {
	return me.matName
}

func (me *Model) render () {
	curTechnique.onRenderMeshModel()
	me.mesh.render()
}

func (me *Model) SetMatName (newMatName string) {
	if newMatName != me.matName { me.mat, me.matName = nga.Materials[newMatName], newMatName }
}
