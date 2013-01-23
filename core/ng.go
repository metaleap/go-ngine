package core

import (
	"fmt"
	"runtime"

	ugl "github.com/go3d/go-glutil"
)

var (
	//	Manages your main-thread's "game loop". You'll need to call it's Loop() method once after go:ngine initialization (see samples).
	Loop EngineLoop

	//	The heart and brain of go:ngine --- a container for all runtime resources and responsible for rendering.
	Core EngineCore

	//	Tracks various go:ngine performance counters over time.
	Stats EngineStats

	//	Your gateway for end-user input (key, mouse etc.) and "output" (window management, not the graphics themselves).
	UserIO EngineUserIO
)

//	Call this to "un-init" go:ngine and to release any and all GPU or RAM resources still allocated.
func Dispose() {
	Core.dispose()
	glDispose()
	UserIO.dispose()
}

//	Initializes go:ngine; this first attempts to initialize OpenGL and then open a window to your supplied specifications with a GL 3.2-or-higher profile.
func Init(options *EngineOptions, winTitle string) (err error) {
	isVerErr, forceContext := false, true || (runtime.GOOS == "darwin")
tryInit:
	if err = UserIO.init(options, winTitle, forceContext); err == nil {
		if err, isVerErr = glInit(); err == nil {
			Stats.reset()
			Loop.init()
			Core.init(options)
			ugl.LogLastError("INIT")
		} else if isVerErr && !forceContext {
			forceContext = true
			UserIO.isGlfwInit, UserIO.isGlfwWindow = false, false
			goto tryInit
		}
	}
	return
}

//	A short-hand for fmt.Sprintf. Feel free to ignore.
func Sfmt(format string, fmtArgs ...interface{}) string {
	return fmt.Sprintf(format, fmtArgs...)
}
