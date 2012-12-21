package assets

import (
	unum "github.com/metaleap/go-util/num"
)

//	Used to declare skinning joints or morph targets.
type ControllerInputs struct {
	//	Extras
	HasExtras
	//	Inputs
	HasInputs
}

//	Describes the data required to blend between sets of static meshes.
type ControllerMorph struct {
	//	Sources
	HasSources
	//	Which blending method to use: true for relative blending, false for normalized blending.
	Relative bool
	//	Refers to the Geometry that describes the base mesh.
	Source RefId
	//	Input meshes (morph targets) to be blended.
	Targets ControllerInputs
}

//	Constructor
func NewControllerMorph() (me *ControllerMorph) {
	me = &ControllerMorph{}
	me.Sources = Sources{}
	return
}

//	Contains vertex and primitive information sufficient to describe blend-weight skinning.
type ControllerSkin struct {
	//	Sources
	HasSources
	//	Provides extra information about the position and orientation of the base mesh before binding.
	BindShapeMatrix unum.Mat4
	//	Describes a per-vertex combination of joints and weights used in this skin.
	//	An index of –1 into the array of joints refers to the bind shape.
	//	Weights should be normalized before use.
	VertexWeights IndexedInputs
	//	Aggregates the per-joint information needed for this skin.
	Joints ControllerInputs
	//	Refers to the base mesh (a static mesh or a morphed mesh).
	//	This also provides the bind-shape of the skinned mesh.
	Source RefId
}

//	Constructor
func NewControllerSkin() (me *ControllerSkin) {
	me = &ControllerSkin{}
	me.BindShapeMatrix.Identity()
	me.Sources = Sources{}
	return
}

//	Defines generic control information for dynamic content.
type ControllerDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	If set, Skin must be nil; declares this a mesh-morphing controller that deforms meshes and blends them.
	Morph *ControllerMorph
	//	If set, Morph must be nil; declares this a vertex-skinning controller that transforms vertices
	//	based on weighted influences to produce a smoothly changing mesh.
	Skin *ControllerSkin
}

//	Initialization
func (me *ControllerDef) Init() {
}

//	Instantiates a controller resource.
type ControllerInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *ControllerDef
	//	Binds a specific material to this controller instantiation.
	BindMaterial *MaterialBinding
	//	Indicates where a Skin controller is to start to search for the joint nodes it needs.
	//	This element is meaningless for Morph controllers.
	SkinSkeletons []string
}

//	Initialization
func (me *ControllerInst) Init() {
}

//#begin-gt _definstlib.gt T:Controller

func newControllerDef(id string) (me *ControllerDef) {
	me = &ControllerDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new ControllerInst instance referencing this ControllerDef definition.
//	Any ControllerInst created by this method will have its Def field readily set to me.
func (me *ControllerDef) NewInst() (inst *ControllerInst) {
	inst = &ControllerInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct ControllerDef
//	according to the current me.DefRef value (by searching AllControllerDefLibs).
//	Then returns me.Def.
//	(Note, every ControllerInst's Def is nil initially, unless it was created via ControllerDef.NewInst().)
func (me *ControllerInst) EnsureDef() *ControllerDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.ControllerDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibControllerDefs libraries associated by their Id.
	AllControllerDefLibs = LibsControllerDef{}

	//	The "default" LibControllerDefs library for ControllerDefs.
	ControllerDefs = AllControllerDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllControllerDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibControllerDefs contained in AllControllerDefLibs) for the ControllerDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) ControllerDef() (def *ControllerDef) {
	id := me.S()
	for _, lib := range AllControllerDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllControllerDefLibs variable:
//	a hash-table that contains LibControllerDefs libraries associated by their Id.
type LibsControllerDef map[string]*LibControllerDefs

//	Creates a new LibControllerDefs library with the specified Id, adds it to this LibsControllerDef, and returns it.
//	If this LibsControllerDef already contains a LibControllerDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains ControllerDefs associated by their Id.
//	To create a new LibControllerDefs library, ONLY use the LibsControllerDef.New() or LibsControllerDef.AddNew() methods.
type LibControllerDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*ControllerDef
}

func newLibControllerDefs(id string) (me *LibControllerDefs) {
	me = &LibControllerDefs{M: map[string]*ControllerDef{}}
	me.Id = id
	return
}

//	Adds the specified ControllerDef definition to this LibControllerDefs, and returns it.
//	If this LibControllerDefs already contains a ControllerDef definition with the same Id, does nothing and returns nil.
func (me *LibControllerDefs) Add(d *ControllerDef) (n *ControllerDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new ControllerDef definition with the specified Id, adds it to this LibControllerDefs, and returns it.
//	If this LibControllerDefs already contains a ControllerDef definition with the specified Id, does nothing and returns nil.
func (me *LibControllerDefs) AddNew(id string) *ControllerDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibControllerDefs) Len() int { return len(me.M) }

//	Creates a new ControllerDef definition with the specified Id and returns it,
//	but does not add it to this LibControllerDefs.
func (me *LibControllerDefs) New(id string) (def *ControllerDef) { def = newControllerDef(id); return }

//	Removes the ControllerDef with the specified Id from this LibControllerDefs.
func (me *LibControllerDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Returns a GetRefSidResolver that looks up and yields the ControllerDef with the specified Id.
func (me *LibControllerDefs) ResolverGetter() GetRefSidResolver {
	return func(id string) RefSidResolver {
		return nil // me.M[id]
	}
}

//	Signals to the core package (or your custom package) that changes have been made to this LibControllerDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibControllerDefs
//	library or its ControllerDef definitions. Also called by the global SyncChanges() function.
func (me *LibControllerDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
