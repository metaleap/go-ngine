package assets

type FormulaDef struct {
	BaseDef
	HasSid
	HasTechniques
	ParamDefs       ParamDefs
	Target          ParamFloat
	TechniqueCommon struct {
		Data interface{}
	}
}

func (me *FormulaDef) init() {
	me.ParamDefs = ParamDefs{}
}

type FormulaInst struct {
	BaseInst
	Def        *FormulaDef
	ParamInsts []*ParamInst
}

func (me *FormulaInst) init() {
}

//#begin-gt _definstlib.gt T:Formula

func newFormulaDef(id string) (me *FormulaDef) {
	me = &FormulaDef{}
	me.ID = id
	me.init()
	return
}

/*
//	Creates and returns a new *FormulaInst* instance referencing this *FormulaDef* definition.
func (me *FormulaDef) NewInst(id string) (inst *FormulaInst) {
	inst = &FormulaInst{Def: me}
	inst.init()
	return
}
*/

var (
	//	A *map* collection that contains *LibFormulaDefs* libraries associated by their *ID*.
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
//	*LibFormulaDefs* libraries associated by their *ID*.
type LibsFormulaDef map[string]*LibFormulaDefs

//	Creates a new *LibFormulaDefs* library with the specified *ID*, adds it to this *LibsFormulaDef*, and returns it.
//	
//	If this *LibsFormulaDef* already contains a *LibFormulaDefs* library with the specified *ID*, does nothing and returns *nil*.
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

//	A library that contains *FormulaDef*s associated by their *ID*. To create a new *LibFormulaDefs* library, ONLY
//	use the *LibsFormulaDef.New()* or *LibsFormulaDef.AddNew()* methods.
type LibFormulaDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*FormulaDef
}

func newLibFormulaDefs(id string) (me *LibFormulaDefs) {
	me = &LibFormulaDefs{M: map[string]*FormulaDef{}}
	me.ID = id
	return
}

//	Adds the specified *FormulaDef* definition to this *LibFormulaDefs*, and returns it.
//	
//	If this *LibFormulaDefs* already contains a *FormulaDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibFormulaDefs) Add(d *FormulaDef) (n *FormulaDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *FormulaDef* definition with the specified *ID*, adds it to this *LibFormulaDefs*, and returns it.
//	
//	If this *LibFormulaDefs* already contains a *FormulaDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibFormulaDefs) AddNew(id string) *FormulaDef { return me.Add(me.New(id)) }

//	Creates a new *FormulaDef* definition with the specified *ID* and returns it, but does not add it to this *LibFormulaDefs*.
func (me *LibFormulaDefs) New(id string) (def *FormulaDef) { def = newFormulaDef(id); return }

//	Removes the *FormulaDef* with the specified *ID* from this *LibFormulaDefs*.
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
