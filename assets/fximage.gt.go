package assets

//	Categorizes the Channels of a FxCreateFormatHint.
type FxFormatChannels int

const (
	_ = iota
	//	Depth map, often used for displacement, parellax, relief, or shadow mapping.
	FxFormatChannelsD FxFormatChannels = iota
	//	Luminance map, often used for light mapping.
	FxFormatChannelsL
	//	Luminance with alpha map, often used for light mapping.
	FxFormatChannelsLa
	//	RGB color map
	FxFormatChannelsRgb
	//	RGB color with alpha map. Often used for color plus transparency
	//	or other things packed into channel A, such as specular power.
	FxFormatChannelsRgba
	//	RGB color with shared exponent for HDR.
	FxFormatChannelsRgbe
)

//	Categorizes the Precision of a FxCreateFormatHint.
type FxFormatPrecision int

const (
	//	Designer does not care as long as it provides "reasonable" precision and performance.
	FxFormatPrecisionDefault FxFormatPrecision = iota
	//	For integers, this typically represents 16 to 32 bits. For floating points, typically
	//	24 to 32 bits.
	FxFormatPrecisionHigh
	//	For integers, this typically represents 8 bits. For floating points, typically 16 bits.
	FxFormatPrecisionLow
	//	Typically 32 bits or 64 bits if available. 64 bits has been separated into its own category
	//	beyond HIGH because it typically has significant performance impact.
	FxFormatPrecisionMax
	//	For integers, this typically represents 8 to 24 bits.
	//	For floating points, typically 16 to 32 bits.
	FxFormatPrecisionMid
)

//	Categorizes the Range of a FxCreateFormatHint.
type FxFormatRange int

const (
	_ = iota
	//	Format should support full floating-point ranges.
	//	High precision is expected to be 32 bits.
	//	Mid precision may be 16 to 32 bits.
	//	Low precision is expected to be 16 bits.
	FxFormatRangeFloat FxFormatRange = iota
	//	Format represents signed integer numbers. For example, 8 bits is -128 to 127.
	FxFormatRangeSint
	//	Format represents a decimal value that remains within the -1 to 1 range.
	//	Implementation could be integer-fixed-point or floating point.
	FxFormatRangeSnorm
	//	Format represent unsigned integer numbers. For example, 8 bits is 0 to 255.
	FxFormatRangeUint
	//	Format represents a decimal value that remains within the 0 to 1 range.
	//	Implementation could be integer-fixed-point or floating point.
	FxFormatRangeUnorm
)

//	Categorizes one of the six sub-images (faces) in a cube-map.
type FxCubeFace int

const (
	//	Cube-map face "X negative"
	FxCubeFaceNegativeX FxCubeFace = 0x8516
	//	Cube-map face "Y negative"
	FxCubeFaceNegativeY FxCubeFace = 0x8518
	//	Cube-map face "Z negative"
	FxCubeFaceNegativeZ FxCubeFace = 0x851A
	//	Cube-map face "X positive"
	FxCubeFacePositiveX FxCubeFace = 0x8515
	//	Cube-map face "Y positive"
	FxCubeFacePositiveY FxCubeFace = 0x8517
	//	Cube-map face "Z positive"
	FxCubeFacePositiveZ FxCubeFace = 0x8519
)

//	Fields shared by FxCreate2D, FxCreate3D and FxCreateCube
type FxCreate struct {
	//	Specifies the length of the 2D array, 3D array or cube-map array.
	ArrayLength uint64
	//	Specifies an image's pixel or compression format.
	//	If not present, the format is assumed to be R8G8B8A8 linear.
	Format *FxCreateFormat
}

//	Assists in the manual creation of a 2D FxImageDef asset.
type FxCreate2D struct {
	//	ArrayLength and Format
	FxCreate
	//	Either Exact or Ratio, but not both, must be present.
	Size struct {
		//	Specifies that the surface should be sized to these exact dimensions.
		Exact *FxCreate2DSizeExact
		//	Specifies that the image size should be relative to the size of the viewport.
		Ratio *FxCreate2DSizeRatio
	}
	//	MIP information. Ignored if Unnormalized is true.
	Mips *FxCreateMips
	//	Unnormalized addressing of texels. (0-W, 0-H).
	Unnormalized bool
	//	Specifies which 2D image to initialize and which MIP level to initialize.
	InitFrom []*FxCreateInitFrom
}

//	Specifies that the surface should be sized to these exact dimensions.
type FxCreate2DSizeExact struct {
	//	width in pixels
	Width uint64
	//	height in pixels
	Height uint64
}

//	Specifies that the image size should be relative to the size of the viewport. For example, 1,1 is the
//	same size as the viewport; 0.5,0.5 is 1/4 the size of the viewport and half as long in either direction.
type FxCreate2DSizeRatio struct {
	//	Relative width where 1.0 represents viewport width.
	Width float64
	//	Relative height where 1.0 represents viewport height.
	Height float64
}

//	Assists in the manual creation of a 3D FxImageDef asset.
type FxCreate3D struct {
	//	ArrayLength and Format
	FxCreate
	//	Specifies that the surface should be sized to these exact dimensions.
	Size struct {
		//	Width in pixels for this 3D texture.
		Width uint64
		//	Height in pixels for this 3D texture.
		Height uint64
		//	Depth in pixels for this 3D texture.
		Depth uint64
	}
	//	MIP information.
	Mips FxCreateMips
	//	Specifies which 3D image to initialize and which MIP level to initialize.
	InitFrom []*FxCreate3DInitFrom
}

//	Initializes an entire 3D texture or portions of a 3D texture from referenced or embedded data.
type FxCreate3DInitFrom struct {
	//	Raw or RefUrl, ArrayIndex and MipIndex
	FxCreateInitFrom
	//	Specifies the slice (depth level) within the MIP that is to be initialized.
	Depth uint64
}

//	Assists in the manual creation of a cube-map FxImageDef asset.
type FxCreateCube struct {
	//	ArrayLength and Format
	FxCreate
	//	Specifies that the cube surfaces should be sized to these exact dimensions.
	Size struct {
		//	Width and height are identical across all faces in a cube-map.
		Width uint64
	}
	//	MIP information.
	Mips FxCreateMips
	//	Specifies which cube image to initialize, which MIP level to initialize,
	//	and which cube face within the MIP that is to be initialized.
	InitFrom []*FxCreateCubeInitFrom
}

//	Initializes an entire cube-map or portions of a cube-map from referenced or embedded data.
type FxCreateCubeInitFrom struct {
	//	Raw or RefUrl, ArrayIndex and MipIndex
	FxCreateInitFrom
	//	Specifies the cube-map face within the MIP that is to be initialized.
	//	Must be one of the FxCubeFace* enumerated constants.
	Face FxCubeFace
}

//	Describes the formatting or memory layout expected of an FxImageDef asset.
type FxCreateFormat struct {
	//	Contains a string representing the profile- and platform-specific texel format
	//	that the author would like this surface to use. If this element is not specified,
	//	or if it is specified but the application cannot process the specified format,
	//	then the application uses the Hint.
	Exact string
	//	If this is not set, then use a common format R8G8B8A8 with linear color gradient, not sRGB.
	Hint *FxCreateFormatHint
}

//	Describes features and characteristics to select an appropriate format for image creation.
type FxCreateFormatHint struct {
	//	Describes the per-texel layout of the format.
	//	Must be one of the FxFormatChannels* enumerated constants.
	Channels FxFormatChannels
	//	Describes the range of texel channel values. Each channel represents a range of values.
	//	Some example ranges are signed or unsigned integers, or are within a clamped range
	//	such as 0.0f to 1.0f, or are a high dynamic range via floating point.
	//	Must be one of the FxFormatRange* enumerated constants.
	Range FxFormatRange
	//	Identifies the precision of the texel channel value.
	//	Must be one of the FxFormatPrecision* enumerated constants.
	Precision FxFormatPrecision
	//	Optional custom / application-specific color-space information.
	Space string
}

//	Initializes an entire image or portions of an image from referenced or embedded data.
type FxCreateInitFrom struct {
	//	Raw and RefUrl
	FxInitFrom
	//	Specifies which array element in the image to initialize (fill).
	ArrayIndex uint64
	//	Specifies which MIP level in the image to initialize.
	MipIndex uint64
}

//	MIP information
type FxCreateMips struct {
	//	Desired number of MIP levels. Special values: 1 is "no MIP levels", 0 is "all MIP levels".
	Levels uint64
	//	If false, initializes higher MIP levels if data does not exist in a file.
	//	If true, no MIP levels are ever automatically initialized.
	NoAutoGen bool
}

//	Initializes an entire image or portions of an image from referenced or embedded data.
type FxImageInitFrom struct {
	//	Raw and RefUrl
	FxInitFrom
	//	If false, initializes higher MIP levels if data does not exist in a file.
	//	If true, no MIP levels are ever automatically initialized.
	NoAutoMip bool
}

//	Constructor
func NewFxImageInitFrom(refUrl string) (me *FxImageInitFrom) {
	me = &FxImageInitFrom{}
	me.RefUrl = refUrl
	return
}

//	Initializes an entire image or portions of an image from referenced or embedded data.
type FxInitFrom struct {
	//	Embedded binary image data; used if RefUrl is empty and Raw.Data is not.
	Raw struct {
		//	Contains the embedded binary image data as a sequence of bytes. Typically contains all
		//	the necessary information including header info such as width and height.
		Data []byte
		//	Specifies which codec decodes the image's descriptions and data.
		//	This is usually its typical file extension, such as "BMP", "JPG", "DDS", "TGA".
		Format string
	}
	//	Contains the URL of a file from which to take initialization data. Assumes the characteristics of the
	//	file: if it is a complex format such as DDS, this might include cube maps, volumes, MIPs, and so on.
	RefUrl string
}

//	Declares the storage for the graphical representation of an object.
type FxImageDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sid
	HasSid
	//	Indicates whether this image represents a render target.
	Renderable struct {
		//	If true, defines the image as a render target, meaning the image can be rendered to.
		Is bool
		//	Indicates whether, when instantiated, the render target is to be shared among all
		//	instances instead of being cloned.
		Shared bool
	}
	//	If set, initializes a custom 2D image by specifying its size, viewport ratio, MIP levels,
	//	normalization, pixel format, and data sources. It also supports arrays of 2D images.
	Create2D *FxCreate2D
	//	If set, initializes a custom 3D image (a volumetric image) by specifying its size, MIP level,
	//	pixel format, and data sources. It also supports arrays of 3D images.
	Create3D *FxCreate3D
	//	If set, initializes the six faces of a cube by specifying its size, MIP level, pixel format,
	//	and data sources. It also supports arrays of images on each of the cube faces.
	//	It also supports arrays of cube images.
	CreateCube *FxCreateCube
	//	If set, initializes the image from a URL (for example, a file) or binary image data.
	InitFrom *FxImageInitFrom
}

//	Initialization
func (me *FxImageDef) Init() {
}

//	Instantiates an image resource.
type FxImageInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *FxImageDef
}

//	Initialization
func (me *FxImageInst) Init() {
}

//	Adds multiple FxImageDefs to this library,
//	with each one's Id and InitFrom.RefUrl set from the specified idRefUrls map.
func (me *LibFxImageDefs) AddFromRefUrls(idRefUrls map[string]string) {
	for imgID, refUrl := range idRefUrls {
		me.AddNew(imgID).InitFrom = NewFxImageInitFrom(refUrl)
	}
}

//#begin-gt _definstlib.gt T:FxImage

func newFxImageDef(id string) (me *FxImageDef) {
	me = &FxImageDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Creates and returns a new FxImageInst instance referencing this FxImageDef definition.
//	Any FxImageInst created by this method will have its Def field readily set to me.
func (me *FxImageDef) NewInst() (inst *FxImageInst) {
	inst = &FxImageInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct FxImageDef
//	according to the current me.DefRef value (by searching AllFxImageDefLibs).
//	Then returns me.Def.
//	(Note, every FxImageInst's Def is nil initially, unless it was created via FxImageDef.NewInst().)
func (me *FxImageInst) EnsureDef() *FxImageDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.FxImageDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibFxImageDefs libraries associated by their Id.
	AllFxImageDefLibs = LibsFxImageDef{}

	//	The "default" LibFxImageDefs library for FxImageDefs.
	FxImageDefs = AllFxImageDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllFxImageDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibFxImageDefs contained in AllFxImageDefLibs) for the FxImageDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) FxImageDef() (def *FxImageDef) {
	id := me.S()
	for _, lib := range AllFxImageDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllFxImageDefLibs variable:
//	a hash-table that contains LibFxImageDefs libraries associated by their Id.
type LibsFxImageDef map[string]*LibFxImageDefs

//	Creates a new LibFxImageDefs library with the specified Id, adds it to this LibsFxImageDef, and returns it.
//	If this LibsFxImageDef already contains a LibFxImageDefs library with the specified Id, does nothing and returns nil.
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

//	A library that contains FxImageDefs associated by their Id.
//	To create a new LibFxImageDefs library, ONLY use the LibsFxImageDef.New() or LibsFxImageDef.AddNew() methods.
type LibFxImageDefs struct {
	BaseLib
	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*FxImageDef
}

func newLibFxImageDefs(id string) (me *LibFxImageDefs) {
	me = &LibFxImageDefs{M: map[string]*FxImageDef{}}
	me.Id = id
	return
}

//	Adds the specified FxImageDef definition to this LibFxImageDefs, and returns it.
//	If this LibFxImageDefs already contains a FxImageDef definition with the same Id, does nothing and returns nil.
func (me *LibFxImageDefs) Add(d *FxImageDef) (n *FxImageDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new FxImageDef definition with the specified Id, adds it to this LibFxImageDefs, and returns it.
//	If this LibFxImageDefs already contains a FxImageDef definition with the specified Id, does nothing and returns nil.
func (me *LibFxImageDefs) AddNew(id string) *FxImageDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibFxImageDefs) Len() int { return len(me.M) }

//	Creates a new FxImageDef definition with the specified Id and returns it,
//	but does not add it to this LibFxImageDefs.
func (me *LibFxImageDefs) New(id string) (def *FxImageDef) { def = newFxImageDef(id); return }

//	Removes the FxImageDef with the specified Id from this LibFxImageDefs.
func (me *LibFxImageDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to the core package (or your custom package) that changes have been made to this LibFxImageDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibFxImageDefs
//	library or its FxImageDef definitions. Also called by the global SyncChanges() function.
func (me *LibFxImageDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
