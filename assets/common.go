package assets

//	Categorizes the kind of a Transform.
type TransformKind int

const (
	//	A position and orientation transformation suitable for aiming a camera.
	TransformKindLookat TransformKind = iota + 1

	//	A transformation that embodies mathematical changes to points within a coordinate system
	//	or the coordinate system itself.
	TransformKindMatrix

	//	A transformation that specifies how to rotate an object around an axis.
	TransformKindRotate

	//	A transformation that specifies how to deform an object along one axis.
	TransformKindSkew

	//	A transformation that specifies how to change an object's size.
	TransformKindScale

	//	A transformation that changes the position of an object in a local coordinate system.
	TransformKindTranslate
)

//	Used in all resources that require asset-management information.
type HasAsset struct {
	//	Resource-specific asset-management information and meta-data.
	Asset *Asset
}

//	Used in all resources that support custom techniques / foreign profiles.
type HasExtras struct {
	//	Custom-technique/foreign-profile meta-data.
	Extras []*Extra
}

//	Used in all FX resources that declare their own parameters.
type HasFxParamDefs struct {
	//	A hash-table containing all parameter declarations of this resource.
	NewParams FxParamDefs
}

//	Used in all resources that declare their own unique identifier.
type HasId struct {
	//	The unique identifier of this resource.
	Id string
}

//	Used in all data consumers that require input connections into a data Source.
type HasInputs struct {
	//	Declares the input semantics of a data Source and connects a consumer to that Source.
	Inputs []*Input
}

//	Used in all resources that support arbitrary pretty-print names/titles.
type HasName struct {
	//	The optional pretty-print name/title of this resource.
	Name string
}

//	Used in all resources that declare their own parameters.
type HasParamDefs struct {
	//	A hash-table containing all parameter declarations of this resource.
	NewParams ParamDefs
}

//	Used in all resources that assign values to other parameters.
type HasParamInsts struct {
	//	A hash-table containing all parameter values assigned by this resource.
	SetParams ParamInsts
}

//	Used in all resources that declare their own scoped identifier.
type HasSid struct {
	//	The Scoped identifier of this resource.
	Sid string
}

//	Used in all resources that provide data arrays.
type HasSources struct {
	//	Provides the bulk of this resource's data.
	Sources Sources
}

//	Used in all resources that support custom techniques / foreign profiles.
type HasTechniques struct {
	//	Custom-technique/foreign-profile content or data.
	Techniques []*Technique
}

//	Resource-specific asset-management information and meta-data.
type Asset struct {
	//	Extras
	HasExtras

	/*
		//	Contains the date and time that the parent element was created.
		//	This is only set-and-retained for imported Collada assets that provide this field,
		//	and is not otherwise used.
		Created string
		//	Contains the date and time that the parent element was last modified.
		//	This is only set-and-retained for imported Collada assets that provide this field,
		//	and is not otherwise used.
		Modified string
		//	Contains a list of words used as search criteria.
		//	This is only set-and-retained for imported Collada assets that provide this field,
		//	and is not otherwise used.
		Keywords string
		//	Contains revision information.
		//	This is only set-and-retained for imported Collada assets that provide this field,
		//	and is not otherwise used.
		Revision string
		//	Contains a description of the topical subject.
		//	This is only set-and-retained for imported Collada assets that provide this field,
		//	and is not otherwise used.
		Subject string
		//	Contains title information.
		//	This is only set-and-retained for imported Collada assets that provide this field,
		//	and is not otherwise used.
		Title string
	*/

	//	Contains descriptive information about the coordinate system of the geometric data. All
	//	coordinates are right-handed by definition. Valid values are "X", "Y" (the default), or "Z".
	UpAxis string

	//	The unit of distance that applies to all spatial measurements within the scope of this resource.
	Unit struct {
		//	How many real-world meters in one distance unit as a floating-point number.
		//	1.0 for meter, 0.01 for centimeter, 1000 for kilometer etc.
		Meter float64

		//	Name of the distance unit, such as "centimeter", "kilometer", "meter", "inch".
		//	Default is "meter".
		Name string
	}
	//	This is only set-and-retained for imported Collada assets that provide this field,
	//	and is not otherwise used.
	Contributors []*AssetContributor

	//	Provides information about the location of the visual scene in physical space.
	//	This is only set-and-retained for imported Collada assets that provide this field,
	//	and is not otherwise used.
	Coverage *AssetGeographicLocation
}

//	Constructor
func NewAsset() (me *Asset) {
	me = &Asset{}
	me.Unit.Meter, me.Unit.Name = 1, "meter"
	return
}

//	Defines authoring information for asset management.
type AssetContributor struct {
	/*
		//	ALL fields are only set-and-retained for imported Collada assets that provide those fields,
		//	and are not otherwise written, read or used.
				Author        string
				AuthorEmail   string
				AuthorWebsite string
				AuthoringTool string
				Comments      string
				Copyright     string
				SourceData    string
	*/
}

//	Provides information about the location of the visual scene in physical space.
//	ALL fields are only set-and-retained for imported Collada assets that provide those fields,
//	and are not otherwise written, read or used.
type AssetGeographicLocation struct {
	// Longitude        float64
	// Latitude         float64
	// Altitude         float64
	// AltitudeAbsolute bool
}

//	Provides arbitrary additional information about or related to its parent resource.
type Extra struct {
	//	Id
	HasId

	//	Name
	HasName

	//	Asset
	HasAsset

	//	Techniques
	HasTechniques

	//	A hint as to the type of information that this particular Extra represents.
	Type string
}

//	Used in various geometry primitives and b-rep resources.
type IndexedInputs struct {
	//	Number of primitives
	Count uint64

	//	Inputs specify how to read data from Sources.
	Inputs []*InputShared

	//	Indices that describe the attributes for a number of primitives.
	//	The indices reference into the Sources that are referenced by the Inputs.
	Indices []uint64

	//	Number of sub-primitives, if used.
	Vcount []int64
}

//	Declares unshared input semantics of a data source and connects a consumer to that source.
type Input struct {
	//	The user-defined meaning of the input connection.
	Semantic string

	//	Refers to the Source for this Input.
	Source RefId
}

//	Declares shared input semantics of a data source and connects a consumer to that source.
type InputShared struct {
	//	Semantic and Source
	Input

	//	The offset into the list of indices.
	Offset uint64

	//	Which inputs to group as a single set.
	//	This is helpful when multiple inputs share the same semantics.
	Set *uint64
}

//	Allows simple association of resources with custom named layers.
type Layers map[string]bool

//	Binds a specific material to a piece of geometry,
//	binding varying and uniform parameters at the same time.
type MaterialBinding struct {
	//	Extras
	HasExtras

	//	Techniques
	HasTechniques

	//	Targets for animation
	Params []*Param

	//	Common-technique profile.
	TC struct {
		//	References to the materials included in this material binding.
		Materials []*FxMaterialInst
	}
}

//	Declares parametric information for its parent resource.
type Param struct {
	//	Name
	HasName

	//	Sid
	HasSid

	//	The user-defined meaning of the parameter.
	Semantic string

	//	The type of the value data. This text string must be understood by the application.
	Type string
}

//	Declares a new parameter for its parent resource, and assigns it an initial value.
type ParamDef struct {
	//	Sid
	HasSid

	//	Initial value for this parameter
	Value interface{}
}

//	A hash-table containing parameter declarations of this resource.
type ParamDefs map[string]*ParamDef

//	Assigns a new value to a previously defined parameter.
type ParamInst struct {
	//	References the identifier of the pre-defined parameter (ParamDef) that will have its value set.
	Ref RefSid

	//	Indicates if the Value is a string referencing the identifier of a connected parameter.
	IsConnectParamRef bool

	//	The new value for the referenced parameter.
	Value interface{}
}

//	A hash-table containing parameter values assigned by this resource.
type ParamInsts map[string]*ParamInst

//	Declares a complete, self-contained base of a scene hierarchy or scene graph.
type Scene struct {
	//	Extras
	HasExtras

	//	Embodies the entire set of information that can be visualized from the contents of a resource.
	Visual *VisualSceneInst

	//	Embodies the entire set of information that can be articulated kinematically from a resource.
	Kinematics *KxSceneInst

	//	Specifies an environment in which physical objects are instantiated and simulated.
	Physics []*PxSceneInst
}

//	Declares platform-specific or program-specific information
//	used to process some portion of the content.
type Technique struct {
	//	The type of profile. This is a vendor-defined character string
	//	that indicates the platform or capability target for the technique.
	Profile string

	//	Arbitrary XML content or meta-data for this Technique.
	Data string
}

//	Represents a single transformation of a specific kind.
type Transform struct {
	//	Sid
	HasSid

	//	The type of this transformation (rotation, skewing, scaling, translation, "look-at", or matrix).
	//	Must be one of the TransformKind* enumerated constants.
	Kind TransformKind

	//	Contains one or more vectors and values representing this transformation.
	//	For TransformKindLookat:
	//		9 values representing three 3D vectors (eye position, interest point, up-axis).
	//	For TransformKindMatrix:
	//		16 values representing one column-major 4x4 matrix.
	//	For TransformKindSkew:
	//		7 values -- one angle in degrees, then two 3D vectors for the axes of rotation and translation.
	//	For TransformKindRotate:
	//		4 values -- one 3D vector specifying the axis of rotation, then one angle in degrees.
	//	For TransformKindTranslate or TransformKindScale:
	//		3 values representing a single 3D vector.
	F []float64
}
