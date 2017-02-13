//	Currently a dummy / no-op shim waiting to be implemented some fine day to provide an SDL implementation of `glctx.CtxProvider`.
package glctx_sdl

import (
	ngctx "github.com/metaleap/go-ngine/glctx"
)

//	Currently returns `nil`.
func New() ngctx.CtxProvider {
	return nil
}
