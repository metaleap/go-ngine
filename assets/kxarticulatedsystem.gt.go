package assets

//	Specifies the parent axis' index in the jointmap.
type KxAxisIndex struct {
	//	If set, specifies the special use of this index.
	Semantic string

	//	If not set, the parent axis will not appear in the jointmap.
	I ParamOrInt
}

//	Specifies the parent axis' soft limits.
type KxAxisLimits struct {
	//	The "minimum" portion of this limits descriptor.
	Min ParamOrFloat

	//	The "maximum" portion of this limits descriptor.
	Max ParamOrFloat
}

//	Binds inputs to kinematics parameters upon instantiation.
type KxBinding struct {
	//	The identifier of the parameter to bind to the new symbol name. Required.
	Symbol string

	//	If set, Value is ignored.
	Param RefParam

	//	Only used if Param is empty.
	Value interface{}
}

//	Specifies additional dynamics information for an effector.
type KxEffector struct {
	//	Sid
	HasSid

	//	Name
	HasName

	//	NewParams
	HasParamDefs

	//	SetParams
	HasParamInsts

	//	Bindings of inputs to kinematics parameters.
	Bindings []*KxBinding

	//	Specifies maximum speed.
	//	The first value is translational (m/sec), the second is rotational (°/sec).
	Speed *ParamOrFloat2

	//	Specifies maximum acceleration.
	//	The first value is translational (m/sec²), the second is rotational (°/sec²).
	Acceleration *ParamOrFloat2

	//	Specifies the maximum deceleration.
	//	The first value is translational (m/sec²), the second is rotational (°/sec²).
	Deceleration *ParamOrFloat2

	//	Specifies the maximum jerk (also called jolt or surge).
	//	The first value is translational (m/sec³), the second is rotational (°/sec³).
	Jerk *ParamOrFloat2
}

//	Constructor
func NewKxEffector() (me *KxEffector) {
	me = &KxEffector{}
	me.NewParams = ParamDefs{}
	me.SetParams = ParamInsts{}
	return
}

//	Contains information for a frame used for kinematics calculation.
type KxFrame struct {
	//	References a KxLink defined in the kinematics model. Optional.
	Link RefSid

	//	Zero or more TransformKindTranslate and/or TransformKindRotate transformations.
	Transforms []*Transform
}

//	Defines the offset frame from the KxFrameOrigin;
//	this offset usually represents the transformation to a work piece.
type KxFrameObject struct {
	//	Link, Transforms
	KxFrame
}

//	Defines the base frame for kinematics calculation.
type KxFrameOrigin struct {
	//	Link, Transforms
	KxFrame
}

//	Defines the offset frame from the KxFrameTip,
//	which usually represents the work point of the end effector (for example, a welding gun).
type KxFrameTcp struct {
	//	Link, Transforms
	KxFrame
}

//	Defines the frame at the end of the kinematics chain.
type KxFrameTip struct {
	//	Link, Transforms
	KxFrame
}

//	Contains axis information to describe the kinematics behavior of an articulated model.
type KxKinematicsAxis struct {
	//	Sid
	HasSid

	//	Name
	HasName

	//	NewParams
	HasParamDefs

	//	The joint axis of an instantiated kinematics model.
	Axis RefSid

	//	Defaults to true.
	Active ParamOrBool

	//	Specifies this axis' indices in the jointmap. If empty, this axis will not appear in the jointmap.
	Indices []*KxAxisIndex

	//	Specifies the soft limits. If not set, the axis is limited only by its physical limits.
	Limits *KxAxisLimits

	//	Defaults to false.
	Locked ParamOrBool

	//	Formulas can be useful to define the behavior of a passive link according to one or more
	//	active axes, or to define dependencies of the soft limits and another joint, for example.
	Formulas []Formula
}

//	Constructor
func NewKxKinematicsAxis() (me *KxKinematicsAxis) {
	me = &KxKinematicsAxis{}
	me.Active.B, me.NewParams = true, ParamDefs{}
	return
}

//	Contains additional information to describe the kinematical behavior of an articulated model.
type KxKinematicsSystem struct {
	//	Techniques
	HasTechniques

	//	The kinematics models to be enhanced with kinematics information.
	Models []*KxModelInst

	//	Common-technique profile
	TC struct {
		//	Kinematics-related information for all axes.
		AxisInfos []*KxKinematicsAxis

		//	Kinematics calculation chain frames
		Frame struct {
			//	Defines the base frame for kinematics calculation.
			Origin KxFrameOrigin

			//	Defines the frame at the end of the kinematics chain.
			Tip KxFrameTip

			//	If set, defines the offset frame from the Tip frame,
			//	which usually represents the work point of the end effector (for example, a welding gun).
			Tcp *KxFrameTcp

			//	If set, defines the offset frame from the Origin frame;
			//	this offset usually represents the transformation to a work piece.
			Object *KxFrameObject
		}
	}
}

//	Contains axis information to describe the motion behavior of an articulated model.
type KxMotionAxis struct {
	//	Sid
	HasSid

	//	Name
	HasName

	//	NewParams
	HasParamDefs

	//	SetParams
	HasParamInsts

	//	References the KxKinematicsAxis of an instantiated kinematics system.
	Axis RefSid

	//	Bindings of inputs to kinematics parameters.
	Bindings []*KxBinding

	//	The maximum permitted speed of the axis in meters per second (m/sec).
	Speed *ParamOrFloat

	//	The maximum permitted acceleration of the axis in m/sec².
	Acceleration *ParamOrFloat

	//	The maximum permitted deceleration of an axis.
	//	If not set, acceleration and deceleration have the same value in m/sec².
	Deceleration *ParamOrFloat

	//	The maximum permitted jerk of an axis in m/sec³.
	Jerk *ParamOrFloat
}

//	Constructor
func NewKxMotionAxis() (me *KxMotionAxis) {
	me = &KxMotionAxis{}
	me.NewParams = ParamDefs{}
	me.SetParams = ParamInsts{}
	return
}

//	Contains additional information to describe the dynamics behaviour of an articulated model.
type KxMotionSystem struct {
	//	Techniques
	HasTechniques

	//	The articulated system to be enhanced with dynamics information.
	ArticulatedSystem *KxArticulatedSystemInst

	//	Common-technique profile
	TC struct {
		//	Dynamics-related information for all axes.
		AxisInfos []*KxMotionAxis

		//	Additional dynamics information
		EffectorInfo *KxEffector
	}
}

//	Categorizes the declaration of generic control information for kinematics systems.
type KxArticulatedSystemDef struct {
	//	Id, Name, Asset, Extras
	BaseDef

	//	If set, Motion must be nil, and this articulated system describes a kinematics system.
	Kinematics *KxKinematicsSystem

	//	If set, Kinematics must be nil, and this articulated system describes a motion system.
	Motion *KxMotionSystem
}

//	Initialization
func (me *KxArticulatedSystemDef) Init() {
}

//	Instantiates a kinematics articulated system resource.
type KxArticulatedSystemInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst

	//	NewParams
	HasParamDefs

	//	SetParams
	HasParamInsts

	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *KxArticulatedSystemDef

	//	Bindings of inputs to kinematics parameters.
	Bindings []*KxBinding
}

//	Initialization
func (me *KxArticulatedSystemInst) Init() {
	me.NewParams = ParamDefs{}
	me.SetParams = ParamInsts{}
}

//#begin-gt _definstlib.gt T:KxArticulatedSystem

func newKxArticulatedSystemDef(id string) (me *KxArticulatedSystemDef) {
	me = &KxArticulatedSystemDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Returns "the default KxArticulatedSystemInst instance" referencing this KxArticulatedSystemDef definition.
//	That instance is created once when this method is first called on me,
//	and will have its Def field readily set to me.
func (me *KxArticulatedSystemDef) DefaultInst() (inst *KxArticulatedSystemInst) {
	if inst = defaultKxArticulatedSystemInsts[me]; inst == nil {
		inst = me.NewInst()
		defaultKxArticulatedSystemInsts[me] = inst
	}
	return
}

//	Creates and returns a new KxArticulatedSystemInst instance referencing this KxArticulatedSystemDef definition.
//	Any KxArticulatedSystemInst created by this method will have its Def field readily set to me.
func (me *KxArticulatedSystemDef) NewInst() (inst *KxArticulatedSystemInst) {
	inst = &KxArticulatedSystemInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct KxArticulatedSystemDef
//	according to the current me.DefRef value (by searching AllKxArticulatedSystemDefLibs).
//	Then returns me.Def.
//	(Note, every KxArticulatedSystemInst's Def is nil initially, unless it was created via KxArticulatedSystemDef.NewInst().)
func (me *KxArticulatedSystemInst) EnsureDef() *KxArticulatedSystemDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.KxArticulatedSystemDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibKxArticulatedSystemDefs libraries associated by their Id.
	AllKxArticulatedSystemDefLibs = LibsKxArticulatedSystemDef{}

	//	The "default" LibKxArticulatedSystemDefs library for KxArticulatedSystemDefs.
	KxArticulatedSystemDefs = AllKxArticulatedSystemDefLibs.AddNew("")

	defaultKxArticulatedSystemInsts = map[*KxArticulatedSystemDef]*KxArticulatedSystemInst{}
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllKxArticulatedSystemDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibKxArticulatedSystemDefs contained in AllKxArticulatedSystemDefLibs) for the KxArticulatedSystemDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) KxArticulatedSystemDef() (def *KxArticulatedSystemDef) {
	id := me.S()
	for _, lib := range AllKxArticulatedSystemDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllKxArticulatedSystemDefLibs variable:
//	a hash-table that contains LibKxArticulatedSystemDefs libraries associated by their Id.
type LibsKxArticulatedSystemDef map[string]*LibKxArticulatedSystemDefs

//	Creates a new LibKxArticulatedSystemDefs library with the specified Id, adds it to this LibsKxArticulatedSystemDef, and returns it.
//	If this LibsKxArticulatedSystemDef already contains a LibKxArticulatedSystemDefs library with the specified Id, does nothing and returns nil.
func (me LibsKxArticulatedSystemDef) AddNew(id string) (lib *LibKxArticulatedSystemDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsKxArticulatedSystemDef) new(id string) (lib *LibKxArticulatedSystemDefs) {
	lib = newLibKxArticulatedSystemDefs(id)
	return
}

//	A library that contains KxArticulatedSystemDefs associated by their Id.
//	To create a new LibKxArticulatedSystemDefs library, ONLY use the LibsKxArticulatedSystemDef.New() or LibsKxArticulatedSystemDef.AddNew() methods.
type LibKxArticulatedSystemDefs struct {
	//	Id, Name
	BaseLib

	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*KxArticulatedSystemDef
}

func newLibKxArticulatedSystemDefs(id string) (me *LibKxArticulatedSystemDefs) {
	me = &LibKxArticulatedSystemDefs{M: map[string]*KxArticulatedSystemDef{}}
	me.BaseLib.init(id)
	return
}

//	Adds the specified KxArticulatedSystemDef definition to this LibKxArticulatedSystemDefs, and returns it.
//	If this LibKxArticulatedSystemDefs already contains a KxArticulatedSystemDef definition with the same Id, does nothing and returns nil.
func (me *LibKxArticulatedSystemDefs) Add(d *KxArticulatedSystemDef) (n *KxArticulatedSystemDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new KxArticulatedSystemDef definition with the specified Id, adds it to this LibKxArticulatedSystemDefs, and returns it.
//	If this LibKxArticulatedSystemDefs already contains a KxArticulatedSystemDef definition with the specified Id, does nothing and returns nil.
func (me *LibKxArticulatedSystemDefs) AddNew(id string) *KxArticulatedSystemDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibKxArticulatedSystemDefs) Len() int { return len(me.M) }

//	Creates a new KxArticulatedSystemDef definition with the specified Id and returns it,
//	but does not add it to this LibKxArticulatedSystemDefs.
func (me *LibKxArticulatedSystemDefs) New(id string) (def *KxArticulatedSystemDef) { def = newKxArticulatedSystemDef(id); return }

//	Removes the KxArticulatedSystemDef with the specified Id from this LibKxArticulatedSystemDefs.
func (me *LibKxArticulatedSystemDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibKxArticulatedSystemDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibKxArticulatedSystemDefs
//	library or its KxArticulatedSystemDef definitions. Also called by the global SyncChanges() function.
func (me *LibKxArticulatedSystemDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
