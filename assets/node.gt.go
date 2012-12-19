package assets

//	Used to recursively define hierarchies of nodes.
type ChildNode struct {
	//	If set, Inst must be nil. An inline node definition.
	Def *NodeDef
	//	If set, Def must be nil. Instantiates a previously defined node.
	Inst *NodeInst
}

//	Declares a point of interest in a scene.
type NodeDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sid
	HasSid
	//	Indicates whether this node is a joint for a skin controller.
	IsSkinJoint bool
	//	The names of the layers to which this node belongs.
	Layers Layers
	//	Any combination of zero or more transformations of any type.
	Transforms []*Transform
	//	Content resources participating in this node.
	Insts struct {
		//	Cameras participating in this node.
		Camera []*CameraInst
		//	Controllers participating in this node.
		Controller []*ControllerInst
		//	Geometries participating in this node.
		Geometry []*GeometryInst
		//	Lights participating in this node.
		Light []*LightInst
	}
	//	Child nodes to recursively define a hierarchy.
	Nodes []ChildNode
}

//	Initialization
func (me *NodeDef) Init() {
	me.Layers = Layers{}
}

//	Instantiates a node resource.
type NodeInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	Optional. The mechanism and use of this attribute is application-defined. For example, it can be used for bounding boxes or level of detail.
	Proxy RefId
}

//	Initialization
func (me *NodeInst) Init() {
}

//#begin-gt _definstlib.gt T:Node

func newNodeDef(id string) (me *NodeDef) {
	me = &NodeDef{}
	me.Id = id
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
	//	A *map* collection that contains *LibNodeDefs* libraries associated by their *Id*.
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
//	*LibNodeDefs* libraries associated by their *Id*.
type LibsNodeDef map[string]*LibNodeDefs

//	Creates a new *LibNodeDefs* library with the specified *Id*, adds it to this *LibsNodeDef*, and returns it.
//	
//	If this *LibsNodeDef* already contains a *LibNodeDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *NodeDef*s associated by their *Id*. To create a new *LibNodeDefs* library, ONLY
//	use the *LibsNodeDef.New()* or *LibsNodeDef.AddNew()* methods.
type LibNodeDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*NodeDef
}

func newLibNodeDefs(id string) (me *LibNodeDefs) {
	me = &LibNodeDefs{M: map[string]*NodeDef{}}
	me.Id = id
	return
}

//	Adds the specified *NodeDef* definition to this *LibNodeDefs*, and returns it.
//	
//	If this *LibNodeDefs* already contains a *NodeDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibNodeDefs) Add(d *NodeDef) (n *NodeDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *NodeDef* definition with the specified *Id*, adds it to this *LibNodeDefs*, and returns it.
//	
//	If this *LibNodeDefs* already contains a *NodeDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibNodeDefs) AddNew(id string) *NodeDef { return me.Add(me.New(id)) }

//	Short-hand for len(lib.M)
func (me *LibNodeDefs) Len() int { return len(me.M) }

//	Creates a new *NodeDef* definition with the specified *Id* and returns it, but does not add it to this *LibNodeDefs*.
func (me *LibNodeDefs) New(id string) (def *NodeDef) { def = newNodeDef(id); return }

//	Removes the *NodeDef* with the specified *Id* from this *LibNodeDefs*.
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
