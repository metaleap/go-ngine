package core

import (
	"log"

	glutil "github.com/go3d/go-util/gl"

	nglcore "github.com/go3d/go-ngine/core/glcore"
)

var (
	AssetRootDirPath = "."
	Loop *tEngineLoop
	Core *tEngineCore
	Stats *tEngineStats
	Windowing = newWindowing()
)

func Dispose () {
	nglcore.LogLastError("ngine.Pre_Dispose")
	if Core != nil { Core.Dispose() }
	nglcore.LogLastError("ngine.Core_Dispose")
	nglcore.Dispose()
	Windowing.Exit()
	Core, Loop, Stats = nil, nil, nil
}

func Init (winWidth, winHeight int, winFullScreen bool, vsync int, assetRootDirPath, winTitle string, onSecTick func ()) error {
	var err error
	if err = Windowing.Init(winWidth, winHeight, winFullScreen, vsync, winTitle); err == nil {
		if err = nglcore.Init(); err == nil {
			AssetRootDirPath, Loop, Stats = assetRootDirPath, newEngineLoop(), &tEngineStats {}
			Core = newEngineCore(winWidth, winHeight)
			Loop.OnSecTick = onSecTick
			Loop.OnLoopHandlers = [] func () {}
			log.Println(glutil.GlConnInfo())
		}
	}
	return err
}
