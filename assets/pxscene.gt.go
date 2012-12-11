package assets

type PxSceneDef struct {
	BaseDef
	HasTechniques
	ForceFields []*PxForceFieldInst
	Models      []*PxModelInst
	TC          struct {
		Gravity  *ScopedVec3
		TimeStep *ScopedFloat
	}
}

func (me *PxSceneDef) Init() {
}

type PxSceneInst struct {
	BaseInst
}

func (me *PxSceneInst) Init() {
}

//#begin-gt _definstlib.gt T:PxScene

func newPxSceneDef(id string) (me *PxSceneDef) {
	me = &PxSceneDef{}
	me.ID = id
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
	//	A *map* collection that contains *LibPxSceneDefs* libraries associated by their *ID*.
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
//	*LibPxSceneDefs* libraries associated by their *ID*.
type LibsPxSceneDef map[string]*LibPxSceneDefs

//	Creates a new *LibPxSceneDefs* library with the specified *ID*, adds it to this *LibsPxSceneDef*, and returns it.
//	
//	If this *LibsPxSceneDef* already contains a *LibPxSceneDefs* library with the specified *ID*, does nothing and returns *nil*.
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

//	A library that contains *PxSceneDef*s associated by their *ID*. To create a new *LibPxSceneDefs* library, ONLY
//	use the *LibsPxSceneDef.New()* or *LibsPxSceneDef.AddNew()* methods.
type LibPxSceneDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*PxSceneDef
}

func newLibPxSceneDefs(id string) (me *LibPxSceneDefs) {
	me = &LibPxSceneDefs{M: map[string]*PxSceneDef{}}
	me.ID = id
	return
}

//	Adds the specified *PxSceneDef* definition to this *LibPxSceneDefs*, and returns it.
//	
//	If this *LibPxSceneDefs* already contains a *PxSceneDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibPxSceneDefs) Add(d *PxSceneDef) (n *PxSceneDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *PxSceneDef* definition with the specified *ID*, adds it to this *LibPxSceneDefs*, and returns it.
//	
//	If this *LibPxSceneDefs* already contains a *PxSceneDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibPxSceneDefs) AddNew(id string) *PxSceneDef { return me.Add(me.New(id)) }

//	Creates a new *PxSceneDef* definition with the specified *ID* and returns it, but does not add it to this *LibPxSceneDefs*.
func (me *LibPxSceneDefs) New(id string) (def *PxSceneDef) { def = newPxSceneDef(id); return }

//	Removes the *PxSceneDef* with the specified *ID* from this *LibPxSceneDefs*.
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
