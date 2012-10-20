package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type tOptions struct {
	DefaultRenderTechnique string

	glBackfaceCulling, glTextureAnisotropy, glTextureFiltering, glTextureMipMaps bool
}

func NewOptions (glBackfaceCulling bool) *tOptions {
	var opt = &tOptions {}
	opt.DefaultRenderTechnique = "unlit"
	opt.glBackfaceCulling = glBackfaceCulling
	return opt
}

func (me *tOptions) GlBackfaceCulling () bool {
	return me.glBackfaceCulling
}

func (me *tOptions) GlTextureAnisotropy () bool {
	return me.glTextureAnisotropy
}

func (me *tOptions) GlTextureFiltering () bool {
	return me.glTextureFiltering
}

func (me *tOptions) GlTextureMipMaps () bool {
	return me.glTextureMipMaps
}

func (me *tOptions) SetGlBackfaceCulling (val bool) {
	if me.glBackfaceCulling = val; me.glBackfaceCulling {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (me *tOptions) SetGlTextureAnisotropy (val bool) {
	if val != me.glTextureAnisotropy {
		
	}
}

func (me *tOptions) ToggleGlBackfaceCulling () {
	me.SetGlBackfaceCulling(!me.glBackfaceCulling)
}
