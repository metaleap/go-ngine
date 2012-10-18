package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type TOptions struct {
	glBackfaceCulling, glTextureAnisotropy, glTextureFiltering, glTextureMipMaps bool
}

func NewOptions (glBackfaceCulling bool) *TOptions {
	var opt = &TOptions {}
	opt.glBackfaceCulling = glBackfaceCulling
	return opt
}

func (me *TOptions) GlBackfaceCulling () bool {
	return me.glBackfaceCulling
}

func (me *TOptions) GlTextureAnisotropy () bool {
	return me.glTextureAnisotropy
}

func (me *TOptions) GlTextureFiltering () bool {
	return me.glTextureFiltering
}

func (me *TOptions) GlTextureMipMaps () bool {
	return me.glTextureMipMaps
}

func (me *TOptions) SetGlBackfaceCulling (val bool) {
	if me.glBackfaceCulling = val; me.glBackfaceCulling {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (me *TOptions) SetGlTextureAnisotropy (val bool) {
	if val != me.glTextureAnisotropy {
		
	}
}

func (me *TOptions) ToggleGlBackfaceCulling () {
	me.SetGlBackfaceCulling(!me.glBackfaceCulling)
}
