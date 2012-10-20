package samplescenes

import (
	"fmt"
	"runtime"

	util "github.com/go3d/go-util"

	ngine "github.com/go3d/go-ngine/core"
)

var (
	MaxKeyHint = len(keyHints) - 1

	ctl *ngine.TController
	curKeyHint = 0
	keyHints = []string {
		"[F2]  --  Toggle Render Technique",
		"[F3]  --  Toggle Backface Culling",
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

func NewMeshCube () *ngine.TMesh {
	var mesh = &ngine.TMesh {}
	mesh.SetVertsCube()
	return mesh
}

func NewMeshFace3 () *ngine.TMesh {
	var mesh = &ngine.TMesh {}
	mesh.SetVertsFace3()
	return mesh
}

func NewMeshFace4 () *ngine.TMesh {
	var mesh = &ngine.TMesh {}
	mesh.SetVertsFace4()
	return mesh
}

func NewMeshPlane () *ngine.TMesh {
	var mesh = &ngine.TMesh {}
	mesh.SetVertsPlane()
	return mesh
}

func NewMeshPyramid () *ngine.TMesh {
	var mesh = &ngine.TMesh {}
	mesh.SetVertsPyramid()
	return mesh
}

func PrintPostLoopSummary () {
	fmt.Printf("Avg. FPS: %v\n", ngine.Stats.FpsOverallAverage())
	if ngine.Stats.TrackGC {
		fmt.Printf("Max. GC: %v\n", ngine.Stats.GcMaxTime)
	}
}

func SamplesMainFunc (loader func ()) {
	runtime.LockOSThread()
	var err error
	defer ngine.Dispose()

	if err = ngine.Init(1280, 720, false, 0, AssetRootDirPath(), "Loading Sample...", SamplesOnSec); err != nil {
		fmt.Printf("ABORT: %v\n", err)
	} else {
		ngine.Stats.TrackGC = true
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
	ctl = ngine.Core.CurCamera.Controller
	return fmt.Sprintf("%v FPS @ %vx%v   |   %s   |   Cam: P=%v D=%v", ngine.Stats.FpsLastSec, ngine.Windowing.WinWidth, ngine.Windowing.WinHeight, keyHints[curKeyHint], ctl.Pos, ctl.Dir)
}
