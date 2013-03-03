package core

import (
	unum "github.com/metaleap/go-util/num"
	usl "github.com/metaleap/go-util/slice"
)

const sceneNodeChildCap = 16

//	Represents a scene graph.
type Scene struct {
	ID int

	allNodes SceneNodeLib

	thrPrep struct {
		copyDone, done bool
	}

	thrRend struct {
		copyDone bool
	}
}

func (me *Scene) dispose() {
	me.allNodes.dispose()
}

func (me *Scene) init() {
	me.allNodes.init()
	root := &me.allNodes[me.allNodes.AddNew()]
	root.childNodeIDs = make([]int, 0, sceneNodeChildCap)
	root.Render.skyMode = true
	root.parentID = -1
}

func (me *Scene) initCamNodeData(nodeID int) {
	var view int
	var rts *RenderTechniqueScene
	for canv := 0; canv < len(Core.Render.Canvases); canv++ {
		for view = 0; view < len(Core.Render.Canvases[canv].Views); view++ {
			if rts = Core.Render.Canvases[canv].Views[view].Technique_Scene(); rts != nil && rts.Camera.sceneID == me.ID {
				rts.Camera.initNodeCamData(&me.allNodes[nodeID])
			}
		}
	}
}

func (me *Scene) AddNewChildNode(parentNodeID int) (childNodeID int) {
	childNodeID = -1
	if me.allNodes.IsOk(parentNodeID) {
		childNodeID = me.allNodes.AddNew()
		me.allNodes[childNodeID].parentID = parentNodeID

		if len(me.allNodes[parentNodeID].childNodeIDs) == cap(me.allNodes[parentNodeID].childNodeIDs) {
			nuCap := sceneNodeChildCap
			if len(me.allNodes[parentNodeID].childNodeIDs) > 0 {
				nuCap = len(me.allNodes[parentNodeID].childNodeIDs) * 2
			}
			nu := make([]int, len(me.allNodes[parentNodeID].childNodeIDs), nuCap)
			copy(nu, me.allNodes[parentNodeID].childNodeIDs)
			me.allNodes[parentNodeID].childNodeIDs = nu
		}
		usl.IntAppendUnique(&me.allNodes[parentNodeID].childNodeIDs, childNodeID)

		me.initCamNodeData(childNodeID)
		me.ApplyNodeTransforms(childNodeID)
	}
	return
}

//	Updates the internal 4x4 transformation matrix for all transformations of the specified
//	node and child-nodes. It is only this matrix that is used by the rendering runtime.
func (me *Scene) ApplyNodeTransforms(nodeID int) {
	if me.allNodes.IsOk(nodeID) {
		//	this node
		var matParent, matTrans, matScale, matRotX, matRotY, matRotZ unum.Mat4
		matScale.Scaling(&me.allNodes[nodeID].Transform.Scale)
		matTrans.Translation(&me.allNodes[nodeID].Transform.Pos)
		matRotX.RotationX(me.allNodes[nodeID].Transform.Rot.X)
		matRotY.RotationY(me.allNodes[nodeID].Transform.Rot.Y)
		matRotZ.RotationZ(me.allNodes[nodeID].Transform.Rot.Z)
		if me.allNodes[nodeID].parentID < 0 {
			matParent.Identity()
		} else {
			matParent.CopyFrom(&me.allNodes[me.allNodes[nodeID].parentID].Transform.thrApp.matModelView)
		}
		me.allNodes[nodeID].Transform.thrApp.matModelView.SetFromMultN(&matParent, &matTrans /*me.Other,*/, &matScale, &matRotX, &matRotY, &matRotZ)
		//	child-nodes
		for i := 0; i < len(me.allNodes[nodeID].childNodeIDs); i++ {
			me.ApplyNodeTransforms(me.allNodes[nodeID].childNodeIDs[i])
		}
	}
}

func (me *Scene) Node(id int) *SceneNode {
	return me.allNodes.get(id)
}

func (me *Scene) ParentNodeID(childNodeID int) (parentID int) {
	if me.allNodes.IsOk(childNodeID) {
		parentID = me.allNodes[childNodeID].parentID
	}
	return
}

func (me *Scene) RemoveNode(fromID int) {
	if fromID > 0 {
		me.allNodes.Remove(fromID, 1)
	}
}

func (me *Scene) Root() *SceneNode {
	return &me.allNodes[0]
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
