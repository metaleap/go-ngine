# core
--
    import "github.com/go3d/go-ngine/core"

The *core* package provides go:ngine client-side core functionality such as
rendering and user I/O.

## Usage

```go
const DefaultBadVersionMessage = `
Minimum required OpenGL version is {MINVER}, but your
current {OS} graphics driver only provides
OpenGL version {CURVER}.

Most likely your machine is just
missing some recent system updates.

*HOW TO RESOLVE*:

1. On the web, search for "downloading & installing the
latest {OS} driver for {GPU}",

2. or simply visit the {VENDOR} website and locate
their "driver downloads" pages to obtain the most
recent driver for {GPU}.
`
```

```go
var (
	Diag = NgDiag{
		LogCategories: DiagLogCatAll,
		LogCategoryNames: map[NgDiagLogCategory]string{
			DiagLogCatMisc:    "[INFO]\t\t",
			DiagLogCatShaders: "[SHADER]\t",
			DiagLogCatImages:  "[IMAGES]\t",
			DiagLogCatMeshes:  "[MESHES]\t",
		},
		LogGLErrorsInLoopOnSec: false,
	}
)
```

#### func  Dispose

```go
func Dispose()
```
Call this to "un-init" go:ngine and to release any and all GPU or RAM resources
still allocated.

#### func  Init

```go
func Init(fullscreen bool, ctx ngctx.CtxProvider) (err error)
```
Initializes go:ngine; this first attempts to initialize OpenGL and then open a
window to your supplied specifications with a GL 3.3-or-higher profile.

#### type Camera

```go
type Camera struct {
	//	Optical and imager properties for this camera.
	Perspective u3d.Perspective

	//	Encapsulates the position and direction of this camera.
	Controller Controller

	Cull struct {
		Frustum bool
	}
}
```


#### func (*Camera) Scene

```go
func (me *Camera) Scene() *Scene
```

#### func (*Camera) SetScene

```go
func (me *Camera) SetScene(sceneID int)
```

#### type Controller

```go
type Controller struct {
	//	The position being manipulated by this Controller.
	//	When manipulating this manually (outside the TurnXyz() / MoveXyz() methods),
	//	do so in between calling the BeginUpdate() and EndUpdate() methods.
	Pos unum.Vec3

	//	Indicates which axis is consider "upward". This is typically
	//	the Y-axis, denoted by the default value (0, 1, 0).
	//	When manipulating this manually (outside the TurnXyz() / MoveXyz() methods),
	//	do so in between calling the BeginUpdate() and EndUpdate() methods.
	UpAxis unum.Vec3

	UpVec unum.Vec3

	//	Defaults to a copy of Options.Cameras.DefaultControllerParams
	Params ControllerParams
}
```

Encapsulates a position-and-direction and provides methods manipulating these
with respect to each other (e.g. "move forward" some entity that is rotated
facing some arbitrary direction).

#### func (*Controller) BeginUpdate

```go
func (me *Controller) BeginUpdate()
```
Temporarily suspends all matrix re-calculations typically occuring inside the
MoveXyz() / TurnXyz() methods. Call this prior to multiple subsequent calls to
any combination of those methods, and/or prior to manually modifying the Pos,
Dir or UpAxis fields of me. Immediately afterwards, be sure to call EndUpdate()
to apply all changes in a final matrix re-calculation.

#### func (*Controller) CopyFrom

```go
func (me *Controller) CopyFrom(copy Controller)
```

#### func (*Controller) Dir

```go
func (me *Controller) Dir() *unum.Vec3
```
The direction being manipulated by this Controller. NOTE: this returns a pointer
to the direction vector to avoid a copy, but it's NOT meant to be modified, as
the vector is re-computed by the TurnFoo() methods.

#### func (*Controller) EndUpdate

```go
func (me *Controller) EndUpdate()
```
Applies all changes made to Pos, Dir or UpAxis since BeginUpdate() was last
called, and recalculates this Controller's final 4x4 transformation matrix. Also
resumes all matrix re-calculations typically occuring inside the MoveXyz() /
TurnXyz() methods that were suspended since BeginUpdate().

#### func (*Controller) MoveBackward

```go
func (me *Controller) MoveBackward()
```
Recomputes Pos with regards to UpAxis and Dir to effect a "move backward".

#### func (*Controller) MoveDown

```go
func (me *Controller) MoveDown()
```
Recomputes Pos with regards to UpAxis to effect a "move downward".

#### func (*Controller) MoveForward

```go
func (me *Controller) MoveForward()
```
Recomputes Pos with regards to UpAxis and Dir to effect a "move forward".

#### func (*Controller) MoveLeft

```go
func (me *Controller) MoveLeft()
```
Recomputes Pos with regards to UpAxis and Dir to effect a "move left-ward".

#### func (*Controller) MoveRight

```go
func (me *Controller) MoveRight()
```
Recomputes Pos with regards to UpAxis and Dir to effect a "move right-ward".

#### func (*Controller) MoveUp

```go
func (me *Controller) MoveUp()
```
Recomputes Pos with regards to UpAxis to effect a "move upward".

#### func (*Controller) StepSizeMove

```go
func (me *Controller) StepSizeMove() float64
```
Returns the current distance that a single MoveXyz() call (per loop iteration)
would move. (Loop.TickDelta * me.Params.MoveSpeed * me.Params.MoveSpeedupFactor)

#### func (*Controller) StepSizeTurn

```go
func (me *Controller) StepSizeTurn() float64
```
Returns the current degrees that a single TurnXyz() call (per loop iteration)
would turn. (Loop.TickDelta * me.Params.TurnSpeed * me.Params.TurnSpeedupFactor)

#### func (*Controller) TurnDown

```go
func (me *Controller) TurnDown()
```
Recomputes Dir with regards to UpAxis and Pos to effect a "turn downward" by
me.StepSizeTurn() degrees.

#### func (*Controller) TurnDownBy

```go
func (me *Controller) TurnDownBy(deg float64)
```
Recomputes Dir with regards to UpAxis and Pos to effect a "turn downward" by the
specified degrees.

#### func (*Controller) TurnLeft

```go
func (me *Controller) TurnLeft()
```
Recomputes Dir with regards to UpAxis and Pos to effect a "turn left-ward" by
me.StepSizeTurn() degrees.

#### func (*Controller) TurnLeftBy

```go
func (me *Controller) TurnLeftBy(deg float64)
```
Recomputes Dir with regards to UpAxis and Pos to effect a "turn left-ward" by
the specified degrees.

#### func (*Controller) TurnRight

```go
func (me *Controller) TurnRight()
```
Recomputes Dir with regards to UpAxis and Pos to effect a "turn right-ward" by
me.StepSizeTurn() degrees.

#### func (*Controller) TurnRightBy

```go
func (me *Controller) TurnRightBy(deg float64)
```
Recomputes Dir with regards to UpAxis and Pos to effect a "turn right-ward" by
the specified degrees.

#### func (*Controller) TurnUp

```go
func (me *Controller) TurnUp()
```
Recomputes Dir with regards to UpAxis and Pos to effect a "turn upward" by
me.StepSizeTurn() degrees.

#### func (*Controller) TurnUpBy

```go
func (me *Controller) TurnUpBy(deg float64)
```
Recomputes Dir with regards to UpAxis and Pos to effect a "turn upward" by the
specified degress.

#### type ControllerParams

```go
type ControllerParams struct {
	//	Speed of "moving" in the MoveXyz() methods, in units per second.
	//	Defaults to 2.
	MoveSpeed float64

	//	A factor multiplied with MoveSpeed in the MoveXyz() methods. Defaults to 1.
	MoveSpeedupFactor float64

	//	Speed of "turning" in the TurnXyz() methods, in degrees per second.
	//	Defaults to 90.
	TurnSpeed float64

	//	A factor multiplied with TurnSpeed in the TurnXyz() methods. Defaults to 1.
	TurnSpeedupFactor float64

	//	The maximum degree that TurnUp() allows. Defaults to 90.
	MaxTurnUp float64

	//	The minimum degree that TurnDown() allows. Defaults to -90.
	MinTurnDown float64
}
```


#### type FxEffect

```go
type FxEffect struct {
	//	An ordered collection of all FxProcs that make up this effect.
	//	When changing the ordering or disabling, enabling or toggling individual FxProcs,
	//	you need to call the FxEffect.UpdateRoutine() method to reflect such changes.
	//	All other dynamic, individual FxProc-specific parameter changes
	//	(colors, image bindings, weights etc.pp.) do not require this.
	FxProcs

	ID int

	KeepProcIDsLast []string
}
```

Declares the visual appearance of a surface. An FxEffect can be reused for
multiple surfaces, it is bound to geometry via an FxMaterial.

#### func (*FxEffect) UpdateRoutine

```go
func (me *FxEffect) UpdateRoutine()
```

#### type FxEffectLib

```go
type FxEffectLib []FxEffect
```

Only used for Core.Libs.Effects

#### func (*FxEffectLib) AddNew

```go
func (me *FxEffectLib) AddNew() (id int)
```

#### func (*FxEffectLib) Compact

```go
func (me *FxEffectLib) Compact()
```

#### func (FxEffectLib) IsOk

```go
func (me FxEffectLib) IsOk(id int) (ok bool)
```

#### func (FxEffectLib) Ok

```go
func (me FxEffectLib) Ok(id int) bool
```

#### func (FxEffectLib) Remove

```go
func (me FxEffectLib) Remove(fromID, num int)
```

#### func (FxEffectLib) Walk

```go
func (me FxEffectLib) Walk(on func(ref *FxEffect))
```

#### type FxImage2D

```go
type FxImage2D struct {
	FxImageBase
	InitFrom FxImageInitFrom
}
```


#### func (*FxImage2D) GpuDelete

```go
func (me *FxImage2D) GpuDelete()
```

#### func (*FxImage2D) GpuSync

```go
func (me *FxImage2D) GpuSync() (err error)
```

#### func (*FxImage2D) Load

```go
func (me *FxImage2D) Load() (err error)
```

#### func (*FxImage2D) Loaded

```go
func (me *FxImage2D) Loaded() bool
```

#### func (*FxImage2D) NoAutoMips

```go
func (me *FxImage2D) NoAutoMips()
```

#### func (*FxImage2D) Unload

```go
func (me *FxImage2D) Unload()
```

#### type FxImage2DLib

```go
type FxImage2DLib []FxImage2D
```

Only used for Core.Libs.Images.Tex2D

#### func (*FxImage2DLib) AddNew

```go
func (me *FxImage2DLib) AddNew() (id int)
```

#### func (*FxImage2DLib) Compact

```go
func (me *FxImage2DLib) Compact()
```

#### func (FxImage2DLib) IsOk

```go
func (me FxImage2DLib) IsOk(id int) (ok bool)
```

#### func (FxImage2DLib) Ok

```go
func (me FxImage2DLib) Ok(id int) bool
```

#### func (FxImage2DLib) Remove

```go
func (me FxImage2DLib) Remove(fromID, num int)
```

#### func (FxImage2DLib) Walk

```go
func (me FxImage2DLib) Walk(on func(ref *FxImage2D))
```

#### type FxImageBase

```go
type FxImageBase struct {
	ID         int
	Preprocess FxImagePreprocess
	Storage    FxImageStorage
}
```


#### func (*FxImageBase) GpuSynced

```go
func (me *FxImageBase) GpuSynced() bool
```

#### type FxImageCube

```go
type FxImageCube struct {
	FxImageBase
	InitFrom [6]FxImageInitFrom
}
```


#### func (*FxImageCube) GpuDelete

```go
func (me *FxImageCube) GpuDelete()
```

#### func (*FxImageCube) GpuSync

```go
func (me *FxImageCube) GpuSync() (err error)
```

#### func (*FxImageCube) Load

```go
func (me *FxImageCube) Load() (err error)
```

#### func (*FxImageCube) Loaded

```go
func (me *FxImageCube) Loaded() bool
```

#### func (*FxImageCube) NoAutoMips

```go
func (me *FxImageCube) NoAutoMips()
```

#### func (*FxImageCube) Unload

```go
func (me *FxImageCube) Unload()
```

#### type FxImageCubeLib

```go
type FxImageCubeLib []FxImageCube
```

Only used for Core.Libs.Images.TexCube

#### func (*FxImageCubeLib) AddNew

```go
func (me *FxImageCubeLib) AddNew() (id int)
```

#### func (*FxImageCubeLib) Compact

```go
func (me *FxImageCubeLib) Compact()
```

#### func (FxImageCubeLib) IsOk

```go
func (me FxImageCubeLib) IsOk(id int) (ok bool)
```

#### func (FxImageCubeLib) Ok

```go
func (me FxImageCubeLib) Ok(id int) bool
```

#### func (FxImageCubeLib) Remove

```go
func (me FxImageCubeLib) Remove(fromID, num int)
```

#### func (FxImageCubeLib) Walk

```go
func (me FxImageCubeLib) Walk(on func(ref *FxImageCube))
```

#### type FxImageInitFrom

```go
type FxImageInitFrom struct {
	RawData []byte
	RefUrl  string
}
```


#### type FxImagePreprocess

```go
type FxImagePreprocess struct {
	FlipY    bool
	ToLinear bool
	ToBgra   bool
}
```


#### type FxImageStorage

```go
type FxImageStorage struct {
	DiskCache struct {
		Enabled      bool
		Compressor   func(w io.WriteCloser) io.WriteCloser
		Decompressor func(r io.ReadCloser) io.ReadCloser
	}
	Gpu struct {
		Bgra    bool
		UintRev bool
	}
}
```


#### type FxMaterial

```go
type FxMaterial struct {
	ID int

	//	This effect is used by default for all faces that do not
	//	match any of the selectors in the FaceEffects field.
	DefaultEffectID int

	//	Associates certain individual faces or tags of faces
	//	with specific effect IDs.
	FaceEffects struct {
		//	Associates face tags with effect IDs.
		ByTag map[string]int

		//	Associates specific face IDs with effect IDs.
		ByID map[string]int
	}
}
```

A material binds effects (FxEffect in Core.Libs.Effects) to geometry (Model or
Node).

#### func (*FxMaterial) HasFaceEffects

```go
func (me *FxMaterial) HasFaceEffects() bool
```

#### type FxMaterialLib

```go
type FxMaterialLib []FxMaterial
```

Only used for Core.Libs.Materials

#### func (*FxMaterialLib) AddNew

```go
func (me *FxMaterialLib) AddNew() (id int)
```

#### func (*FxMaterialLib) Compact

```go
func (me *FxMaterialLib) Compact()
```

#### func (FxMaterialLib) IsOk

```go
func (me FxMaterialLib) IsOk(id int) (ok bool)
```

#### func (FxMaterialLib) Ok

```go
func (me FxMaterialLib) Ok(id int) bool
```

#### func (FxMaterialLib) Remove

```go
func (me FxMaterialLib) Remove(fromID, num int)
```

#### func (FxMaterialLib) Walk

```go
func (me FxMaterialLib) Walk(on func(ref *FxMaterial))
```

#### type FxProc

```go
type FxProc struct {
	Enabled bool

	Color struct {
		Rgb ugl.GlVec3
	}

	Tex struct {
		ImageID int

		Sampler ugl.Sampler
	}
}
```


#### func (*FxProc) Color_SetRgb

```go
func (me *FxProc) Color_SetRgb(rgb ...gl.Float) *FxProc
```

#### func (*FxProc) IsColor

```go
func (me *FxProc) IsColor() bool
```

#### func (*FxProc) IsCoords

```go
func (me *FxProc) IsCoords() bool
```

#### func (*FxProc) IsGamma

```go
func (me *FxProc) IsGamma() bool
```

#### func (*FxProc) IsGrayscale

```go
func (me *FxProc) IsGrayscale() bool
```

#### func (*FxProc) IsOrangify

```go
func (me *FxProc) IsOrangify() bool
```

#### func (*FxProc) IsTex

```go
func (me *FxProc) IsTex() bool
```

#### func (*FxProc) IsTex2D

```go
func (me *FxProc) IsTex2D() bool
```

#### func (*FxProc) IsTexCube

```go
func (me *FxProc) IsTexCube() bool
```

#### func (*FxProc) SetMixWeight

```go
func (me *FxProc) SetMixWeight(weight float64)
```

#### func (*FxProc) Tex_SetImageID

```go
func (me *FxProc) Tex_SetImageID(imageID int) *FxProc
```

#### func (*FxProc) Toggle

```go
func (me *FxProc) Toggle()
```

#### type FxProcs

```go
type FxProcs []FxProc
```

Used for FxEffect.Procs and Camera.Rendering.FxProcs.

#### func (FxProcs) Disable

```go
func (me FxProcs) Disable(procID string, n int)
```
Disables the nth (0-based) FxProc with the specified procID, or all FxProcs with
the specified procID if n < 0. The procID must be one of the
Core.Render.Fx.KnownProcIDs. For this change to be applied, call
FxEffect.UpdateRoutine() subsequently.

#### func (FxProcs) DisableColor

```go
func (me FxProcs) DisableColor(n int)
```
Convenience short-hand for me.Disable("Color", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (FxProcs) DisableCoords

```go
func (me FxProcs) DisableCoords(n int)
```
Convenience short-hand for me.Disable("Coords", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (FxProcs) DisableGamma

```go
func (me FxProcs) DisableGamma(n int)
```
Convenience short-hand for me.Disable("Gamma", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (FxProcs) DisableGrayscale

```go
func (me FxProcs) DisableGrayscale(n int)
```
Convenience short-hand for me.Disable("Grayscale", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (FxProcs) DisableOrangify

```go
func (me FxProcs) DisableOrangify(n int)
```
Convenience short-hand for me.Disable("Orangify", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (FxProcs) DisableTex2D

```go
func (me FxProcs) DisableTex2D(n int)
```
Convenience short-hand for me.Disable("Tex2D", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (FxProcs) DisableTexCube

```go
func (me FxProcs) DisableTexCube(n int)
```
Convenience short-hand for me.Disable("TexCube", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) Enable

```go
func (me *FxProcs) Enable(procID string, n int) (proc *FxProc)
```
Enables the nth (0-based) FxProc with the specified procID, or all FxProcs with
the specified procID if n < 0. If me has no FxProc with the specified procID,
appends a new one. The procID must be one of the Core.Render.Fx.KnownProcIDs.
For this change to be applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) EnableColor

```go
func (me *FxProcs) EnableColor(n int) *FxProc
```
Convenience short-hand for me.Enable("Color", n). For this change to be applied,
call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) EnableCoords

```go
func (me *FxProcs) EnableCoords(n int) *FxProc
```
Convenience short-hand for me.Enable("Coords", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) EnableGamma

```go
func (me *FxProcs) EnableGamma(n int) *FxProc
```
Convenience short-hand for me.Enable("Gamma", n). For this change to be applied,
call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) EnableGrayscale

```go
func (me *FxProcs) EnableGrayscale(n int) *FxProc
```
Convenience short-hand for me.Enable("Grayscale", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) EnableOrangify

```go
func (me *FxProcs) EnableOrangify(n int) *FxProc
```
Convenience short-hand for me.Enable("Orangify", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) EnableTex2D

```go
func (me *FxProcs) EnableTex2D(n int) *FxProc
```
Convenience short-hand for me.Enable("Tex2D", n). For this change to be applied,
call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) EnableTexCube

```go
func (me *FxProcs) EnableTexCube(n int) *FxProc
```
Convenience short-hand for me.Enable("TexCube", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (FxProcs) EnsureLast

```go
func (me FxProcs) EnsureLast(lastProcIDs ...string)
```

#### func (FxProcs) Get

```go
func (me FxProcs) Get(procID string, n int) (proc *FxProc)
```
Returns the nth (0-based) FxProc with the specified procID. The procID must be
one of the Core.Render.Fx.KnownProcIDs.

#### func (FxProcs) GetColor

```go
func (me FxProcs) GetColor(n int) *FxProc
```
Convenience short-hand for me.Get("Color", n).

#### func (FxProcs) GetCoords

```go
func (me FxProcs) GetCoords(n int) *FxProc
```
Convenience short-hand for me.Get("Coords", n).

#### func (FxProcs) GetGamma

```go
func (me FxProcs) GetGamma(n int) *FxProc
```
Convenience short-hand for me.Get("Gamma", n).

#### func (FxProcs) GetGrayscale

```go
func (me FxProcs) GetGrayscale(n int) *FxProc
```
Convenience short-hand for me.Get("Grayscale", n).

#### func (FxProcs) GetOrangify

```go
func (me FxProcs) GetOrangify(n int) *FxProc
```
Convenience short-hand for me.Get("Orangify", n).

#### func (FxProcs) GetTex2D

```go
func (me FxProcs) GetTex2D(n int) *FxProc
```
Convenience short-hand for me.Get("Tex2D", n).

#### func (FxProcs) GetTexCube

```go
func (me FxProcs) GetTexCube(n int) *FxProc
```
Convenience short-hand for me.Get("TexCube", n).

#### func (*FxProcs) Toggle

```go
func (me *FxProcs) Toggle(procID string, n int)
```
Toggles the nth (0-based) FxProc with the specified procID, or all FxProcs with
the specified procID if n < 0. If me has no FxProc with the specified procID,
appends a new one. The procID must be one of the Core.Render.Fx.KnownProcIDs.
For this change to be applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) ToggleColor

```go
func (me *FxProcs) ToggleColor(n int)
```
Convenience short-hand for me.Toggle("Color", n). For this change to be applied,
call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) ToggleCoords

```go
func (me *FxProcs) ToggleCoords(n int)
```
Convenience short-hand for me.Toggle("Coords", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) ToggleGamma

```go
func (me *FxProcs) ToggleGamma(n int)
```
Convenience short-hand for me.Toggle("Gamma", n). For this change to be applied,
call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) ToggleGrayscale

```go
func (me *FxProcs) ToggleGrayscale(n int)
```
Convenience short-hand for me.Toggle("Grayscale", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) ToggleOrangify

```go
func (me *FxProcs) ToggleOrangify(n int)
```
Convenience short-hand for me.Toggle("Orangify", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) ToggleTex2D

```go
func (me *FxProcs) ToggleTex2D(n int)
```
Convenience short-hand for me.Toggle("Tex2D", n). For this change to be applied,
call FxEffect.UpdateRoutine() subsequently.

#### func (*FxProcs) ToggleTexCube

```go
func (me *FxProcs) ToggleTexCube(n int)
```
Convenience short-hand for me.Toggle("TexCube", n). For this change to be
applied, call FxEffect.UpdateRoutine() subsequently.

#### type LibElemIDChangedHandler

```go
type LibElemIDChangedHandler func(oldNewIDs map[int]int)
```


#### type LibElemIDChangedHandlers

```go
type LibElemIDChangedHandlers []LibElemIDChangedHandler
```


#### func (*LibElemIDChangedHandlers) Add

```go
func (me *LibElemIDChangedHandlers) Add(fn LibElemIDChangedHandler)
```

#### type Mesh

```go
type Mesh struct {
	ID             int
	DefaultModelID int
	Name           string
}
```


#### func (*Mesh) GpuDelete

```go
func (me *Mesh) GpuDelete()
```

#### func (*Mesh) GpuUpload

```go
func (me *Mesh) GpuUpload() (err error)
```

#### func (*Mesh) GpuUploaded

```go
func (me *Mesh) GpuUploaded() bool
```

#### func (*Mesh) Load

```go
func (me *Mesh) Load(provider u3d.MeshProvider) (err error)
```

#### func (*Mesh) Loaded

```go
func (me *Mesh) Loaded() bool
```

#### func (*Mesh) Unload

```go
func (me *Mesh) Unload()
```

#### type MeshBuffer

```go
type MeshBuffer struct {
	Name string
}
```


#### func (*MeshBuffer) Add

```go
func (me *MeshBuffer) Add(meshID int) (err error)
```

#### func (*MeshBuffer) Remove

```go
func (me *MeshBuffer) Remove(meshID int)
```

#### type MeshBufferLib

```go
type MeshBufferLib []*MeshBuffer
```

Only used for Core.Mesh.Buffers

#### func (*MeshBufferLib) AddNew

```go
func (me *MeshBufferLib) AddNew(name string, capacity int32) (buf *MeshBuffer, err error)
```

#### func (MeshBufferLib) FloatsPerVertex

```go
func (_ MeshBufferLib) FloatsPerVertex() int32
```

#### func (MeshBufferLib) IsOk

```go
func (me MeshBufferLib) IsOk(id int) bool
```

#### func (MeshBufferLib) MemSizePerIndex

```go
func (_ MeshBufferLib) MemSizePerIndex() int32
```

#### func (MeshBufferLib) MemSizePerVertex

```go
func (_ MeshBufferLib) MemSizePerVertex() int32
```

#### func (*MeshBufferLib) Remove

```go
func (me *MeshBufferLib) Remove(fromID, num int)
```

#### func (MeshBufferLib) Walk

```go
func (me MeshBufferLib) Walk(on func(ref *MeshBuffer))
```

#### type MeshLib

```go
type MeshLib []Mesh
```

Only used for Core.Libs.Meshes

#### func (*MeshLib) AddNew

```go
func (me *MeshLib) AddNew() (id int)
```

#### func (*MeshLib) AddNewAndLoad

```go
func (me *MeshLib) AddNewAndLoad(name string, meshProvider u3d.MeshProvider) (meshID int, err error)
```

#### func (*MeshLib) Compact

```go
func (me *MeshLib) Compact()
```

#### func (MeshLib) GpuSync

```go
func (_ MeshLib) GpuSync() (err error)
```

#### func (MeshLib) IsOk

```go
func (me MeshLib) IsOk(id int) (ok bool)
```

#### func (MeshLib) MeshCube

```go
func (_ MeshLib) MeshCube() u3d.MeshProvider
```

#### func (MeshLib) MeshPlane

```go
func (_ MeshLib) MeshPlane() u3d.MeshProvider
```

#### func (MeshLib) MeshPyramid

```go
func (_ MeshLib) MeshPyramid() u3d.MeshProvider
```

#### func (MeshLib) MeshQuad

```go
func (_ MeshLib) MeshQuad() u3d.MeshProvider
```

#### func (MeshLib) MeshTri

```go
func (_ MeshLib) MeshTri() u3d.MeshProvider
```

#### func (MeshLib) Ok

```go
func (me MeshLib) Ok(id int) bool
```

#### func (MeshLib) Remove

```go
func (me MeshLib) Remove(fromID, num int)
```

#### func (MeshLib) Walk

```go
func (me MeshLib) Walk(on func(ref *Mesh))
```

#### type Model

```go
type Model struct {
	ID    int
	MatID int
	Name  string
}
```

A Model is a parameterized instantiation of its parent Mesh geometry with unique
appearance, material or other properties.

#### type ModelLib

```go
type ModelLib []Model
```

Only used for Core.Libs.Models

#### func (*ModelLib) AddNew

```go
func (me *ModelLib) AddNew() (id int)
```

#### func (*ModelLib) Compact

```go
func (me *ModelLib) Compact()
```

#### func (ModelLib) IsOk

```go
func (me ModelLib) IsOk(id int) (ok bool)
```

#### func (ModelLib) Ok

```go
func (me ModelLib) Ok(id int) bool
```

#### func (ModelLib) Remove

```go
func (me ModelLib) Remove(fromID, num int)
```

#### func (ModelLib) Walk

```go
func (me ModelLib) Walk(on func(ref *Model))
```

#### type NgCore

```go
type NgCore struct {
	Libs NgLibs
	Mesh struct {
		Buffers MeshBufferLib
		Desc    struct {
			Cube, Plane, Pyramid, Quad, Tri u3d.MeshProvider
		}
	}
	Render struct {
		Canvases RenderCanvasLib
		Fx       struct {
			KnownProcIDs []string
			Samplers     struct {
				NoFilteringClamp    ugl.Sampler
				FullFilteringRepeat ugl.Sampler
				FullFilteringClamp  ugl.Sampler
			}
		}
		KnownTechniques map[string]renderTechniqueProvider
	}
}
```

NgCore is a singleton type, only used for the core.Core package-global exported
variable. It is only aware of that instance and does not support any other
NgCore instances.

```go
var (
	//	The heart and brain of go:ngine --- a container for all runtime resources and responsible for rendering.
	Core NgCore
)
```

#### func (*NgCore) GpuSyncImageLibs

```go
func (_ *NgCore) GpuSyncImageLibs() (err error)
```

#### type NgDiag

```go
type NgDiag struct {
	LogCategories          NgDiagLogCategory
	LogCategoryNames       map[NgDiagLogCategory]string
	LogGLErrorsInLoopOnSec bool
}
```

NgDiag is a singleton type, only used for the core.Diag package-global exported
variable. It is only aware of that instance and does not support any other
NgDiag instances.

#### func (*NgDiag) Log

```go
func (_ *NgDiag) Log(cat NgDiagLogCategory, fmt string, fmtArgs ...interface{})
```

#### func (*NgDiag) LogErr

```go
func (_ *NgDiag) LogErr(err error)
```

#### func (*NgDiag) LogIfGlErr

```go
func (_ *NgDiag) LogIfGlErr(fmt string, fmtArgs ...interface{})
```

#### func (*NgDiag) LogImages

```go
func (_ *NgDiag) LogImages(fmt string, fmtArgs ...interface{})
```

#### func (*NgDiag) LogMeshes

```go
func (_ *NgDiag) LogMeshes(fmt string, fmtArgs ...interface{})
```

#### func (*NgDiag) LogMisc

```go
func (_ *NgDiag) LogMisc(fmt string, fmtArgs ...interface{})
```

#### func (*NgDiag) LogShaders

```go
func (_ *NgDiag) LogShaders(fmt string, fmtArgs ...interface{})
```

#### type NgDiagLogCategory

```go
type NgDiagLogCategory int
```


```go
const (
	DiagLogCatMisc    NgDiagLogCategory = 1
	DiagLogCatMeshes  NgDiagLogCategory = 2
	DiagLogCatShaders NgDiagLogCategory = 4
	DiagLogCatImages  NgDiagLogCategory = 8
	DiagLogCatAll     NgDiagLogCategory = DiagLogCatMeshes | DiagLogCatMisc | DiagLogCatShaders | DiagLogCatImages
)
```

#### type NgLibs

```go
type NgLibs struct {
	Effects   FxEffectLib
	Materials FxMaterialLib
	Images    struct {
		SplashScreen FxImage2D
		TexCube      FxImageCubeLib
		Tex2D        FxImage2DLib
	}
	Meshes MeshLib
	Models ModelLib
	Scenes SceneLib
}
```

Only used for Core.Libs.

#### func (*NgLibs) UpdateIDRef

```go
func (_ *NgLibs) UpdateIDRef(oldNewIDs map[int]int, ptr *int)
```

#### func (*NgLibs) UpdateIDRefs

```go
func (_ *NgLibs) UpdateIDRefs(oldNewIDs map[int]int, ptrs ...*int)
```

#### func (*NgLibs) UpdateIDRefsIn

```go
func (_ *NgLibs) UpdateIDRefsIn(oldNewIDs map[int]int, slice []int)
```

#### func (*NgLibs) UpdatedIDRef

```go
func (_ *NgLibs) UpdatedIDRef(oldNewIDs map[int]int, in int) (out int)
```

#### type NgLoop

```go
type NgLoop struct {
	//	Set to true by Loop.Run(). Set to false to stop looping.
	Running bool

	Delay time.Duration

	MaxIterations float64

	On struct {
		//	While Loop.Run() is running, this callback is invoked (in its own "app thread")
		//	every loop iteration (ie. once per frame).
		//	This callback may run in parallel with On.EverySec(), but never with On.WinThread().
		AppThread func()

		//	While Loop.Run() is running, this callback is invoked (on the main windowing thread)
		//	every loop iteration (ie. once per frame).
		//	This callback is guaranteed to never run in parallel with
		//	(and always after) the On.AppThread() and On.EverySec() callbacks.
		WinThread func()

		//	While Loop.Run() is running, this callback is invoked (on the main windowing thread)
		//	at least and at most once per second, a useful entry point for non-real-time periodically recurring code.
		//	Caution: unlike On.WinThread(), this callback runs in parallel with your On.AppThread() callback.
		EverySec func()
	}

	Tick struct {
		//	The tick-time when the Loop.On.EverySec() callback was last invoked.
		PrevSec int

		//	While Loop.Run() is running, is set to the current "tick-time":
		//	the time in seconds expired ever since Loop.Run() was last called.
		Now float64

		//	While Loop.Run() is running, is set to the previous tick-time.
		Prev float64

		//	The delta between Tick.Prev and Tick.Now.
		Delta float64
	}
}
```

NgLoop is a singleton type, only used for the Loop variable. It is only aware of
that instance and does not support any other NgLoop instances.

```go
var (
	//	Manages your main-thread's render loop.
	//	Call it's Run() method once after go:ngine initialization (see examples).
	Loop NgLoop
)
```

#### func (*NgLoop) Run

```go
func (_ *NgLoop) Run()
```
Initiates a rendering loop. This method returns only when the loop is stopped
for whatever reason.

(Before entering the loop, this method performs a one-off GC invokation.)

#### func (*NgLoop) Time

```go
func (_ *NgLoop) Time() float64
```
Returns the number of seconds expired ever since Loop.Run() was last called.

#### type NgOptions

```go
type NgOptions struct {
	AppDir struct {
		//	The base directory path for app file paths.
		BasePath string

		Temp struct {
			BaseName       string
			ShaderSources  string
			CachedTextures string
		}
	}

	Cameras struct {
		DefaultControllerParams ControllerParams
		PerspectiveDefaults     u3d.Perspective
	}

	Initialization struct {
		GlContext struct {
			CoreProfile struct {
				//	Required on Mac OS X, not necessary elsewhere.
				//	While potentially slightly beneficial with recent GL drivers,
				//	might also fail with a select few rather outdated ones.
				//	Defaults to true on Mac OS X, to false elsewhere.
				ForceFirst bool

				//	While required on Mac OS X, really not recommended elsewhere
				//	(at present). Defaults to true on Mac OS X, to false elsewhere.
				//	Only used if go:ngine requests the creation of a GL core-profile context.
				ForwardCompat bool

				//	Defaults to the newest GL version currently supported by the GL
				//	binding used by go:ngine. The binding adaptively uses features
				//	of GL versions newer than 3.3 only if they are available, so this
				//	is a most strongly recommended default for release apps. But for
				//	testing, this is useful to test performance in older GL versions.
				//	Must be one of the values in glutil.KnownVersions.
				VersionHint float64
			}

			//	Defaults to the DefaultBadVersionMessage constant. If using a custom
			//	string, you can use the same placeholders as that one.
			BadVersionMessage string
		}
		DefaultCanvas struct {
			GammaViaShader bool
			SplashImage    []byte
		}
		Window struct {
			//	Defaults: Color.R=8 Color.G=8 Color.B=8 Color.A=0 Depth=8 Stencil=0.
			//	These defaults are reasonable when using a render-to-texture off-screen
			//	RenderCanvas. Otherwise, may want to bump Depth to at least 24 or 32.
			//	Depth shouldn't be 0 as this causes some Intel HD drivers to bug out badly.
			BufSizes ngctx.BufferBits
		}
	}

	Libs struct {
		InitialCap   int
		GrowCapBy    int
		OnIDsChanged struct {
			Effects LibElemIDChangedHandlers
			Images  struct {
				TexCube LibElemIDChangedHandlers
				Tex2D   LibElemIDChangedHandlers
			}
			Materials LibElemIDChangedHandlers
			Meshes    LibElemIDChangedHandlers
			Models    LibElemIDChangedHandlers
			Scenes    LibElemIDChangedHandlers
		}
	}

	Loop struct {
		//	Controls whether and how often the Garbage Collector
		//	is invoked during the Loop.
		GcEvery struct {
			//	Defaults to false. If true, GC will
			//	be invoked every frame during the Loop.
			Frame bool

			//	Defaults to true. If true, GC will be invoked at
			//	least and at most once per second during the Loop.
			Sec bool
		}
	}

	Rendering struct {
		DefaultBatcher    RenderBatcher
		DefaultClearColor ugl.GlVec4
	}

	Textures struct {
		Storage FxImageStorage
	}
}
```

Only used for the Options variable.

```go
var (
	Options NgOptions
)
```

#### type NgStats

```go
type NgStats struct {
	//	Gives the total number of frames rendered during the "previous"
	//	(not the current) second. Good enough for just a simple-minded FPS indicator.
	FpsLastSec int

	//	This TimingStats instance combines all the individual FrameFoo fields
	//	to track over time (both average and maximum) total cost per frame.
	Frame TimingStats

	//	"Rendering" consists of a CPU-side and a GPU-side cost.
	//	This TimingStats instance combines both to track over time
	//	(both average and maximum) total rendering cost per frame.
	FrameRenderBoth TimingStats

	//	The CPU-side cost of rendering comprises sending pre-batched
	//	rendering commands (prepared by the "prep" stage) to the GPU.
	//	This TimingStats instance tracks over time (both average
	//	and maximum) CPU-side rendering cost per frame.
	FrameRenderCpu TimingStats

	//	The GPU-side cost of rendering comprises execution of all draw calls
	//	sent by the CPU-side, plus waiting for V-sync if enabled.
	//	This TimingStats instance tracks over time (both average
	//	and maximum) GPU-side rendering cost per frame.
	FrameRenderGpu TimingStats

	//	"Prep code" comprises all go:ngine logic executed every frame in parallel to cull
	//	geometry and prepare a batch of rendering commands for the next (not current) frame.
	//	This TimingStats instance tracks over time (both average and maximum) "prep code" cost per frame.
	FramePrepThread TimingStats

	//	"App code" comprises (mostly user-specific) logic executed every frame in parallel in
	//	your Loop.OnAppThread() callback. Such code may freely modify dynamic Cameras, Nodes etc.
	//	Unlike OnWinThread() code, "app code" always runs in its own thread in parallel to the prep and main threads.
	//	This TimingStats instance tracks over time (both average and maximum) "app code" cost per frame.
	FrameAppThread TimingStats

	//	Tracks over the time (both average and maximum) cost for Loop.On.EverySec() and,
	//	if Diag.LogGLErrorsInLoopOnSec is true, Diag.LogIfGlErr().
	FrameOnSec TimingStats

	//	"Windowing/GPU/IO code" comprises user-specific logic executed every frame via your own
	//	Loop.OnWinThread() callback. This should be kept to a minimum to fully enjoy
	//	the benefits of multi-threading. Main use-cases are calls resulting in GPU state
	//	changes (such as toggling effects in Core.Render.PostFx) and working with UserIO
	//	to poll for user input -- but do consider executing resulting logic in your OnAppThread().
	//	This TimingStats instance tracks over time (both average and maximum) "input code" cost per frame.
	FrameWinThread TimingStats

	//	When CPU-side rendering is completed, Loop waits for the app thread and prep thread
	//	to finish (either before or after GPU-side rendering depending on Loop.SwapLast).
	//	It then moves "prep results" to the render thread and "app results" to the prep thread.
	//	This TimingStats instance tracks over time (both average and maximum) "thread sync" cost per frame.
	FrameThreadSync TimingStats

	//	During the Loop, the Go Garbge Collector is invoked at least and at most once per second.
	//
	//	Forcing GC "that often" practically guarantees it will almost never have so much work to do as to
	//	noticably block user interaction --- typically well below 10ms, most often around 1ms.
	//
	//	This TimingStats instance over time tracks the maximum and average time spent on that
	//	1x-per-second-during-Loop GC invokation (but does not track any other GC invokations).
	Gc TimingStats

	Programs struct {
		NumProgsCompiled int
		TotalTimeCost    int64
	}
}
```

Consider NgStats a "Singleton" type, only valid use is the core.Stats global
variable. Tracks various go:ngine performance indicators over time.

```go
var (
	//	Tracks various go:ngine performance counters over time.
	Stats NgStats
)
```

#### func (*NgStats) AverageFps

```go
func (_ *NgStats) AverageFps() float64
```
Returns the average number of frames-per-second since Loop.Loop() was last
called.

#### func (*NgStats) TotalFrames

```go
func (_ *NgStats) TotalFrames() float64
```

#### type NgUserIO

```go
type NgUserIO struct {
	//	Minimum delay for NgUserIO.KeyToggled() method, in seconds. Defaults to 0.15.
	KeyToggleMinDelay float64

	Window WindowOptions
}
```

Consider NgUserIO a "Singleton" type, only valid use is the core.UserIO global
variable. Your gateway for end-user input (key, mouse etc.) and "output" (window
management, not the graphics themselves).

```go
var (
	//	Your gateway for end-user input (key, mouse etc.) and "output" (window management, not the graphics themselves).
	UserIO NgUserIO
)
```

#### func (*NgUserIO) IifKeyF

```go
func (_ *NgUserIO) IifKeyF(key int, ifTrue, ifFalse float64) float64
```
Returns ifTrue if the specified key is pressed, otherwise returns ifFalse.

#### func (*NgUserIO) KeyPressed

```go
func (_ *NgUserIO) KeyPressed(key int) bool
```
Returns true if the specified key is pressed.

#### func (*NgUserIO) KeyPressedWhich

```go
func (_ *NgUserIO) KeyPressedWhich(keys ...int) int
```
Returns the first in keys that is pressed.

#### func (*NgUserIO) KeyToggled

```go
func (_ *NgUserIO) KeyToggled(key int) bool
```
Returns true if the specified key has been "toggled", ie. its pressed-state
changed within the last me.KeyToggleMinDelay seconds.

#### func (*NgUserIO) KeysPressedAll2

```go
func (_ *NgUserIO) KeysPressedAll2(k1, k2 int) bool
```
Returns true if both specified keys are pressed.

#### func (*NgUserIO) KeysPressedAll3

```go
func (_ *NgUserIO) KeysPressedAll3(k1, k2, k3 int) bool
```
Returns true if all three specified keys are pressed.

#### func (*NgUserIO) KeysPressedAny2

```go
func (_ *NgUserIO) KeysPressedAny2(k1, k2 int) bool
```
Returns true if any of the two specified keys is pressed.

#### func (*NgUserIO) KeysPressedAny3

```go
func (_ *NgUserIO) KeysPressedAny3(k1, k2, k3 int) bool
```
Returns true if any of the three specified keys is pressed.

#### type RenderBatchCriteria

```go
type RenderBatchCriteria int
```


```go
const (
	BatchByProgram RenderBatchCriteria = 0
	BatchByTexture RenderBatchCriteria = 1
	BatchByBuffer  RenderBatchCriteria = 2
)
```

#### type RenderBatcher

```go
type RenderBatcher struct {
	Enabled  bool
	Priority [numPrios]RenderBatchCriteria
}
```


#### type RenderCanvas

```go
type RenderCanvas struct {
	//	This should be an non-negative integer, it's a float64 just to avoid a
	//	type conversion. How often this RenderCanvas is included in rendering:
	//	1 = every frame (this is the default value)
	//	2 = every 2nd frame
	//	3, 5, 8 etc. = every 3rd, 5th, 8th etc. frame
	//	0 = this RenderCanvas is disabled for rendering
	EveryNthFrame float64

	Views RenderViewLib

	Srgb bool
}
```

Represents a surface (texture framebuffer) that can be rendered to.

#### func (*RenderCanvas) AddNewView

```go
func (me *RenderCanvas) AddNewView(renderTechnique string) (view *RenderView)
```

#### func (*RenderCanvas) CurrentAbsoluteSize

```go
func (me *RenderCanvas) CurrentAbsoluteSize() (width, height int)
```

#### func (*RenderCanvas) SetSize

```go
func (me *RenderCanvas) SetSize(relative bool, width, height float64)
```
Sets the 2 dimensions of this render canvas. If relative is true, width and
height are interpreted relative to the resolution of the OpenGL context's
default framebuffer, with 1 being 100%. Otherwise, width and height are absolute
pixel dimensions.

#### type RenderCanvasLib

```go
type RenderCanvasLib []*RenderCanvas
```

Only used for Core.Render.Canvases

#### func (*RenderCanvasLib) AddNew

```go
func (me *RenderCanvasLib) AddNew(relative bool, width, height float64) (canv *RenderCanvas)
```

#### func (RenderCanvasLib) IsOk

```go
func (me RenderCanvasLib) IsOk(id int) bool
```

#### func (*RenderCanvasLib) Remove

```go
func (me *RenderCanvasLib) Remove(fromID, num int)
```

#### func (RenderCanvasLib) Walk

```go
func (me RenderCanvasLib) Walk(on func(ref *RenderCanvas))
```

#### type RenderTechnique

```go
type RenderTechnique interface {
	// contains filtered or unexported methods
}
```


#### type RenderTechniqueQuad

```go
type RenderTechniqueQuad struct {
	Effect FxEffect
}
```


#### type RenderTechniqueScene

```go
type RenderTechniqueScene struct {
	Batch  RenderBatcher
	Camera Camera
}
```


#### func (*RenderTechniqueScene) ApplyCameraPerspective

```go
func (me *RenderTechniqueScene) ApplyCameraPerspective()
```
Applies changes made to the Enabled, FovY, ZNear and/or ZFar parameters in
me.Camera.Perspective.

#### func (*RenderTechniqueScene) NumDrawCalls

```go
func (me *RenderTechniqueScene) NumDrawCalls() int
```

#### func (*RenderTechniqueScene) ToggleBatching

```go
func (me *RenderTechniqueScene) ToggleBatching()
```

#### type RenderView

```go
type RenderView struct {
	Enabled bool

	FxProcs FxProcs

	//	The device-relative or absolute view-port for this view.
	Port RenderViewport

	RenderStates ugl.RenderStatesBag

	Technique RenderTechnique
}
```


#### func (*RenderView) Technique_Quad

```go
func (me *RenderView) Technique_Quad() (tech *RenderTechniqueQuad)
```

#### func (*RenderView) Technique_Scene

```go
func (me *RenderView) Technique_Scene() (tech *RenderTechniqueScene)
```

#### type RenderViewLib

```go
type RenderViewLib []*RenderView
```

Only used for Core.Render.Canvases[id].Views

#### func (RenderViewLib) IsOk

```go
func (me RenderViewLib) IsOk(id int) bool
```

#### func (*RenderViewLib) Remove

```go
func (me *RenderViewLib) Remove(fromID, num int)
```

#### func (RenderViewLib) Walk

```go
func (me RenderViewLib) Walk(on func(ref *RenderView))
```

#### type RenderViewport

```go
type RenderViewport struct {
}
```

Only used for RenderView.Port

#### func (*RenderViewport) SetAbsolute

```go
func (me *RenderViewport) SetAbsolute(x, y, width, height int)
```
Sets the absolute viewport origin and dimensions in pixels.

#### func (*RenderViewport) SetRelative

```go
func (me *RenderViewport) SetRelative(x, y, width, height float64)
```
Sets the RenderCanvas-relative viewport origin and dimensions, with the value
1.0 representing the maximum extent of the viewport on that respective axis.

#### type Scene

```go
type Scene struct {
	ID int
}
```

Represents a scene graph.

#### func (*Scene) AddNewChildNode

```go
func (me *Scene) AddNewChildNode(parentNodeID, meshID int) (childNodeID int)
```

#### func (*Scene) ApplyNodeTransforms

```go
func (me *Scene) ApplyNodeTransforms(nodeID int)
```
Updates the internal 4x4 transformation matrix for all transformations of the
specified node and child-nodes. It is only this matrix that is used by the
rendering runtime.

#### func (*Scene) Node

```go
func (me *Scene) Node(id int) *SceneNode
```

#### func (*Scene) NumNodes

```go
func (me *Scene) NumNodes() int
```

#### func (*Scene) ParentNodeID

```go
func (me *Scene) ParentNodeID(childNodeID int) (parentID int)
```

#### func (*Scene) RemoveNode

```go
func (me *Scene) RemoveNode(fromID int)
```

#### func (*Scene) Root

```go
func (me *Scene) Root() *SceneNode
```

#### func (*Scene) SetNodeMeshID

```go
func (me *Scene) SetNodeMeshID(nodeID, meshID int)
```

#### type SceneLib

```go
type SceneLib []Scene
```

Only used for Core.Libs.Scenes

#### func (*SceneLib) AddNew

```go
func (me *SceneLib) AddNew() (id int)
```

#### func (*SceneLib) Compact

```go
func (me *SceneLib) Compact()
```

#### func (SceneLib) IsOk

```go
func (me SceneLib) IsOk(id int) (ok bool)
```

#### func (SceneLib) Ok

```go
func (me SceneLib) Ok(id int) bool
```

#### func (SceneLib) Remove

```go
func (me SceneLib) Remove(fromID, num int)
```

#### func (SceneLib) Walk

```go
func (me SceneLib) Walk(on func(ref *Scene))
```

#### type SceneNode

```go
type SceneNode struct {
	ID        int
	Transform SceneNodeTransform

	Render struct {
		Cull struct {
			Frustum bool
		}
		Enabled bool
		MatID   int
		ModelID int
	}
}
```


#### type SceneNodeLib

```go
type SceneNodeLib []SceneNode
```

Only used for Core.Scenes[id].allNodes

#### func (*SceneNodeLib) AddNew

```go
func (me *SceneNodeLib) AddNew() (id int)
```

#### func (*SceneNodeLib) Compact

```go
func (me *SceneNodeLib) Compact()
```

#### func (SceneNodeLib) IsOk

```go
func (me SceneNodeLib) IsOk(id int) (ok bool)
```

#### func (SceneNodeLib) Ok

```go
func (me SceneNodeLib) Ok(id int) bool
```

#### func (SceneNodeLib) Remove

```go
func (me SceneNodeLib) Remove(fromID, num int)
```

#### func (SceneNodeLib) Walk

```go
func (me SceneNodeLib) Walk(on func(ref *SceneNode))
```

#### type SceneNodeTransform

```go
type SceneNodeTransform struct {
	//	Translation of the from origin.
	Pos unum.Vec3

	//	Rotation for each axis in radians.
	Rot unum.Vec3

	//	Scaling of this node, if any. Defaults to (1, 1, 1) for no scaling.
	Scale unum.Vec3
}
```

Represents one or more transformations of a Node. This is only used by Node
objects, which initialize their SceneNodeTransform with the proper defaults and
associate themselves with their SceneNodeTransform. (Any other
SceneNodeTransform are invalid.)

A single SceneNodeTransform encapsulates an unexported 4x4 matrix that is
recalculated whenever its parent Node.ApplyTransform() method is called.

#### func (*SceneNodeTransform) AddRot

```go
func (me *SceneNodeTransform) AddRot(rot *unum.Vec3)
```

#### func (*SceneNodeTransform) AddRotXYZ

```go
func (me *SceneNodeTransform) AddRotXYZ(x, y, z float64)
```

#### func (*SceneNodeTransform) SetPos

```go
func (me *SceneNodeTransform) SetPos(posX, posY, posZ float64)
```

#### func (*SceneNodeTransform) SetRot

```go
func (me *SceneNodeTransform) SetRot(radX, radY, radZ float64)
```

#### func (*SceneNodeTransform) SetScale

```go
func (me *SceneNodeTransform) SetScale(s float64)
```

#### func (*SceneNodeTransform) SetScaleXyz

```go
func (me *SceneNodeTransform) SetScaleXyz(x, y, z float64)
```

#### func (*SceneNodeTransform) StepDelta

```go
func (me *SceneNodeTransform) StepDelta(deltaPerSecond float64) float64
```
Returns the result of multiplying deltaPerSecond with EngineLoop.TickDelta.

#### type TimingStats

```go
type TimingStats struct {
}
```

Helps track average and maximum cost for a variety of performance indicators.

#### func (*TimingStats) Average

```go
func (me *TimingStats) Average() float64
```
Returns the average cost tracked by this performance indicator.

#### func (*TimingStats) Max

```go
func (me *TimingStats) Max() float64
```
Returns the maximum cost tracked by this performance indicator.

#### type WindowOptions

```go
type WindowOptions struct {
	//	Defaults to a function that returns true to allow closing the window.
	OnCloseRequested func() bool

	//	Minimum delay, in seconds, to wait after the last window-resize event received from
	//	the OS before notifying the rendering runtime of the new window dimensions.
	//	Defaults to 0.15.
	ResizeMinDelay float64

	//	Number of samples (0, 2, 4, 8...)
	MultiSampling int
}
```


#### func (*WindowOptions) Created

```go
func (me *WindowOptions) Created() bool
```

#### func (*WindowOptions) Fullscreen

```go
func (me *WindowOptions) Fullscreen() bool
```

#### func (*WindowOptions) Height

```go
func (me *WindowOptions) Height() int
```
Returns the height of the window in pixels.

#### func (*WindowOptions) SetSize

```go
func (me *WindowOptions) SetSize(width, height int)
```

#### func (*WindowOptions) SetSwapInterval

```go
func (me *WindowOptions) SetSwapInterval(newSwap int)
```

#### func (*WindowOptions) SetTitle

```go
func (me *WindowOptions) SetTitle(newTitle string)
```
Sets the window title to newTitle.

#### func (*WindowOptions) SwapInterval

```go
func (me *WindowOptions) SwapInterval() int
```

#### func (*WindowOptions) Title

```go
func (me *WindowOptions) Title() string
```

#### func (*WindowOptions) Width

```go
func (me *WindowOptions) Width() int
```
Returns the width of the window in pixels.

--
**godocdown** http://github.com/robertkrimen/godocdown
