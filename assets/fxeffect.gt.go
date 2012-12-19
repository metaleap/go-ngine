package assets

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

const (
	//	Takes the transparency information from the color's alpha channel, where the value 1.0 is opaque.
	FX_COLOR_TEXTURE_OPAQUE_A_ZERO = 0
	//	Takes the transparency information from the color's red, green, and blue channels, where the value 0.0 is opaque, with each channel modulated independently.
	FX_COLOR_TEXTURE_OPAQUE_A_ONE = 1
	//	Takes the transparency information from the color’s alpha channel, where the value 0.0 is opaque.
	FX_COLOR_TEXTURE_OPAQUE_RGB_ZERO = 2
	//	Takes the transparency information from the color’s red, green, and blue channels, where the value 1.0 is opaque, with each channel modulated independently.
	FX_COLOR_TEXTURE_OPAQUE_RGB_ONE = 3

	FX_PASS_PROGRAM_SHADER_STAGE_TESSELATION = 0
	FX_PASS_PROGRAM_SHADER_STAGE_VERTEX      = 1
	FX_PASS_PROGRAM_SHADER_STAGE_GEOMETRY    = 2
	FX_PASS_PROGRAM_SHADER_STAGE_FRAGMENT    = 3
	FX_PASS_PROGRAM_SHADER_STAGE_COMPUTE     = 4
)

//	Annotations communicate metadata from the Effect Runtime to the application only and are not otherwise interpreted within the *assets* package.
type FxAnnotation struct {
	//	Name
	HasName
	//	Value
	Value interface{}
}

//	Describes color attributes of fixed-function shaders inside FxProfileCommon effects.
type FxColorOrTexture struct {
	//	Specifies from which channel to take transparency information. Can be any of the FX_COLOR_TEXTURE_OPAQUE_* enumerated constants.
	Opaque int
	//	If set, describes he literal color of this value.
	Color *ugfx.Rgba32
	//	If set, refers to a previously-defined parameter in the current scope that provides 4 float values definining the literal color of this value.
	ParamRef string
	//	If set, refers to a previously-defined FxSampler with type FX_SAMPLER_TYPE_2D.
	Texture *FxTexture
}

//	Declares a new parameter for its parent FX-related resource, and assigns it an initial value.
type FxParamDef struct {
	//	Sid and Value
	ParamDef
	//	Application-specific FX metadata
	Annotations []*FxAnnotation
	//	Specifies constant, external, or uniform parameters.
	Modifier string
	//	Provides metadata that describes the purpose of a parameter declaration.
	Semantic string
}

//	A hash-table containing parameter declarations of this FX-related resource.
type FxParamDefs map[string]*FxParamDef

//	Provides a static declaration of all the render states, shaders, and settings for one rendering pipeline.
type FxPass struct {
	//	Sid
	HasSid
	//	Custom-profile/foreign-technique meta-data
	HasExtras
	//	Application-specific FX metadata
	Annotations []*FxAnnotation
	//	Contains all rendering states to set up for this pass.
	States map[string]*FxPassState
	//	Links multiple shaders together to produce a pipeline for geometry processing.
	Program *FxPassProgram
	//	Contains evaluation elements for this rendering pass.
	Evaluate *FxPassEvaluation
}

//	Constructor
func NewFxPass() (me *FxPass) {
	me = &FxPass{States: map[string]*FxPassState{}}
	return
}

//	Contains evaluation elements for a rendering pass.
type FxPassEvaluation struct {
	//	Instructs the FX Runtime what kind of geometry to submit.
	Draw string
	//	Color-information render target
	Color struct {
		//	Specifies whether this render target image is to be cleared, and which value to use.
		Clear *FxPassEvaluationClearColor
		//	Specifies which FxImageDef will receive the color information from the output of this pass.
		Target *FxPassEvaluationTarget
	}
	//	Depth-information render target
	Depth struct {
		//	Specifies whether this render target image is to be cleared, and which value to use.
		Clear *FxPassEvaluationClearDepth
		//	Specifies which FxImageDef will receive the depth information from the output of this pass.
		Target *FxPassEvaluationTarget
	}
	//	Stencil-information render target
	Stencil struct {
		//	Specifies whether this render target image is to be cleared, and which value to use.
		Clear *FxPassEvaluationClearStencil
		//	Specifies which FxImageDef will receive the stencil information from the output of this pass.
		Target *FxPassEvaluationTarget
	}
}

//	Specifies whether a color-information render target image is to be cleared, and which value to use.
type FxPassEvaluationClearColor struct {
	//	Default clear-color value
	ugfx.Rgba32
	//	Which of the multiple render targets is being set. The default is 0.
	Index uint64
}

//	Specifies whether a depth-information render target image is to be cleared, and which value to use.
type FxPassEvaluationClearDepth struct {
	//	Default clear-depth value
	F float64
	//	Which of the multiple render targets is being set. The default is 0.
	Index uint64
}

//	Specifies whether a stencil-information render target image is to be cleared, and which value to use.
type FxPassEvaluationClearStencil struct {
	//	Default clear-stencil value
	B byte
	//	Which of the multiple render targets is being set. The default is 0.
	Index uint64
}

//	Specifies which FxImageDef will receive the information from the output of a pass.
type FxPassEvaluationTarget struct {
	//	Indexes one of the Multiple Render Targets.
	Index uint64
	//	Indexes a sub-image inside a target surface, specifically: a layer of a 3D texture.
	Slice uint64
	//	Indexes a sub-image inside a target surface, specifically: a single MIP-map level.
	Mip uint64
	//	Indexes a sub-image inside a target surface, specifically: a unique cube face. Can be any of the FX_CUBE_FACE_* enumerated constants.
	CubeFace int
	//	If set, this render target references a sampler parameter to determine which image to use.
	SamplerParamRef string
	//	If set, this render target directly instantiates a renderable image.
	Image *FxImageInst
}

//	Constructor
func NewFxPassEvaluationTarget() (me *FxPassEvaluationTarget) {
	me = &FxPassEvaluationTarget{Index: 1}
	return
}

//	Links multiple shaders together to produce a pipeline for geometry processing.
type FxPassProgram struct {
	//	Information for binding the shader variables to effect parameters.
	BindAttributes []*FxPassProgramBindAttribute
	//	Binds a uniform shader variable to a parameter or a value.
	BindUniforms []*FxPassProgramBindUniform
	//	Setup and compilation information for shaders such as vertex and pixel shaders.
	Shaders []*FxPassProgramShader
}

//	Binds semantics to vertex attribute inputs of a shader.
type FxPassProgramBindAttribute struct {
	//	The identifier for a vertex attribute variable in the shader (a formal function parameter or in-scope global).
	Symbol string
	//	Contains an alternative name to the attribute variable for semantic binding to geometry vertex inputs.
	Semantic string
}

//	Binds values to uniform inputs of a shader or binds values to effect parameters upon instantiation.
type FxPassProgramBindUniform struct {
	//	The identifier for a uniform input parameter to the shader (a formal function parameter or in-scope global) that will be bound to an external resource.
	Symbol string
	//	If set, refers to a previously defined parameter providing the uniform value to be bound.
	ParamRef string
	//	If set, the uniform value to be bound.
	Value interface{}
}

//	Declares and prepares a shader for execution in the rendering pipeline of a pass.
type FxPassProgramShader struct {
	//	Custom-profile/foreign-technique meta-data
	HasExtras
	//	In which pipeline stage this programmable shader is designed to execute. Can be any of the FX_PASS_PROGRAM_SHADER_STAGE_* enumerated constants.
	Stage int
	//	Concatenates the source code for the shader from one or more sources.
	Sources []FxPassProgramShaderSources
}

//	Contains either code or an import reference.
type FxPassProgramShaderSources struct {
	//	The code or import reference.
	S string
	//	If true, S is an import reference; otherwise, S is code.
	IsImportRef bool
}

//	Represents a rendering state for a pass.
type FxPassState struct {
	//	If set, the value for this rendering state.
	Value string
	//	If set, refers to a previously defined parameter providing the value for this rendering state.
	ParamRef string
	//	State-specific optional index attribute.
	Index float64
}

//	An FX profile represents a shader-based rendering pipeline.
type FxProfile struct {
	//	Id
	HasId
	//	Resource-specific asset-management information and meta-data.
	HasAsset
	//	Custom-profile/foreign-technique meta-data
	HasExtras
	//	A hash-table containing parameter declarations of this profile.
	HasFxParamDefs
	//	If set, this FxProfile represents a common, fixed-function shader pipeline.
	Common *FxProfileCommon
	//	If set, this FxProfile represents an OpenGL Shading Language (GLSL) pipeline.
	GlSl *FxProfileGlSl
}

func NewProfile() (me *FxProfile) {
	me = &FxProfile{}
	me.NewParams = FxParamDefs{}
	return
}

//	This FX profile provides platform-independent declarations for the common, fixed-function shader.
type FxProfileCommon struct {
	//	Declares the only technique for this effect.
	Technique FxTechniqueCommon
}

//	This FX profile provides platform-specific declarations for the OpenGL Shading Language.
type FxProfileGlSl struct {
	//	The type of platform. This is a vendor-defined character string that indicates the platform or capability target for the technique. Defaults to "PC".
	Platform string
	//	GLSL shader sources
	CodesIncludes []FxProfileGlSlCodeInclude
	//	Declares the techniques for this effect.
	Techniques []*FxTechniqueGlsl
}

//	Constructor
func NewFxProfileGlSl() (me *FxProfileGlSl) {
	me = &FxProfileGlSl{Platform: "PC"}
	return
}

//	GLSL shader sources
type FxProfileGlSlCodeInclude struct {
	//	Source code or include reference
	ScopedString
	//	Indicates whether ScopedString is an import reference (true) or source code (false).
	IsInclude bool
}

//	Holds a description of the textures, samplers, shaders, parameters, and passes necessary for rendering this effect using one method.
type FxTechnique struct {
	//	Id
	HasId
	//	Sid
	HasSid
	//	Asset
	HasAsset
	//	Extras
	HasExtras
}

//	Holds a description of the textures, samplers, shaders, parameters, and passes necessary for rendering this effect within an FxProfileCommon.
type FxTechniqueCommon struct {
	//	Id, Sid, Asset, Extras
	FxTechnique
	//	Produces a shaded surface with a Blinn BRDF approximation.
	Blinn *FxTechniqueCommonBlinn
	//	Produces a constantly shaded surface that is independent of lighting.
	Constant *FxTechniqueCommonConstant
	//	Produces a constantly shaded surface that is independent of lighting.
	Lambert *FxTechniqueCommonLambert
	//	Produces a shaded surface where the specular reflection is shaded according the Phong BRDF approximation.
	Phong *FxTechniqueCommonPhong
}

//	Produces a shaded surface with a Blinn BRDF approximation.
type FxTechniqueCommonBlinn struct {
	//	Ambient, Diffuse, Emission, Reflective, Reflectivity, Transparent, Transparency, IndexOfRefraction
	FxTechniqueCommonLambert
	//	Declares the color of light specularly reflected from the surface of this object.
	Specular *FxColorOrTexture
	//	Declares the specularity or roughness of the specular reflection lobe.
	Shininess *ParamScopedFloat
}

//	Produces a constantly shaded surface that is independent of lighting.
type FxTechniqueCommonConstant struct {
	//	Declares the amount of light emitted from the surface of this object
	Emission *FxColorOrTexture
	//	Declares the color of a perfect mirror reflection.
	Reflective *FxColorOrTexture
	//	Declares the amount of perfect mirror reflection to be added to the reflected light as a value between 0.0 and 1.0.
	Reflectivity *ParamScopedFloat
	//	Declares the color of perfectly refracted light.
	Transparent *FxColorOrTexture
	//	Declares the amount of perfectly refracted light added to the reflected color as a scalar value between 0.0 and 1.0.
	Transparency *ParamScopedFloat
	//	Declares the index of refraction for perfectly refracted light as a single scalar index.
	IndexOfRefraction *ParamScopedFloat
}

//	Produces a constantly shaded surface that is independent of lighting.
type FxTechniqueCommonLambert struct {
	//	Emission, Reflective, Reflectivity, Transparent, Transparency, IndexOfRefraction
	FxTechniqueCommonConstant
	//	Declares the amount of ambient light reflected from the surface of this object.
	Ambient *FxColorOrTexture
	//	Declares the amount of light diffusely reflected from the surface of this object.
	Diffuse *FxColorOrTexture
}

//	Produces a shaded surface where the specular reflection is shaded according the Phong BRDF approximation.
type FxTechniqueCommonPhong struct {
	//	Specular, Shininess, Ambient, Diffuse, Emission, Reflective, Reflectivity, Transparent, Transparency, IndexOfRefraction
	FxTechniqueCommonBlinn
}

//	Holds a description of the textures, samplers, shaders, parameters, and passes necessary for rendering this effect within an FxProfileGlsl.
type FxTechniqueGlsl struct {
	//	Id, Sid, Asset, Extras
	FxTechnique
	//	Application-specific FX metadata
	Annotations []*FxAnnotation
	//	Static declarations of all the render states, shaders, and settings for the rendering pipeline.
	Passes []*FxPass
}

//	Used in FxColorOrTexture instances that refer to a texture image instead of a literal color value.
type FxTexture struct {
	//	Extras
	HasExtras
	//	References a previously defined FxSampler of type FX_SAMPLER_TYPE_2D.
	Sampler2D string
	//	A semantic token, which will be referenced within FxMaterialBinding to bind an array of texture-coordinates from a geometry instance to the sampler.
	TexCoord string
}

//	Defines the equations necessary for the visual appearance of geometry and/or screen-space image processing.
type FxEffectDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	NewParams
	HasFxParamDefs
	//	Application-specific FX metadata
	Annotations []*FxAnnotation
	//	Rendering pipeline(s).
	Profiles []*FxProfile
}

//	Initialization
func (me *FxEffectDef) Init() {
	me.NewParams = FxParamDefs{}
}

//	Instantiates an effect resource.
type FxEffectInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	SetParams
	HasParamInsts
	//	Platform-specific hints of which techniques to use in this effect.
	TechniqueHints []*FxEffectInstTechniqueHint
}

//	Initialization
func (me *FxEffectInst) Init() {
}

//	Adds a hint for a platform of which technique to use in this effect.
type FxEffectInstTechniqueHint struct {
	//	Defines a string that specifies for which platform this hint is intended. Optional.
	Platform string
	//	A reference to the name of the platform. Required.
	Ref string
	//	A string that specifies for which API profile this hint is intended. It is the name of the profile within the effect that contains the technique. Optional. If set, can be "COMMON" or "GLSL".
	Profile string
}

//#begin-gt _definstlib.gt T:FxEffect

func newFxEffectDef(id string) (me *FxEffectDef) {
	me = &FxEffectDef{}
	me.Id = id
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
	//	A *map* collection that contains *LibFxEffectDefs* libraries associated by their *Id*.
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
//	*LibFxEffectDefs* libraries associated by their *Id*.
type LibsFxEffectDef map[string]*LibFxEffectDefs

//	Creates a new *LibFxEffectDefs* library with the specified *Id*, adds it to this *LibsFxEffectDef*, and returns it.
//	
//	If this *LibsFxEffectDef* already contains a *LibFxEffectDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *FxEffectDef*s associated by their *Id*. To create a new *LibFxEffectDefs* library, ONLY
//	use the *LibsFxEffectDef.New()* or *LibsFxEffectDef.AddNew()* methods.
type LibFxEffectDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*FxEffectDef
}

func newLibFxEffectDefs(id string) (me *LibFxEffectDefs) {
	me = &LibFxEffectDefs{M: map[string]*FxEffectDef{}}
	me.Id = id
	return
}

//	Adds the specified *FxEffectDef* definition to this *LibFxEffectDefs*, and returns it.
//	
//	If this *LibFxEffectDefs* already contains a *FxEffectDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibFxEffectDefs) Add(d *FxEffectDef) (n *FxEffectDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *FxEffectDef* definition with the specified *Id*, adds it to this *LibFxEffectDefs*, and returns it.
//	
//	If this *LibFxEffectDefs* already contains a *FxEffectDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibFxEffectDefs) AddNew(id string) *FxEffectDef { return me.Add(me.New(id)) }

//	Short-hand for len(lib.M)
func (me *LibFxEffectDefs) Len() int { return len(me.M) }

//	Creates a new *FxEffectDef* definition with the specified *Id* and returns it, but does not add it to this *LibFxEffectDefs*.
func (me *LibFxEffectDefs) New(id string) (def *FxEffectDef) { def = newFxEffectDef(id); return }

//	Removes the *FxEffectDef* with the specified *Id* from this *LibFxEffectDefs*.
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
