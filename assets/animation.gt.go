package assets

type AnimationChannel struct {
	Source string
	Target string
}

type AnimationSampler struct {
	HasID
	PreBehavior  string
	PostBehavior string
	Inputs       []*Input
}

type AnimationDef struct {
	BaseDef
	AnimationDefs []*AnimationDef
	Channels      []*AnimationChannel
	Samplers      []*AnimationSampler
	Sources       Sources
}

func (me *AnimationDef) Init() {
	me.Sources = Sources{}
}

type AnimationInst struct {
	BaseInst
}

func (me *AnimationInst) init() {
}

//#begin-gt _definstlib.gt T:Animation

func newAnimationDef(id string) (me *AnimationDef) {
	me = &AnimationDef{}
	me.ID = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *AnimationInst* instance referencing this *AnimationDef* definition.
func (me *AnimationDef) NewInst(id string) (inst *AnimationInst) {
	inst = &AnimationInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibAnimationDefs* libraries associated by their *ID*.
	AllAnimationDefLibs = LibsAnimationDef{}

	//	The "default" *LibAnimationDefs* library for *AnimationDef*s.
	AnimationDefs = AllAnimationDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllAnimationDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllAnimationDefLibs* variable: a *map* collection that contains
//	*LibAnimationDefs* libraries associated by their *ID*.
type LibsAnimationDef map[string]*LibAnimationDefs

//	Creates a new *LibAnimationDefs* library with the specified *ID*, adds it to this *LibsAnimationDef*, and returns it.
//	
//	If this *LibsAnimationDef* already contains a *LibAnimationDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsAnimationDef) AddNew(id string) (lib *LibAnimationDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsAnimationDef) new(id string) (lib *LibAnimationDefs) {
	lib = newLibAnimationDefs(id)
	return
}

//	A library that contains *AnimationDef*s associated by their *ID*. To create a new *LibAnimationDefs* library, ONLY
//	use the *LibsAnimationDef.New()* or *LibsAnimationDef.AddNew()* methods.
type LibAnimationDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*AnimationDef
}

func newLibAnimationDefs(id string) (me *LibAnimationDefs) {
	me = &LibAnimationDefs{M: map[string]*AnimationDef{}}
	me.ID = id
	return
}

//	Adds the specified *AnimationDef* definition to this *LibAnimationDefs*, and returns it.
//	
//	If this *LibAnimationDefs* already contains a *AnimationDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibAnimationDefs) Add(d *AnimationDef) (n *AnimationDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *AnimationDef* definition with the specified *ID*, adds it to this *LibAnimationDefs*, and returns it.
//	
//	If this *LibAnimationDefs* already contains a *AnimationDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibAnimationDefs) AddNew(id string) *AnimationDef { return me.Add(me.New(id)) }

//	Creates a new *AnimationDef* definition with the specified *ID* and returns it, but does not add it to this *LibAnimationDefs*.
func (me *LibAnimationDefs) New(id string) (def *AnimationDef) { def = newAnimationDef(id); return }

//	Removes the *AnimationDef* with the specified *ID* from this *LibAnimationDefs*.
func (me *LibAnimationDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibAnimationDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibAnimationDefs* library or its *AnimationDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibAnimationDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
