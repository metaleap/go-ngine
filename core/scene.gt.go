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
	Core.Rendering.Canvases.Walk(nil, func(cam *Camera) {
		if cam.scene() == me {
			cam.SetScene(-1)
		}
	})
}

func (me *Scene) init() {
	me.RootNode = *newNode("", "", "", nil, me)
}

//#begin-gt -gen-lib2.gt T:Scene L:Scenes

//	Only used for Core.Libs.Scenes.
type SceneLib []Scene

func (_ SceneLib) AddNew() (ref *Scene) {
	me, id := &Core.Libs.Scenes, -1
	for i, _ := range *me {
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

func (_ SceneLib) Compact() {
	var (
		before, after []Scene
		ref           *Scene
		oldID         int
	)
	me := &Core.Libs.Scenes
	for i, _ := range *me {
		if (*me)[i].ID < 0 {
			before, after = (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	changed := make(map[int]int, len(*me))
	for i, _ := range *me {
		if (*me)[i].ID != i {
			ref = &(*me)[i]
			oldID, ref.ID = ref.ID, i
			changed[oldID] = i
		}
	}
	if len(changed) > 0 {
		Core.Libs.Scenes.onSceneIDsChanged(changed)
		Options.Libs.OnIDsChanged.Scenes(changed)
	}
}

func (_ SceneLib) ctor() {
	me := &Core.Libs.Scenes
	*me = make(SceneLib, 0, Options.Libs.InitialCap)
}

func (_ SceneLib) dispose() {
	me := &Core.Libs.Scenes
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (_ SceneLib) Get(id int) (ref *Scene) {
	if id >= 0 && id < len(Core.Libs.Scenes) {
		if ref = &Core.Libs.Scenes[id]; ref.ID != id {
			ref = nil
		}
	}
	return
}

func (_ SceneLib) Has(id int) (has bool) {
	if id >= 0 && id < len(Core.Libs.Scenes) {
		has = Core.Libs.Scenes[id].ID == id
	}
	return
}

func (_ SceneLib) Remove(fromID, num int) {
	me := &Core.Libs.Scenes
	if l := len(*me); fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			(*me)[id].dispose()
			changed[id], (*me)[id].ID = -1, -1
		}
		Core.Libs.Scenes.onSceneIDsChanged(changed)
		Options.Libs.OnIDsChanged.Scenes(changed)
	}
}

//#end-gt
