package core

import (
	ugl "github.com/go3d/go-opengl/util"
)

type FxImage interface {
	Load() error
	Loaded() bool
	GpuSync() error
}

type FxImageStorage struct {
	UintRev bool
}

type FxImageBase struct {
	PreProcess struct {
		FlipY    bool
		ToLinear bool
		ToBgra   bool
	}
	Storage FxImageStorage

	glTex    *ugl.TextureBase
	glSynced bool
}

func (me *FxImageBase) dispose() {
	me.GpuDelete()
}

func (me *FxImageBase) init(glTex *ugl.TextureBase) {
	me.glTex = glTex
	me.PreProcess.ToLinear, me.PreProcess.FlipY, me.PreProcess.ToBgra = true, true, true
	me.Storage = Core.Options.Textures.FxImageStorage
}

func (me *FxImageBase) gpuSync(tex ugl.Texture) (err error) {
	if err = tex.Recreate(); err == nil {
		me.glSynced = true
	}
	return
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
