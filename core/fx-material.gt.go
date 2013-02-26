package core

//	A material binds effects (FxEffect in Core.Libs.Effects) to geometry (Model or Node).
type FxMaterial struct {
	ID int

	//	This effect is used by default for all faces that do not
	//	match any of the selectors in the FaceEffects field.
	DefaultEffectID string

	//	Associates certain individual faces or tags of faces
	//	with specific effect IDs.
	FaceEffects struct {
		//	Associates face tags with effect IDs.
		ByTag map[string]string

		//	Associates specific face IDs with effect IDs.
		ByID map[string]string
	}
}

func (me *FxMaterial) dispose() {
}

func (me *FxMaterial) faceEffect(face *meshRawFace) *FxEffect {
	return Core.Libs.Effects[me.faceEffectID(face)]
}

func (me *FxMaterial) faceEffectID(face *meshRawFace) (fxID string) {
	if fxID = me.FaceEffects.ByID[face.base.ID]; len(fxID) == 0 {
		for _, faceTag := range face.base.Tags {
			if fxID = me.FaceEffects.ByTag[faceTag]; len(fxID) > 0 {
				return
			}
		}
		fxID = me.DefaultEffectID
	}
	return
}

func (me *FxMaterial) init() {
	me.FaceEffects.ByID = map[string]string{}
	me.FaceEffects.ByTag = map[string]string{}
}

func (me *FxMaterial) HasFaceEffects() bool {
	return len(me.FaceEffects.ByID) > 0 || len(me.FaceEffects.ByTag) > 0
}

//#begin-gt -gen-lib2.gt T:FxMaterial L:Materials

//	Only used for Core.Libs.Materials.
type FxMaterialLib []FxMaterial

func (_ FxMaterialLib) AddNew() (ref *FxMaterial) {
	me, id := &Core.Libs.Materials, -1
	for i, _ := range *me {
		if (*me)[i].ID < 0 {
			id = i
			break
		}
	}
	if id < 0 {
		if id = len(*me); id == cap(*me) {
			nu := make(FxMaterialLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, FxMaterial{})
	}
	ref = &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (_ FxMaterialLib) Compact() {
	var (
		before, after []FxMaterial
		ref           *FxMaterial
		oldID         int
	)
	me := &Core.Libs.Materials
	for i, _ := range *me {
		if (*me)[i].ID < 0 {
			before, after = (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	changed := make(map[int]int, len(*me))
	for i, _ := range *me {
		if (*me)[i].ID != i {
			ref = &(*me)[i]
			oldID, ref.ID = ref.ID, i
			changed[oldID] = i
		}
	}
	if len(changed) > 0 {
		Core.Libs.Materials.onFxMaterialIDsChanged(changed)
		Options.Libs.OnIDsChanged.Materials(changed)
	}
}

func (_ FxMaterialLib) ctor() {
	me := &Core.Libs.Materials
	*me = make(FxMaterialLib, 0, Options.Libs.InitialCap)
}

func (_ FxMaterialLib) dispose() {
	me := &Core.Libs.Materials
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (_ FxMaterialLib) Get(id int) (ref *FxMaterial) {
	if id >= 0 && id < len(Core.Libs.Materials) {
		if ref = &Core.Libs.Materials[id]; ref.ID != id {
			ref = nil
		}
	}
	return
}

func (_ FxMaterialLib) Has(id int) (has bool) {
	if id >= 0 && id < len(Core.Libs.Materials) {
		has = Core.Libs.Materials[id].ID == id
	}
	return
}

func (_ FxMaterialLib) Remove(fromID, num int) {
	me := &Core.Libs.Materials
	if l := len(*me); fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			(*me)[id].dispose()
			changed[id], (*me)[id].ID = -1, -1
		}
		Core.Libs.Materials.onFxMaterialIDsChanged(changed)
		Options.Libs.OnIDsChanged.Materials(changed)
	}
}

//#end-gt
