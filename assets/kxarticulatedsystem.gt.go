package assets

const (
	KX_FRAME_TYPE_OBJECT = 0
	KX_FRAME_TYPE_ORIGIN = iota
	KX_FRAME_TYPE_TCP    = iota
	KX_FRAME_TYPE_TIP    = iota
)

type KxBind struct {
	Symbol   string
	ParamRef string
	Value    interface{}
}

type KxArticulatedSystemAxisIndex struct {
	Semantic string
	I        ParamInt
}

type KxArticulatedSystemAxisLimits struct {
	Min ParamFloat
	Max ParamFloat
}

type KxArticulatedSystemEffector struct {
	HasSid
	HasName
	HasParamDefs
	HasParamInsts
	Bindings     []*KxBind
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
	HasTechniques
	Models []*KxModelInst
	TC     struct {
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
	Indices   []*KxArticulatedSystemAxisIndex
	Locked    ParamBool
	Limits    *KxArticulatedSystemAxisLimits
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
	HasTechniques
	ArticulatedSystem *KxArticulatedSystemInst
	TC                struct {
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
	Bindings     []*KxBind
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

type KxArticulatedSystemDef struct {
	BaseDef
	Kinematics *KxArticulatedSystemKinematics
	Motion     *KxArticulatedSystemMotion
}

func (me *KxArticulatedSystemDef) Init() {
}

type KxArticulatedSystemInst struct {
	BaseInst
	HasParamDefs
	HasParamInsts
	Bindings []*KxBind
}

func (me *KxArticulatedSystemInst) Init() {
	me.NewParams = ParamDefs{}
}

//#begin-gt _definstlib.gt T:KxArticulatedSystem

func newKxArticulatedSystemDef(id string) (me *KxArticulatedSystemDef) {
	me = &KxArticulatedSystemDef{}
	me.ID = id
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
	//	A *map* collection that contains *LibKxArticulatedSystemDefs* libraries associated by their *ID*.
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
//	*LibKxArticulatedSystemDefs* libraries associated by their *ID*.
type LibsKxArticulatedSystemDef map[string]*LibKxArticulatedSystemDefs

//	Creates a new *LibKxArticulatedSystemDefs* library with the specified *ID*, adds it to this *LibsKxArticulatedSystemDef*, and returns it.
//	
//	If this *LibsKxArticulatedSystemDef* already contains a *LibKxArticulatedSystemDefs* library with the specified *ID*, does nothing and returns *nil*.
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

//	A library that contains *KxArticulatedSystemDef*s associated by their *ID*. To create a new *LibKxArticulatedSystemDefs* library, ONLY
//	use the *LibsKxArticulatedSystemDef.New()* or *LibsKxArticulatedSystemDef.AddNew()* methods.
type LibKxArticulatedSystemDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*KxArticulatedSystemDef
}

func newLibKxArticulatedSystemDefs(id string) (me *LibKxArticulatedSystemDefs) {
	me = &LibKxArticulatedSystemDefs{M: map[string]*KxArticulatedSystemDef{}}
	me.ID = id
	return
}

//	Adds the specified *KxArticulatedSystemDef* definition to this *LibKxArticulatedSystemDefs*, and returns it.
//	
//	If this *LibKxArticulatedSystemDefs* already contains a *KxArticulatedSystemDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibKxArticulatedSystemDefs) Add(d *KxArticulatedSystemDef) (n *KxArticulatedSystemDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *KxArticulatedSystemDef* definition with the specified *ID*, adds it to this *LibKxArticulatedSystemDefs*, and returns it.
//	
//	If this *LibKxArticulatedSystemDefs* already contains a *KxArticulatedSystemDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibKxArticulatedSystemDefs) AddNew(id string) *KxArticulatedSystemDef { return me.Add(me.New(id)) }

//	Creates a new *KxArticulatedSystemDef* definition with the specified *ID* and returns it, but does not add it to this *LibKxArticulatedSystemDefs*.
func (me *LibKxArticulatedSystemDefs) New(id string) (def *KxArticulatedSystemDef) { def = newKxArticulatedSystemDef(id); return }

//	Removes the *KxArticulatedSystemDef* with the specified *ID* from this *LibKxArticulatedSystemDefs*.
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
