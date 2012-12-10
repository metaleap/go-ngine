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
	Inputs []*Input
}

type GeometryMesh struct {
	HasExtras
	ConvexHullOf string
	Sources      Sources
	Vertices     *GeometryVertices
	Primitives   []*GeometryPrimitives
}

func NewGeometryMesh() (me *GeometryMesh) {
	me = &GeometryMesh{Sources: Sources{}}
	return
}

type GeometryPrimitives struct {
	HasExtras
	HasName
	Count     uint64
	Material  string
	Inputs    []*InputShared
	Indices   []int64
	Type      int
	PolyHoles []*GeometryPolygonHole
}

type GeometryPolygonHole struct {
	Indices []int64
	Holes   [][]uint64
}

type GeometrySpline struct {
	HasExtras
	Closed          bool
	Sources         Sources
	ControlVertices GeometryControlVertices
}

type GeometryVertices struct {
	HasID
	HasName
	HasExtras
	Inputs []*Input
}

func NewGeometrySpline() (me *GeometrySpline) {
	me = &GeometrySpline{Sources: Sources{}}
	return
}

type GeometryDef struct {
	BaseDef
	Brep   *GeometryBrep
	Mesh   *GeometryMesh
	Spline *GeometrySpline
}

func (me *GeometryDef) init() {
}

type GeometryInst struct {
	BaseInst
	BindMaterial *BindMaterial
}

func (me *GeometryInst) init() {
}

//#begin-gt _definstlib.gt T:Geometry

func newGeometryDef(id string) (me *GeometryDef) {
	me = &GeometryDef{}
	me.ID = id
	me.Base.init()
	me.init()
	return
}

/*
//	Creates and returns a new *GeometryInst* instance referencing this *GeometryDef* definition.
func (me *GeometryDef) NewInst(id string) (inst *GeometryInst) {
	inst = &GeometryInst{Def: me}
	inst.init()
	return
}
*/

var (
	//	A *map* collection that contains *LibGeometryDefs* libraries associated by their *ID*.
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
//	*LibGeometryDefs* libraries associated by their *ID*.
type LibsGeometryDef map[string]*LibGeometryDefs

//	Creates a new *LibGeometryDefs* library with the specified *ID*, adds it to this *LibsGeometryDef*, and returns it.
//	
//	If this *LibsGeometryDef* already contains a *LibGeometryDefs* library with the specified *ID*, does nothing and returns *nil*.
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

//	A library that contains *GeometryDef*s associated by their *ID*. To create a new *LibGeometryDefs* library, ONLY
//	use the *LibsGeometryDef.New()* or *LibsGeometryDef.AddNew()* methods.
type LibGeometryDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*GeometryDef
}

func newLibGeometryDefs(id string) (me *LibGeometryDefs) {
	me = &LibGeometryDefs{M: map[string]*GeometryDef{}}
	me.ID = id
	return
}

//	Adds the specified *GeometryDef* definition to this *LibGeometryDefs*, and returns it.
//	
//	If this *LibGeometryDefs* already contains a *GeometryDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibGeometryDefs) Add(d *GeometryDef) (n *GeometryDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *GeometryDef* definition with the specified *ID*, adds it to this *LibGeometryDefs*, and returns it.
//	
//	If this *LibGeometryDefs* already contains a *GeometryDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibGeometryDefs) AddNew(id string) *GeometryDef { return me.Add(me.New(id)) }

//	Creates a new *GeometryDef* definition with the specified *ID* and returns it, but does not add it to this *LibGeometryDefs*.
func (me *LibGeometryDefs) New(id string) (def *GeometryDef) { def = newGeometryDef(id); return }

//	Removes the *GeometryDef* with the specified *ID* from this *LibGeometryDefs*.
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
