package core

import (
	"image"
	"strings"

	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
	nga "github.com/go3d/go-ngine/assets"
)

var (
	asyncTextures = map[*Texture]bool{}
)

type textures map[string]*Texture

func (me textures) add(def *nga.FxImageDef) (item *Texture) {
	item = newTexture(def)
	me[def.ID] = item
	return
}

func (me textures) NewParams(filter bool, filterAnisotropy float64) *textureParams {
	return newTextureParams(filter, filterAnisotropy)
}

func (me textures) syncAssetChanges() {
	var (
		item *Texture
		id   string
	)
	for _, def := range nga.FxImageDefs.M {
		if item = me[def.ID]; item == nil {
			item = me.add(def)
		}
	}
	for id, item = range me {
		if nga.FxImageDefs.M[item.ID] == nil {
			delete(me, id)
			item.dispose()
		}
	}
}

type Texture struct {
	*nga.FxImageDef
	LastError error
	Params    *textureParams

	img                                                       image.Image
	gpuSynced, noMipMap                                       bool
	glTex                                                     gl.Uint
	glTexWidth, glTexHeight, glTexLevels                      gl.Sizei
	glPixPointer                                              gl.Pointer
	glSizedInternalFormat, glPixelDataFormat, glPixelDataType gl.Enum
}

func newTexture(def *nga.FxImageDef) (me *Texture) {
	me = &Texture{}
	me.FxImageDef = def
	me.FxImageDef.OnSync = func() {
		me.load()
		me.GpuSync()
	}
	me.Params = Core.Options.DefaultTextureParams
	return
}

func (me *Texture) dispose() {
	me.Unload()
	me.GpuDelete()
}

func (me *Texture) GpuDelete() {
	if me.glTex != 0 {
		gl.DeleteTextures(1, &me.glTex)
		me.glTex = 0
	}
}

func (me *Texture) GpuSync() {
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
		if !me.noMipMap {
			gl.GenerateMipmap(gl.TEXTURE_2D)
		}
	}
	me.gpuSynced = true
}

func (me *Texture) GpuSynced() bool {
	return me.gpuSynced
}

func (me *Texture) load_OnImg(img image.Image, err error) error {
	me.gpuSynced, me.img = false, nil
	if err != nil {
		me.LastError = err
	}
	if me.Loaded() {
		me.Unload()
	}
	if me.img = img; me.img != nil {
		conv := false
		switch me.img.(type) {
		case *image.YCbCr, *image.Paletted:
			conv = true
		}
		if conv {
			nuW, nuH := me.img.Bounds().Dx(), me.img.Bounds().Dy()
			nuImage := image.NewRGBA(image.Rect(0, 0, nuW, nuH))
			for x := 0; x < nuW; x++ {
				for y := 0; y < nuH; y++ {
					nuImage.Set(x, y, me.img.At(x, y))
				}
			}
			me.img = nuImage
		}
	}
	return err
}

func (me *Texture) load() {
	prov, arg, remote := me.provider()
	if remote {
		me.loadAsync(prov, arg)
	} else {
		me.load_OnImg(prov(arg))
	}
}

func (me *Texture) loadAsync(prov TextureProvider, arg interface{}) {
	me.gpuSynced = false
	me.Unload()
	asyncTextures[me] = false
	go func() {
		if err := me.load_OnImg(prov(arg)); err != nil {
			//	mark as "done anyway" in the async queue.
			asyncTextures[me] = true
		}
	}()
}

func (me *Texture) Loaded() bool {
	return me.img != nil
}

func (me *Texture) provider() (prov TextureProvider, arg interface{}, remote bool) {
	if me.FxImageDef != nil {
		if me.FxImageDef.Init.From != nil {
			if len(me.FxImageDef.Init.From.Raw.Data) > 0 {
				prov, arg = TextureProviders.IoReader, me.FxImageDef.Init.From.Raw.Data
			} else if len(me.FxImageDef.Init.From.RefUrl) > 0 {
				if remote = strings.Contains(me.FxImageDef.Init.From.RefUrl, "://"); remote {
					prov, arg = TextureProviders.RemoteFile, me.FxImageDef.Init.From.RefUrl
				} else {
					prov, arg = TextureProviders.LocalFile, me.FxImageDef.Init.From.RefUrl
				}
			}
		}
	}
	return
}

func (me *Texture) SuppressMipMaps() {
	me.noMipMap = true
}

func (me *Texture) Unload() {
	me.img, me.glPixPointer = nil, gl.Pointer(nil)
}
