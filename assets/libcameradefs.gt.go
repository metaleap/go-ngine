package assets

//#begin-gt lib_tmpl.gt T:CameraDef

var (
	AllCameraDefLibs = LibsCameraDef {}
	CameraDefs = AllCameraDefLibs.AddNew("")
)

type LibsCameraDef map[string]*LibCameraDefs

	func (me LibsCameraDef) AddNew (id string) (lib *LibCameraDefs) {
		if me[id] != nil { return }
		lib = me.New(id)
		me[id] = lib
		return
	}

	func (me LibsCameraDef) New (id string) (lib *LibCameraDefs) {
		lib = newLibCameraDefs(id)
		return
	}

type LibCameraDefs struct {
	baseLib
	M map[string]*CameraDef
}

	func newLibCameraDefs (id string) (me *LibCameraDefs) {
		me = &LibCameraDefs {}
		me.base.init(id)
		me.M = map[string]*CameraDef {}
		return
	}

	func (me *LibCameraDefs) Add (def *CameraDef) *CameraDef {
		if me.M[def.ID] != nil { return nil }
		me.M[def.ID] = def
		return def
	}

	func (me *LibCameraDefs) AddNew (id string) (nd *CameraDef) {
		if me.M[id] != nil { return }
		nd = me.New(id)
		me.M[id] = nd
		return
	}

	func (me *LibCameraDefs) Get (id string) *CameraDef {
		return me.M[id]
	}

	func (me *LibCameraDefs) New (id string) (d *CameraDef) {
		d = newCameraDef(id)
		return
	}

	func (me *LibCameraDefs) Remove (id string) {
		delete(me.M, id)
	}

//#end-gt
