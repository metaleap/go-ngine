package assets

type PxModelDef struct {
	BaseDef
	RigidBodies      []*PxRigidBodyDef
	RigidConstraints []*PxRigidConstraintDef
	Insts            []*PxModelInst
}

func (me *PxModelDef) Init() {
}

type PxModelInst struct {
	BaseInst
	ParentRef        string
	ForceFields      []*PxForceFieldInst
	RigidBodies      []*PxRigidBodyInst
	RigidConstraints []*PxRigidConstraintInst
}

func (me *PxModelInst) Init() {
}

//#begin-gt _definstlib.gt T:PxModel

func newPxModelDef(id string) (me *PxModelDef) {
	me = &PxModelDef{}
	me.ID = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *PxModelInst* instance referencing this *PxModelDef* definition.
func (me *PxModelDef) NewInst(id string) (inst *PxModelInst) {
	inst = &PxModelInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibPxModelDefs* libraries associated by their *ID*.
	AllPxModelDefLibs = LibsPxModelDef{}

	//	The "default" *LibPxModelDefs* library for *PxModelDef*s.
	PxModelDefs = AllPxModelDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllPxModelDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllPxModelDefLibs* variable: a *map* collection that contains
//	*LibPxModelDefs* libraries associated by their *ID*.
type LibsPxModelDef map[string]*LibPxModelDefs

//	Creates a new *LibPxModelDefs* library with the specified *ID*, adds it to this *LibsPxModelDef*, and returns it.
//	
//	If this *LibsPxModelDef* already contains a *LibPxModelDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsPxModelDef) AddNew(id string) (lib *LibPxModelDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsPxModelDef) new(id string) (lib *LibPxModelDefs) {
	lib = newLibPxModelDefs(id)
	return
}

//	A library that contains *PxModelDef*s associated by their *ID*. To create a new *LibPxModelDefs* library, ONLY
//	use the *LibsPxModelDef.New()* or *LibsPxModelDef.AddNew()* methods.
type LibPxModelDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*PxModelDef
}

func newLibPxModelDefs(id string) (me *LibPxModelDefs) {
	me = &LibPxModelDefs{M: map[string]*PxModelDef{}}
	me.ID = id
	return
}

//	Adds the specified *PxModelDef* definition to this *LibPxModelDefs*, and returns it.
//	
//	If this *LibPxModelDefs* already contains a *PxModelDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibPxModelDefs) Add(d *PxModelDef) (n *PxModelDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *PxModelDef* definition with the specified *ID*, adds it to this *LibPxModelDefs*, and returns it.
//	
//	If this *LibPxModelDefs* already contains a *PxModelDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibPxModelDefs) AddNew(id string) *PxModelDef { return me.Add(me.New(id)) }

//	Creates a new *PxModelDef* definition with the specified *ID* and returns it, but does not add it to this *LibPxModelDefs*.
func (me *LibPxModelDefs) New(id string) (def *PxModelDef) { def = newPxModelDef(id); return }

//	Removes the *PxModelDef* with the specified *ID* from this *LibPxModelDefs*.
func (me *LibPxModelDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibPxModelDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibPxModelDefs* library or its *PxModelDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibPxModelDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
