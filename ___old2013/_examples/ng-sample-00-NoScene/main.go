// Shows only the splash screen without loading or rendering any scene data.
package main

import (
	apputil "github.com/metaleap/go-ngine/___old2013/_examples/shared-utils"
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
