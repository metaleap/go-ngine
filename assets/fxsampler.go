package assets

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

const (
	//	No MIP-mapped minification.
	FX_SAMPLER_FILTER_NONE = 20
	//	Bilinear filtering.
	FX_SAMPLER_FILTER_NEAREST = 0x2600
	//	Trilinear filtering.
	FX_SAMPLER_FILTER_LINEAR = 0x2601
	//	Compensates for distortion caused by the difference in angle between the polygon and the plane of the screen.
	FX_SAMPLER_FILTER_ANISOTROPIC = 21

	//	Declares a one-dimensional texture sampler.
	FX_SAMPLER_TYPE_1D = 0x8B5D
	//	Declares a two-dimensional texture sampler.
	FX_SAMPLER_TYPE_2D = 0x8B5E
	//	Declares a three-dimensional texture sampler.
	FX_SAMPLER_TYPE_3D = 0x8B5F
	//	Declares a texture sampler for cube maps.
	FX_SAMPLER_TYPE_CUBE = 0x8B60
	//	Declares a texture sampler for depth maps.
	FX_SAMPLER_TYPE_DEPTH = 30
	//	Declares a rectangular texture sampler.
	FX_SAMPLER_TYPE_RECT = 31

	//	Ignores the integer part of texture coordinates, using only the fractional part. Tiles the texture at every integer junction. For example, for u values between 0 and 3, the texture is repeated three times; no mirroring is performed.
	FX_SAMPLER_WRAP_WRAP = 0x2901
	//	First mirrors the texture coordinate. The mirrored coordinate is then clamped as described for FX_SAMPLER_WRAP_CLAMP. Flips the texture at every integer junction. For u values between 0 and 1, for example, the texture is addressed normally; between 1 and 2, the texture is flipped (mirrored); between 2 and 3, the texture is normal again; and so on.
	FX_SAMPLER_WRAP_MIRROR = 0x8370
	//	Clamps texture coordinates at all MIPmap levels such that the texture filter never samples a border texel.
	FX_SAMPLER_WRAP_CLAMP = 0x812F
	//	Clamps texture coordinates at all MIPmaps such that the texture filter always samples border texels for fragments whose corresponding texture coordinate is sufficiently far outside the range [0, 1]. Much like FX_SAMPLER_WRAP_CLAMP, except texture coordinates outside the range [0.0, 1.0] are set to the border color.
	FX_SAMPLER_WRAP_BORDER = 0x812D
	//	Takes the absolute value of the texture coordinate (thus, mirroring around 0), and then clamps to the maximum value.
	FX_SAMPLER_WRAP_MIRROR_ONCE = 40
)

var (
	//	Default texture minification, magnification and MIP-mapping.
	DefaultFxSamplerFiltering = &FxSamplerFiltering{
		FilterMag:     FX_SAMPLER_FILTER_LINEAR,
		FilterMin:     FX_SAMPLER_FILTER_LINEAR,
		FilterMip:     FX_SAMPLER_FILTER_LINEAR,
		MaxAnisotropy: 1,
	}
	//	Default texture repeating and clamping.
	DefaultFxSamplerWrapping = &FxSamplerWrapping{
		BorderColor: ugfx.Rgba32{R: 0, G: 0, B: 0, A: 1},
		WrapS:       FX_SAMPLER_WRAP_WRAP,
		WrapT:       FX_SAMPLER_WRAP_WRAP,
		WrapP:       FX_SAMPLER_WRAP_WRAP,
	}
)

//	Declares a texture sampler.
type FxSampler struct {
	//	Extras
	HasExtras
	//	Filtering and Wrapping
	FxSamplerStates
	//	If set, instantiates a default image from which the sampler is to consume.
	Image *FxImageInst
	//	Indicates the type of this texture sampler. Must be one of the FX_SAMPLER_TYPE_* enumerated constants.
	Type int
}

//	Constructor
func NewFxSampler() (me *FxSampler) {
	me = &FxSampler{}
	me.FxSamplerStates.Filtering = DefaultFxSamplerFiltering
	me.FxSamplerStates.Wrapping = DefaultFxSamplerWrapping
	return
}

//	Instantiates an image targeted for samplers.
type FxSamplerImage struct {
	//	Sid, Name, Extras, DefRef
	FxImageInst
}

//	Controls texture minification, magnification and MIP-mapping.
type FxSamplerFiltering struct {
	//	Magnification filter. Must be one of the FX_SAMPLER_FILTER_* enumerated constants.
	FilterMag int
	//	Minification filter. Must be one of the FX_SAMPLER_FILTER_* enumerated constants.
	FilterMin int
	//	Mip-mapping filter. Must be one of the FX_SAMPLER_FILTER_* enumerated constants.
	FilterMip int
	//	The number of samples that can be used durring anisotropic filtering.
	MaxAnisotropy uint32
	//	The maximum number of progressive levels that the sampler will evaluate.
	MipMaxLevel uint8
	//	The minimum progressive levels to begin to evaluate.
	MipMinLevel uint8
	//	Biases the gamma (level of detail parameter) that is used by the sampler to evaluate the MIPmap chain.
	MipBias float64
}

//	Allows users to modify an effect’s sampler state from a material.
type FxSamplerStates struct {
	//	Extras
	HasExtras
	//	Controls texture minification, magnification and MIP-mapping.
	Filtering *FxSamplerFiltering
	//	Controls texture repeating and clamping.
	Wrapping *FxSamplerWrapping
}

//	Constructor
func NewFxSamplerStates() (me *FxSamplerStates) {
	me = &FxSamplerStates{Filtering: DefaultFxSamplerFiltering, Wrapping: DefaultFxSamplerWrapping}
	return
}

//	Controls texture repeating and clamping.
type FxSamplerWrapping struct {
	//	When reading past the edge of the texture address space based on the wrap modes involving clamps, this color takes over.
	BorderColor ugfx.Rgba32
	//	Controls texture repeating and clamping of the S coordinate. Must be one of the FX_SAMPLER_WRAP_* enumerated constants.
	WrapS int
	//	Controls texture repeating and clamping of the T coordinate. Must be one of the FX_SAMPLER_WRAP_* enumerated constants.
	WrapT int
	//	Controls texture repeating and clamping of the P coordinate. Must be one of the FX_SAMPLER_WRAP_* enumerated constants.
	WrapP int
}
