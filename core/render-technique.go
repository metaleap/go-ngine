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
	prog.SetUnifLocations(unifs...)
	prog.SetAttrLocations(attrs...)
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
