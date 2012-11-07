package core

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type assetFileOpener func (string) (io.ReadCloser, error)

type assetManager struct {
}

	func newAssetManager () *assetManager {
		var assMan = &assetManager {}
		return assMan
	}

	func (me *assetManager) LoadAssets (reader io.Reader) (*AssetProvider, error) {
		return newAssetProvider(reader)
	}

	func (me *assetManager) OpenLocalFile (absoluteOrAssetRootRelativeFilePath string) (io.ReadCloser, error) {
		return os.Open(me.ResolveLocalFilePath(absoluteOrAssetRootRelativeFilePath))
	}

	func (me *assetManager) OpenRemoteFile (fileUrl string) (rc io.ReadCloser, err error) {
		var resp *http.Response
		resp, err = http.Get(fileUrl)
		if resp != nil { rc = resp.Body }
		return
	}

	func (me *assetManager) ResolveLocalFilePath (absoluteOrAssetRootRelativeFilePath string) string {
		if strings.HasPrefix(absoluteOrAssetRootRelativeFilePath, "/") { return absoluteOrAssetRootRelativeFilePath }
		return filepath.Join(Core.Options.AssetRootDirPath, absoluteOrAssetRootRelativeFilePath)
	}
