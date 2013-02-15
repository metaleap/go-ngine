package core

import (
	"fmt"
	"path/filepath"

	ugl "github.com/go3d/go-opengl/util"
	uio "github.com/metaleap/go-util/io"
)

//	Call this to "un-init" go:ngine and to release any and all GPU or RAM resources still allocated.
func Dispose() {
	Core.dispose()
	glDispose()
	UserIO.dispose()
}

//	Initializes go:ngine; this first attempts to initialize OpenGL and then open a window to your supplied specifications with a GL 3.3-or-higher profile.
func Init(options *EngineOptions) (err error) {
	var (
		glVerIndex = len(ugl.KnownVersions) - 1
		badVer     string
		glVer      float64
	)
	Core.Options = *options
	options = nil
	if len(Diag.WriteTmpFilesTo.BaseDirName) > 0 {
		for _, diagTmpDirName := range []string{Diag.WriteTmpFilesTo.ShaderPrograms} {
			if err = uio.ClearDirectory(filepath.Join(Core.Options.AppDir.BasePath, Diag.WriteTmpFilesTo.BaseDirName, diagTmpDirName)); err != nil {
				return
			}
		}
	}
tryInit:
	if Core.Options.Initialization.GlContext.CoreProfile.ForceFirst {
		for i, v := range ugl.KnownVersions {
			if v == Core.Options.Initialization.GlContext.CoreProfile.VersionHint {
				glVerIndex = i
				break
			}
		}
		glVer = ugl.KnownVersions[glVerIndex]
	}
	if err = UserIO.init(glVer); err == nil {
		if err, badVer = glInit(); err == nil && len(badVer) == 0 {
			Stats.reset()
			Loop.init()
			Core.init()
			ugl.LogLastError("INIT")
		} else if len(badVer) > 0 && !Core.Options.Initialization.GlContext.CoreProfile.ForceFirst {
			Core.Options.Initialization.GlContext.CoreProfile.ForceFirst = true
			UserIO.isGlfwInit, UserIO.Window.isCreated = false, false
			goto tryInit
		}
	} else if Core.Options.Initialization.GlContext.CoreProfile.ForceFirst && (glVerIndex > 0) {
		glVerIndex--
		UserIO.isGlfwInit, UserIO.Window.isCreated = false, false
		goto tryInit
	} else {
		badVer = glc.lastBadVer
	}
	if len(badVer) > 0 {
		err = errf(glVersionErrorMessage(glMinVerStr, badVer))
	}
	return
}

func errf(format string, fmtArgs ...interface{}) error {
	return fmt.Errorf(format, fmtArgs...)
}

func strf(format string, fmtArgs ...interface{}) string {
	return fmt.Sprintf(format, fmtArgs...)
}
