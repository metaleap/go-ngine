package core

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	ugl "github.com/go3d/go-opengl/util"
	ugfx "github.com/metaleap/go-util/gfx"
)

type FxImage2D struct {
	FxImageBase
	FlipY           bool
	ConvertToLinear bool

	glTex ugl.Texture2D
	img   image.Image
}

func (me *FxImage2D) init() {
	me.FlipY, me.ConvertToLinear = true, true
	me.glTex.Init()
	me.FxImageBase.init(&me.glTex.TextureBase)
}

func (me *FxImage2D) dispose() {
	me.FxImageBase.dispose()
	me.Unload()
}

func (me *FxImage2D) GpuSync() (err error) {
	if err = me.glTex.PrepFromImage(me.img); err == nil {
		err = me.gpuSync(&me.glTex)
	}
	return
}

func (me *FxImage2D) Load() (err error) {
	var img image.Image
	prov, arg, _ := me.provider()
	img, err = prov(arg)
	me.Unload()
	if img != nil && err == nil {
		//	If img is YCbCr or Paletted, convert to RGBA:
		switch img.(type) {
		case *image.YCbCr, *image.Paletted:
			b := img.Bounds()
			nuImage := image.NewRGBA(b)
			draw.Draw(nuImage, b, img, b.Min, draw.Src)
			img = nuImage
		}
		if me.FlipY || me.ConvertToLinear {
			pic := ugfx.CloneImage(img, !me.FlipY)
			if me.FlipY {
				ugfx.FlipVertical(img, pic)
			}
			if me.ConvertToLinear {
				ugfx.SrgbToLinear(pic, pic)
			}
			img = pic
		}
		me.img = img
		if me.OnLoad != nil {
			me.OnLoad(img, err)
		}
	}
	if err != nil {
		Diag.LogErr(err)
	}
	return
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
