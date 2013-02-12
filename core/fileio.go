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

func (me *fileIO) openLocalFile(absOrRelFilePath string) (io.ReadCloser, error) {
	return os.Open(me.resolveLocalFilePath(absOrRelFilePath))
}

func (me *fileIO) openRemoteFile(fileUrl string) (rc io.ReadCloser, err error) {
	var resp *http.Response
	resp, err = http.Get(fileUrl)
	if resp != nil {
		rc = resp.Body
	}
	return
}

func (me *fileIO) resolveLocalFilePath(absOrRelFilePath string) string {
	if strings.HasPrefix(absOrRelFilePath, "/") {
		return absOrRelFilePath
	}
	return filepath.Join(Core.Options.AppDir.BasePath, absOrRelFilePath)
}
