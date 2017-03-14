package glctx_glfw3

import (
	glfw "github.com/go-gl/glfw/v3.0/glfw"
)

type window struct {
	*glfw.Window
}

func newWindow(win *glfw.Window) (me *window) {
	me = &window{}
	me.Window = win
	me.Window.SetInputMode(glfw.StickyKeys, 0)
	me.Window.MakeContextCurrent()
	return
}

func (me *window) CallbackWindowClose(f func()) {
	me.Window.SetCloseCallback(func(_ *glfw.Window) {
		f()
		me.SetShouldClose(true)
	})
}

func (me *window) CallbackWindowSize(f func(int, int)) {
	me.Window.SetSizeCallback(func(_ *glfw.Window, w, h int) {
		f(w, h)
	})
}

func (me *window) Close() {
	me.Window.Destroy()
}

func (me *window) InputMode(flag, value int) {
	me.Window.SetInputMode(glfw.InputMode(flag), value)
}

func (me *window) Key(key int) int {
	return int(me.Window.GetKey(glfw.Key(key)))
}

func (me *window) Size() (width, height int) {
	return me.Window.GetSize()
}
