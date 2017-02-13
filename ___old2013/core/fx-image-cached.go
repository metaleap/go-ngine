package core

import (
	"encoding/base64"
	"encoding/binary"
	"hash/fnv"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/metaleap/go-util-gfx"
	"github.com/go-util-hash"
)

type fxImageCached struct {
	pix      []byte
	bounds   [2]uint64
	needImg  bool
	fullPath string
	src      os.FileInfo
}

func newFxImageCached(init *FxImageInitFrom, fxImg *FxImageBase) (me *fxImageCached) {
	var (
		dst os.FileInfo
		err error
	)
	dirPath := Core.fileIO.resolveLocalFilePath(filepath.Join(Options.AppDir.Temp.BaseName, Options.AppDir.Temp.CachedTextures))
	hash, _ := uhash.WriteAndSum(fnv.New64a(), []byte(strf("%s_%t_%t_%t_%t", init.RefUrl, fxImg.Preprocess.FlipY, fxImg.Preprocess.ToBgra, fxImg.Preprocess.ToLinear, fxImg.Storage.Gpu.Bgra)), nil)
	fileName := base64.URLEncoding.EncodeToString(hash)
	me = &fxImageCached{needImg: true, fullPath: filepath.Join(dirPath, fileName)}

	if me.src, err = os.Stat(Core.fileIO.resolveLocalFilePath(init.RefUrl)); err != nil {
		me.src = nil
	}
	if dst, err = os.Stat(me.fullPath); err != nil {
		dst = nil
	}

	// if the cached-file exists -- can we use it?
	if dst != nil && dst.Size() > 0 {
		if me.src == nil || dst.ModTime().UnixNano() > me.src.ModTime().UnixNano() {
			// cached-file is newer than source file: check.
			var file *os.File
			if file, err = os.Open(me.fullPath); err == nil {
				// ceched-file can be opened: check.
				defer file.Close()
				unpacker := fxImg.Storage.DiskCache.Decompressor(file)
				defer unpacker.Close()
				if err = binary.Read(unpacker, binary.LittleEndian, &me.bounds); err == nil {
					var meta [3]int64
					if err = binary.Read(unpacker, binary.LittleEndian, &meta); err == nil {
						if me.src == nil || (me.src.ModTime().UnixNano() == meta[1] && me.src.Size() == meta[2]) {
							// source file size and mod-time matches those remembered in cached-file: check.
							if me.pix = make([]byte, meta[0]); meta[0] > 0 {
								if err = binary.Read(unpacker, binary.LittleEndian, me.pix); err == nil {
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

// Returns nil, but this part of *fxImageCached implementing image.Image isn't ever called currently, anyway
func (me *fxImageCached) At(x, y int) (col color.Color) {
	return nil
}

//	Implementing image.Image.Bounds()
func (me *fxImageCached) Bounds() (r image.Rectangle) {
	r.Max.X, r.Max.Y = int(me.bounds[0]), int(me.bounds[1])
	return
}

// Returns nil, but this part of *fxImageCached implementing image.Image isn't ever called currently, anyway
func (me *fxImageCached) ColorModel() color.Model {
	return nil
}

//	Used by ugl.TextureXYZ.PrepFromImage/s()
func (me *fxImageCached) Pix() []byte {
	return me.pix
}

func (me *fxImageCached) setImg(img image.Image, fxImg *FxImageBase) (err error) {
	me.bounds[0], me.bounds[1] = uint64(img.Bounds().Dx()), uint64(img.Bounds().Dy())
	_, me.pix = ugfx.CreateLike(img, true)
	if len(me.fullPath) > 0 {
		var file *os.File
		if file, err = os.Create(me.fullPath); err == nil {
			defer file.Close()
			packer := fxImg.Storage.DiskCache.Compressor(file)
			defer packer.Close()
			if err = binary.Write(packer, binary.LittleEndian, &me.bounds); err == nil {
				if err = binary.Write(packer, binary.LittleEndian, [3]int64{int64(len(me.pix)), me.src.ModTime().UnixNano(), me.src.Size()}); err == nil {
					err = binary.Write(packer, binary.LittleEndian, me.pix)
				}
			}
		}
	}
	return
}
