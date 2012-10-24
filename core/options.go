package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type tOptions struct {
	AssetRootDirPath, DefaultRenderTechnique string
	DefaultTextureParams *tTextureParams

	glBackfaceCulling, glTextureAnisotropy, glTextureFiltering, glTextureMipMaps, winFullScreen bool
	winHeight, winSwapInterval, winWidth int
}

func NewOptions (assetRootDirPath string, winWidth, winHeight, winSwapInterval int, winFullScreen bool) *tOptions {
	var opt = &tOptions {}
	opt.DefaultTextureParams = newTextureParams(true, 6)
	opt.AssetRootDirPath, opt.DefaultRenderTechnique = assetRootDirPath, "unlit"
	opt.winWidth, opt.winHeight, opt.winSwapInterval, opt.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
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

func (me *tOptions) WinHeight () int {
	return me.winHeight
}

func (me *tOptions) WinWidth () int {
	return me.winWidth
}
