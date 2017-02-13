package exampleutils

import (
	ng "github.com/metaleap/go-ngine/___old2013/core"
)

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

		tmp int
	}
)

func init() {
	Keys.Pressed = map[int]bool{}
	Keys.CheckForPressed = []int{
		KeyLalt, KeyLshift, KeyRshift, 'W', 'A', 'S', 'D',
		KeyLeft, KeyRight, KeyUp, KeyDown,
		KeyPagedown, KeyPageup, KeyKP9, KeyKP3,
	}
}

//	Key hints are shown in the window title bar, a new one every few seconds.
//	This function adds a new hint to be included in that rotation.
func AddKeyHint(key, hint string) {
	KeyHints = append(KeyHints, "["+key+"]  --  "+hint)
	MaxKeyHint = len(KeyHints) - 1
}

//	To be called every frame (by the parent example app, ONLY in onWinThread()!) to collect key-press
//	states for controlling (in HandleCamCtlKeys() in onAppThread()!) SceneCam.Controller to move or rotate SceneCam.
func CheckCamCtlKeys() {
	if !Paused {
		for _, Keys.tmp = range Keys.CheckForPressed {
			Keys.Pressed[Keys.tmp] = ng.UserIO.KeyPressed(Keys.tmp)
		}
	}
}

//	To be called every frame (by the parent example app, ONLY in onWinThread()!) to
//	check AND handle key-toggle states for F2, F3 etc. function keys and Esc.
func CheckAndHandleToggleKeys() {
	in := &ng.UserIO
	if in.KeyToggled(KeyEsc) {
		PauseResume()
	}
	if SceneView != nil && in.KeyToggled(KeyF2) {
		SceneView.RenderStates.FaceCulling = !SceneView.RenderStates.FaceCulling
	}
	if in.KeyToggled(KeyF3) {
		toggleRetro()
	}
	if in.KeyToggled(KeyF4) {
		toggleTexturing()
	}
	if in.KeyToggled(KeyF5) {
		toggleBatching()
	}
	if in.KeysPressedAll2(KeyLctrl, 'Q') {
		ng.Loop.Running = false
	}
}

//	To be called every frame (by the parent example app, ONLY in onAppThread()!) to process key-press
//	states (from the previous CheckCamCtlKeys() call) for controlling SceneCam.Controller to move or rotate SceneCam.
func HandleCamCtlKeys() {
	if SceneCam.Controller.Params.MoveSpeedupFactor = 1; !Paused {
		camCtl := &SceneCam.Controller
		if Keys.Pressed[KeyLshift] {
			camCtl.Params.MoveSpeedupFactor = 10
		} else if Keys.Pressed[KeyRshift] {
			camCtl.Params.MoveSpeedupFactor = 100
		} else if Keys.Pressed[KeyLalt] {
			camCtl.Params.MoveSpeedupFactor = 0.1
		}
		if Keys.Pressed[KeyUp] {
			camCtl.MoveForward()
		}
		if Keys.Pressed[KeyDown] {
			camCtl.MoveBackward()
		}
		if Keys.Pressed['A'] {
			camCtl.MoveLeft()
		}
		if Keys.Pressed['D'] {
			camCtl.MoveRight()
		}
		if Keys.Pressed['W'] {
			camCtl.MoveUp()
		}
		if Keys.Pressed['S'] {
			camCtl.MoveDown()
		}
		if Keys.Pressed[KeyLeft] {
			camCtl.TurnLeft()
		}
		if Keys.Pressed[KeyRight] {
			camCtl.TurnRight()
		}
		if Keys.Pressed[KeyPageup] || Keys.Pressed[KeyKP9] {
			camCtl.TurnUp()
		}
		if Keys.Pressed[KeyPagedown] || Keys.Pressed[KeyKP3] {
			camCtl.TurnDown()
		}
	}
}
