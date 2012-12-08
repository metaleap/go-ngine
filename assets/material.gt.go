package assets

//	Defines the equations necessary for the visual appearance of geometry and screen-space image processing.
type MaterialDef struct {
	BaseDef
}

func (me *MaterialDef) init() {
}

//	An instance referencing a material definition.
type MaterialInst struct {
	BaseInst

	//	The material definition referenced by this instance.
	Def *MaterialDef
}

func (me *MaterialInst) init() {
}

//#begin-gt _definstlib.gt T:Material

func newMaterialDef(id string) (me *MaterialDef) {
	me = &MaterialDef{}
	me.ID = id
	me.init()
	return
}

//	Creates and returns a new *MaterialInst* instance referencing this *MaterialDef* definition.
func (me *MaterialDef) NewInst(id string) (inst *MaterialInst) {
	inst = &MaterialInst{Def: me}
	inst.init()
	return
}

var (
	//	A *map* collection that contains *LibMaterialDefs* libraries associated by their *ID*.
	AllMaterialDefLibs = LibsMaterialDef{}

	//	The "default" *LibMaterialDefs* library for *MaterialDef*s.
	MaterialDefs = AllMaterialDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllMaterialDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllMaterialDefLibs* variable: a *map* collection that contains
//	*LibMaterialDefs* libraries associated by their *ID*.
type LibsMaterialDef map[string]*LibMaterialDefs

//	Creates a new *LibMaterialDefs* library with the specified *ID*, adds it to this *LibsMaterialDef*, and returns it.
//	
//	If this *LibsMaterialDef* already contains a *LibMaterialDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsMaterialDef) AddNew(id string) (lib *LibMaterialDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsMaterialDef) new(id string) (lib *LibMaterialDefs) {
	lib = newLibMaterialDefs(id)
	return
}

//	A library that contains *MaterialDef*s associated by their *ID*. To create a new *LibMaterialDefs* library, ONLY
//	use the *LibsMaterialDef.New()* or *LibsMaterialDef.AddNew()* methods.
type LibMaterialDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*MaterialDef
}

func newLibMaterialDefs(id string) (me *LibMaterialDefs) {
	me = &LibMaterialDefs{M: map[string]*MaterialDef{}}
	me.ID = id
	return
}

//	Adds the specified *MaterialDef* definition to this *LibMaterialDefs*, and returns it.
//	
//	If this *LibMaterialDefs* already contains a *MaterialDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibMaterialDefs) Add(d *MaterialDef) (n *MaterialDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *MaterialDef* definition with the specified *ID*, adds it to this *LibMaterialDefs*, and returns it.
//	
//	If this *LibMaterialDefs* already contains a *MaterialDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibMaterialDefs) AddNew(id string) *MaterialDef { return me.Add(me.New(id)) }

//	Creates a new *MaterialDef* definition with the specified *ID* and returns it, but does not add it to this *LibMaterialDefs*.
func (me *LibMaterialDefs) New(id string) (def *MaterialDef) { def = newMaterialDef(id); return }

//	Removes the *MaterialDef* with the specified *ID* from this *LibMaterialDefs*.
func (me *LibMaterialDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibMaterialDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibMaterialDefs* library or its *MaterialDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibMaterialDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
