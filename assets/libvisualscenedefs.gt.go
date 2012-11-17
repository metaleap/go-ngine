package assets

//#begin-gt lib_tmpl.gt T:VisualSceneDef

var (
	AllVisualSceneDefLibs = LibsVisualSceneDef {}
	VisualSceneDefs = AllVisualSceneDefLibs.AddNew("")
)

type LibsVisualSceneDef map[string]*LibVisualSceneDefs

	func (me LibsVisualSceneDef) AddNew (id string) (lib *LibVisualSceneDefs) {
		if me[id] != nil { return }
		lib = me.New(id)
		me[id] = lib
		return
	}

	func (me LibsVisualSceneDef) New (id string) (lib *LibVisualSceneDefs) {
		lib = newLibVisualSceneDefs(id)
		return
	}

type LibVisualSceneDefs struct {
	baseLib
	M map[string]*VisualSceneDef
}

	func newLibVisualSceneDefs (id string) (me *LibVisualSceneDefs) {
		me = &LibVisualSceneDefs {}
		me.base.init(id)
		me.M = map[string]*VisualSceneDef {}
		return
	}

	func (me *LibVisualSceneDefs) Add (def *VisualSceneDef) *VisualSceneDef {
		if me.M[def.ID] != nil { return nil }
		me.M[def.ID] = def
		return def
	}

	func (me *LibVisualSceneDefs) AddNew (id string) (nd *VisualSceneDef) {
		if me.M[id] != nil { return }
		nd = me.New(id)
		me.M[id] = nd
		return
	}

	func (me *LibVisualSceneDefs) Get (id string) *VisualSceneDef {
		return me.M[id]
	}

	func (me *LibVisualSceneDefs) New (id string) (d *VisualSceneDef) {
		d = newVisualSceneDef(id)
		return
	}

	func (me *LibVisualSceneDefs) Remove (id string) {
		delete(me.M, id)
	}

//#end-gt
