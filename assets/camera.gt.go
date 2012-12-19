package assets

//	Represents the image sensor of a camera (for example, film or CCD).
type CameraImager struct {
	//	Custom-profile/foreign-technique meta-data
	HasExtras
	//	Custom-profile/foreign-technique support
	HasTechniques
}

//	Represents the apparatus on a camera that projects the image onto the image sensor.
type CameraOptics struct {
	//	Custom-profile/foreign-technique meta-data
	HasExtras
	//	Custom-profile/foreign-technique support
	HasTechniques
	//	Common-technique profile.
	TC struct {
		//	Aspect ratio of the field of view.
		AspectRatio *ScopedFloat
		//	Distance to the far clipping plane.
		Zfar ScopedFloat
		//	Distance to the near clipping plane.
		Znear ScopedFloat
		//	Orthographic projection type. To use Perspective instead, also set this to nil.
		Orthographic *CameraOrthographic
		//	Perspective projection type. To use Orthographic instead, also set this to nil.
		Perspective *CameraPerspective
	}
}

//	Describes the field of view of an orthographic camera.
type CameraOrthographic struct {
	//	Horizontal magnification of the view.
	MagX *ScopedFloat
	//	Vertical magnification of the view.
	MagY *ScopedFloat
}

//	Describes the field of view of a perspective camera.
type CameraPerspective struct {
	//	Horizontal field of view in degrees.
	FovX *ScopedFloat
	//	Vertical field of view in degrees.
	FovY *ScopedFloat
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
}

//	Initialization
func (me *CameraInst) Init() {
}

//#begin-gt _definstlib.gt T:Camera

func newCameraDef(id string) (me *CameraDef) {
	me = &CameraDef{}
	me.Id = id
	me.Base.init()
	me.Init()
	return
}

/*
//	Creates and returns a new *CameraInst* instance referencing this *CameraDef* definition.
func (me *CameraDef) NewInst(id string) (inst *CameraInst) {
	inst = &CameraInst{Def: me}
	inst.Init()
	return
}
*/

var (
	//	A *map* collection that contains *LibCameraDefs* libraries associated by their *Id*.
	AllCameraDefLibs = LibsCameraDef{}

	//	The "default" *LibCameraDefs* library for *CameraDef*s.
	CameraDefs = AllCameraDefLibs.AddNew("")
)

func init() {
	syncHandlers = append(syncHandlers, func() {
		for _, lib := range AllCameraDefLibs {
			lib.SyncChanges()
		}
	})
}

//	The underlying type of the global *AllCameraDefLibs* variable: a *map* collection that contains
//	*LibCameraDefs* libraries associated by their *Id*.
type LibsCameraDef map[string]*LibCameraDefs

//	Creates a new *LibCameraDefs* library with the specified *Id*, adds it to this *LibsCameraDef*, and returns it.
//	
//	If this *LibsCameraDef* already contains a *LibCameraDefs* library with the specified *Id*, does nothing and returns *nil*.
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

//	A library that contains *CameraDef*s associated by their *Id*. To create a new *LibCameraDefs* library, ONLY
//	use the *LibsCameraDef.New()* or *LibsCameraDef.AddNew()* methods.
type LibCameraDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*CameraDef
}

func newLibCameraDefs(id string) (me *LibCameraDefs) {
	me = &LibCameraDefs{M: map[string]*CameraDef{}}
	me.Id = id
	return
}

//	Adds the specified *CameraDef* definition to this *LibCameraDefs*, and returns it.
//	
//	If this *LibCameraDefs* already contains a *CameraDef* definition with the same *Id*, does nothing and returns *nil*.
func (me *LibCameraDefs) Add(d *CameraDef) (n *CameraDef) {
	if me.M[d.Id] == nil {
		n, me.M[d.Id] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *CameraDef* definition with the specified *Id*, adds it to this *LibCameraDefs*, and returns it.
//	
//	If this *LibCameraDefs* already contains a *CameraDef* definition with the specified *Id*, does nothing and returns *nil*.
func (me *LibCameraDefs) AddNew(id string) *CameraDef { return me.Add(me.New(id)) }

//	Short-hand for len(lib.M)
func (me *LibCameraDefs) Len() int { return len(me.M) }

//	Creates a new *CameraDef* definition with the specified *Id* and returns it, but does not add it to this *LibCameraDefs*.
func (me *LibCameraDefs) New(id string) (def *CameraDef) { def = newCameraDef(id); return }

//	Removes the *CameraDef* with the specified *Id* from this *LibCameraDefs*.
func (me *LibCameraDefs) Remove(id string) { delete(me.M, id); me.SetDirty() }

//	Signals to *core* (or your custom package) that changes have been made to this *LibCameraDefs* that need to be picked up.
//	Call this after you have made any number of changes to this *LibCameraDefs* library or its *CameraDef* definitions.
//	Also called by the global *SyncChanges()* function.
func (me *LibCameraDefs) SyncChanges() {
	me.BaseLib.Base.SyncChanges()
	for _, def := range me.M {
		def.BaseDef.Base.SyncChanges()
	}
}

//#end-gt
