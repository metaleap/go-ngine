package exampleutils

import (
	"compress/flate"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	ng "github.com/go3d/go-ngine/core"
	ugo "github.com/metaleap/go-util"
	unum "github.com/metaleap/go-util/num"
)

var (
	//	Change to higher value to check out your splash-screen
	ArtificialSplashScreenDelay = 0 * time.Second

	//	Optionally set this to a callback function to be
	//	invoked every second on the windowing main thread.
	OnSec = func() {}

	//	The RenderCanvas the example scene is initially being rendered to.
	//	This is an off-screen "render-to-texture" RenderCanvas.
	SceneCanvas *ng.RenderCanvas

	//	The primary scene-rendering camera, rendering to SceneCanvas.
	SceneCam *ng.Camera

	//	Do not set this field directly, only use PauseResume() to
	//	toggle it and effect the associated render-state changes.

	//	Unlike the off-screen (render-to-texture) SceneCanvas above,
	//	this RenderCanvas epresents the actual screen/window.
	PostFxCanvas *ng.RenderCanvas

	//	Takes the image rendered to SceneCanvas, may
	//	post-process it or not, and blits it to PostFxCanvas.
	PostFxCam *ng.Camera

	Paused bool

	retro  bool
	numCgo struct {
		preLoop  int64
		postLoop int64
	}
	winTitle struct {
		appName        string
		cw, ch         int
		camPos, camDir unum.Vec3
	}

	curKeyHint = 0
	sec        = 0
)

//	Returns the base path of the "app dir" for our example apps, in this case: $GOPATH/src/github.com/go3d/go-ngine/_examples/-app-data
func AppDirBasePath() string {
	return ugo.GopathSrcGithub("go3d", "go-ngine", "_examples", "-app-data")
}

//	Returns the window title to be set by onSec().
func appWindowTitle() string {
	winTitle.cw, winTitle.ch = ng.UserIO.Window.Width(), ng.UserIO.Window.Height()
	if SceneCanvas != nil {
		winTitle.cw, winTitle.ch = SceneCanvas.CurrentAbsoluteSize()
	}
	if SceneCam != nil {
		winTitle.camPos, winTitle.camDir = SceneCam.Controller.Pos, *SceneCam.Controller.Dir()
	}
	return fmt.Sprintf("%s   |   %v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", winTitle.appName, ng.Stats.FpsLastSec, winTitle.cw, winTitle.ch, KeyHints[curKeyHint], winTitle.camPos.String(), winTitle.camDir.String())
}

func onSec() {
	if sec++; sec == 3 {
		sec = 0
		if curKeyHint++; (curKeyHint > MaxKeyHint) || (curKeyHint >= (len(KeyHints))) {
			curKeyHint = 0
		}
	}
	ng.UserIO.Window.SetTitle(appWindowTitle())
	OnSec()
}

//	Called by each example-app's func main(). Initializes go:ngine, sets SceneCam/SceneCanvas/PostFxCam/PostFxCanvas etc., calls the specified setupExampleScene function, then enters The Loop.
func Main(setupExampleScene, onAppThread, onWinThread func()) {
	//	by design, go:ngine doesn't do this for you automagically:
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())

	//	can set window options before it is created
	win := &ng.UserIO.Window
	//	release apps shouldn't do this, but during dev/test we want to observe max fps:
	win.SetSwapInterval(0)
	winFullscreen := false
	win.SetSize(1280, 720)
	winTitle.appName = filepath.Base(os.Args[0])
	win.SetTitle(fmt.Sprintf("Loading \"%s\" example app... (%v CPU cores)", winTitle.appName, runtime.GOMAXPROCS(0)))

	opt := &ng.Options

	//	While the default for this (true on Macs only) is reasonable for release apps at present,
	//	here we force core profile to verify all of go:ngine's GL code is fully core-profile compliant
	opt.Initialization.GlContext.CoreProfile.ForceFirst = true

	//	Release apps shouldn't do this, but here we're verifying everything runs in the oldest-supported GL version:
	opt.Initialization.GlContext.CoreProfile.VersionHint = 3.3

	//	Worth toggling this every once in a while just to see whether it makes a perf diff at all...
	realThreads := true
	opt.Loop.ForceThreads.App, opt.Loop.ForceThreads.Prep = realThreads, realThreads
	opt.Loop.GcEvery.Frame = true

	opt.AppDir.BasePath = AppDirBasePath()
	opt.AppDir.Temp.BaseName = filepath.Join("_tmp", filepath.Base(os.Args[0]))
	opt.AppDir.Temp.ShaderSources, opt.AppDir.Temp.CachedTextures = "glsl", "tex"
	// but for now, we don't need separate per-app tmp dirs:
	opt.AppDir.Temp.BaseName = "_tmp"

	if compressCachedTextures := true; compressCachedTextures {
		opt.Textures.Storage.DiskCache.Compressor = func(w io.WriteCloser) (wc io.WriteCloser) {
			var err error
			if wc, err = flate.NewWriter(w, 9); err != nil {
				panic(err)
			}
			return
		}
		opt.Textures.Storage.DiskCache.Decompressor = func(r io.ReadCloser) io.ReadCloser {
			return flate.NewReader(r)
		}
	}

	//	STEP 1: init go:ngine
	err := ng.Init(winFullscreen)
	if err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		// defer ng.Dispose()

		//	STEP 2: post-init, pre-loop setup
		ng.Loop.On.EverySec, ng.Loop.On.AppThread, ng.Loop.On.WinThread = onSec, onAppThread, onWinThread

		PostFxCanvas = &ng.Core.Rendering.Canvases[0]
		PostFxCam = &PostFxCanvas.Cameras[0]

		if setupExampleScene != nil {
			SceneCanvas = ng.Core.Rendering.Canvases.AddNew()
			SceneCam = SceneCanvas.AddNewCamera3D()
			SceneCam.Rendering.States.ClearColor.Set(0.5, 0.6, 0.85, 1)
			setupExampleScene()
			if err = ng.Core.Libs.Meshes.GpuSync(); err != nil {
				panic(err)
			}
			ng.Core.GpuSyncImageLibs()
		}
		time.Sleep(ArtificialSplashScreenDelay)
		numCgo.preLoop = runtime.NumCgoCall()

		//	STEP 3: enter... Da Loop.
		ng.Loop.Run()
		numCgo.postLoop = runtime.NumCgoCall()
		PrintPostLoopSummary() // don't wanna defer this: useless when exit-on-error
	}
}
