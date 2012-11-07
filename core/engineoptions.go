package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

type engineOptions struct {
	AssetRootDirPath, DefaultRenderTechnique string
	DefaultTextureParams *textureParams
	GpuMemMeshes, GpuMemTextures uint

	glClearColor ugl.GlVec4
	glTextureAnisotropy, winFullScreen bool
	winHeight, winSwapInterval, winWidth int
}

func NewOptions (assetRootDirPath string, winWidth, winHeight, winSwapInterval int, winFullScreen bool) *engineOptions {
	var opt = &engineOptions {}
	opt.glClearColor = ugl.GlVec4 { 0, 0, 0, 1 }
	opt.DefaultTextureParams = newTextureParams(true, 6)
	opt.AssetRootDirPath, opt.DefaultRenderTechnique = assetRootDirPath, "rt_unlit_colored"
	opt.winWidth, opt.winHeight, opt.winSwapInterval, opt.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
	opt.GpuMemMeshes, opt.GpuMemTextures = 128, 256
	return opt
}

func (me *engineOptions) GlClearColor () ugl.GlVec4 {
	return me.glClearColor
}

func (me *engineOptions) SetGlClearColor (col ugl.GlVec4) {
	me.glClearColor = col
	gl.ClearColor(col[0], col[1], col[2], col[3])
}
