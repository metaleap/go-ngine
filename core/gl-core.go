package core

import (
	"runtime"
	"time"

	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
	ugo "github.com/metaleap/go-util"
	ustr "github.com/metaleap/go-util/str"
)

const (
	glMinVer    = 3.3
	glMinVerStr = "3.3"
)

var (
	glc struct {
		lastBadVer string
		isInit     bool
		progMan    ugl.ProgramManager
		shaderMan  shaderComposer
	}
)

func init() {
	glc.lastBadVer = "0"
}

func glDispose() {
	if glc.isInit {
		glc.isInit = false
		glc.progMan.Reset()
	}
}

func glInit() (err error, badVer string) {
	if !glc.isInit {
		if !gl.Util.Init() {
			badVer = glc.lastBadVer
			return
		} else {
			ugl.Init()
			if !ugl.VersionMatch(glMinVer) {
				badVer = fmtStr("%v.%v", ugl.Support.GlVersion.MajorMinor[0], ugl.Support.GlVersion.MajorMinor[1])
				glc.lastBadVer = badVer
				return
			} else {
				var dur time.Duration
				gl.FrontFace(gl.CCW)
				gl.CullFace(gl.BACK)
				Diag.LogMisc(ugl.Util.ConnInfo())
				if dur, err = glc.progMan.MakeProgramsFromRawSources(true); err == nil {
					Diag.LogShaders("Total shader compilation time for all %v programs: %v\n", len(glc.progMan.Programs), dur)
				}
			}
		}
		if err == nil {
			err = ugl.Util.Error("ng-gl-core.Init")
		}
		if err == nil {
			glc.isInit = true
		}
	}
	return
}

func glVersionErrorMessage(minVer, curVer string) string {
	return ustr.Replace(Core.Options.Initialization.GlContext.BadVersionMessage, map[string]string{
		"{MINVER}": minVer,
		"{CURVER}": curVer,
		"{OS}":     ugo.OSName(runtime.GOOS),
		"{GPU}":    ugl.Util.Str(gl.RENDERER),
		"{VENDOR}": ugl.Util.Str(gl.VENDOR),
	})
}
