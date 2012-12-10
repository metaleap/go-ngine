package assets

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

const (
	FX_COLOR_TEXTURE_OPAQUE_A_ZERO   = 0
	FX_COLOR_TEXTURE_OPAQUE_A_ONE    = 1
	FX_COLOR_TEXTURE_OPAQUE_RGB_ZERO = 2
	FX_COLOR_TEXTURE_OPAQUE_RGB_ONE  = 3

	FX_PASS_PROGRAM_SHADER_STAGE_TESSELATION = 0
	FX_PASS_PROGRAM_SHADER_STAGE_VERTEX      = 1
	FX_PASS_PROGRAM_SHADER_STAGE_GEOMETRY    = 2
	FX_PASS_PROGRAM_SHADER_STAGE_FRAGMENT    = 3
	FX_PASS_PROGRAM_SHADER_STAGE_COMPUTE     = 4
)

type FxAnnotation struct {
	HasName
	Value interface{}
}

type FxColorOrTexture struct {
	Opaque   int
	Color    *ugfx.Rgba32
	ParamRef string
	Texture  *FxTexture
}

type FxParamDef struct {
	ParamDef
	Annotations        []*FxAnnotation
	Modifier, Semantic string
}

type FxParamDefs map[string]*FxParamDef

type FxPass struct {
	HasSid
	HasExtras
	Annotations []*FxAnnotation
	States      map[string]*FxPassState
	Program     *FxPassProgram
	Evaluate    *FxPassEvaluation
}

type FxPassEvaluation struct {
	Draw  string
	Color struct {
		Clear  *FxPassEvaluationClearColor
		Target *FxPassEvaluationTarget
	}
	Depth struct {
		Clear  *FxPassEvaluationClearDepth
		Target *FxPassEvaluationTarget
	}
	Stencil struct {
		Clear  *FxPassEvaluationClearStencil
		Target *FxPassEvaluationTarget
	}
}

type FxPassEvaluationClearColor struct {
	ugfx.Rgba32
	Index uint64
}

type FxPassEvaluationClearDepth struct {
	F     float64
	Index uint64
}

type FxPassEvaluationClearStencil struct {
	B     byte
	Index uint64
}

type FxPassEvaluationTarget struct {
	Index           uint64
	Slice           uint64
	Mip             uint64
	CubeFace        int
	SamplerParamRef string
	Image           *FxImageInst
}

type FxPassProgram struct {
	BindAttributes []*FxPassProgramBindAttribute
	BindUniforms   []*FxPassProgramBindUniform
	Shaders        []*FxPassProgramShader
}

type FxPassProgramBindAttribute struct {
	Symbol   string
	Semantic string
}

type FxPassProgramBindUniform struct {
	Symbol   string
	ParamRef string
	Value    interface{}
}

type FxPassProgramShader struct {
	HasExtras
	Stage   int
	Sources struct {
		Entry string
		All   []FxPassProgramShaderSources
	}
}

type FxPassProgramShaderSources struct {
	S           string
	IsImportRef bool
}

type FxPassState struct {
	Value    interface{}
	ParamRef string
	Index    int64
}

type FxProfile struct {
	HasID
	HasAsset
	HasExtras
	NewParams FxParamDefs
}

type FxProfileCommon struct {
	FxProfile
	Technique FxTechniqueCommon
}

func NewFxProfileCommon() (me *FxProfileCommon) {
	me = &FxProfileCommon{}
	me.NewParams = FxParamDefs{}
	return
}

type FxProfileGlSl struct {
	FxProfile
	Platform      string
	CodesIncludes []FxProfileGlSlCodeInclude
	Techniques    []*FxTechniqueGlsl
}

type FxProfileGlSlCodeInclude struct {
	ScopedString
	IsInclude bool
}

func NewFxProfileGlSl() (me *FxProfileGlSl) {
	me = &FxProfileGlSl{}
	me.NewParams = FxParamDefs{}
	return
}

type FxTechnique struct {
	HasID
	HasSid
	HasAsset
	HasExtras
}

type FxTechniqueCommon struct {
	FxTechnique
	Blinn    *FxTechniqueCommonBlinn
	Constant *FxTechniqueCommonConstant
	Lambert  *FxTechniqueCommonLambert
	Phong    *FxTechniqueCommonPhong
}

type FxTechniqueCommonBlinn struct {
	FxTechniqueCommonLambert
	Specular  *FxColorOrTexture
	Shininess *ParamScopedFloat
}

type FxTechniqueCommonConstant struct {
	Emission          *FxColorOrTexture
	Reflective        *FxColorOrTexture
	Reflectivity      *ParamScopedFloat
	Transparent       *FxColorOrTexture
	Transparency      *ParamScopedFloat
	IndexOfRefraction *ParamScopedFloat
}

type FxTechniqueCommonLambert struct {
	FxTechniqueCommonConstant
	Ambient *FxColorOrTexture
	Diffuse *FxColorOrTexture
}

type FxTechniqueCommonPhong struct {
	FxTechniqueCommonBlinn
}

type FxTechniqueGlsl struct {
	FxTechnique
	Annotations []*FxAnnotation
	Passes      []*FxPass
}

type FxTexture struct {
	Sampler2D string
	TexCoord  string
}

type FxEffectDef struct {
	BaseDef
	Annotations []*FxAnnotation
	NewParams   FxParamDefs
	Profiles    struct {
		GlSl   []*FxProfileGlSl
		Common []*FxProfileCommon
	}
}

func (me *FxEffectDef) Init() {
	me.NewParams = FxParamDefs{}
}

type FxEffectInst struct {
	BaseInst
	SetParams      []*ParamInst
	TechniqueHints []*FxEffectInstTechniqueHint
}

func (me *FxEffectInst) init() {
}

type FxEffectInstTechniqueHint struct {
	Platform string
	Ref      string
	Profile  string
}

//#begin-gt _definstlib.gt T:FxEffect

func newFxEffectDef(id string) (me *FxEffectDef) {
	me = &FxEffectDef{}
	me.ID = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *FxEffectInst* instance referencing this *FxEffectDef* definition.
func (me *FxEffectDef) NewInst(id string) (inst *FxEffectInst) {
	inst = &FxEffectInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibFxEffectDefs* libraries associated by their *ID*.
	AllFxEffectDefLibs = LibsFxEffectDef{}

	//	The "default" *LibFxEffectDefs* library for *FxEffectDef*s.
	FxEffectDefs = AllFxEffectDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllFxEffectDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllFxEffectDefLibs* variable: a *map* collection that contains
//	*LibFxEffectDefs* libraries associated by their *ID*.
type LibsFxEffectDef map[string]*LibFxEffectDefs

//	Creates a new *LibFxEffectDefs* library with the specified *ID*, adds it to this *LibsFxEffectDef*, and returns it.
//	
//	If this *LibsFxEffectDef* already contains a *LibFxEffectDefs* library with the specified *ID*, does nothing and returns *nil*.
func (me LibsFxEffectDef) AddNew(id string) (lib *LibFxEffectDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsFxEffectDef) new(id string) (lib *LibFxEffectDefs) {
	lib = newLibFxEffectDefs(id)
	return
}

//	A library that contains *FxEffectDef*s associated by their *ID*. To create a new *LibFxEffectDefs* library, ONLY
//	use the *LibsFxEffectDef.New()* or *LibsFxEffectDef.AddNew()* methods.
type LibFxEffectDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*FxEffectDef
}

func newLibFxEffectDefs(id string) (me *LibFxEffectDefs) {
	me = &LibFxEffectDefs{M: map[string]*FxEffectDef{}}
	me.ID = id
	return
}

//	Adds the specified *FxEffectDef* definition to this *LibFxEffectDefs*, and returns it.
//	
//	If this *LibFxEffectDefs* already contains a *FxEffectDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibFxEffectDefs) Add(d *FxEffectDef) (n *FxEffectDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *FxEffectDef* definition with the specified *ID*, adds it to this *LibFxEffectDefs*, and returns it.
//	
//	If this *LibFxEffectDefs* already contains a *FxEffectDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibFxEffectDefs) AddNew(id string) *FxEffectDef { return me.Add(me.New(id)) }

//	Creates a new *FxEffectDef* definition with the specified *ID* and returns it, but does not add it to this *LibFxEffectDefs*.
func (me *LibFxEffectDefs) New(id string) (def *FxEffectDef) { def = newFxEffectDef(id); return }

//	Removes the *FxEffectDef* with the specified *ID* from this *LibFxEffectDefs*.
func (me *LibFxEffectDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibFxEffectDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibFxEffectDefs* library or its *FxEffectDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibFxEffectDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
