package assets

import (
	"fmt"
	"time"
)

var (
	//	This callback, set by *core* (or your custom package), gets called before SyncChanges() proceeds with syncing.
	OnBeforeSyncAll func()

	//	This callback, set by *core* (or your custom package), gets called after SyncChanges() has finished syncing.
	OnAfterSyncAll func()

	//	Your default unit-in-meters for geometry, coordinates and transformations.
	//	If a unit represents:
	//	- a meter, set to 1;
	//	- a centimeter, set to 0.01;
	//	- a kilometer, set to 1000;
	//	- an inch, set to 0.02539999969303608... etc.
	//	The *assets* package does not support multiple different or individual per-asset units.
	UnitInMeters float64 = 1

	syncHandlers []func()
)

func init() {
	OnBeforeSyncAll = func() {}
	OnAfterSyncAll = func() {}
}

func now() int64 {
	return time.Now().UnixNano()
}

func sfmt(format string, fmtArgs ...interface{}) string {
	return fmt.Sprintf(format, fmtArgs...)
}

//	Signals to *core* (or your custom package) that changes have been made that need to be picked up. Call this after you have made any number of changes to your Defs, Insts or Libs.
func SyncChanges() {
	OnBeforeSyncAll()
	for _, syncer := range syncHandlers {
		syncer()
	}
	OnAfterSyncAll()
}
