package core

import (
	ugl "github.com/go3d/go-opengl/util"
	ugfx "github.com/metaleap/go-util/gfx"
)

//	Declares a source of color information:
//	either pointing to an FxImage2D or a 64-bit RGBA color.
type FxColorOrTexture struct {
	//	If set, the color is sourced from this value.
	Color *ugfx.Rgba64

	//	If Color is nil, the color is sourced from a 2D texture image.
	Texture struct {
		//	The ID of the FxImage2D in Core.Libs.Images.I2D
		Image2ID string

		//	Describes filtering and wrapping when sampling the texture image.
		//	Defaults to Core.Rendering.Samplers.FullFilteringRepeat.
		Sampler *ugl.Sampler
	}
}

//	Creates and returns a new FxColorOrTexture that points to the specified color.
func NewFxColor(rgba ...float64) (me *FxColorOrTexture) {
	me = &FxColorOrTexture{Color: ugfx.NewRgba64(rgba...)}
	return
}

//	Creates and returns a new FxColorOrTexture that points to the specified texture image.
//	If sampler is nil, Core.Rendering.Samplers.FullFilteringRepeat is used for me.Texture.Sampler.
func NewFxTexture(image2ID string, sampler *ugl.Sampler) (me *FxColorOrTexture) {
	me = &FxColorOrTexture{}
	if sampler == nil {
		sampler = &Core.Rendering.Samplers.FullFilteringRepeat
	}
	me.Texture.Image2ID, me.Texture.Sampler = image2ID, sampler
	return
}
