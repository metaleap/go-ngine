//	Implements a `CtxProvider` for GLFW 3.x.
package glctx_glfw3

import (
	"fmt"

	glfw "github.com/go-gl/glfw/v3.0/glfw"
	ngctx "github.com/metaleap/go-ngine/glctx"
)

type context struct {
	lastErr struct {
		code glfw.ErrorCode
		msg  string
	}
}

//	Returns a new `CtxProvider` for GLFW 3.x.
func New() ngctx.CtxProvider {
	me := &context{}
	me.lastErr.code = -1
	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		me.lastErr.code, me.lastErr.msg = err, desc
	})
	return me
}

func (me *context) Hint(flag, value int) {
	glfw.WindowHint(glfw.Hint(flag), value)
}

func (me *context) Init() (err error) {
	if !glfw.Init() {
		err = fmt.Errorf("GLFW3.Init() failed: #%v %s", me.lastErr.code, me.lastErr.msg)
	}
	return
}

func (me *context) Window(winf *ngctx.WinProfile, bufSize *ngctx.BufferBits, ctxProf *ngctx.CtxProfile) (window ngctx.Window, err error) {
	glfw.WindowHint(glfw.Samples, winf.MultiSampling)
	glfw.WindowHint(glfw.RedBits, bufSize.Color.R)
	glfw.WindowHint(glfw.GreenBits, bufSize.Color.G)
	glfw.WindowHint(glfw.BlueBits, bufSize.Color.B)
	glfw.WindowHint(glfw.AlphaBits, bufSize.Color.A)
	glfw.WindowHint(glfw.DepthBits, bufSize.Depth)
	glfw.WindowHint(glfw.StencilBits, bufSize.Stencil)
	glfw.WindowHint(glfw.ContextVersionMajor, ctxProf.Version.Major)
	glfw.WindowHint(glfw.ContextVersionMinor, ctxProf.Version.Minor)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	if ctxProf.ForwardCompat {
		glfw.WindowHint(glfw.OpenglForwardCompatible, 1)
	}
	var mon *glfw.Monitor
	if winf.FullScreen {
		mon, err = glfw.GetPrimaryMonitor()
	}
	if err == nil {
		var win *glfw.Window
		if win, err = glfw.CreateWindow(winf.Width, winf.Height, winf.Title, mon, nil); win != nil {
			window = newWindow(win)
			if winf.FullScreen {
				win.SetInputMode(glfw.Cursor, glfw.CursorHidden)
			}
		}
	}
	return
}

func (me *context) PollEvents() {
	glfw.PollEvents()
}

func (me *context) SetSwapInterval(interval int) {
	glfw.SwapInterval(interval)
}

func (me *context) SetTime(t float64) {
	glfw.SetTime(t)
}

func (me *context) Terminate() {
	glfw.Terminate()
}

func (me *context) Time() float64 {
	return glfw.GetTime()
}
