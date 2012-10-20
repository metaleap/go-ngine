package core

import (
	// gl "github.com/chsc/gogl/gl42"

	// glutil "github.com/go3d/go-util/gl"
)

type TTexture struct {
	loaded, synced bool
}

func (me *TTexture) Delete () {
	
}

func (me *TTexture) Dispose () {

}

func (me *TTexture) Loaded () bool {
	return me.loaded
}

func (me *TTexture) Unload () {
	me.loaded = false
}
