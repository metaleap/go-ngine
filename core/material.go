package core

type materials map[string]*Material

func (me materials) New(texName string) *Material {
	var mat = &Material{}
	mat.TexName = texName
	return mat
}

func (me materials) Set(name, texName string) *Material {
	var mat = me.New(texName)
	me[name] = mat
	return mat
}

type Material struct {
	TexName string
}
