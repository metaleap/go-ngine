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
func Init() (err error) {
	var (
		glVerIndex = len(ugl.KnownVersions) - 1
		badVer     string
		glVer      float64
	)
	if len(Options.AppDir.Temp.BaseName) > 0 {
		for _, diagTmpDirName := range []string{Options.AppDir.Temp.ShaderSources} {
			if err = uio.ClearDirectory(Core.fileIO.resolveLocalFilePath(filepath.Join(Options.AppDir.Temp.BaseName, diagTmpDirName))); err != nil {
				return
			}
		}
	}
	if Options.Initialization.GlContext.CoreProfile.ForceFirst {
		for i, v := range ugl.KnownVersions {
			if v == Options.Initialization.GlContext.CoreProfile.VersionHint {
				glVerIndex = i
				break
			}
		}
	}
tryInit:
	glVer = ugl.KnownVersions[glVerIndex]
	if err = UserIO.init(glVer); err == nil {
		if err, badVer = glInit(); err == nil && len(badVer) == 0 {
			Stats.reset()
			Loop.init()
			Core.init()
			Diag.LogIfGlErr("INIT")
		} else if len(badVer) > 0 && !Options.Initialization.GlContext.CoreProfile.ForceFirst {
			Options.Initialization.GlContext.CoreProfile.ForceFirst = true
			UserIO.isGlfwInit, UserIO.Window.isCreated = false, false
			goto tryInit
		}
	} else if Options.Initialization.GlContext.CoreProfile.ForceFirst && (glVerIndex > 0) {
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
