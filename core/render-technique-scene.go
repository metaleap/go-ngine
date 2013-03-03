package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueScene struct {
	renderTechniqueBase
	// Batch RenderBatch

	Camera Camera
}

func newRenderTechniqueScene(view *RenderView) RenderTechnique {
	me := &RenderTechniqueScene{}
	view.RenderStates.DepthTest = true
	me.init("Scene", view)
	me.Camera.init()
	// me.Batch.init(me)
	// me.Batch.Enabled = true
	me.ApplyCameraProjection()
	return me
}

//	Applies changes made to the Enabled, FovY, ZNear and/or ZFar parameters in me.Camera.Perspective.
func (me *RenderTechniqueScene) ApplyCameraProjection() {
	if me.Camera.Perspective.Enabled {
		me.Camera.thrApp.matProj.Perspective(me.Camera.Perspective.FovY, me.view.Port.aspect, me.Camera.Perspective.ZNear, me.Camera.Perspective.ZFar)
	} else {
		me.Camera.thrApp.matProj.Identity()
	}
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
