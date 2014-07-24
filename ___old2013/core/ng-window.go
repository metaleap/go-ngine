package core

import (
	ngctx "github.com/go3d/go-ngine/glctx"
)

type WindowOptions struct {
	//	Defaults to a function that returns true to allow closing the window.
	OnCloseRequested func() bool

	//	Minimum delay, in seconds, to wait after the last window-resize event received from
	//	the OS before notifying the rendering runtime of the new window dimensions.
	//	Defaults to 0.15.
	ResizeMinDelay float64

	//	Number of samples (0, 2, 4, 8...)
	MultiSampling int

	win                   ngctx.Window
	fullscreen, isCreated bool
	width, height, swap   int
	title                 string
	lastResize            float64
}

func (me *WindowOptions) Created() bool {
	return UserIO.isCtxInit && me.isCreated && me.win != nil
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
		me.win.SetSize(width, height)
	}
}

func (me *WindowOptions) SetSwapInterval(newSwap int) {
	me.swap = newSwap
	if me.Created() {
		UserIO.ctx.SetSwapInterval(me.swap)
	}
}

//	Sets the window title to newTitle.
func (me *WindowOptions) SetTitle(newTitle string) {
	me.title = newTitle
	if me.Created() {
		me.win.SetTitle(newTitle)
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

func glctxOnWindowClose() {
	if UserIO.Window.OnCloseRequested() {
		Loop.Running = false
	}
}

func glctxOnWindowResize(width, height int) {
	UserIO.Window.width, UserIO.Window.height = width, height
	UserIO.Window.lastResize = Loop.Tick.Now
}
