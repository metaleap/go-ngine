package core

import (
	glfw "github.com/jteeuwen/glfw"

	util "github.com/go3d/go-util"
)

type tEngineWindowing struct {
	KeyToggleMinDelay float64
	WinWidth, WinHeight int

	isGlfwInit, isGlfwWindow, togglePress bool
	lastToggles map[int]float64
}

func newWindowing () *tEngineWindowing {
	var curGlfw = &tEngineWindowing {}
	curGlfw.lastToggles = map[int]float64 {}
	curGlfw.KeyToggleMinDelay = 0.25
	return curGlfw
}

func (me *tEngineWindowing) Exit () {
	if (me.isGlfwWindow) {
		me.isGlfwWindow = false
		glfw.CloseWindow()
	}
	if (me.isGlfwInit) {
		me.isGlfwInit = false
		glfw.Terminate()
	}
}

func (me *tEngineWindowing) Init (winWidth, winHeight int, winFullScreen bool, vsync int, winTitle string) error {
	var err error
	if (!me.isGlfwInit) {
		if err = glfw.Init(); err == nil {
			me.isGlfwInit = true
		}
	}
	if (me.isGlfwInit && !me.isGlfwWindow) {
		glfw.OpenWindowHint(glfw.FsaaSamples, 0) // AA is a pluggable post-processing shader, because super-/multi-sampling is simply highly undesirable
		if err = glfw.OpenWindow(winWidth, winHeight, 8, 8, 8, 0, 24, 8, util.Ifi(winFullScreen, glfw.Fullscreen, glfw.Windowed)); err == nil {
			me.WinWidth, me.WinHeight = winWidth, winHeight
			me.isGlfwWindow = true
		}
	}
	if (me.isGlfwWindow) {
		me.SetTitle(winTitle)
		glfw.SetSwapInterval(vsync)
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
	Windowing.WinWidth, Windowing.WinHeight = width, height
	if (Core != nil) { Core.ResizeView(width, height) }
}

func (me *tEngineWindowing) IifKey (key int, ifTrue, ifFalse float64) float64 {
	if me.KeyPressed(key) { return ifTrue }
	return ifFalse
}

func (me *tEngineWindowing) KeyPressed (key int) bool {
	return glfw.Key(key) == glfw.KeyPress
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

func (me *tEngineWindowing) OnLoop () {
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
