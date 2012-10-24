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

func glInit() error {
	var minMatch = "3_2"
	var err error
	var vPos int
	var vMatch = "VERSION_"
	var vMessage = `Minimum required OpenGL version is %v, but
your graphics-card driver's OpenGL version is: %v.

*HOW TO FIX THIS*:
Simply visit the <%v> website.
Look for their "driver downloads" pages and follow their
instructions to find & download the newest driver version
for: <%v>.
`
	if !glIsInit {
		if err = gl.Init(); err != nil {
			// 	check for a message such as "unable to initialize VERSION_4_0"
			if vPos = strings.Index(err.Error(), vMatch); vPos >= 0 {
				vMatch = err.Error()[vPos+len(vMatch):]
				if vMatch > minMatch {
					err = nil
				} else {
					if vMatch > "1_0" { vMatch = glutil.GlStr(gl.VERSION) }
					err = fmt.Errorf(vMessage, strings.Replace(minMatch, "_", ".", -1), strings.Replace(vMatch, "_", ".", -1), glutil.GlStr(gl.VENDOR), glutil.GlStr(gl.RENDERER))
				}
			}
		}
		if err == nil {
			glutil.SetVersion()
			glIsInit = true
			gl.ClearColor(0, 0.05, 0.25, 1)
			gl.Enable(gl.DEPTH_TEST)
			gl.FrontFace(gl.CCW)
			gl.CullFace(gl.BACK)
			gl.Enable(gl.CULL_FACE)
			if err = glShaderMan.compileAll(); err == nil {
			}
		}
		if err == nil { err = glutil.LastError("nglcore.Init") }
	}
	return err
}

func glLogLastError(step string, fmtArgs ...interface{}) {
	var err = glutil.LastError(step, fmtArgs...)
	if err != nil { log.Println(err.Error()) }
}
