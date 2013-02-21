package core

import (
	"bytes"
	"image"
	"image/draw"
	"io"
	"strings"

	ugfx "github.com/metaleap/go-util/gfx"
)

type fxImageRaw struct {
	data []byte
}

type FxImageInitFrom struct {
	RawData []byte
	RefUrl  string
}

func (me *FxImageInitFrom) loadImage(fxImg *FxImageBase) (img image.Image, err error) {
	// var rawImg *fxImageRaw
	prov, arg, _ := me.provider()

	img, err = prov(arg)
	if img != nil && err == nil {
		switch img.(type) {
		case *image.YCbCr, *image.Paletted:
			rect := img.Bounds()
			tmpImg := image.NewRGBA(rect)
			draw.Draw(tmpImg, rect, img, rect.Min, draw.Src)
			img = tmpImg
		}
		if fxImg.PreProcess.FlipY || fxImg.PreProcess.ToLinear || fxImg.PreProcess.ToBgra {
			pic := ugfx.CloneImage(img, false)
			ugfx.PreprocessImage(img, pic, fxImg.PreProcess.FlipY, fxImg.PreProcess.ToBgra, fxImg.PreProcess.ToLinear)
			img = pic
		}
	}
	return
}

func (me *FxImageInitFrom) provider() (prov fxImageProvider, arg interface{}, remote bool) {
	if len(me.RawData) > 0 {
		prov, arg = fxImageProviderIoReader, bytes.NewReader(me.RawData)
	} else if len(me.RefUrl) > 0 {
		if remote = strings.Contains(me.RefUrl, "://"); remote {
			prov, arg = fxImageProviderRemoteFile, me.RefUrl
		} else {
			prov, arg = fxImageProviderLocalFile, me.RefUrl
		}
	}
	return
}

type fxImageProvider func(interface{}) (image.Image, error)

func fxImageProviderIoReader(arg interface{}) (img image.Image, err error) {
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

func fxImageProviderLocalFile(arg interface{}) (img image.Image, err error) {
	var rc io.ReadCloser
	if rc, err = Core.fileIO.openLocalFile(arg.(string)); err == nil {
		img, err = fxImageProviderIoReader(rc)
	}
	return
}

func fxImageProviderRemoteFile(arg interface{}) (img image.Image, err error) {
	var rc io.ReadCloser
	if rc, err = Core.fileIO.openRemoteFile(arg.(string)); err == nil {
		img, err = fxImageProviderIoReader(rc)
	}
	return
}
