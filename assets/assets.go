package assets

import (
	"fmt"
)

var (
	//	This callback, set by the core package (or your custom package),
	//	gets called before SyncChanges() proceeds with syncing.
	OnBeforeSyncAll func()
	//	This callback, set by the core package (or your custom package),
	//	gets called after SyncChanges() has finished syncing.
	OnAfterSyncAll func()

	syncHandlers []func()
)

func init() {
	OnBeforeSyncAll = func() {}
	OnAfterSyncAll = func() {}
}

func sfmt(format string, fmtArgs ...interface{}) string {
	return fmt.Sprintf(format, fmtArgs...)
}

//	Returns a ScopedFloat with the specified value and no Sid.
func SidF(f float64) (sf *SidFloat) {
	sf = &SidFloat{F: f}
	return
}

//	Signals to the core package (or your custom package) that changes have been made that need to be
//	picked up. Call this after you have made any number of changes to your Defs, Insts or Libs.
func SyncChanges() {
	OnBeforeSyncAll()
	for _, syncer := range syncHandlers {
		syncer()
	}
	OnAfterSyncAll()
}
