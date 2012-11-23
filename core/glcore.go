package core

import (
	"fmt"
	"log"
	"strings"

	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

var (
	glBackfaceCulling = true
	glIsInit          bool
	glShaderMan       = newShaderManager()
)

func glDispose() {
	if glIsInit {
		glIsInit = false
		glShaderMan.dispose()
	}
}

func glInit() (err error, isVerErr bool) {
	var minMatch = "3_2"
	var vPos int
	var vMatch = "VERSION_"
	var vMessage = `Minimum required OpenGL version is %v, but your
graphics-card driver (or your OS) currently
only provides OpenGL version %v.

*HOW TO FIX THIS*:
Simply visit the <%v> website.
Look for their "driver downloads" pages and follow their
instructions to find & download the newest driver version
for: <%v>.
`
	var makeVerErr = func(curVer string) error {
		isVerErr = true
		return fmt.Errorf(vMessage, strings.Replace(minMatch, "_", ".", -1), curVer, ugl.GlStr(gl.VENDOR), ugl.GlStr(gl.RENDERER))
	}
	if !glIsInit {
		if err = gl.Init(); err != nil {
			// 	check for a message such as "unable to initialize VERSION_4_0"
			if vPos = strings.Index(err.Error(), vMatch); vPos >= 0 {
				vMatch = err.Error()[vPos+len(vMatch):]
				if vMatch > minMatch {
					err = nil
				} else {
					if vMatch > "1_0" {
						vMatch = ugl.GlStr(gl.VERSION)
					}
					err = makeVerErr(strings.Replace(vMatch, "_", ".", -1))
				}
			}
		}
		if err == nil {
			if ugl.SetVersion(); !ugl.VersionMatch(3.2) {
				err = makeVerErr(Sfmt("%v.%v", ugl.GlVersionMajorMinor[0], ugl.GlVersionMajorMinor[1]))
			} else {
				gl.ClearColor(0, 0, 0, 1)
				gl.Enable(gl.DEPTH_TEST)
				gl.FrontFace(gl.CCW)
				gl.CullFace(gl.BACK)
				if glBackfaceCulling {
					gl.Enable(gl.CULL_FACE)
				}
				log.Println(ugl.GlConnInfo())
				if err = glShaderMan.compileAll(); err == nil {
				}
			}
		}
		if err == nil {
			err = ugl.LastError("nglcore.Init")
		}
		if err == nil {
			glIsInit = true
		}
	}
	return
}

func glLogLastError(step string, fmtArgs ...interface{}) {
	logError(ugl.LastError(step, fmtArgs...))
}

func glSetBackfaceCulling(val bool) {
	if val != glBackfaceCulling {
		if glBackfaceCulling = val; glBackfaceCulling {
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}
	}
}
