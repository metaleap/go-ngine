package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
	util "github.com/metaleap/go-util"
)

type textureParams struct {
	aniso  float64
	filter bool

	glAniso                  gl.Float
	glFilterMag, glFilterMin gl.Int
}

func newTextureParams(filter bool, filterAnisotropy float64) (me *textureParams) {
	me = &textureParams{filter: filter, aniso: filterAnisotropy}
	return
}

func (me *textureParams) apply(tex *Texture) {
	if me.glAniso > 0 {
		gl.TexParameterf(gl.TEXTURE_2D, ugl.TEXTURE_MAX_ANISOTROPY_EXT, me.glAniso)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, me.glFilterMag)
	if tex.noMipMap {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, ugl.Ifi(me.filter, gl.LINEAR, gl.NEAREST))
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, me.glFilterMin)
	}
}

func (me *textureParams) Filter() bool {
	return me.filter
}

func (me *textureParams) FilterAnisotropy() float64 {
	return me.aniso
}

func (me *textureParams) gpuSync() {
	for _, tex := range Core.Libs.Textures {
		if (tex.Params == me) && tex.GpuSynced() {
			gl.BindTexture(gl.TEXTURE_2D, tex.glTex)
			me.apply(tex)
		}
	}
}

func (me *textureParams) Set(filter bool, filterAnisotropy float64) {
	var (
		glAniso gl.Float
		changed bool
	)
	if ugl.MaxTextureAnisotropy >= 1 {
		glAniso = ugl.Clamp(gl.Float(filterAnisotropy), 1, ugl.MaxTextureAnisotropy)
	} else {
		glAniso, filterAnisotropy = 0, 0
	}
	glMag, glMin := ugl.Ifi(filter, gl.LINEAR, gl.NEAREST), ugl.Ifi(filter, gl.LINEAR_MIPMAP_LINEAR, gl.NEAREST_MIPMAP_NEAREST)
	me.filter, me.aniso = filter, filterAnisotropy
	if glAniso != me.glAniso {
		changed, me.glAniso = true, glAniso
	}
	if glMin != me.glFilterMin {
		changed, me.glFilterMin = true, glMin
	}
	if glMag != me.glFilterMag {
		changed, me.glFilterMag = true, glMag
	}
	if changed && glIsInit {
		me.gpuSync()
	}
}

func (me *textureParams) setAgain() {
	me.Set(me.filter, me.aniso)
}

func (me *textureParams) SetFilter(filter bool) {
	me.Set(filter, me.FilterAnisotropy())
}

func (me *textureParams) SetFilterAnisotropy(filterAnisotropy float64) {
	me.Set(me.Filter(), filterAnisotropy)
}

func (me *textureParams) ToggleFilter() {
	me.SetFilter(!me.Filter())
}

func (me *textureParams) ToggleFilterAnisotropy() {
	me.Set(me.Filter(), util.Ifd(me.aniso == float64(ugl.MaxTextureAnisotropy), 1, me.aniso+1))
}
