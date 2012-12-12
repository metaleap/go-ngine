package assets

type VisualSceneEvaluation struct {
	BaseDef
	HasSid
	Disabled     bool
	RenderPasses []*VisualSceneRendering
}

type VisualSceneRendering struct {
	HasName
	HasSid
	HasExtras
	CameraNode   string
	Layers       Layers
	MaterialInst *VisualSceneRenderingMaterialInst
}

func NewVisualSceneRendering() (me *VisualSceneRendering) {
	me = &VisualSceneRendering{Layers: Layers{}}
	return
}

type VisualSceneRenderingMaterialInst struct {
	HasExtras
	Bindings          []*FxMaterialBinding
	OverrideTechnique struct {
		Ref  string
		Pass string
	}
}

type VisualSceneDef struct {
	BaseDef
	Nodes       []*NodeDef
	Evaluations []*VisualSceneEvaluation
}

func (me *VisualSceneDef) Init() {
}

type VisualSceneInst struct {
	BaseInst
}

func (me *VisualSceneInst) Init() {
}

//#begin-gt _definstlib.gt T:VisualScene

func newVisualSceneDef(id string) (me *VisualSceneDef) {
	me = &VisualSceneDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *VisualSceneInst* instance referencing this *VisualSceneDef* definition.
func (me *VisualSceneDef) NewInst(id string) (inst *VisualSceneInst) {
	inst = &VisualSceneInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibVisualSceneDefs* libraries associated by their *Id*.
	AllVisualSceneDefLibs = LibsVisualSceneDef{}

	//	The "default" *LibVisualSceneDefs* library for *VisualSceneDef*s.
	VisualSceneDefs = AllVisualSceneDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllVisualSceneDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllVisualSceneDefLibs* variable: a *map* collection that contains
//	*LibVisualSceneDefs* libraries associated by their *Id*.
type LibsVisualSceneDef map[string]*LibVisualSceneDefs

//	Creates a new *LibVisualSceneDefs* library with the specified *Id*, adds it to this *LibsVisualSceneDef*, and returns it.
//	
//	If this *LibsVisualSceneDef* already contains a *LibVisualSceneDefs* library with the specified *Id*, does nothing and returns *nil*.
func (me LibsVisualSceneDef) AddNew(id string) (lib *LibVisualSceneDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsVisualSceneDef) new(id string) (lib *LibVisualSceneDefs) {
	lib = newLibVisualSceneDefs(id)
	return
}

//	A library that contains *VisualSceneDef*s associated by their *Id*. To create a new *LibVisualSceneDefs* library, ONLY
//	use the *LibsVisualSceneDef.New()* or *LibsVisualSceneDef.AddNew()* methods.
type LibVisualSceneDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*VisualSceneDef
}

func newLibVisualSceneDefs(id string) (me *LibVisualSceneDefs) {
	me = &LibVisualSceneDefs{M: map[string]*VisualSceneDef{}}
	me.Id = id
	return
}

//	Adds the specified *VisualSceneDef* definition to this *LibVisualSceneDefs*, and returns it.
//	
//	If this *LibVisualSceneDefs* already contains a *VisualSceneDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibVisualSceneDefs) Add(d *VisualSceneDef) (n *VisualSceneDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *VisualSceneDef* definition with the specified *Id*, adds it to this *LibVisualSceneDefs*, and returns it.
//	
//	If this *LibVisualSceneDefs* already contains a *VisualSceneDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibVisualSceneDefs) AddNew(id string) *VisualSceneDef { return me.Add(me.New(id)) }

//	Creates a new *VisualSceneDef* definition with the specified *Id* and returns it, but does not add it to this *LibVisualSceneDefs*.
func (me *LibVisualSceneDefs) New(id string) (def *VisualSceneDef) { def = newVisualSceneDef(id); return }

//	Removes the *VisualSceneDef* with the specified *Id* from this *LibVisualSceneDefs*.
func (me *LibVisualSceneDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibVisualSceneDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibVisualSceneDefs* library or its *VisualSceneDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibVisualSceneDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
