package assetutil

import (
	nga "github.com/go3d/go-ngine/assets"
)

//	Creates and returns a new FxSampler of assets.FxSamplerKind2D, sampling the specified
//	image with the specified filtering (or assets.DefaultFxSamplerFiltering if nil) and the
//	specified wrapping (or assets.DefaultFxSamplerWrapping if nil).
func NewFxSampler2D(image *nga.FxImageInst, filtering *nga.FxSamplerFiltering, wrapping *nga.FxSamplerWrapping) (me *nga.FxSampler) {
	me = nga.NewFxSampler()
	if filtering != nil {
		me.Filtering = filtering
	}
	if wrapping != nil {
		me.Wrapping = wrapping
	}
	me.Kind = nga.FxSamplerKind2D
	me.Image = image
	return
}
