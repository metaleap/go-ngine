package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type TNode struct {
	Disabled bool
	SubNodes map[string]*TNode
	Transform *TTransform

	mat *TMaterial
	mesh *TMesh
	glVertTexCoordsBuf gl.Uint
	curKey, matKey, meshKey string
	curSubNode, parentNode *TNode
}

func newNode (meshKey string, parent *TNode) *TNode {
	var node = &TNode {}
	node.parentNode = parent
	node.SubNodes = map[string]*TNode {}
	node.SetMeshKey(meshKey)
	node.Transform = NewTransform(node)
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
	if me.glVertTexCoordsBuf > 0 { gl.DeleteBuffers(1, &me.glVertTexCoordsBuf) }
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

func (me *TNode) SetMatKey (newMatKey string, texCoords []gl.Float) {
	if newMatKey != me.matKey {
		me.mat, me.matKey = Core.Materials[newMatKey], newMatKey
	}
	if (me.mat != nil) && (me.glVertTexCoordsBuf == 0) && (len(texCoords) > 0) {
		gl.GenBuffers(1, &me.glVertTexCoordsBuf)
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glVertTexCoordsBuf)
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(texCoords)), gl.Pointer(&texCoords[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}
}

func (me *TNode) SetMeshKey (newMeshKey string) {
	if newMeshKey != me.meshKey {
		me.mesh, me.meshKey = Core.Meshes[newMeshKey], newMeshKey
	}
}

func (me *TNode) transform () *TTransform {
	return me.Transform
}

func (me *TNode) transformChildrenUpdateMatrices () {
	for _, me.curSubNode = range me.SubNodes {
		me.curSubNode.Transform.updateMatrices()
	}
}

func (me *TNode) transformParent () ITransformable {
	return me.parentNode
}
