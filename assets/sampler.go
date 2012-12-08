package assets

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

const (
	SAMPLER_FILTER_NONE        = 0
	SAMPLER_FILTER_NEAREST     = 0x2600
	SAMPLER_FILTER_LINEAR      = 0x2601
	SAMPLER_FILTER_ANISOTROPIC = 1

	SAMPLER_WRAP_WRAP        = 0x2901
	SAMPLER_WRAP_MIRROR      = 0x8370
	SAMPLER_WRAP_CLAMP       = 0x812F
	SAMPLER_WRAP_BORDER      = 0x812D
	SAMPLER_WRAP_MIRROR_ONCE = 0
)

var (
	DefaultSamplerFiltering = &SamplerFiltering{
		Mag:           SAMPLER_FILTER_LINEAR,
		Min:           SAMPLER_FILTER_LINEAR,
		Mip:           SAMPLER_FILTER_LINEAR,
		MaxAnisotropy: 1,
	}
	DefaultSamplerWrapping = &SamplerWrapping{
		BorderColor: ugfx.Rgba32{R: 0, G: 0, B: 0, A: 1},
		WrapS:       SAMPLER_WRAP_WRAP,
		WrapT:       SAMPLER_WRAP_WRAP,
		WrapP:       SAMPLER_WRAP_WRAP,
	}
)

type SamplerFiltering struct {
	Mag, Min, Mip            int
	MaxAnisotropy            uint32
	MipMaxLevel, MipMinLevel uint8
	MipBias                  float64
}

type SamplerWrapping struct {
	BorderColor         ugfx.Rgba32
	WrapS, WrapT, WrapP int
}

type Sampler struct {
	Image     *ImageInst
	Filtering *SamplerFiltering
	Wrapping  *SamplerWrapping
}

func NewSampler() (me *Sampler) {
	me = &Sampler{Filtering: DefaultSamplerFiltering, Wrapping: DefaultSamplerWrapping}
	return
}
