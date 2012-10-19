package client

import (
	glfw "github.com/jteeuwen/glfw"

	util "github.com/go3d/go-util"
)

type TWindowing struct {
	WinWidth, WinHeight int
	OnWinResizing, OnWinResized func (int, int)
	KeyToggleMinDelay float64

	isGlfwInit, isGlfwWindow, togglePress bool
	lastToggles map[int]float64
}

func newWindowing () *TWindowing {
	var win = &TWindowing {}
	win.lastToggles = map[int]float64 {}
	win.KeyToggleMinDelay = 0.25
	return win
}

func (me *TWindowing) glfwExit () {
	if (me.isGlfwWindow) {
		me.isGlfwWindow = false
		glfw.CloseWindow()
	}
	if (me.isGlfwInit) {
		me.isGlfwInit = false
		glfw.Terminate()
	}
}

func (me *TWindowing) glfwInit (winWidth, winHeight int, winFullScreen bool, vsync int, winTitle string) error {
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

func (me *TWindowing) glfwOnLoop () {
	if (glfw.Key(glfw.KeyEsc) == glfw.KeyPress) || (glfw.WindowParam(glfw.Opened) != 1) {
		Loop.IsLooping = false
	} else {
		glfw.SwapBuffers()
	}
}

func (me *TWindowing) IfKey (key int, ifTrue, ifFalse float64) float64 {
	if me.KeyPressed(key) { return ifTrue }
	return ifFalse
}

func (me *TWindowing) KeyPressed (key int) bool {
	return glfw.Key(key) == glfw.KeyPress
}

func (me *TWindowing) KeysPressedAll2 (k1, k2 int) bool {
	return me.KeyPressed(k1) && me.KeyPressed(k2)
}

func (me *TWindowing) KeysPressedAny2 (k1, k2 int) bool {
	return me.KeyPressed(k1) || me.KeyPressed(k2)
}

func (me *TWindowing) KeysPressedAll3 (k1, k2, k3 int) bool {
	return me.KeyPressed(k1) && me.KeyPressed(k2) && me.KeyPressed(k3)
}

func (me *TWindowing) KeysPressedAny3 (k1, k2, k3 int) bool {
	return me.KeyPressed(k1) || me.KeyPressed(k2) || me.KeyPressed(k3)
}

func (me *TWindowing) KeyToggled (key int) bool {
	if me.togglePress = me.KeyPressed(key); me.togglePress && ((Loop.TickNow - me.lastToggles[key]) > me.KeyToggleMinDelay) {
		me.lastToggles[key] = Loop.TickNow
		return true
	}
	return false
}

func (me *TWindowing) SetTitle (newTitle string) {
	glfw.SetWindowTitle(newTitle)
}

func glfwOnWindowClose () int {
	Loop.IsLooping = false
	return 1
}

func glfwOnWindowResize (width, height int) {
	if (Windowing.OnWinResizing != nil) { Windowing.OnWinResizing(width, height) }
	Windowing.WinWidth, Windowing.WinHeight = width, height
	if (Core != nil) { Core.ResizeView(width, height) }
	if (Windowing.OnWinResized != nil) { Windowing.OnWinResized(Windowing.WinWidth, Windowing.WinHeight) }
}
