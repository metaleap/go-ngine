package core

import (
	"fmt"
	"log"
	"runtime"
)

var (
	//	Manages your main-thread's "game loop". You'll need to call it's Loop() method once after go:ngine initialization (see samples).
	Loop *EngineLoop

	//	The heart and brain of go:ngine --- a container for all runtime resources and responsible for rendering.
	Core *EngineCore

	//	Tracks various go:ngine performance counters over time.
	Stats *EngineStats

	//	Your gateway for end-user input (key, mouse etc.) and "output" (window management, not the graphics themselves).
	UserIO = newUserIO()
)

//	Call this to "un-init" go:ngine and to release any and all GPU or RAM resources still allocated.
func Dispose() {
	if Core != nil {
		Core.dispose()
	}
	glDispose()
	UserIO.dispose()
	Core, Loop, Stats = nil, nil, nil
}

//	Initializes go:ngine; this first attempts to initialize OpenGL and then open a window to your supplied specifications with a GL 3.2-or-higher profile.
func Init(options *EngineOptions, winTitle string) (err error) {
	isVerErr, forceContext := false, (runtime.GOOS == "darwin")
tryInit:
	if err = UserIO.init(options, winTitle, forceContext); err == nil {
		if err, isVerErr = glInit(); err == nil {
			Loop, Stats = newEngineLoop(), newEngineStats()
			newEngineCore(options)
		} else if isVerErr && !forceContext {
			forceContext = true
			UserIO.isGlfwInit, UserIO.isGlfwWindow = false, false
			goto tryInit
		}
	}
	return
}

func logError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

//	A short-hand for fmt.Sprintf. Feel free to ignore.
func Sfmt(format string, fmtArgs ...interface{}) string {
	return fmt.Sprintf(format, fmtArgs...)
}
