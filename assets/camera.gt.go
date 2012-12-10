package assets

type CameraImager struct {
	HasTechniques
	HasExtras
}

type CameraCommon struct {
	AspectRatio *ScopedFloat
	Zfar        ScopedFloat
	Znear       ScopedFloat
}

type CameraOptics struct {
	HasExtras
	HasTechniques
	TechniqueCommon struct {
		CameraCommon
		Orthographic *CameraOrthographic
		Perspective  *CameraPerspective
	}
}

type CameraOrthographic struct {
	MagX *ScopedFloat
	MagY *ScopedFloat
}

type CameraPerspective struct {
	FovX *ScopedFloat
	FovY *ScopedFloat
}

//	Defines a perspective or orthographic camera. Only perspective cameras are supported at this point.
type CameraDef struct {
	BaseDef
	Optics CameraOptics
	Imager *CameraImager
}

func (me *CameraDef) Init() {
}

type CameraInst struct {
	BaseInst
}

func (me *CameraInst) init() {
}

//#begin-gt _definstlib.gt T:Camera

func newCameraDef(id string) (me *CameraDef) {
	me = &CameraDef{}
	me.ID = id
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
	//	A *map* collection that contains *LibCameraDefs* libraries associated by their *ID*.
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
//	*LibCameraDefs* libraries associated by their *ID*.
type LibsCameraDef map[string]*LibCameraDefs

//	Creates a new *LibCameraDefs* library with the specified *ID*, adds it to this *LibsCameraDef*, and returns it.
//	
//	If this *LibsCameraDef* already contains a *LibCameraDefs* library with the specified *ID*, does nothing and returns *nil*.
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

//	A library that contains *CameraDef*s associated by their *ID*. To create a new *LibCameraDefs* library, ONLY
//	use the *LibsCameraDef.New()* or *LibsCameraDef.AddNew()* methods.
type LibCameraDefs struct {
	BaseLib
	//	The underlying *map* collection. NOTE: this is for easier read-access and range-iteration -- DO NOT
	//	write to *M*, instead use the *Add()*, *AddNew()*, *Remove()* methods ONLY or bugs WILL ensue.
	M map[string]*CameraDef
}

func newLibCameraDefs(id string) (me *LibCameraDefs) {
	me = &LibCameraDefs{M: map[string]*CameraDef{}}
	me.ID = id
	return
}

//	Adds the specified *CameraDef* definition to this *LibCameraDefs*, and returns it.
//	
//	If this *LibCameraDefs* already contains a *CameraDef* definition with the same *ID*, does nothing and returns *nil*.
func (me *LibCameraDefs) Add(d *CameraDef) (n *CameraDef) {
	if me.M[d.ID] == nil {
		n, me.M[d.ID] = d, d
		me.SetDirty()
	}
	return
}

//	Creates a new *CameraDef* definition with the specified *ID*, adds it to this *LibCameraDefs*, and returns it.
//	
//	If this *LibCameraDefs* already contains a *CameraDef* definition with the specified *ID*, does nothing and returns *nil*.
func (me *LibCameraDefs) AddNew(id string) *CameraDef { return me.Add(me.New(id)) }

//	Creates a new *CameraDef* definition with the specified *ID* and returns it, but does not add it to this *LibCameraDefs*.
func (me *LibCameraDefs) New(id string) (def *CameraDef) { def = newCameraDef(id); return }

//	Removes the *CameraDef* with the specified *ID* from this *LibCameraDefs*.
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
