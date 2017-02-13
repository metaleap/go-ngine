package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	uio "github.com/metaleap/go-util-fs"

	"github.com/ftrvxmtrx/tga"
)

var (
	fSrcDirPath = flag.String("src", "Q:\\oga\\yughues", "TGA src dir root path")
	fOutDirPath = flag.String("out", "Q:\\oga\\yughues_png", "PNG out dir root path")
)

func processFile(srcFilePath string) (keepWalking bool) {
	if strings.HasSuffix(srcFilePath, ".tga") {
		var (
			outFilePath      = strings.Replace(strings.Replace(srcFilePath, *fSrcDirPath, *fOutDirPath, -1), ".tga", ".png", -1)
			dirPath          = filepath.Dir(outFilePath)
			tgaFile, pngFile *os.File
			img              image.Image
			err              error
		)
		if err = uio.EnsureDirExists(dirPath); err != nil {
			log.Printf("ERR MkDir %v: %v\n", dirPath, err)
		} else if uio.DirExists(dirPath) {
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
	}
	return true
}

func main() {
	flag.Parse()
	if errs := uio.NewDirWalker(true, nil, processFile).Walk(*fSrcDirPath); len(errs) > 0 {
		panic(errs[0])
	}
}
