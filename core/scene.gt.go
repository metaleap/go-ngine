package core

//	Represents a scene graph.
type Scene struct {
	//	The root Node for this scene graph.
	RootNode Node
}

func (me *Scene) dispose() {
	Core.Rendering.Canvases.Walk(nil, func(cam *Camera) {
		if cam.scene == me {
			cam.setScene(nil)
		}
	})
}

func (me *Scene) init() {
	me.RootNode = *newNode("", "", "", nil, me)
}

func (me LibScenes) Walk(onNode func(*Node)) {
	for _, scene := range me {
		scene.RootNode.Walk(onNode)
	}
}

//#begin-gt -gen-lib.gt T:Scene

//	Initializes and returns a new Scene with default parameters.
func NewScene() (me *Scene) {
	me = &Scene{}
	me.init()
	return
}

//	A hash-table of Scenes associated by IDs. Only for use in Core.Libs.
type LibScenes map[string]*Scene

//	Creates and initializes a new Scene with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibScenes) AddNew(id string) (obj *Scene) {
	obj = NewScene()
	me[id] = obj
	return
}

func (me *LibScenes) ctor() {
	*me = LibScenes{}
}

func (me *LibScenes) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibScenes) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
