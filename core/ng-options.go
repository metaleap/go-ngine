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
	AppDir struct {
		//	The base directory path for app file paths.
		BasePath string
	}

	Initialization struct {
		GlContext struct {
			CoreProfile struct {
				//	Required on Mac OS X, not necessary elsewhere.
				//	While potentially slightly beneficial with recent GL drivers,
				//	might also fail with a select few rather outdated ones.
				//	Defaults to true on Mac OS X, to false elsewhere.
				ForceFirst bool

				//	While required on Mac OS X, really not recommended elsewhere
				//	(at present). Defaults to true on Mac OS X, to false elsewhere.
				ForwardCompat bool

				//	Defaults to the newest GL version currently supported by the GL
				//	binding used by go:ngine. The binding adaptively uses features
				//	of GL versions newer than 3.3 only if they are available, so this
				//	is a most strongly recommended default for release apps. But for
				//	testing, this is useful to test performance in older GL versions.
				//	Must be one of the values in glutil.KnownVersions.
				VersionHint float64
			}

			//	Defaults to the DefaultBadVersionMessage constant. If using a custom
			//	string, you can use the same placeholders as that one.
			BadVersionMessage string
		}
		Window struct {
			//	Defaults: R=8 G=8 B=8 A=0 D=8 S=0.
			//	These defaults are reasonable when using a render-to-texture off-screen
			//	RenderCanvas. Otherwise, may want to bump D to at least 24 or 32.
			//	D shouldn't be 0 as this causes some Intel HD drivers to bug out badly.
			Rbits, Gbits, Bbits, Abits, DepthBits, StencilBits int
		}
	}

	Loop struct {
		//	By default, the app and prep "threads" being invoked every frame
		//	during Loop() are just normal go-routines that may or may not
		//	correspond to real separate OS thread contexts.
		ForceThreads struct {
			//	If true, each On.AppThread() go-routine invokation locks its own exclusive thread context.
			App bool

			//	If true, each prep-stage go-routine invokation locks its own exclusive thread context.
			Prep bool
		}

		//	By default, during Loop() GC is invoked at least and at most once per second.
		//	If GcEveryFrame is true, it will be invoked every frame.
		GcEveryFrame bool
	}

	Misc struct {
		DefaultControllerParams *ControllerParams
	}

	Rendering struct {
		DefaultClearColor ugl.GlVec4

		//	Default render technique for a Camera created via RenderCanvas.AddNewCamera2D().
		//	Defaults to "rt_scene".
		DefaultTechnique2D string

		//	Default render technique for a Camera created via RenderCanvas.AddNewCamera3D().
		//	Defaults to "rt_scene".
		DefaultTechnique3D string

		//	Default render technique for a Camera created via RenderCanvas.AddNewCameraQuad().
		//	Defaults to "rt_quad".
		DefaultTechniqueQuad string
	}
}

//	Allocates, initializes and returns a new core.EngineOptions instance.
func NewEngineOptions(appDirBasePath string, windowOptions *WindowOptions) (me *EngineOptions) {
	me = &EngineOptions{}
	me.AppDir.BasePath = appDirBasePath
	me.Misc.DefaultControllerParams = NewControllerParams()
	init, isMac, initGl := &me.Initialization, runtime.GOOS == "darwin", &me.Initialization.GlContext
	initGl.CoreProfile.ForceFirst, initGl.CoreProfile.ForwardCompat, initGl.BadVersionMessage = isMac, isMac, DefaultBadVersionMessage
	initGl.CoreProfile.VersionHint = ugl.KnownVersions[len(ugl.KnownVersions)-1]
	init.Window.Rbits, init.Window.Gbits, init.Window.Bbits, init.Window.DepthBits = 8, 8, 8, 8
	rend := &me.Rendering
	rend.DefaultClearColor = ugl.GlVec4{0, 0, 0, 1}
	rend.DefaultTechnique2D, rend.DefaultTechnique3D, rend.DefaultTechniqueQuad = "rt_scene", "rt_scene", "rt_quad"
	UserIO.Window = *windowOptions
	return
}
