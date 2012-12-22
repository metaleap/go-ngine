package assets

import (
	ugfx "github.com/metaleap/go-util/gfx"
)

type FxTextureOpaque int

const (
	//	Takes the transparency information from the color's alpha channel,
	//	where the value 1.0 is opaque. This is the default.
	FxTextureOpaqueA1 FxTextureOpaque = iota
	//	Takes the transparency information from the color's alpha channel,
	//	where the value 0.0 is opaque.
	FxTextureOpaqueA0
	//	Takes the transparency information from the color's red, green, and blue channels,
	//	where the value 0.0 is opaque, with each channel modulated independently.
	FxTextureOpaqueRgb0
	//	Takes the transparency information from the color's red, green, and blue channels,
	//	where the value 1.0 is opaque, with each channel modulated independently.
	FxTextureOpaqueRgb1
)

type FxShaderStage int

const (
	_ = iota
	//	This programmable shader is designed to execute in the Tessellation pipeline stage.
	FxShaderStageTessellation FxShaderStage = iota
	//	This programmable shader is designed to execute in the Vertex pipeline stage.
	FxShaderStageVertex
	//	This programmable shader is designed to execute in the Geometry pipeline stage.
	FxShaderStageGeometry
	//	This programmable shader is designed to execute in the Fragment pipeline stage.
	FxShaderStageFragment
	//	This programmable shader is designed to execute in the Compute pipeline stage.
	FxShaderStageCompute
)

//	Annotations communicate metadata from the Effect Runtime to the application only
//	and are not otherwise interpreted within the *assets* package.
type FxAnnotation struct {
	//	Name
	HasName
	//	Annotation value.
	Value interface{}
}

//	Used to describe the literal color of an FxColorOrTexture.
type FxColor struct {
	//	Sid
	HasSid
	//	Describes the literal color of the parent FxColorOrTexture.
	Color ugfx.Rgba32
}

//	Describes color attributes of fixed-function shaders inside FxProfileCommon effects.
type FxColorOrTexture struct {
	//	Specifies from which channel to take transparency information.
	//	Must be one of the FxTextureOpaque* enumerated constants.
	Opaque FxTextureOpaque
	//	If set, refers to a previously-defined parameter in the current scope that provides
	//	four float values describing the literal color of this value.
	ParamRef RefParam
	//	If set, describes he literal color of this value.
	Color *FxColor
	//	If set, refers to a previously-defined FxSampler with a Kind of FxSamplerKind2D.
	Texture *FxTexture
}

//	Adds a hint for a platform of which technique to use in this effect.
type FxEffectInstTechniqueHint struct {
	//	Defines a string that specifies for which platform this hint is intended. Optional.
	Platform string
	//	A reference to the name of the platform. Required.
	Ref string
	//	Specifies for which API profile this hint is intended.
	//	Optional. If set, must be either "COMMON" or "GLSL".
	Profile string
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

//	Provides a static declaration of all the render states, shaders, and settings
//	for one rendering pipeline.
type FxPass struct {
	//	Sid
	HasSid
	//	Extras
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
	//	Indexes a sub-image inside a target surface, specifically: a unique cube face.
	//	Must be one of the FxCubeFace* enumerated constants.
	CubeFace FxCubeFace
	//	If set, Image is ignored; this render target references a sampler parameter to determine which image to use.
	Sampler RefParam
	//	If set (and Sampler is empty), this render target directly instantiates a renderable image.
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
	//	The identifier for a vertex attribute variable in the shader
	//	(a formal function parameter or in-scope global).
	Symbol string
	//	Contains an alternative name to the attribute variable
	//	for semantic binding to geometry vertex inputs.
	Semantic string
}

//	Binds values to uniform inputs of a shader or binds values to effect parameters upon instantiation.
type FxPassProgramBindUniform struct {
	//	The identifier for a uniform input parameter to the shader
	//	(a formal function parameter or in-scope global) that will be bound to an external resource.
	Symbol string
	//	If set, refers to a previously defined parameter providing the uniform value to be bound.
	ParamRef RefParam
	//	If set, the uniform value to be bound.
	Value interface{}
}

//	Declares and prepares a shader for execution in the rendering pipeline of a pass.
type FxPassProgramShader struct {
	//	Extras
	HasExtras
	//	In which pipeline stage this programmable shader is designed to execute.
	//	Must be one of the FxShaderStage* enumerated constants.
	Stage FxShaderStage
	//	Concatenates the source code for the shader from one or more sources.
	Sources []FxPassProgramShaderSources
}

//	Contains either code or an import reference.
type FxPassProgramShaderSources struct {
	//	If true, S is an import reference; otherwise, S is code.
	IsImportRef bool
	//	The code or import reference.
	S string
}

//	Represents a rendering state for a pass.
type FxPassState struct {
	//	If set, Value is ignored; refers to a previously defined parameter providing the value for this rendering state.
	Param RefParam
	//	If set (and Param is empty), the value for this rendering state.
	Value string
	//	State-specific optional index attribute.
	Index float64
}

//	An FX profile represents a shader-based rendering pipeline.
type FxProfile struct {
	//	Id
	HasId
	//	Asset
	HasAsset
	//	Extras
	HasExtras
	//	NewParams
	HasFxParamDefs
	//	If set, Glsl must be nil, and this FxProfile represents a common, fixed-function shader pipeline.
	Common *FxProfileCommon
	//	If set, Common must be nil, and this FxProfile represents an OpenGL Shading Language pipeline.
	Glsl *FxProfileGlsl
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
type FxProfileGlsl struct {
	//	The type of platform. This is a vendor-defined character string that
	//	indicates the platform or capability target for the technique. Defaults to "PC".
	Platform string
	//	GLSL shader sources
	CodesIncludes []FxProfileGlslCodeInclude
	//	Declares the techniques for this effect.
	Techniques FxGlslTechniques
}

//	Constructor
func NewFxProfileGlsl() (me *FxProfileGlsl) {
	me = &FxProfileGlsl{Platform: "PC", Techniques: FxGlslTechniques{}}
	return
}

//	GLSL shader sources
type FxProfileGlslCodeInclude struct {
	//	Source code or include reference
	SidString
	//	Indicates whether SidString is an import reference (true) or source code (false).
	IsInclude bool
}

//	A hash-table of GLSL techniques mapped to their scoped identifiers.
type FxGlslTechniques map[string]*FxTechniqueGlsl

//	Holds a description of the textures, samplers, shaders, parameters, and passes
//	necessary for rendering this effect using one method.
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

//	Holds a description of the textures, samplers, shaders, parameters, and passes
//	necessary for rendering this effect within an FxProfileCommon.
type FxTechniqueCommon struct {
	//	Id, Sid, Asset, Extras
	FxTechnique
	//	Produces a shaded surface with a Blinn BRDF approximation.
	Blinn *FxTechniqueCommonBlinn
	//	Produces a constantly shaded surface that is independent of lighting.
	Constant *FxTechniqueCommonConstant
	//	Produces a constantly shaded surface that is independent of lighting.
	Lambert *FxTechniqueCommonLambert
	//	Produces a shaded surface with a Phong BRDF approximation.
	Phong *FxTechniqueCommonPhong
}

//	Produces a shaded surface with a Blinn BRDF approximation.
type FxTechniqueCommonBlinn struct {
	//	Ambient, Diffuse, Emission, Reflective, Reflectivity, Transparent, Transparency, IndexOfRefraction
	FxTechniqueCommonLambert
	//	Declares the color of light specularly reflected from the surface of this object.
	Specular *FxColorOrTexture
	//	Declares the specularity or roughness of the specular reflection lobe.
	Shininess *ParamOrSidFloat
}

//	Produces a constantly shaded surface that is independent of lighting.
type FxTechniqueCommonConstant struct {
	//	Declares the amount of light emitted from the surface of this object
	Emission *FxColorOrTexture
	//	Declares the color of a perfect mirror reflection.
	Reflective *FxColorOrTexture
	//	Declares the amount of perfect mirror reflection to be added to the reflected light
	//	as a value between 0.0 and 1.0.
	Reflectivity *ParamOrSidFloat
	//	Declares the color of perfectly refracted light.
	Transparent *FxColorOrTexture
	//	Declares the amount of perfectly refracted light added to the reflected color
	//	as a scalar value between 0.0 and 1.0.
	Transparency *ParamOrSidFloat
	//	Declares the index of refraction for perfectly refracted light as a single scalar index.
	IndexOfRefraction *ParamOrSidFloat
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

//	Produces a shaded surface with a Phong BRDF approximation.
type FxTechniqueCommonPhong struct {
	//	Specular, Shininess, Ambient, Diffuse, Emission, Reflective, Reflectivity, Transparent, Transparency, IndexOfRefraction
	FxTechniqueCommonBlinn
}

//	Holds a description of the textures, samplers, shaders, parameters, and passes
//	necessary for rendering this effect within an FxProfileGlsl.
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
	//	References a previously defined FxSampler of Kind FxSamplerKind2D.
	Sampler2D RefParam
	//	A semantic token, which will be referenced within FxMaterialBinding
	//	to bind an array of texture-coordinates from a geometry instance to the sampler.
	TexCoord string
}

//	Defines the equations necessary for the visual appearance of geometry or screen-space image processing.
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
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *FxEffectDef
	//	Platform-specific hints of which techniques to use in this effect.
	TechniqueHints []*FxEffectInstTechniqueHint
}

//	Initialization
func (me *FxEffectInst) Init() {
	me.SetParams = ParamInsts{}
}

//#begin-gt _definstlib.gt T:FxEffect

func newFxEffectDef(id string) (me *FxEffectDef) {
	me = &FxEffectDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new FxEffectInst instance referencing this FxEffectDef definition.
//	Any FxEffectInst created by this method will have its Def field readily set to me.
func (me *FxEffectDef) NewInst() (inst *FxEffectInst) {
	inst = &FxEffectInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct FxEffectDef
//	according to the current me.DefRef value (by searching AllFxEffectDefLibs).
//	Then returns me.Def.
//	(Note, every FxEffectInst's Def is nil initially, unless it was created via FxEffectDef.NewInst().)
func (me *FxEffectInst) EnsureDef() *FxEffectDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.FxEffectDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibFxEffectDefs libraries associated by their Id.
	AllFxEffectDefLibs = LibsFxEffectDef{}

	//	The "default" LibFxEffectDefs library for FxEffectDefs.
	FxEffectDefs = AllFxEffectDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllFxEffectDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibFxEffectDefs contained in AllFxEffectDefLibs) for the FxEffectDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) FxEffectDef() (def *FxEffectDef) {
	id := me.S()
	for _, lib := range AllFxEffectDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllFxEffectDefLibs variable:
//	a hash-table that contains LibFxEffectDefs libraries associated by their Id.
type LibsFxEffectDef map[string]*LibFxEffectDefs

//	Creates a new LibFxEffectDefs library with the specified Id, adds it to this LibsFxEffectDef, and returns it.
//	If this LibsFxEffectDef already contains a LibFxEffectDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains FxEffectDefs associated by their Id.
//	To create a new LibFxEffectDefs library, ONLY use the LibsFxEffectDef.New() or LibsFxEffectDef.AddNew() methods.
type LibFxEffectDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*FxEffectDef
}

func newLibFxEffectDefs(id string) (me *LibFxEffectDefs) {
	me = &LibFxEffectDefs{M: map[string]*FxEffectDef{}}
	me.Id = id
	return
}

//	Adds the specified FxEffectDef definition to this LibFxEffectDefs, and returns it.
//	If this LibFxEffectDefs already contains a FxEffectDef definition with the same Id, does nothing and returns nil.
func (me *LibFxEffectDefs) Add(d *FxEffectDef) (n *FxEffectDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new FxEffectDef definition with the specified Id, adds it to this LibFxEffectDefs, and returns it.
//	If this LibFxEffectDefs already contains a FxEffectDef definition with the specified Id, does nothing and returns nil.
func (me *LibFxEffectDefs) AddNew(id string) *FxEffectDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibFxEffectDefs) Len() int { return len(me.M) }

//	Creates a new FxEffectDef definition with the specified Id and returns it,
//	but does not add it to this LibFxEffectDefs.
func (me *LibFxEffectDefs) New(id string) (def *FxEffectDef) { def = newFxEffectDef(id); return }

//	Removes the FxEffectDef with the specified Id from this LibFxEffectDefs.
func (me *LibFxEffectDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibFxEffectDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibFxEffectDefs
//	library or its FxEffectDef definitions. Also called by the global SyncChanges() function.
func (me *LibFxEffectDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
