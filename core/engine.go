package core

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
	var isVerErr bool
	var forceContext = false
	tryInit:
	if err = Windowing.init(options, winTitle, forceContext); err == nil {
		if err, isVerErr = glInit(); err == nil {
			Loop, Stats, Core = newEngineLoop(), newEngineStats(), newEngineCore(options)
		} else if isVerErr && !forceContext {
			forceContext = true
			Windowing.isGlfwInit, Windowing.isGlfwWindow = false, false
			goto tryInit
		}
	}
	return err
}
