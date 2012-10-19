package nglfw

import (
	glfw "github.com/jteeuwen/glfw"

	util "github.com/go3d/go-util"

	nclient "github.com/go3d/go-ngine/client"
)

var (
	curGlfw *TWindowing
)

type TWindowing struct {
	KeyToggleMinDelay float64

	isGlfwInit, isGlfwWindow, togglePress bool
	lastToggles map[int]float64
	winWidth, winHeight int
}

func New () *TWindowing {
	curGlfw = &TWindowing {}
	curGlfw.lastToggles = map[int]float64 {}
	curGlfw.KeyToggleMinDelay = 0.25
	return curGlfw
}

func (me *TWindowing) Exit () {
	if (me.isGlfwWindow) {
		me.isGlfwWindow = false
		glfw.CloseWindow()
	}
	if (me.isGlfwInit) {
		me.isGlfwInit = false
		glfw.Terminate()
	}
}

func (me *TWindowing) Init (winWidth, winHeight int, winFullScreen bool, vsync int, winTitle string) error {
	var err error
	if (!me.isGlfwInit) {
		if err = glfw.Init(); err == nil {
			me.isGlfwInit = true
		}
	}
	if (me.isGlfwInit && !me.isGlfwWindow) {
		glfw.OpenWindowHint(glfw.FsaaSamples, 0) // AA is a pluggable post-processing shader, because super-/multi-sampling is simply highly undesirable
		if err = glfw.OpenWindow(winWidth, winHeight, 8, 8, 8, 0, 24, 8, util.Ifi(winFullScreen, glfw.Fullscreen, glfw.Windowed)); err == nil {
			me.winWidth, me.winHeight = winWidth, winHeight
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
	nclient.Loop.IsLooping = false
	return 1
}

func glfwOnWindowResize (width, height int) {
	curGlfw.winWidth, curGlfw.winHeight = width, height
	if (nclient.Core != nil) { nclient.Core.ResizeView(width, height) }
}

func (me *TWindowing) IifKey (key int, ifTrue, ifFalse float64) float64 {
	if me.KeyPressed(key) { return ifTrue }
	return ifFalse
}

func (me *TWindowing) KeyPressed (key int) bool {
	return glfw.Key(key) == glfw.KeyPress
}

func (me *TWindowing) KeysPressedAll2 (k1, k2 int) bool {
	return me.KeyPressed(k1) && me.KeyPressed(k2)
}

func (me *TWindowing) KeysPressedAll3 (k1, k2, k3 int) bool {
	return me.KeyPressed(k1) && me.KeyPressed(k2) && me.KeyPressed(k3)
}

func (me *TWindowing) KeysPressedAny2 (k1, k2 int) bool {
	return me.KeyPressed(k1) || me.KeyPressed(k2)
}

func (me *TWindowing) KeysPressedAny3 (k1, k2, k3 int) bool {
	return me.KeyPressed(k1) || me.KeyPressed(k2) || me.KeyPressed(k3)
}

func (me *TWindowing) KeyToggled (key int) bool {
	if me.togglePress = me.KeyPressed(key); me.togglePress && ((nclient.Loop.TickNow - me.lastToggles[key]) > me.KeyToggleMinDelay) {
		me.lastToggles[key] = nclient.Loop.TickNow
		return true
	}
	return false
}

func (me *TWindowing) OnLoop () {
	if (glfw.Key(glfw.KeyEsc) == glfw.KeyPress) || (glfw.WindowParam(glfw.Opened) != 1) {
		nclient.Loop.IsLooping = false
	} else {
		glfw.SwapBuffers()
	}
}

func (me *TWindowing) SetTitle (newTitle string) {
	glfw.SetWindowTitle(newTitle)
}

func (me *TWindowing) Time () float64 {
	return glfw.Time()
}

func (me *TWindowing) WinHeight () int {
	return me.winHeight
}

func (me *TWindowing) WinWidth () int {
	return me.winWidth
}
