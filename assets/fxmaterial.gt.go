package assets

type FxMaterialDef struct {
	BaseDef
	Effect FxEffectInst
}

func (me *FxMaterialDef) Init() {
}

type FxMaterialInst struct {
	BaseInst
	Symbol           string
	Binds            []*FxMaterialInstBind
	BindVertexInputs []*FxMaterialInstBindVertexInput
}

func (me *FxMaterialInst) init() {
}

type FxMaterialInstBind struct {
	Semantic string
	Target   string
}

type FxMaterialInstBindVertexInput struct {
	Semantic      string
	InputSemantic string
	InputSet      *uint64
}

//#begin-gt _definstlib.gt T:FxMaterial

func newFxMaterialDef(id string) (me *FxMaterialDef) {
	me = &FxMaterialDef{}
	me.ID = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *FxMaterialInst* instance referencing this *FxMaterialDef* definition.
func (me *FxMaterialDef) NewInst(id string) (inst *FxMaterialInst) {
	inst = &FxMaterialInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibFxMaterialDefs* libraries associated by their *ID*.
	AllFxMaterialDefLibs = LibsFxMaterialDef{}

	//	The "default" *LibFxMaterialDefs* library for *FxMaterialDef*s.
	FxMaterialDefs = AllFxMaterialDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllFxMaterialDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllFxMaterialDefLibs* variable: a *map* collection that contains
//	*LibFxMaterialDefs* libraries associated by their *ID*.
type LibsFxMaterialDef map[string]*LibFxMaterialDefs

//	Creates a new *LibFxMaterialDefs* library with the specified *ID*, adds it to this *LibsFxMaterialDef*, and returns it.
//	
//	If this *LibsFxMaterialDef* already contains a *LibFxMaterialDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsFxMaterialDef) AddNew(id string) (lib *LibFxMaterialDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsFxMaterialDef) new(id string) (lib *LibFxMaterialDefs) {
	lib = newLibFxMaterialDefs(id)
	return
}

//	A library that contains *FxMaterialDef*s associated by their *ID*. To create a new *LibFxMaterialDefs* library, ONLY
//	use the *LibsFxMaterialDef.New()* or *LibsFxMaterialDef.AddNew()* methods.
type LibFxMaterialDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*FxMaterialDef
}

func newLibFxMaterialDefs(id string) (me *LibFxMaterialDefs) {
	me = &LibFxMaterialDefs{M: map[string]*FxMaterialDef{}}
	me.ID = id
	return
}

//	Adds the specified *FxMaterialDef* definition to this *LibFxMaterialDefs*, and returns it.
//	
//	If this *LibFxMaterialDefs* already contains a *FxMaterialDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibFxMaterialDefs) Add(d *FxMaterialDef) (n *FxMaterialDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *FxMaterialDef* definition with the specified *ID*, adds it to this *LibFxMaterialDefs*, and returns it.
//	
//	If this *LibFxMaterialDefs* already contains a *FxMaterialDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibFxMaterialDefs) AddNew(id string) *FxMaterialDef { return me.Add(me.New(id)) }

//	Creates a new *FxMaterialDef* definition with the specified *ID* and returns it, but does not add it to this *LibFxMaterialDefs*.
func (me *LibFxMaterialDefs) New(id string) (def *FxMaterialDef) { def = newFxMaterialDef(id); return }

//	Removes the *FxMaterialDef* with the specified *ID* from this *LibFxMaterialDefs*.
func (me *LibFxMaterialDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibFxMaterialDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibFxMaterialDefs* library or its *FxMaterialDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibFxMaterialDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
