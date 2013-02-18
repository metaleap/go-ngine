package core

//	A material binds effects (FxEffect in Core.Libs.Effects) to geometry (Model or Node).
type FxMaterial struct {
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
		for _, thrRend.tmpFaceTag = range face.base.Tags {
			if fxID = me.FaceEffects.ByTag[thrRend.tmpFaceTag]; len(fxID) > 0 {
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

//#begin-gt -gen-lib.gt T:FxMaterial

//	Initializes and returns a new FxMaterial with default parameters.
func NewFxMaterial() (me *FxMaterial) {
	me = &FxMaterial{}
	me.init()
	return
}

//	A hash-table of FxMaterials associated by IDs. Only for use in Core.Libs.
type LibFxMaterials map[string]*FxMaterial

//	Creates and initializes a new FxMaterial with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibFxMaterials) AddNew(id string) (obj *FxMaterial) {
	obj = NewFxMaterial()
	me[id] = obj
	return
}

func (me *LibFxMaterials) ctor() {
	*me = LibFxMaterials{}
}

func (me *LibFxMaterials) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibFxMaterials) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
