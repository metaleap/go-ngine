package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueScene struct {
	renderTechniqueBase
}

func newRenderTechniqueScene() RenderTechnique {
	me := &RenderTechniqueScene{}
	me.init("Scene")
	return me
}

func (me *RenderTechniqueScene) initMeshBuffer(meshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	atts = append(atts,
		ugl.NewVertexAttribPointer("att_vec3_Pos", 3, 8*4, gl.Ptr(nil)),
		ugl.NewVertexAttribPointer("att_vec2_Tex2D", 2, 8*4, gl.Util.PtrOffset(nil, 3*4)),
	)
	return
}
