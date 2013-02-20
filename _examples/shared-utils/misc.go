package exampleutils

import (
	"fmt"
	"runtime"
	"time"

	ng "github.com/go3d/go-ngine/core"
)

var (
	SkyNode *ng.Node
)

//	Pauses rendering or resumes from the current pause.
//	When paused, the frame last rendered is frozen and rendered in a gray-scale effect.
func PauseResume() {
	tech := PostFxCam.RenderTechniqueQuad()
	tech.Effect.Ops.ToggleGrayscale(-1)
	tech.Effect.UpdateRoutine()
	if Paused = !Paused; Paused {
		if SceneCanvas != nil {
			SceneCanvas.EveryNthFrame = 0
		}
	} else {
		if SceneCanvas != nil {
			SceneCanvas.EveryNthFrame = 1
		}
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
	fmt.Printf("Shaders: compiled %v GLSL shader programs over time, which took %v in total.\n", ng.Stats.Programs.NumProgsCompiled, time.Duration(ng.Stats.Programs.TotalTimeCost))
	cgoPerFrame, numTotalFrames := int64(0), ng.Stats.TotalFrames()
	if numTotalFrames != 0 {
		cgoPerFrame = runtime.NumCgoCall() / int64(ng.Stats.TotalFrames())
	}
	fmt.Printf("CGO calls: pre-loop init %v, loop %v (avg. %v/frame)\n\n", numCgo.preLoop, numCgo.postLoop-numCgo.preLoop, cgoPerFrame)
}

func ToggleBatching() {
	if SceneCam != nil {
		SceneCam.RenderTechniqueScene().ToggleBatching()
	}
	if RearView.Cam != nil {
		RearView.Cam.RenderTechniqueScene().ToggleBatching()
	}
}

//	Toggles "retro mode" for the example app.
//	If retro is on, the resolution of the main canvas is 1/4th of the window resolution.
func ToggleRetro() {
	if SceneCanvas != nil && !Paused {
		if retro = !retro; retro {
			SceneCanvas.SetSize(true, 0.25, 0.25)
		} else {
			SceneCanvas.SetSize(true, 1, 1)
		}
	}
}

func ToggleTexturing() {
	for _, fx := range ng.Core.Libs.Effects {
		fx.Ops.ToggleColored(-1)
		fx.Ops.Toggle("Tex*", -1)
		fx.UpdateRoutine()
	}
}
