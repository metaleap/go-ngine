package core

import (
	_ "github.com/ftrvxmtrx/tga"

	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type TextureProvider func (arg interface {}) (image.Image, error)

type textureProviders struct {
	IoReader, LocalFile, RemoteFile TextureProvider
}

var (
	//	A collection of all "texture providers" known to go:ngine.
	//	Usage: core.TextureProviders.IoReader, core.TextureProviders.LocalFile, or core.TextureProviders.RemoteFile
	TextureProviders = &textureProviders { textureProviderIoReader, textureProviderLocalFile, textureProviderRemoteFile }
)

func textureProviderIoReader (arg interface{}) (img image.Image, err error) {
	if closer, isCl := arg.(io.Closer); isCl { defer closer.Close() }
	if reader, isR  := arg.(io.Reader); isR { img, _, err = image.Decode(reader) }
	return
}

func textureProviderLocalFile (arg interface{}) (img image.Image, err error) {
	var rc io.ReadCloser
	if rc, err = Core.AssetManager.OpenLocalFile(arg.(string)); err == nil {
		img, err = textureProviderIoReader(rc)
	}
	return
}

func textureProviderRemoteFile (arg interface{}) (img image.Image, err error) {
	var rc io.ReadCloser
	if rc, err = Core.AssetManager.OpenRemoteFile(arg.(string)); err == nil {
		img, err = textureProviderIoReader(rc)
	}
	return
}
