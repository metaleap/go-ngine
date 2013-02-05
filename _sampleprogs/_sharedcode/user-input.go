package samplescenes

import (
	glfw "github.com/go-gl/glfw"
	ng "github.com/go3d/go-ngine/core"
)

var (
	//	The maximum index for KeyHints when cycling through it in OnSec()
	MaxKeyHint = len(KeyHints) - 1

	//	OnSec() changes the window title every second to display FPS etc.
	//	Also every 4 seconds shows the next one in a number of "key hints" defined here:
	KeyHints = []string{
		"[F2]  --  Toggle Backface Culling",
		"[F3]  --  Pause/Resume",
		"[F4]  --  Retro Mode",
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
		glfw.KeyLalt, glfw.KeyLshift, glfw.KeyRshift, 'W', 'A', 'S', 'D',
		glfw.KeyLeft, glfw.KeyRight, glfw.KeyUp, glfw.KeyDown,
		glfw.KeyPagedown, glfw.KeyPageup, glfw.KeyKP9, glfw.KeyKP3,
	}
}

//	Key hints are shown in the window title bar, a new one every few seconds.
//	This function adds a new hint to be included in that rotation.
func AddKeyHint(key, hint string) {
	KeyHints = append(KeyHints, "["+key+"]  --  "+hint)
	MaxKeyHint = len(KeyHints) - 1
}

//	To be called every frame (by the parent example app, ONLY in onWinThread()!) to collect key-press
//	states for controlling (in HandleCamCtlKeys() in onAppThread()!) CamCtl to move or rotate Cam.
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
	if in.KeyToggled(glfw.KeyEsc) {
		ng.Loop.Stop()
	}
	if in.KeyToggled(glfw.KeyF2) {
		Cam.Rendering.States.FaceCulling = !Cam.Rendering.States.FaceCulling
	}
	if in.KeyToggled(glfw.KeyF3) {
		PauseResume()
	}
	if in.KeyToggled(glfw.KeyF4) {
		ToggleRetro()
	}
}

//	To be called every frame (by the parent example app, ONLY in onAppThread()!) to process key-press
//	states (from the previous CheckCamCtlKeys() call) for controlling CamCtl to move or rotate Cam.
func HandleCamCtlKeys() {
	if CamCtl.Params.MoveSpeedupFactor = 1; !Paused {
		if Keys.Pressed[glfw.KeyLshift] {
			CamCtl.Params.MoveSpeedupFactor = 10
		} else if Keys.Pressed[glfw.KeyRshift] {
			CamCtl.Params.MoveSpeedupFactor = 100
		} else if Keys.Pressed[glfw.KeyLalt] {
			CamCtl.Params.MoveSpeedupFactor = 0.1
		}
		if Keys.Pressed[glfw.KeyUp] {
			CamCtl.MoveForward()
		}
		if Keys.Pressed[glfw.KeyDown] {
			CamCtl.MoveBackward()
		}
		if Keys.Pressed['A'] {
			CamCtl.MoveLeft()
		}
		if Keys.Pressed['D'] {
			CamCtl.MoveRight()
		}
		if Keys.Pressed['W'] {
			CamCtl.MoveUp()
		}
		if Keys.Pressed['S'] {
			CamCtl.MoveDown()
		}
		if Keys.Pressed[glfw.KeyLeft] {
			CamCtl.TurnLeft()
		}
		if Keys.Pressed[glfw.KeyRight] {
			CamCtl.TurnRight()
		}
		if Keys.Pressed[glfw.KeyPageup] || Keys.Pressed[glfw.KeyKP9] {
			CamCtl.TurnUp()
		}
		if Keys.Pressed[glfw.KeyPagedown] || Keys.Pressed[glfw.KeyKP3] {
			CamCtl.TurnDown()
		}
	}
}
