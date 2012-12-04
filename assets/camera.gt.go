package assets

//	Defines a perspective or orthographic camera. Only perspective cameras are supported at this point.
type CameraDef struct {
	BaseDef

	// FovX float64

	//	Vertical field-of-view for perspective camera.
	FovY float64

	//	Horizontal magnification for orthographic camera.
	MagX float64

	//	Verticial magnification for orthographic camera.
	MagY float64

	//	The distance of the lens to the far-plane. Camera cannot see anything behind the far-plane.
	Zfar float64

	//	The distance of the lens to the near-plane. Camera cannot see anything in front of the near-plane.
	Znear float64

	//	Specifies whether this camera is an orthographic (rather than a perspective) camera.
	Ortho bool
}

func (me *CameraDef) init() {
}

/*
//	Sets the *FovX* field for this *CameraDef* and calls its *Base.SetDirty()* to register the change for syncing.
func (me *CameraDef) SetFovX(v float64) {
	if me.FovX != v {
		me.FovX = v
		me.SetDirty()
	}
}
*/

//	Sets the *FovY* field for this *CameraDef* and calls its *Base.SetDirty()* to register the change for syncing.
func (me *CameraDef) SetFovY(v float64) {
	if me.FovY != v {
		me.FovY = v
		me.SetDirty()
	}
}

//	Sets the *MagX* field for this *CameraDef* and calls its *Base.SetDirty()* to register the change for syncing.
func (me *CameraDef) SetMagX(v float64) {
	if me.MagX != v {
		me.MagX = v
		me.SetDirty()
	}
}

//	Sets the *MagY* field for this *CameraDef* and calls its *Base.SetDirty()* to register the change for syncing.
func (me *CameraDef) SetMagY(v float64) {
	if me.MagY != v {
		me.MagY = v
		me.SetDirty()
	}
}

//	Sets the *Ortho* field for this *CameraDef* and calls its *Base.SetDirty()* to register the change for syncing.
func (me *CameraDef) SetOrtho(v bool) {
	if me.Ortho != v {
		me.Ortho = v
		me.SetDirty()
	}
}

//	Sets the *Zfar* field for this *CameraDef* and calls its *Base.SetDirty()* to register the change for syncing.
func (me *CameraDef) SetZfar(v float64) {
	if me.Zfar != v {
		me.Zfar = v
		me.SetDirty()
	}
}

//	Sets the *Znear* field for this *CameraDef* and calls its *Base.SetDirty()* to register the change for syncing.
func (me *CameraDef) SetZnear(v float64) {
	if me.Znear != v {
		me.Znear = v
		me.SetDirty()
	}
}

//	An instance referencing a camera definition.
type CameraInst struct {
	BaseInst

	//	The camera definition referenced by this instance.
	Def *CameraDef
}

func (me *CameraInst) init() {
}

//#begin-gt _definstlib.gt T:Camera

func newCameraDef(id string) (me *CameraDef) {
	me = &CameraDef{}
	me.BaseDef.init(id)
	me.init()
	return
}

//	Creates and returns a new *CameraInst* instance referencing this *CameraDef* definition.
func (me *CameraDef) NewInst(id string) (inst *CameraInst) {
	inst = &CameraInst{Def: me}
	inst.Base.init(id)
	inst.init()
	return
}

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
	me.Base.init(id)
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
