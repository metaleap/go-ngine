package core

type tScenes map[string]*TScene

type TScene struct {
	RootNode *TNode
}

func (me *TScene) Dispose () {
	if me.RootNode != nil {
		me.RootNode.Dispose()
	}
}

func NewScene () *TScene {
	var scene = &TScene {}
	scene.RootNode = newNode("", nil)
	return scene
}
