package core

type sceneChildNodes []int

//#begin-gt -gen-lib.gt T:SceneNode L:Core.Scenes[id].allNodes

//	Only used for Core.Scenes[id].allNodes
type SceneNodeLib []SceneNode

func (me *SceneNodeLib) AddNew() (id int) {
	id = -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			id = i
			break
		}
	}
	if id == -1 {
		if id = len(*me); id == cap(*me) {
			nu := make(SceneNodeLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, SceneNode{})
	}
	ref := &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *SceneNodeLib) Compact() {
	var (
		before, after []SceneNode
		ref           *SceneNode
		oldID, i      int
		compact       bool
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			compact, before, after = true, (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	if compact {
		changed := make(map[int]int, len(*me))
		for i = 0; i < len(*me); i++ {
			if ref = &(*me)[i]; ref.ID != i {
				oldID, ref.ID = ref.ID, i
				changed[oldID] = i
			}
		}
		if len(changed) > 0 {
			me.onSceneNodeIDsChanged(changed)
		}
	}
}

func (me *SceneNodeLib) init() {
	*me = make(SceneNodeLib, 0, Options.Libs.InitialCap)
}

func (me *SceneNodeLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me SceneNodeLib) get(id int) (ref *SceneNode) {
	if me.IsOk(id) {
		ref = &me[id]
	}
	return
}

func (me SceneNodeLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me SceneNodeLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me SceneNodeLib) Remove(fromID, num int) {
	if l := len(me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onSceneNodeIDsChanged(changed)
	}
}

func (me SceneNodeLib) Walk(on func(ref *SceneNode)) {
	for id := 0; id < len(me); id++ {
		if me.Ok(id) {
			on(&me[id])
		}
	}
}

//#end-gt
