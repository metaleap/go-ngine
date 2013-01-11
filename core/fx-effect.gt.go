package core

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

type FxColorOrTexture struct {
	Color   *ugfx.Rgba64
	Texture struct {
		Image2dID string
		Sampler   *FxSampler
	}
}

type FxEffect struct {
	Diffuse *FxColorOrTexture
}

func NewFxEffect() (me *FxEffect) {
	return
}

func (me *FxEffect) dispose() {
}

//#begin-gt -gen-lib.gt T:FxEffect

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
