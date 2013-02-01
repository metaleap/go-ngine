package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	uio "github.com/metaleap/go-util/io"
)

var (
	srcDirPath  = "Q:\\oga\\yughues\\"
	prevDirPath = "Q:\\oga\\prevs\\"
)

func processPreview(walker *uio.DirWalker, filePath string, isDir bool) bool {
	if strings.HasSuffix(filePath, ".jpg") {
		var (
			relPath     = strings.Replace(strings.Replace(filePath, srcDirPath, "", -1), string(os.PathSeparator), "_", -1)
			newFilePath = filepath.Join(prevDirPath, relPath)
		)
		if err := uio.CopyFile(filePath, newFilePath); err != nil {
			log.Printf("ERR CopyFile(%v --> %v) %v", filePath, newFilePath, err)
		}
	}
	return true
}

func main() {
	if errs := uio.NewDirWalker(nil, processPreview).Walk(srcDirPath); len(errs) > 0 {
		panic(errs[0])
	}
}
