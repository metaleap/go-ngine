package assets

type OldMaterials map[string]*Material

	func (me OldMaterials) New (texName string) *Material {
		var mat = &Material {}
		mat.TexName = texName
		return mat
	}

	func (me OldMaterials) Set (name, texName string) *Material {
		var mat = me.New(texName)
		me[name] = mat
		return mat
	}

type Material struct {
	TexName string
}
