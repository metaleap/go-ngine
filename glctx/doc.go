//	Provides GL context creation for go:ngine.
//
//	go:ngine `core` packages are themselves toolkit-agnostic: you pass a `glctx.CtxProvider` during
//	initialization that takes care of creating the requested windowed or full-screen GL context.
//
//	Sub-packages `glfw2` and `glfw3` provide ready-made `CtxProvider` implementations
//	for GLFW 2.x and GLFW 3.x respectively. Sub-package `sdl` is currently a
//	dummy / no-op shim waiting to be implemented some fine day to provide an SDL implementation.
//
//	If your GL context creation needs are more exotic than those (or even to re-use an
//	existing GL context), simply implement your own `CtxProvider`.
package glctx
