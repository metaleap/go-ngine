package core

type scenes map[string]*Scene

type Scene struct {
	RootNode *Node
}

func NewScene () *Scene {
	var scene = &Scene {}
	scene.RootNode = newNode("", "", "", nil)
	return scene
}
