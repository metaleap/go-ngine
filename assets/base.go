package assets

//	Allows for notifying outside packages of changes.
type BaseSync struct {
	//	This callback, set by the core package (or your custom package), gets called by the
	//	SyncChanges() method. This is the ultimate point in the sync chain where the core package
	//	(or your custom package) picks up the changed contents of this resource.
	//	If the parent is a Lib then this gets called after all its Defs have synced.
	OnSync func()

	dirty bool
}

//	You NEED to call this method if you modified this Def or Inst by setting its fields directly
//	(instead of using the provided SetFoo() or SetFieldX() methods) for your changes
//	to be picked up by the core package (or your custom package).
func (me *BaseSync) SetDirty() {
	me.dirty = true
}

//	If field does not equal val, sets field to val and calls SetDirty().
func (me *BaseSync) SetFieldB(field *bool, val bool) {
	if *field != val {
		*field = val
		me.SetDirty()
	}
}

//	If field does not equal val, sets field to val and calls SetDirty().
func (me *BaseSync) SetFieldF(field *float64, val float64) {
	if *field != val {
		*field = val
		me.SetDirty()
	}
}

func (me *BaseSync) init() {
	me.OnSync = func() {}
	me.SetDirty()
}

//	Signals to the core package (or your custom package) that changes have been made to this
//	Def, Inst or Lib resource that need to be picked up. Call this after you have made a number
//	of changes to this this resource. Also called by the global SyncChanges() function.
func (me *BaseSync) SyncChanges() {
	if me.dirty {
		me.dirty = false
		me.OnSync()
	}
}

//	Provides a common base for resource definitions.
type BaseDef struct {
	//	Syncability
	BaseSync
	//	Id
	HasId
	//	Name
	HasName
	//	Asset
	HasAsset
	//	Extras
	HasExtras
}

//	Provides a common base for resource instantiations.
type BaseInst struct {
	//	Syncability
	BaseSync
	//	Name
	HasName
	//	Sid
	HasSid
	//	Extras
	HasExtras
	//	The Id of the resource definition referenced by this instance.
	DefRef RefId
}

//	Provides a common base for resource libraries.
type BaseLib struct {
	//	Syncability
	BaseSync
	//	Id
	HasId
	//	Name
	HasName
}
