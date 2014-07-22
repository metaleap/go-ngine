# glctx
--
    import "github.com/go3d/go-ngine/glctx"

Provides GL context creation for go:ngine.

go:ngine `core` packages are themselves toolkit-agnostic: you pass a
`glctx.CtxProvider` during initialization that takes care of creating the
requested windowed or full-screen GL context.

Sub-packages `glfw2` and `glfw3` provide ready-made `CtxProvider`
implementations for GLFW 2.x and GLFW 3.x respectively. Sub-package `sdl` is
currently a dummy / no-op shim waiting to be implemented some fine day to
provide an SDL implementation.

If your GL context creation needs are more exotic than those (or even to re-use
an existing GL context), simply implement your own `CtxProvider`.

## Usage

#### type BufferBits

```go
type BufferBits struct {
	//	Color buffer
	Color struct{ R, G, B, A int }

	//	Depth buffer
	Depth int

	//	Stencil buffer
	Stencil int
}
```

Passed to `CtxProvider.Window` method. Defines bit-depths for various GL context
buffers.

#### type CtxProfile

```go
type CtxProfile struct {
	//	If `true`, an OpenGL Core Profile context is created.
	//	If `false`, an OpenGL Compatibility Profile context is created.
	CoreProfile bool

	//	Whether to create a "forward-compatible" OpenGL context.
	ForwardCompat bool

	//	Required version of the OpenGL context to be created.
	Version struct{ Major, Minor int }
}
```

Passed to `CtxProvider.Window` method. Declares GL context creation
requirements.

#### type CtxProvider

```go
type CtxProvider interface {
	//	Arbitrary configuration flags or hints for GL context and window creation via `Window` method.
	Hint(flag, value int)

	//	Call this always before working with a `CtxProvider`.
	Init() error

	//	Creates both the GL context specified by `ctxInfo` and the associated
	//	window (or full-screen) specified by `winInfo`.
	Window(winInfo *WinProfile, bufSize *BufferBits, ctxInfo *CtxProfile) (Window, error)

	//	Call this always before checking for fresh user input in a `Window`.
	PollEvents()

	//	Typically sets V-sync interval.
	SetSwapInterval(int)

	//	Resets the "time in seconds since `Init`".
	SetTime(float64)

	//	Releases resources associated with this `CtxProvider` and invalidates it for further use.
	//	All `Window`s should be `Close`d before you `Terminate`.
	Terminate()

	//	Returns the "time in seconds since `Init`".
	Time() float64
}
```


#### type WinProfile

```go
type WinProfile struct {
	//	Window dimensions or full-screen resolution.
	Width, Height int

	//	If `false`, a window is created; if `true`, a full-screen monitor is backing the newly created GL context.
	FullScreen bool

	//	Window title
	Title string

	//	Number of samples (0, 2, 4, 8...)
	MultiSampling int
}
```

Passed to `CtxProvider.Window` method. Declares window creation parameters.

#### type Window

```go
type Window interface {
	//	Lets you specify a call-back handler called on `Window.Close`.
	CallbackWindowClose(func())

	//	Lets you specify a call-back handler called when the window is resized.
	CallbackWindowSize(func(int, int))

	//	Closes this `Window` and destroys the associated OpenGL context.
	Close()

	//	Sets the specified `flag` input mode to the specified `value`.
	InputMode(flag, value int)

	//	If the specified key is pressed, should return 1; else should return 0.
	//	Both the argument and the return value are completely implementation-specific however.
	Key(int) int

	//	Changes the dimensions of this window or the resolution of this full-screen monitor.
	SetSize(width, height int)

	//	Changes the title of the opened window.
	SetTitle(string)

	//	Should return `true` if the user performed a "window-closing interaction".
	ShouldClose() bool

	//	Returns the dimensions of this window or the resolution of this full-screen monitor.
	Size() (width, height int)

	// SwapBuffers swaps the back and front color buffers of the window.
	SwapBuffers()
}
```

Returned by `CtxProvider.Window` method.

--
**godocdown** http://github.com/robertkrimen/godocdown
