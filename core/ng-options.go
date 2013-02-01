package core

import (
	"runtime"

	ugl "github.com/go3d/go-opengl/util"
)

//	Consider EngineOptions a "Singleton" type, only valid use is the one instance you created for core.Init().
//	Various "global" (rather than use-case-specific) options.
type EngineOptions struct {
	//	The base directory path for asset file paths.
	AssetRootDirPath string

	Initialization struct {
		GlCoreContext bool
	}

	Misc struct {
		DefaultControllerParams *ControllerParams
	}

	Rendering struct {
		DefaultClearColor ugl.GlVec4

		//	Name for the default render technique of a Camera2D,
		//	defaults to the currently only implementation "rt_unlit3".
		DefaultTechnique2D string

		//	Name for the default render technique of a Camera3D,
		//	defaults to the currently only implementation "rt_unlit3".
		DefaultTechnique3D string

		PostFx struct {
			TextureRect bool
		}
	}

	glTextureAnisotropy, winFullScreen   bool
	winHeight, winSwapInterval, winWidth int
}

//	Allocates, initializes and returns a new core.EngineOptions instance.
func NewEngineOptions(assetRootDirPath string, winWidth, winHeight, winSwapInterval int, winFullScreen bool) (me *EngineOptions) {
	me = &EngineOptions{AssetRootDirPath: assetRootDirPath}
	me.Misc.DefaultControllerParams = NewControllerParams()
	me.Initialization.GlCoreContext = (runtime.GOOS == "darwin")
	me.Rendering.DefaultClearColor = ugl.GlVec4{0, 0, 0, 1}
	me.Rendering.DefaultTechnique2D, me.Rendering.DefaultTechnique3D = "rt_unlit3", "rt_unlit3"
	me.winWidth, me.winHeight, me.winSwapInterval, me.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
	return
}
