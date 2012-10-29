package core

type tScenes map[string]*TScene

type TScene struct {
	RootNode *TNode
}

func NewScene () *TScene {
	var scene = &TScene {}
	scene.RootNode = newNode("", "", "", nil)
	return scene
}
