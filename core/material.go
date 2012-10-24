package core

import (
)

type tMaterials map[string]*TMaterial

	func (me *tMaterials) New (texKey string) *TMaterial {
		var mat = &TMaterial {}
		mat.texKey = texKey
		return mat
	}

type TMaterial struct {
	texKey string
}
