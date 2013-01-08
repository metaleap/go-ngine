package core

import (
	"runtime"

	glfw "github.com/go-gl/glfw"

	util "github.com/metaleap/go-util"
)

//	Consider EngineUserIO a "Singleton" type, only valid use is the core.UserIO global variable.
//	Your gateway for end-user input (key, mouse etc.) and "output" (window management, not the graphics themselves).
type EngineUserIO struct {
	//	Minimum delay for EngineUserIO.KeyToggled() method, in seconds. Defaults to 0.25.
	KeyToggleMinDelay float64

	isGlfwInit, isGlfwWindow, togglePress bool
	lastToggles                           map[int]float64
	keyWhich                              int
}

func newUserIO() (me *EngineUserIO) {
	me = &EngineUserIO{lastToggles: map[int]float64{}, KeyToggleMinDelay: 0.25}
	return
}

func (me *EngineUserIO) dispose() {
	if me.isGlfwWindow {
		me.isGlfwWindow = false
		glfw.CloseWindow()
	}
	if me.isGlfwInit {
		me.isGlfwInit = false
		glfw.Terminate()
	}
}

func (me *EngineUserIO) init(opt *EngineOptions, winTitle string, forceContext bool) (err error) {
	if !me.isGlfwInit {
		if err = glfw.Init(); err == nil {
			me.isGlfwInit = true
		}
	}
	if me.isGlfwInit && !me.isGlfwWindow {
		glfw.OpenWindowHint(glfw.FsaaSamples, 0) // AA will be a pluggable post-processing shader
		if forceContext {
			glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
			glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
			glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
			if runtime.GOOS == "darwin" {
				glfw.OpenWindowHint(glfw.OpenGLForwardCompat, 1)
			}
		}
		if err = glfw.OpenWindow(opt.winWidth, opt.winHeight, 8, 8, 8, 0, 24, 8, util.Ifi(opt.winFullScreen, glfw.Fullscreen, glfw.Windowed)); err == nil {
			opt.winWidth, opt.winHeight = glfw.WindowSize()
			me.isGlfwWindow = true
		}
	}
	if me.isGlfwWindow {
		me.SetWinTitle(winTitle)
		glfw.SetSwapInterval(opt.winSwapInterval)
		glfw.SetWindowCloseCallback(glfwOnWindowClose)
		glfw.SetWindowSizeCallback(glfwOnWindowResize)
		// glfw.Disable(glfw.MouseCursor)
		glfw.Enable(glfw.StickyKeys)
	}
	return
}

func glfwOnWindowClose() int {
	Loop.IsLooping = false
	return 1
}

func glfwOnWindowResize(width, height int) {
	Core.resizeView(width, height)
}

//	Returns ifTrue if the specified key is pressed, otherwise returns ifFalse.
func (me *EngineUserIO) IifKeyF(key int, ifTrue, ifFalse float64) float64 {
	if me.KeyPressed(key) {
		return ifTrue
	}
	return ifFalse
}

//	Returns true if the specified key is pressed.
func (me *EngineUserIO) KeyPressed(key int) bool {
	return glfw.Key(key) == glfw.KeyPress
}

//	Returns the first in keys that is pressed.
func (me *EngineUserIO) KeyPressedWhich(keys ...int) int {
	for _, me.keyWhich = range keys {
		if me.KeyPressed(me.keyWhich) {
			return me.keyWhich
		}
	}
	return 0
}

//	Returns true if both specified keys are pressed.
func (me *EngineUserIO) KeysPressedAll2(k1, k2 int) bool {
	return me.KeyPressed(k1) && me.KeyPressed(k2)
}

//	Returns true if all three specified keys are pressed.
func (me *EngineUserIO) KeysPressedAll3(k1, k2, k3 int) bool {
	return me.KeyPressed(k1) && me.KeyPressed(k2) && me.KeyPressed(k3)
}

//	Returns true if any of the two specified keys is pressed.
func (me *EngineUserIO) KeysPressedAny2(k1, k2 int) bool {
	return me.KeyPressed(k1) || me.KeyPressed(k2)
}

//	Returns true if any of the three specified keys is pressed.
func (me *EngineUserIO) KeysPressedAny3(k1, k2, k3 int) bool {
	return me.KeyPressed(k1) || me.KeyPressed(k2) || me.KeyPressed(k3)
}

//	Returns true if the specified key has been "toggled", ie. its pressed-state changed within the last me.KeyToggleMinDelay seconds.
func (me *EngineUserIO) KeyToggled(key int) bool {
	if me.togglePress = me.KeyPressed(key); me.togglePress && ((Loop.TickNow - me.lastToggles[key]) > me.KeyToggleMinDelay) {
		me.lastToggles[key] = Loop.TickNow
		return true
	}
	return false
}

//	Sets the window title to newTitle.
func (me *EngineUserIO) SetWinTitle(newTitle string) {
	glfw.SetWindowTitle(newTitle)
}

//	Returns the height of the window in pixels.
func (me *EngineUserIO) WinHeight() int {
	return Core.Options.winHeight
}

//	Returns the width of the window in pixels.
func (me *EngineUserIO) WinWidth() int {
	return Core.Options.winWidth
}
