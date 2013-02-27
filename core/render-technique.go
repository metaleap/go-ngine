package core

type RenderTechniqueProvider func() RenderTechnique

type RenderTechnique interface {
	dispose()
	name() string
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
