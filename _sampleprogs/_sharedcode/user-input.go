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

	Keys struct {
		Pressed map[int]bool

		CheckFor struct {
			Pressed []int
		}
		tmp int
	}
)

func init() {
	Keys.Pressed = map[int]bool{}
	Keys.CheckFor.Pressed = []int{
		glfw.KeyLalt, glfw.KeyLshift, glfw.KeyRshift, 'W', 'A', 'S', 'D',
		glfw.KeyLeft, glfw.KeyRight, glfw.KeyUp, glfw.KeyDown,
		glfw.KeyPagedown, glfw.KeyPageup, glfw.KeyKP9, glfw.KeyKP3,
	}
}

//	Called every frame (by the parent example app) to check the state for keys controlling CamCtl to move or rotate Cam.
func CheckCamCtlKeys() {
	if !Paused {
		for _, Keys.tmp = range Keys.CheckFor.Pressed {
			Keys.Pressed[Keys.tmp] = ng.UserIO.KeyPressed(Keys.tmp)
		}
	}
}

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
