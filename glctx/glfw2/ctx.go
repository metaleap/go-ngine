package glctx_glfw2

import (
	glfw "github.com/go-gl/glfw"
	ngctx "github.com/go3d/go-ngine/glctx"
)

// if FS glfw.Disable(glfw.MouseCursor)

type context struct {
}

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

func (me *context) Window(width, height int, title string, bufSize *ngctx.BufferBits, ctxProf *ngctx.CtxProfile) (win ngctx.Window, err error) {
	glfw.OpenWindowHint(glfw.FsaaSamples, 0) // AA will be a pluggable post-processing shader
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, ctxProf.Version.Major)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, ctxProf.Version.Minor)
	glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	if ctxProf.ForwardCompatibility {
		glfw.OpenWindowHint(glfw.OpenGLForwardCompat, 1)
	}
	if err = glfw.OpenWindow(width, height, bufSize.Color.R, bufSize.Color.G, bufSize.Color.B, bufSize.Color.A, bufSize.Depth, bufSize.Stencil, glfw.Windowed); err == nil {
		win = newWindow()
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
