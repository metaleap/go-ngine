package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

//	Consider EngineOptions a "Singleton" type, only valid use is the one instance you created for core.Init().
//	Various "global" (rather than use-case-specific) options.
type EngineOptions struct {
	//	The base directory path for asset file paths.
	AssetRootDirPath string

	//	Name for the default render technique, currently "rt_unlit_colored"
	DefaultRenderTechnique string

	//	All textures default to these params unless they have their own params set.
	DefaultTextureParams *textureParams

	glClearColor                         ugl.GlVec4
	glTextureAnisotropy, winFullScreen   bool
	winHeight, winSwapInterval, winWidth int
}

//	Allocates, initializes and returns a new core.EngineOptions instance.
func NewEngineOptions(assetRootDirPath string, winWidth, winHeight, winSwapInterval int, winFullScreen bool) (me *EngineOptions) {
	me = &EngineOptions{AssetRootDirPath: assetRootDirPath, DefaultRenderTechnique: "rt_unlit_colored"}
	me.glClearColor = ugl.GlVec4{0, 0, 0, 1}
	me.DefaultTextureParams = newTextureParams(true, 8)
	me.winWidth, me.winHeight, me.winSwapInterval, me.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
	return
}

//	Returns the current OpenGL "clear color".
func (me *EngineOptions) GlClearColor() ugl.GlVec4 {
	return me.glClearColor
}

//	Changes the current OpenGL "clear color"
func (me *EngineOptions) SetGlClearColor(col ugl.GlVec4) {
	me.glClearColor = col
	gl.ClearColor(col[0], col[1], col[2], col[3])
}
