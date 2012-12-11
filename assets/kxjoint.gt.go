package assets

const (
	KX_JOINT_TYPE_PRISMATIC = 0
	KX_JOINT_TYPE_REVOLUTE  = iota
)

type KxJoint struct {
	HasSid
	Axis struct {
		HasSid
		HasName
		F Float3
	}
	Limits *KxJointLimits
	Type   int
}

type KxJointLimits struct {
	Min *ScopedFloat
	Max *ScopedFloat
}

type KxJointDef struct {
	BaseDef
	HasSid
	All []*KxJoint
}

func (me *KxJointDef) Init() {
}

type KxJointInst struct {
	BaseInst
}

func (me *KxJointInst) Init() {
}

//#begin-gt _definstlib.gt T:KxJoint

func newKxJointDef(id string) (me *KxJointDef) {
	me = &KxJointDef{}
	me.ID = id
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
	//	A *map* collection that contains *LibKxJointDefs* libraries associated by their *ID*.
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
//	*LibKxJointDefs* libraries associated by their *ID*.
type LibsKxJointDef map[string]*LibKxJointDefs

//	Creates a new *LibKxJointDefs* library with the specified *ID*, adds it to this *LibsKxJointDef*, and returns it.
//	
//	If this *LibsKxJointDef* already contains a *LibKxJointDefs* library with the specified *ID*, does nothing and returns *nil*.
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

//	A library that contains *KxJointDef*s associated by their *ID*. To create a new *LibKxJointDefs* library, ONLY
//	use the *LibsKxJointDef.New()* or *LibsKxJointDef.AddNew()* methods.
type LibKxJointDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*KxJointDef
}

func newLibKxJointDefs(id string) (me *LibKxJointDefs) {
	me = &LibKxJointDefs{M: map[string]*KxJointDef{}}
	me.ID = id
	return
}

//	Adds the specified *KxJointDef* definition to this *LibKxJointDefs*, and returns it.
//	
//	If this *LibKxJointDefs* already contains a *KxJointDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibKxJointDefs) Add(d *KxJointDef) (n *KxJointDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *KxJointDef* definition with the specified *ID*, adds it to this *LibKxJointDefs*, and returns it.
//	
//	If this *LibKxJointDefs* already contains a *KxJointDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibKxJointDefs) AddNew(id string) *KxJointDef { return me.Add(me.New(id)) }

//	Creates a new *KxJointDef* definition with the specified *ID* and returns it, but does not add it to this *LibKxJointDefs*.
func (me *LibKxJointDefs) New(id string) (def *KxJointDef) { def = newKxJointDef(id); return }

//	Removes the *KxJointDef* with the specified *ID* from this *LibKxJointDefs*.
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
