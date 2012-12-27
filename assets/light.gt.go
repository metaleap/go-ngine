package assets

//	Describes an ambient light source.
//	An ambient light is one that lights everything evenly, regardless of location or orientation.
type LightAmbient struct {
	//	Color
	LightBase
}

//	Describes how the intensity of a light source is attenuated.
type LightAttenuation struct {
	//	Constant light attenuation. Defaults to 1.
	Constant SidFloat

	//	Linear light attenuation.
	Linear SidFloat

	//	Quadratic light attenuation.
	Quadratic SidFloat
}

//	Constructor
func NewLightAttenuation() (me *LightAttenuation) {
	me = &LightAttenuation{}
	me.Constant.F = 1
	return
}

//	Contains three floating-point numbers specifying the color of a light.
type LightBase struct {
	//	Three floating-point numbers specifying the color of this light.
	Color Float3
}

//	Describes a directional light source.
//	A directional light is one that lights everything from the same direction, regardless of location.
//	The light's default direction vector in local coordinates is [0,0,-1], pointing down the negative z axis.
//	The actual direction of the light is defined by the transform of the node where the light is instantiated.
type LightDirectional struct {
	//	Color
	LightBase
}

//	Describes a point light source.
//	A point light source radiates light in all directions from a known location in space.
//	The position of the light is defined by the transform of the node in which it is instantiated.
type LightPoint struct {
	//	Color
	LightBase

	//	The intensity of a point light source is attenuated as the distance to the light source increases.
	Attenuation LightAttenuation
}

//	Constructor
func NewLightPoint() (me *LightPoint) {
	me = &LightPoint{}
	me.Attenuation.Constant.F = 1
	return
}

//	Describes a spot light source.
//	A spot light source radiates light in one direction in a cone shape from a known location in space.
//	The light's default direction vector in local coordinates is [0,0,-1], pointing down the negative z axis.
//	The actual direction of the light is defined by the transform of the node in which the light is instantiated.
type LightSpot struct {
	//	Color
	LightBase

	//	 The intensity of a spot light is also attenuated as the distance to the light source increases.
	Attenuation LightAttenuation

	//	The light's intensity is also attenuated as the radiation angle increases away from the direction of the light source.
	Falloff struct {
		//	Fall-off angle. Defaults to 180.
		Angle SidFloat

		//	Fall-off exponent.
		Exponent SidFloat
	}
}

//	Constructor
func NewLightSpot() (me *LightSpot) {
	me = &LightSpot{}
	me.Attenuation.Constant.F = 1
	me.Falloff.Angle.F = 180
	return
}

//	Declares a light source that illuminates a scene.
type LightDef struct {
	//	Id, Name, Asset, Extras
	BaseDef

	//	Techniques
	HasTechniques

	//	Common-technique profile. At least and at most one of its fields should ever be set.
	TC struct {
		//	If set, this light declares an ambient light.
		Ambient *LightAmbient

		//	If set, this light declares a directional light.
		Directional *LightDirectional

		//	If set, this light declares a point light.
		Point *LightPoint

		//	If set, this light declares a spot light.
		Spot *LightSpot
	}
}

//	Initialization
func (me *LightDef) Init() {
}

//	Instantiates a light resource.
type LightInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst

	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *LightDef
}

//	Initialization
func (me *LightInst) Init() {
}

//#begin-gt _definstlib.gt T:Light

func newLightDef(id string) (me *LightDef) {
	me = &LightDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new LightInst instance referencing this LightDef definition.
//	Any LightInst created by this method will have its Def field readily set to me.
func (me *LightDef) NewInst() (inst *LightInst) {
	inst = &LightInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct LightDef
//	according to the current me.DefRef value (by searching AllLightDefLibs).
//	Then returns me.Def.
//	(Note, every LightInst's Def is nil initially, unless it was created via LightDef.NewInst().)
func (me *LightInst) EnsureDef() *LightDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.LightDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibLightDefs libraries associated by their Id.
	AllLightDefLibs = LibsLightDef{}

	//	The "default" LibLightDefs library for LightDefs.
	LightDefs = AllLightDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllLightDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibLightDefs contained in AllLightDefLibs) for the LightDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) LightDef() (def *LightDef) {
	id := me.S()
	for _, lib := range AllLightDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllLightDefLibs variable:
//	a hash-table that contains LibLightDefs libraries associated by their Id.
type LibsLightDef map[string]*LibLightDefs

//	Creates a new LibLightDefs library with the specified Id, adds it to this LibsLightDef, and returns it.
//	If this LibsLightDef already contains a LibLightDefs library with the specified Id, does nothing and returns nil.
func (me LibsLightDef) AddNew(id string) (lib *LibLightDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsLightDef) new(id string) (lib *LibLightDefs) {
	lib = newLibLightDefs(id)
	return
}

//	A library that contains LightDefs associated by their Id.
//	To create a new LibLightDefs library, ONLY use the LibsLightDef.New() or LibsLightDef.AddNew() methods.
type LibLightDefs struct {
	//	Id, Name
	BaseLib

	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*LightDef
}

func newLibLightDefs(id string) (me *LibLightDefs) {
	me = &LibLightDefs{M: map[string]*LightDef{}}
	me.BaseLib.init(id)
	return
}

//	Adds the specified LightDef definition to this LibLightDefs, and returns it.
//	If this LibLightDefs already contains a LightDef definition with the same Id, does nothing and returns nil.
func (me *LibLightDefs) Add(d *LightDef) (n *LightDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new LightDef definition with the specified Id, adds it to this LibLightDefs, and returns it.
//	If this LibLightDefs already contains a LightDef definition with the specified Id, does nothing and returns nil.
func (me *LibLightDefs) AddNew(id string) *LightDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibLightDefs) Len() int { return len(me.M) }

//	Creates a new LightDef definition with the specified Id and returns it,
//	but does not add it to this LibLightDefs.
func (me *LibLightDefs) New(id string) (def *LightDef) { def = newLightDef(id); return }

//	Removes the LightDef with the specified Id from this LibLightDefs.
func (me *LibLightDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

func (me *LibLightDefs) resolver(part0 string) refSidResolver {
	return me.M[part0]
}

func (me *LibLightDefs) resolverRootIsLib() bool {
	return true
}

//	Signals to the core package (or your custom package) that changes have been made to this LibLightDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibLightDefs
//	library or its LightDef definitions. Also called by the global SyncChanges() function.
func (me *LibLightDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
