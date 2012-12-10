package assets

type NodeDef struct {
	BaseDef
	HasSid
	IsJoint    bool
	Layers     Layers
	Transforms []*Transform
	Insts      struct {
		Camera     []*CameraInst
		Controller []*ControllerInst
		Geometry   []*GeometryInst
		Light      []*LightInst
	}
	Nodes struct {
		Defs  []*NodeDef
		Insts []*NodeInst
	}
}

func (me *NodeDef) Init() {
	me.Layers = Layers{}
}

type NodeInst struct {
	BaseInst
}

func (me *NodeInst) init() {
}

//#begin-gt _definstlib.gt T:Node

func newNodeDef(id string) (me *NodeDef) {
	me = &NodeDef{}
	me.ID = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *NodeInst* instance referencing this *NodeDef* definition.
func (me *NodeDef) NewInst(id string) (inst *NodeInst) {
	inst = &NodeInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibNodeDefs* libraries associated by their *ID*.
	AllNodeDefLibs = LibsNodeDef{}

	//	The "default" *LibNodeDefs* library for *NodeDef*s.
	NodeDefs = AllNodeDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllNodeDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllNodeDefLibs* variable: a *map* collection that contains
//	*LibNodeDefs* libraries associated by their *ID*.
type LibsNodeDef map[string]*LibNodeDefs

//	Creates a new *LibNodeDefs* library with the specified *ID*, adds it to this *LibsNodeDef*, and returns it.
//	
//	If this *LibsNodeDef* already contains a *LibNodeDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsNodeDef) AddNew(id string) (lib *LibNodeDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsNodeDef) new(id string) (lib *LibNodeDefs) {
	lib = newLibNodeDefs(id)
	return
}

//	A library that contains *NodeDef*s associated by their *ID*. To create a new *LibNodeDefs* library, ONLY
//	use the *LibsNodeDef.New()* or *LibsNodeDef.AddNew()* methods.
type LibNodeDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*NodeDef
}

func newLibNodeDefs(id string) (me *LibNodeDefs) {
	me = &LibNodeDefs{M: map[string]*NodeDef{}}
	me.ID = id
	return
}

//	Adds the specified *NodeDef* definition to this *LibNodeDefs*, and returns it.
//	
//	If this *LibNodeDefs* already contains a *NodeDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibNodeDefs) Add(d *NodeDef) (n *NodeDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *NodeDef* definition with the specified *ID*, adds it to this *LibNodeDefs*, and returns it.
//	
//	If this *LibNodeDefs* already contains a *NodeDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibNodeDefs) AddNew(id string) *NodeDef { return me.Add(me.New(id)) }

//	Creates a new *NodeDef* definition with the specified *ID* and returns it, but does not add it to this *LibNodeDefs*.
func (me *LibNodeDefs) New(id string) (def *NodeDef) { def = newNodeDef(id); return }

//	Removes the *NodeDef* with the specified *ID* from this *LibNodeDefs*.
func (me *LibNodeDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibNodeDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibNodeDefs* library or its *NodeDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibNodeDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
