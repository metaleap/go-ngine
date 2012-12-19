package assets

const (
	//	Defines a single translational degree of freedom of a joint.
	KX_JOINT_TYPE_PRISMATIC = 1
	//	Defines a single rotational degree of freedom of a joint.
	KX_JOINT_TYPE_REVOLUTE = iota
)

//	Primitive (simple) joints are joints with one degree of freedom (one given axis) and are used to construct more complex joint types (compound joints) that consist of multiple primitives, each representing an axis.
type KxJoint struct {
	//	Sid
	HasSid
	//	Must be one of the KX_JOINT_TYPE_* enumerated constants.
	Type int
	//	Specifies the axis of the degree of freedom.
	Axis struct {
		//	Name
		HasName
		//	Sid, V
		ScopedVec3
	}
	//	If set, these specified limits are physical limits.
	Limits *KxJointLimits
}

//	Declares a primitive/simple joint as fully limited (if Min and Max are both set), partially limited (if either Min or Max is nil, but not both) or unlimited (if Min and Max are nil).
type KxJointLimits struct {
	//	If set, the "minimum" portion of this joint limitation.
	Min *ScopedFloat
	//	If set, the "maximum" portion of this joint limitation.
	Max *ScopedFloat
}

//	Defines a single complex/compound joint with one or more degrees of freedom.
type KxJointDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sid
	HasSid
	//	Primitive (simple) joints are joints with one degree of freedom (one given axis) and are used to construct more complex joint types (compound joints) that consist of multiple primitives, each representing an axis.
	All []*KxJoint
}

//	Initialization
func (me *KxJointDef) Init() {
}

//	Instantiates a kinematics joint resource.
type KxJointInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
}

//	Initialization
func (me *KxJointInst) Init() {
}

//#begin-gt _definstlib.gt T:KxJoint

func newKxJointDef(id string) (me *KxJointDef) {
	me = &KxJointDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *KxJointInst* instance referencing this *KxJointDef* definition.
func (me *KxJointDef) NewInst(id string) (inst *KxJointInst) {
	inst = &KxJointInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibKxJointDefs* libraries associated by their *Id*.
	AllKxJointDefLibs = LibsKxJointDef{}

	//	The "default" *LibKxJointDefs* library for *KxJointDef*s.
	KxJointDefs = AllKxJointDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllKxJointDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllKxJointDefLibs* variable: a *map* collection that contains
//	*LibKxJointDefs* libraries associated by their *Id*.
type LibsKxJointDef map[string]*LibKxJointDefs

//	Creates a new *LibKxJointDefs* library with the specified *Id*, adds it to this *LibsKxJointDef*, and returns it.
//	
//	If this *LibsKxJointDef* already contains a *LibKxJointDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *KxJointDef*s associated by their *Id*. To create a new *LibKxJointDefs* library, ONLY
//	use the *LibsKxJointDef.New()* or *LibsKxJointDef.AddNew()* methods.
type LibKxJointDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*KxJointDef
}

func newLibKxJointDefs(id string) (me *LibKxJointDefs) {
	me = &LibKxJointDefs{M: map[string]*KxJointDef{}}
	me.Id = id
	return
}

//	Adds the specified *KxJointDef* definition to this *LibKxJointDefs*, and returns it.
//	
//	If this *LibKxJointDefs* already contains a *KxJointDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibKxJointDefs) Add(d *KxJointDef) (n *KxJointDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *KxJointDef* definition with the specified *Id*, adds it to this *LibKxJointDefs*, and returns it.
//	
//	If this *LibKxJointDefs* already contains a *KxJointDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibKxJointDefs) AddNew(id string) *KxJointDef { return me.Add(me.New(id)) }

//	Short-hand for len(lib.M)
func (me *LibKxJointDefs) Len() int { return len(me.M) }

//	Creates a new *KxJointDef* definition with the specified *Id* and returns it, but does not add it to this *LibKxJointDefs*.
func (me *LibKxJointDefs) New(id string) (def *KxJointDef) { def = newKxJointDef(id); return }

//	Removes the *KxJointDef* with the specified *Id* from this *LibKxJointDefs*.
func (me *LibKxJointDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibKxJointDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibKxJointDefs* library or its *KxJointDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibKxJointDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
