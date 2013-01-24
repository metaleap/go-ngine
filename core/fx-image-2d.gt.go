package core

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	ugl "github.com/go3d/go-glutil"
	ugo "github.com/metaleap/go-util"
)

type FxImage2D struct {
	FxImageBase

	glTex ugl.Texture2D
	img   image.Image
}

func (me *FxImage2D) init() {
	me.glTex.Init()
	me.FxImageBase.init(&me.glTex.TextureBase)
}

func (me *FxImage2D) dispose() {
	me.FxImageBase.dispose()
	me.Unload()
}

func (me *FxImage2D) GpuSync() (err error) {
	if err = ugl.ImageTextureProperties(me.img, &me.glTex.Width, &me.glTex.Height, &me.glTex.MipMap.NumLevels, &me.glTex.SizedInternalFormat, &me.glTex.PixelData.Format, &me.glTex.PixelData.Type, &me.glTex.PixelData.Ptr); err == nil {
		me.gpuSync(&me.glTex)
	}
	return
}

func (me *FxImage2D) Load() {
	var (
		loadAsync func()
		err       error
		img       image.Image
	)
	prov, arg, _ := me.provider()
	if me.AsyncNumAttempts != 0 {
		asyncResources[me] = false
		loadAsync = func() {
			for i := 0; (i < me.AsyncNumAttempts) || (me.AsyncNumAttempts < 0); i++ {
				if img, err = prov(arg); err == nil {
					break
				}
			}
			asyncResources[me] = true
			me.load_OnImg(img, err, true)
		}
		go loadAsync()
	} else {
		img, err = prov(arg)
		me.load_OnImg(img, err, false)
	}
}

func (me *FxImage2D) load_OnImg(img image.Image, err error, async bool) {
	me.Unload()
	if me.img = img; me.img != nil {
		//	If img is YCbCr or Paletted, convert to RGBA:
		switch me.img.(type) {
		case *image.YCbCr, *image.Paletted:
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
	if me.OnLoad != nil {
		me.OnLoad(img, err, async)
	}
	if err != nil {
		ugo.LogError(err)
	}
}

func (me *FxImage2D) Loaded() bool {
	return (me.img != nil)
}

func (me *FxImage2D) provider() (prov fxImage2DProvider, arg interface{}, remote bool) {
	if len(me.InitFrom.RawData) > 0 {
		prov, arg = fxImage2DProviderIoReader, bytes.NewReader(me.InitFrom.RawData)
	} else if len(me.InitFrom.RefUrl) > 0 {
		if remote = strings.Contains(me.InitFrom.RefUrl, "://"); remote {
			prov, arg = fxImage2DProviderRemoteFile, me.InitFrom.RefUrl
		} else {
			prov, arg = fxImage2DProviderLocalFile, me.InitFrom.RefUrl
		}
	}
	return
}

func (me *FxImage2D) IsRemote() (remote bool) {
	_, _, remote = me.provider()
	return
}

func (me *FxImage2D) Unload() {
	me.img, me.glSynced = nil, false
}

type fxImage2DProvider func(interface{}) (image.Image, error)

func fxImage2DProviderIoReader(arg interface{}) (img image.Image, err error) {
	var (
		ok     bool
		closer io.Closer
		reader io.Reader
	)
	if closer, ok = arg.(io.Closer); ok {
		defer closer.Close()
	}
	if reader, ok = arg.(io.Reader); ok {
		img, _, err = image.Decode(reader)
	}
	return
}

func fxImage2DProviderLocalFile(arg interface{}) (img image.Image, err error) {
	var rc io.ReadCloser
	if rc, err = Core.fileIO.openLocalFile(arg.(string)); err == nil {
		img, err = fxImage2DProviderIoReader(rc)
	}
	return
}

func fxImage2DProviderRemoteFile(arg interface{}) (img image.Image, err error) {
	var rc io.ReadCloser
	if rc, err = Core.fileIO.openRemoteFile(arg.(string)); err == nil {
		img, err = fxImage2DProviderIoReader(rc)
	}
	return
}

//#begin-gt -gen-lib.gt T:FxImage2D

//	Initializes and returns a new FxImage2D with default parameters.
func NewFxImage2D() (me *FxImage2D) {
	me = &FxImage2D{}
	me.init()
	return
}

//	A hash-table of FxImage2Ds associated by IDs. Only for use in Core.Libs.
type LibFxImage2Ds map[string]*FxImage2D

//	Creates and initializes a new FxImage2D with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibFxImage2Ds) AddNew(id string) (obj *FxImage2D) {
	obj = NewFxImage2D()
	me[id] = obj
	return
}

func (me *LibFxImage2Ds) ctor() {
	*me = LibFxImage2Ds{}
}

func (me *LibFxImage2Ds) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibFxImage2Ds) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
