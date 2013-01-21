package core

/*

render sorting per pass:

- by VAO (first 3D, then 2D etc)
- by material shader
- by texture (uniform values)
- by other setunif
- by mesh (so multiple mesh-insts with identical material/texture are rendered together)

*/

type RenderBatch struct {
	scene *Scene
}

func newRenderBatch(scene *Scene) (me *RenderBatch) {
	me = &RenderBatch{scene: scene}
	return
}

func (me *RenderBatch) Update() {

}
