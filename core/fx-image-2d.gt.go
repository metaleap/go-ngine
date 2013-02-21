package core

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	ugl "github.com/go3d/go-opengl/util"
)

type FxImage2D struct {
	FxImageBase
	InitFrom FxImageInitFrom

	glTex ugl.Texture2D
	img   image.Image
}

func (me *FxImage2D) init() {
	me.glTex.Init()
	me.FxImageBase.init(&me.glTex.TextureBase)
}

func (me *FxImage2D) dispose() {
	me.FxImageBase.dispose()
	me.Unload()
}

func (me *FxImage2D) GpuSync() (err error) {
	if err = me.glTex.PrepFromImage(me.Storage.Bgra, me.Storage.UintRev, me.img); err == nil {
		err = me.gpuSync(&me.glTex)
	}
	return
}

func (me *FxImage2D) Load() (err error) {
	me.Unload()
	me.img, err = me.InitFrom.loadImage(&me.FxImageBase)
	return
}

func (me *FxImage2D) Loaded() bool {
	return me.img != nil
}

func (me *FxImage2D) Unload() {
	me.img, me.glSynced = nil, false
}

//#begin-gt -gen-lib.gt T:FxImage2D

//	Initializes and returns a new FxImage2D with default parameters.
func NewFxImage2D() (me *FxImage2D) {
	me = &FxImage2D{}
	me.init()
	return
}

//	A hash-table of FxImage2Ds associated by IDs. Only for use in Core.Libs.
type LibFxImage2Ds map[string]*FxImage2D

//	Creates and initializes a new FxImage2D with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibFxImage2Ds) AddNew(id string) (obj *FxImage2D) {
	obj = NewFxImage2D()
	me[id] = obj
	return
}

func (me *LibFxImage2Ds) ctor() {
	*me = LibFxImage2Ds{}
}

func (me *LibFxImage2Ds) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibFxImage2Ds) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
