package assets

//	Declares a complete, self-contained base of a Scene hierarchy or Scene graph. Currently just defined by a
//	Visual Scene, later to be augmented by optional "kinematics scenes" and/or "physics scenes".
type Scene struct {
	//	The Visual Scene associated with this Scene.
	VisualSceneInst *VisualSceneInst
}

func newScene(id string) (me *Scene) {
	me = &Scene{}
	return
}
