package core

import (
	"strings"
	"time"

	ugl "github.com/go3d/go-opengl/util"
)

//	Temporary concoction
type PostFxEffect struct {
	//	later configure effect params in here
}

type RenderTechniqueQuad struct {
	renderTechniqueBase
	effects map[string]*PostFxEffect
	glVao   ugl.VertexArray
}

func newRenderTechniqueQuad(progName string) RenderTechnique {
	me := &RenderTechniqueQuad{effects: map[string]*PostFxEffect{}}
	me.setProg(progName)
	me.glVao.Create()
	me.glVao.Setup(nil)
	return me
}

func (me *RenderTechniqueQuad) dispose() {
	me.glVao.Dispose()
}

//	Switches me to a postfx shader program that has all effects enabled as specified
//	by all previous DisableEffect() / EnableEffect() / ToggleEffect() calls since the last ApplyEffects() call.
//	If that shader program does not yet exist, builds it. If that fails, a non-nil err is returned.
func (me *RenderTechniqueQuad) ApplyEffects() (err error) {
	var (
		dur   time.Duration
		pname = Core.Options.Rendering.DefaultTechniqueQuad
	)
	for defName, _ := range glc.progMan.Defines {
		if strings.HasPrefix(defName, "FX_") {
			delete(glc.progMan.Defines, defName)
		}
	}
	for name, _ := range me.effects {
		glc.progMan.Defines["FX_"+name] = 1
		pname += ("__" + name)
	}
	if glc.progMan.CloneRawSources(Core.Options.Rendering.DefaultTechniqueQuad, pname) {
		if dur, err = glc.progMan.MakeProgramsFromRawSources(true, pname); err == nil {
			Diag.LogShaders("Built new shader program '%s' in %v", pname, dur)
			Stats.addProgCompile(1, dur.Nanoseconds())
		}
	}
	if err == nil {
		me.setProg(pname)
	}
	thrRend.curTechnique, thrRend.curProg = nil, nil
	return
}

//	Deactivates the specified post-processing full-screen effect.
//	After all necessary calls to DisableEffect() / EnableEffect() / ToggleEffect(), be sure to call ApplyEffects() once.
func (me *RenderTechniqueQuad) DisableEffect(name string) {
	delete(me.effects, name)
}

//	Activates the specified post-processing full-screen effect.
//	After all necessary calls to DisableEffect() / EnableEffect() / ToggleEffect(), be sure to call ApplyEffects() once.
func (me *RenderTechniqueQuad) EnableEffect(name string) (effect *PostFxEffect) {
	if effect = me.effects[name]; effect == nil {
		effect = &PostFxEffect{}
		me.effects[name] = effect
	}
	return
}

func (me *RenderTechniqueQuad) setProg(progName string) {
	me.renderTechniqueBase.setProg(progName, []string{"uni_Tex0"}, nil)
}

//	Activates or deactivates the specified post-processing full-screen effect.
//	After all necessary calls to DisableEffect() / EnableEffect() / ToggleEffect(), be sure to call ApplyEffects() once.
//	Returns whether this call has activated (true) or deactivated (false) the specified effect.
func (me *RenderTechniqueQuad) ToggleEffect(name string) (enabled bool) {
	if effect := me.effects[name]; effect == nil {
		enabled, effect = true, &PostFxEffect{}
		me.effects[name] = effect
	} else {
		delete(me.effects, name)
	}
	return
}
