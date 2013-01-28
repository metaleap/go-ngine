package core

import (
	"strings"
	"time"

	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

const postFxProgName = "postfx"

type PostFxEffect struct {
	//	later configure effect params in here
}

//	ONLY used for Core.Rendering.PostFx.
//	Represents the final shading stage in a rendered frame and is always used, but initially all effects are disabled.
//	Always takes the image produced by the Core.Rendering.Canvases.Main render canvas and blits it to the screen,
//	applying the currently enabled effects, if any.
type PostFx struct {
	effects           map[string]*PostFxEffect
	glVao             ugl.VertexArray
	glWidth, glHeight gl.Sizei
	prog              *ugl.Program
}

func (me *PostFx) dispose() {
	me.glVao.Dispose()
}

//	Switches me to a postfx shader program that has all effects enabled as specified
//	by all previous DisableEffect() / EnableEffect() / ToggleEffect() calls since the last ApplyEffects() call.
//	If that shader program does not yet exist, builds it. If that fails, a non-nil err is returned.
func (me *PostFx) ApplyEffects() (err error) {
	var (
		dur   time.Duration
		pname = postFxProgName
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
	if glc.progMan.CloneRawSources(postFxProgName, pname) {
		if dur, err = glc.progMan.MakeProgramsFromRawSources(true, pname); err == nil {
			Diag.LogShaders("Built new shader program '%s' in %v", pname, dur)
		}
	}
	if err == nil {
		me.setProg(pname)
	}
	return
}

//	Deactivates the specified post-processing full-screen effect.
//	After all necessary calls to DisableEffect() / EnableEffect() / ToggleEffect(), be sure to call ApplyEffects() once.
func (me *PostFx) DisableEffect(name string) {
	delete(me.effects, name)
}

//	Activates the specified post-processing full-screen effect.
//	After all necessary calls to DisableEffect() / EnableEffect() / ToggleEffect(), be sure to call ApplyEffects() once.
func (me *PostFx) EnableEffect(name string) (effect *PostFxEffect) {
	if effect = me.effects[name]; effect == nil {
		effect = &PostFxEffect{}
		me.effects[name] = effect
	}
	return
}

func (me *PostFx) init() {
	me.effects = map[string]*PostFxEffect{}
	me.glVao.Create()
	me.setProg(postFxProgName)
}

func (me *PostFx) setProg(progName string) {
	me.prog = glc.progMan.Programs[progName]
	me.prog.SetUnifLocations("uTexRendering")
}

func (me *PostFx) render() {
	curProg, curMat, curTechnique, curMatKey = nil, nil, nil, ""
	Core.Rendering.states.DisableDepthTest()
	Core.Rendering.states.DisableFaceCulling()
	//Core.Rendering.Samplers.NoFilteringClamp.Bind(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, me.glWidth, me.glHeight)
	// ugl.LogLastError("pre-clrscr")
	gl.Clear(gl.COLOR_BUFFER_BIT)
	// ugl.LogLastError("post-clrscr")
	me.glVao.Bind()
	me.prog.Use()
	mainCanvas.frameBuf.BindTexture(0)
	gl.Uniform1i(me.prog.UnifLocs["uTexRendering"], 0)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}

//	Activates or deactivates the specified post-processing full-screen effect.
//	After all necessary calls to DisableEffect() / EnableEffect() / ToggleEffect(), be sure to call ApplyEffects() once.
//	Returns whether this call has activated (true) or deactivated (false) the specified effect.
func (me *PostFx) ToggleEffect(name string) (enabled bool) {
	if effect := me.effects[name]; effect == nil {
		enabled, effect = true, &PostFxEffect{}
		me.effects[name] = effect
	} else {
		delete(me.effects, name)
	}
	return
}
