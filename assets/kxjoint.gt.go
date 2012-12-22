package assets

//	Categorizes the kind of a KxJoint.
type KxJointKind int

const (
	_ = iota
	//	Defines a single translational degree of freedom of a joint.
	KxJointKindPrismatic KxJointKind = iota
	//	Defines a single rotational degree of freedom of a joint.
	KxJointKindRevolute
)

//	Primitive (simple) joints are joints with one degree of freedom (one given axis) and
//	are used to construct more complex joint types (compound joints) that consist of
//	multiple primitives, each representing an axis.
type KxJoint struct {
	//	Sid
	HasSid
	//	Must be one of the KxJointKind* enumerated constants.
	Kind KxJointKind
	//	Specifies the axis of the degree of freedom.
	Axis struct {
		//	Name
		HasName
		//	Sid, V
		SidVec3
	}
	//	If set, these specified limits are physical limits.
	Limits *KxJointLimits
}

//	Declares a primitive/simple joint as fully limited (if Min and Max are both set),
//	partially limited (if either Min or Max is nil, but not both) or unlimited (if Min and Max are nil).
type KxJointLimits struct {
	//	If set, the "minimum" portion of this joint limitation.
	Min *SidFloat
	//	If set, the "maximum" portion of this joint limitation.
	Max *SidFloat
}

//	Defines a single complex/compound joint with one or more degrees of freedom.
type KxJointDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sid
	HasSid
	//	Primitive (simple) joints are joints with one degree of freedom (one given axis) and are
	//	used to construct more complex joint types (compound joints) that
	//	consist of multiple primitives, each representing an axis.
	All []*KxJoint
}

//	Initialization
func (me *KxJointDef) Init() {
}

//	Instantiates a kinematics joint resource.
type KxJointInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *KxJointDef
}

//	Initialization
func (me *KxJointInst) Init() {
}

//#begin-gt _definstlib.gt T:KxJoint

func newKxJointDef(id string) (me *KxJointDef) {
	me = &KxJointDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new KxJointInst instance referencing this KxJointDef definition.
//	Any KxJointInst created by this method will have its Def field readily set to me.
func (me *KxJointDef) NewInst() (inst *KxJointInst) {
	inst = &KxJointInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct KxJointDef
//	according to the current me.DefRef value (by searching AllKxJointDefLibs).
//	Then returns me.Def.
//	(Note, every KxJointInst's Def is nil initially, unless it was created via KxJointDef.NewInst().)
func (me *KxJointInst) EnsureDef() *KxJointDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.KxJointDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibKxJointDefs libraries associated by their Id.
	AllKxJointDefLibs = LibsKxJointDef{}

	//	The "default" LibKxJointDefs library for KxJointDefs.
	KxJointDefs = AllKxJointDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllKxJointDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibKxJointDefs contained in AllKxJointDefLibs) for the KxJointDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) KxJointDef() (def *KxJointDef) {
	id := me.S()
	for _, lib := range AllKxJointDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllKxJointDefLibs variable:
//	a hash-table that contains LibKxJointDefs libraries associated by their Id.
type LibsKxJointDef map[string]*LibKxJointDefs

//	Creates a new LibKxJointDefs library with the specified Id, adds it to this LibsKxJointDef, and returns it.
//	If this LibsKxJointDef already contains a LibKxJointDefs library with the specified Id, does nothing and returns nil.
func (me LibsKxJointDef) AddNew(id string) (lib *LibKxJointDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsKxJointDef) new(id string) (lib *LibKxJointDefs) {
	lib = newLibKxJointDefs(id)
	return
}

//	A library that contains KxJointDefs associated by their Id.
//	To create a new LibKxJointDefs library, ONLY use the LibsKxJointDef.New() or LibsKxJointDef.AddNew() methods.
type LibKxJointDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*KxJointDef
}

func newLibKxJointDefs(id string) (me *LibKxJointDefs) {
	me = &LibKxJointDefs{M: map[string]*KxJointDef{}}
	me.Id = id
	return
}

//	Adds the specified KxJointDef definition to this LibKxJointDefs, and returns it.
//	If this LibKxJointDefs already contains a KxJointDef definition with the same Id, does nothing and returns nil.
func (me *LibKxJointDefs) Add(d *KxJointDef) (n *KxJointDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new KxJointDef definition with the specified Id, adds it to this LibKxJointDefs, and returns it.
//	If this LibKxJointDefs already contains a KxJointDef definition with the specified Id, does nothing and returns nil.
func (me *LibKxJointDefs) AddNew(id string) *KxJointDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibKxJointDefs) Len() int { return len(me.M) }

//	Creates a new KxJointDef definition with the specified Id and returns it,
//	but does not add it to this LibKxJointDefs.
func (me *LibKxJointDefs) New(id string) (def *KxJointDef) { def = newKxJointDef(id); return }

//	Removes the KxJointDef with the specified Id from this LibKxJointDefs.
func (me *LibKxJointDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

func (me *LibKxJointDefs) resolver(part0 string) refSidResolver {
	return me.M[part0]
}

//	Signals to the core package (or your custom package) that changes have been made to this LibKxJointDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibKxJointDefs
//	library or its KxJointDef definitions. Also called by the global SyncChanges() function.
func (me *LibKxJointDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
