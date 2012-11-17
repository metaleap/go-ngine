package assets

import (
)

type Scene struct {
	base
	VisualSceneInst *VisualSceneInst
}

func newScene (id string) (me *Scene) {
	me = &Scene {}
	me.ID = id
	return
}
