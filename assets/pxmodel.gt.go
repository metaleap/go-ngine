package assets

//	Allows for building complex combinations of rigid bodies and constraints
//	that may be instantiated multiple times.
type PxModelDef struct {
	//	Id, Name, Asset, Extras
	BaseDef

	//	Contains zero or more rigid bodies participating in this physics model.
	RigidBodies PxRigidBodyDefs

	//	Contains zero or more rigid constraints participating in this physics model.
	RigidConstraints PxRigidConstraintDefs

	//	Child physics models participating in this physics model, with optional property overrides.
	Insts []*PxModelInst
}

//	Initialization
func (me *PxModelDef) Init() {
	me.RigidBodies = PxRigidBodyDefs{}
	me.RigidConstraints = PxRigidConstraintDefs{}
}

//	Embeds a physics model inside another physics model or
//	instantiates a physics model within a physics scene.
type PxModelInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst

	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *PxModelDef

	//	Points to the Id of a node in the visual scene. This allows a physics model to be instantiated
	//	under a specific transform node, which will dictate the initial position and orientation,
	//	and could be animated to influence kinematic rigid bodies. Optional.
	//	By default, the physics model is instantiated under the world, rather than a specific transform node.
	//	This parameter is only meaningful when the parent element of the current physics model is a physics scene.
	Parent RefId

	//	Zero or more force fields influencing this physics model.
	ForceFields []*PxForceFieldInst

	//	Contains instances of those rigid bodies included in the instantiated physics model that should
	//	have some properties overridden, or should be linked with transform nodes in the visual scene.
	RigidBodies []*PxRigidBodyInst

	//	Contains instances of those rigid constraints included in the instantiated
	//	physics model that should have some properties overridden.
	RigidConstraints []*PxRigidConstraintInst
}

//	Initialization
func (me *PxModelInst) Init() {
}

//#begin-gt _definstlib.gt T:PxModel

func newPxModelDef(id string) (me *PxModelDef) {
	me = &PxModelDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Returns "the default PxModelInst instance" referencing this PxModelDef definition.
//	That instance is created once when this method is first called on me,
//	and will have its Def field readily set to me.
func (me *PxModelDef) DefaultInst() (inst *PxModelInst) {
	if inst = defaultPxModelInsts[me]; inst == nil {
		inst = me.NewInst()
		defaultPxModelInsts[me] = inst
	}
	return
}

//	Creates and returns a new PxModelInst instance referencing this PxModelDef definition.
//	Any PxModelInst created by this method will have its Def field readily set to me.
func (me *PxModelDef) NewInst() (inst *PxModelInst) {
	inst = &PxModelInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct PxModelDef
//	according to the current me.DefRef value (by searching AllPxModelDefLibs).
//	Then returns me.Def.
//	(Note, every PxModelInst's Def is nil initially, unless it was created via PxModelDef.NewInst().)
func (me *PxModelInst) EnsureDef() *PxModelDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.PxModelDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibPxModelDefs libraries associated by their Id.
	AllPxModelDefLibs = LibsPxModelDef{}

	//	The "default" LibPxModelDefs library for PxModelDefs.
	PxModelDefs = AllPxModelDefLibs.AddNew("")

	defaultPxModelInsts = map[*PxModelDef]*PxModelInst{}
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllPxModelDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibPxModelDefs contained in AllPxModelDefLibs) for the PxModelDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) PxModelDef() (def *PxModelDef) {
	id := me.S()
	for _, lib := range AllPxModelDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllPxModelDefLibs variable:
//	a hash-table that contains LibPxModelDefs libraries associated by their Id.
type LibsPxModelDef map[string]*LibPxModelDefs

//	Creates a new LibPxModelDefs library with the specified Id, adds it to this LibsPxModelDef, and returns it.
//	If this LibsPxModelDef already contains a LibPxModelDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains PxModelDefs associated by their Id.
//	To create a new LibPxModelDefs library, ONLY use the LibsPxModelDef.New() or LibsPxModelDef.AddNew() methods.
type LibPxModelDefs struct {
	//	Id, Name
	BaseLib

	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*PxModelDef
}

func newLibPxModelDefs(id string) (me *LibPxModelDefs) {
	me = &LibPxModelDefs{M: map[string]*PxModelDef{}}
	me.BaseLib.init(id)
	return
}

//	Adds the specified PxModelDef definition to this LibPxModelDefs, and returns it.
//	If this LibPxModelDefs already contains a PxModelDef definition with the same Id, does nothing and returns nil.
func (me *LibPxModelDefs) Add(d *PxModelDef) (n *PxModelDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new PxModelDef definition with the specified Id, adds it to this LibPxModelDefs, and returns it.
//	If this LibPxModelDefs already contains a PxModelDef definition with the specified Id, does nothing and returns nil.
func (me *LibPxModelDefs) AddNew(id string) *PxModelDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibPxModelDefs) Len() int { return len(me.M) }

//	Creates a new PxModelDef definition with the specified Id and returns it,
//	but does not add it to this LibPxModelDefs.
func (me *LibPxModelDefs) New(id string) (def *PxModelDef) { def = newPxModelDef(id); return }

//	Removes the PxModelDef with the specified Id from this LibPxModelDefs.
func (me *LibPxModelDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

func (me *LibPxModelDefs) resolver(part0 string) refSidResolver {
	return me.M[part0]
}

func (me *LibPxModelDefs) resolverRootIsLib() bool {
	return true
}

//	Signals to the core package (or your custom package) that changes have been made to this LibPxModelDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibPxModelDefs
//	library or its PxModelDef definitions. Also called by the global SyncChanges() function.
func (me *LibPxModelDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
