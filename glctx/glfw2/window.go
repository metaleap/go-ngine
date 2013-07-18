package glctx_glfw2

import (
	glfw "github.com/go-gl/glfw"
)

type window struct {
}

func newWindow() (me *window) {
	me = &window{}
	return
}

func (me *window) CallbackWindowClose(f func()) {
	glfw.SetWindowCloseCallback(func() int {
		f()
		return 0
	})
}

func (me *window) CallbackWindowSize(f func(int, int)) {
	glfw.SetWindowSizeCallback(f)
}

func (me *window) Close() {
	glfw.CloseWindow()
}

func (me *window) InputMode(flag, value int) {
	if value == 0 {
		glfw.Disable(flag)
	} else {
		glfw.Enable(flag)
	}
}

func (me *window) Key(key int) int {
	return glfw.Key(key)
}

func (me *window) SetSize(width, height int) {
	glfw.SetWindowSize(width, height)
}

func (me *window) SetTitle(title string) {
	glfw.SetWindowTitle(title)
}

func (me *window) ShouldClose() bool {
	return glfw.WindowParam(glfw.Opened) != 1
}

func (me *window) Size() (width, height int) {
	return glfw.WindowSize()
}

func (me *window) SwapBuffers() {
	glfw.SwapBuffers()
}
