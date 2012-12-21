package assets

//	Declares information specifying how to evaluate a visual scene.
type VisualSceneEvaluation struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sid
	HasSid
	//	Whether evaluation is enabled. Disabling evaluation can be useful for debugging.
	Disabled bool
	//	Describes effects passes to render a scene.
	RenderPasses []*VisualSceneRendering
}

//	Describes one effect pass to evaluate a scene.
type VisualSceneRendering struct {
	//	Name
	HasName
	//	Sid
	HasSid
	//	Extras
	HasExtras
	//	Refers to a NodeDef that contains a camera describing
	//	the viewpoint from which to render this compositing step. Optional.
	CameraNode RefId
	//	Specifies which layer or layers to render in this compositing step while evaluating the scene.
	Layers Layers
	//	If set, specifies which effect to render in this compositing step while evaluating the scene.
	MaterialInst *VisualSceneRenderingMaterialInst
}

//	Constructor
func NewVisualSceneRendering() (me *VisualSceneRendering) {
	me = &VisualSceneRendering{Layers: Layers{}}
	return
}

//	Instantiates a material resource for a screen effect.
type VisualSceneRenderingMaterialInst struct {
	//	Extras
	HasExtras
	//	Binds values to effect parameters upon instantiation.
	Bindings []*FxBinding
	//	Target specific techniques and passes inside a material
	//	rather than having to split the effects techniques and passes into multiple effects.
	OverrideTechnique struct {
		//	Specifies the Sid of a Technique
		Ref RefSid
		//	Specifies the Sid of one FxPass to execute.
		//	If not specified, then all of the Technique's passes are used.
		Pass RefSid
	}
}

//	Embodies the entire set of information that can be visualized from the contents of a resource.
type VisualSceneDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	A scene graph containing nodes of visual information and related data.
	Nodes []*NodeDef
	//	Specifies how to evaluate this visual scene.
	Evaluations []*VisualSceneEvaluation
}

//	Initialization
func (me *VisualSceneDef) Init() {
}

//	Instantiates a visual scene resource.
type VisualSceneInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *VisualSceneDef
}

//	Initialization
func (me *VisualSceneInst) Init() {
}

//#begin-gt _definstlib.gt T:VisualScene

func newVisualSceneDef(id string) (me *VisualSceneDef) {
	me = &VisualSceneDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new VisualSceneInst instance referencing this VisualSceneDef definition.
//	Any VisualSceneInst created by this method will have its Def field readily set to me.
func (me *VisualSceneDef) NewInst() (inst *VisualSceneInst) {
	inst = &VisualSceneInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct VisualSceneDef
//	according to the current me.DefRef value (by searching AllVisualSceneDefLibs).
//	Then returns me.Def.
//	(Note, every VisualSceneInst's Def is nil initially, unless it was created via VisualSceneDef.NewInst().)
func (me *VisualSceneInst) EnsureDef() *VisualSceneDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.VisualSceneDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibVisualSceneDefs libraries associated by their Id.
	AllVisualSceneDefLibs = LibsVisualSceneDef{}

	//	The "default" LibVisualSceneDefs library for VisualSceneDefs.
	VisualSceneDefs = AllVisualSceneDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllVisualSceneDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibVisualSceneDefs contained in AllVisualSceneDefLibs) for the VisualSceneDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) VisualSceneDef() (def *VisualSceneDef) {
	id := me.S()
	for _, lib := range AllVisualSceneDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllVisualSceneDefLibs variable:
//	a hash-table that contains LibVisualSceneDefs libraries associated by their Id.
type LibsVisualSceneDef map[string]*LibVisualSceneDefs

//	Creates a new LibVisualSceneDefs library with the specified Id, adds it to this LibsVisualSceneDef, and returns it.
//	If this LibsVisualSceneDef already contains a LibVisualSceneDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains VisualSceneDefs associated by their Id.
//	To create a new LibVisualSceneDefs library, ONLY use the LibsVisualSceneDef.New() or LibsVisualSceneDef.AddNew() methods.
type LibVisualSceneDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*VisualSceneDef
}

func newLibVisualSceneDefs(id string) (me *LibVisualSceneDefs) {
	me = &LibVisualSceneDefs{M: map[string]*VisualSceneDef{}}
	me.Id = id
	return
}

//	Adds the specified VisualSceneDef definition to this LibVisualSceneDefs, and returns it.
//	If this LibVisualSceneDefs already contains a VisualSceneDef definition with the same Id, does nothing and returns nil.
func (me *LibVisualSceneDefs) Add(d *VisualSceneDef) (n *VisualSceneDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new VisualSceneDef definition with the specified Id, adds it to this LibVisualSceneDefs, and returns it.
//	If this LibVisualSceneDefs already contains a VisualSceneDef definition with the specified Id, does nothing and returns nil.
func (me *LibVisualSceneDefs) AddNew(id string) *VisualSceneDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibVisualSceneDefs) Len() int { return len(me.M) }

//	Creates a new VisualSceneDef definition with the specified Id and returns it,
//	but does not add it to this LibVisualSceneDefs.
func (me *LibVisualSceneDefs) New(id string) (def *VisualSceneDef) { def = newVisualSceneDef(id); return }

//	Removes the VisualSceneDef with the specified Id from this LibVisualSceneDefs.
func (me *LibVisualSceneDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Returns a GetRefSidResolver that looks up and yields the VisualSceneDef with the specified Id.
func (me *LibVisualSceneDefs) ResolverGetter() GetRefSidResolver {
	return func(id string) RefSidResolver {
		return nil // me.M[id]
	}
}

//	Signals to the core package (or your custom package) that changes have been made to this LibVisualSceneDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibVisualSceneDefs
//	library or its VisualSceneDef definitions. Also called by the global SyncChanges() function.
func (me *LibVisualSceneDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
