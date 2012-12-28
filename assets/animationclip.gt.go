package assets

//	Defines a section of a set of animation curves and/or formulas
//	to be used together as an animation clip.
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

//	Instantiates an animation clip resource.
type AnimationClipInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst

	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *AnimationClipDef
}

//	Initialization
func (me *AnimationClipInst) Init() {
}

//#begin-gt _definstlib.gt T:AnimationClip

func newAnimationClipDef(id string) (me *AnimationClipDef) {
	me = &AnimationClipDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Returns "the default AnimationClipInst instance" referencing this AnimationClipDef definition.
//	That instance is created once when this method is first called on me,
//	and will have its Def field readily set to me.
func (me *AnimationClipDef) DefaultInst() (inst *AnimationClipInst) {
	if inst = defaultAnimationClipInsts[me]; inst == nil {
		inst = me.NewInst()
		defaultAnimationClipInsts[me] = inst
	}
	return
}

//	Creates and returns a new AnimationClipInst instance referencing this AnimationClipDef definition.
//	Any AnimationClipInst created by this method will have its Def field readily set to me.
func (me *AnimationClipDef) NewInst() (inst *AnimationClipInst) {
	inst = &AnimationClipInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct AnimationClipDef
//	according to the current me.DefRef value (by searching AllAnimationClipDefLibs).
//	Then returns me.Def.
//	(Note, every AnimationClipInst's Def is nil initially, unless it was created via AnimationClipDef.NewInst().)
func (me *AnimationClipInst) EnsureDef() *AnimationClipDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.AnimationClipDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibAnimationClipDefs libraries associated by their Id.
	AllAnimationClipDefLibs = LibsAnimationClipDef{}

	//	The "default" LibAnimationClipDefs library for AnimationClipDefs.
	AnimationClipDefs = AllAnimationClipDefLibs.AddNew("")

	defaultAnimationClipInsts = map[*AnimationClipDef]*AnimationClipInst{}
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllAnimationClipDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibAnimationClipDefs contained in AllAnimationClipDefLibs) for the AnimationClipDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) AnimationClipDef() (def *AnimationClipDef) {
	id := me.S()
	for _, lib := range AllAnimationClipDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllAnimationClipDefLibs variable:
//	a hash-table that contains LibAnimationClipDefs libraries associated by their Id.
type LibsAnimationClipDef map[string]*LibAnimationClipDefs

//	Creates a new LibAnimationClipDefs library with the specified Id, adds it to this LibsAnimationClipDef, and returns it.
//	If this LibsAnimationClipDef already contains a LibAnimationClipDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains AnimationClipDefs associated by their Id.
//	To create a new LibAnimationClipDefs library, ONLY use the LibsAnimationClipDef.New() or LibsAnimationClipDef.AddNew() methods.
type LibAnimationClipDefs struct {
	//	Id, Name
	BaseLib

	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*AnimationClipDef
}

func newLibAnimationClipDefs(id string) (me *LibAnimationClipDefs) {
	me = &LibAnimationClipDefs{M: map[string]*AnimationClipDef{}}
	me.BaseLib.init(id)
	return
}

//	Adds the specified AnimationClipDef definition to this LibAnimationClipDefs, and returns it.
//	If this LibAnimationClipDefs already contains a AnimationClipDef definition with the same Id, does nothing and returns nil.
func (me *LibAnimationClipDefs) Add(d *AnimationClipDef) (n *AnimationClipDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new AnimationClipDef definition with the specified Id, adds it to this LibAnimationClipDefs, and returns it.
//	If this LibAnimationClipDefs already contains a AnimationClipDef definition with the specified Id, does nothing and returns nil.
func (me *LibAnimationClipDefs) AddNew(id string) *AnimationClipDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibAnimationClipDefs) Len() int { return len(me.M) }

//	Creates a new AnimationClipDef definition with the specified Id and returns it,
//	but does not add it to this LibAnimationClipDefs.
func (me *LibAnimationClipDefs) New(id string) (def *AnimationClipDef) { def = newAnimationClipDef(id); return }

//	Removes the AnimationClipDef with the specified Id from this LibAnimationClipDefs.
func (me *LibAnimationClipDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibAnimationClipDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibAnimationClipDefs
//	library or its AnimationClipDef definitions. Also called by the global SyncChanges() function.
func (me *LibAnimationClipDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
