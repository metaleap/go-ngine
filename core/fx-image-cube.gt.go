package core

import (
	"image"

	ugl "github.com/go3d/go-opengl/util"
)

type FxImageCube struct {
	FxImageBase
	InitFrom [6]FxImageInitFrom

	glTex  ugl.TextureCube
	images [6]image.Image
}

func (me *FxImageCube) init() {
	me.glTex.Init()
	me.FxImageBase.init()
	me.Preprocess.FlipY = false
}

func (me *FxImageCube) dispose() {
	me.Unload()
}

func (me *FxImageCube) GpuDelete() {
	me.FxImageBase.gpuDelete(&me.glTex)
}

func (me *FxImageCube) GpuSync() (err error) {
	if err = me.glTex.PrepFromImages(me.Storage.Gpu.Bgra, me.Storage.Gpu.UintRev, me.images[:]...); err == nil {
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

func (me *FxImageCube) NoAutoMips() {
	me.FxImageBase.noAutoMips(&me.glTex.TextureBase)
}

func (me *FxImageCube) Unload() {
	me.glSynced = false
	for i := 0; i < len(me.images); i++ {
		me.images[i] = nil
	}
}



//#begin-gt -gen-lib.gt T:FxImageCube L:Core.Libs.Images.TexCube

//	Only used for Core.Libs.Images.TexCube
type FxImageCubeLib []FxImageCube

func (me *FxImageCubeLib) AddNew() (id int) {
	id = -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			id = i
			break
		}
	}
	if id == -1 {
		if id = len(*me); id == cap(*me) {
			nu := make(FxImageCubeLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, FxImageCube{})
	}
	ref := &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *FxImageCubeLib) Compact() {
	var (
		before, after []FxImageCube
		ref           *FxImageCube
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
			me.onFxImageCubeIDsChanged(changed)
		}
	}
}

func (me *FxImageCubeLib) init() {
	*me = make(FxImageCubeLib, 0, Options.Libs.InitialCap)
}

func (me *FxImageCubeLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me FxImageCubeLib) get(id int) (ref *FxImageCube) {
	if me.IsOk(id) {
		ref = &me[id]
	}
	return
}

func (me FxImageCubeLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me FxImageCubeLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me FxImageCubeLib) Remove(fromID, num int) {
	if l := len(me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onFxImageCubeIDsChanged(changed)
	}
}

func (me FxImageCubeLib) Walk(on func(ref *FxImageCube)) {
	for id := 0; id < len(me); id++ {
		if me.Ok(id) {
			on(&me[id])
		}
	}
}

//#end-gt
