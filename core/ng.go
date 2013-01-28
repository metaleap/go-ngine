package core

import (
	"fmt"

	ugl "github.com/go3d/go-glutil"
)

func init() {
	Diag.LogCategories = LogCatAll
}

//	Call this to "un-init" go:ngine and to release any and all GPU or RAM resources still allocated.
func Dispose() {
	Core.dispose()
	glDispose()
	UserIO.dispose()
}

//	Initializes go:ngine; this first attempts to initialize OpenGL and then open a window to your supplied specifications with a GL 3.3-or-higher profile.
func Init(options *EngineOptions, initialWinTitle string) (err error) {
	var (
		isVerErr   = false
		glVerIndex = len(ugl.KnownVersions) - 1
		glVer      float64
	)
	Core.Options = *options
tryInit:
	if options.Initialization.GlCoreContext {
		glVer = ugl.KnownVersions[glVerIndex]
	}
	if err = UserIO.init(options, initialWinTitle, glVer); err == nil {
		if err, isVerErr = glInit(); err == nil {
			Stats.reset()
			Loop.init()
			Core.init()
			ugl.LogLastError("INIT")
		} else if isVerErr && !options.Initialization.GlCoreContext {
			options.Initialization.GlCoreContext = true
			UserIO.isGlfwInit, UserIO.isGlfwWindow = false, false
			goto tryInit
		}
	} else if options.Initialization.GlCoreContext && (glVerIndex > 0) {
		glVerIndex--
		UserIO.isGlfwInit, UserIO.isGlfwWindow = false, false
		goto tryInit
	}
	return
}

func fmtErr(format string, fmtArgs ...interface{}) error {
	return fmt.Errorf(format, fmtArgs...)
}

func fmtStr(format string, fmtArgs ...interface{}) string {
	return fmt.Sprintf(format, fmtArgs...)
}
