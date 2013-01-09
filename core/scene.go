package core

//	A hash-table of Scenes associated with their IDs.
//	Only used for Core.Libs.Scenes.
type Scenes map[string]*Scene

//	Represents a scene graph.
type Scene struct {
	//	The root Node for this scene graph.
	RootNode Node
}

//	Initializes and returns a new, empty Scene.
func NewScene() (me *Scene) {
	me = &Scene{}
	me.RootNode = *newNode("", "", "", nil)
	return
}
