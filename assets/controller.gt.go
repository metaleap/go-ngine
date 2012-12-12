package assets

import (
	unum "github.com/metaleap/go-util/num"
)

type ControllerInputs struct {
	HasExtras
	HasInputs
}

type ControllerMorph struct {
	HasSources
	Relative bool
	Source   string
	Targets  ControllerInputs
}

func NewControllerMorph() (me *ControllerMorph) {
	me = &ControllerMorph{}
	me.Sources = Sources{}
	return
}

type ControllerSkin struct {
	HasSources
	BindShapeMatrix unum.Mat4
	VertexWeights   IndexedInputsV
	Joints          ControllerInputs
	Source          string
}

func NewControllerSkin() (me *ControllerSkin) {
	me = &ControllerSkin{}
	me.BindShapeMatrix.Identity()
	me.Sources = Sources{}
	return
}

type ControllerDef struct {
	BaseDef
	Morph *ControllerMorph
	Skin  *ControllerSkin
}

func (me *ControllerDef) Init() {
}

type ControllerInst struct {
	BaseInst
	BindMaterial *BindMaterial
	Skeletons    []string
}

func (me *ControllerInst) Init() {
}

//#begin-gt _definstlib.gt T:Controller

func newControllerDef(id string) (me *ControllerDef) {
	me = &ControllerDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *ControllerInst* instance referencing this *ControllerDef* definition.
func (me *ControllerDef) NewInst(id string) (inst *ControllerInst) {
	inst = &ControllerInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibControllerDefs* libraries associated by their *Id*.
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
//	*LibControllerDefs* libraries associated by their *Id*.
type LibsControllerDef map[string]*LibControllerDefs

//	Creates a new *LibControllerDefs* library with the specified *Id*, adds it to this *LibsControllerDef*, and returns it.
//	
//	If this *LibsControllerDef* already contains a *LibControllerDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *ControllerDef*s associated by their *Id*. To create a new *LibControllerDefs* library, ONLY
//	use the *LibsControllerDef.New()* or *LibsControllerDef.AddNew()* methods.
type LibControllerDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*ControllerDef
}

func newLibControllerDefs(id string) (me *LibControllerDefs) {
	me = &LibControllerDefs{M: map[string]*ControllerDef{}}
	me.Id = id
	return
}

//	Adds the specified *ControllerDef* definition to this *LibControllerDefs*, and returns it.
//	
//	If this *LibControllerDefs* already contains a *ControllerDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibControllerDefs) Add(d *ControllerDef) (n *ControllerDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *ControllerDef* definition with the specified *Id*, adds it to this *LibControllerDefs*, and returns it.
//	
//	If this *LibControllerDefs* already contains a *ControllerDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibControllerDefs) AddNew(id string) *ControllerDef { return me.Add(me.New(id)) }

//	Creates a new *ControllerDef* definition with the specified *Id* and returns it, but does not add it to this *LibControllerDefs*.
func (me *LibControllerDefs) New(id string) (def *ControllerDef) { def = newControllerDef(id); return }

//	Removes the *ControllerDef* with the specified *Id* from this *LibControllerDefs*.
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
