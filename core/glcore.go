package core

import (
	"fmt"
	"log"
	"strings"

	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

var (
	glIsInit bool
	glShaderMan = newShaderManager()
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
	var makeVerErr = func (curVer string) error {
		isVerErr = true
		return fmt.Errorf(vMessage, strings.Replace(minMatch, "_", ".", -1), curVer, glutil.GlStr(gl.VENDOR), glutil.GlStr(gl.RENDERER))
	}
	if !glIsInit {
		if err = gl.Init(); err != nil {
			// 	check for a message such as "unable to initialize VERSION_4_0"
			if vPos = strings.Index(err.Error(), vMatch); vPos >= 0 {
				vMatch = err.Error()[vPos+len(vMatch):]
				if vMatch > minMatch {
					err = nil
				} else {
					if vMatch > "1_0" { vMatch = glutil.GlStr(gl.VERSION) }
					err = makeVerErr(strings.Replace(vMatch, "_", ".", -1))
				}
			}
		}
		if err == nil {
			if glutil.SetVersion(); !glutil.VersionMatch(3, 2) {
				err = makeVerErr(fmt.Sprintf("%v.%v", glutil.Version[0], glutil.Version[1]))
			} else {
				glIsInit = true
				gl.ClearColor(0, 0.05, 0.25, 1)
				gl.Enable(gl.DEPTH_TEST)
				gl.FrontFace(gl.CCW)
				gl.CullFace(gl.BACK)
				gl.Enable(gl.CULL_FACE)
				log.Println(glutil.GlConnInfo())
				if err = glShaderMan.compileAll(); err == nil {
				}
			}
		}
		if err == nil { err = glutil.LastError("nglcore.Init") }
	}
	return
}

func glLogLastError(step string, fmtArgs ... interface {}) {
	LogError(glutil.LastError(step, fmtArgs ...))
}
