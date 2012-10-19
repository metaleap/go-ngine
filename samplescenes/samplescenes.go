package samplescenes

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	util "github.com/go3d/go-util"

	ngine "github.com/go3d/go-ngine/client"
	ncore "github.com/go3d/go-ngine/client/core"
)

var (
	ctl *ncore.TController
)

func AssetRootDirPath () string {
	return util.BaseCodePath("go-ngine", "_sampleprogs", "_sharedassets")
}

func NewMeshCube () *ncore.TMesh {
	var mesh = &ncore.TMesh {}
	mesh.SetVertsCube()
	return mesh
}

func NewMeshFace3 () *ncore.TMesh {
	var mesh = &ncore.TMesh {}
	mesh.SetVertsFace3()
	return mesh
}

func NewMeshFace4 () *ncore.TMesh {
	var mesh = &ncore.TMesh {}
	mesh.SetVertsFace4()
	return mesh
}

func NewMeshPlane () *ncore.TMesh {
	var mesh = &ncore.TMesh {}
	mesh.SetVertsPlane()
	return mesh
}

func NewMeshPyramid () *ncore.TMesh {
	var mesh = &ncore.TMesh {}
	mesh.SetVertsPyramid()
	return mesh
}

func PrintPostLoopSummary () {
	fmt.Printf("Avg. FPS: %v\n", ngine.Stats.FpsAvg())
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
	ngine.Windowing.SetTitle(WindowTitle())
}

func WindowTitle () string {
	ctl = ngine.Core.CurCamera.Controller
	return fmt.Sprintf("%v FPS @ %vx%v | F2: Toggle Render Technique | Cam: P=%v D=%v", ngine.Stats.Fps, ngine.Windowing.WinWidth, ngine.Windowing.WinHeight, ctl.Pos, ctl.Dir)
}
