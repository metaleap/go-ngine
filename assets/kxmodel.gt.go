package assets

type KxAttachment struct {
	Joint string
	Transforms []*Transform
	Link *KxLink
}

type KxLink struct {
	HasSid
	HasName
	Transforms []*Transform
	Attachments struct {
		Full []*KxAttachment
		Start []*KxAttachment
		End []*KxAttachment
	}
}

type KxModelDef struct {
	BaseDef
	HasTechniques
	TechniqueCommon struct {
		NewParams ParamDefs
		Links []*KxLink
		Formulas struct {
			Defs []*FormulaDef
			Insts []*FormulaInst
		}
	}
}

func (me *KxModelDef) init() {
	me.TechniqueCommon.NewParams = ParamDefs {}
}

type KxModelInst struct {
	BaseInst
	Def *KxModelDef
	Bindings []*KxBind
	NewParams ParamDefs
	SetParams []*ParamInst
}

func (me *KxModelInst) init() {
	me.NewParams = ParamDefs {}
}

//#begin-gt _definstlib.gt T:KxModel

func newKxModelDef(id string) (me *KxModelDef) {
	me = &KxModelDef{}
	me.ID = id
	me.Base.init()
	me.init()
	return
}

/*
//	Creates and returns a new *KxModelInst* instance referencing this *KxModelDef* definition.
func (me *KxModelDef) NewInst(id string) (inst *KxModelInst) {
	inst = &KxModelInst{Def: me}
	inst.init()
	return
}
*/

var (
	//	A *map* collection that contains *LibKxModelDefs* libraries associated by their *ID*.
	AllKxModelDefLibs = LibsKxModelDef{}

	//	The "default" *LibKxModelDefs* library for *KxModelDef*s.
	KxModelDefs = AllKxModelDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllKxModelDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllKxModelDefLibs* variable: a *map* collection that contains
//	*LibKxModelDefs* libraries associated by their *ID*.
type LibsKxModelDef map[string]*LibKxModelDefs

//	Creates a new *LibKxModelDefs* library with the specified *ID*, adds it to this *LibsKxModelDef*, and returns it.
//	
//	If this *LibsKxModelDef* already contains a *LibKxModelDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsKxModelDef) AddNew(id string) (lib *LibKxModelDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsKxModelDef) new(id string) (lib *LibKxModelDefs) {
	lib = newLibKxModelDefs(id)
	return
}

//	A library that contains *KxModelDef*s associated by their *ID*. To create a new *LibKxModelDefs* library, ONLY
//	use the *LibsKxModelDef.New()* or *LibsKxModelDef.AddNew()* methods.
type LibKxModelDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*KxModelDef
}

func newLibKxModelDefs(id string) (me *LibKxModelDefs) {
	me = &LibKxModelDefs{M: map[string]*KxModelDef{}}
	me.ID = id
	return
}

//	Adds the specified *KxModelDef* definition to this *LibKxModelDefs*, and returns it.
//	
//	If this *LibKxModelDefs* already contains a *KxModelDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibKxModelDefs) Add(d *KxModelDef) (n *KxModelDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *KxModelDef* definition with the specified *ID*, adds it to this *LibKxModelDefs*, and returns it.
//	
//	If this *LibKxModelDefs* already contains a *KxModelDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibKxModelDefs) AddNew(id string) *KxModelDef { return me.Add(me.New(id)) }

//	Creates a new *KxModelDef* definition with the specified *ID* and returns it, but does not add it to this *LibKxModelDefs*.
func (me *LibKxModelDefs) New(id string) (def *KxModelDef) { def = newKxModelDef(id); return }

//	Removes the *KxModelDef* with the specified *ID* from this *LibKxModelDefs*.
func (me *LibKxModelDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibKxModelDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibKxModelDefs* library or its *KxModelDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibKxModelDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
