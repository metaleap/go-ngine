package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

type FxProc struct {
	Enabled bool

	weight    gl.Float
	procIndex int
	procID    string
	unifNames map[[2]string]string

	Color struct {
		Rgb ugl.GlVec3
	}

	Tex struct {
		ImageID int
		glUnitU gl.Uint
		glUnitI gl.Int
		glUnitE gl.Enum

		Sampler ugl.Sampler
	}
}

func (me *FxProc) init(procID string, procIndex int) {
	me.procID, me.procIndex, me.Enabled, me.weight = procID, -1, true, 1
	me.setProcIndex(procIndex)
	if me.IsTex() {
		me.Tex.ImageID = -1
		me.Tex.Sampler = Core.Render.Fx.Samplers.FullFilteringRepeat
	}
}

func (me *FxProc) qualifiers(inout string) (q string) {
	return
}

func (me *FxProc) setProcIndex(index int) {
	if index != me.procIndex {
		me.procIndex = index
		me.unifNames = map[[2]string]string{}
		if me.IsTex() {
			me.Tex.glUnitI = gl.Int(index)
			me.Tex.glUnitU = gl.Uint(index)
			me.Tex.glUnitE = gl.Enum(gl.TEXTURE0 + index)
		}
	}
}

func (me *FxProc) unifName(t, n string) (un string) {
	k := [2]string{t, n}
	if un = me.unifNames[k]; len(un) == 0 {
		un = strf("uni_%s_%s%d_%s", t, me.procID, me.procIndex, n)
		me.unifNames[k] = un
	}
	return
}

func (me *FxProc) use() {
	thrRend.curProg.Uniform1f(me.unifName("float", "MixWeight"), me.weight)
	if me.IsColor() {
		thrRend.curProg.Uniform3fv(me.unifName("vec3", "Rgb"), 1, &me.Color.Rgb[0])
	}
	if me.IsTex() {
		me.Tex.Sampler.Bind(me.Tex.glUnitU)
		ugl.Cache.ActiveTexture(me.Tex.glUnitE)
		if me.IsTex2D() {
			thrRend.curProg.Uniform1i(me.unifName("sampler2D", "Img"), me.Tex.glUnitI)
			if Core.Libs.Images.Tex2D.IsOk(me.Tex.ImageID) {
				Core.Libs.Images.Tex2D[me.Tex.ImageID].glTex.Bind()
			}
		} else if me.IsTexCube() {
			thrRend.curProg.Uniform1i(me.unifName("samplerCube", "Img"), me.Tex.glUnitI)
			if Core.Libs.Images.TexCube.IsOk(me.Tex.ImageID) {
				Core.Libs.Images.TexCube[me.Tex.ImageID].glTex.Bind()
			}
		}
	}
}

func (me *FxProc) SetMixWeight(weight float64) {
	me.weight = gl.Float(weight)
}

func (me *FxProc) Toggle() {
	me.Enabled = !me.Enabled
}

func (me *FxProc) Color_SetRgb(rgb ...gl.Float) *FxProc {
	me.Color.Rgb.Set(rgb...)
	return me
}

func (me *FxProc) Tex_SetImageID(imageID int) *FxProc {
	me.Tex.ImageID = imageID
	return me
}

func (me *FxProc) IsTex() bool {
	return me.IsTex2D() || me.IsTexCube()
}
