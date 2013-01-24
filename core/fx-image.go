package core

import (
	ugl "github.com/go3d/go-glutil"
)

type FxImageBase struct {
	OnAsyncDone      func()
	OnLoad           FxImageOnLoad
	AsyncNumAttempts int
	InitFrom         struct {
		RawData []byte
		RefUrl  string
	}

	glTex    *ugl.TextureBase
	glSynced bool
}

func (me *FxImageBase) dispose() {
	me.GpuDelete()
}

func (me *FxImageBase) init(glTex *ugl.TextureBase) {
	me.glTex = glTex
}

func (me *FxImageBase) gpuSync(tex ugl.Texture) {
	tex.Recreate()
	me.glSynced = true
}

func (me *FxImageBase) GpuDelete() {
	me.glTex.Dispose()
	me.glSynced = false
}

func (me *FxImageBase) GpuSynced() bool {
	return me.glSynced
}

func (me *FxImageBase) NoAutoMips() {
	me.glTex.MipMap.AutoGen = false
}

func (me *FxImageBase) onAsyncDone() {
	if me.OnAsyncDone != nil {
		me.OnAsyncDone()
	}
}

type FxImageOnLoad func(img interface{}, err error, async bool)
