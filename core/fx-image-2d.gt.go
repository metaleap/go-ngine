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
	me.FxImageBase.init()
}

func (me *FxImage2D) dispose() {
	me.Unload()
}

func (me *FxImage2D) GpuDelete() {
	me.FxImageBase.gpuDelete(&me.glTex)
}

func (me *FxImage2D) GpuSync() (err error) {
	if err = me.glTex.PrepFromImage(me.Storage.Gpu.Bgra, me.Storage.Gpu.UintRev, me.img); err == nil {
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

func (me *FxImage2D) NoAutoMips() {
	me.FxImageBase.noAutoMips(&me.glTex.TextureBase)
}

func (me *FxImage2D) Unload() {
	me.img, me.glSynced = nil, false
}






//#begin-gt -gen-lib.gt T:FxImage2D L:Core.Libs.Images.Tex2D

//	Only used for Core.Libs.Images.Tex2D
type FxImage2DLib []FxImage2D

func (me *FxImage2DLib) AddNew() (id int) {
	id = -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			id = i
			break
		}
	}
	if id == -1 {
		if id = len(*me); id == cap(*me) {
			nu := make(FxImage2DLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, FxImage2D{})
	}
	ref := &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *FxImage2DLib) Compact() {
	var (
		before, after []FxImage2D
		ref           *FxImage2D
		oldID, i      int
		compact       bool
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			compact, before, after = true, (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	if compact {
		changed := make(map[int]int, len(*me))
		for i = 0; i < len(*me); i++ {
			if ref = &(*me)[i]; ref.ID != i {
				oldID, ref.ID = ref.ID, i
				changed[oldID] = i
			}
		}
		if len(changed) > 0 {
			me.onFxImage2DIDsChanged(changed)
		}
	}
}

func (me *FxImage2DLib) init() {
	*me = make(FxImage2DLib, 0, Options.Libs.InitialCap)
}

func (me *FxImage2DLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me FxImage2DLib) get(id int) (ref *FxImage2D) {
	if me.IsOk(id) {
		ref = &me[id]
	}
	return
}

func (me FxImage2DLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me FxImage2DLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me FxImage2DLib) Remove(fromID, num int) {
	if l := len(me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onFxImage2DIDsChanged(changed)
	}
}

func (me FxImage2DLib) Walk(on func(ref *FxImage2D)) {
	for id := 0; id < len(me); id++ {
		if me.Ok(id) {
			on(&me[id])
		}
	}
}

//#end-gt
