package assets

//	Allows for building complex combinations of rigid bodies and constraints that may be instantiated multiple times.
type PxModelDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Contains zero or more rigid bodies participating in this physics model.
	RigidBodies []*PxRigidBodyDef
	//	Contains zero or more rigid constraints participating in this physics model.
	RigidConstraints []*PxRigidConstraintDef
	//	Child physics models participating in this physics model, with optional property overrides.
	Insts []*PxModelInst
}

//	Initialization
func (me *PxModelDef) Init() {
}

//	Embeds a physics model inside another physics model or instantiates a physics model within a physics scene.
type PxModelInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	Points to the Id of a node in the visual scene. This allows a physics model to be instantiated under a specific transform node, which will dictate the initial position and orientation, and could be animated to influence kinematic rigid bodies. Optional.
	//	By default, the physics model is instantiated under the world, rather than a specific transform node. This parameter is only meaningful when the parent element of the current physics model is a physics scene.
	Parent RefId
	//	Zero or more force fields influencing this physics model.
	ForceFields []*PxForceFieldInst
	//	Contains instances of those rigid bodies included in the instantiated physics model that should have some properties overridden, or should be linked with transform nodes in the visual scene.
	RigidBodies []*PxRigidBodyInst
	//	Contains instances of those rigid constraints included in the instantiated physics model that should have some properties overridden.
	RigidConstraints []*PxRigidConstraintInst
}

//	Initialization
func (me *PxModelInst) Init() {
}

//#begin-gt _definstlib.gt T:PxModel

func newPxModelDef(id string) (me *PxModelDef) {
	me = &PxModelDef{}
	me.Id = id
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
	//	A *map* collection that contains *LibPxModelDefs* libraries associated by their *Id*.
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
//	*LibPxModelDefs* libraries associated by their *Id*.
type LibsPxModelDef map[string]*LibPxModelDefs

//	Creates a new *LibPxModelDefs* library with the specified *Id*, adds it to this *LibsPxModelDef*, and returns it.
//	
//	If this *LibsPxModelDef* already contains a *LibPxModelDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *PxModelDef*s associated by their *Id*. To create a new *LibPxModelDefs* library, ONLY
//	use the *LibsPxModelDef.New()* or *LibsPxModelDef.AddNew()* methods.
type LibPxModelDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*PxModelDef
}

func newLibPxModelDefs(id string) (me *LibPxModelDefs) {
	me = &LibPxModelDefs{M: map[string]*PxModelDef{}}
	me.Id = id
	return
}

//	Adds the specified *PxModelDef* definition to this *LibPxModelDefs*, and returns it.
//	
//	If this *LibPxModelDefs* already contains a *PxModelDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibPxModelDefs) Add(d *PxModelDef) (n *PxModelDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *PxModelDef* definition with the specified *Id*, adds it to this *LibPxModelDefs*, and returns it.
//	
//	If this *LibPxModelDefs* already contains a *PxModelDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibPxModelDefs) AddNew(id string) *PxModelDef { return me.Add(me.New(id)) }

//	Short-hand for len(lib.M)
func (me *LibPxModelDefs) Len() int { return len(me.M) }

//	Creates a new *PxModelDef* definition with the specified *Id* and returns it, but does not add it to this *LibPxModelDefs*.
func (me *LibPxModelDefs) New(id string) (def *PxModelDef) { def = newPxModelDef(id); return }

//	Removes the *PxModelDef* with the specified *Id* from this *LibPxModelDefs*.
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
