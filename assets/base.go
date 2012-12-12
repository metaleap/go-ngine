package assets

//	Provides a common base for *Def*s, *Inst*s and *Lib*s.
type Base struct {
	//	This callback, set by *core* (or your custom package), gets called by the *SyncChanges()* method.
	//	This is the ultimate point in the sync chain where *core* (or your custom package) picks up the changed
	//	contents of this *Def*, *Inst* or *Lib*. If this is a *Lib* this gets called after all *Defs* in it
	//	have synced.
	OnSync func()

	dirty bool
}

//	You NEED to call this method if you modified this *Def* or *Inst* by setting its fields directly
//	(instead of using the provided *SetFoo()* methods) for your changes to be picked up by *core* (or your custom package).
func (me *Base) SetDirty() {
	me.dirty = true
}

//	If field does not equal val, sets field to val and calls SetDirty().
func (me *Base) SetFieldB(field *bool, val bool) {
	if *field != val {
		*field = val
		me.SetDirty()
	}
}

//	If field does not equal val, sets field to val and calls SetDirty().
func (me *Base) SetFieldF(field *float64, val float64) {
	if *field != val {
		*field = val
		me.SetDirty()
	}
}

func (me *Base) init() {
	me.OnSync = func() {}
	me.SetDirty()
}

//	Signals to *core* (or your custom package) that changes have been made to this *Def*, *Inst* or *Lib* that need to be picked up.
//	Call this after you have made any number of changes to this this *Def*, *Inst* or *Lib*.
//	Also called by the global *SyncChanges()* function.
func (me *Base) SyncChanges() {
	if me.dirty {
		me.dirty = false
		me.OnSync()
	}
}

//	Provides a common base for *Def*s.
type BaseDef struct {
	//	Syncability
	Base
	//	Unique identifier
	HasId
	//	Pretty-print name/title
	HasName
	//	Asset meta-data
	HasAsset
	//	Custom-technique/foreign-profile meta-data
	HasExtras
}

//	Provides a common base for *Inst*s.
type BaseInst struct {
	//	Syncability
	Base
	//	Pretty-print name/title
	HasName
	//	Scoped identifier
	HasSid
	//	Custom-technique/foreign-profile meta-data
	HasExtras
	//	The unique ID of the definition that this instance refers to.
	DefRef RefId
}

//	Provides a common base for *Lib*s.
type BaseLib struct {
	//	Syncability
	Base
	//	Unique identifier
	HasId
	//	Pretty-print name/title
	HasName
}
