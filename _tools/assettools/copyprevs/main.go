package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	ioutil "github.com/go3d/go-util/io"
)

var (
	srcDirPath = "Q:\\oga\\yughues\\"
	prevDirPath = "Q:\\oga\\prevs\\"
)

func processPreview (filePath string, recurse bool) bool {
	var relPath = strings.Replace(strings.Replace(filePath, srcDirPath, "", -1), string(os.PathSeparator), "_", -1)
	var newFilePath = filepath.Join(prevDirPath, relPath)
	if err := ioutil.CopyFile(filePath, newFilePath); err != nil {
		log.Printf("ERR CopyFile(%v --> %v) %v", filePath, newFilePath, err)
	}
	return false
}

func main () {
	ioutil.WalkDirectory(srcDirPath, ".jpg", processPreview, true)
}
