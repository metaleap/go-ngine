package assets

type AnimationClipDef struct {
	BaseDef
	Start      float64
	End        float64
	Animations []*AnimationInst
	Formulas   []*FormulaInst
}

func (me *AnimationClipDef) init() {
}

//#begin-gt _definstlib.gt T:AnimationClip

func newAnimationClipDef(id string) (me *AnimationClipDef) {
	me = &AnimationClipDef{}
	me.ID = id
	me.Base.init()
	me.init()
	return
}

/*
//	Creates and returns a new *AnimationClipInst* instance referencing this *AnimationClipDef* definition.
func (me *AnimationClipDef) NewInst(id string) (inst *AnimationClipInst) {
	inst = &AnimationClipInst{Def: me}
	inst.init()
	return
}
*/

var (
	//	A *map* collection that contains *LibAnimationClipDefs* libraries associated by their *ID*.
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
//	*LibAnimationClipDefs* libraries associated by their *ID*.
type LibsAnimationClipDef map[string]*LibAnimationClipDefs

//	Creates a new *LibAnimationClipDefs* library with the specified *ID*, adds it to this *LibsAnimationClipDef*, and returns it.
//	
//	If this *LibsAnimationClipDef* already contains a *LibAnimationClipDefs* library with the specified *ID*, does nothing and returns *nil*.
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

//	A library that contains *AnimationClipDef*s associated by their *ID*. To create a new *LibAnimationClipDefs* library, ONLY
//	use the *LibsAnimationClipDef.New()* or *LibsAnimationClipDef.AddNew()* methods.
type LibAnimationClipDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*AnimationClipDef
}

func newLibAnimationClipDefs(id string) (me *LibAnimationClipDefs) {
	me = &LibAnimationClipDefs{M: map[string]*AnimationClipDef{}}
	me.ID = id
	return
}

//	Adds the specified *AnimationClipDef* definition to this *LibAnimationClipDefs*, and returns it.
//	
//	If this *LibAnimationClipDefs* already contains a *AnimationClipDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibAnimationClipDefs) Add(d *AnimationClipDef) (n *AnimationClipDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *AnimationClipDef* definition with the specified *ID*, adds it to this *LibAnimationClipDefs*, and returns it.
//	
//	If this *LibAnimationClipDefs* already contains a *AnimationClipDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibAnimationClipDefs) AddNew(id string) *AnimationClipDef { return me.Add(me.New(id)) }

//	Creates a new *AnimationClipDef* definition with the specified *ID* and returns it, but does not add it to this *LibAnimationClipDefs*.
func (me *LibAnimationClipDefs) New(id string) (def *AnimationClipDef) { def = newAnimationClipDef(id); return }

//	Removes the *AnimationClipDef* with the specified *ID* from this *LibAnimationClipDefs*.
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
