package core

//	A material binds effects (FxEffect in Core.Libs.Effects) to geometry (Model or Node).
type FxMaterial struct {
	//	This effect is used by default for all faces that do not fall
	//	under any of the associations in the FaceEffects field.
	DefaultEffectID string

	//	Associates certain individual faces or classes of faces
	//	with specific effect IDs.
	FaceEffects struct {
		//	Associates face classes with effect IDs.
		ByClass map[string]string

		//	Associates specific face IDs with effect IDs.
		ByID map[string]string
	}
}

func (me *FxMaterial) dispose() {
}

func (me *FxMaterial) init() {
	me.FaceEffects.ByID = map[string]string{}
	me.FaceEffects.ByClass = map[string]string{}
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

//#end-gt
