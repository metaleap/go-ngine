package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
	ugfx "github.com/metaleap/go-util/gfx"
)

var (
	FxSamplerFilteringHighMipped = NewFxSamplerFiltering(8, true, true, true)
	FxSamplerWrappingRepeat      = ugfx.NewSamplerWrapping(ugfx.WrapKindRepeat, nil)
)

type FxSampler struct {
	Filtering *FxSamplerFiltering
	Wrapping  *ugfx.SamplerWrapping
}

func NewFxSampler(filtering *FxSamplerFiltering, wrapping *ugfx.SamplerWrapping) (me *FxSampler) {
	me = &FxSampler{Filtering: filtering, Wrapping: wrapping}
	return
}

type FxSamplerFiltering struct {
	glAniso       gl.Float
	glMag, glMin  gl.Int
	aniso         float64
	mag, min, mip bool
}

func NewFxSamplerFiltering(aniso float64, mag, min, mip bool) (me *FxSamplerFiltering) {
	me = &FxSamplerFiltering{}
	me.SetAnisotropy(aniso)
	me.SetFilteringMag(mag)
	me.SetFilteringMin(min, mip)
	return
}

func (me *FxSamplerFiltering) Anisotropy() float64 {
	return me.aniso
}

func (me *FxSamplerFiltering) Filtering() (mag, min, mip bool) {
	mag, min, mip = me.mag, me.min, me.mip
	return
}

func (me *FxSamplerFiltering) SetAnisotropy(aniso float64) {
	if me.aniso != aniso {
		if me.aniso = aniso; me.aniso > 0 {
			me.glAniso = gl.Float(me.aniso)
		} else {
			me.glAniso = ugl.MaxTextureAnisotropy
		}
	}
}

func (me *FxSamplerFiltering) SetFiltering(mag, min, mip bool) {
	me.SetFilteringMag(mag)
	me.SetFilteringMin(min, mip)
}

func (me *FxSamplerFiltering) SetFilteringMag(filter bool) {
	if me.mag != filter {
		me.mag, me.glMag = filter, ugl.Ifi(filter, gl.LINEAR, gl.NEAREST)
	}
}

func (me *FxSamplerFiltering) SetFilteringMin(min, mip bool) {
	if (me.min != min) || (me.mip != mip) {
		if me.min, me.mip = min, mip; me.mip {
			me.glMin = ugl.Ifi(me.min, gl.LINEAR_MIPMAP_LINEAR, gl.NEAREST_MIPMAP_NEAREST)
		} else {
			me.glMin = ugl.Ifi(me.min, gl.LINEAR, gl.NEAREST)
		}
	}
}
