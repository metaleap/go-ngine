package core

type FxMaterial struct {
	EffectID string
}

func NewFxMaterial() (me *FxMaterial) {
	return
}

func (me *FxMaterial) dispose() {
}

//#begin-gt -gen-lib.gt T:FxMaterial

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
