package assets

const (
	//	Defines the base frame for kinematics calculation.
	KX_FRAME_TYPE_ORIGIN = 1
	//	Defines the frame at the end of the kinematics chain.
	KX_FRAME_TYPE_TIP = iota
	//	Defines the offset frame from the kinematics KX_FRAME_TYPE_TIP frame, which usually represents the work point of the end effector (for example, a welding gun).
	KX_FRAME_TYPE_TCP = iota
	//	Defines the offset frame from the kinematics KX_FRAME_TYPE_ORIGIN frame; this offset usually represents the transformation to a work piece.
	KX_FRAME_TYPE_OBJECT = iota
)

//	Specifies the parent axis’ index in the jointmap.
type KxAxisIndex struct {
	//	If set, specifies the special use of this index.
	Semantic string
	//	If not set, the parent axis will not appear in the jointmap.
	I ParamInt
}

//	Specifies the parent axis' soft limits.
type KxAxisLimits struct {
	//	The "minimum" portion of this limits descriptor.
	Min ParamFloat
	//	The "maximum" portion of this limits descriptor.
	Max ParamFloat
}

//	Binds inputs to kinematics parameters upon instantiation.
type KxBinding struct {
	//	The identifier of the parameter to bind to the new symbol name. Required.
	Symbol string
	//	If set, Value is ignored.
	ParamRef RefParam
	//	Only used if ParamRef is empty.
	Value interface{}
}

type KxArticulatedSystemEffector struct {
	HasSid
	HasName
	HasParamDefs
	HasParamInsts
	Bindings     []*KxBinding
	Speed        *ParamFloat2
	Acceleration *ParamFloat2
	Deceleration *ParamFloat2
	Jerk         *ParamFloat2
}

func NewKxArticulatedSystemEffector() (me *KxArticulatedSystemEffector) {
	me = &KxArticulatedSystemEffector{}
	me.NewParams = ParamDefs{}
	return
}

type KxArticulatedSystemKinematics struct {
	//	Techniques
	HasTechniques
	Models []*KxModelInst
	//	Common-technique profile
	TC struct {
		AxisInfos []*KxArticulatedSystemKinematicsAxis
		Frame     struct {
			Origin KxArticulatedSystemKinematicsFrame
			Tip    KxArticulatedSystemKinematicsFrame
			Tcp    *KxArticulatedSystemKinematicsFrame
			Object *KxArticulatedSystemKinematicsFrame
		}
	}
}

type KxArticulatedSystemKinematicsAxis struct {
	HasSid
	HasName
	HasParamDefs
	JointAxis string
	Active    ParamBool
	Indices   []*KxAxisIndex
	Limits    *KxAxisLimits
	Locked    ParamBool
	Formulas  struct {
		Defs  []*FormulaDef
		Insts []*FormulaInst
	}
}

func NewKxArticulatedSystemKinematicsAxis() (me *KxArticulatedSystemKinematicsAxis) {
	me = &KxArticulatedSystemKinematicsAxis{Active: ParamBool{true, ""}}
	me.NewParams = ParamDefs{}
	return
}

type KxArticulatedSystemKinematicsFrame struct {
	Link       string
	Type       int
	Transforms []*Transform
}

type KxArticulatedSystemMotion struct {
	//	Techniques
	HasTechniques
	ArticulatedSystem *KxArticulatedSystemInst
	//	Common-technique profile
	TC struct {
		AxisInfos     []*KxArticulatedSystemMotionAxis
		EffectorInfos []*KxArticulatedSystemEffector
	}
}

type KxArticulatedSystemMotionAxis struct {
	HasSid
	HasName
	HasParamDefs
	HasParamInsts
	Axis         string
	Bindings     []*KxBinding
	Speed        *ParamFloat
	Acceleration *ParamFloat
	Deceleration *ParamFloat
	Jerk         *ParamFloat
}

func NewKxArticulatedSystemMotionAxis() (me *KxArticulatedSystemMotionAxis) {
	me = &KxArticulatedSystemMotionAxis{}
	me.NewParams = ParamDefs{}
	return
}

//	Categorizes the declaration of generic control information for kinematics systems.
type KxArticulatedSystemDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	If set, Motion must be nil, and this articulated system describes a kinematics system.
	Kinematics *KxArticulatedSystemKinematics
	//	If set, Kinematics must be nil, and this articulated system describes a motion system.
	Motion *KxArticulatedSystemMotion
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
	//	Bindings of inputs to kinematics parameters.
	Bindings []*KxBinding
}

//	Initialization
func (me *KxArticulatedSystemInst) Init() {
	me.NewParams = ParamDefs{}
}

//#begin-gt _definstlib.gt T:KxArticulatedSystem

func newKxArticulatedSystemDef(id string) (me *KxArticulatedSystemDef) {
	me = &KxArticulatedSystemDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *KxArticulatedSystemInst* instance referencing this *KxArticulatedSystemDef* definition.
func (me *KxArticulatedSystemDef) NewInst(id string) (inst *KxArticulatedSystemInst) {
	inst = &KxArticulatedSystemInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibKxArticulatedSystemDefs* libraries associated by their *Id*.
	AllKxArticulatedSystemDefLibs = LibsKxArticulatedSystemDef{}

	//	The "default" *LibKxArticulatedSystemDefs* library for *KxArticulatedSystemDef*s.
	KxArticulatedSystemDefs = AllKxArticulatedSystemDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllKxArticulatedSystemDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllKxArticulatedSystemDefLibs* variable: a *map* collection that contains
//	*LibKxArticulatedSystemDefs* libraries associated by their *Id*.
type LibsKxArticulatedSystemDef map[string]*LibKxArticulatedSystemDefs

//	Creates a new *LibKxArticulatedSystemDefs* library with the specified *Id*, adds it to this *LibsKxArticulatedSystemDef*, and returns it.
//	
//	If this *LibsKxArticulatedSystemDef* already contains a *LibKxArticulatedSystemDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *KxArticulatedSystemDef*s associated by their *Id*. To create a new *LibKxArticulatedSystemDefs* library, ONLY
//	use the *LibsKxArticulatedSystemDef.New()* or *LibsKxArticulatedSystemDef.AddNew()* methods.
type LibKxArticulatedSystemDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*KxArticulatedSystemDef
}

func newLibKxArticulatedSystemDefs(id string) (me *LibKxArticulatedSystemDefs) {
	me = &LibKxArticulatedSystemDefs{M: map[string]*KxArticulatedSystemDef{}}
	me.Id = id
	return
}

//	Adds the specified *KxArticulatedSystemDef* definition to this *LibKxArticulatedSystemDefs*, and returns it.
//	
//	If this *LibKxArticulatedSystemDefs* already contains a *KxArticulatedSystemDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibKxArticulatedSystemDefs) Add(d *KxArticulatedSystemDef) (n *KxArticulatedSystemDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *KxArticulatedSystemDef* definition with the specified *Id*, adds it to this *LibKxArticulatedSystemDefs*, and returns it.
//	
//	If this *LibKxArticulatedSystemDefs* already contains a *KxArticulatedSystemDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibKxArticulatedSystemDefs) AddNew(id string) *KxArticulatedSystemDef {
	return me.Add(me.New(id))
}

//	Short-hand for len(lib.M)
func (me *LibKxArticulatedSystemDefs) Len() int { return len(me.M) }

//	Creates a new *KxArticulatedSystemDef* definition with the specified *Id* and returns it, but does not add it to this *LibKxArticulatedSystemDefs*.
func (me *LibKxArticulatedSystemDefs) New(id string) (def *KxArticulatedSystemDef) {
	def = newKxArticulatedSystemDef(id)
	return
}

//	Removes the *KxArticulatedSystemDef* with the specified *Id* from this *LibKxArticulatedSystemDefs*.
func (me *LibKxArticulatedSystemDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibKxArticulatedSystemDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibKxArticulatedSystemDefs* library or its *KxArticulatedSystemDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibKxArticulatedSystemDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
