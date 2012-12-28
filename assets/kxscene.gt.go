package assets

//	Binds a kinematics model to a node. The description of a kinematics model is
//	completely independent of any visual information, but for calculation the position is important.
type KxModelBinding struct {
	//	A reference to a node.
	Node RefId

	//	Refers to the kinematics model being bound.
	//	Only either SidRef or ParamRef, but not both, must be specified.
	Model struct {
		//	If set, ParamRef must be empty.
		//	The Sid path to the kinematics model to bind to the node.
		SidRef RefSid

		//	If set, SidRef must be empty.
		//	The parameter of the kinematics model that is defined in the instantiated kinematics scene.
		ParamRef RefParam
	}
}

//	Binds a joint axis of a kinematics model to a single transformation of a node. By binding a joint axis
//	to a transformation of a node, it is possible to synchronize a kinematics scene with a visual scene.
type KxJointAxisBinding struct {
	//	A reference to a transformation of a node.
	Target RefSid

	//	If set, Value is ignored. Specifies an axis of a kinematics model.
	Axis ParamOrRefSid

	//	Only used if Axis is empty. Specifies a value of the axis.
	Value ParamOrFloat
}

//	Embodies the entire set of kinematics information that can be articulated from a resource.
type KxSceneDef struct {
	//	Id, Name, Asset, Extras
	BaseDef

	//	Zero or more kinematics models participating in this kinematics scene.
	Models []*KxModelInst

	//	Zero or more articulated systems participating in this kinematics scene.
	ArticulatedSystems []*KxArticulatedSystemInst
}

//	Initialization
func (me *KxSceneDef) Init() {
}

//	Instantiates a kinematics scene resource.
type KxSceneInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst

	//	NewParams
	HasParamDefs

	//	SetParams
	HasParamInsts

	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *KxSceneDef

	//	Zero or more bindings of kinematics models to nodes.
	ModelBindings []*KxModelBinding

	//	Zero or more bindings of kinematics models' joint axes to single node transformations.
	JointAxisBindings []*KxJointAxisBinding
}

//	Initialization
func (me *KxSceneInst) Init() {
	me.NewParams = ParamDefs{}
	me.SetParams = ParamInsts{}
}

//#begin-gt _definstlib.gt T:KxScene

func newKxSceneDef(id string) (me *KxSceneDef) {
	me = &KxSceneDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Returns "the default KxSceneInst instance" referencing this KxSceneDef definition.
//	That instance is created once when this method is first called on me,
//	and will have its Def field readily set to me.
func (me *KxSceneDef) DefaultInst() (inst *KxSceneInst) {
	if inst = defaultKxSceneInsts[me]; inst == nil {
		inst = me.NewInst()
		defaultKxSceneInsts[me] = inst
	}
	return
}

//	Creates and returns a new KxSceneInst instance referencing this KxSceneDef definition.
//	Any KxSceneInst created by this method will have its Def field readily set to me.
func (me *KxSceneDef) NewInst() (inst *KxSceneInst) {
	inst = &KxSceneInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct KxSceneDef
//	according to the current me.DefRef value (by searching AllKxSceneDefLibs).
//	Then returns me.Def.
//	(Note, every KxSceneInst's Def is nil initially, unless it was created via KxSceneDef.NewInst().)
func (me *KxSceneInst) EnsureDef() *KxSceneDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.KxSceneDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibKxSceneDefs libraries associated by their Id.
	AllKxSceneDefLibs = LibsKxSceneDef{}

	//	The "default" LibKxSceneDefs library for KxSceneDefs.
	KxSceneDefs = AllKxSceneDefLibs.AddNew("")

	defaultKxSceneInsts = map[*KxSceneDef]*KxSceneInst{}
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllKxSceneDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibKxSceneDefs contained in AllKxSceneDefLibs) for the KxSceneDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) KxSceneDef() (def *KxSceneDef) {
	id := me.S()
	for _, lib := range AllKxSceneDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllKxSceneDefLibs variable:
//	a hash-table that contains LibKxSceneDefs libraries associated by their Id.
type LibsKxSceneDef map[string]*LibKxSceneDefs

//	Creates a new LibKxSceneDefs library with the specified Id, adds it to this LibsKxSceneDef, and returns it.
//	If this LibsKxSceneDef already contains a LibKxSceneDefs library with the specified Id, does nothing and returns nil.
func (me LibsKxSceneDef) AddNew(id string) (lib *LibKxSceneDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsKxSceneDef) new(id string) (lib *LibKxSceneDefs) {
	lib = newLibKxSceneDefs(id)
	return
}

//	A library that contains KxSceneDefs associated by their Id.
//	To create a new LibKxSceneDefs library, ONLY use the LibsKxSceneDef.New() or LibsKxSceneDef.AddNew() methods.
type LibKxSceneDefs struct {
	//	Id, Name
	BaseLib

	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*KxSceneDef
}

func newLibKxSceneDefs(id string) (me *LibKxSceneDefs) {
	me = &LibKxSceneDefs{M: map[string]*KxSceneDef{}}
	me.BaseLib.init(id)
	return
}

//	Adds the specified KxSceneDef definition to this LibKxSceneDefs, and returns it.
//	If this LibKxSceneDefs already contains a KxSceneDef definition with the same Id, does nothing and returns nil.
func (me *LibKxSceneDefs) Add(d *KxSceneDef) (n *KxSceneDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new KxSceneDef definition with the specified Id, adds it to this LibKxSceneDefs, and returns it.
//	If this LibKxSceneDefs already contains a KxSceneDef definition with the specified Id, does nothing and returns nil.
func (me *LibKxSceneDefs) AddNew(id string) *KxSceneDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibKxSceneDefs) Len() int { return len(me.M) }

//	Creates a new KxSceneDef definition with the specified Id and returns it,
//	but does not add it to this LibKxSceneDefs.
func (me *LibKxSceneDefs) New(id string) (def *KxSceneDef) { def = newKxSceneDef(id); return }

//	Removes the KxSceneDef with the specified Id from this LibKxSceneDefs.
func (me *LibKxSceneDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

func (me *LibKxSceneDefs) resolver(part0 string) refSidResolver {
	return me.M[part0]
}

func (me *LibKxSceneDefs) resolverRootIsLib() bool {
	return true
}

//	Signals to the core package (or your custom package) that changes have been made to this LibKxSceneDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibKxSceneDefs
//	library or its KxSceneDef definitions. Also called by the global SyncChanges() function.
func (me *LibKxSceneDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
