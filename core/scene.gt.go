package core

//	Represents a scene graph.
type Scene struct {
	//	The root Node for this scene graph.
	RootNode Node
}

//	Makes, Init()s and returns a new Scene with default parameters.
func NewScene() (me *Scene) {
	me = new(Scene)
	me.Init()
	return
}

func (me *Scene) Dispose() {
	Core.Rendering.Canvases.Walk(nil, func(cam *Camera) {
		if cam.scene == me {
			cam.SetScene(nil)
		}
	})
}

func (me *Scene) Init() {
	me.RootNode = *newNode("", "", "", nil, me)
}
