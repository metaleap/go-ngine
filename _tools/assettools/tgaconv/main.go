package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	ioutil "github.com/go3d/go-util/io"

	"github.com/ftrvxmtrx/tga"
)

var (
	fSrcDirPath = flag.String("src", "C:\\oga\\yughues", "TGA src dir root path")
	fRecurse = flag.Bool("rec", false, "recurse into sub-dirs")
	fOutDirPath = flag.String("out", "C:\\oga\\yughues_png", "PNG out dir root path")
)

func processFile (srcFilePath string, recurse bool) bool {
	var outFilePath = strings.Replace(strings.Replace(srcFilePath, *fSrcDirPath, *fOutDirPath, -1), ".tga", ".png", -1)
	var err error
	var dirPath = filepath.Dir(outFilePath)
	var tgaFile, pngFile *os.File
	var img image.Image
	if err = ioutil.EnsureDirExists(dirPath); err != nil {
		log.Printf("ERR MkDir %v: %v\n", dirPath, err)
	} else if ioutil.DirExists(dirPath) {
		if tgaFile, err = os.OpenFile(srcFilePath, os.O_RDONLY, os.ModePerm); err != nil {
			log.Printf("ERR OpenFile %v: %v\n", srcFilePath, err)
		} else {
			defer tgaFile.Close()
			if img, err = tga.Decode(tgaFile); err != nil {
				log.Printf("ERR Decode %v: %v\n", srcFilePath, err)
			} else {
				if pngFile, err = os.Create(outFilePath); err != nil {
					log.Printf("ERR CreateFile %v: %v\n", outFilePath, err)
				} else {
					defer pngFile.Close()
					if err = png.Encode(pngFile, img); err != nil {
						log.Printf("ERR Encode %v: %v\n", outFilePath, err)
					} else {
						log.Printf("DONE %v\n", outFilePath)
					}
				}
			}
		}
	}
	return recurse
}

func main () {
	flag.Parse()
	ioutil.WalkDirectory(*fSrcDirPath, ".tga", processFile, *fRecurse)
}
