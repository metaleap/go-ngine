package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueScene struct {
	renderTechniqueBase
	Batch RenderBatch
}

func newRenderTechniqueScene(cam *Camera) RenderTechnique {
	me := &RenderTechniqueScene{}
	me.init("Scene", cam)
	me.Batch.init(me)
	me.Batch.Enabled = true
	return me
}

func (me *RenderTechniqueScene) ToggleBatching() {
	me.Batch.Enabled = !me.Batch.Enabled
}

func (me *RenderTechniqueScene) vertexAttribPointers(meshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	atts = append(atts,
		ugl.NewVertexAttribPointer("att_vec3_Pos", 3, 8*4, ugl.PtrNil),
		ugl.NewVertexAttribPointer("att_vec2_Tex2D", 2, 8*4, gl.Util.PtrOffset(nil, 3*4)),
	)
	return
}
