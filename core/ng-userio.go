package core

import (
	"runtime"

	glfw "github.com/go-gl/glfw"
	ugl "github.com/go3d/go-opengl/util"
	ugo "github.com/metaleap/go-util"
)

var (
	//	Your gateway for end-user input (key, mouse etc.) and "output" (window management, not the graphics themselves).
	UserIO EngineUserIO
)

//	Consider EngineUserIO a "Singleton" type, only valid use is the core.UserIO global variable.
//	Your gateway for end-user input (key, mouse etc.) and "output" (window management, not the graphics themselves).
type EngineUserIO struct {
	//	Minimum delay for EngineUserIO.KeyToggled() method, in seconds. Defaults to 0.25.
	KeyToggleMinDelay float64

	//	Minimum delay, in seconds, to wait after the last window-resize event received from
	//	the OS before notifying the rendering runtime of the new window dimensions.
	//	Defaults to 0.25.
	WinResizeMinDelay float64

	isGlfwInit, isGlfwWindow, togglePress bool
	keyWhich                              int
	lastWinResize                         float64
	lastToggles                           map[int]float64
}

func (_ *EngineUserIO) dispose() {
	if UserIO.isGlfwWindow {
		UserIO.isGlfwWindow = false
		glfw.CloseWindow()
	}
	if UserIO.isGlfwInit {
		UserIO.isGlfwInit = false
		glfw.Terminate()
	}
}

func (_ *EngineUserIO) init(winTitle string, forceContextVersion float64) (err error) {
	opt := &Core.Options
	UserIO.KeyToggleMinDelay, UserIO.WinResizeMinDelay, UserIO.lastToggles = 0.25, 0.25, map[int]float64{}
	if !UserIO.isGlfwInit {
		if err = glfw.Init(); err == nil {
			UserIO.isGlfwInit = true
		}
	}
	if UserIO.isGlfwInit && !UserIO.isGlfwWindow {
		glfw.OpenWindowHint(glfw.FsaaSamples, 0) // AA will be a pluggable post-processing shader
		if forceContextVersion > 0 {
			major, minor := ugl.VersionMajorMinor(forceContextVersion)
			glfw.OpenWindowHint(glfw.OpenGLVersionMajor, major)
			glfw.OpenWindowHint(glfw.OpenGLVersionMinor, minor)
			glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
			if runtime.GOOS == "darwin" {
				glfw.OpenWindowHint(glfw.OpenGLForwardCompat, 1)
			}
		}
		winInit := &opt.Initialization.Window
		if err = glfw.OpenWindow(opt.winWidth, opt.winHeight, winInit.Rbits, winInit.Gbits, winInit.Bbits, winInit.Abits, winInit.DepthBits, winInit.StencilBits, ugo.Ifi(opt.winFullScreen, glfw.Fullscreen, glfw.Windowed)); err == nil {
			opt.winWidth, opt.winHeight = glfw.WindowSize()
			UserIO.isGlfwWindow = true
		}
	}
	if UserIO.isGlfwWindow {
		UserIO.SetWinTitle(winTitle)
		glfw.SetSwapInterval(opt.winSwapInterval)
		glfw.SetWindowCloseCallback(glfwOnWindowClose)
		glfw.SetWindowSizeCallback(glfwOnWindowResize)
		// glfw.Disable(glfw.MouseCursor)
		glfw.Disable(glfw.AutoPollEvents)
		glfw.Enable(glfw.StickyKeys)
	}
	return
}

//	just a GLFW event callback without creating a closure
func glfwOnWindowClose() int {
	Loop.Looping = false
	//	Return 0 to cancel the user's onWindowClose -- we close manually.
	//	If we returned 1 to confirm the Close, the loop would continue
	//	doing GL stuff against the destroyed window and record a GL error.
	return 0
}

//	just a GLFW event callback without creating a closure
func glfwOnWindowResize(width, height int) {
	Core.Options.winWidth, Core.Options.winHeight = width, height
	UserIO.lastWinResize = Loop.Tick.Now
}

//	Returns ifTrue if the specified key is pressed, otherwise returns ifFalse.
func (_ *EngineUserIO) IifKeyF(key int, ifTrue, ifFalse float64) float64 {
	if UserIO.KeyPressed(key) {
		return ifTrue
	}
	return ifFalse
}

//	Returns true if the specified key is pressed.
func (_ *EngineUserIO) KeyPressed(key int) bool {
	return glfw.Key(key) == glfw.KeyPress
}

//	Returns the first in keys that is pressed.
func (_ *EngineUserIO) KeyPressedWhich(keys ...int) int {
	for _, UserIO.keyWhich = range keys {
		if UserIO.KeyPressed(UserIO.keyWhich) {
			return UserIO.keyWhich
		}
	}
	return 0
}

//	Returns true if both specified keys are pressed.
func (_ *EngineUserIO) KeysPressedAll2(k1, k2 int) bool {
	return UserIO.KeyPressed(k1) && UserIO.KeyPressed(k2)
}

//	Returns true if all three specified keys are pressed.
func (_ *EngineUserIO) KeysPressedAll3(k1, k2, k3 int) bool {
	return UserIO.KeyPressed(k1) && UserIO.KeyPressed(k2) && UserIO.KeyPressed(k3)
}

//	Returns true if any of the two specified keys is pressed.
func (_ *EngineUserIO) KeysPressedAny2(k1, k2 int) bool {
	return UserIO.KeyPressed(k1) || UserIO.KeyPressed(k2)
}

//	Returns true if any of the three specified keys is pressed.
func (_ *EngineUserIO) KeysPressedAny3(k1, k2, k3 int) bool {
	return UserIO.KeyPressed(k1) || UserIO.KeyPressed(k2) || UserIO.KeyPressed(k3)
}

//	Returns true if the specified key has been "toggled", ie. its pressed-state changed within the last me.KeyToggleMinDelay seconds.
func (_ *EngineUserIO) KeyToggled(key int) bool {
	if UserIO.togglePress = UserIO.KeyPressed(key); UserIO.togglePress && ((Loop.Tick.Now - UserIO.lastToggles[key]) > UserIO.KeyToggleMinDelay) {
		UserIO.lastToggles[key] = Loop.Tick.Now
		return true
	}
	return false
}

//	Sets the window title to newTitle.
func (_ *EngineUserIO) SetWinTitle(newTitle string) {
	glfw.SetWindowTitle(newTitle)
}

//	Returns the height of the window in pixels.
func (_ *EngineUserIO) WinHeight() int {
	return Core.Options.winHeight
}

//	Returns the width of the window in pixels.
func (_ *EngineUserIO) WinWidth() int {
	return Core.Options.winWidth
}
