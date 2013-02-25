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
	me.Unload()
	me.glTex.Init()
	me.FxImageBase.init(&me.glTex.TextureBase)
	me.PreProcess.FlipY = false
}

func (me *FxImageCube) dispose() {
	me.FxImageBase.dispose()
	me.Unload()
}

func (me *FxImageCube) GpuSync() (err error) {
	if err = me.glTex.PrepFromImages(me.Storage.Gpu.Bgra, me.Storage.Gpu.UintRev, me.images...); err == nil {
		err = me.gpuSync(&me.glTex)
	}
	return
}

func (me *FxImageCube) Load() (err error) {
	me.Unload()
	for i := 0; i < 6; i++ {
		if me.images[i], err = me.InitFrom[i].loadImage(&me.FxImageBase); err != nil {
			break
		}
	}
	return
}

func (me *FxImageCube) Loaded() bool {
	for i := 0; i < 6; i++ {
		if me.images[i] == nil {
			return false
		}
	}
	return true
}

func (me *FxImageCube) Unload() {
	me.glSynced, me.images = false, make([]image.Image, 6)
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
	*me = make(LibFxImageCubes, 50)
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
