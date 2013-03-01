package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueScene struct {
	renderTechniqueBase
	// Batch RenderBatch
}

func newRenderTechniqueScene() RenderTechnique {
	me := &RenderTechniqueScene{}
	me.init("Scene")
	// me.Batch.init(me)
	// me.Batch.Enabled = true
	return me
}

func (me *RenderTechniqueScene) ToggleBatching() {
	// me.Batch.Enabled = !me.Batch.Enabled
}

func (me *RenderTechniqueScene) vertexAttribPointers() (atts []ugl.VertexAttribPointer) {
	atts = make([]ugl.VertexAttribPointer, 2)
	atts[0].Init("att_vec3_Pos", 3, 8*4, ugl.PtrNil)
	atts[1].Init("att_vec2_Tex2D", 2, 8*4, gl.Util.PtrOffset(nil, 3*4))
	return
}
