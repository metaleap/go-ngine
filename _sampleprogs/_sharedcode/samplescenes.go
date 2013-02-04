package samplescenes

import (
	"fmt"
	"runtime"

	ng "github.com/go3d/go-ngine/core"
	ugo "github.com/metaleap/go-util"
)

var (
	//	The primary camera
	Cam *ng.Camera

	//	Controller for the primary camera
	CamCtl *ng.Controller

	Paused, retro bool

	curKeyHint = 0
	sec        = 0
)

//	Creates a new core.Scene, adds it to the Scenes library under
//	the specified ID, and returns it.
func AddScene(id string, mainCamScene bool) (me *ng.Scene) {
	me = ng.NewScene()
	if ng.Core.Libs.Scenes[id] = me; mainCamScene {
		Cam.SetScene(id)
	}
	return
}

//	Sets up plain-color effects/materials with the specified IDs.
//	For each ID (such as "foo" and "bar"):
//	-	creates an ng.FxEffect with ID "fx_{ID}" (ie. "fx_foo" and "fx_bar") and adds it
//	to ng.Core.Libs.Effects; its Diffuse field pointing to the color
//	-	creates an ng.FxMaterial with ID "mat_{ID}" (ie. "mat_foo" and "mat_bar") and
//	adds it to ng.Core.Libs.Materials; its DefaultEffectID pointing to the ng.FxEffect.
func AddColorMaterials(idsColors map[string][]float64) {
	for id, col := range idsColors {
		ng.Core.Libs.Effects.AddNew("fx_" + id).Diffuse = ng.NewFxColor(col...)
		ng.Core.Libs.Materials.AddNew("mat_" + id).DefaultEffectID = "fx_" + id
	}
}

//	Sets up textures and associated effects/materials with the specified IDs and image URLs.
//	For each ID (such as "foo" and "bar"):
//	-	creates an ng.FxImage2D with ID "img_{ID}" (ie. "img_foo" and "img_bar") and
//	adds it to ng.Core.Libs.Images.I2D
//	-	creates an ng.FxEffect with ID "fx_{ID}" (ie. "fx_foo" and "fx_bar") and adds it
//	to ng.Core.Libs.Effects; its Diffuse field pointing to the ng.FxImage2D
//	-	creates an ng.FxMaterial with ID "mat_{ID}" (ie. "mat_foo" and "mat_bar") and
//	adds it to ng.Core.Libs.Materials; its DefaultEffectID pointing to the ng.FxEffect.
func AddTextureMaterials(idsUrls map[string]string) {
	for id, refUrl := range idsUrls {
		img := ng.Core.Libs.Images.I2D.AddNew("img_" + id)
		if img.InitFrom.RefUrl = refUrl; img.IsRemote() {
			img.AsyncNumAttempts = -1
			img.OnAsyncDone = func() {
				if img.Loaded() {
					if err := img.GpuSync(); err != nil {
						panic(err)
					}
				}
			}
		}
		img.OnLoad = func(image interface{}, err error, async bool) {
			if (err == nil) && (image != nil) && !async {
				err = img.GpuSync()
			}
			if err != nil {
				panic(err)
			}
		}
		ng.Core.Libs.Effects.AddNew("fx_" + id).Diffuse = ng.NewFxTexture("img_"+id, nil)
		ng.Core.Libs.Materials.AddNew("mat_" + id).DefaultEffectID = "fx_" + id
	}
}

//	Returns the "asset root directory" path for go:ngine, in this case: $GOPATH/src/github.com/go3d/go-ngine/_sampleprogs/_sharedassets
func AssetRootDirPath() string {
	return ugo.GopathSrcGithub("go3d", "go-ngine", "_sampleprogs", "_sharedassets")
}

//	Called every second by go:ngine's *core* package. Refreshes the window title and doing so, every 3 seconds shows the next one entry in KeyHints.
func OnSec() {
	if sec++; sec == 3 {
		sec = 0
		if curKeyHint++; (curKeyHint > MaxKeyHint) || (curKeyHint >= (len(KeyHints))) {
			curKeyHint = 0
		}
	}
	ng.UserIO.SetWinTitle(WindowTitle())
}

func PauseResume() {
	canv := ng.Core.Rendering.Canvases.Main()
	if Paused = ng.Core.Rendering.PostFx.ToggleEffect("Grayscale"); Paused {
		canv.EveryNthFrame = 0
	} else {
		canv.EveryNthFrame = 1
	}
	if err := ng.Core.Rendering.PostFx.ApplyEffects(); err != nil {
		ugo.LogError(err)
	}
}

//	Prints a summary of go:ngine's *Stats* performance counters when the parent example app exits.
func PrintPostLoopSummary() {
	printStatSummary := func(name string, timing *ng.TimingStats) {
		fmt.Printf("%v:\t\tAvg=%3.5f secs\tMax=%3.5f secs\n", name, timing.Average(), timing.Max())
	}
	fmt.Printf("Average FPS:\t\t%v (total %v over %6.2fsec.)\n", ng.Stats.AverageFps(), ng.Stats.TotalFrames(), ng.Loop.Time())
	printStatSummary("Frame Full Loop", &ng.Stats.Frame)
	printStatSummary("Frame OnAppThread", &ng.Stats.FrameAppThread)
	printStatSummary("Frame OnWinThread", &ng.Stats.FrameWinThread)
	printStatSummary("Frame Prep Thread", &ng.Stats.FramePrepThread)
	printStatSummary("Frame Thread Sync", &ng.Stats.FrameThreadSync)
	printStatSummary("Frame Render (CPU)", &ng.Stats.FrameRenderCpu)
	printStatSummary("Frame Render (GPU)", &ng.Stats.FrameRenderGpu)
	printStatSummary("Frame Render Both", &ng.Stats.FrameRenderBoth)
	printStatSummary("GC (max 1x/sec)", &ng.Stats.Gc)
	fmt.Printf("CGO calls: %v\n\n", runtime.NumCgoCall())
}

//	The *func main()* implementation for the parent example app. Initializes go:ngine, sets Cam and CamCtl, calls the specified assetLoader function, then enters The Loop.
func SamplesMainFunc(assetLoader func()) {
	runtime.LockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())

	width, height, fullscreen := 1280, 720, false
	// width, height, fullscreen := 1920, 1080, true
	opt := ng.NewEngineOptions(AssetRootDirPath(), width, height, 0, fullscreen)
	opt.Initialization.GlCoreContext = true
	// opt.Rendering.PostFx.TextureRect = true

	if err := ng.Init(opt, fmt.Sprintf("Loading sample app... (%v CPU cores)", runtime.GOMAXPROCS(0))); err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		defer ng.Dispose()
		ng.Loop.OnSec = OnSec
		Cam = ng.Core.Rendering.Canvases[0].Cameras[0]
		Cam.Rendering.States.ClearColor.Set(0.5, 0.6, 0.85, 1)
		CamCtl = &Cam.Controller
		assetLoader()
		ng.Core.SyncUpdates()
		// ng.Loop.SwapLast = true // because our OnAppThread() isn't doing any heavy stuff
		ng.Loop.Loop()
		PrintPostLoopSummary()
	}
}

func ToggleRetro() {
	if retro = !retro; retro {
		ng.Core.Rendering.Canvases.Main().SetSize(true, 0.25, 0.25)
	} else {
		ng.Core.Rendering.Canvases.Main().SetSize(true, 1, 1)
	}
}

//	Returns the window title to be set by OnSec().
func WindowTitle() string {
	cw, ch := ng.Core.Rendering.Canvases.Main().CurrentAbsoluteSize()
	return fmt.Sprintf("%v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", ng.Stats.FpsLastSec, cw, ch, KeyHints[curKeyHint], CamCtl.Pos.String(), CamCtl.Dir().String())
}
