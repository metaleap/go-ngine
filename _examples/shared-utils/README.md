# exampleutils
--
    import "github.com/go3d/go-ngine/_examples/shared-utils"

Some common base functionality shared by all the go:ngine example apps in
_examples.

For full-screen instead of windowed, change `WinFullscreen`.

To switch GL context providers (GLFW3 vs. GLFW2 vs. SDL), check out the
`use-*.g*` files.

## Usage

```go
const (
	KeyLalt     int = int(glfw.KeyLeftAlt)
	KeyLshift   int = int(glfw.KeyLeftShift)
	KeyRshift   int = int(glfw.KeyRightShift)
	KeyLeft     int = int(glfw.KeyLeft)
	KeyRight    int = int(glfw.KeyRight)
	KeyUp       int = int(glfw.KeyUp)
	KeyDown     int = int(glfw.KeyDown)
	KeyPagedown int = int(glfw.KeyPageDown)
	KeyPageup   int = int(glfw.KeyPageUp)
	KeyKP9      int = int(glfw.KeyKp9)
	KeyKP3      int = int(glfw.KeyKp3)
	KeyEsc      int = int(glfw.KeyEscape)
	KeyF1       int = int(glfw.KeyF1)
	KeyF2       int = int(glfw.KeyF2)
	KeyF3       int = int(glfw.KeyF3)
	KeyF4       int = int(glfw.KeyF4)
	KeyF5       int = int(glfw.KeyF5)
	KeyF6       int = int(glfw.KeyF6)
	KeyF7       int = int(glfw.KeyF7)
	KeyF8       int = int(glfw.KeyF8)
	KeyF9       int = int(glfw.KeyF9)
	KeyF10      int = int(glfw.KeyF10)
	KeyF11      int = int(glfw.KeyF11)
	KeyF12      int = int(glfw.KeyF12)
	KeyLctrl    int = int(glfw.KeyLeftControl)
)
```

```go
var (
	//	Change this before calling `Main` for full-screen mode.
	WinFullscreen = false

	//	Change to higher value to check out your splash-screen
	ArtificialSplashScreenDelay = 0 * time.Second

	//	Optionally set this to a callback function to be
	//	invoked every second on the windowing main thread.
	OnSec = func() {}

	//	The RenderCanvas the example scene is initially being rendered to.
	//	This is an off-screen "render-to-texture" RenderCanvas.
	SceneCanvas *ng.RenderCanvas

	SceneCam *ng.Camera

	SceneView *ng.RenderView

	//	Unlike the off-screen (render-to-texture) SceneCanvas above,
	//	this RenderCanvas epresents the actual screen/window.
	PostFxCanvas *ng.RenderCanvas

	//	Takes the image rendered to SceneCanvas, may
	//	post-process it or not, and blits it to PostFxCanvas.
	PostFxView *ng.RenderView

	Paused bool
)
```

```go
var (
	//	The maximum index for KeyHints when cycling through it in OnSec()
	MaxKeyHint = len(KeyHints) - 1

	//	OnSec() changes the window title every second to display FPS etc.
	//	Also every 4 seconds shows the next one in a number of "key hints" defined here:
	KeyHints = []string{
		"[Ctrl+Q]  --  Quit",
		"[Esc]  --  Pause/Resume",
		"[F2]  --  Toggle Backface Culling",
		"[F3]  --  Toggle Retro Mode",
		"[F4]  --  Toggle Texturing",
		"[F5]  --  Toggle Batching",
		"[W][S]  --  Camera rise / sink",
		"[A][D]  --  Camera strafe left / right",
		"[<][>]  --  Camera turn left / right",
		"[^][v]  --  Camera move forward / backward",
		"[PgUp][PgDn]  --  Camera turn up / down",
		"[Alt][LShift][RShift]  --  Camera move-speed x 0.1 / 10 / 100",
	}

	//	Because checking for key-presses happens in onWinThread() but handling such key presses
	//	(typically resulting in matrix remultiplications and other CPU work) happens in onAppThread,
	//	this struct helps keep track of both which keys to check in the first place, and their
	//	latest press-state.
	Keys struct {
		//	Updated by CheckCamCtlKeys(), contains the latest press-states of the keys in CheckForPressed.
		Pressed map[int]bool

		//	Contains the keys that CheckCamCtlKeys() will poll and update in Pressed.
		CheckForPressed []int
	}
)
```

```go
var LibIDs struct {
	Fx, Img2D, ImgCube, Mat map[string]int
}
```

#### func  AddKeyHint

```go
func AddKeyHint(key, hint string)
```
Key hints are shown in the window title bar, a new one every few seconds. This
function adds a new hint to be included in that rotation.

#### func  AddMainScene

```go
func AddMainScene() (me *ng.Scene)
```

#### func  AddNode

```go
func AddNode(scene *ng.Scene, parentNodeID, meshID, matID, modelID int) (node *ng.SceneNode)
```

#### func  AddSkyMesh

```go
func AddSkyMesh(scene *ng.Scene, meshID int)
```

#### func  AddTextureMaterials

```go
func AddTextureMaterials(idsUrls map[string]string)
```
Sets up textures and associated effects/materials with the specified IDs and
image URLs. For each ID (such as "foo" and "bar"): - creates an ng.FxImage2D
with ID "img_{ID}" (ie. "img_foo" and "img_bar") and adds it to
ng.Core.Libs.Images.I2D - creates an ng.FxEffect with ID "fx_{ID}" (ie. "fx_foo"
and "fx_bar") and adds it to ng.Core.Libs.Effects; its Diffuse field pointing to
the ng.FxImage2D - creates an ng.FxMaterial with ID "mat_{ID}" (ie. "mat_foo"
and "mat_bar") and adds it to ng.Core.Libs.Materials; its DefaultEffectID
pointing to the ng.FxEffect.

#### func  AppDirBasePath

```go
func AppDirBasePath() string
```
Returns the base path of the "app dir" for our example apps, in this case:
$GOPATH/src/github.com/go3d/go-ngine/_examples/-app-data

#### func  CheckAndHandleToggleKeys

```go
func CheckAndHandleToggleKeys()
```
To be called every frame (by the parent example app, ONLY in onWinThread()!) to
check AND handle key-toggle states for F2, F3 etc. function keys and Esc.

#### func  CheckCamCtlKeys

```go
func CheckCamCtlKeys()
```
To be called every frame (by the parent example app, ONLY in onWinThread()!) to
collect key-press states for controlling (in HandleCamCtlKeys() in
onAppThread()!) SceneCam.Controller to move or rotate SceneCam.

#### func  HandleCamCtlKeys

```go
func HandleCamCtlKeys()
```
To be called every frame (by the parent example app, ONLY in onAppThread()!) to
process key-press states (from the previous CheckCamCtlKeys() call) for
controlling SceneCam.Controller to move or rotate SceneCam.

#### func  Main

```go
func Main(setupExampleScene, onAppThread, onWinThread func())
```
Called by each example-app's func main(). Initializes go:ngine, calls the
specified setupExampleScene function, then enters The Loop.

#### func  PauseResume

```go
func PauseResume()
```
Pauses rendering or resumes from the current pause. When paused, the frame last
rendered is frozen and rendered in a gray-scale effect.

#### func  PrintPostLoopSummary

```go
func PrintPostLoopSummary()
```
Prints a summary of go:ngine's *Stats* performance counters when the parent
example app exits.

#### type Gui2D

```go
type Gui2D struct {
	View                 *ng.RenderView
	Cam                  *ng.Camera
	CatNodeID, DogNodeID int
}
```

A fake "2D GUI" concoction. There will be much better ngine-provided support for
stuff like this. Has two textured quads, a cat and a dog one, shows them both
animated and overlapping inside a red 64x64 px square in the bottom left canvas
corner.

#### func (*Gui2D) Setup

```go
func (me *Gui2D) Setup() (err error)
```

#### type RearMirror

```go
type RearMirror struct {
	Cam  *ng.Camera
	View *ng.RenderView
}
```

A rather simple "rear-view mirror" view that can be added to the
example-program's main render canvas.

```go
var RearView RearMirror
```
Not actively used unless RearView.Setup() is called.

#### func (*RearMirror) OnApp

```go
func (me *RearMirror) OnApp()
```
Copies the main camera's current position and direction and reverses its by 180
degrees to achieve a rear-view mirror.

#### func (*RearMirror) OnWin

```go
func (me *RearMirror) OnWin()
```
Syncs the rear-view camera's render states with the main camera's.

#### func (*RearMirror) Setup

```go
func (me *RearMirror) Setup(sceneID int)
```
Adds the rear-view mirror's camera to the main render canvas, at 1/3rd of its
width and 1/4th of its height.

#### func (*RearMirror) Toggle

```go
func (me *RearMirror) Toggle()
```
Enables or disables this rear-view mirror.

--
**godocdown** http://github.com/robertkrimen/godocdown
