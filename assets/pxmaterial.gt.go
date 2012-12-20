package assets

//	Defines the physical properties of an object.
type PxMaterialDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Techniques
	HasTechniques
	//	Common-technique profile
	TC struct {
		//	The dynamic friction coefficient.
		DynamicFriction SidFloat
		//	The proportion of the kinetic energy preserved in the impact
		//	(typically ranges from 0.0 to 1.0). Also known as "bounciness" or "elasticity."
		Restitution SidFloat
		//	The static friction coefficient.
		StaticFriction SidFloat
	}
}

//	Initialization
func (me *PxMaterialDef) Init() {
}

//	Lets a shape specify its surface properties using a previously defined physics material.
type PxMaterialInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default and meant to be set ONLY by the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *PxMaterialDef
}

//	Initialization
func (me *PxMaterialInst) Init() {
}

//#begin-gt _definstlib.gt T:PxMaterial

func newPxMaterialDef(id string) (me *PxMaterialDef) {
	me = &PxMaterialDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new PxMaterialInst instance referencing this PxMaterialDef definition.
func (me *PxMaterialDef) NewInst() (inst *PxMaterialInst) {
	inst = &PxMaterialInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is dirty or me.Def is nil, sets me.Def to the correct PxMaterialDef
//	according to the current me.DefRef value (by searching AllPxMaterialDefLibs).
//	Then returns me.Def.
func (me *PxMaterialInst) EnsureDef() *PxMaterialDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.PxMaterialDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibPxMaterialDefs libraries associated by their Id.
	AllPxMaterialDefLibs = LibsPxMaterialDef{}

	//	The "default" LibPxMaterialDefs library for PxMaterialDefs.
	PxMaterialDefs = AllPxMaterialDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllPxMaterialDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibPxMaterialDefs contained in AllPxMaterialDefLibs) for the PxMaterialDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) PxMaterialDef() (def *PxMaterialDef) {
	id := me.S()
	for _, lib := range AllPxMaterialDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllPxMaterialDefLibs variable:
//	a hash-table that contains LibPxMaterialDefs libraries associated by their Id.
type LibsPxMaterialDef map[string]*LibPxMaterialDefs

//	Creates a new LibPxMaterialDefs library with the specified Id, adds it to this LibsPxMaterialDef, and returns it.
//	If this LibsPxMaterialDef already contains a LibPxMaterialDefs library with the specified Id, does nothing and returns nil.
func (me LibsPxMaterialDef) AddNew(id string) (lib *LibPxMaterialDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsPxMaterialDef) new(id string) (lib *LibPxMaterialDefs) {
	lib = newLibPxMaterialDefs(id)
	return
}

//	A library that contains PxMaterialDefs associated by their Id.
//	To create a new LibPxMaterialDefs library, ONLY use the LibsPxMaterialDef.New() or LibsPxMaterialDef.AddNew() methods.
type LibPxMaterialDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*PxMaterialDef
}

func newLibPxMaterialDefs(id string) (me *LibPxMaterialDefs) {
	me = &LibPxMaterialDefs{M: map[string]*PxMaterialDef{}}
	me.Id = id
	return
}

//	Adds the specified PxMaterialDef definition to this LibPxMaterialDefs, and returns it.
//	If this LibPxMaterialDefs already contains a PxMaterialDef definition with the same Id, does nothing and returns nil.
func (me *LibPxMaterialDefs) Add(d *PxMaterialDef) (n *PxMaterialDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new PxMaterialDef definition with the specified Id, adds it to this LibPxMaterialDefs, and returns it.
//	If this LibPxMaterialDefs already contains a PxMaterialDef definition with the specified Id, does nothing and returns nil.
func (me *LibPxMaterialDefs) AddNew(id string) *PxMaterialDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibPxMaterialDefs) Len() int { return len(me.M) }

//	Creates a new PxMaterialDef definition with the specified Id and returns it,
//	but does not add it to this LibPxMaterialDefs.
func (me *LibPxMaterialDefs) New(id string) (def *PxMaterialDef) { def = newPxMaterialDef(id); return }

//	Removes the PxMaterialDef with the specified Id from this LibPxMaterialDefs.
func (me *LibPxMaterialDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibPxMaterialDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibPxMaterialDefs
//	library or its PxMaterialDef definitions. Also called by the global SyncChanges() function.
func (me *LibPxMaterialDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
