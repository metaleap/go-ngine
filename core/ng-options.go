package core

import (
	"runtime"

	ugl "github.com/go3d/go-glutil"
)

//	Consider EngineOptions a "Singleton" type, only valid use is the one instance you created for core.Init().
//	Various "global" (rather than use-case-specific) options.
type EngineOptions struct {
	//	The base directory path for asset file paths.
	AssetRootDirPath string

	Initialization struct {
		GlCoreContext bool
	}

	Rendering struct {
		DefaultClearColor ugl.GlVec4

		//	Name for the default render technique, defaults to the only
		//	currently supported value "rt_unlit".
		DefaultTechnique string
	}

	glTextureAnisotropy, winFullScreen   bool
	winHeight, winSwapInterval, winWidth int
}

//	Allocates, initializes and returns a new core.EngineOptions instance.
func NewEngineOptions(assetRootDirPath string, winWidth, winHeight, winSwapInterval int, winFullScreen bool) (me *EngineOptions) {
	me = &EngineOptions{AssetRootDirPath: assetRootDirPath}
	me.Initialization.GlCoreContext = (runtime.GOOS == "darwin")
	me.Rendering.DefaultClearColor = ugl.GlVec4{0, 0, 0, 1}
	me.Rendering.DefaultTechnique = "rt_unlit"
	me.winWidth, me.winHeight, me.winSwapInterval, me.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
	return
}
