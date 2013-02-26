package core

//	A hash-table of child nodes. Only used for Node.ChildNodes.
type Nodes struct {
	//	The underlying hash-table. NOT to be modified directly.
	//	ONLY use the methods defined on the Nodes type to add,
	//	remove or move nodes.
	M map[string]*Node

	owner *Node
}

func (me *Nodes) init(owner *Node) {
	me.owner = owner
	me.M = map[string]*Node{}
}

//	Removes node from its previous parent Node (if any)
//	and adds it to me.M under its ID.
func (me *Nodes) Add(node *Node) {
	if node.parentNode != nil {
		node.parentNode.ChildNodes.Remove(node.id)
	}
	node.parentNode = me.owner
	me.M[node.id] = node
}

//	Creates a new Node with the specified ID, binds it to the
//	specified Mesh and Model, adds it to me.M and returns it.
func (me *Nodes) AddNew(id string, meshID int, modelID string) (node *Node) {
	node = newNode(id, meshID, modelID, me.owner, me.owner.rootScene)
	me.Add(node)
	return
}

//	Removes the Node with the specified ID from me.M.
func (me *Nodes) Remove(id string) {
	if node := me.M[id]; node != nil {
		node.parentNode = nil
	}
	delete(me.M, id)
}
