# glctx
--
    import "github.com/go3d/go-ngine/glctx"


## Usage

#### type BufferBits

```go
type BufferBits struct {
	Color struct {
		R, G, B, A int
	}
	Depth   int
	Stencil int
}
```


#### type CtxProfile

```go
type CtxProfile struct {
	CoreProfile, ForwardCompatibility bool
	Version                           struct {
		Major, Minor int
	}
}
```


#### type CtxProvider

```go
type CtxProvider interface {
	Hint(flag, value int)
	Init() error
	Window(winInfo *WinProfile, bufSize *BufferBits, ctxInfo *CtxProfile) (Window, error)
	PollEvents()
	SetSwapInterval(int)
	SetTime(float64)
	Terminate()
	Time() float64
}
```


#### type WinProfile

```go
type WinProfile struct {
	Width, Height int
	FullScreen    bool
	Title         string
}
```


#### type Window

```go
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
```

--
**godocdown** http://github.com/robertkrimen/godocdown
