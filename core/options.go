package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

type tOptions struct {
	AssetRootDirPath, DefaultRenderTechnique string
	DefaultTextureParams *tTextureParams
	GpuMemMeshes, GpuMemTextures uint

	glClearColor glutil.TGlVec4
	glTextureAnisotropy, winFullScreen bool
	winHeight, winSwapInterval, winWidth int
}

func NewOptions (assetRootDirPath string, winWidth, winHeight, winSwapInterval int, winFullScreen bool) *tOptions {
	var opt = &tOptions {}
	opt.glClearColor = glutil.TGlVec4 { 0, 0, 0, 1 }
	opt.DefaultTextureParams = newTextureParams(true, 6)
	opt.AssetRootDirPath, opt.DefaultRenderTechnique = assetRootDirPath, "rt_unlit_colored"
	opt.winWidth, opt.winHeight, opt.winSwapInterval, opt.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
	opt.GpuMemMeshes, opt.GpuMemTextures = 128, 256
	return opt
}

func (me *tOptions) GlClearColor () glutil.TGlVec4 {
	return me.glClearColor
}

func (me *tOptions) SetGlClearColor (col glutil.TGlVec4) {
	me.glClearColor = col
	gl.ClearColor(col[0], col[1], col[2], col[3])
}
