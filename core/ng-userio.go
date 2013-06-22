package core

import (
	glfw "github.com/go-gl/glfw"
	ugo "github.com/go-utils/ugo"
	ugl "github.com/go3d/go-opengl/util"
)

var (
	//	Your gateway for end-user input (key, mouse etc.) and "output" (window management, not the graphics themselves).
	UserIO NgUserIO
)

//	Consider NgUserIO a "Singleton" type, only valid use is the core.UserIO global variable.
//	Your gateway for end-user input (key, mouse etc.) and "output" (window management, not the graphics themselves).
type NgUserIO struct {
	//	Minimum delay for NgUserIO.KeyToggled() method, in seconds. Defaults to 0.15.
	KeyToggleMinDelay float64

	Window WindowOptions

	isGlfwInit, togglePress bool
	keyWhich                int
	lastToggles             map[int]float64
}

func (_ *NgUserIO) dispose() {
	if UserIO.Window.isCreated {
		UserIO.Window.isCreated = false
		glfw.CloseWindow()
	}
	if UserIO.isGlfwInit {
		UserIO.isGlfwInit = false
		glfw.Terminate()
	}
}

func (_ *NgUserIO) init(forceContextVersion float64) (err error) {
	UserIO.KeyToggleMinDelay, UserIO.lastToggles = 0.15, make(map[int]float64, 80)
	if !UserIO.isGlfwInit {
		if err = glfw.Init(); err == nil {
			UserIO.isGlfwInit = true
		}
	}
	if UserIO.isGlfwInit && !UserIO.Window.isCreated {
		glfw.OpenWindowHint(glfw.FsaaSamples, 0) // AA will be a pluggable post-processing shader
		if forceContextVersion > 0 {
			major, minor := ugl.VersionMajorMinor(forceContextVersion)
			glfw.OpenWindowHint(glfw.OpenGLVersionMajor, major)
			glfw.OpenWindowHint(glfw.OpenGLVersionMinor, minor)
			glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
			if Options.Initialization.GlContext.CoreProfile.ForwardCompat {
				glfw.OpenWindowHint(glfw.OpenGLForwardCompat, 1)
			}
		}
		err = UserIO.recreateWin()
	}
	return
}

func (_ *NgUserIO) recreateWin() (err error) {
	winInit := &Options.Initialization.Window
	if UserIO.Window.isCreated {
		glfw.CloseWindow()
	}
	if UserIO.Window.isCreated, err = false, glfw.OpenWindow(UserIO.Window.width, UserIO.Window.height, winInit.Rbits, winInit.Gbits, winInit.Bbits, winInit.Abits, winInit.DepthBits, winInit.StencilBits, ugo.Ifi(UserIO.Window.fullscreen, glfw.Fullscreen, glfw.Windowed)); err == nil {
		UserIO.Window.width, UserIO.Window.height = glfw.WindowSize()
		UserIO.Window.isCreated = true
		UserIO.Window.SetTitle(UserIO.Window.title)
		UserIO.Window.SetSwapInterval(UserIO.Window.swap)
		glfw.SetWindowCloseCallback(glfwOnWindowClose)
		glfw.SetWindowSizeCallback(glfwOnWindowResize)
		// glfw.Disable(glfw.MouseCursor)
		glfw.Disable(glfw.AutoPollEvents)
		glfw.Disable(glfw.StickyKeys)
	}
	return
}

//	Returns ifTrue if the specified key is pressed, otherwise returns ifFalse.
func (_ *NgUserIO) IifKeyF(key int, ifTrue, ifFalse float64) float64 {
	if UserIO.KeyPressed(key) {
		return ifTrue
	}
	return ifFalse
}

//	Returns true if the specified key is pressed.
func (_ *NgUserIO) KeyPressed(key int) bool {
	return glfw.Key(key) == glfw.KeyPress
}

//	Returns the first in keys that is pressed.
func (_ *NgUserIO) KeyPressedWhich(keys ...int) int {
	for _, UserIO.keyWhich = range keys {
		if UserIO.KeyPressed(UserIO.keyWhich) {
			return UserIO.keyWhich
		}
	}
	return 0
}

//	Returns true if both specified keys are pressed.
func (_ *NgUserIO) KeysPressedAll2(k1, k2 int) bool {
	return UserIO.KeyPressed(k1) && UserIO.KeyPressed(k2)
}

//	Returns true if all three specified keys are pressed.
func (_ *NgUserIO) KeysPressedAll3(k1, k2, k3 int) bool {
	return UserIO.KeyPressed(k1) && UserIO.KeyPressed(k2) && UserIO.KeyPressed(k3)
}

//	Returns true if any of the two specified keys is pressed.
func (_ *NgUserIO) KeysPressedAny2(k1, k2 int) bool {
	return UserIO.KeyPressed(k1) || UserIO.KeyPressed(k2)
}

//	Returns true if any of the three specified keys is pressed.
func (_ *NgUserIO) KeysPressedAny3(k1, k2, k3 int) bool {
	return UserIO.KeyPressed(k1) || UserIO.KeyPressed(k2) || UserIO.KeyPressed(k3)
}

//	Returns true if the specified key has been "toggled", ie. its pressed-state changed within the last me.KeyToggleMinDelay seconds.
func (_ *NgUserIO) KeyToggled(key int) bool {
	if UserIO.togglePress = UserIO.KeyPressed(key); UserIO.togglePress && ((Loop.Tick.Now - UserIO.lastToggles[key]) > UserIO.KeyToggleMinDelay) {
		UserIO.lastToggles[key] = Loop.Tick.Now
		return true
	}
	return false
}
