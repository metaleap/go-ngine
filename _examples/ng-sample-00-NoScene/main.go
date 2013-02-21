package main

import (
	apputil "github.com/go3d/go-ngine/_examples/shared-utils"
)

func main() {
	apputil.MaxKeyHint = 1
	//	You better check out this function, it's part of the "minimal go:ngine setup":
	apputil.Main(nil, func() {}, onWinThread)
}

//	called once per frame in main thread
func onWinThread() {
	apputil.CheckAndHandleToggleKeys()
}
