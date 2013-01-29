package core

//	A hash-table of Models associated with their ID.
//	Used only for Mesh.Models.
type Models map[string]*Model

//	Returns the default Model (with ID "") for the parent Mesh.
func (me Models) Default() *Model {
	return me[""]
}

//	A Model is a parameterized instantiation of its parent Mesh geometry
//	with unique appearance, material or other properties.
//	
//	Each Mesh provides at least one Model, the "default model" (with ID ""),
//	accessible via someMesh.Models.Default(). To create new models for a Mesh,
//	call someMesh.Models["sourceModelID"].Clone("newModelID").
type Model struct {
	matID string
	id    string
	mat   *FxMaterial
	mesh  *Mesh
}

func newModel(id string, mesh *Mesh) (me *Model) {
	me = &Model{id: id, mesh: mesh}
	return
}

//	Creates a copy of me and adds it to the parent Mesh's Models
//	hash-table under the specified newModelID.
func (me *Model) Clone(newModelID string) (clonedModel *Model) {
	if (newModelID != me.id) && (me.mesh.Models[newModelID] == nil) {
		clonedModel = newModel(newModelID, me.mesh)
		me.mesh.Models[newModelID] = clonedModel
	}
	return
}

func (me *Model) MatID() string {
	return me.matID
}

func (me *Model) SetMatID(newMatID string) {
	if newMatID != me.matID {
		me.mat, me.matID = Core.Libs.Materials[newMatID], newMatID
	}
}
