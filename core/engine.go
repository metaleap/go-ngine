package core

import (
	"log"

	glutil "github.com/go3d/go-util/gl"

	nglcore "github.com/go3d/go-ngine/core/glcore"
)

var (
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
	Windowing.dispose()
	Core, Loop, Stats = nil, nil, nil
}

func Init (options *tOptions, winTitle string) error {
	var err error
	if err = Windowing.init(options, winTitle); err == nil {
		if err = nglcore.Init(); err == nil {
			Loop, Stats, Core = newEngineLoop(), newEngineStats(), newEngineCore(options)
			log.Println(glutil.GlConnInfo())
		}
	}
	return err
}
