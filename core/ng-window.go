package core

import (
	glfw "github.com/go-gl/glfw"
)

type WindowOptions struct {
	//	Defaults to a function that returns true to allow closing the window.
	OnCloseRequested func() bool

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
	if UserIO.Window.OnCloseRequested() {
		if !Loop.Running {
			doit = 1
		}
		//	If looping, return 0 for now to "cancel the Close" but also stop loop to close cleanly.
		//	If we returned 1 here to "accept the Close", the loop would continue
		//	doing GL stuff against the destroyed window and record a GL error.
		Loop.Running = false
	}
	return
}

//	just a GLFW event callback without creating a closure
func glfwOnWindowResize(width, height int) {
	UserIO.Window.width, UserIO.Window.height = width, height
	UserIO.Window.lastResize = Loop.Tick.Now
}
