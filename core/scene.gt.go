package core

//	Represents a scene graph.
type Scene struct {
	//	The root Node for this scene graph.
	RootNode Node
}

func (me *Scene) dispose() {
}

func (me *Scene) init() {
	me.RootNode = *newNode("", "", "", nil)
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

//#end-gt
