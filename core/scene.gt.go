package core

import (
	"github.com/go-utils/uslice"
)

const sceneNodeChildCap = 16

//	Represents a scene graph.
type Scene struct {
	ID int

	allNodes  SceneNodeLib
	nodeCount int

	thrPrep struct {
		copyDone, done bool
	}

	thrRend struct {
		copyDone bool
	}
}

func (me *Scene) dispose() {
	me.allNodes.dispose()
	me.nodeCount = 0
}

func (me *Scene) init() {
	me.allNodes.init()
	root := &me.allNodes[me.allNodes.AddNew()]
	me.nodeCount = 1
	root.parentID, root.childNodeIDs, root.Render.skyMode, root.Render.Cull.Frustum = -1, make([]int, 0, sceneNodeChildCap), true, false
}

func (me *Scene) AddNewChildNode(parentNodeID, meshID int) (childNodeID int) {
	childNodeID = -1
	if me.allNodes.IsOk(parentNodeID) {
		me.nodeCount++
		childNodeID = me.allNodes.AddNew()
		me.allNodes[childNodeID].parentID, me.allNodes[childNodeID].Render.meshID = parentNodeID, meshID

		if len(me.allNodes[parentNodeID].childNodeIDs) == cap(me.allNodes[parentNodeID].childNodeIDs) {
			if len(me.allNodes[parentNodeID].childNodeIDs) > 0 {
				uslice.IntSetCap(&me.allNodes[parentNodeID].childNodeIDs, 2*len(me.allNodes[parentNodeID].childNodeIDs))
			} else {
				uslice.IntSetCap(&me.allNodes[parentNodeID].childNodeIDs, sceneNodeChildCap)
			}
		}
		uslice.IntAppendUnique(&me.allNodes[parentNodeID].childNodeIDs, childNodeID)
		me.ApplyNodeTransforms(childNodeID)

		var view int
		var rts *RenderTechniqueScene
		for canv := 0; canv < len(Core.Render.Canvases); canv++ {
			for view = 0; view < len(Core.Render.Canvases[canv].Views); view++ {
				if rts = Core.Render.Canvases[canv].Views[view].Technique_Scene(); rts != nil && rts.Camera.sceneID == me.ID {
					rts.Camera.initNodeCamData(me.allNodes, childNodeID)
				}
			}
		}
	}
	return
}

func (me *Scene) Node(id int) *SceneNode {
	return me.allNodes.get(id)
}

func (me *Scene) NumNodes() int {
	return me.nodeCount
}

func (me *Scene) ParentNodeID(childNodeID int) (parentID int) {
	if me.allNodes.IsOk(childNodeID) {
		parentID = me.allNodes[childNodeID].parentID
	}
	return
}

func (me *Scene) RemoveNode(fromID int) {
	if fromID > 0 {
		for i := 0; i < len(me.allNodes); i++ {
			if me.allNodes.Ok(i) && me.allNodes[i].parentID == fromID {
				me.RemoveNode(i)
			}
		}
		me.allNodes.Remove(fromID, 1)
		me.nodeCount--
	}
}

func (me *Scene) Root() *SceneNode {
	return &me.allNodes[0]
}

func (me *Scene) SetNodeMeshID(nodeID, meshID int) {
	if me.allNodes.IsOk(nodeID) {
		me.allNodes[nodeID].Render.meshID = meshID
		me.ApplyNodeTransforms(nodeID)
	}
}






//#begin-gt -gen-lib.gt T:Scene L:Core.Libs.Scenes

//	Only used for Core.Libs.Scenes
type SceneLib []Scene

func (me *SceneLib) AddNew() (id int) {
	id = -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			id = i
			break
		}
	}
	if id == -1 {
		if id = len(*me); id == cap(*me) {
			nu := make(SceneLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, Scene{})
	}
	ref := &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *SceneLib) Compact() {
	var (
		before, after []Scene
		ref           *Scene
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
			me.onSceneIDsChanged(changed)
		}
	}
}

func (me *SceneLib) init() {
	*me = make(SceneLib, 0, Options.Libs.InitialCap)
}

func (me *SceneLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me SceneLib) get(id int) (ref *Scene) {
	if me.IsOk(id) {
		ref = &me[id]
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
	if l := len(me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onSceneIDsChanged(changed)
	}
}

func (me SceneLib) Walk(on func(ref *Scene)) {
	for id := 0; id < len(me); id++ {
		if me.Ok(id) {
			on(&me[id])
		}
	}
}

//#end-gt
