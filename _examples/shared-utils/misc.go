package exampleutils

import (
	"fmt"
	"runtime"
	"time"

	ng "github.com/go3d/go-ngine/core"
)

//	Pauses rendering or resumes from the current pause.
//	When paused, the frame last rendered is frozen and rendered in a gray-scale effect.
func PauseResume() {
	tech := PostFxCam.Rendering.Technique.(*ng.RenderTechniqueQuad)
	if Paused = tech.ToggleEffect("Grayscale"); Paused {
		if SceneCanvas != nil {
			SceneCanvas.EveryNthFrame = 0
		}
	} else {
		if SceneCanvas != nil {
			SceneCanvas.EveryNthFrame = 1
		}
	}
	if err := tech.ApplyEffects(); err != nil {
		ng.Diag.LogErr(err)
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
	fmt.Printf("Shaders: compiled %v GLSL programs over time, which took %v in total.\n", ng.Stats.Programs.NumProgsCompiled, time.Duration(ng.Stats.Programs.TotalTimeCost))
	fmt.Printf("CGO calls: pre-loop init %v, loop %v (avg. %v/frame)\n\n", numCgo.preLoop, numCgo.postLoop-numCgo.preLoop, runtime.NumCgoCall()/int64(ng.Stats.TotalFrames()))
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
