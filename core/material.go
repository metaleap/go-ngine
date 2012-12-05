package core

type materials map[string]*Material

func (me materials) New(texName string) (mat *Material) {
	mat = &Material{TexName: texName}
	return
}

func (me materials) Set(name, texName string) (mat *Material) {
	mat = me.New(texName)
	me[name] = mat
	return
}

type Material struct {
	TexName string
}
