package glctx

type BufferBits struct {
	Color struct {
		R, G, B, A int
	}
	Depth   int
	Stencil int
}

type CtxProvider interface {
	Hint(flag, value int)
	Init() error
	Window(width, height int, title string, bufSize *BufferBits, ctxInfo *CtxProfile) (Window, error)
	PollEvents()
	SetSwapInterval(int)
	SetTime(float64)
	Terminate()
	Time() float64
}

type CtxProfile struct {
	CoreProfile, ForwardCompatibility bool
	Version                           struct {
		Major, Minor int
	}
}

type Window interface {
	CallbackWindowClose(func())
	CallbackWindowSize(func(int, int))
	Close()
	InputMode(flag, value int)
	Key(int) int
	SetSize(width, height int)
	SetTitle(string)
	ShouldClose() bool
	Size() (width, height int)
	SwapBuffers()
}
