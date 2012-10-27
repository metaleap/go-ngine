package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type TNode struct {
	Disabled bool
	SubNodes map[string]*TNode
	Transform *tTransform

	mat *TMaterial
	mesh *TMesh
	curKey, matKey, meshKey string
	curSubNode, parentNode *TNode
}

func newNode (meshKey string, parent *TNode) *TNode {
	var node = &TNode {}
	node.parentNode = parent
	node.SubNodes = map[string]*TNode {}
	node.SetMeshKey(meshKey)
	node.Transform = newTransform(node)
	return node
}

func (me *TNode) AddSubNodes (meshKeys ... string) {
	for _, meshKey := range meshKeys { me.SubNodes[meshKey] = newNode(meshKey, me) }
}

func (me *TNode) AddSubNodesNamed (nodeAndMeshKeys map[string]string) {
	for nodeKey, meshKey := range nodeAndMeshKeys { me.SubNodes[nodeKey] = newNode(meshKey, me) }
}

func (me *TNode) Dispose () {
	for _, subNode := range me.SubNodes { subNode.Dispose() }
	me.SubNodes = map[string]*TNode {}
}

func (me *TNode) MatKey () string {
	return me.matKey
}

func (me *TNode) MeshKey () string {
	return me.meshKey
}

func (me *TNode) render () {
	if (!me.Disabled) {
		curNode = me
		if (me.mesh != nil) {
			curTechnique.OnRenderNode()
			gl.UniformMatrix4fv(curProg.UnifLocs["uMatModelView"], 1, gl.FALSE, &me.Transform.glMatModelView[0])
			me.mesh.render()
		}
		for me.curKey, me.curSubNode = range me.SubNodes {
			me.curSubNode.render()
		}
	}
}

func (me *TNode) SetMatKey (newMatKey string) {
	if newMatKey != me.matKey {
		me.mat, me.matKey = Core.Materials[newMatKey], newMatKey
	}
}

func (me *TNode) SetMeshKey (newMeshKey string) {
	if newMeshKey != me.meshKey {
		me.mesh, me.meshKey = Core.Meshes[newMeshKey], newMeshKey
	}
}

func (me *TNode) transform () *tTransform {
	return me.Transform
}

func (me *TNode) transformChildrenUpdateMatrices () {
	for _, me.curSubNode = range me.SubNodes {
		me.curSubNode.Transform.updateMatrices()
	}
}

func (me *TNode) transformParent () iTransformable {
	return me.parentNode
}
