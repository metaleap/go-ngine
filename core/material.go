package core

type Materials map[string]*Material

func (me Materials) New(texName string) (mat *Material) {
	mat = &Material{TexName: texName}
	return
}

func (me Materials) Set(name, texName string) (mat *Material) {
	mat = me.New(texName)
	me[name] = mat
	return
}

type Material struct {
	TexName string
}
