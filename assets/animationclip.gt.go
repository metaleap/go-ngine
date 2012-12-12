package assets

//	Defines a section of a set of animation curves and/or formulas to be used together as an animation clip.
type AnimationClipDef struct {
	//	Id, Name, Asset, Extras
	BaseDef

	//	The time in seconds of the beginning of the clip.
	Start float64

	//	The time in seconds of the end of the clip.
	End float64

	//	The animation instances contributing to this animation clip.
	Animations []*AnimationInst

	//	Any formulas used in this animation clip.
	Formulas []*FormulaInst
}

//	Initialization
func (me *AnimationClipDef) Init() {
}

//#begin-gt _definstlib.gt T:AnimationClip

func newAnimationClipDef(id string) (me *AnimationClipDef) {
	me = &AnimationClipDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *AnimationClipInst* instance referencing this *AnimationClipDef* definition.
func (me *AnimationClipDef) NewInst(id string) (inst *AnimationClipInst) {
	inst = &AnimationClipInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibAnimationClipDefs* libraries associated by their *Id*.
	AllAnimationClipDefLibs = LibsAnimationClipDef{}

	//	The "default" *LibAnimationClipDefs* library for *AnimationClipDef*s.
	AnimationClipDefs = AllAnimationClipDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllAnimationClipDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllAnimationClipDefLibs* variable: a *map* collection that contains
//	*LibAnimationClipDefs* libraries associated by their *Id*.
type LibsAnimationClipDef map[string]*LibAnimationClipDefs

//	Creates a new *LibAnimationClipDefs* library with the specified *Id*, adds it to this *LibsAnimationClipDef*, and returns it.
//	
//	If this *LibsAnimationClipDef* already contains a *LibAnimationClipDefs* library with the specified *Id*, does nothing and returns *nil*.
func (me LibsAnimationClipDef) AddNew(id string) (lib *LibAnimationClipDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsAnimationClipDef) new(id string) (lib *LibAnimationClipDefs) {
	lib = newLibAnimationClipDefs(id)
	return
}

//	A library that contains *AnimationClipDef*s associated by their *Id*. To create a new *LibAnimationClipDefs* library, ONLY
//	use the *LibsAnimationClipDef.New()* or *LibsAnimationClipDef.AddNew()* methods.
type LibAnimationClipDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*AnimationClipDef
}

func newLibAnimationClipDefs(id string) (me *LibAnimationClipDefs) {
	me = &LibAnimationClipDefs{M: map[string]*AnimationClipDef{}}
	me.Id = id
	return
}

//	Adds the specified *AnimationClipDef* definition to this *LibAnimationClipDefs*, and returns it.
//	
//	If this *LibAnimationClipDefs* already contains a *AnimationClipDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibAnimationClipDefs) Add(d *AnimationClipDef) (n *AnimationClipDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *AnimationClipDef* definition with the specified *Id*, adds it to this *LibAnimationClipDefs*, and returns it.
//	
//	If this *LibAnimationClipDefs* already contains a *AnimationClipDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibAnimationClipDefs) AddNew(id string) *AnimationClipDef { return me.Add(me.New(id)) }

//	Creates a new *AnimationClipDef* definition with the specified *Id* and returns it, but does not add it to this *LibAnimationClipDefs*.
func (me *LibAnimationClipDefs) New(id string) (def *AnimationClipDef) { def = newAnimationClipDef(id); return }

//	Removes the *AnimationClipDef* with the specified *Id* from this *LibAnimationClipDefs*.
func (me *LibAnimationClipDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibAnimationClipDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibAnimationClipDefs* library or its *AnimationClipDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibAnimationClipDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
