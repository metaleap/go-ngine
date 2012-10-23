package core

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
)

type FTextureProvider func (args ... interface {}) (image.Image, error)

type tTextureProviders struct {
	LocalFile, RemoteFile, Uint8Array FTextureProvider
}

var (
	textureProviders = &tTextureProviders { textureProviderLocalFile, textureProviderRemoteFile, textureProviderUint8Array }
)

func textureProviderLocalFile (args ... interface {}) (img image.Image, err error) {
	var file *os.File
	if file, err = os.Open(args[0].(string)); err != nil { return }
	defer file.Close()
	return textureProviderReader(file)
}

func textureProviderReader (args ... interface {}) (img image.Image, err error) {
	img, _, err = image.Decode(args[0].(io.Reader))
	return
}

func textureProviderRemoteFile (args ... interface {}) (img image.Image, err error) {
	var resp *http.Response
	if resp, err = http.Get(args[0].(string)); err != nil { return }
	defer resp.Body.Close()
	return textureProviderReader(resp.Body)
}

func textureProviderUint8Array (args ... interface {}) (image.Image, error) {
	return nil, nil
}
