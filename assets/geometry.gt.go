package assets

const (
	//	Organizes vertices into individual lines.
	GEOMETRY_PRIMITIVE_TYPE_LINES = 0x0001
	//	Organizes vertices into connected line-strips.
	GEOMETRY_PRIMITIVE_TYPE_LINE_STRIPS = 0x0003
	//	Organizes vertices into individual polygons that may contain holes.
	GEOMETRY_PRIMITIVE_TYPE_POLYGONS = 2
	//	Organizes vertices into individual polygons that cannot contain holes.
	GEOMETRY_PRIMITIVE_TYPE_POLYLIST = 7
	//	Organizes vertices into individual triangles.
	GEOMETRY_PRIMITIVE_TYPE_TRIANGLES = 0x0004
	//	Organizes vertices into fan-connected triangles.
	GEOMETRY_PRIMITIVE_TYPE_TRIFANS = 0x0006
	//	Organizes vertices into strip-connected triangles.
	GEOMETRY_PRIMITIVE_TYPE_TRISTRIPS = 0x0005
)

//	Describes the control vertices of a spline.
type GeometryControlVertices struct {
	//	Extras
	HasExtras
	//	Inputs
	HasInputs
}

//	Describes basic geometric meshes using vertex and primitive information.
type GeometryMesh struct {
	//	Extras
	HasExtras
	//	Sources
	HasSources
	//	Refers to a GeometryDef described by a GeometryMesh.
	//	If specified, compute the convex hull of the specified mesh.
	ConvexHullOf RefId
	//	Describes the mesh-vertex attributes and establishes their topological identity.
	//	Required if ConvexHullOf is empty.
	Vertices *GeometryVertices
	//	Geometric primitives, which assemble values from the inputs into vertex attribute data.
	Primitives []*GeometryPrimitives
}

//	Constructor
func NewGeometryMesh() (me *GeometryMesh) {
	me = &GeometryMesh{}
	me.Sources = Sources{}
	return
}

//	Geometric primitives, which assemble values from inputs into vertex attribute data.
type GeometryPrimitives struct {
	//	Extras
	HasExtras
	//	Name
	HasName
	//	When at least one input is present, one input must specify its Semantic as "VERTEX".
	IndexedInputs
	//	Must be one of the GEOMETRY_PRIMITIVE_TYPE_* enumerated constants.
	Type int
	//	Declares a symbol for a material. This symbol is bound to a material at the time of instantiation.
	//	Optional. If not specified then the lighting and shading results are application defined.
	Material string
	//	If Type is GEOMETRY_PRIMITIVE_TYPE_POLYGONS, describes the polygons that contain one or more holes.
	PolyHoles []*GeometryPolygonHole
}

//	Describes a polygon that contains one or more holes.
type GeometryPolygonHole struct {
	//	Specifies the vertex attributes (indices) for an individual polygon.
	Indices []uint64
	//	Specifies the indices of a hole in the polygon specified by Indices.
	Holes [][]uint64
}

//	Describes a multisegment spline with control vertex and segment information.
type GeometrySpline struct {
	//	Extras
	HasExtras
	//	Sources
	HasSources
	//	Whether there is a segment connecting the first and last control vertices.
	//	The default is false, indicating that the spline is open.
	Closed bool
	//	Describes the control vertices of the spline.
	ControlVertices GeometryControlVertices
}

//	Constructor
func NewGeometrySpline() (me *GeometrySpline) {
	me = &GeometrySpline{}
	me.Sources = Sources{}
	return
}

//	Declares the attributes and identity of mesh-vertices.
//	The mesh-vertices represent the position (identity) of the vertices comprising the mesh
//	and other vertex attributes that are invariant to tessellation.
type GeometryVertices struct {
	//	Id
	HasId
	//	Name
	HasName
	//	Extras
	HasExtras
	//	Inputs
	HasInputs
}

//	Describes the visual shape and appearance of an object in a scene.
type GeometryDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	If set, Mesh and Spline must be nil, and the GeometryDef is described by this B-rep structure.
	Brep *GeometryBrep
	//	If set, Brep and Spline must be nil, and the GeometryDef is described by this mesh structure.
	Mesh *GeometryMesh
	//	If set, Mesh and Brep must be nil, and the GeometryDef is described by this multi-segment spline.
	Spline *GeometrySpline
}

//	Initialization
func (me *GeometryDef) Init() {
}

//	Instantiates a geometry resource.
type GeometryInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default and meant to be set ONLY by the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *GeometryDef
	//	Binds material symbols to material instances. This allows a single geometry
	//	to be instantiated into a scene multiple times each with a different appearance.
	MaterialBinding *MaterialBinding
}

//	Initialization
func (me *GeometryInst) Init() {
}

//#begin-gt _definstlib.gt T:Geometry

func newGeometryDef(id string) (me *GeometryDef) {
	me = &GeometryDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new GeometryInst instance referencing this GeometryDef definition.
func (me *GeometryDef) NewInst() (inst *GeometryInst) {
	inst = &GeometryInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is dirty or me.Def is nil, sets me.Def to the correct GeometryDef
//	according to the current me.DefRef value (by searching AllGeometryDefLibs).
//	Then returns me.Def.
func (me *GeometryInst) EnsureDef() *GeometryDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.GeometryDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibGeometryDefs libraries associated by their Id.
	AllGeometryDefLibs = LibsGeometryDef{}

	//	The "default" LibGeometryDefs library for GeometryDefs.
	GeometryDefs = AllGeometryDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllGeometryDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) GeometryDef() (def *GeometryDef) {
	id := me.S()
	for _, lib := range AllGeometryDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllGeometryDefLibs variable:
//	a hash-table that contains LibGeometryDefs libraries associated by their Id.
type LibsGeometryDef map[string]*LibGeometryDefs

//	Creates a new LibGeometryDefs library with the specified Id, adds it to this LibsGeometryDef, and returns it.
//	If this LibsGeometryDef already contains a LibGeometryDefs library with the specified Id, does nothing and returns nil.
func (me LibsGeometryDef) AddNew(id string) (lib *LibGeometryDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsGeometryDef) new(id string) (lib *LibGeometryDefs) {
	lib = newLibGeometryDefs(id)
	return
}

//	A library that contains GeometryDefs associated by their Id.
//	To create a new LibGeometryDefs library, ONLY use the LibsGeometryDef.New() or LibsGeometryDef.AddNew() methods.
type LibGeometryDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*GeometryDef
}

func newLibGeometryDefs(id string) (me *LibGeometryDefs) {
	me = &LibGeometryDefs{M: map[string]*GeometryDef{}}
	me.Id = id
	return
}

//	Adds the specified GeometryDef definition to this LibGeometryDefs, and returns it.
//	If this LibGeometryDefs already contains a GeometryDef definition with the same Id, does nothing and returns nil.
func (me *LibGeometryDefs) Add(d *GeometryDef) (n *GeometryDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new GeometryDef definition with the specified Id, adds it to this LibGeometryDefs, and returns it.
//	If this LibGeometryDefs already contains a GeometryDef definition with the specified Id, does nothing and returns nil.
func (me *LibGeometryDefs) AddNew(id string) *GeometryDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibGeometryDefs) Len() int { return len(me.M) }

//	Creates a new GeometryDef definition with the specified Id and returns it,
//	but does not add it to this LibGeometryDefs.
func (me *LibGeometryDefs) New(id string) (def *GeometryDef) { def = newGeometryDef(id); return }

//	Removes the GeometryDef with the specified Id from this LibGeometryDefs.
func (me *LibGeometryDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibGeometryDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibGeometryDefs
//	library or its GeometryDef definitions. Also called by the global SyncChanges() function.
func (me *LibGeometryDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
