package core

import (
	"runtime"

	glfw "github.com/go-gl/glfw"

	util "github.com/go3d/go-util"
)

type tEngineWindowing struct {
	KeyToggleMinDelay float64

	isGlfwInit, isGlfwWindow, togglePress bool
	lastToggles map[int]float64
	keyWhich int
}

func newWindowing () *tEngineWindowing {
	var curGlfw = &tEngineWindowing {}
	curGlfw.lastToggles = map[int]float64 {}
	curGlfw.KeyToggleMinDelay = 0.25
	return curGlfw
}

func (me *tEngineWindowing) dispose () {
	if (me.isGlfwWindow) {
		me.isGlfwWindow = false
		glfw.CloseWindow()
	}
	if (me.isGlfwInit) {
		me.isGlfwInit = false
		glfw.Terminate()
	}
}

func (me *tEngineWindowing) init (opt *tOptions, winTitle string) error {
	var err error
	if (!me.isGlfwInit) {
		if err = glfw.Init(); err == nil {
			me.isGlfwInit = true
		}
	}
	if (me.isGlfwInit && !me.isGlfwWindow) {
		glfw.OpenWindowHint(glfw.FsaaSamples, 0) // AA will be a pluggable post-processing shader
		if (runtime.GOOS == "darwin") {
			glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
			glfw.OpenWindowHint(glfw.OpenGLForwardCompat, 1)
			glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
			glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
		}
		if err = glfw.OpenWindow(opt.winWidth, opt.winHeight, 8, 8, 8, 0, 24, 8, util.Ifi(opt.winFullScreen, glfw.Fullscreen, glfw.Windowed)); err == nil {
			opt.winWidth, opt.winHeight = glfw.WindowSize()
			me.isGlfwWindow = true
		}
	}
	if (me.isGlfwWindow) {
		me.SetTitle(winTitle)
		glfw.SetSwapInterval(opt.winSwapInterval)
		glfw.SetWindowCloseCallback(glfwOnWindowClose)
		glfw.SetWindowSizeCallback(glfwOnWindowResize)
		// glfw.Disable(glfw.MouseCursor)
		glfw.Enable(glfw.StickyKeys)
	}
	return err
}

func glfwOnWindowClose () int {
	Loop.IsLooping = false
	return 1
}

func glfwOnWindowResize (width, height int) {
	if (Core != nil) { Core.resizeView(width, height) }
}

func (me *tEngineWindowing) IifKey (key int, ifTrue, ifFalse float64) float64 {
	if me.KeyPressed(key) { return ifTrue }
	return ifFalse
}

func (me *tEngineWindowing) KeyPressed (key int) bool {
	return glfw.Key(key) == glfw.KeyPress
}

func (me *tEngineWindowing) KeyPressedWhich (keys ... int) int {
	for _, me.keyWhich = range keys {
		if me.KeyPressed(me.keyWhich) {
			return me.keyWhich
		}
	}
	return 0
}

func (me *tEngineWindowing) KeysPressedAll2 (k1, k2 int) bool {
	return me.KeyPressed(k1) && me.KeyPressed(k2)
}

func (me *tEngineWindowing) KeysPressedAll3 (k1, k2, k3 int) bool {
	return me.KeyPressed(k1) && me.KeyPressed(k2) && me.KeyPressed(k3)
}

func (me *tEngineWindowing) KeysPressedAny2 (k1, k2 int) bool {
	return me.KeyPressed(k1) || me.KeyPressed(k2)
}

func (me *tEngineWindowing) KeysPressedAny3 (k1, k2, k3 int) bool {
	return me.KeyPressed(k1) || me.KeyPressed(k2) || me.KeyPressed(k3)
}

func (me *tEngineWindowing) KeyToggled (key int) bool {
	if me.togglePress = me.KeyPressed(key); me.togglePress && ((Loop.TickNow - me.lastToggles[key]) > me.KeyToggleMinDelay) {
		me.lastToggles[key] = Loop.TickNow
		return true
	}
	return false
}

func (me *tEngineWindowing) onLoop () {
	if (glfw.Key(glfw.KeyEsc) == glfw.KeyPress) || (glfw.WindowParam(glfw.Opened) != 1) {
		Loop.IsLooping = false
	} else {
		glfw.SwapBuffers()
	}
}

func (me *tEngineWindowing) SetTitle (newTitle string) {
	glfw.SetWindowTitle(newTitle)
}

func (me *tEngineWindowing) Time () float64 {
	return glfw.Time()
}

func (me *tEngineWindowing) WinHeight () int {
	return Core.Options.winHeight
}

func (me *tEngineWindowing) WinWidth () int {
	return Core.Options.winWidth
}
