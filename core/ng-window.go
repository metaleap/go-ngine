package core

import (
	glfw "github.com/go-gl/glfw"
)

//	Only used for core.UserIO.Window after initialization, BUT: core.UserIO.Window gets
//	initialized from a copy of your NewWindowOptions() that you pass to your NewEngineOptions().
type WindowOptions struct {
	//	Minimum delay, in seconds, to wait after the last window-resize event received from
	//	the OS before notifying the rendering runtime of the new window dimensions.
	//	Defaults to 0.15.
	ResizeMinDelay float64

	fullscreen, isCreated bool
	width, height, swap   int
	title                 string
	lastResize            float64
}

func (me *WindowOptions) Created() bool {
	return UserIO.isGlfwInit && me.isCreated
}

func (me *WindowOptions) Fullscreen() bool {
	return me.fullscreen
}

//	Returns the height of the window in pixels.
func (me *WindowOptions) Height() int {
	return me.height
}

func (me *WindowOptions) SetSize(width, height int) {
	if me.width, me.height = width, height; me.Created() {
		glfw.SetWindowSize(width, height)
	}
}

func (me *WindowOptions) SetSwapInterval(newSwap int) {
	me.swap = newSwap
	if me.Created() {
		glfw.SetSwapInterval(me.swap)
	}
}

//	Sets the window title to newTitle.
func (me *WindowOptions) SetTitle(newTitle string) {
	me.title = newTitle
	if me.Created() {
		glfw.SetWindowTitle(newTitle)
	}
}

func (me *WindowOptions) SwapInterval() int {
	return me.swap
}

func (me *WindowOptions) Title() string {
	return me.title
}

//	Returns the width of the window in pixels.
func (me *WindowOptions) Width() int {
	return me.width
}

//	just a GLFW event callback without creating a closure
func glfwOnWindowClose() (doit int) {
	if !Loop.Running {
		doit = 1
	}
	//	If looping, return 0 to cancel the user's onWindowClose -- we close manually.
	//	If we returned 1 to confirm the Close, the loop would continue
	//	doing GL stuff against the destroyed window and record a GL error.
	Loop.Running = false
	return
}

//	just a GLFW event callback without creating a closure
func glfwOnWindowResize(width, height int) {
	UserIO.Window.width, UserIO.Window.height = width, height
	UserIO.Window.lastResize = Loop.Tick.Now
}
