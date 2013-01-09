package core

//	Consider EngineOptions a "Singleton" type, only valid use is the one instance you created for core.Init().
//	Various "global" (rather than use-case-specific) options.
type EngineOptions struct {
	//	The base directory path for asset file paths.
	AssetRootDirPath string

	//	Name for the default render technique, currently "rt_unlit_colored"
	DefaultRenderTechnique string

	//	All textures default to these params unless they have their own params set.
	DefaultTextureParams *textureParams

	glTextureAnisotropy, winFullScreen   bool
	winHeight, winSwapInterval, winWidth int
}

//	Allocates, initializes and returns a new core.EngineOptions instance.
func NewEngineOptions(assetRootDirPath string, winWidth, winHeight, winSwapInterval int, winFullScreen bool) (me *EngineOptions) {
	me = &EngineOptions{AssetRootDirPath: assetRootDirPath, DefaultRenderTechnique: "rt_unlit_colored"}
	me.DefaultTextureParams = newTextureParams(true, 8)
	me.winWidth, me.winHeight, me.winSwapInterval, me.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
	return
}
