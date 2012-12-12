package assets

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

const (
	FX_SAMPLER_FILTER_NONE        = 0
	FX_SAMPLER_FILTER_NEAREST     = 0x2600
	FX_SAMPLER_FILTER_LINEAR      = 0x2601
	FX_SAMPLER_FILTER_ANISOTROPIC = 1

	FX_SAMPLER_TYPE_1D    = 0x8B5D
	FX_SAMPLER_TYPE_2D    = 0x8B5E
	FX_SAMPLER_TYPE_3D    = 0x8B5F
	FX_SAMPLER_TYPE_CUBE  = 0x8B60
	FX_SAMPLER_TYPE_DEPTH = 80
	FX_SAMPLER_TYPE_RECT  = 81

	FX_SAMPLER_WRAP_WRAP        = 0x2901
	FX_SAMPLER_WRAP_MIRROR      = 0x8370
	FX_SAMPLER_WRAP_CLAMP       = 0x812F
	FX_SAMPLER_WRAP_BORDER      = 0x812D
	FX_SAMPLER_WRAP_MIRROR_ONCE = 0
)

var (
	DefaultFxSamplerFiltering = &FxSamplerFiltering{
		FilterMag:     FX_SAMPLER_FILTER_LINEAR,
		FilterMin:     FX_SAMPLER_FILTER_LINEAR,
		FilterMip:     FX_SAMPLER_FILTER_LINEAR,
		MaxAnisotropy: 1,
	}
	DefaultFxSamplerWrapping = &FxSamplerWrapping{
		BorderColor: ugfx.Rgba32{R: 0, G: 0, B: 0, A: 1},
		WrapS:       FX_SAMPLER_WRAP_WRAP,
		WrapT:       FX_SAMPLER_WRAP_WRAP,
		WrapP:       FX_SAMPLER_WRAP_WRAP,
	}
)

type FxSampler struct {
	HasExtras
	Image     *FxImageInst
	Filtering *FxSamplerFiltering
	Wrapping  *FxSamplerWrapping
	Type      int
}

func NewFxSampler() (me *FxSampler) {
	me = &FxSampler{Filtering: DefaultFxSamplerFiltering, Wrapping: DefaultFxSamplerWrapping}
	return
}

type FxSamplerImage struct {
	FxImageInst
}

type FxSamplerFiltering struct {
	FilterMag, FilterMin, FilterMip int
	MaxAnisotropy                   uint32
	MipMaxLevel, MipMinLevel        uint8
	MipBias                         float64
}

type FxSamplerStates struct {
	FxSampler
}

func NewFxSamplerStates() (me *FxSamplerStates) {
	me = &FxSamplerStates{}
	me.FxSampler = *NewFxSampler()
	return
}

type FxSamplerWrapping struct {
	BorderColor         ugfx.Rgba32
	WrapS, WrapT, WrapP int
}
