package samplescenes

import (
	glfw "github.com/go-gl/glfw"
	ng "github.com/go3d/go-ngine/core"
)

var (
	Keys struct {
		Pressed map[int]bool
		Toggled map[int]bool

		CheckFor struct {
			Pressed []int
			Toggled []int
		}
		tmp int
	}
)

func init() {
	Keys.Pressed = map[int]bool{}
	Keys.Toggled = map[int]bool{}
	Keys.CheckFor.Toggled = []int{glfw.KeyF2, glfw.KeyF3, glfw.KeyEsc}
	Keys.CheckFor.Pressed = []int{
		glfw.KeyLalt, glfw.KeyLshift, glfw.KeyRshift, 'W', 'A', 'S', 'D',
		glfw.KeyLeft, glfw.KeyRight, glfw.KeyUp, glfw.KeyDown,
		glfw.KeyPagedown, glfw.KeyPageup, glfw.KeyKP9, glfw.KeyKP3,
	}
}

//	Called every frame (by the parent example app) to check the state for keys controlling CamCtl to move or rotate Cam.
func CheckCamCtlKeys() {
	for _, Keys.tmp = range Keys.CheckFor.Pressed {
		Keys.Pressed[Keys.tmp] = ng.UserIO.KeyPressed(Keys.tmp)
	}
}

//	Called every frame (by the parent example app) to check "toggle keys" that toggle certain options.
func CheckToggleKeys() {
	for _, Keys.tmp = range Keys.CheckFor.Toggled {
		Keys.Toggled[Keys.tmp] = ng.UserIO.KeyToggled(Keys.tmp)
	}
}

func HandleCamCtlKeys() {
	if CamCtl.MoveSpeedupFactor = 1; !paused {
		if Keys.Pressed[glfw.KeyLshift] {
			CamCtl.MoveSpeedupFactor = 10
		} else if Keys.Pressed[glfw.KeyRshift] {
			CamCtl.MoveSpeedupFactor = 100
		} else if Keys.Pressed[glfw.KeyLalt] {
			CamCtl.MoveSpeedupFactor = 0.1
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

func HandleToggleKeys() {
	if Keys.Toggled[glfw.KeyEsc] {
		ng.Loop.Stop()
	}
	if Keys.Toggled[glfw.KeyF2] {
		Cam.Rendering.States.FaceCulling = !Cam.Rendering.States.FaceCulling
	}
	if Keys.Toggled[glfw.KeyF3] {
		PauseResume()
	}
}
