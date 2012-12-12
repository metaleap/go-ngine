package assets

import (
	"fmt"
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
	//	This is ONLY used when importing assets that specify their own unit-in-meters, those will be re-scaled to this unit.
	//	If you need to customize this value, do so before populating the *assets* package's libraries.
	UnitInMeters float64 = 1

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
func Scopedf(f float64) (sf ScopedFloat) {
	sf.F = f
	return
}

//	Signals to *core* (or your custom package) that changes have been made that need to be picked up. Call this after you have made any number of changes to your Defs, Insts or Libs.
func SyncChanges() {
	OnBeforeSyncAll()
	for _, syncer := range syncHandlers {
		syncer()
	}
	OnAfterSyncAll()
}
