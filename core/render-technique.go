package core

type renderTechniqueProvider func(cam *Camera) RenderTechnique

type RenderTechnique interface {
	dispose()
	name() string
	render()
}

type renderTechniqueBase struct {
	cam   *Camera
	tname string
}

func (me *renderTechniqueBase) dispose() {
}

func (me *renderTechniqueBase) init(name string, cam *Camera) {
	me.cam, me.tname = cam, name
}

func (me *renderTechniqueBase) name() string {
	return me.tname
}
