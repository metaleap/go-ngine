package core

import (
	"image"

	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

var (
	asyncTextures = map[*Texture]bool {}
)

type textures map[string]*Texture

	func (me textures) Load (loadProvider TextureProvider, providerArgs ... interface {}) *Texture {
		var tex = newTexture()
		tex.Load(loadProvider, providerArgs ...)
		return tex
	}

	func (me textures) LoadAndSet (name string, loadProvider TextureProvider, providerArgs ... interface {}) *Texture {
		return me.Set(name, me.Load(loadProvider, providerArgs ...))
	}

	func (me textures) LoadAsync (loadProvider TextureProvider, providerArgs ... interface {}) *Texture {
		var tex = newTexture()
		tex.LoadAsync(loadProvider, providerArgs ...)
		return tex
	}

	func (me textures) LoadAsyncAndSet (name string, loadProvider TextureProvider, providerArgs ... interface {}) *Texture {
		return me.Set(name, me.LoadAsync(loadProvider, providerArgs ...))
	}

	func (me textures) NewParams (filter bool, filterAnisotropy float64) *textureParams {
		return newTextureParams(filter, filterAnisotropy)
	}

	func (me textures) Set (name string, tex *Texture) *Texture {
		me[name] = tex
		return tex
	}

type Texture struct {
	LastError error
	Params *textureParams

	img image.Image
	gpuSynced, noMipMap bool
	glTex gl.Uint
	glTexWidth, glTexHeight, glTexLevels gl.Sizei
	glPixPointer gl.Pointer
	glSizedInternalFormat, glPixelDataFormat, glPixelDataType gl.Enum
}

	func newTexture () *Texture {
		var tex = &Texture {}
		tex.Params = Core.Options.DefaultTextureParams
		return tex
	}

	func (me *Texture) GpuDelete () {
		if me.glTex != 0 {
			gl.DeleteTextures(1, &me.glTex)
			me.glTex = 0
		}
	}

	func (me *Texture) GpuSync () {
		me.gpuSynced = false
		me.GpuDelete()
		gl.GenTextures(1, &me.glTex)
		gl.BindTexture(gl.TEXTURE_2D, me.glTex)
		defer gl.BindTexture(gl.TEXTURE_2D, 0)
		me.Params.apply(me)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		if me.img != nil {
			me.glPixPointer = ugl.ImageTextureProperties(me.img, &me.glTexWidth, &me.glTexHeight, &me.glTexLevels, &me.glSizedInternalFormat, &me.glPixelDataFormat, &me.glPixelDataType)
			if ugl.IsGl42 {
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

	func (me *Texture) GpuSynced () bool {
		return me.gpuSynced
	}

	func (me *Texture) load_OnImg (img image.Image, err error) error {
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

	func (me *Texture) Load (provider TextureProvider, args ... interface {}) {
		me.load_OnImg(provider(args ...))
	}

	func (me *Texture) LoadAsync (provider TextureProvider, args ... interface {}) {
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

	func (me *Texture) Loaded () bool {
		return me.img != nil
	}

	func (me *Texture) SuppressMipMaps () {
		me.noMipMap = true
	}

	func (me *Texture) Unload () {
		me.img, me.glPixPointer = nil, gl.Pointer(nil)
	}
