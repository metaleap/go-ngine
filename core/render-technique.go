package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

var (
	techs map[string]renderTechnique
)

type renderTechnique interface {
	initMeshBuffer(*MeshBuffer) []*ugl.VertexAttribPointer
	name() string
	onRenderNode()
}

type techniqueCtor func(string) renderTechnique

func initTechniques() {
	techs = map[string]renderTechnique{}
	for techName, techMaker := range map[string]techniqueCtor{"rt_quad": newRenderTechniqueQuad, "rt_unlit": newRenderTechniqueUnlit} {
		techs[techName] = techMaker(techName)
	}
}

type renderTechniqueBase struct {
	prog *ugl.Program
}

func (me *renderTechniqueBase) initMeshBuffer(MeshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	return
}

func (me *renderTechniqueBase) onRenderNode() {
}

func (me *renderTechniqueBase) name() string {
	return me.prog.Name
}

func (me *renderTechniqueBase) setProg(name string, unifs []string, attrs []string) {
	prog := glc.progMan.Programs[name]
	if len(unifs) > 0 {
		prog.SetUnifLocations(unifs...)
	}
	if len(attrs) > 0 {
		prog.SetAttrLocations(attrs...)
	}
	me.prog = prog
}

type renderTechniqueQuad struct {
	renderTechniqueBase
}

func newRenderTechniqueQuad(progName string) renderTechnique {
	me := &renderTechniqueQuad{}
	me.renderTechniqueBase.setProg(progName, []string{"uTexRendering"}, nil)
	return me
}

type renderTechniqueUnlit struct {
	renderTechniqueBase
}

func newRenderTechniqueUnlit(progName string) renderTechnique {
	me := &renderTechniqueUnlit{}
	me.renderTechniqueBase.setProg(progName, []string{"uMatModelProj", "uDiffuse"}, []string{"aPos", "aTexCoords"})
	return me
}

func (me *renderTechniqueUnlit) initMeshBuffer(meshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	atts = append(atts,
		ugl.NewVertexAttribPointer("aPos", me.prog.AttrLocs["aPos"], 3, 8*4, gl.Ptr(nil)),
		ugl.NewVertexAttribPointer("aTexCoords", me.prog.AttrLocs["aTexCoords"], 2, 8*4, gl.Util.PtrOffset(nil, 3*4)),
	)
	return
}

func (me *renderTechniqueUnlit) onRenderNode() {
	if thrRend.tmpMat = thrRend.curNode.EffectiveMaterial(); thrRend.tmpMat != thrRend.curMat {
		if thrRend.curMat = thrRend.tmpMat; thrRend.curMat != nil {
			thrRend.tmpEffect = Core.Libs.Effects[thrRend.curMat.DefaultEffectID]
			Core.Libs.Images.I2D[thrRend.tmpEffect.Diffuse.Texture.Image2ID].glTex.Bind()
			gl.Uniform1i(thrRend.curProg.UnifLocs["uDiffuse"], 0)
		}
	}
}
