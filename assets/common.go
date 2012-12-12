package assets

import (
	xmlx "github.com/jteeuwen/go-pkg-xmlx"
)

const (
	//	A position and orientation transformation suitable for aiming a camera.
	TRANSFORM_TYPE_LOOKAT = 1
	//	A transformation that embodies mathematical changes to points within a coordinate system or the coordinate system itself.
	TRANSFORM_TYPE_MATRIX = iota
	//	A transformation that specifies how to rotate an object around an axis.
	TRANSFORM_TYPE_ROTATE = iota
	//	A transformation that specifies how to deform an object along one axis.
	TRANSFORM_TYPE_SKEW = iota
	//	A transformation that specifies how to change an object's size.
	TRANSFORM_TYPE_SCALE = iota
	//	A transformation that changes the position of an object in a local coordinate system.
	TRANSFORM_TYPE_TRANSLATE = iota
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
	//	The optional pretty-print name/title of this *Def*, *Inst* or *Lib*.
	Name string
}

//	Used in all resources that declare their own parameters.
type HasParamDefs struct {
	//	A hash-table containing all parameter declarations of this resource.
	NewParams ParamDefs
}

//	Used in all resources that assign values to other parameters.
type HasParamInsts struct {
	//	All parameter values assigned by this resource.
	SetParams []*ParamInst
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
	//	Custom-technique/foreign-profile meta-data.
	HasExtras
	//	Contains the date and time that the parent element was created.
	//	This is only set-and-retained for imported Collada assets that provide this field, and is not otherwise used.
	Created string
	//	Contains the date and time that the parent element was last modified.
	//	This is only set-and-retained for imported Collada assets that provide this field, and is not otherwise used.
	Modified string
	//	Contains a list of words used as search criteria.
	//	This is only set-and-retained for imported Collada assets that provide this field, and is not otherwise used.
	Keywords string
	//	Contains revision information.
	//	This is only set-and-retained for imported Collada assets that provide this field, and is not otherwise used.
	Revision string
	//	Contains a description of the topical subject.
	//	This is only set-and-retained for imported Collada assets that provide this field, and is not otherwise used.
	Subject string
	//	Contains title information.
	//	This is only set-and-retained for imported Collada assets that provide this field, and is not otherwise used.
	Title string
	//	Contains descriptive information about the coordinate system of the geometric data. All
	//	coordinates are right-handed by definition. Valid values are "X", "Y" (the default), or "Z".
	UpAxis string
	//	The unit of distance that applies to all spatial measurements within the scope of this resource.
	Unit struct {
		//	How many real-world meters in one distance unit as a floating-point number. 1.0 for meter, 0.01 for centimeter, 1000 for kilometer etc.
		Meter float64
		//	Name of the distance unit, such as "centimeter", "kilometer", "meter", "inch". Default is "meter".
		Name string
	}
	//	This is only set-and-retained for imported Collada assets that provide this field, and is not otherwise used.
	Contributors []*AssetContributor
	//	Provides information about the location of the visual scene in physical space.
	//	This is only set-and-retained for imported Collada assets that provide this field, and is not otherwise used.
	Coverage *AssetGeographicLocation
}

//	Constructor
func NewAsset() (me *Asset) {
	me = &Asset{}
	me.Unit.Meter, me.Unit.Name = 1, "meter"
	return
}

//	Defines authoring information for asset management.
//	ALL fields are only set-and-retained for imported Collada assets that provide those field, and are not otherwise written, read or used.
type AssetContributor struct {
	Author        string
	AuthorEmail   string
	AuthorWebsite string
	AuthoringTool string
	Comments      string
	Copyright     string
	SourceData    string
}

//	Provides information about the location of the visual scene in physical space.
//	ALL fields are only set-and-retained for imported Collada assets that provide those field, and are not otherwise written, read or used.
type AssetGeographicLocation struct {
	Longitude        float64
	Latitude         float64
	Altitude         float64
	AltitudeAbsolute bool
}

type BindMaterial struct {
	HasExtras
	HasTechniques
	Params []*Param
	TC     struct {
		Materials []*FxMaterialInst
	}
}

type Extra struct {
	HasId
	HasName
	HasAsset
	HasTechniques
	Type string
}

type IndexedInputs struct {
	Count   uint64
	Inputs  []*InputShared
	Indices []uint64
}

type IndexedInputsV struct {
	IndexedInputs
	Vcount []int64
}

type Input struct {
	Semantic string
	Source   string
}

type InputShared struct {
	Input
	Offset uint64
	Set    *uint64
}

type Layers map[string]bool

type Param struct {
	HasName
	HasSid
	Semantic string
	Type     string
}

type ParamDef struct {
	HasSid
	Value interface{}
}

type ParamDefs map[string]*ParamDef

type ParamInst struct {
	Ref               string
	IsConnectParamRef bool
	Value             interface{}
}

type RefId string

func (me RefId) S() string {
	return string(me)
}

func (me *RefId) Set(v string) {
	*me = RefId(v)
}

type RefSid string

func (me RefSid) S() string {
	return string(me)
}

func (me *RefSid) Set(v string) {
	*me = RefSid(v)
}

type Technique struct {
	Profile string
	Data    []*xmlx.Node
}

type Transform struct {
	HasSid
	Type int
	F    []float64
}
