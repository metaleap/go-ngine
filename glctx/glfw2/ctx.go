//	Implements a `CtxProvider` for GLFW 2.x.
package glctx_glfw2

import (
	glfw "github.com/go-gl/glfw"
	ngctx "github.com/metaleap/go-ngine/glctx"
)

type context struct {
}

//	Returns a new `CtxProvider` for GLFW 2.x.
func New() ngctx.CtxProvider {
	return &context{}
}

func (me *context) Hint(flag, value int) {
	glfw.OpenWindowHint(flag, value)
}

func (me *context) Init() (err error) {
	if err = glfw.Init(); err == nil {
		glfw.Disable(glfw.StickyKeys)
		glfw.Disable(glfw.AutoPollEvents)
	}
	return
}

func (me *context) Window(winf *ngctx.WinProfile, bufSize *ngctx.BufferBits, ctxProf *ngctx.CtxProfile) (win ngctx.Window, err error) {
	glfw.OpenWindowHint(glfw.FsaaSamples, 0) // AA will be a pluggable post-processing shader
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, ctxProf.Version.Major)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, ctxProf.Version.Minor)
	glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	if ctxProf.ForwardCompat {
		glfw.OpenWindowHint(glfw.OpenGLForwardCompat, 1)
	}
	winMode := glfw.Windowed
	if winf.FullScreen {
		winMode = glfw.Fullscreen
	}
	if err = glfw.OpenWindow(winf.Width, winf.Height, bufSize.Color.R, bufSize.Color.G, bufSize.Color.B, bufSize.Color.A, bufSize.Depth, bufSize.Stencil, winMode); err == nil {
		win = newWindow()
		win.SetTitle(winf.Title)
		if winMode == glfw.Fullscreen {
			glfw.Disable(glfw.MouseCursor)
		}
	}
	return
}

func (me *context) PollEvents() {
	glfw.PollEvents()
}

func (me *context) SetSwapInterval(interval int) {
	glfw.SetSwapInterval(interval)
}

func (me *context) SetTime(t float64) {
	glfw.SetTime(t)
}

func (me *context) Terminate() {
	glfw.Terminate()
}

func (me *context) Time() float64 {
	return glfw.Time()
}
