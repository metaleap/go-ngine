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
	model *TModel
	curKey, matName, meshName, modelName, name string
	curSubNode, parentNode *TNode
}

func newNode (nodeName, meshName, modelName string, parent *TNode) *TNode {
	var node = &TNode {}
	node.name = nodeName
	node.parentNode = parent
	node.SubNodes = map[string]*TNode {}
	node.SetMeshModelName(meshName, modelName)
	node.Transform = newTransform(node)
	return node
}

func (me *TNode) AttachSubNode (node *TNode) {
	me.SubNodes[node.name] = node
	node.parentNode = me
}

func (me *TNode) DetachFromParent () {
	if me.parentNode != nil { me.parentNode.RemoveSubNode(me.name) }
}

func (me *TNode) MakeSubNode (nodeName, meshName, modelName string) *TNode {
	me.SubNodes[nodeName] = newNode(nodeName, meshName, modelName, me)
	return me.SubNodes[nodeName]
}

func (me *TNode) MakeSubNodes (nodeMeshModelNames ... string) {
	for i := 2; i < len(nodeMeshModelNames); i += 3 {
		me.MakeSubNode(nodeMeshModelNames[i - 2], nodeMeshModelNames[i - 1], nodeMeshModelNames[i])
	}
}

func (me *TNode) Material () *TMaterial {
	if me.mat != nil { return me.mat }
	return me.model.mat
}

func (me *TNode) MatName () string {
	return me.matName
}

func (me *TNode) MeshName () string {
	return me.meshName
}

func (me *TNode) MeshModelName () string {
	return me.modelName
}

func (me *TNode) RemoveSubNode (name string) {
	if node := me.SubNodes[name]; node != nil { node.parentNode = nil }
	delete(me.SubNodes, name)
}

func (me *TNode) render () {
	if (!me.Disabled) {
		curNode, curMesh, curModel = me, me.mesh, me.model
		if (me.model != nil) {
			curTechnique.onRenderNode()
			gl.UniformMatrix4fv(curProg.UnifLocs["uMatModelView"], 1, gl.FALSE, &me.Transform.glMatModelView[0])
			me.model.render()
		}
		for me.curKey, me.curSubNode = range me.SubNodes {
			me.curSubNode.render()
		}
	}
}

func (me *TNode) SetMatName (newMatName string) {
	if newMatName != me.matName {
		me.mat, me.matName = Core.Materials[newMatName], newMatName
	}
}

func (me *TNode) SetMeshModelName (meshName, modelName string) {
	if meshName != me.meshName {
		me.mesh, me.meshName = Core.Meshes[meshName], meshName
		me.model, me.modelName = me.mesh.Models.Default(), ""
	}
	if modelName != me.modelName {
		me.model, me.modelName = me.mesh.Models[modelName], modelName
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
