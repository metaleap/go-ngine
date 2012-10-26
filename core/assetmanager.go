package core

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FAssetFileOpener func (string) (io.ReadCloser, error)

type tModelFile struct {
}

	func loadModelFile () *tModelFile {
		return nil
	}

type tAssetManager struct {
}

	func newAssetManager () *tAssetManager {
		var assMan = &tAssetManager {}
		return assMan
	}

	func (me *tAssetManager) LoadModelFile () *tModelFile {
		return nil
	}

	func (me *tAssetManager) OpenLocalFile (absoluteOrAssetRootRelativeFilePath string) (io.ReadCloser, error) {
		return os.Open(me.ResolveLocalFilePath(absoluteOrAssetRootRelativeFilePath))
	}

	func (me *tAssetManager) OpenRemoteFile (fileUrl string) (io.ReadCloser, error) {
		var resp *http.Response
		var err error
		resp, err = http.Get(fileUrl)
		return resp.Body, err
	}

	func (me *tAssetManager) ResolveLocalFilePath (absoluteOrAssetRootRelativeFilePath string) string {
		if strings.HasPrefix(absoluteOrAssetRootRelativeFilePath, "/") { return absoluteOrAssetRootRelativeFilePath }
		return filepath.Join(Core.Options.AssetRootDirPath, absoluteOrAssetRootRelativeFilePath)
	}
