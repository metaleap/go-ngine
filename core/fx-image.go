package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

type FxImage struct {
	OnLoad   FxImageOnLoad
	InitFrom struct {
		RawData []byte
		RefUrl  string
	}

	glSynced, noAutoMips                                      bool
	glTex                                                     gl.Uint
	glPixPointer                                              gl.Pointer
	glTexWidth, glTexHeight, glTexDepth, glTexLevels          gl.Sizei
	glSizedInternalFormat, glPixelDataFormat, glPixelDataType gl.Enum
}

func (me *FxImage) dispose() {
	me.GpuDelete()
}

func (me *FxImage) init() {
	me.glPixPointer = gl.Pointer(nil)
}

func (me *FxImage) gpuSync(canMip bool, glTarget gl.Enum) {
	me.GpuDelete()
	gl.GenTextures(1, &me.glTex)
	gl.BindTexture(glTarget, me.glTex)
	defer gl.BindTexture(glTarget, 0)
	if ugl.IsGl42 {
		switch glTarget {
		case gl.TEXTURE_2D:
			gl.TexStorage2D(glTarget, me.glTexLevels, me.glSizedInternalFormat, me.glTexWidth, me.glTexHeight)
			gl.TexSubImage2D(glTarget, 0, 0, 0, me.glTexWidth, me.glTexHeight, me.glPixelDataFormat, me.glPixelDataType, me.glPixPointer)
		}
	} else {
		switch glTarget {
		case gl.TEXTURE_2D:
			gl.TexImage2D(glTarget, 0, gl.Int(me.glSizedInternalFormat), me.glTexWidth, me.glTexHeight, 0, me.glPixelDataFormat, me.glPixelDataType, me.glPixPointer)
		}
	}
	if canMip && !me.noAutoMips {
		gl.GenerateMipmap(glTarget)
	}
	me.glSynced = true
}

func (me *FxImage) GpuDelete() {
	if me.glTex != 0 {
		gl.DeleteTextures(1, &me.glTex)
		me.glTex, me.glSynced = 0, false
	}
}

func (me *FxImage) GpuSynced() bool {
	return me.glSynced
}

func (me *FxImage) NoAutoMips() {
	me.noAutoMips = true
}

type FxImageOnLoad func(img interface{}, err error, async bool)
