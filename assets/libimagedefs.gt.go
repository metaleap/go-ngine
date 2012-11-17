package assets

//#begin-gt lib_tmpl.gt T:ImageDef

var (
	AllImageDefLibs = LibsImageDef {}
	ImageDefs = AllImageDefLibs.AddNew("")
)

type LibsImageDef map[string]*LibImageDefs

	func (me LibsImageDef) AddNew (id string) (lib *LibImageDefs) {
		if me[id] != nil { return }
		lib = me.New(id)
		me[id] = lib
		return
	}

	func (me LibsImageDef) New (id string) (lib *LibImageDefs) {
		lib = newLibImageDefs(id)
		return
	}

type LibImageDefs struct {
	baseLib
	M map[string]*ImageDef
}

	func newLibImageDefs (id string) (me *LibImageDefs) {
		me = &LibImageDefs {}
		me.base.init(id)
		me.M = map[string]*ImageDef {}
		return
	}

	func (me *LibImageDefs) Add (def *ImageDef) *ImageDef {
		if me.M[def.ID] != nil { return nil }
		me.M[def.ID] = def
		return def
	}

	func (me *LibImageDefs) AddNew (id string) (nd *ImageDef) {
		if me.M[id] != nil { return }
		nd = me.New(id)
		me.M[id] = nd
		return
	}

	func (me *LibImageDefs) Get (id string) *ImageDef {
		return me.M[id]
	}

	func (me *LibImageDefs) New (id string) (d *ImageDef) {
		d = newImageDef(id)
		return
	}

	func (me *LibImageDefs) Remove (id string) {
		delete(me.M, id)
	}

//#end-gt
