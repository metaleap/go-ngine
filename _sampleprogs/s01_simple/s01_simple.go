package main

import (
	"fmt"
	"runtime"

	ngine "github.com/metaleap/go-ngine/client"
	ngine_samplescenes "github.com/metaleap/go-ngine/samplescenes"
)

func ngineOnSec () {
	ngine.Windowing.SetTitle(ngine_samplescenes.WindowTitle())
}

func main () {
	runtime.LockOSThread()
	var err error
	defer ngine.Dispose()

	if err = ngine.Init(1920, 1080, false, 0, ngine_samplescenes.AssetRootDirPath(), "Loading Sample...", ngineOnSec); err != nil {
		fmt.Printf("ABORT: %v\n", err)
	} else {
		ngine.Stats.TrackGC = true
		ngine_samplescenes.LoadSampleScene_01_TriQuad()
		ngine.Loop.Loop()
		ngine_samplescenes.PrintPostLoopSummary()
	}
}
