package glctx_glfw3

import (
	"fmt"

	glfw "github.com/go-gl/glfw3"
	ngctx "github.com/go3d/go-ngine/glctx"
)

// if FS glfw.Disable(glfw.MouseCursor)

type context struct {
	lastErr struct {
		code glfw.ErrorCode
		msg  string
	}
}

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

func (me *context) Window(width, height int, title string, bufSize *ngctx.BufferBits, ctxProf *ngctx.CtxProfile) (window ngctx.Window, err error) {
	glfw.WindowHint(glfw.Samples, 0) // AA will be a pluggable post-processing shader
	glfw.WindowHint(glfw.RedBits, bufSize.Color.R)
	glfw.WindowHint(glfw.GreenBits, bufSize.Color.G)
	glfw.WindowHint(glfw.BlueBits, bufSize.Color.B)
	glfw.WindowHint(glfw.AlphaBits, bufSize.Color.A)
	glfw.WindowHint(glfw.DepthBits, bufSize.Depth)
	glfw.WindowHint(glfw.StencilBits, bufSize.Stencil)
	glfw.WindowHint(glfw.ContextVersionMajor, ctxProf.Version.Major)
	glfw.WindowHint(glfw.ContextVersionMinor, ctxProf.Version.Minor)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	if ctxProf.ForwardCompatibility {
		glfw.WindowHint(glfw.OpenglForwardCompatible, 1)
	}
	var win *glfw.Window
	if win, err = glfw.CreateWindow(width, height, title, nil, nil); err == nil {
		window = newWindow(win)
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
