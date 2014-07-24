package core

//	A material binds effects (FxEffect in Core.Libs.Effects) to geometry (Model or Node).
type FxMaterial struct {
	ID int

	//	This effect is used by default for all faces that do not
	//	match any of the selectors in the FaceEffects field.
	DefaultEffectID int

	//	Associates certain individual faces or tags of faces
	//	with specific effect IDs.
	FaceEffects struct {
		//	Associates face tags with effect IDs.
		ByTag map[string]int

		//	Associates specific face IDs with effect IDs.
		ByID map[string]int
	}
}

func (me *FxMaterial) init() {
	me.DefaultEffectID = -1
	me.FaceEffects.ByID = make(map[string]int, Options.Libs.InitialCap)
	me.FaceEffects.ByTag = make(map[string]int, Options.Libs.InitialCap)
}

func (me *FxMaterial) dispose() {
}

func (me *FxMaterial) faceEffect(face *meshRawFace) *FxEffect {
	return Core.Libs.Effects.get(me.faceEffectID(face))
}

func (me *FxMaterial) faceEffectID(face *meshRawFace) (fxID int) {
	var ok bool
	if fxID, ok = me.FaceEffects.ByID[face.base.ID]; (!ok) || fxID < 0 {
		for _, faceTag := range face.base.Tags {
			if fxID, ok = me.FaceEffects.ByTag[faceTag]; ok && fxID >= 0 {
				return
			}
		}
		fxID = me.DefaultEffectID
	}
	return
}

func (me *FxMaterial) HasFaceEffects() bool {
	return len(me.FaceEffects.ByID) > 0 || len(me.FaceEffects.ByTag) > 0
}











//#begin-gt -gen-lib.gt T:FxMaterial L:Core.Libs.Materials

//	Only used for Core.Libs.Materials
type FxMaterialLib []FxMaterial

func (me *FxMaterialLib) AddNew() (id int) {
	id = -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			id = i
			break
		}
	}
	if id == -1 {
		if id = len(*me); id == cap(*me) {
			nu := make(FxMaterialLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, FxMaterial{})
	}
	ref := &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *FxMaterialLib) Compact() {
	var (
		before, after []FxMaterial
		ref           *FxMaterial
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
			me.onFxMaterialIDsChanged(changed)
		}
	}
}

func (me *FxMaterialLib) init() {
	*me = make(FxMaterialLib, 0, Options.Libs.InitialCap)
}

func (me *FxMaterialLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me FxMaterialLib) get(id int) (ref *FxMaterial) {
	if me.IsOk(id) {
		ref = &me[id]
	}
	return
}

func (me FxMaterialLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me FxMaterialLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me FxMaterialLib) Remove(fromID, num int) {
	if l := len(me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onFxMaterialIDsChanged(changed)
	}
}

func (me FxMaterialLib) Walk(on func(ref *FxMaterial)) {
	for id := 0; id < len(me); id++ {
		if me.Ok(id) {
			on(&me[id])
		}
	}
}

//#end-gt
