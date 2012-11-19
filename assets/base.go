package assets

//	Provides a common base for *Def*s, *Inst*s and *Lib*s.
type Base struct {
	//	The unique identifier of this *Def*, *Inst* or *Lib*.
	ID string

	//	This callback, set by *core* (or your custom package), gets called by the *SyncChanges()* method.
	//	This is the ultimate point in the sync chain where *core* (or your custom package) picks up the changed
	//	contents of this *Def*, *Inst* or *Lib*. If this is a *Lib* this gets called after all *Defs* in it
	//	have synced.
	OnSync func ()

	timeLastChanged, timeLastSynced int64
}

	//	You NEED to call this method if you modified this *Def* or *Inst* by setting its fields directly
	//	(instead of using the provided *SetXyz()* methods) for your changes to be picked up by *core* (or your custom package).
	func (me *Base) SetDirty () {
		me.timeLastChanged = now()
	}

	func (me *Base) init (id string) {
		me.ID = id
		me.OnSync = func () {}
		me.SetDirty()
	}

	//	Signals to *core* (or your custom package) that changes have been made to this *Def*, *Inst* or *Lib* that need to be picked up.
	//	Call this after you have made any number of changes to this this *Def*, *Inst* or *Lib*.
	//	Also called by the global *SyncChanges()* function.
	func (me *Base) SyncChanges () {
		if me.timeLastChanged > me.timeLastSynced {
			me.OnSync()
			me.timeLastSynced = now()
		}
	}

//	Provides a common base for *Def*s.
type BaseDef struct {
	//	Provides ID and syncing
	Base
}

//	Provides a common base for *Inst*s.
type BaseInst struct {
	//	Provides ID and syncing
	Base
}

//	Provides a common base for *Lib*s.
type BaseLib struct {
	//	Provides ID and syncing
	Base
}
