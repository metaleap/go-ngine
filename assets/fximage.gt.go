package assets

const (
	FX_CREATE_CUBE_FACE_NEGATIVE_X          = 0x8516
	FX_CREATE_CUBE_FACE_NEGATIVE_Y          = 0x8518
	FX_CREATE_CUBE_FACE_NEGATIVE_Z          = 0x851A
	FX_CREATE_CUBE_FACE_POSITIVE_X          = 0x8515
	FX_CREATE_CUBE_FACE_POSITIVE_Y          = 0x8517
	FX_CREATE_CUBE_FACE_POSITIVE_Z          = 0x8519
	FX_CREATE_FORMAT_HINT_CHANNELS_RGB      = 0
	FX_CREATE_FORMAT_HINT_CHANNELS_RGBA     = iota
	FX_CREATE_FORMAT_HINT_CHANNELS_RGBE     = iota
	FX_CREATE_FORMAT_HINT_CHANNELS_LUM      = iota
	FX_CREATE_FORMAT_HINT_CHANNELS_LUMA     = iota
	FX_CREATE_FORMAT_HINT_CHANNELS_DEPTH    = iota
	FX_CREATE_FORMAT_HINT_RANGE_SNORM       = 0
	FX_CREATE_FORMAT_HINT_RANGE_UNORM       = iota
	FX_CREATE_FORMAT_HINT_RANGE_SINT        = iota
	FX_CREATE_FORMAT_HINT_RANGE_UINT        = iota
	FX_CREATE_FORMAT_HINT_RANGE_FLOAT       = iota
	FX_CREATE_FORMAT_HINT_PRECISION_DEFAULT = 0
	FX_CREATE_FORMAT_HINT_PRECISION_LOW     = iota
	FX_CREATE_FORMAT_HINT_PRECISION_MID     = iota
	FX_CREATE_FORMAT_HINT_PRECISION_HIGH    = iota
	FX_CREATE_FORMAT_HINT_PRECISION_MAX     = iota
)

type FxCreate2D struct {
	FxCreateCommon
	Size struct {
		Exact *FxCreate2DSizeExact
		Ratio *FxCreate2DSizeRatio
	}
	Mips         *FxCreateMips
	Unnormalized bool
	InitFrom     []*FxCreateInitFrom
}

type FxCreate2DSizeExact struct {
	Width  uint64
	Height uint64
}

type FxCreate2DSizeRatio struct {
	Width  float64
	Height float64
}

type FxCreate3D struct {
	FxCreateCommon
	Size struct {
		Width  uint64
		Height uint64
		Depth  uint64
	}
	Mips     FxCreateMips
	InitFrom []*FxCreate3DInitFrom
}

type FxCreate3DInitFrom struct {
	FxCreateInitFrom
	Depth uint64
}

type FxCreateCommon struct {
	ArrayLength uint64
	Format      *FxCreateFormat
}

type FxCreateCube struct {
	FxCreateCommon
	Size struct {
		Width uint64
	}
	Mips     FxCreateMips
	InitFrom []*FxCreateCubeInitFrom
}

type FxCreateCubeInitFrom struct {
	FxCreateInitFrom
	Face int
}

type FxCreateFormat struct {
	Exact string
	Hint  *FxCreateFormatHint
}

type FxCreateFormatHint struct {
	Channels  int
	Range     int
	Precision int
	Space     string
}

type FxCreateInitFrom struct {
	FxInitFrom
	ArrayIndex uint64
	MipIndex   uint64
}

type FxCreateMips struct {
	Levels    uint64
	NoAutoGen bool
}

type FxImageInitFrom struct {
	FxInitFrom
	NoAutoMip bool
}

type FxInitFrom struct {
	Raw struct {
		Data   []byte
		Format string
	}
	RefUrl string
}

type FxImageDef struct {
	BaseDef
	HasSid
	Renderable struct {
		Is     bool
		Shared bool
	}
	Init struct {
		Create2D   *FxCreate2D
		Create3D   *FxCreate3D
		CreateCube *FxCreateCube
		From       *FxImageInitFrom
	}
}

func (me *FxImageDef) init() {
}

type FxImageInst struct {
	BaseInst
	Def *FxImageDef
}

func (me *FxImageInst) init() {
}

//#begin-gt _definstlib.gt T:FxImage

func newFxImageDef(id string) (me *FxImageDef) {
	me = &FxImageDef{}
	me.ID = id
	me.init()
	return
}

/*
//	Creates and returns a new *FxImageInst* instance referencing this *FxImageDef* definition.
func (me *FxImageDef) NewInst(id string) (inst *FxImageInst) {
	inst = &FxImageInst{Def: me}
	inst.init()
	return
}
*/

var (
	//	A *map* collection that contains *LibFxImageDefs* libraries associated by their *ID*.
	AllFxImageDefLibs = LibsFxImageDef{}

	//	The "default" *LibFxImageDefs* library for *FxImageDef*s.
	FxImageDefs = AllFxImageDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllFxImageDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllFxImageDefLibs* variable: a *map* collection that contains
//	*LibFxImageDefs* libraries associated by their *ID*.
type LibsFxImageDef map[string]*LibFxImageDefs

//	Creates a new *LibFxImageDefs* library with the specified *ID*, adds it to this *LibsFxImageDef*, and returns it.
//	
//	If this *LibsFxImageDef* already contains a *LibFxImageDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsFxImageDef) AddNew(id string) (lib *LibFxImageDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsFxImageDef) new(id string) (lib *LibFxImageDefs) {
	lib = newLibFxImageDefs(id)
	return
}

//	A library that contains *FxImageDef*s associated by their *ID*. To create a new *LibFxImageDefs* library, ONLY
//	use the *LibsFxImageDef.New()* or *LibsFxImageDef.AddNew()* methods.
type LibFxImageDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*FxImageDef
}

func newLibFxImageDefs(id string) (me *LibFxImageDefs) {
	me = &LibFxImageDefs{M: map[string]*FxImageDef{}}
	me.ID = id
	return
}

//	Adds the specified *FxImageDef* definition to this *LibFxImageDefs*, and returns it.
//	
//	If this *LibFxImageDefs* already contains a *FxImageDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibFxImageDefs) Add(d *FxImageDef) (n *FxImageDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *FxImageDef* definition with the specified *ID*, adds it to this *LibFxImageDefs*, and returns it.
//	
//	If this *LibFxImageDefs* already contains a *FxImageDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibFxImageDefs) AddNew(id string) *FxImageDef { return me.Add(me.New(id)) }

//	Creates a new *FxImageDef* definition with the specified *ID* and returns it, but does not add it to this *LibFxImageDefs*.
func (me *LibFxImageDefs) New(id string) (def *FxImageDef) { def = newFxImageDef(id); return }

//	Removes the *FxImageDef* with the specified *ID* from this *LibFxImageDefs*.
func (me *LibFxImageDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibFxImageDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibFxImageDefs* library or its *FxImageDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibFxImageDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
