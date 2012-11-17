package assets

//#begin-gt lib_tmpl.gt T:NodeDef

var (
	AllNodeDefLibs = LibsNodeDef {}
	NodeDefs = AllNodeDefLibs.AddNew("")
)

type LibsNodeDef map[string]*LibNodeDefs

	func (me LibsNodeDef) AddNew (id string) (lib *LibNodeDefs) {
		if me[id] != nil { return }
		lib = me.New(id)
		me[id] = lib
		return
	}

	func (me LibsNodeDef) New (id string) (lib *LibNodeDefs) {
		lib = newLibNodeDefs(id)
		return
	}

type LibNodeDefs struct {
	baseLib
	M map[string]*NodeDef
}

	func newLibNodeDefs (id string) (me *LibNodeDefs) {
		me = &LibNodeDefs {}
		me.base.init(id)
		me.M = map[string]*NodeDef {}
		return
	}

	func (me *LibNodeDefs) Add (def *NodeDef) *NodeDef {
		if me.M[def.ID] != nil { return nil }
		me.M[def.ID] = def
		return def
	}

	func (me *LibNodeDefs) AddNew (id string) (nd *NodeDef) {
		if me.M[id] != nil { return }
		nd = me.New(id)
		me.M[id] = nd
		return
	}

	func (me *LibNodeDefs) Get (id string) *NodeDef {
		return me.M[id]
	}

	func (me *LibNodeDefs) New (id string) (d *NodeDef) {
		d = newNodeDef(id)
		return
	}

	func (me *LibNodeDefs) Remove (id string) {
		delete(me.M, id)
	}

//#end-gt
