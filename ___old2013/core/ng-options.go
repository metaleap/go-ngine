package core

import (
	"io"
	"runtime"

	u3d "github.com/metaleap/go-util-3d"
	ugl "github.com/metaleap/go-opengl/util"

	ngctx "github.com/metaleap/go-ngine/glctx"
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

var (
	Options NgOptions
)

//	Only used for the Options variable.
type NgOptions struct {
	AppDir struct {
		//	The base directory path for app file paths.
		BasePath string

		Temp struct {
			BaseName       string
			ShaderSources  string
			CachedTextures string
		}
	}

	Cameras struct {
		DefaultControllerParams ControllerParams
		PerspectiveDefaults     u3d.Perspective
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
				//	Only used if go:ngine requests the creation of a GL core-profile context.
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
		DefaultCanvas struct {
			GammaViaShader bool
			SplashImage    []byte
		}
		Window struct {
			//	Defaults: Color.R=8 Color.G=8 Color.B=8 Color.A=0 Depth=8 Stencil=0.
			//	These defaults are reasonable when using a render-to-texture off-screen
			//	RenderCanvas. Otherwise, may want to bump Depth to at least 24 or 32.
			//	Depth shouldn't be 0 as this causes some Intel HD drivers to bug out badly.
			BufSizes ngctx.BufferBits
		}
	}

	Libs struct {
		InitialCap   int
		GrowCapBy    int
		OnIDsChanged struct {
			Effects LibElemIDChangedHandlers
			Images  struct {
				TexCube LibElemIDChangedHandlers
				Tex2D   LibElemIDChangedHandlers
			}
			Materials LibElemIDChangedHandlers
			Meshes    LibElemIDChangedHandlers
			Models    LibElemIDChangedHandlers
			Scenes    LibElemIDChangedHandlers
		}
	}

	Loop struct {
		//	Controls whether and how often the Garbage Collector
		//	is invoked during the Loop.
		GcEvery struct {
			//	Defaults to false. If true, GC will
			//	be invoked every frame during the Loop.
			Frame bool

			//	Defaults to true. If true, GC will be invoked at
			//	least and at most once per second during the Loop.
			Sec bool
		}
	}

	Rendering struct {
		DefaultBatcher    RenderBatcher
		DefaultClearColor ugl.GlVec4
	}

	Textures struct {
		Storage FxImageStorage
	}
}

func init() {
	o := &Options
	o.Textures.Storage.Gpu.UintRev, o.Textures.Storage.Gpu.Bgra = true, true
	o.Textures.Storage.DiskCache.Enabled = true
	o.Textures.Storage.DiskCache.Compressor = func(w io.WriteCloser) io.WriteCloser { return w }
	o.Textures.Storage.DiskCache.Decompressor = func(r io.ReadCloser) io.ReadCloser { return r }
	o.Cameras.DefaultControllerParams.initDefaults()
	persp := &o.Cameras.PerspectiveDefaults
	persp.FovY.Deg, persp.ZFar, persp.ZNear, persp.Enabled = 37.8493, 30000, 0.3, true
	o.Loop.GcEvery.Sec = true
	o.Libs.InitialCap, o.Libs.GrowCapBy = 16, 32

	//	Set all ID-changed handlers to empty funcs so we don't need to check for nil
	init, isMac, initGl := &o.Initialization, runtime.GOOS == "darwin", &o.Initialization.GlContext
	initGl.CoreProfile.ForceFirst, initGl.CoreProfile.ForwardCompat, initGl.BadVersionMessage = isMac, isMac, DefaultBadVersionMessage
	initGl.CoreProfile.VersionHint = ugl.KnownVersions[len(ugl.KnownVersions)-1]
	buf := &init.Window.BufSizes
	buf.Color.R, buf.Color.G, buf.Color.B, buf.Depth = 8, 8, 8, 8

	rend := &o.Rendering
	rend.DefaultBatcher.Enabled = true
	rend.DefaultBatcher.Priority[0] = BatchByProgram
	rend.DefaultBatcher.Priority[1] = BatchByTexture
	rend.DefaultBatcher.Priority[2] = BatchByBuffer
	rend.DefaultClearColor = ugl.GlVec4{0, 0, 0, 1}

	win := &UserIO.Window
	win.OnCloseRequested = func() bool { return true }
	win.title, win.width, win.height, win.swap, win.ResizeMinDelay = "go:ngine", 1024, 576, 1, 0.15
}
