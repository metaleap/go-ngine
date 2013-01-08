package core

type Scenes map[string]*Scene

type Scene struct {
	RootNode *Node
}

func NewScene() (me *Scene) {
	me = &Scene{}
	me.RootNode = newNode("", "", "", nil)
	return
}
