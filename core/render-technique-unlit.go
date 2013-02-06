package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechniqueUnlit struct {
	renderTechniqueBase
}

func newRenderTechniqueUnlit(progName string) RenderTechnique {
	me := &RenderTechniqueUnlit{}
	me.renderTechniqueBase.setProg(progName, []string{"uMatModelProj", "uDiffuse"}, []string{"aPos", "aTexCoords"})
	return me
}

func (me *RenderTechniqueUnlit) initMeshBuffer(meshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	atts = append(atts,
		ugl.NewVertexAttribPointer("aPos", me.prog.AttrLocs["aPos"], 3, 8*4, gl.Ptr(nil)),
		ugl.NewVertexAttribPointer("aTexCoords", me.prog.AttrLocs["aTexCoords"], 2, 8*4, gl.Util.PtrOffset(nil, 3*4)),
	)
	return
}

func (me *RenderTechniqueUnlit) onRenderNode() {
	if thrRend.tmpMat = thrRend.curNode.EffectiveMaterial(); thrRend.tmpMat != thrRend.curMat {
		if thrRend.curMat = thrRend.tmpMat; thrRend.curMat != nil {
			thrRend.tmpEffect = Core.Libs.Effects[thrRend.curMat.DefaultEffectID]
			Core.Libs.Images.I2D[thrRend.tmpEffect.Diffuse.Texture.Image2ID].glTex.Bind()
			gl.Uniform1i(thrRend.curProg.UnifLocs["uDiffuse"], 0)
		}
	}
}
