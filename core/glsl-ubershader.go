package core

type shaderComposer struct {
	rawSources map[string]string
}

func (me *shaderComposer) init() {
	me.rawSources = map[string]string{}
}
