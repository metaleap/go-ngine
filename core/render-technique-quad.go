package core

import (
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueQuad struct {
	Effect FxEffect
	renderTechniqueBase
	glVao ugl.VertexArray
}

func newRenderTechniqueQuad(progName string) RenderTechnique {
	me := &RenderTechniqueQuad{}
	me.init("Quad")
	me.Effect.init()
	tex2d := me.Effect.Ops.EnableTex2D(-1)
	tex2d.Sampler = &Core.Rendering.Fx.Samplers.NoFilteringClamp
	tex2d.bindTex = func(_ string) {}
	me.Effect.UpdateRoutine()
	me.glVao.Create()
	me.glVao.Setup(nil)
	return me
}

func (me *RenderTechniqueQuad) dispose() {
	me.glVao.Dispose()
}
