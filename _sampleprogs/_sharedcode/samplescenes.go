package samplescenes

import (
	"fmt"
	"runtime"

	glfw "github.com/go-gl/glfw"

	ugl "github.com/go3d/go-glutil"
	util "github.com/metaleap/go-util"

	ngine "github.com/go3d/go-ngine/core"
)

var (
	Cam *ngine.Camera
	CamCtl *ngine.Controller
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
	return util.BaseCodePath("go3d", "go-ngine", "_sampleprogs", "_sharedassets")
}

func CheckCamCtlKeys () {
	CamCtl.MoveSpeedupFactor = 1
	if ngine.UserIO.KeyPressed(glfw.KeyLshift) {
		CamCtl.MoveSpeedupFactor = 10
	} else if ngine.UserIO.KeyPressed(glfw.KeyRshift) {
		CamCtl.MoveSpeedupFactor = 100
	} else if ngine.UserIO.KeyPressed(glfw.KeyLalt) {
		CamCtl.MoveSpeedupFactor = 0.1
	}
	if ngine.UserIO.KeyPressed(glfw.KeyUp) { CamCtl.MoveForward() }
	if ngine.UserIO.KeyPressed(glfw.KeyDown) { CamCtl.MoveBackward() }
	if ngine.UserIO.KeyPressed('A') { CamCtl.MoveLeft() }
	if ngine.UserIO.KeyPressed('D') { CamCtl.MoveRight() }
	if ngine.UserIO.KeyPressed('W') { CamCtl.MoveUp() }
	if ngine.UserIO.KeyPressed('S') { CamCtl.MoveDown() }
	if ngine.UserIO.KeyPressed(glfw.KeyLeft) { CamCtl.TurnLeft() }
	if ngine.UserIO.KeyPressed(glfw.KeyRight) { CamCtl.TurnRight() }
	if ngine.UserIO.KeysPressedAny2(glfw.KeyPageup, glfw.KeyKP9) { CamCtl.TurnUp() }
	if ngine.UserIO.KeysPressedAny2(glfw.KeyPagedown, glfw.KeyKP3) { CamCtl.TurnDown() }
}

func CheckToggleKeys () {
	if ngine.UserIO.KeyPressed(glfw.KeyEsc) { ngine.Loop.Stop() }
	if ngine.UserIO.KeyToggled(glfw.KeyF2) { Cam.ToggleTechnique() }
	if ngine.UserIO.KeyToggled(glfw.KeyF3) { Cam.Options.ToggleGlBackfaceCulling() }
	if ngine.UserIO.KeyToggled(glfw.KeyF4) { ngine.Core.Options.DefaultTextureParams.ToggleFilter() }
	if ngine.UserIO.KeyToggled(glfw.KeyF5) { ngine.Core.Options.DefaultTextureParams.ToggleFilterAnisotropy() }
}

func PrintPostLoopSummary () {
	var printStatSummary = func (name string, timing *ngine.TimingStats) {
		fmt.Printf("%v:\tAvg=%3.5f secs\tMax=%3.5f secs\n", name, timing.Average(), timing.Max())
	}
	fmt.Printf("Average FPS:\t\t%v\n", ngine.Stats.AverageFps())
	printStatSummary("Frame Full Loop", ngine.Stats.Frame)
	printStatSummary("Frame Render (CPU)", ngine.Stats.FrameRenderCpu)
	printStatSummary("Frame Render (GPU)", ngine.Stats.FrameRenderGpu)
	printStatSummary("Frame Render Both", ngine.Stats.FrameRenderBoth)
	printStatSummary("Frame Core Code", ngine.Stats.FrameCoreCode)
	printStatSummary("Frame User Code", ngine.Stats.FrameUserCode)
	printStatSummary("GC (max 1x/sec)", ngine.Stats.Gc)
}

func SamplesMainFunc (loader func ()) {
	runtime.LockOSThread()
	var err error
	defer ngine.Dispose()

	if err = ngine.Init(ngine.NewOptions(AssetRootDirPath(), 1280, 720, 0, false), "Loading Sample..."); err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		ngine.Core.Options.SetGlClearColor(ugl.GlVec4 { 0.75, 0.75, 0.97, 1 })
		ngine.Loop.OnSecTick = SamplesOnSec
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
	ngine.UserIO.SetWinTitle(WindowTitle())
}

func WindowTitle () string {
	return fmt.Sprintf("%v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", ngine.Stats.FpsLastSec, ngine.UserIO.WinWidth(), ngine.UserIO.WinHeight(), keyHints[curKeyHint], CamCtl.Pos, CamCtl.Dir)
}
