package core

import (
	"path"
	"runtime"
	"time"

	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
	ugo "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
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
		shaderMan  uberShader
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
				badVer = strf("%v.%v", ugl.Support.GlVersion.MajorMinor[0], ugl.Support.GlVersion.MajorMinor[1])
				glc.lastBadVer = badVer
				return
			} else {
				var dur time.Duration
				gl.FrontFace(gl.CCW)
				gl.CullFace(gl.BACK)
				gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)
				Diag.LogMisc(ugl.Util.ConnInfo())
				if len(glc.progMan.Programs) > 0 {
					if dur, err = glcProgsMake(true); err == nil {
						Diag.LogShaders("Total shader compilation time for all %v auxiliary (non-ubershader) programs: %v\n", len(glc.progMan.Programs), dur)
						Stats.addProgCompile(len(glc.progMan.Programs), dur.Nanoseconds())
					}
				}
				if err == nil {
					glc.shaderMan.loadFromRawSources()
					glc.shaderMan.processFuncs()
				}
			}
		}
		if err == nil {
			err = ugl.Util.LastError("ng-gl-core.Init")
		}
		if err == nil {
			glc.isInit = true
		}
	}
	return
}

func glcProgsMake(forceAll bool, forceSome ...string) (dur time.Duration, err error) {
	dur, err = glc.progMan.MakeProgramsFromRawSources(forceAll, forceSome...)
	if len(Options.AppDir.Temp.ShaderSources) > 0 {
		var src string
		if len(forceSome) == 0 {
			forceSome = glc.progMan.Names
		}
		for ext, sources := range map[string]map[string]string{
			".glcs": glc.progMan.FinalRealSources.Compute,
			".glfs": glc.progMan.FinalRealSources.Fragment,
			".glgs": glc.progMan.FinalRealSources.Geometry,
			".glhs": glc.progMan.FinalRealSources.TessCtl,
			".glds": glc.progMan.FinalRealSources.TessEval,
			".glvs": glc.progMan.FinalRealSources.Vertex,
		} {
			for _, progName := range forceSome {
				if src = sources[progName]; len(src) > 0 {
					src = "/*\tTemp file written at runtime for diagnostic purposes.\n\tThe following is a runtime-generated GLSL source string,\n\texactly as it was sent to the GL for compilation and linking: */\n" + src
					if err = uio.WriteTextFile(Core.fileIO.resolveLocalFilePath(path.Join(Options.AppDir.Temp.BaseName, Options.AppDir.Temp.ShaderSources, progName+ext)), src); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

func glVersionErrorMessage(minVer, curVer string) string {
	return ustr.Replace(Options.Initialization.GlContext.BadVersionMessage, map[string]string{
		"{MINVER}": minVer,
		"{CURVER}": curVer,
		"{OS}":     ugo.OSName(runtime.GOOS),
		"{GPU}":    ugl.Util.Str(gl.RENDERER),
		"{VENDOR}": ugl.Util.Str(gl.VENDOR),
	})
}
