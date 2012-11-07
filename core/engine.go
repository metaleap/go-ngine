package core

import (
	"fmt"
	"log"
)

var (
	Loop *engineLoop
	Core *engineCore
	Stats *engineStats
	UserIO = newUserIO()
)

func Dispose () {
	if Core != nil { Core.Dispose() }
	glDispose()
	UserIO.dispose()
	Core, Loop, Stats = nil, nil, nil
}

func Init (options *engineOptions, winTitle string) error {
	var err error
	var isVerErr bool
	var forceContext = false
	tryInit:
	if err = UserIO.init(options, winTitle, forceContext); err == nil {
		if err, isVerErr = glInit(); err == nil {
			Loop, Stats, Core = newEngineLoop(), newEngineStats(), newEngineCore(options)
		} else if isVerErr && !forceContext {
			forceContext = true
			UserIO.isGlfwInit, UserIO.isGlfwWindow = false, false
			goto tryInit
		}
	}
	return err
}

func Fmt (format string, fmtArgs ... interface {}) string {
	return fmt.Sprintf(format, fmtArgs ...)
}

func LogError (err error) {
	if err != nil { log.Println(err.Error()) }
}

func LogLastGlError (step string, fmtArgs ... interface {}) {
	glLogLastError(step, fmtArgs ...)
}
