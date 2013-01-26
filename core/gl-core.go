package core

import (
	"fmt"
	"log"
	"strings"
	"time"

	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

var (
	glIsInit  bool
	glProgMan = ugl.NewProgramManager()
)

func glDispose() {
	if glIsInit {
		glIsInit = false
		glProgMan.Dispose()
	}
}

func glInit() (err error, isVerErr bool) {
	const (
		minMatch = "3_3"
		vMessage = `Minimum required OpenGL version is %v, but your
graphics-card driver (or your OS) currently
only provides OpenGL version %v.

*HOW TO FIX THIS*:
Simply visit the <%v> website.
Look for their "driver downloads" pages and follow their
instructions to find & download the newest driver version
for: <%v>.
`
	)
	vMatch := "VERSION_"
	makeVerErr := func(curVer string) error {
		isVerErr = true
		return fmt.Errorf(vMessage, strings.Replace(minMatch, "_", ".", -1), curVer, ugl.Gl.Str(gl.VENDOR), ugl.Gl.Str(gl.RENDERER))
	}
	if !glIsInit {
		if err = gl.Init(); err != nil {
			// 	check for a message such as "unable to initialize VERSION_4_0"
			if vPos := strings.Index(err.Error(), vMatch); vPos >= 0 {
				vMatch = err.Error()[vPos+len(vMatch):]
				if vMatch > minMatch {
					err = nil
				} else {
					if vMatch > "1_0" {
						vMatch = ugl.Gl.Str(gl.VERSION)
					}
					err = makeVerErr(strings.Replace(vMatch, "_", ".", -1))
				}
			}
		}
		if err == nil {
			ugl.Init()
			if !ugl.VersionMatch(3.3) {
				err = makeVerErr(Sfmt("%v.%v", ugl.Support.GlVersion.MajorMinor[0], ugl.Support.GlVersion.MajorMinor[1]))
			} else {
				var dur time.Duration
				gl.FrontFace(gl.CCW)
				gl.CullFace(gl.BACK)
				log.Println(ugl.Gl.ConnInfo())
				if dur, err = glProgMan.MakeAllProgramsFromRawSources(); err == nil {
					log.Printf("Total shader compilation time for all %v programs: %v\n", len(glProgMan.Programs), dur)
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
