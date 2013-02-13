package core

type RenderTechnique interface {
	dispose()
	name() string
	onRenderNode()
	render()
}

type renderTechniqueBase struct {
	tname string
}

func (me *renderTechniqueBase) dispose() {
}

func (me *renderTechniqueBase) init(name string) {
	me.tname = name
}

func (me *renderTechniqueBase) name() string {
	return me.tname
}

func (me *renderTechniqueBase) onRenderNode() {
}

//	Used only in Core.Rendering.Techniques.
type RenderTechniques map[string]RenderTechnique

func (me *RenderTechniques) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	*me = RenderTechniques{}
}
