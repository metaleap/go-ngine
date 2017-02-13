package core

import (
	"sync"

	gl "github.com/metaleap/go-opengl/core"
	ugl "github.com/metaleap/go-opengl/util"
)

type RenderTechniqueScene struct {
	renderTechniqueBase
	numDrawCalls int

	Batch  RenderBatcher
	Camera Camera

	thrPrep struct {
		sync.WaitGroup
		batch renderBatchList
	}
	thrRend struct {
		batch renderBatchList
	}
}

func newRenderTechniqueScene(view *RenderView) RenderTechnique {
	me := &RenderTechniqueScene{}
	view.RenderStates.DepthTest = true
	me.init("Scene", view)
	me.Batch = Options.Rendering.DefaultBatcher
	me.Camera.init()
	me.ApplyCameraPerspective()
	return me
}

//	Applies changes made to the Enabled, FovY, ZNear and/or ZFar parameters in me.Camera.Perspective.
func (me *RenderTechniqueScene) ApplyCameraPerspective() {
	me.Camera.viewportAspectRatio = me.view.Port.aspect
	me.Camera.applyPerspective()
}

func (me *RenderTechniqueScene) NumDrawCalls() int {
	return me.numDrawCalls
}

func (me *RenderTechniqueScene) ToggleBatching() {
	me.Batch.Enabled = !me.Batch.Enabled
}

func (me *RenderTechniqueScene) vertexAttribPointers() (atts []ugl.VertexAttribPointer) {
	atts = make([]ugl.VertexAttribPointer, 2)
	atts[0].Init("att_vec3_Pos", 3, 8*4, ugl.PtrNil)
	atts[1].Init("att_vec2_Tex2D", 2, 8*4, gl.Util.PtrOffset(nil, 3*4))
	return
}
