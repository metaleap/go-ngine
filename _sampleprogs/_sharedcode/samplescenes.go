package samplescenes

import (
	"fmt"
	"runtime"

	glfw "github.com/go-gl/glfw"

	ugl "github.com/go3d/go-glutil"
	util "github.com/metaleap/go-util"

	ng "github.com/go3d/go-ngine/core"
	nga "github.com/go3d/go-ngine/assets"
)

var (
	Cam *ng.Camera
	CamCtl *ng.Controller
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
	return util.BaseCodePathGithub("go3d", "go-ngine", "_sampleprogs", "_sharedassets")
}

func CheckCamCtlKeys () {
	CamCtl.MoveSpeedupFactor = 1
	if ng.UserIO.KeyPressed(glfw.KeyLshift) {
		CamCtl.MoveSpeedupFactor = 10
	} else if ng.UserIO.KeyPressed(glfw.KeyRshift) {
		CamCtl.MoveSpeedupFactor = 100
	} else if ng.UserIO.KeyPressed(glfw.KeyLalt) {
		CamCtl.MoveSpeedupFactor = 0.1
	}
	if ng.UserIO.KeyPressed(glfw.KeyUp) { CamCtl.MoveForward() }
	if ng.UserIO.KeyPressed(glfw.KeyDown) { CamCtl.MoveBackward() }
	if ng.UserIO.KeyPressed('A') { CamCtl.MoveLeft() }
	if ng.UserIO.KeyPressed('D') { CamCtl.MoveRight() }
	if ng.UserIO.KeyPressed('W') { CamCtl.MoveUp() }
	if ng.UserIO.KeyPressed('S') { CamCtl.MoveDown() }
	if ng.UserIO.KeyPressed(glfw.KeyLeft) { CamCtl.TurnLeft() }
	if ng.UserIO.KeyPressed(glfw.KeyRight) { CamCtl.TurnRight() }
	if ng.UserIO.KeysPressedAny2(glfw.KeyPageup, glfw.KeyKP9) { CamCtl.TurnUp() }
	if ng.UserIO.KeysPressedAny2(glfw.KeyPagedown, glfw.KeyKP3) { CamCtl.TurnDown() }
}

func CheckToggleKeys () {
	if ng.UserIO.KeyPressed(glfw.KeyEsc) { ng.Loop.Stop() }
	if ng.UserIO.KeyToggled(glfw.KeyF2) { Cam.ToggleTechnique() }
	if ng.UserIO.KeyToggled(glfw.KeyF3) { Cam.Options.ToggleGlBackfaceCulling() }
	if ng.UserIO.KeyToggled(glfw.KeyF4) { ng.Core.Options.DefaultTextureParams.ToggleFilter() }
	if ng.UserIO.KeyToggled(glfw.KeyF5) { ng.Core.Options.DefaultTextureParams.ToggleFilterAnisotropy() }
}

func PrintPostLoopSummary () {
	var printStatSummary = func (name string, timing *ng.TimingStats) {
		fmt.Printf("%v:\tAvg=%3.5f secs\tMax=%3.5f secs\n", name, timing.Average(), timing.Max())
	}
	fmt.Printf("Average FPS:\t\t%v\n", ng.Stats.AverageFps())
	printStatSummary("Frame Full Loop", ng.Stats.Frame)
	printStatSummary("Frame Render (CPU)", ng.Stats.FrameRenderCpu)
	printStatSummary("Frame Render (GPU)", ng.Stats.FrameRenderGpu)
	printStatSummary("Frame Render Both", ng.Stats.FrameRenderBoth)
	printStatSummary("Frame Core Code", ng.Stats.FrameCoreCode)
	printStatSummary("Frame User Code", ng.Stats.FrameUserCode)
	printStatSummary("GC (max 1x/sec)", ng.Stats.Gc)
}

func SamplesMainFunc (loader func ()) {
	runtime.LockOSThread()
	var err error
	defer ng.Dispose()

	if err = ng.Init(ng.NewOptions(AssetRootDirPath(), 1280, 720, 0, false), "Loading Sample..."); err != nil {
		fmt.Printf("ABORT:\n%v\n", err)
	} else {
		ng.Core.Options.SetGlClearColor(ugl.GlVec4 { 0.75, 0.75, 0.97, 1 })
		ng.Loop.OnSecTick = SamplesOnSec
		camDef := nga.CameraDefs.AddNew("")
		camDef.FovOrMagY, camDef.ZnearPlane, camDef.ZfarPlane = 37.8493, 0.3, 30000
		ng.Core.SyncAssetDefs()
		Cam = ng.Core.Cameras[""]
		CamCtl = Cam.Controller
		loader()
		ng.Loop.Loop()
		PrintPostLoopSummary()
	}
}

func SamplesOnSec () {
	if sec++; sec == 4 {
		sec = 0
		if curKeyHint++; (curKeyHint > MaxKeyHint) || (curKeyHint >= (len(keyHints))) { curKeyHint = 0 }
	}
	ng.UserIO.SetWinTitle(WindowTitle())
}

func WindowTitle () string {
	return fmt.Sprintf("%v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", ng.Stats.FpsLastSec, ng.UserIO.WinWidth(), ng.UserIO.WinHeight(), keyHints[curKeyHint], CamCtl.Pos, CamCtl.Dir)
}
