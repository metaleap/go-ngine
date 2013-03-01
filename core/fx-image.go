package core

import (
	"io"

	ugl "github.com/go3d/go-opengl/util"
)

type fxImage interface {
	Load() error
	Loaded() bool
	GpuSync() error
}

type FxImagePreprocess struct {
	FlipY    bool
	ToLinear bool
	ToBgra   bool
}

type FxImageStorage struct {
	DiskCache struct {
		Enabled      bool
		Compressor   func(w io.WriteCloser) io.WriteCloser
		Decompressor func(r io.ReadCloser) io.ReadCloser
	}
	Gpu struct {
		Bgra    bool
		UintRev bool
	}
}

type FxImageBase struct {
	ID         int
	Preprocess FxImagePreprocess
	Storage    FxImageStorage

	glSynced bool
}

func (me *FxImageBase) init() {
	me.Storage = Options.Textures.Storage
	me.Preprocess.ToLinear, me.Preprocess.FlipY, me.Preprocess.ToBgra = true, true, me.Storage.Gpu.Bgra
}

func (me *FxImageBase) gpuDelete(glTex ugl.Texture) {
	glTex.Dispose()
	me.glSynced = false
}

func (me *FxImageBase) gpuSync(tex ugl.Texture) (err error) {
	if err = tex.Recreate(); err == nil {
		me.glSynced = true
	}
	return
}

func (me *FxImageBase) GpuSynced() bool {
	return me.glSynced
}

func (me *FxImageBase) needPreproc() bool {
	return me.Preprocess.FlipY || me.Preprocess.ToBgra || me.Preprocess.ToLinear
}

func (me *FxImageBase) noAutoMips(glTex *ugl.TextureBase) {
	glTex.MipMap.AutoGen = false
}
