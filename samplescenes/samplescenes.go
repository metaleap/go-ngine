package samplescenes

import (
	"fmt"
	"os"
	"path/filepath"

	ngine "github.com/metaleap/go-ngine/client"
	ncore "github.com/metaleap/go-ngine/client/core"
)

var (
	ctl *ncore.TController
)

func AssetRootDirPath () string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "metaleap", "go-ngine", "_sampleprogs", "_sharedassets")
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

func WindowTitle () string {
	ctl = ngine.Core.CurCamera.Controller
	return fmt.Sprintf("%v FPS @ %vx%v | Cam: P=%v D=%v", ngine.Stats.Fps, ngine.Windowing.WinWidth, ngine.Windowing.WinHeight, ctl.Pos, ctl.Dir)
}
