//	First off, the go:ngine 3 RULES of bug-free development:
//	
//	1. you do NOT set exported fields directly, those are provided for simplified read-access. If a field is
//	meant to be modifiable, there will be a *SetXyz()* method or it will be documented explicitly as directly
//	modifiable, but such instances will be rare.
//	
//	2. you do NOT instantiate exported struct types directly, as in "new(Material)" or "&Material {}". Many of
//	those are exported only for documentation, but are to be instantiated only inside the go:ngine package.
//	For types to be instantiated by package-external code (ie. your code), go:ngine packages provide constructor
//	functions -- often in collections, ie. mat := Materials.New("arg") instead of, say, mat := NewMaterial("arg").
//	
//	3. Those 2 are *default* assumptions and modes of operation -- there are a few "if you know what you're
//	doing" exemptions and those are (or will be) explicitly documented as such.
//	
//	________
//	
//	TODO pkg doc
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
func Dispose () {
	if Core != nil { Core.Dispose() }
	glDispose()
	UserIO.dispose()
	Core, Loop, Stats = nil, nil, nil
}

//	Initializes go:ngine; this first attempts to initialize OpenGL and then open a window to your supplied specifications with a GL 3.2-or-higher profile.
func Init (options *EngineOptions, winTitle string) error {
	var err error
	var isVerErr bool
	var forceContext = (runtime.GOOS == "darwin")
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
	return err
}

func logError (err error) {
	if err != nil { log.Println(err.Error()) }
}

//	A short-hand for fmt.Sprintf. Feel free to ignore.
func Sfmt (format string, fmtArgs ... interface {}) string {
	return fmt.Sprintf(format, fmtArgs ...)
}
