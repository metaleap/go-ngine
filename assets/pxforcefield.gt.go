package assets

//	Provides a general container for force fields.
//	Force fields affect physical objects, such as rigid bodies, and
//	may be instantiated under a physics scene or a physics model instance.
type PxForceFieldDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Techniques
	HasTechniques
}

//	Initialization
func (me *PxForceFieldDef) Init() {
}

//	Instantiates a force field resource.
type PxForceFieldInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
}

//	Initialization
func (me *PxForceFieldInst) Init() {
}

//#begin-gt _definstlib.gt T:PxForceField

func newPxForceFieldDef(id string) (me *PxForceFieldDef) {
	me = &PxForceFieldDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

/*
//	Creates and returns a new PxForceFieldInst instance referencing this PxForceFieldDef definition.
func (me *PxForceFieldDef) NewInst(id string) (inst *PxForceFieldInst) {
	inst = &PxForceFieldInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A hash-table that contains LibPxForceFieldDefs libraries associated by their Id.
	AllPxForceFieldDefLibs = LibsPxForceFieldDef{}

	//	The "default" LibPxForceFieldDefs library for PxForceFieldDefs.
	PxForceFieldDefs = AllPxForceFieldDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllPxForceFieldDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global AllPxForceFieldDefLibs variable:
//	a hash-table that contains LibPxForceFieldDefs libraries associated by their Id.
type LibsPxForceFieldDef map[string]*LibPxForceFieldDefs

//	Creates a new LibPxForceFieldDefs library with the specified Id, adds it to this LibsPxForceFieldDef, and returns it.
//	If this LibsPxForceFieldDef already contains a LibPxForceFieldDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains PxForceFieldDefs associated by their Id.
//	To create a new LibPxForceFieldDefs library, ONLY use the LibsPxForceFieldDef.New() or LibsPxForceFieldDef.AddNew() methods.
type LibPxForceFieldDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*PxForceFieldDef
}

func newLibPxForceFieldDefs(id string) (me *LibPxForceFieldDefs) {
	me = &LibPxForceFieldDefs{M: map[string]*PxForceFieldDef{}}
	me.Id = id
	return
}

//	Adds the specified PxForceFieldDef definition to this LibPxForceFieldDefs, and returns it.
//	If this LibPxForceFieldDefs already contains a PxForceFieldDef definition with the same Id, does nothing and returns nil.
func (me *LibPxForceFieldDefs) Add(d *PxForceFieldDef) (n *PxForceFieldDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new PxForceFieldDef definition with the specified Id, adds it to this LibPxForceFieldDefs, and returns it.
//	If this LibPxForceFieldDefs already contains a PxForceFieldDef definition with the specified Id, does nothing and returns nil.
func (me *LibPxForceFieldDefs) AddNew(id string) *PxForceFieldDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibPxForceFieldDefs) Len() int { return len(me.M) }

//	Creates a new PxForceFieldDef definition with the specified Id and returns it,
//	but does not add it to this LibPxForceFieldDefs.
func (me *LibPxForceFieldDefs) New(id string) (def *PxForceFieldDef) { def = newPxForceFieldDef(id); return }

//	Removes the PxForceFieldDef with the specified Id from this LibPxForceFieldDefs.
func (me *LibPxForceFieldDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibPxForceFieldDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibPxForceFieldDefs
//	library or its PxForceFieldDef definitions. Also called by the global SyncChanges() function.
func (me *LibPxForceFieldDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
