package core

import (
)

type tMaterials map[string]*TMaterial

	func (me tMaterials) New (texKey string) *TMaterial {
		var mat = &TMaterial {}
		mat.texKey = texKey
		return mat
	}

	func (me tMaterials) Set (name, texName string) *TMaterial {
		var mat = me.New(texName)
		me[name] = mat
		return mat
	}

type TMaterial struct {
	texKey string
}
