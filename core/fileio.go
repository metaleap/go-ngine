package core

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type fileIO struct {
}

func newFileIO() (me *fileIO) {
	me = &fileIO{}
	return
}

func (me *fileIO) openLocalFile(absoluteOrAssetRootRelativeFilePath string) (io.ReadCloser, error) {
	return os.Open(me.resolveLocalFilePath(absoluteOrAssetRootRelativeFilePath))
}

func (me *fileIO) openRemoteFile(fileUrl string) (rc io.ReadCloser, err error) {
	var resp *http.Response
	resp, err = http.Get(fileUrl)
	if resp != nil {
		rc = resp.Body
	}
	return
}

func (me *fileIO) resolveLocalFilePath(absoluteOrAssetRootRelativeFilePath string) string {
	if strings.HasPrefix(absoluteOrAssetRootRelativeFilePath, "/") {
		return absoluteOrAssetRootRelativeFilePath
	}
	return filepath.Join(Core.Options.AssetRootDirPath, absoluteOrAssetRootRelativeFilePath)
}
