package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/metaleap/go-util/gl"
)

type TMaterial struct {
	glTexture gl.Uint
}

func NewMaterialFromLocalTextureImageFile (filePath string) *TMaterial {
	var mat = &TMaterial {}
	mat.glTexture = glutil.MakeTextureFromImageFile(filePath, gl.REPEAT, gl.REPEAT, gl.LINEAR_MIPMAP_LINEAR, gl.LINEAR, true)
	return mat
}

func (me *TMaterial) Dispose () {
	gl.DeleteTextures(1, &me.glTexture)
}
