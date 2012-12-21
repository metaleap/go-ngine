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
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *NodeDef
	//	Optional. The mechanism and use of this attribute is application-defined.
	//	For example, it can be used for bounding boxes or level of detail.
	Proxy RefId
}

//	Initialization
func (me *NodeInst) Init() {
}

//#begin-gt _definstlib.gt T:Node

func newNodeDef(id string) (me *NodeDef) {
	me = &NodeDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new NodeInst instance referencing this NodeDef definition.
//	Any NodeInst created by this method will have its Def field readily set to me.
func (me *NodeDef) NewInst() (inst *NodeInst) {
	inst = &NodeInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct NodeDef
//	according to the current me.DefRef value (by searching AllNodeDefLibs).
//	Then returns me.Def.
//	(Note, every NodeInst's Def is nil initially, unless it was created via NodeDef.NewInst().)
func (me *NodeInst) EnsureDef() *NodeDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.NodeDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibNodeDefs libraries associated by their Id.
	AllNodeDefLibs = LibsNodeDef{}

	//	The "default" LibNodeDefs library for NodeDefs.
	NodeDefs = AllNodeDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllNodeDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibNodeDefs contained in AllNodeDefLibs) for the NodeDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) NodeDef() (def *NodeDef) {
	id := me.S()
	for _, lib := range AllNodeDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllNodeDefLibs variable:
//	a hash-table that contains LibNodeDefs libraries associated by their Id.
type LibsNodeDef map[string]*LibNodeDefs

//	Creates a new LibNodeDefs library with the specified Id, adds it to this LibsNodeDef, and returns it.
//	If this LibsNodeDef already contains a LibNodeDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains NodeDefs associated by their Id.
//	To create a new LibNodeDefs library, ONLY use the LibsNodeDef.New() or LibsNodeDef.AddNew() methods.
type LibNodeDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*NodeDef
}

func newLibNodeDefs(id string) (me *LibNodeDefs) {
	me = &LibNodeDefs{M: map[string]*NodeDef{}}
	me.Id = id
	return
}

//	Adds the specified NodeDef definition to this LibNodeDefs, and returns it.
//	If this LibNodeDefs already contains a NodeDef definition with the same Id, does nothing and returns nil.
func (me *LibNodeDefs) Add(d *NodeDef) (n *NodeDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new NodeDef definition with the specified Id, adds it to this LibNodeDefs, and returns it.
//	If this LibNodeDefs already contains a NodeDef definition with the specified Id, does nothing and returns nil.
func (me *LibNodeDefs) AddNew(id string) *NodeDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibNodeDefs) Len() int { return len(me.M) }

//	Creates a new NodeDef definition with the specified Id and returns it,
//	but does not add it to this LibNodeDefs.
func (me *LibNodeDefs) New(id string) (def *NodeDef) { def = newNodeDef(id); return }

//	Removes the NodeDef with the specified Id from this LibNodeDefs.
func (me *LibNodeDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Returns a GetRefSidResolver that looks up and yields the NodeDef with the specified Id.
func (me *LibNodeDefs) ResolverGetter() GetRefSidResolver {
	return func(id string) RefSidResolver {
		return nil // me.M[id]
	}
}

//	Signals to the core package (or your custom package) that changes have been made to this LibNodeDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibNodeDefs
//	library or its NodeDef definitions. Also called by the global SyncChanges() function.
func (me *LibNodeDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
