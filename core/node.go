package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type subNodes struct {
	M     map[string]*Node
	owner *Node
}

func newSubNodes(owner *Node) (nodes *subNodes) {
	nodes = &subNodes{owner: owner, M: map[string]*Node{}}
	return
}

func (me *subNodes) Add(node *Node) {
	if node.parentNode != nil {
		node.parentNode.SubNodes.Remove(node.name)
	}
	node.parentNode = me.owner
	me.M[node.name] = node
}

func (me *subNodes) Get(names ...string) (nodes []*Node) {
	nodes = make([]*Node, len(names))
	for curIndex, curStr = range names {
		nodes[curIndex] = me.M[curStr]
	}
	return
}

func (me *subNodes) Make(nodeName, meshName, modelName string) (node *Node) {
	node = newNode(nodeName, meshName, modelName, me.owner)
	me.Add(node)
	return
}

func (me *subNodes) MakeN(nodeMeshModelNames ...string) {
	for i := 2; i < len(nodeMeshModelNames); i += 3 {
		me.Make(nodeMeshModelNames[i-2], nodeMeshModelNames[i-1], nodeMeshModelNames[i])
	}
}

func (me *subNodes) Remove(name string) {
	if node := me.M[name]; node != nil {
		node.parentNode = nil
	}
	delete(me.M, name)
}

type Node struct {
	Disabled  bool
	SubNodes  *subNodes
	Transform *Transforms

	mat                                        *Material
	mesh                                       *Mesh
	model                                      *Model
	curKey, matName, meshName, modelName, name string
	curSubNode, parentNode                     *Node
}

func newNode(nodeName, meshName, modelName string, parent *Node) (me *Node) {
	me = &Node{name: nodeName, parentNode: parent}
	me.SubNodes = newSubNodes(me)
	me.SetMeshModelName(meshName, modelName)
	me.Transform = newTransforms(me)
	return
}

func (me *Node) Material() *Material {
	if me.mat != nil {
		return me.mat
	}
	return me.model.mat
}

func (me *Node) MatName() string {
	return me.matName
}

func (me *Node) MeshName() string {
	return me.meshName
}

func (me *Node) MeshModelName() string {
	return me.modelName
}

func (me *Node) render() {
	if !me.Disabled {
		curNode, curMesh, curModel = me, me.mesh, me.model
		if me.model != nil {
			curTechnique.onRenderNode()
			gl.UniformMatrix4fv(curProg.UnifLocs["uMatModelView"], 1, gl.FALSE, &me.Transform.glMatModelView[0])
			me.model.render()
		}
		for me.curKey, me.curSubNode = range me.SubNodes.M {
			me.curSubNode.render()
		}
	}
}

func (me *Node) SetMatName(newMatName string) {
	if newMatName != me.matName {
		me.mat, me.matName = Core.Materials[newMatName], newMatName
	}
}

func (me *Node) SetMeshModelName(meshName, modelName string) {
	if meshName != me.meshName {
		me.mesh, me.meshName = Core.Meshes[meshName], meshName
		me.model, me.modelName = me.mesh.Models.Default(), ""
	}
	if modelName != me.modelName {
		me.model, me.modelName = me.mesh.Models[modelName], modelName
	}
}

func (me *Node) Transforms() *Transforms {
	return me.Transform
}

func (me *Node) ChildrenUpdateMatrices() {
	for _, me.curSubNode = range me.SubNodes.M {
		me.curSubNode.Transform.UpdateMatrices()
	}
}

func (me *Node) Parent() Transformable {
	return me.parentNode
}
