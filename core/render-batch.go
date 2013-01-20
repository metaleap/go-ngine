package core

type RenderBatch struct {
	scene *Scene
}

func newRenderBatch(scene *Scene) (me *RenderBatch) {
	me = &RenderBatch{scene: scene}
	return
}

func (me *RenderBatch) Update() {

}
