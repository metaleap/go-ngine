package core

import (
	"path"
	"runtime"
	"time"

	"github.com/metaleap/go-util-fs"
	"github.com/metaleap/go-util-misc"
	"github.com/metaleap/go-util-str"
	gl "github.com/metaleap/go-opengl/core"
	ugl "github.com/metaleap/go-opengl/util"
)

const (
	glMinVer    = 3.3
	glMinVerStr = "3.3"
)

type glMan struct {
	lastBadVer string
	isInit     bool
	progs      ugl.ProgramManager
	uber       uberShaders
}

var ogl glMan

func init() {
	ogl.lastBadVer = "0"
}

func (_ *glMan) dispose() {
	if ogl.isInit {
		ogl.isInit = false
		ogl.progs.Dispose()
	}
}

func (_ *glMan) init() (err error, badVer string) {
	if !ogl.isInit {
		if !gl.Util.Init() {
			badVer = ogl.lastBadVer
			return
		} else {
			ugl.Init()
			if !ugl.VersionMatch(glMinVer) {
				badVer = strf("%v.%v", ugl.Support.GlVersion.MajorMinor[0], ugl.Support.GlVersion.MajorMinor[1])
				ogl.lastBadVer = badVer
				return
			} else {
				var dur time.Duration
				gl.FrontFace(gl.CCW)
				gl.CullFace(gl.BACK)
				gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)
				Diag.LogMisc(ugl.Util.ConnInfo())
				if len(ogl.progs.All) > 0 {
					if dur, err = ogl.makeProgs(true); err == nil {
						Diag.LogShaders("Total shader compilation time for all %v auxiliary (non-ubershader) programs: %v\n", len(ogl.progs.All), dur)
						Stats.addProgCompile(len(ogl.progs.All), dur.Nanoseconds())
					}
				}
				if err == nil {
					ogl.uber.loadFromRawSources()
					ogl.uber.processFuncs()
				}
			}
		}
		if err == nil {
			err = ugl.Util.LastError("ng-gl-core.Init")
		}
		if err == nil {
			ogl.isInit = true
		}
	}
	return
}

func (_ *glMan) makeProgs(forceAll bool, forceSome ...string) (dur time.Duration, err error) {
	var progsMade []bool
	dur, progsMade, err = ogl.progs.MakeProgramsFromRawSources(nil, forceAll, forceSome...)
	if len(Options.AppDir.Temp.ShaderSources) > 0 {
		writeSrc := func(i int, ext, src string) {
			if len(src) > 0 {
				ufs.WriteTextFile(Core.fileIO.resolveLocalFilePath(path.Join(Options.AppDir.Temp.BaseName, Options.AppDir.Temp.ShaderSources, ogl.progs.All[i].Name+ext)), src)
			}
		}
		for i := 0; i < len(ogl.progs.All); i++ {
			if progsMade[i] {
				writeSrc(i, ".glcs", ogl.progs.All[i].Sources.Out.Compute)
				writeSrc(i, ".glfs", ogl.progs.All[i].Sources.Out.Fragment)
				writeSrc(i, ".glgs", ogl.progs.All[i].Sources.Out.Geometry)
				writeSrc(i, ".glhs", ogl.progs.All[i].Sources.Out.TessCtl)
				writeSrc(i, ".glds", ogl.progs.All[i].Sources.Out.TessEval)
				writeSrc(i, ".glvs", ogl.progs.All[i].Sources.Out.Vertex)
			}
		}
	}
	return
}

func (_ *glMan) versionErrorMessage(minVer, curVer string) string {
	return ustr.Replace(Options.Initialization.GlContext.BadVersionMessage, map[string]string{
		"{MINVER}": minVer,
		"{CURVER}": curVer,
		"{OS}":     ugo.OSName(runtime.GOOS),
		"{GPU}":    ugl.Util.Str(gl.RENDERER),
		"{VENDOR}": ugl.Util.Str(gl.VENDOR),
	})
}
