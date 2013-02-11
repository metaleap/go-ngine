package core

type FxEffectShader struct {
	shaderIDs []string
	pname     string
}

func NewFxEffectShader(shaderIDs ...string) (me *FxEffectShader) {
	me = &FxEffectShader{shaderIDs: shaderIDs[:]}
	for _, shid := range shaderIDs {
		me.pname += ("_" + shid)
	}
	return
}

type FxShader struct {
	FuncName string
}

func NewFxShader(funcName string) (me *FxShader) {
	me = &FxShader{FuncName: "fx_" + funcName}
	return
}

type FxEffect struct {
	OldDiffuse   *FxColorOrTexture
	EffectShader *FxEffectShader
}

func (me *FxEffect) dispose() {
}

func (me *FxEffect) init() {
}

//#begin-gt -gen-lib.gt T:FxEffect

//	Initializes and returns a new FxEffect with default parameters.
func NewFxEffect() (me *FxEffect) {
	me = &FxEffect{}
	me.init()
	return
}

//	A hash-table of FxEffects associated by IDs. Only for use in Core.Libs.
type LibFxEffects map[string]*FxEffect

//	Creates and initializes a new FxEffect with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibFxEffects) AddNew(id string) (obj *FxEffect) {
	obj = NewFxEffect()
	me[id] = obj
	return
}

func (me *LibFxEffects) ctor() {
	*me = LibFxEffects{}
}

func (me *LibFxEffects) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibFxEffects) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
