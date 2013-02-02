package core

import (
	"time"

	gl "github.com/go3d/go-opengl/core"

	ugl "github.com/go3d/go-opengl/util"
)

var (
	glc struct {
		isInit  bool
		progMan ugl.ProgramManager
	}
)

func glDispose() {
	if glc.isInit {
		glc.isInit = false
		glc.progMan.Reset()
	}
}

func glInit() (err error, isVerErr bool) {
	const (
		vMessage = `Minimum required OpenGL version is %s, but your
graphics driver (or your computer's OS) currently
only provides OpenGL version %s.

Most likely your computer is just missing some
recent system updates.

*HOW TO RESOLVE*:
Google for "how to find and download the latest
driver for %s",
or simply visit the <%s> website,
look for their "driver downloads" pages and follow
their instructions to find & download the latest
driver for: <%s>.
`
	)
	makeVerErr := func(curVer string) error {
		isVerErr = true
		return fmtErr(vMessage, "3.3", curVer, ugl.Gl.Str(gl.VENDOR)+" "+ugl.Gl.Str(gl.RENDERER), ugl.Gl.Str(gl.VENDOR), ugl.Gl.Str(gl.RENDERER))
	}
	if !glc.isInit {
		if !gl.Util.Init() {
			err = fmtErr("OpenGL 3.3+ initialization failed.")
		} else {
			// return makeVerErr(fmtStr("%v.%v", ugl.Support.GlVersion.MajorMinor[0], ugl.Support.GlVersion.MajorMinor[1])), true
			ugl.Init()
			if !ugl.VersionMatch(3.3) {
				err = makeVerErr(fmtStr("%v.%v", ugl.Support.GlVersion.MajorMinor[0], ugl.Support.GlVersion.MajorMinor[1]))
			} else {
				var dur time.Duration
				gl.FrontFace(gl.CCW)
				gl.CullFace(gl.BACK)
				Diag.LogMisc(ugl.Gl.ConnInfo())
				if Core.Options.Rendering.PostFx.TextureRect {
					glc.progMan.Defines["PostFx_TextureRect"] = 1
				}
				if dur, err = glc.progMan.MakeProgramsFromRawSources(true); err == nil {
					Diag.LogShaders("Total shader compilation time for all %v programs: %v\n", len(glc.progMan.Programs), dur)
				}
			}
		}
		if err == nil {
			err = gl.Util.Error("ng-gl-core.Init")
		}
		if err == nil {
			glc.isInit = true
		}
	}
	return
}
