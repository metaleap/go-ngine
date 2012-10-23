package core

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

type FTextureProvider func (args ... interface {}) (image.Image, error)

type tTextureProviders struct {
	IoReader, LocalFile, RemoteFile FTextureProvider
}

var (
	textureProviders = &tTextureProviders { textureProviderIoReader, textureProviderLocalFile, textureProviderRemoteFile }
)

func textureProviderIoReader (args ... interface {}) (img image.Image, err error) {
	img, _, err = image.Decode(args[0].(io.Reader))
	return
}

func textureProviderLocalFile (args ... interface {}) (img image.Image, err error) {
	var file *os.File
	var fpath = args[0].(string)
	if !strings.HasSuffix(fpath, "/") { fpath = Core.AssetManager.LocalFilePath(fpath) }
	if file, err = os.Open(fpath); err != nil { return }
	defer file.Close()
	return textureProviderIoReader(file)
}

func textureProviderRemoteFile (args ... interface {}) (img image.Image, err error) {
	var resp *http.Response
	if resp, err = http.Get(args[0].(string)); err != nil { return }
	defer resp.Body.Close()
	return textureProviderIoReader(resp.Body)
}
