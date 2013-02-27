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
		uberShader uberShader
	}
)

func init() {
	glc.lastBadVer = "0"
}

func glDispose() {
	if glc.isInit {
		glc.isInit = false
		glc.progMan.Dispose()
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
					glc.uberShader.loadFromRawSources()
					glc.uberShader.processFuncs()
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
	var progsMade []bool
	dur, progsMade, err = glc.progMan.MakeProgramsFromRawSources(nil, forceAll, forceSome...)
	if len(Options.AppDir.Temp.ShaderSources) > 0 {
		writeSrc := func(i int, ext, src string) {
			if len(src) > 0 {
				uio.WriteTextFile(Core.fileIO.resolveLocalFilePath(path.Join(Options.AppDir.Temp.BaseName, Options.AppDir.Temp.ShaderSources, glc.progMan.Programs[i].Name+ext)), src)
			}
		}
		for i := 0; i < len(glc.progMan.Programs); i++ {
			if progsMade[i] {
				writeSrc(i, ".glcs", glc.progMan.Programs[i].Sources.Out.Compute)
				writeSrc(i, ".glfs", glc.progMan.Programs[i].Sources.Out.Fragment)
				writeSrc(i, ".glgs", glc.progMan.Programs[i].Sources.Out.Geometry)
				writeSrc(i, ".glhs", glc.progMan.Programs[i].Sources.Out.TessCtl)
				writeSrc(i, ".glds", glc.progMan.Programs[i].Sources.Out.TessEval)
				writeSrc(i, ".glvs", glc.progMan.Programs[i].Sources.Out.Vertex)
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
