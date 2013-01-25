package core

import (
	ugl "github.com/go3d/go-glutil"
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

type FxEffect struct {
	//	The source for diffuse color information. Required.
	Diffuse *FxColorOrTexture
}

func (me *FxEffect) dispose() {
}

func (me *FxEffect) init() {
}

//#begin-gt -gen-lib.gt T:FxEffect

//	Initializes and returns a new FxEffect with default parameters.
func NewFxEffect() (me *FxEffect) {
	me = &FxEffect{}
	me.init()
	return
}

//	A hash-table of FxEffects associated by IDs. Only for use in Core.Libs.
type LibFxEffects map[string]*FxEffect

//	Creates and initializes a new FxEffect with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibFxEffects) AddNew(id string) (obj *FxEffect) {
	obj = NewFxEffect()
	me[id] = obj
	return
}

func (me *LibFxEffects) ctor() {
	*me = LibFxEffects{}
}

func (me *LibFxEffects) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibFxEffects) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
