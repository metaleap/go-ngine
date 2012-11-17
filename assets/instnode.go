package assets

type NodeInst struct {
	baseInst
	Def *NodeDef
}

func newNodeInst (def *NodeDef, id string) (me *NodeInst) {
	me = &NodeInst { Def: def }
	me.base.init(id)
	return
}
