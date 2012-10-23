package core

import (
)

type TMaterial struct {
	texKey string
}

func NewMaterial (texKey string) *TMaterial {
	var mat = &TMaterial {}
	mat.texKey = texKey
	return mat
}
