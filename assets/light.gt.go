package assets

//	Describes an ambient light source.
type LightAmbient struct {
	//	Declares the parameters required to describe an ambient light source. An
	//	ambient light is one that lights everything evenly, regardless of location or orientation.
	LightBase
}

type LightAttenuation struct {
	Constant  ScopedFloat
	Linear    ScopedFloat
	Quadratic ScopedFloat
}

type LightBase struct {
	Color Float3
}

type LightDirectional struct {
	LightBase
}

type LightPoint struct {
	LightBase
	Attenuation LightAttenuation
}

type LightSpot struct {
	LightBase
	Attenuation LightAttenuation
	Falloff     struct {
		Angle    ScopedFloat
		Exponent ScopedFloat
	}
}

func NewLightSpot() (me *LightSpot) {
	me = &LightSpot{}
	me.Falloff.Angle.F = 180
	return
}

type LightDef struct {
	BaseDef
	HasTechniques
	TC struct {
		Ambient     *LightAmbient
		Directional *LightDirectional
		Point       *LightPoint
		Spot        *LightSpot
	}
}

func (me *LightDef) Init() {
}

type LightInst struct {
	BaseInst
}

func (me *LightInst) Init() {
}

//#begin-gt _definstlib.gt T:Light

func newLightDef(id string) (me *LightDef) {
	me = &LightDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *LightInst* instance referencing this *LightDef* definition.
func (me *LightDef) NewInst(id string) (inst *LightInst) {
	inst = &LightInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibLightDefs* libraries associated by their *Id*.
	AllLightDefLibs = LibsLightDef{}

	//	The "default" *LibLightDefs* library for *LightDef*s.
	LightDefs = AllLightDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllLightDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllLightDefLibs* variable: a *map* collection that contains
//	*LibLightDefs* libraries associated by their *Id*.
type LibsLightDef map[string]*LibLightDefs

//	Creates a new *LibLightDefs* library with the specified *Id*, adds it to this *LibsLightDef*, and returns it.
//	
//	If this *LibsLightDef* already contains a *LibLightDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *LightDef*s associated by their *Id*. To create a new *LibLightDefs* library, ONLY
//	use the *LibsLightDef.New()* or *LibsLightDef.AddNew()* methods.
type LibLightDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*LightDef
}

func newLibLightDefs(id string) (me *LibLightDefs) {
	me = &LibLightDefs{M: map[string]*LightDef{}}
	me.Id = id
	return
}

//	Adds the specified *LightDef* definition to this *LibLightDefs*, and returns it.
//	
//	If this *LibLightDefs* already contains a *LightDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibLightDefs) Add(d *LightDef) (n *LightDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *LightDef* definition with the specified *Id*, adds it to this *LibLightDefs*, and returns it.
//	
//	If this *LibLightDefs* already contains a *LightDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibLightDefs) AddNew(id string) *LightDef { return me.Add(me.New(id)) }

//	Creates a new *LightDef* definition with the specified *Id* and returns it, but does not add it to this *LibLightDefs*.
func (me *LibLightDefs) New(id string) (def *LightDef) { def = newLightDef(id); return }

//	Removes the *LightDef* with the specified *Id* from this *LibLightDefs*.
func (me *LibLightDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibLightDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibLightDefs* library or its *LightDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibLightDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
