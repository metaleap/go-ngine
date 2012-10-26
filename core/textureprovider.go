package core

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type FTextureProvider func (args ... interface {}) (image.Image, error)

type tTextureProviders struct {
	IoReader, LocalFile, RemoteFile FTextureProvider
}

var (
	TextureProviders = &tTextureProviders { textureProviderIoReader, textureProviderLocalFile, textureProviderRemoteFile }
)

func textureProviderIoReader (args ... interface {}) (img image.Image, err error) {
	if cl, isCl := args[0].(io.Closer); isCl { defer cl.Close() }
	img, _, err = image.Decode(args[0].(io.Reader))
	return
}

func textureProviderLocalFile (args ... interface {}) (img image.Image, err error) {
	return textureProviderIoReader(Core.AssetManager.OpenLocalFile(args[0].(string)))
}

func textureProviderRemoteFile (args ... interface {}) (img image.Image, err error) {
	return textureProviderIoReader(Core.AssetManager.OpenRemoteFile(args[0].(string)))
}
