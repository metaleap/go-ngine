package client

import (
)

type IWindowing interface {
	Exit()
	Init (winWidth, winHeight int, winFullScreen bool, vsync int, winTitle string) error
	KeyPressed (key int) bool
	KeysPressedAll2 (k1, k2 int) bool
	KeysPressedAll3 (k1, k2, k3 int) bool
	KeysPressedAny2 (k1, k2 int) bool
	KeysPressedAny3 (k1, k2, k3 int) bool
	KeyToggled (key int) bool
	OnLoop ()
	SetTitle (newTitle string)
	Time () float64
	WinHeight () int
	WinWidth () int
}
