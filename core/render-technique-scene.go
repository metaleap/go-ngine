package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueScene struct {
	renderTechniqueBase
}

func newRenderTechniqueScene(progName string) RenderTechnique {
	me := &RenderTechniqueScene{}
	me.init("Scene")
	return me
}

func (me *RenderTechniqueScene) initMeshBuffer(meshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	atts = append(atts,
		ugl.NewVertexAttribPointer("att_vec3_Pos", thrRend.curProg.AttrLocs["att_vec3_Pos"], 3, 8*4, gl.Ptr(nil)),
		ugl.NewVertexAttribPointer("att_vec2_Tex2D", thrRend.curProg.AttrLocs["att_vec2_Tex2D"], 2, 8*4, gl.Util.PtrOffset(nil, 3*4)),
	)
	return
}

func (me *RenderTechniqueScene) onRenderNode() {
	if thrRend.tmpMat = thrRend.curNode.EffectiveMaterial(); thrRend.tmpMat != thrRend.curMat {
		if thrRend.curMat = thrRend.tmpMat; thrRend.curMat != nil {
			thrRend.tmpTech, thrRend.tmpEffect = me, Core.Libs.Effects[thrRend.curMat.DefaultEffectID]
			Core.useTechFx()
			thrRend.curEffect.use()
			gl.UniformMatrix4fv(thrRend.curProg.UnifLocs["uni_mat4_VertexMatrix"], 1, gl.FALSE, &thrRend.curNode.thrRend.matProjs[thrRend.curCam][0])
		}
	}
}
