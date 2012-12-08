package assets

//	Defines the equations necessary for the visual appearance of geometry and screen-space image processing.
type EffectDef struct {
	BaseDef
	Params FxParamDefs
}

func (me *EffectDef) init() {
	me.Params = FxParamDefs{}
}

//	An instance referencing a Effect definition.
type EffectInst struct {
	BaseInst

	//	The Effect definition referenced by this instance.
	Def *EffectDef
}

func (me *EffectInst) init() {
}

//#begin-gt _definstlib.gt T:Effect

func newEffectDef(id string) (me *EffectDef) {
	me = &EffectDef{}
	me.BaseDef.init(id)
	me.init()
	return
}

//	Creates and returns a new *EffectInst* instance referencing this *EffectDef* definition.
func (me *EffectDef) NewInst(id string) (inst *EffectInst) {
	inst = &EffectInst{Def: me}
	inst.Base.init(id)
	inst.init()
	return
}

var (
	//	A *map* collection that contains *LibEffectDefs* libraries associated by their *ID*.
	AllEffectDefLibs = LibsEffectDef{}

	//	The "default" *LibEffectDefs* library for *EffectDef*s.
	EffectDefs = AllEffectDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllEffectDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllEffectDefLibs* variable: a *map* collection that contains
//	*LibEffectDefs* libraries associated by their *ID*.
type LibsEffectDef map[string]*LibEffectDefs

//	Creates a new *LibEffectDefs* library with the specified *ID*, adds it to this *LibsEffectDef*, and returns it.
//	
//	If this *LibsEffectDef* already contains a *LibEffectDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsEffectDef) AddNew(id string) (lib *LibEffectDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsEffectDef) new(id string) (lib *LibEffectDefs) {
	lib = newLibEffectDefs(id)
	return
}

//	A library that contains *EffectDef*s associated by their *ID*. To create a new *LibEffectDefs* library, ONLY
//	use the *LibsEffectDef.New()* or *LibsEffectDef.AddNew()* methods.
type LibEffectDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*EffectDef
}

func newLibEffectDefs(id string) (me *LibEffectDefs) {
	me = &LibEffectDefs{M: map[string]*EffectDef{}}
	me.Base.init(id)
	return
}

//	Adds the specified *EffectDef* definition to this *LibEffectDefs*, and returns it.
//	
//	If this *LibEffectDefs* already contains a *EffectDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibEffectDefs) Add(d *EffectDef) (n *EffectDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *EffectDef* definition with the specified *ID*, adds it to this *LibEffectDefs*, and returns it.
//	
//	If this *LibEffectDefs* already contains a *EffectDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibEffectDefs) AddNew(id string) *EffectDef { return me.Add(me.New(id)) }

//	Creates a new *EffectDef* definition with the specified *ID* and returns it, but does not add it to this *LibEffectDefs*.
func (me *LibEffectDefs) New(id string) (def *EffectDef) { def = newEffectDef(id); return }

//	Removes the *EffectDef* with the specified *ID* from this *LibEffectDefs*.
func (me *LibEffectDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibEffectDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibEffectDefs* library or its *EffectDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibEffectDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
