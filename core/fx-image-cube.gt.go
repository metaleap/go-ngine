package core

import (
	"image"

	ugl "github.com/go3d/go-opengl/util"
)

type FxImageCube struct {
	FxImageBase
	InitFrom [6]FxImageInitFrom

	glTex  ugl.TextureCube
	images []image.Image
}

func (me *FxImageCube) init() {
	me.images = make([]image.Image, 6)
	me.glTex.Init()
	me.FxImageBase.init(&me.glTex.TextureBase)
	me.PreProcess.FlipY = false
}

func (me *FxImageCube) dispose() {
	me.FxImageBase.dispose()
	me.Unload()
}

func (me *FxImageCube) GpuSync() (err error) {
	if err = me.glTex.PrepFromImages(true, me.Storage.UintRev, me.images...); err == nil {
		err = me.gpuSync(&me.glTex)
	}
	return
}

func (me *FxImageCube) Load() (err error) {
	var (
		img image.Image
	)
	me.Unload()
	for i := 0; i < 6; i++ {
		img, err = me.InitFrom[i].loadImage(&me.FxImageBase)
		if err == nil {
			me.images[i] = img
		} else {
			break
		}
	}
	return
}

func (me *FxImageCube) Loaded() bool {
	for i := 0; i < len(me.images); i++ {
		if me.images[i] == nil {
			return false
		}
	}
	return true
}

func (me *FxImageCube) Unload() {
	me.glSynced = false
	for i := 0; i < len(me.images); i++ {
		me.images[i] = nil
	}
}

//#begin-gt -gen-lib.gt T:FxImageCube

//	Initializes and returns a new FxImageCube with default parameters.
func NewFxImageCube() (me *FxImageCube) {
	me = &FxImageCube{}
	me.init()
	return
}

//	A hash-table of FxImageCubes associated by IDs. Only for use in Core.Libs.
type LibFxImageCubes map[string]*FxImageCube

//	Creates and initializes a new FxImageCube with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibFxImageCubes) AddNew(id string) (obj *FxImageCube) {
	obj = NewFxImageCube()
	me[id] = obj
	return
}

func (me *LibFxImageCubes) ctor() {
	*me = LibFxImageCubes{}
}

func (me *LibFxImageCubes) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibFxImageCubes) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
