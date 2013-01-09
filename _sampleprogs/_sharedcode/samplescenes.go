package samplescenes

import (
	"fmt"
	"runtime"

	glfw "github.com/go-gl/glfw"

	util "github.com/metaleap/go-util"

	ng "github.com/go3d/go-ngine/core"
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
		"[F2]  --  Toggle Render Technique",
		"[F3]  --  Toggle Backface Culling",
		"[F4]  --  Toggle Texture Filtering",
		"[F5]  --  Increase Filtering Anisotropy",
		"[W][S]  --  Camera rise / sink",
		"[A][D]  --  Camera strafe left / right",
		"[<][>]  --  Camera turn left / right",
		"[^][v]  --  Camera move forward / backward",
		"[PgUp][PgDn]  --  Camera turn up / down",
		"[Alt][LShift][RShift]  --  Camera move-speed x 0.1 / 10 / 100",
	}

	curKeyHint = 0
	sec        = 0
)

//	Creates a new core.Scene, adds it to the Scenes library under
//	the specified ID, and returns it.
func AddScene(id string) (me *ng.Scene) {
	me = ng.NewScene()
	ng.Core.Libs.Scenes[""] = me
	return
}

//	Sets up texture materials (associated and 2D samplers) with the specified IDs and image URLs.
//	For each ID (such as "foo" and "bar"):
//	-	creates an assets.FxImageDef with ID "tex_ID" (ie. "tex_foo" and "tex_bar")
//	-	creates an assets.FxEffectDef with ID "fx_ID" (ie. "fx_foo" and "fx_bar")
//	-	creates an assets.FxMaterialDef with ID "mat_ID" (ie. "mat_foo" and "mat_bar")
//	-	and links them all together to create a ready-to-use texture material.
func AddTextureMaterials(idsUrls map[string]string) {
	var (
	// image *ng.Texture
	// effect *nga.FxEffectDef
	)
	for id, refUrl := range idsUrls {
		ng.Core.Libs.Textures.AddNew("tex_"+id, refUrl)
		// image = nga.FxImageDefs.Add(ngau.NewFxImageDef("tex_"+id, refUrl))

		// effect = nga.FxEffectDefs.Add(ngau.NewFxEffectDef("fx_"+id, true, false))
		// effect.Profiles[0].NewParams.Set("sampler_"+id, ngau.NewFxSampler2D(image.DefaultInst(), nil, nil))
		// effect.Profiles[0].Common.Technique.Diffuse = ngau.NewFxColorOrTexture(ngau.NewFxTexture("sampler_"+id, "TEX0"), nil, "")

		// nga.FxMaterialDefs.AddNew("mat_" + id).Effect.DefRef.SetIdRef("fx_" + id)

		ng.Core.Libs.Materials["mat_"+id] = ng.Core.Libs.Materials.New("tex_" + id)
	}
}

//	Returns the "asset root directory" path for go:ngine, in this case: $GOPATH/src/github.com/go3d/go-ngine/_sampleprogs/_sharedassets
func AssetRootDirPath() string {
	return util.GopathSrcGithub("go3d", "go-ngine", "_sampleprogs", "_sharedassets")
}

//	Called every frame (by the parent example app) to check the state for keys controlling CamCtl to move or rotate Cam.
func CheckCamCtlKeys() {
	CamCtl.MoveSpeedupFactor = 1
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
		Cam.ToggleTechnique()
	}
	if ng.UserIO.KeyToggled(glfw.KeyF3) {
		ng.Core.Rendering.States.ToggleFaceCulling()
	}
	if ng.UserIO.KeyToggled(glfw.KeyF4) {
		ng.Core.Options.DefaultTextureParams.ToggleFilter()
	}
	if ng.UserIO.KeyToggled(glfw.KeyF5) {
		ng.Core.Options.DefaultTextureParams.ToggleFilterAnisotropy()
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
	fmt.Printf("CGO calls: %v, Goroutines: %v", runtime.NumCgoCall(), runtime.NumGoroutine())
}

//	The *func main()* implementation for the parent example app. Initializes go:ngine, sets Cam and CamCtl, calls the specified assetLoader function, then enters the Loop.
func SamplesMainFunc(assetLoader func()) {
	runtime.LockOSThread()

	if err := ng.Init(ng.NewEngineOptions(AssetRootDirPath(), 1280, 720, 0, false), "Loading Sample..."); err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		defer ng.Dispose()
		ng.Core.Rendering.States.SetClearColor(0.75, 0.75, 0.97, 1)
		ng.Loop.OnSec = OnSec
		Cam = ng.Core.Libs.Cameras.AddNew("")
		CamCtl = &Cam.Controller
		assetLoader()
		for _, canvas := range ng.Core.Rendering.Canvases {
			canvas.SetCameraIDs("")
		}
		ng.Core.SyncUpdates()
		ng.Loop.Loop()
		PrintPostLoopSummary()
	}
}

//	Called every second by go:ngine's *core* package. Refreshes the window title and doing so, every 4 seconds shows the next one entry in KeyHints.
func OnSec() {
	if sec++; sec == 4 {
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
