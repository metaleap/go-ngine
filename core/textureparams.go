package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
	util "github.com/go3d/go-util"

	nglcore "github.com/go3d/go-ngine/core/glcore"
)

type tTextureParams struct {
	TextureProviders *tTextureProviders

	aniso float64
	filter bool

	glAniso gl.Float
	glFilterMag, glFilterMin gl.Int
}

func NewTextureParams (filter bool, filterAnisotropy float64) *tTextureParams {
	var tp = &tTextureParams {}
	tp.TextureProviders = textureProviders
	tp.filter, tp.aniso = filter, filterAnisotropy
	return tp
}

func (me *tTextureParams) apply (tex *TTexture) {
	if me.glAniso > 0 { gl.TexParameterf(gl.TEXTURE_2D, glutil.TEXTURE_MAX_ANISOTROPY_EXT, me.glAniso) }
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, me.glFilterMag)
	if tex.noMipMap {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, glutil.IfI(me.filter, gl.LINEAR, gl.NEAREST))
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, me.glFilterMin)
	}
}

func (me *tTextureParams) Filter () bool {
	return me.filter
}

func (me *tTextureParams) FilterAnisotropy () float64 {
	return me.aniso
}

func (me *tTextureParams) gpuSync () {
	for _, tex := range Core.Textures {
		if (tex.Params == me) && tex.GpuSynced() {
			gl.BindTexture(gl.TEXTURE_2D, tex.glTex)
			me.apply(tex)
		} else {
		}
	}
}

func (me *tTextureParams) Set (filter bool, filterAnisotropy float64) {
	var glAniso gl.Float
	var glMin, glMag gl.Int
	var changed bool
	if glutil.MaxTextureAnisotropy() > 1 {
		glAniso = glutil.ClampF(gl.Float(filterAnisotropy), 1, glutil.MaxTextureAnisotropy())
	} else {
		glAniso, filterAnisotropy = 0, 0
	}
	glMag = glutil.IfI(filter, gl.LINEAR, gl.NEAREST)
	glMin = glutil.IfI(filter, gl.LINEAR_MIPMAP_LINEAR, gl.NEAREST_MIPMAP_NEAREST)
	me.filter, me.aniso = filter, filterAnisotropy
	if glAniso != me.glAniso { changed, me.glAniso = true, glAniso }
	if glMin != me.glFilterMin { changed, me.glFilterMin = true, glMin }
	if glMag != me.glFilterMag { changed, me.glFilterMag = true, glMag }
	if changed && nglcore.IsGlInit { me.gpuSync() }
}

func (me *tTextureParams) setAgain () {
	me.Set(me.filter, me.aniso)
}

func (me *tTextureParams) SetFilter (filter bool) {
	me.Set(filter, me.FilterAnisotropy())
}

func (me *tTextureParams) SetFilterAnisotropy (filterAnisotropy float64) {
	me.Set(me.Filter(), filterAnisotropy)
}

func (me *tTextureParams) ToggleFilter () {
	me.SetFilter(!me.Filter())
}

func (me *tTextureParams) ToggleFilterAnisotropy () {
	me.Set(me.Filter(), util.Ifd(me.aniso == float64(glutil.MaxTextureAnisotropy()), 1, me.aniso + 1))
}
