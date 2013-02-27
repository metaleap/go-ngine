package core

//	Represents a scene graph.
type Scene struct {
	ID int

	//	The root Node for this scene graph.
	RootNode Node
}

func NewScene() (me *Scene) {
	me = new(Scene)
	me.init()
	return
}

func (me *Scene) dispose() {
}

func (me *Scene) init() {
	me.RootNode = *newNode("", -1, nil, me)
}

//#begin-gt -gen-lib.gt T:Scene L:Scenes

//	Only used for Core.Libs.Scenes.
type SceneLib []Scene

func (me *SceneLib) AddNew() (ref *Scene) {
	id := -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			id = i
			break
		}
	}
	if id < 0 {
		if id = len(*me); id == cap(*me) {
			nu := make(SceneLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, Scene{})
	}
	ref = &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *SceneLib) Compact() {
	var (
		before, after []Scene
		ref           *Scene
		oldID, i      int
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			before, after = (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	changed := make(map[int]int, len(*me))
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID != i {
			ref = &(*me)[i]
			oldID, ref.ID = ref.ID, i
			changed[oldID] = i
		}
	}
	if len(changed) > 0 {
		me.onSceneIDsChanged(changed)
		Options.Libs.OnIDsChanged.Scenes(changed)
	}
}

func (me *SceneLib) ctor() {
	*me = make(SceneLib, 0, Options.Libs.InitialCap)
}

func (me *SceneLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me SceneLib) Get(id int) (ref *Scene) {
	if id > -1 && id < len(me) {
		if ref = &me[id]; ref.ID != id {
			ref = nil
		}
	}
	return
}

func (me SceneLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me SceneLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me SceneLib) Remove(fromID, num int) {
	if l := len(me); fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onSceneIDsChanged(changed)
		Options.Libs.OnIDsChanged.Scenes(changed)
	}
}

func (me SceneLib) Walk(on func(ref *Scene)) {
	for id := 0; id < len(me); id++ {
		if me[id].ID > -1 {
			on(&me[id])
		}
	}
}

//#end-gt
