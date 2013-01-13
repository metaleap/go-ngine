package core

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

type FxColorOrTexture struct {
	Color   *ugfx.Rgba64
	Texture struct {
		Image2ID string
		Sampler  *FxSampler
	}
}

func NewFxColor(col *ugfx.Rgba64) (me *FxColorOrTexture) {
	me = &FxColorOrTexture{Color: col}
	return
}

func NewFxTexture(image2ID string, sampler *FxSampler) (me *FxColorOrTexture) {
	me = &FxColorOrTexture{}
	if sampler == nil {
		sampler = FxSamplerDefault
	}
	me.Texture.Image2ID, me.Texture.Sampler = image2ID, sampler
	return
}

type FxEffect struct {
	Ambient *FxColorOrTexture
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

//#end-gt
