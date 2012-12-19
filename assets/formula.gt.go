package assets

import (
	xmlx "github.com/jteeuwen/go-pkg-xmlx"
)

//	Represents either a formula definition or a formula instance.
type Formula struct {
	//	If set, Inst must be nil.
	Def *FormulaDef
	//	If set, Def must be nil.
	Inst *FormulaInst
}

//	There are many ways to describe a formula. Like COLLADA, the *assets* package uses MathML as its common technique.
type FormulaDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sid
	HasSid
	//	Formula parameter definitions
	HasParamDefs
	//	Custom-profile/foreign-technique support
	HasTechniques
	//	A parameter that specifies the result variable of the formula.
	Target ParamFloat
	//	Common-technique profile.
	TC struct {
		//	Any valid MathML (content) XML defining this formula.
		MathML []*xmlx.Node
	}
}

//	Initialization
func (me *FormulaDef) Init() {
	me.NewParams = ParamDefs{}
}

//	Instantiates a formula resource.
type FormulaInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	Specifies the source (for arguments) or the destination (for the result) of the instantiated formula.
	HasParamInsts
}

//	Initialization
func (me *FormulaInst) Init() {
}

//#begin-gt _definstlib.gt T:Formula

func newFormulaDef(id string) (me *FormulaDef) {
	me = &FormulaDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *FormulaInst* instance referencing this *FormulaDef* definition.
func (me *FormulaDef) NewInst(id string) (inst *FormulaInst) {
	inst = &FormulaInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibFormulaDefs* libraries associated by their *Id*.
	AllFormulaDefLibs = LibsFormulaDef{}

	//	The "default" *LibFormulaDefs* library for *FormulaDef*s.
	FormulaDefs = AllFormulaDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllFormulaDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllFormulaDefLibs* variable: a *map* collection that contains
//	*LibFormulaDefs* libraries associated by their *Id*.
type LibsFormulaDef map[string]*LibFormulaDefs

//	Creates a new *LibFormulaDefs* library with the specified *Id*, adds it to this *LibsFormulaDef*, and returns it.
//	
//	If this *LibsFormulaDef* already contains a *LibFormulaDefs* library with the specified *Id*, does nothing and returns *nil*.
func (me LibsFormulaDef) AddNew(id string) (lib *LibFormulaDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsFormulaDef) new(id string) (lib *LibFormulaDefs) {
	lib = newLibFormulaDefs(id)
	return
}

//	A library that contains *FormulaDef*s associated by their *Id*. To create a new *LibFormulaDefs* library, ONLY
//	use the *LibsFormulaDef.New()* or *LibsFormulaDef.AddNew()* methods.
type LibFormulaDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*FormulaDef
}

func newLibFormulaDefs(id string) (me *LibFormulaDefs) {
	me = &LibFormulaDefs{M: map[string]*FormulaDef{}}
	me.Id = id
	return
}

//	Adds the specified *FormulaDef* definition to this *LibFormulaDefs*, and returns it.
//	
//	If this *LibFormulaDefs* already contains a *FormulaDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibFormulaDefs) Add(d *FormulaDef) (n *FormulaDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *FormulaDef* definition with the specified *Id*, adds it to this *LibFormulaDefs*, and returns it.
//	
//	If this *LibFormulaDefs* already contains a *FormulaDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibFormulaDefs) AddNew(id string) *FormulaDef { return me.Add(me.New(id)) }

//	Short-hand for len(lib.M)
func (me *LibFormulaDefs) Len() int { return len(me.M) }

//	Creates a new *FormulaDef* definition with the specified *Id* and returns it, but does not add it to this *LibFormulaDefs*.
func (me *LibFormulaDefs) New(id string) (def *FormulaDef) { def = newFormulaDef(id); return }

//	Removes the *FormulaDef* with the specified *Id* from this *LibFormulaDefs*.
func (me *LibFormulaDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibFormulaDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibFormulaDefs* library or its *FormulaDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibFormulaDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
