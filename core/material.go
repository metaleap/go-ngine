package core

import (
)

type materials map[string]*Material

	func (me materials) New (texName string) *Material {
		var mat = &Material {}
		mat.texName = texName
		return mat
	}

	func (me materials) Set (name, texName string) *Material {
		var mat = me.New(texName)
		me[name] = mat
		return mat
	}

type Material struct {
	texName string
}
