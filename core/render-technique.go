package core

import (
	ugl "github.com/go3d/go-opengl/util"
)

type RenderTechnique interface {
	dispose()
	initMeshBuffer(*MeshBuffer) []*ugl.VertexAttribPointer
	onRenderNode()
	program() *ugl.Program
	render()
}

type renderTechniqueBase struct {
	prog *ugl.Program
}

func (me *renderTechniqueBase) dispose() {
}

func (me *renderTechniqueBase) initMeshBuffer(MeshBuffer *MeshBuffer) (atts []*ugl.VertexAttribPointer) {
	return
}

func (me *renderTechniqueBase) onRenderNode() {
}

func (me *renderTechniqueBase) program() *ugl.Program {
	return me.prog
}

func (me *renderTechniqueBase) setProg(name string, unifs []string, attrs []string) {
	prog := glc.progMan.Programs[name]
	if len(unifs) > 0 && len(prog.UnifLocs) == 0 {
		prog.SetUnifLocations(unifs...)
	}
	if len(attrs) > 0 && len(prog.AttrLocs) == 0 {
		prog.SetAttrLocations(attrs...)
	}
	me.prog = prog
}

//	Used only in Core.Rendering.Techniques.
type RenderTechniques map[string]RenderTechnique

func (me *RenderTechniques) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	*me = RenderTechniques{}
}
