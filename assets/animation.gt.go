package assets

const (
	//	The before and after behaviors are not defined.
	ANIM_SAMPLER_BEHAVIOR_UNDEFINED = 0
	//	The value for the first (PreBehavior) or last (PostBehavior) is returned.
	ANIM_SAMPLER_BEHAVIOR_CONSTANT = iota
	//	The key is mapped in the [first_key , last_key] interval so that the animation cycles.
	ANIM_SAMPLER_BEHAVIOR_CYCLE = iota
	//	The animation continues indefinitely.
	ANIM_SAMPLER_BEHAVIOR_CYCLE_RELATIVE = iota
	//	The value follows the line given by the last two keys in the sample.
	ANIM_SAMPLER_BEHAVIOR_GRADIENT = iota
	//	The key is mapped in the [first_key , last_key] interval so that the animation oscillates.
	ANIM_SAMPLER_BEHAVIOR_OSCILLATE = iota
)

//	Declares an output channel of an animation.
type AnimationChannel struct {
	//	Refers to the source animation sampler.
	Source RefId
	//	Refers to the resource property bound to the output of the sampler.
	Target RefSid
}

//	Declares an interpolation sampling function for an animation.
type AnimationSampler struct {
	//	Id
	HasId
	//	Inputs. These describe sampling points.
	//	At least one of the Inputs must have its Semantic set to "INTERPOLATION".
	HasInputs
	//	Indicates what the sampled value should be before the first key.
	//	Valid values are the ANIM_SAMPLER_BEHAVIOR_* enumerated constants.
	PreBehavior int
	//	Indicates what the sampled value should be after the last key.
	//	Valid values are the ANIM_SAMPLER_BEHAVIOR_* enumerated constants.
	PostBehavior int
}

//	Categorizes the declaration of animation information.
type AnimationDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sources
	HasSources
	//	Allows the formation of a hierarchy of related animations.
	AnimationDefs []*AnimationDef
	//	Describes output channels for the animation.
	Channels []*AnimationChannel
	//	Describes the interpolation sampling functions for the animation.
	Samplers []*AnimationSampler
}

//	Initialization
func (me *AnimationDef) Init() {
	me.Sources = Sources{}
}

//	Instantiates an Animation resource.
type AnimationInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default and meant to be set ONLY by the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *AnimationDef
}

//	Initialization
func (me *AnimationInst) Init() {
}

//#begin-gt _definstlib.gt T:Animation

func newAnimationDef(id string) (me *AnimationDef) {
	me = &AnimationDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new AnimationInst instance referencing this AnimationDef definition.
func (me *AnimationDef) NewInst() (inst *AnimationInst) {
	inst = &AnimationInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is dirty or me.Def is nil, sets me.Def to the correct AnimationDef
//	according to the current me.DefRef value (by searching AllAnimationDefLibs).
//	Then returns me.Def.
func (me *AnimationInst) EnsureDef() *AnimationDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.AnimationDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibAnimationDefs libraries associated by their Id.
	AllAnimationDefLibs = LibsAnimationDef{}

	//	The "default" LibAnimationDefs library for AnimationDefs.
	AnimationDefs = AllAnimationDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllAnimationDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibAnimationDefs contained in AllAnimationDefLibs) for the AnimationDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) AnimationDef() (def *AnimationDef) {
	id := me.S()
	for _, lib := range AllAnimationDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllAnimationDefLibs variable:
//	a hash-table that contains LibAnimationDefs libraries associated by their Id.
type LibsAnimationDef map[string]*LibAnimationDefs

//	Creates a new LibAnimationDefs library with the specified Id, adds it to this LibsAnimationDef, and returns it.
//	If this LibsAnimationDef already contains a LibAnimationDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains AnimationDefs associated by their Id.
//	To create a new LibAnimationDefs library, ONLY use the LibsAnimationDef.New() or LibsAnimationDef.AddNew() methods.
type LibAnimationDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*AnimationDef
}

func newLibAnimationDefs(id string) (me *LibAnimationDefs) {
	me = &LibAnimationDefs{M: map[string]*AnimationDef{}}
	me.Id = id
	return
}

//	Adds the specified AnimationDef definition to this LibAnimationDefs, and returns it.
//	If this LibAnimationDefs already contains a AnimationDef definition with the same Id, does nothing and returns nil.
func (me *LibAnimationDefs) Add(d *AnimationDef) (n *AnimationDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new AnimationDef definition with the specified Id, adds it to this LibAnimationDefs, and returns it.
//	If this LibAnimationDefs already contains a AnimationDef definition with the specified Id, does nothing and returns nil.
func (me *LibAnimationDefs) AddNew(id string) *AnimationDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibAnimationDefs) Len() int { return len(me.M) }

//	Creates a new AnimationDef definition with the specified Id and returns it,
//	but does not add it to this LibAnimationDefs.
func (me *LibAnimationDefs) New(id string) (def *AnimationDef) { def = newAnimationDef(id); return }

//	Removes the AnimationDef with the specified Id from this LibAnimationDefs.
func (me *LibAnimationDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibAnimationDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibAnimationDefs
//	library or its AnimationDef definitions. Also called by the global SyncChanges() function.
func (me *LibAnimationDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
