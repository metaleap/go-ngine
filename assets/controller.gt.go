package assets

type ControllerMorph struct {
	Relative bool
	Source   string
	Sources  Sources
	Targets  struct {
		HasExtras
		Inputs []*Input
	}
}

func NewControllerMorph() (me *ControllerMorph) {
	me = &ControllerMorph{Sources: Sources{}}
	return
}

type ControllerSkin struct {
	Source          string
	BindShapeMatrix *Float4x4
	Sources         Sources
	Joints          struct {
		HasExtras
		Inputs []*Input
	}
	VertexWeights IndexedInputs
}

func NewControllerSkin() (me *ControllerSkin) {
	me = &ControllerSkin{Sources: Sources{}}
	return
}

type ControllerDef struct {
	BaseDef
	Morph *ControllerMorph
	Skin  *ControllerSkin
}

func (me *ControllerDef) init() {
}

type ControllerInst struct {
	BaseInst
	Def          *ControllerDef
	BindMaterial *BindMaterial
	Skeletons    []string
}

func (me *ControllerInst) init() {
}

//#begin-gt _definstlib.gt T:Controller

func newControllerDef(id string) (me *ControllerDef) {
	me = &ControllerDef{}
	me.ID = id
	me.Base.init()
	me.init()
	return
}

/*
//	Creates and returns a new *ControllerInst* instance referencing this *ControllerDef* definition.
func (me *ControllerDef) NewInst(id string) (inst *ControllerInst) {
	inst = &ControllerInst{Def: me}
	inst.init()
	return
}
*/

var (
	//	A *map* collection that contains *LibControllerDefs* libraries associated by their *ID*.
	AllControllerDefLibs = LibsControllerDef{}

	//	The "default" *LibControllerDefs* library for *ControllerDef*s.
	ControllerDefs = AllControllerDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllControllerDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllControllerDefLibs* variable: a *map* collection that contains
//	*LibControllerDefs* libraries associated by their *ID*.
type LibsControllerDef map[string]*LibControllerDefs

//	Creates a new *LibControllerDefs* library with the specified *ID*, adds it to this *LibsControllerDef*, and returns it.
//	
//	If this *LibsControllerDef* already contains a *LibControllerDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsControllerDef) AddNew(id string) (lib *LibControllerDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsControllerDef) new(id string) (lib *LibControllerDefs) {
	lib = newLibControllerDefs(id)
	return
}

//	A library that contains *ControllerDef*s associated by their *ID*. To create a new *LibControllerDefs* library, ONLY
//	use the *LibsControllerDef.New()* or *LibsControllerDef.AddNew()* methods.
type LibControllerDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*ControllerDef
}

func newLibControllerDefs(id string) (me *LibControllerDefs) {
	me = &LibControllerDefs{M: map[string]*ControllerDef{}}
	me.ID = id
	return
}

//	Adds the specified *ControllerDef* definition to this *LibControllerDefs*, and returns it.
//	
//	If this *LibControllerDefs* already contains a *ControllerDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibControllerDefs) Add(d *ControllerDef) (n *ControllerDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *ControllerDef* definition with the specified *ID*, adds it to this *LibControllerDefs*, and returns it.
//	
//	If this *LibControllerDefs* already contains a *ControllerDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibControllerDefs) AddNew(id string) *ControllerDef { return me.Add(me.New(id)) }

//	Creates a new *ControllerDef* definition with the specified *ID* and returns it, but does not add it to this *LibControllerDefs*.
func (me *LibControllerDefs) New(id string) (def *ControllerDef) { def = newControllerDef(id); return }

//	Removes the *ControllerDef* with the specified *ID* from this *LibControllerDefs*.
func (me *LibControllerDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibControllerDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibControllerDefs* library or its *ControllerDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibControllerDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
