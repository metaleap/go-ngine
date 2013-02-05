package core

import (
	"runtime"

	ugl "github.com/go3d/go-opengl/util"
)

const DefaultBadVersionMessage = `
Minimum required OpenGL version is {MINVER}, but your
current {OS} graphics driver only provides
OpenGL version {CURVER}.

Most likely your machine is just
missing some recent system updates.

*HOW TO RESOLVE*:

1. On the web, search for "downloading & installing the
latest {OS} driver for {GPU}",

2. or simply visit the {VENDOR} website and locate
their "driver downloads" pages to obtain the most
recent driver for {GPU}.
`

//	Only valid initialization for an EngineOptions is NewEngineOptions().
//	This object is passed to the core.Init() function which immediately copies it to Core.Options.
//	All subsequent option reads and writes are exclusively from and to Core.Options, the copy you initially
//	received from NewEngineOptions() is no longer used. Slightly confusing, but there's a reason for this design.
type EngineOptions struct {
	//	The base directory path for asset file paths.
	AssetRootDirPath string

	Initialization struct {
		GlContext struct {
			//	Required on Mac OS X, not necessary elsewhere.
			//	While potentially slightly beneficial with recent GL drivers, might also fail with somewhat outdated ones.
			//	Defaults to true on Mac OS X, to false elsewhere.
			CoreProfile bool

			//	Required on Mac OS X, not recommended elsewhere.
			//	Defaults to true on Mac OS X, to false elsewhere.
			ForwardCompat bool
		}
		BadVersionMessage string
		Window            struct {
			Rbits, Gbits, Bbits, Abits, DepthBits, StencilBits int
		}
	}

	Misc struct {
		DefaultControllerParams *ControllerParams
	}

	Rendering struct {
		DefaultClearColor ugl.GlVec4

		//	Name for the default render technique of a Camera2D,
		//	defaults to the currently only implementation "rt_unlit".
		DefaultTechnique2D string

		//	Name for the default render technique of a Camera3D,
		//	defaults to the currently only implementation "rt_unlit".
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
	init, isMac := &me.Initialization, runtime.GOOS == "darwin"
	init.GlContext.CoreProfile, init.GlContext.ForwardCompat, init.BadVersionMessage = isMac, isMac, DefaultBadVersionMessage
	init.Window.Rbits, init.Window.Gbits, init.Window.Bbits = 8, 8, 8
	// this depth-bits should be 0 really: since there's no depth involved in the final postfx-pass -- but then Intel cards bug out badly
	init.Window.DepthBits = 8
	rend := &me.Rendering
	rend.DefaultClearColor = ugl.GlVec4{0, 0, 0, 1}
	rend.DefaultTechnique2D, rend.DefaultTechnique3D = "rt_unlit", "rt_unlit"
	me.winWidth, me.winHeight, me.winSwapInterval, me.winFullScreen = winWidth, winHeight, winSwapInterval, winFullScreen
	return
}
