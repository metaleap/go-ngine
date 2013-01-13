//	The *core* package provides go:ngine client-side core functionality such as rendering and user I/O.
package core

type asyncResource interface {
	onAsyncDone()
}

type ctorable interface {
	ctor()
}

type disposable interface {
	dispose()
}
