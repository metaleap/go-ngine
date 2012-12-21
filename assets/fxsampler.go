package assets

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

//	Categorizes the kinds of filtering used in FxSamplerFiltering.
type FxFilterKind int

const (
	//	Bilinear filtering.
	FxFilterKindNearest FxFilterKind = 0x2600
	//	Trilinear filtering.
	FxFilterKindLinear FxFilterKind = 0x2601
	//	Compensates for distortion caused by the difference in angle between a polygon and the view plane.
	FxFilterKindAnisotropic FxFilterKind = 21
	//	No MIP-mapped minification.
	FxFilterKindMipNone FxFilterKind = 22
)

//	Categorizes the kind of an FxSampler.
type FxSamplerKind int

const (
	//	Declares a one-dimensional texture sampler.
	FxSamplerKind1D FxSamplerKind = 0x8B5D
	//	Declares a two-dimensional texture sampler.
	FxSamplerKind2D FxSamplerKind = 0x8B5E
	//	Declares a three-dimensional texture sampler.
	FxSamplerKind3D FxSamplerKind = 0x8B5F
	//	Declares a texture sampler for cube maps.
	FxSamplerKindCube FxSamplerKind = 0x8B60
	//	Declares a texture sampler for depth maps.
	FxSamplerKindDepth FxSamplerKind = 31
	//	Declares a rectangular texture sampler.
	FxSamplerKindRect FxSamplerKind = 32
)

//	Categorizes the kinds of wrapping used in FxSamplerWrapping.
type FxWrapKind int

const (
	//	Ignores the integer part of texture coordinates, using only the fractional part and tiling the
	//	texture at every integer junction. For example, for u values between 0 and 3, the texture is
	//	repeated three times; no mirroring is performed.
	FxWrapKindWrap FxWrapKind = 0x2901
	//	First mirrors the texture coordinate. The mirrored coordinate is then clamped as described for
	//	FxWrapKindClamp. Flips the texture at every integer junction. For u values between 0 and 1,
	//	for example, the texture is addressed normally; between 1 and 2, the texture is flipped (mirrored);
	//	between 2 and 3, the texture is normal again; and so on.
	FxWrapKindMirror FxWrapKind = 0x8370
	//	Clamps texture coordinates at all MIPmap levels such that
	//	the texture filter never samples a border texel.
	FxWrapKindClamp FxWrapKind = 0x812F
	//	Clamps texture coordinates at all MIPmaps such that the texture filter always samples border
	//	texels for fragments whose corresponding texture coordinate is sufficiently far outside
	//	the range [0, 1]. Much like FxWrapKindClamp, except texture coordinates outside
	//	the range [0.0, 1.0] are set to the border color.
	FxWrapKindBorder FxWrapKind = 0x812D
	//	Takes the absolute value of the texture coordinate (thus, mirroring around 0),
	//	and then clamps to the maximum value.
	FxWrapKindMirrorOnce FxWrapKind = 41
)

var (
	//	Default texture minification, magnification and MIP-mapping.
	DefaultFxSamplerFiltering = &FxSamplerFiltering{
		FilterMag:     FxFilterKindLinear,
		FilterMin:     FxFilterKindLinear,
		FilterMip:     FxFilterKindLinear,
		MaxAnisotropy: 1,
	}
	//	Default texture repeating and clamping.
	DefaultFxSamplerWrapping = &FxSamplerWrapping{
		BorderColor: ugfx.Rgba32{R: 0, G: 0, B: 0, A: 1},
		WrapS:       FxWrapKindWrap,
		WrapT:       FxWrapKindWrap,
		WrapP:       FxWrapKindWrap,
	}
)

//	Declares a texture sampler.
type FxSampler struct {
	//	Extras
	HasExtras
	//	Filtering, Wrapping
	FxSamplerStates
	//	If set, instantiates a default image from which the sampler is to consume.
	Image *FxImageInst
	//	Indicates the kind of this texture sampler.
	//	Must be one of the FxSamplerKind* enumerated constants.
	Kind FxSamplerKind
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
	//	Magnification filter. Must be one of the FxFilterKind* enumerated constants.
	FilterMag FxFilterKind
	//	Minification filter. Must be one of the FxFilterKind* enumerated constants.
	FilterMin FxFilterKind
	//	Mip-mapping filter. Must be one of the FxFilterKind* enumerated constants.
	FilterMip FxFilterKind
	//	The number of samples that can be used durring anisotropic filtering.
	MaxAnisotropy uint32
	//	The maximum number of progressive levels that the sampler will evaluate.
	MipMaxLevel uint8
	//	The minimum progressive levels to begin to evaluate.
	MipMinLevel uint8
	//	Biases the gamma (level of detail parameter) used by the sampler to evaluate the MIPmap chain.
	MipBias float64
}

//	Allows users to modify an effect's sampler state from a material.
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
	//	When reading past the edge of the texture address space
	//	based on the wrap modes involving clamps, this color takes over.
	BorderColor ugfx.Rgba32
	//	Controls texture repeating and clamping of the S coordinate.
	//	Must be one of the FxWrapKind* enumerated constants.
	WrapS FxWrapKind
	//	Controls texture repeating and clamping of the T coordinate.
	//	Must be one of the FxWrapKind* enumerated constants.
	WrapT FxWrapKind
	//	Controls texture repeating and clamping of the P coordinate.
	//	Must be one of the FxWrapKind* enumerated constants.
	WrapP FxWrapKind
}
