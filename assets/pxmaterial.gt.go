package assets

type PxMaterialDef struct {
	BaseDef
	HasTechniques
	TechniqueCommon struct {
		DynamicFriction ScopedFloat
		Restitution     ScopedFloat
		StaticFriction  ScopedFloat
	}
}

func (me *PxMaterialDef) init() {
}

type PxMaterialInst struct {
	BaseInst

	Def *PxMaterialDef
}

func (me *PxMaterialInst) init() {
}

//#begin-gt _definstlib.gt T:PxMaterial

func newPxMaterialDef(id string) (me *PxMaterialDef) {
	me = &PxMaterialDef{}
	me.ID = id
	me.init()
	return
}

/*
//	Creates and returns a new *PxMaterialInst* instance referencing this *PxMaterialDef* definition.
func (me *PxMaterialDef) NewInst(id string) (inst *PxMaterialInst) {
	inst = &PxMaterialInst{Def: me}
	inst.init()
	return
}
*/

var (
	//	A *map* collection that contains *LibPxMaterialDefs* libraries associated by their *ID*.
	AllPxMaterialDefLibs = LibsPxMaterialDef{}

	//	The "default" *LibPxMaterialDefs* library for *PxMaterialDef*s.
	PxMaterialDefs = AllPxMaterialDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllPxMaterialDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllPxMaterialDefLibs* variable: a *map* collection that contains
//	*LibPxMaterialDefs* libraries associated by their *ID*.
type LibsPxMaterialDef map[string]*LibPxMaterialDefs

//	Creates a new *LibPxMaterialDefs* library with the specified *ID*, adds it to this *LibsPxMaterialDef*, and returns it.
//	
//	If this *LibsPxMaterialDef* already contains a *LibPxMaterialDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsPxMaterialDef) AddNew(id string) (lib *LibPxMaterialDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsPxMaterialDef) new(id string) (lib *LibPxMaterialDefs) {
	lib = newLibPxMaterialDefs(id)
	return
}

//	A library that contains *PxMaterialDef*s associated by their *ID*. To create a new *LibPxMaterialDefs* library, ONLY
//	use the *LibsPxMaterialDef.New()* or *LibsPxMaterialDef.AddNew()* methods.
type LibPxMaterialDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*PxMaterialDef
}

func newLibPxMaterialDefs(id string) (me *LibPxMaterialDefs) {
	me = &LibPxMaterialDefs{M: map[string]*PxMaterialDef{}}
	me.ID = id
	return
}

//	Adds the specified *PxMaterialDef* definition to this *LibPxMaterialDefs*, and returns it.
//	
//	If this *LibPxMaterialDefs* already contains a *PxMaterialDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibPxMaterialDefs) Add(d *PxMaterialDef) (n *PxMaterialDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *PxMaterialDef* definition with the specified *ID*, adds it to this *LibPxMaterialDefs*, and returns it.
//	
//	If this *LibPxMaterialDefs* already contains a *PxMaterialDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibPxMaterialDefs) AddNew(id string) *PxMaterialDef { return me.Add(me.New(id)) }

//	Creates a new *PxMaterialDef* definition with the specified *ID* and returns it, but does not add it to this *LibPxMaterialDefs*.
func (me *LibPxMaterialDefs) New(id string) (def *PxMaterialDef) { def = newPxMaterialDef(id); return }

//	Removes the *PxMaterialDef* with the specified *ID* from this *LibPxMaterialDefs*.
func (me *LibPxMaterialDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibPxMaterialDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibPxMaterialDefs* library or its *PxMaterialDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibPxMaterialDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
