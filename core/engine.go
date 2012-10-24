package core

import (
	"log"

	glutil "github.com/go3d/go-util/gl"
)

var (
	Loop *tEngineLoop
	Core *tEngineCore
	Stats *tEngineStats
	Windowing = newWindowing()
)

func Dispose () {
	glLogLastError("ngine.Pre_Dispose")
	if Core != nil { Core.Dispose() }
	glLogLastError("ngine.Core_Dispose")
	glDispose()
	Windowing.dispose()
	Core, Loop, Stats = nil, nil, nil
}

func Init (options *tOptions, winTitle string) error {
	var err error
	if err = Windowing.init(options, winTitle); err == nil {
		if err = glInit(); err == nil {
			Loop, Stats, Core = newEngineLoop(), newEngineStats(), newEngineCore(options)
			log.Println(glutil.GlConnInfo())
		}
	}
	return err
}
