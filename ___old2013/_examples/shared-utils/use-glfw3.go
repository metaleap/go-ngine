package exampleutils

import (
	glfw "github.com/go-gl/glfw/v3.0/glfw"
	glctx "github.com/metaleap/go-ngine/glctx"
	ngctx "github.com/metaleap/go-ngine/glctx/glfw3"
)

const (
	KeyLalt     int = int(glfw.KeyLeftAlt)
	KeyLshift   int = int(glfw.KeyLeftShift)
	KeyRshift   int = int(glfw.KeyRightShift)
	KeyLeft     int = int(glfw.KeyLeft)
	KeyRight    int = int(glfw.KeyRight)
	KeyUp       int = int(glfw.KeyUp)
	KeyDown     int = int(glfw.KeyDown)
	KeyPagedown int = int(glfw.KeyPageDown)
	KeyPageup   int = int(glfw.KeyPageUp)
	KeyKP9      int = int(glfw.KeyKp9)
	KeyKP3      int = int(glfw.KeyKp3)
	KeyEsc      int = int(glfw.KeyEscape)
	KeyF1       int = int(glfw.KeyF1)
	KeyF2       int = int(glfw.KeyF2)
	KeyF3       int = int(glfw.KeyF3)
	KeyF4       int = int(glfw.KeyF4)
	KeyF5       int = int(glfw.KeyF5)
	KeyF6       int = int(glfw.KeyF6)
	KeyF7       int = int(glfw.KeyF7)
	KeyF8       int = int(glfw.KeyF8)
	KeyF9       int = int(glfw.KeyF9)
	KeyF10      int = int(glfw.KeyF10)
	KeyF11      int = int(glfw.KeyF11)
	KeyF12      int = int(glfw.KeyF12)
	KeyLctrl    int = int(glfw.KeyLeftControl)
)

func newGlCtx() glctx.CtxProvider {
	return ngctx.New()
}
