package core

import (
	"encoding/binary"
	"image"
	"image/color"
	"os"
	"path/filepath"

	ugo "github.com/metaleap/go-util"
	ugfx "github.com/metaleap/go-util/gfx"
)

type fxImageRaw struct {
	pix                         []byte
	bounds                      [2]uint64
	needImg                     bool
	dirPath, fileName, fullPath string
	src                         os.FileInfo
}

func newFxImageRaw(init *FxImageInitFrom, fxImg *FxImageBase) (me *fxImageRaw, err error) {
	var dst os.FileInfo
	me = &fxImageRaw{
		needImg: true, dirPath: Core.fileIO.resolveLocalFilePath(filepath.Join(Options.AppDir.Temp.BaseName, Options.AppDir.Temp.Textures)),
		fileName:/*filepath.Base(init.RefUrl) + "-" +*/ ugo.Hash(strf("%s_%t_%t_%t_%t", init.RefUrl, fxImg.PreProcess.FlipY, fxImg.PreProcess.ToBgra, fxImg.PreProcess.ToLinear, fxImg.Storage.Bgra)),
	}
	me.fullPath = filepath.Join(me.dirPath, me.fileName)

	if me.src, err = os.Stat(Core.fileIO.resolveLocalFilePath(init.RefUrl)); err != nil {
		me.src, err = nil, nil
	}
	if dst, err = os.Stat(me.fullPath); err != nil {
		err, dst = nil, nil
	}

	// if the cached-file exists -- can we use it?
	if dst != nil && dst.Size() > 0 {
		if me.src == nil || dst.ModTime().UnixNano() > me.src.ModTime().UnixNano() {
			// cached-file is newer than source file: check.
			var file *os.File
			if file, err = os.Open(me.fullPath); err == nil {
				// ceched-file can be opened: check.
				defer file.Close()
				if err = binary.Read(file, binary.LittleEndian, &me.bounds); err == nil {
					var meta [3]int64
					if err = binary.Read(file, binary.LittleEndian, &meta); err == nil {
						if me.src == nil || (me.src.ModTime().UnixNano() == meta[1] && me.src.Size() == meta[2]) {
							// source file size and mod-time matches those remembered in cached-file: check.
							if me.pix = make([]byte, meta[0]); meta[0] > 0 {
								if err = binary.Read(file, binary.LittleEndian, me.pix); err == nil {
									// no need to load and process the source file, cached-file loaded successfully
									me.needImg = false
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func (me *fxImageRaw) At(x, y int) (col color.Color) {
	// dirty, but this part of *fxImageRaw implementing image.Image isn't ever called currently, anyway
	return nil
}

func (me *fxImageRaw) Bounds() (r image.Rectangle) {
	r.Max.X, r.Max.Y = int(me.bounds[0]), int(me.bounds[1])
	return
}

func (me *fxImageRaw) ColorModel() color.Model {
	// dirty, but this part of *fxImageRaw implementing image.Image isn't ever called currently, anyway
	return nil
}

func (me *fxImageRaw) Pix() []byte {
	return me.pix
}

func (me *fxImageRaw) setImg(img image.Image) (err error) {
	me.bounds[0], me.bounds[1] = uint64(img.Bounds().Dx()), uint64(img.Bounds().Dy())
	_, me.pix = ugfx.CloneImage(img, true)
	if len(me.fullPath) > 0 {
		var file *os.File
		if file, err = os.Create(me.fullPath); err == nil {
			defer file.Close()
			if err = binary.Write(file, binary.LittleEndian, &me.bounds); err == nil {
				if err = binary.Write(file, binary.LittleEndian, [3]int64{int64(len(me.pix)), me.src.ModTime().UnixNano(), me.src.Size()}); err == nil {
					err = binary.Write(file, binary.LittleEndian, me.pix)
				}
			}
		}
	}
	return
}
