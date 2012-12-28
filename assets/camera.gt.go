package assets

//	Represents the image sensor of a camera (for example, film or CCD).
type CameraImager struct {
	//	Extras
	HasExtras

	//	Techniques
	HasTechniques
}

//	Represents the apparatus on a camera that projects the image onto the image sensor.
type CameraOptics struct {
	//	Extras
	HasExtras

	//	Techniques
	HasTechniques

	//	Common-technique profile.
	TC struct {
		//	Aspect ratio of the field of view.
		AspectRatio *SidFloat

		//	Distance to the far clipping plane.
		Zfar SidFloat

		//	Distance to the near clipping plane.
		Znear SidFloat

		//	Orthographic projection type. To use Perspective instead, also set this to nil.
		Orthographic *CameraOrthographic

		//	Perspective projection type. To use Orthographic instead, also set this to nil.
		Perspective *CameraPerspective
	}
}

//	Describes the field of view of an orthographic camera.
type CameraOrthographic struct {
	//	Horizontal magnification of the view.
	MagX *SidFloat

	//	Vertical magnification of the view.
	MagY *SidFloat
}

//	Describes the field of view of a perspective camera.
type CameraPerspective struct {
	//	Horizontal field of view in degrees.
	FovX *SidFloat

	//	Vertical field of view in degrees.
	FovY *SidFloat
}

//	Declares a view of the visual scene hierarchy or scene graph.
type CameraDef struct {
	//	Id, Name, Asset, Extras
	BaseDef

	//	Describes the field of view and viewing frustum using canonical parameters.
	Optics CameraOptics

	//	Represents the image sensor of a camera.
	Imager *CameraImager
}

//	Initialization
func (me *CameraDef) Init() {
}

//	Instantiates a camera resource.
type CameraInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst

	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default (unless created via Def.NewInst()) and meant to be set ONLY by
	//	the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *CameraDef
}

//	Initialization
func (me *CameraInst) Init() {
}

//#begin-gt _definstlib.gt T:Camera

func newCameraDef(id string) (me *CameraDef) {
	me = &CameraDef{}
	me.Id = id
	me.BaseSync.init()
	me.Init()
	return
}

//	Returns "the default CameraInst instance" referencing this CameraDef definition.
//	That instance is created once when this method is first called on me,
//	and will have its Def field readily set to me.
func (me *CameraDef) DefaultInst() (inst *CameraInst) {
	if inst = defaultCameraInsts[me]; inst == nil {
		inst = me.NewInst()
		defaultCameraInsts[me] = inst
	}
	return
}

//	Creates and returns a new CameraInst instance referencing this CameraDef definition.
//	Any CameraInst created by this method will have its Def field readily set to me.
func (me *CameraDef) NewInst() (inst *CameraInst) {
	inst = &CameraInst{Def: me}
	inst.DefRef = RefId(me.Id)
	inst.Init()
	return
}

//	If me is "dirty" or me.Def is nil, sets me.Def to the correct CameraDef
//	according to the current me.DefRef value (by searching AllCameraDefLibs).
//	Then returns me.Def.
//	(Note, every CameraInst's Def is nil initially, unless it was created via CameraDef.NewInst().)
func (me *CameraInst) EnsureDef() *CameraDef {
	if (me.Def == nil) || me.dirty {
		me.Def = me.DefRef.CameraDef()
	}
	return me.Def
}

var (
	//	A hash-table that contains LibCameraDefs libraries associated by their Id.
	AllCameraDefLibs = LibsCameraDef{}

	//	The "default" LibCameraDefs library for CameraDefs.
	CameraDefs = AllCameraDefLibs.AddNew("")

	defaultCameraInsts = map[*CameraDef]*CameraInst{}
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllCameraDefLibs {
			lib.SyncChanges()
		}
	})
}

//	Searches (all LibCameraDefs contained in AllCameraDefLibs) for the CameraDef
//	whose Id is referenced by me, returning the first match found.
func (me RefId) CameraDef() (def *CameraDef) {
	id := me.S()
	for _, lib := range AllCameraDefLibs {
		if def = lib.M[id]; def != nil {
			return
		}
	}
	return
}

//	The underlying type of the global AllCameraDefLibs variable:
//	a hash-table that contains LibCameraDefs libraries associated by their Id.
type LibsCameraDef map[string]*LibCameraDefs

//	Creates a new LibCameraDefs library with the specified Id, adds it to this LibsCameraDef, and returns it.
//	If this LibsCameraDef already contains a LibCameraDefs library with the specified Id, does nothing and returns nil.
func (me LibsCameraDef) AddNew(id string) (lib *LibCameraDefs) {
	if me[id] != nil {
		return
	}
	lib = me.new(id)
	me[id] = lib
	return
}

func (me LibsCameraDef) new(id string) (lib *LibCameraDefs) {
	lib = newLibCameraDefs(id)
	return
}

//	A library that contains CameraDefs associated by their Id.
//	To create a new LibCameraDefs library, ONLY use the LibsCameraDef.New() or LibsCameraDef.AddNew() methods.
type LibCameraDefs struct {
	//	Id, Name
	BaseLib

	//	The underlying hash-table. NOTE -- this is for easier read-access and range-iteration:
	//	DO NOT write to M, instead use the Add(), AddNew(), Remove() methods ONLY or bugs WILL ensue.
	M map[string]*CameraDef
}

func newLibCameraDefs(id string) (me *LibCameraDefs) {
	me = &LibCameraDefs{M: map[string]*CameraDef{}}
	me.BaseLib.init(id)
	return
}

//	Adds the specified CameraDef definition to this LibCameraDefs, and returns it.
//	If this LibCameraDefs already contains a CameraDef definition with the same Id, does nothing and returns nil.
func (me *LibCameraDefs) Add(d *CameraDef) (n *CameraDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new CameraDef definition with the specified Id, adds it to this LibCameraDefs, and returns it.
//	If this LibCameraDefs already contains a CameraDef definition with the specified Id, does nothing and returns nil.
func (me *LibCameraDefs) AddNew(id string) *CameraDef { return me.Add(me.New(id)) }

//	Convenience short-hand for len(lib.M)
func (me *LibCameraDefs) Len() int { return len(me.M) }

//	Creates a new CameraDef definition with the specified Id and returns it,
//	but does not add it to this LibCameraDefs.
func (me *LibCameraDefs) New(id string) (def *CameraDef) { def = newCameraDef(id); return }

//	Removes the CameraDef with the specified Id from this LibCameraDefs.
func (me *LibCameraDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

func (me *LibCameraDefs) resolver(part0 string) refSidResolver {
	return me.M[part0]
}

func (me *LibCameraDefs) resolverRootIsLib() bool {
	return true
}

//	Signals to the core package (or your custom package) that changes have been made to this LibCameraDefs
//	that need to be picked up. Call this after you have made a number of changes to this LibCameraDefs
//	library or its CameraDef definitions. Also called by the global SyncChanges() function.
func (me *LibCameraDefs) SyncChanges() {
	me.BaseLib.BaseSync.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.BaseSync.SyncChanges()
	}
}

//#end-gt
