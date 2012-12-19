package assets

//	Specifies an environment in which physical objects are instantiated and simulated.
type PxSceneDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Techniques
	HasTechniques
	//	Force fields influencing this physics scene.
	ForceFields []*PxForceFieldInst
	//	Refers to the rigid bodies and constraints participating in this scene.
	Models []*PxModelInst
	//	Common-technique profile
	TC struct {
		//	If set, a vector representation of this physics scene’s gravity force field. It is given as a denormalized direction vector of three floating-point values that indicate both the magnitude and direction of acceleration caused by the field.
		Gravity *ScopedVec3
		//	If set, the integration time step, measured in seconds, of the physics scene. This value is engine-specific. If omitted, the physics engine’s default is used.
		TimeStep *ScopedFloat
	}
}

//	Initialization
func (me *PxSceneDef) Init() {
}

//	Instantiates a physics scene resource.
type PxSceneInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
}

//	Initialization
func (me *PxSceneInst) Init() {
}

//#begin-gt _definstlib.gt T:PxScene

func newPxSceneDef(id string) (me *PxSceneDef) {
	me = &PxSceneDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *PxSceneInst* instance referencing this *PxSceneDef* definition.
func (me *PxSceneDef) NewInst(id string) (inst *PxSceneInst) {
	inst = &PxSceneInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibPxSceneDefs* libraries associated by their *Id*.
	AllPxSceneDefLibs = LibsPxSceneDef{}

	//	The "default" *LibPxSceneDefs* library for *PxSceneDef*s.
	PxSceneDefs = AllPxSceneDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllPxSceneDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllPxSceneDefLibs* variable: a *map* collection that contains
//	*LibPxSceneDefs* libraries associated by their *Id*.
type LibsPxSceneDef map[string]*LibPxSceneDefs

//	Creates a new *LibPxSceneDefs* library with the specified *Id*, adds it to this *LibsPxSceneDef*, and returns it.
//	
//	If this *LibsPxSceneDef* already contains a *LibPxSceneDefs* library with the specified *Id*, does nothing and returns *nil*.
func (me LibsPxSceneDef) AddNew(id string) (lib *LibPxSceneDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsPxSceneDef) new(id string) (lib *LibPxSceneDefs) {
	lib = newLibPxSceneDefs(id)
	return
}

//	A library that contains *PxSceneDef*s associated by their *Id*. To create a new *LibPxSceneDefs* library, ONLY
//	use the *LibsPxSceneDef.New()* or *LibsPxSceneDef.AddNew()* methods.
type LibPxSceneDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*PxSceneDef
}

func newLibPxSceneDefs(id string) (me *LibPxSceneDefs) {
	me = &LibPxSceneDefs{M: map[string]*PxSceneDef{}}
	me.Id = id
	return
}

//	Adds the specified *PxSceneDef* definition to this *LibPxSceneDefs*, and returns it.
//	
//	If this *LibPxSceneDefs* already contains a *PxSceneDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibPxSceneDefs) Add(d *PxSceneDef) (n *PxSceneDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *PxSceneDef* definition with the specified *Id*, adds it to this *LibPxSceneDefs*, and returns it.
//	
//	If this *LibPxSceneDefs* already contains a *PxSceneDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibPxSceneDefs) AddNew(id string) *PxSceneDef { return me.Add(me.New(id)) }

//	Short-hand for len(lib.M)
func (me *LibPxSceneDefs) Len() int { return len(me.M) }

//	Creates a new *PxSceneDef* definition with the specified *Id* and returns it, but does not add it to this *LibPxSceneDefs*.
func (me *LibPxSceneDefs) New(id string) (def *PxSceneDef) { def = newPxSceneDef(id); return }

//	Removes the *PxSceneDef* with the specified *Id* from this *LibPxSceneDefs*.
func (me *LibPxSceneDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibPxSceneDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibPxSceneDefs* library or its *PxSceneDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibPxSceneDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
