package assets

type KxSceneDef struct {
	BaseDef
	Models             []*KxModelInst
	ArticulatedSystems []*KxArticulatedSystemInst
}

func (me *KxSceneDef) Init() {
}

type KxSceneInst struct {
	BaseInst
	NewParams     ParamDefs
	SetParams     []*ParamInst
	ModelBindings []*KxSceneInstBindModel
	JointAxes     []*KxSceneInstJointAxis
}

func (me *KxSceneInst) init() {
	me.NewParams = ParamDefs{}
}

type KxSceneInstBindModel struct {
	Node string
	Ref  struct {
		ModelSid string
		Param    string
	}
}

type KxSceneInstJointAxis struct {
	Target string
	Axis   ParamString
	Value  ParamFloat
}

//#begin-gt _definstlib.gt T:KxScene

func newKxSceneDef(id string) (me *KxSceneDef) {
	me = &KxSceneDef{}
	me.ID = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *KxSceneInst* instance referencing this *KxSceneDef* definition.
func (me *KxSceneDef) NewInst(id string) (inst *KxSceneInst) {
	inst = &KxSceneInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibKxSceneDefs* libraries associated by their *ID*.
	AllKxSceneDefLibs = LibsKxSceneDef{}

	//	The "default" *LibKxSceneDefs* library for *KxSceneDef*s.
	KxSceneDefs = AllKxSceneDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllKxSceneDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllKxSceneDefLibs* variable: a *map* collection that contains
//	*LibKxSceneDefs* libraries associated by their *ID*.
type LibsKxSceneDef map[string]*LibKxSceneDefs

//	Creates a new *LibKxSceneDefs* library with the specified *ID*, adds it to this *LibsKxSceneDef*, and returns it.
//	
//	If this *LibsKxSceneDef* already contains a *LibKxSceneDefs* library with the specified *ID*, does nothing and returns *nil*.
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

//	A library that contains *KxSceneDef*s associated by their *ID*. To create a new *LibKxSceneDefs* library, ONLY
//	use the *LibsKxSceneDef.New()* or *LibsKxSceneDef.AddNew()* methods.
type LibKxSceneDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*KxSceneDef
}

func newLibKxSceneDefs(id string) (me *LibKxSceneDefs) {
	me = &LibKxSceneDefs{M: map[string]*KxSceneDef{}}
	me.ID = id
	return
}

//	Adds the specified *KxSceneDef* definition to this *LibKxSceneDefs*, and returns it.
//	
//	If this *LibKxSceneDefs* already contains a *KxSceneDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibKxSceneDefs) Add(d *KxSceneDef) (n *KxSceneDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *KxSceneDef* definition with the specified *ID*, adds it to this *LibKxSceneDefs*, and returns it.
//	
//	If this *LibKxSceneDefs* already contains a *KxSceneDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibKxSceneDefs) AddNew(id string) *KxSceneDef { return me.Add(me.New(id)) }

//	Creates a new *KxSceneDef* definition with the specified *ID* and returns it, but does not add it to this *LibKxSceneDefs*.
func (me *LibKxSceneDefs) New(id string) (def *KxSceneDef) { def = newKxSceneDef(id); return }

//	Removes the *KxSceneDef* with the specified *ID* from this *LibKxSceneDefs*.
func (me *LibKxSceneDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibKxSceneDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibKxSceneDefs* library or its *KxSceneDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibKxSceneDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
