package core

import (
	"image"

	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

var (
	asyncTextures = map[*TTexture]bool {}
)

type tTextures map[string]*TTexture

	func (me *tTextures) Load (loadProvider FTextureProvider, providerArgs ... interface {}) *TTexture {
		var tex = newTexture()
		tex.Load(loadProvider, providerArgs ...)
		return tex
	}

	func (me *tTextures) LoadAsync (loadProvider FTextureProvider, providerArgs ... interface {}) *TTexture {
		var tex = newTexture()
		tex.LoadAsync(loadProvider, providerArgs ...)
		return tex
	}

	func (me *tTextures) NewParams (filter bool, filterAnisotropy float64) *tTextureParams {
		return newTextureParams(filter, filterAnisotropy)
	}

type TTexture struct {
	LastError error
	Params *tTextureParams

	img image.Image
	gpuSynced, noMipMap bool
	glTex gl.Uint
	glTexWidth, glTexHeight, glTexLevels gl.Sizei
	glPixPointer gl.Pointer
	glSizedInternalFormat, glPixelDataFormat, glPixelDataType gl.Enum
}

	func newTexture () *TTexture {
		var tex = &TTexture {}
		tex.Params = Core.Options.DefaultTextureParams
		return tex
	}

	func (me *TTexture) GpuDelete () {
		if me.glTex != 0 {
			gl.DeleteTextures(1, &me.glTex)
			me.glTex = 0
		}
	}

	func (me *TTexture) GpuSync () {
		me.gpuSynced = false
		me.GpuDelete()
		gl.GenTextures(1, &me.glTex)
		gl.BindTexture(gl.TEXTURE_2D, me.glTex)
		defer gl.BindTexture(gl.TEXTURE_2D, 0)
		me.Params.apply(me)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		if me.img != nil {
			me.glPixPointer = glutil.ImageTextureProperties(me.img, &me.glTexWidth, &me.glTexHeight, &me.glTexLevels, &me.glSizedInternalFormat, &me.glPixelDataFormat, &me.glPixelDataType)
			if glutil.IsGl43 {
				gl.TexStorage2D(gl.TEXTURE_2D, me.glTexLevels, me.glSizedInternalFormat, me.glTexWidth, me.glTexHeight)
				gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, me.glTexWidth, me.glTexHeight, me.glPixelDataFormat, me.glPixelDataType, me.glPixPointer)
			} else {
				gl.TexImage2D(gl.TEXTURE_2D, 0, gl.Int(me.glSizedInternalFormat), me.glTexWidth, me.glTexHeight, 0, me.glPixelDataFormat, me.glPixelDataType, me.glPixPointer)
				glLogLastError("ttex.gpusync(9)")
			}
			if !me.noMipMap { gl.GenerateMipmap(gl.TEXTURE_2D) }
		}
		me.gpuSynced = true
	}

	func (me *TTexture) GpuSynced () bool {
		return me.gpuSynced
	}

	func (me *TTexture) load_OnImg (img image.Image, err error) error {
		var nuW, nuH int
		var conv = false
		var nuImage *image.RGBA
		me.gpuSynced, me.img = false, nil
		if err != nil { me.LastError = err }
		if me.Loaded() { me.Unload() }
		if me.img = img; me.img != nil {
			switch me.img.(type) {
			case *image.YCbCr, *image.Paletted:
				conv = true
			}
			if conv {
				nuW, nuH = me.img.Bounds().Dx(), me.img.Bounds().Dy()
				nuImage = image.NewRGBA(image.Rect(0, 0, nuW, nuH))
				for x := 0; x < nuW; x++ { for y := 0; y < nuH; y++ { nuImage.Set(x, y, me.img.At(x, y)) } }
				me.img = nuImage
			}
		}
		return err
	}

	func (me *TTexture) Load (provider FTextureProvider, args ... interface {}) {
		me.load_OnImg(provider(args ...))
	}

	func (me *TTexture) LoadAsync (provider FTextureProvider, args ... interface {}) {
		me.gpuSynced = false
		me.Unload()
		asyncTextures[me] = false
		go func () {
			if err := me.load_OnImg(provider(args ...)); err != nil {
				//	mark as "done" anyway in the async queue.
				asyncTextures[me] = true
			}
		} ()
	}

	func (me *TTexture) Loaded () bool {
		return me.img != nil
	}

	func (me *TTexture) SuppressMipMaps () {
		me.noMipMap = true
	}

	func (me *TTexture) Unload () {
		me.img, me.glPixPointer = nil, gl.Pointer(nil)
	}
