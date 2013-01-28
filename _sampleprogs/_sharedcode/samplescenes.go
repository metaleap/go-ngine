package samplescenes

import (
	"fmt"
	"runtime"

	glfw "github.com/go-gl/glfw"
	ng "github.com/go3d/go-ngine/core"
	ugo "github.com/metaleap/go-util"
)

var (
	//	The primary camera
	Cam *ng.Camera

	//	Controller for the primary camera
	CamCtl *ng.Controller

	//	The maximum index for KeyHints when cycling through it in OnSec()
	MaxKeyHint = len(KeyHints) - 1

	//	OnSec() changes the window title every second to display FPS etc.
	//	Also every 4 seconds shows the next one in a number of "key hints" defined here:
	KeyHints = []string{
		"[F2]  --  Toggle Backface Culling",
		"[F3]  --  Pause/Resume",
		"[W][S]  --  Camera rise / sink",
		"[A][D]  --  Camera strafe left / right",
		"[<][>]  --  Camera turn left / right",
		"[^][v]  --  Camera move forward / backward",
		"[PgUp][PgDn]  --  Camera turn up / down",
		"[Alt][LShift][RShift]  --  Camera move-speed x 0.1 / 10 / 100",
	}

	curKeyHint = 0
	paused     = false
	sec        = 0
)

//	Creates a new core.Scene, adds it to the Scenes library under
//	the specified ID, and returns it.
func AddScene(id string) (me *ng.Scene) {
	me = ng.NewScene()
	ng.Core.Libs.Scenes[id] = me
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
					img.GpuSync()
				}
			}
		}
		img.OnLoad = func(image interface{}, err error, async bool) {
			if (err == nil) && (image != nil) && !async {
				img.GpuSync()
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

//	Called every frame (by the parent example app) to check the state for keys controlling CamCtl to move or rotate Cam.
func CheckCamCtlKeys() {
	if CamCtl.MoveSpeedupFactor = 1; paused {
		return
	}
	if ng.UserIO.KeyPressed(glfw.KeyLshift) {
		CamCtl.MoveSpeedupFactor = 10
	} else if ng.UserIO.KeyPressed(glfw.KeyRshift) {
		CamCtl.MoveSpeedupFactor = 100
	} else if ng.UserIO.KeyPressed(glfw.KeyLalt) {
		CamCtl.MoveSpeedupFactor = 0.1
	}
	if ng.UserIO.KeyPressed(glfw.KeyUp) {
		CamCtl.MoveForward()
	}
	if ng.UserIO.KeyPressed(glfw.KeyDown) {
		CamCtl.MoveBackward()
	}
	if ng.UserIO.KeyPressed('A') {
		CamCtl.MoveLeft()
	}
	if ng.UserIO.KeyPressed('D') {
		CamCtl.MoveRight()
	}
	if ng.UserIO.KeyPressed('W') {
		CamCtl.MoveUp()
	}
	if ng.UserIO.KeyPressed('S') {
		CamCtl.MoveDown()
	}
	if ng.UserIO.KeyPressed(glfw.KeyLeft) {
		CamCtl.TurnLeft()
	}
	if ng.UserIO.KeyPressed(glfw.KeyRight) {
		CamCtl.TurnRight()
	}
	if ng.UserIO.KeysPressedAny2(glfw.KeyPageup, glfw.KeyKP9) {
		CamCtl.TurnUp()
	}
	if ng.UserIO.KeysPressedAny2(glfw.KeyPagedown, glfw.KeyKP3) {
		CamCtl.TurnDown()
	}
}

//	Called every frame (by the parent example app) to check "toggle keys" that toggle certain options.
func CheckToggleKeys() {
	if ng.UserIO.KeyPressed(glfw.KeyEsc) {
		ng.Loop.Stop()
	}
	if ng.UserIO.KeyToggled(glfw.KeyF2) {
		Cam.Rendering.States.FaceCulling = !Cam.Rendering.States.FaceCulling
	}
	if ng.UserIO.KeyToggled(glfw.KeyF3) {
		PauseResume()
	}
}

func PauseResume() {
	canv := ng.Core.Rendering.Canvases.Main()
	if paused = ng.Core.Rendering.PostFx.ToggleEffect("Grayscale"); paused {
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
		fmt.Printf("%v:\tAvg=%3.5f secs\tMax=%3.5f secs\n", name, timing.Average(), timing.Max())
	}
	fmt.Printf("Average FPS:\t\t%v (total %v over %6.2fsec.)\n", ng.Stats.AverageFps(), ng.Stats.TotalFrames(), ng.Loop.Time())
	printStatSummary("Frame Full Loop", &ng.Stats.Frame)
	printStatSummary("Frame Render (CPU)", &ng.Stats.FrameRenderCpu)
	printStatSummary("Frame Render (GPU)", &ng.Stats.FrameRenderGpu)
	printStatSummary("Frame Render Both", &ng.Stats.FrameRenderBoth)
	printStatSummary("Frame Core Code", &ng.Stats.FrameCoreCode)
	printStatSummary("Frame User Code", &ng.Stats.FrameUserCode)
	printStatSummary("GC (max 1x/sec)", &ng.Stats.Gc)
	fmt.Printf("CGO calls: %v, Goroutines: %v\n\n", runtime.NumCgoCall(), runtime.NumGoroutine())

}

//	The *func main()* implementation for the parent example app. Initializes go:ngine, sets Cam and CamCtl, calls the specified assetLoader function, then enters The Loop.
func SamplesMainFunc(assetLoader func()) {
	runtime.LockOSThread()

	opt := ng.NewEngineOptions(AssetRootDirPath(), 1280, 720, 0, false)
	opt.Diagnostics.Log.MeshUploadInfo = true

	if err := ng.Init(opt, "Loading Sample App..."); err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		defer ng.Dispose()
		ng.Loop.OnSec = OnSec
		Cam = ng.Core.Rendering.Canvases[0].Cameras[0]
		Cam.Rendering.States.ClearColor.Set(0.5, 0.6, 0.85, 1)
		CamCtl = &Cam.Controller
		assetLoader()
		ng.Core.SyncUpdates()
		ng.Loop.Loop()
		PrintPostLoopSummary()
	}
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

//	Returns the window title to be set by OnSec().
func WindowTitle() string {
	return ng.Sfmt("%v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", ng.Stats.FpsLastSec, ng.UserIO.WinWidth(), ng.UserIO.WinHeight(), KeyHints[curKeyHint], CamCtl.Pos, CamCtl.Dir)
}
