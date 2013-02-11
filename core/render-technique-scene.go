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
	me.renderTechniqueBase.setProg(progName, []string{"uni_VertexMatrix", "uni_Tex0"}, []string{"att_Pos", "att_Tex0"})
	return me
}

func (me *RenderTechniqueScene) initMeshBuffer(meshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	atts = append(atts,
		ugl.NewVertexAttribPointer("att_Pos", me.prog.AttrLocs["att_Pos"], 3, 8*4, gl.Ptr(nil)),
		ugl.NewVertexAttribPointer("att_Tex0", me.prog.AttrLocs["att_Tex0"], 2, 8*4, gl.Util.PtrOffset(nil, 3*4)),
	)
	return
}

func (me *RenderTechniqueScene) onRenderNode() {
	if thrRend.tmpMat = thrRend.curNode.EffectiveMaterial(); thrRend.tmpMat != thrRend.curMat {
		if thrRend.curMat = thrRend.tmpMat; thrRend.curMat != nil {
			thrRend.tmpEffect = Core.Libs.Effects[thrRend.curMat.DefaultEffectID]
			Core.Libs.Images.I2D[thrRend.tmpEffect.OldDiffuse.Texture.Image2ID].glTex.Bind()
			gl.Uniform1i(thrRend.curProg.UnifLocs["uni_Tex0"], 0)
		}
	}
}
