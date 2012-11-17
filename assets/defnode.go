package assets

type NodeDef struct {
	baseDef
	Layers map[string]bool
	NodeDefs []*NodeDef
	NodeInsts []*NodeInst
	Transforms *Transforms
}

func newNodeDef (id string) (me *NodeDef) {
	me = &NodeDef {}
	me.base.init(id)
	me.Layers = map[string]bool {}
	me.NodeDefs = []*NodeDef {}
	me.NodeInsts = []*NodeInst {}
	me.Transforms = newTransforms()
	return
}

func (me *NodeDef) NewInst (id string) *NodeInst {
	return newNodeInst(me, id)
}
