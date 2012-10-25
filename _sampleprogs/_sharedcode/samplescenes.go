package samplescenes

import (
	"fmt"
	"runtime"

	glfw "github.com/go-gl/glfw"

	util "github.com/go3d/go-util"

	ngine "github.com/go3d/go-ngine/core"
)

var (
	Cam *ngine.TCamera
	CamCtl *ngine.TController
	MaxKeyHint = len(keyHints) - 1

	curKeyHint = 0
	keyHints = []string {
		"[F2]  --  Toggle Render Technique",
		"[F3]  --  Toggle Backface Culling",
		"[F4]  --  Toggle Texture Filtering",
		"[F5]  --  Increase Filtering Anisotropy",
		"[W][S]  --  Camera rise / sink",
		"[A][D]  --  Camera strafe left / right",
		"[<][>]  --  Camera turn left / right",
		"[^][v]  --  Camera move forward / backward",
		"[PgUp][PgDown]  --  Camera turn up / down",
		"[Alt][LShift][RShift]  --  Camera move-speed x 0.1 / 10 / 100",
	}
	sec = 0
)

func AssetRootDirPath () string {
	return util.BaseCodePath("go-ngine", "_sampleprogs", "_sharedassets")
}

func CheckCamCtlKeys () {
	CamCtl.MoveSpeedupFactor = 1
	if ngine.Windowing.KeyPressed(glfw.KeyLshift) {
		CamCtl.MoveSpeedupFactor = 10
	} else if ngine.Windowing.KeyPressed(glfw.KeyRshift) {
		CamCtl.MoveSpeedupFactor = 100
	} else if ngine.Windowing.KeyPressed(glfw.KeyLalt) {
		CamCtl.MoveSpeedupFactor = 0.1
	}
	if ngine.Windowing.KeyPressed(glfw.KeyUp) { CamCtl.MoveForward() }
	if ngine.Windowing.KeyPressed(glfw.KeyDown) { CamCtl.MoveBackward() }
	if ngine.Windowing.KeyPressed('A') { CamCtl.MoveLeft() }
	if ngine.Windowing.KeyPressed('D') { CamCtl.MoveRight() }
	if ngine.Windowing.KeyPressed('W') { CamCtl.MoveUp() }
	if ngine.Windowing.KeyPressed('S') { CamCtl.MoveDown() }
	if ngine.Windowing.KeyPressed(glfw.KeyLeft) { CamCtl.TurnLeft() }
	if ngine.Windowing.KeyPressed(glfw.KeyRight) { CamCtl.TurnRight() }
	if ngine.Windowing.KeysPressedAny2(glfw.KeyPageup, glfw.KeyKP9) { CamCtl.TurnUp() }
	if ngine.Windowing.KeysPressedAny2(glfw.KeyPagedown, glfw.KeyKP3) { CamCtl.TurnDown() }
}

func CheckToggleKeys () {
	if ngine.Windowing.KeyToggled(glfw.KeyF2) { Cam.ToggleTechnique() }
	if ngine.Windowing.KeyToggled(glfw.KeyF3) { ngine.Core.Options.ToggleGlBackfaceCulling() }
	if ngine.Windowing.KeyToggled(glfw.KeyF4) { ngine.Core.Options.DefaultTextureParams.ToggleFilter() }
	if ngine.Windowing.KeyToggled(glfw.KeyF5) { ngine.Core.Options.DefaultTextureParams.ToggleFilterAnisotropy() }
}

func NewMaterialFromLocalTextureImageFile (assetRootRelativeFilePath string) *ngine.TMaterial {
	ngine.Core.Textures[assetRootRelativeFilePath] = ngine.Core.Textures.NewLoad(false, ngine.Core.Textures.Providers().LocalFile, assetRootRelativeFilePath)
	return ngine.Core.Materials.New(assetRootRelativeFilePath)
}

func NewMaterialFromRemoteTextureImageFile (fileUrl string) *ngine.TMaterial {
	ngine.Core.Textures[fileUrl] = ngine.Core.Textures.NewLoad(true, ngine.Core.Textures.Providers().RemoteFile, fileUrl)
	return ngine.Core.Materials.New(fileUrl)
}

func PrintPostLoopSummary () {
	fmt.Printf("Avg. FPS: %v\n", ngine.Stats.FpsOverallAverage())
	if ngine.Stats.TrackGC {
		fmt.Printf("GC: avg=%v max=%v\n", ngine.Stats.GcOverallAverage(), ngine.Stats.GcMaxTime)
	}
}

func SamplesMainFunc (loader func ()) {
	runtime.LockOSThread()
	var err error
	defer ngine.Dispose()

	if err = ngine.Init(ngine.NewOptions(AssetRootDirPath(), 1280, 720, 0, false), "Loading Sample..."); err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		ngine.Loop.OnSecTick = SamplesOnSec
		ngine.Stats.TrackGC = true
		Cam = ngine.Core.Canvases[ngine.Core.DefaultCanvasIndex].Cameras[0]
		CamCtl = Cam.Controller
		loader()
		ngine.Loop.Loop()
		PrintPostLoopSummary()
	}
}

func SamplesOnSec () {
	if sec++; sec == 4 {
		sec = 0
		if curKeyHint++; (curKeyHint > MaxKeyHint) || (curKeyHint >= (len(keyHints))) { curKeyHint = 0 }
	}
	ngine.Windowing.SetTitle(WindowTitle())
}

func WindowTitle () string {
	return fmt.Sprintf("%v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", ngine.Stats.FpsLastSec, ngine.Windowing.WinWidth(), ngine.Windowing.WinHeight(), keyHints[curKeyHint], CamCtl.Pos, CamCtl.Dir)
}
