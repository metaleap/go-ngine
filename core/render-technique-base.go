package core

type renderTechniqueProvider func(view *RenderView) RenderTechnique

type RenderTechnique interface {
	copyAppToPrep()
	copyPrepToRend()
	dispose()
	name() string
	onPrep()
	render()
}

type renderTechniqueBase struct {
	view  *RenderView
	tname string
}

func (me *renderTechniqueBase) copyAppToPrep() {
}

func (me *renderTechniqueBase) copyPrepToRend() {
}

func (me *renderTechniqueBase) dispose() {
}

func (me *renderTechniqueBase) init(name string, view *RenderView) {
	me.view, me.tname = view, name
}

func (me *renderTechniqueBase) name() string {
	return me.tname
}

func (me *renderTechniqueBase) onPrep() {
}
