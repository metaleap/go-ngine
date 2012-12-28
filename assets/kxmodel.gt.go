package assets

//	Categorizes the kind of a KxAttachment.
type KxAttachmentKind int

const (
	//	Connects two links, describing a real parent-child dependency between them.
	KxAttachmentKindFull KxAttachmentKind = iota + 1

	//	Connects two links and defines one end of a closed loop.
	KxAttachmentKindStart

	//	Defines one end of the closed loop in an attachment.
	KxAttachmentKindEnd
)

//	Connects links or define ends of closed loops.
type KxAttachment struct {
	//	Must be one of the KxAttachmentKind* enumerated constants.
	Kind KxAttachmentKind

	//	Refers to the KxJoint that connects the parent with the child link. Required.
	Joint RefSid

	//	Zero or more TransformKindRotate and/or TransformKindTranslate transformations.
	Transforms []*Transform

	//	If Kind is KxAttachmentKindFull, specifies the child link in this parent-child dependency.
	Link *KxLink
}

//	Represents a rigid kinematical object without mass whose motion is constrained by one or more joints.
type KxLink struct {
	//	Sid
	HasSid

	//	Name
	HasName

	//	Zero or more TransformKindRotate and/or TransformKindTranslate transformations.
	Transforms []*Transform

	//	The attachments that make up this link.
	Attachments []*KxAttachment
}

//	Categorizes the declaration of kinematical information, containing declarations of
//	joints, links, and attachment points. A kinematics model is focused on strict
//	kinematics description "in zero position", without any additional physical descriptions.
type KxModelDef struct {
	//	Id, Name, Asset, Extras
	BaseDef

	//	Techniques
	HasTechniques

	//	Common-technique profile
	TC struct {
		//	NewParams
		HasParamDefs

		//	The kinematics chain.
		Links []*KxLink

		//	Specifies dependencies among the joints.
		Formulas []Formula
	}
}

//	Initialization
func (me *KxModelDef) Init() {
	me.TC.NewParams = ParamDefs{}
}

//	Instantiates a kinematics model resource.
type KxModelInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst

	//	NewParams
	HasParamDefs

	//	SetParams
	HasParamInsts

	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *KxModelDef

	//	Bindings of inputs to kinematics parameters.
	Bindings []*KxBinding
}

//	Initialization
func (me *KxModelInst) Init() {
	me.NewParams = ParamDefs{}
	me.SetParams = ParamInsts{}
}

//#begin-gt _definstlib.gt T:KxModel

func newKxModelDef(id string) (me *KxModelDef) {
	me = &KxModelDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Returns "the default KxModelInst instance" referencing this KxModelDef definition.
//	That instance is created once when this method is first called on me,
//	and will have its Def field readily set to me.
func (me *KxModelDef) DefaultInst() (inst *KxModelInst) {
	if inst = defaultKxModelInsts[me]; inst == nil {
		inst = me.NewInst()
		defaultKxModelInsts[me] = inst
	}
	return
}

//	Creates and returns a new KxModelInst instance referencing this KxModelDef definition.
//	Any KxModelInst created by this method will have its Def field readily set to me.
func (me *KxModelDef) NewInst() (inst *KxModelInst) {
	inst = &KxModelInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct KxModelDef
//	according to the current me.DefRef value (by searching AllKxModelDefLibs).
//	Then returns me.Def.
//	(Note, every KxModelInst's Def is nil initially, unless it was created via KxModelDef.NewInst().)
func (me *KxModelInst) EnsureDef() *KxModelDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.KxModelDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibKxModelDefs libraries associated by their Id.
	AllKxModelDefLibs = LibsKxModelDef{}

	//	The "default" LibKxModelDefs library for KxModelDefs.
	KxModelDefs = AllKxModelDefLibs.AddNew("")

	defaultKxModelInsts = map[*KxModelDef]*KxModelInst{}
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllKxModelDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibKxModelDefs contained in AllKxModelDefLibs) for the KxModelDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) KxModelDef() (def *KxModelDef) {
	id := me.S()
	for _, lib := range AllKxModelDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllKxModelDefLibs variable:
//	a hash-table that contains LibKxModelDefs libraries associated by their Id.
type LibsKxModelDef map[string]*LibKxModelDefs

//	Creates a new LibKxModelDefs library with the specified Id, adds it to this LibsKxModelDef, and returns it.
//	If this LibsKxModelDef already contains a LibKxModelDefs library with the specified Id, does nothing and returns nil.
func (me LibsKxModelDef) AddNew(id string) (lib *LibKxModelDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsKxModelDef) new(id string) (lib *LibKxModelDefs) {
	lib = newLibKxModelDefs(id)
	return
}

//	A library that contains KxModelDefs associated by their Id.
//	To create a new LibKxModelDefs library, ONLY use the LibsKxModelDef.New() or LibsKxModelDef.AddNew() methods.
type LibKxModelDefs struct {
	//	Id, Name
	BaseLib

	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*KxModelDef
}

func newLibKxModelDefs(id string) (me *LibKxModelDefs) {
	me = &LibKxModelDefs{M: map[string]*KxModelDef{}}
	me.BaseLib.init(id)
	return
}

//	Adds the specified KxModelDef definition to this LibKxModelDefs, and returns it.
//	If this LibKxModelDefs already contains a KxModelDef definition with the same Id, does nothing and returns nil.
func (me *LibKxModelDefs) Add(d *KxModelDef) (n *KxModelDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new KxModelDef definition with the specified Id, adds it to this LibKxModelDefs, and returns it.
//	If this LibKxModelDefs already contains a KxModelDef definition with the specified Id, does nothing and returns nil.
func (me *LibKxModelDefs) AddNew(id string) *KxModelDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibKxModelDefs) Len() int { return len(me.M) }

//	Creates a new KxModelDef definition with the specified Id and returns it,
//	but does not add it to this LibKxModelDefs.
func (me *LibKxModelDefs) New(id string) (def *KxModelDef) { def = newKxModelDef(id); return }

//	Removes the KxModelDef with the specified Id from this LibKxModelDefs.
func (me *LibKxModelDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibKxModelDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibKxModelDefs
//	library or its KxModelDef definitions. Also called by the global SyncChanges() function.
func (me *LibKxModelDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
