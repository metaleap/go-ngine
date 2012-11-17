package assets

//#begin-gt lib_tmpl.gt T:Scene

var (
	AllSceneLibs = LibsScene {}
	Scenes = AllSceneLibs.AddNew("")
)

type LibsScene map[string]*LibScenes

	func (me LibsScene) AddNew (id string) (lib *LibScenes) {
		if me[id] != nil { return }
		lib = me.New(id)
		me[id] = lib
		return
	}

	func (me LibsScene) New (id string) (lib *LibScenes) {
		lib = newLibScenes(id)
		return
	}

type LibScenes struct {
	baseLib
	M map[string]*Scene
}

	func newLibScenes (id string) (me *LibScenes) {
		me = &LibScenes {}
		me.base.init(id)
		me.M = map[string]*Scene {}
		return
	}

	func (me *LibScenes) Add (def *Scene) *Scene {
		if me.M[def.ID] != nil { return nil }
		me.M[def.ID] = def
		return def
	}

	func (me *LibScenes) AddNew (id string) (nd *Scene) {
		if me.M[id] != nil { return }
		nd = me.New(id)
		me.M[id] = nd
		return
	}

	func (me *LibScenes) Get (id string) *Scene {
		return me.M[id]
	}

	func (me *LibScenes) New (id string) (d *Scene) {
		d = newScene(id)
		return
	}

	func (me *LibScenes) Remove (id string) {
		delete(me.M, id)
	}

//#end-gt
