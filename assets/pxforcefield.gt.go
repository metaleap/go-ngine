package assets

type PxForceFieldDef struct {
	BaseDef
	HasTechniques
}

func (me *PxForceFieldDef) Init() {
}

type PxForceFieldInst struct {
	BaseInst
}

func (me *PxForceFieldInst) Init() {
}

//#begin-gt _definstlib.gt T:PxForceField

func newPxForceFieldDef(id string) (me *PxForceFieldDef) {
	me = &PxForceFieldDef{}
	me.ID = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *PxForceFieldInst* instance referencing this *PxForceFieldDef* definition.
func (me *PxForceFieldDef) NewInst(id string) (inst *PxForceFieldInst) {
	inst = &PxForceFieldInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibPxForceFieldDefs* libraries associated by their *ID*.
	AllPxForceFieldDefLibs = LibsPxForceFieldDef{}

	//	The "default" *LibPxForceFieldDefs* library for *PxForceFieldDef*s.
	PxForceFieldDefs = AllPxForceFieldDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllPxForceFieldDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllPxForceFieldDefLibs* variable: a *map* collection that contains
//	*LibPxForceFieldDefs* libraries associated by their *ID*.
type LibsPxForceFieldDef map[string]*LibPxForceFieldDefs

//	Creates a new *LibPxForceFieldDefs* library with the specified *ID*, adds it to this *LibsPxForceFieldDef*, and returns it.
//	
//	If this *LibsPxForceFieldDef* already contains a *LibPxForceFieldDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsPxForceFieldDef) AddNew(id string) (lib *LibPxForceFieldDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsPxForceFieldDef) new(id string) (lib *LibPxForceFieldDefs) {
	lib = newLibPxForceFieldDefs(id)
	return
}

//	A library that contains *PxForceFieldDef*s associated by their *ID*. To create a new *LibPxForceFieldDefs* library, ONLY
//	use the *LibsPxForceFieldDef.New()* or *LibsPxForceFieldDef.AddNew()* methods.
type LibPxForceFieldDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*PxForceFieldDef
}

func newLibPxForceFieldDefs(id string) (me *LibPxForceFieldDefs) {
	me = &LibPxForceFieldDefs{M: map[string]*PxForceFieldDef{}}
	me.ID = id
	return
}

//	Adds the specified *PxForceFieldDef* definition to this *LibPxForceFieldDefs*, and returns it.
//	
//	If this *LibPxForceFieldDefs* already contains a *PxForceFieldDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibPxForceFieldDefs) Add(d *PxForceFieldDef) (n *PxForceFieldDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *PxForceFieldDef* definition with the specified *ID*, adds it to this *LibPxForceFieldDefs*, and returns it.
//	
//	If this *LibPxForceFieldDefs* already contains a *PxForceFieldDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibPxForceFieldDefs) AddNew(id string) *PxForceFieldDef { return me.Add(me.New(id)) }

//	Creates a new *PxForceFieldDef* definition with the specified *ID* and returns it, but does not add it to this *LibPxForceFieldDefs*.
func (me *LibPxForceFieldDefs) New(id string) (def *PxForceFieldDef) { def = newPxForceFieldDef(id); return }

//	Removes the *PxForceFieldDef* with the specified *ID* from this *LibPxForceFieldDefs*.
func (me *LibPxForceFieldDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibPxForceFieldDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibPxForceFieldDefs* library or its *PxForceFieldDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibPxForceFieldDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
