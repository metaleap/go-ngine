package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type tOptions struct {
	AssetRootDirPath, DefaultRenderTechnique string
	DefaultTextureParams *tTextureParams
	GpuMemMeshes, GpuMemTextures uint

	glBackfaceCulling, glTextureAnisotropy, winFullScreen bool
	winHeight, winSwapInterval, winWidth int
}

func NewOptions (assetRootDirPath string, winWidth, winHeight, winSwapInterval int, winFullScreen bool) *tOptions {
	var opt = &tOptions {}
	opt.DefaultTextureParams = newTextureParams(true, 6)
	opt.AssetRootDirPath, opt.DefaultRenderTechnique = assetRootDirPath, "rt_unlit_colored"
	opt.winWidth, opt.winHeight, opt.winSwapInterval, opt.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
	opt.GpuMemMeshes, opt.GpuMemTextures = 128, 256
	return opt
}

func (me *tOptions) GlBackfaceCulling () bool {
	return me.glBackfaceCulling
}

func (me *tOptions) SetGlBackfaceCulling (val bool) {
	if me.glBackfaceCulling = val; me.glBackfaceCulling {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (me *tOptions) ToggleGlBackfaceCulling () {
	me.SetGlBackfaceCulling(!me.glBackfaceCulling)
}
