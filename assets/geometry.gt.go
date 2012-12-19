package assets

const (
	GEOMETRY_PRIMITIVE_TYPE_LINES       = 0x0001
	GEOMETRY_PRIMITIVE_TYPE_LINE_STRIPS = 0x0003
	GEOMETRY_PRIMITIVE_TYPE_POLYGONS    = 2
	GEOMETRY_PRIMITIVE_TYPE_POLYLIST    = 7
	GEOMETRY_PRIMITIVE_TYPE_TRIANGLES   = 0x0004
	GEOMETRY_PRIMITIVE_TYPE_TRIFANS     = 0x0006
	GEOMETRY_PRIMITIVE_TYPE_TRISTRIPS   = 0x0005
)

type GeometryControlVertices struct {
	HasExtras
	HasInputs
}

type GeometryMesh struct {
	HasExtras
	HasSources
	ConvexHullOf string
	Vertices     *GeometryVertices
	Primitives   []*GeometryPrimitives
}

func NewGeometryMesh() (me *GeometryMesh) {
	me = &GeometryMesh{}
	me.Sources = Sources{}
	return
}

type GeometryPrimitives struct {
	HasExtras
	HasName
	IndexedInputs
	Material  string
	Type      int
	PolyHoles []*GeometryPolygonHole
}

type GeometryPolygonHole struct {
	Indices []uint64
	Holes   [][]uint64
}

type GeometrySpline struct {
	HasExtras
	HasSources
	Closed          bool
	ControlVertices GeometryControlVertices
}

func NewGeometrySpline() (me *GeometrySpline) {
	me = &GeometrySpline{}
	me.Sources = Sources{}
	return
}

type GeometryVertices struct {
	HasId
	HasName
	HasExtras
	HasInputs
}

type GeometryDef struct {
	BaseDef
	Brep   *GeometryBrep
	Mesh   *GeometryMesh
	Spline *GeometrySpline
}

func (me *GeometryDef) Init() {
}

type GeometryInst struct {
	BaseInst
	BindMaterial *MaterialBinding
}

func (me *GeometryInst) Init() {
}

//#begin-gt _definstlib.gt T:Geometry

func newGeometryDef(id string) (me *GeometryDef) {
	me = &GeometryDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *GeometryInst* instance referencing this *GeometryDef* definition.
func (me *GeometryDef) NewInst(id string) (inst *GeometryInst) {
	inst = &GeometryInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibGeometryDefs* libraries associated by their *Id*.
	AllGeometryDefLibs = LibsGeometryDef{}

	//	The "default" *LibGeometryDefs* library for *GeometryDef*s.
	GeometryDefs = AllGeometryDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllGeometryDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllGeometryDefLibs* variable: a *map* collection that contains
//	*LibGeometryDefs* libraries associated by their *Id*.
type LibsGeometryDef map[string]*LibGeometryDefs

//	Creates a new *LibGeometryDefs* library with the specified *Id*, adds it to this *LibsGeometryDef*, and returns it.
//	
//	If this *LibsGeometryDef* already contains a *LibGeometryDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *GeometryDef*s associated by their *Id*. To create a new *LibGeometryDefs* library, ONLY
//	use the *LibsGeometryDef.New()* or *LibsGeometryDef.AddNew()* methods.
type LibGeometryDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*GeometryDef
}

func newLibGeometryDefs(id string) (me *LibGeometryDefs) {
	me = &LibGeometryDefs{M: map[string]*GeometryDef{}}
	me.Id = id
	return
}

//	Adds the specified *GeometryDef* definition to this *LibGeometryDefs*, and returns it.
//	
//	If this *LibGeometryDefs* already contains a *GeometryDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibGeometryDefs) Add(d *GeometryDef) (n *GeometryDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *GeometryDef* definition with the specified *Id*, adds it to this *LibGeometryDefs*, and returns it.
//	
//	If this *LibGeometryDefs* already contains a *GeometryDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibGeometryDefs) AddNew(id string) *GeometryDef { return me.Add(me.New(id)) }

//	Short-hand for len(lib.M)
func (me *LibGeometryDefs) Len() int { return len(me.M) }

//	Creates a new *GeometryDef* definition with the specified *Id* and returns it, but does not add it to this *LibGeometryDefs*.
func (me *LibGeometryDefs) New(id string) (def *GeometryDef) { def = newGeometryDef(id); return }

//	Removes the *GeometryDef* with the specified *Id* from this *LibGeometryDefs*.
func (me *LibGeometryDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibGeometryDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibGeometryDefs* library or its *GeometryDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibGeometryDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
