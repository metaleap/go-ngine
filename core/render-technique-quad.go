package core

import (
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueQuad struct {
	renderTechniqueBase
	glVao ugl.VertexArray
}

func newRenderTechniqueQuad(progName string) RenderTechnique {
	me := &RenderTechniqueQuad{}
	me.init("Quad")
	me.glVao.Create()
	me.glVao.Setup(nil)
	return me
}

func (me *RenderTechniqueQuad) dispose() {
	me.glVao.Dispose()
}
