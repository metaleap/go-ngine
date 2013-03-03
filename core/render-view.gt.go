package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

type RenderView struct {
	Enabled bool

	FxProcs FxProcs

	//	The device-relative or absolute view-port for this view.
	Port RenderViewport

	RenderStates ugl.RenderStatesBag

	Technique RenderTechnique
}

func (me *RenderView) init() {
	me.Enabled = true
	me.RenderStates.FaceCulling, me.RenderStates.StencilTest = false, false
	me.RenderStates.ClearColor = Options.Rendering.DefaultClearColor
}

func (me *RenderView) dispose() {
	me.Technique.dispose()
}

func (me *RenderView) setup(canv *RenderCanvas, technique string) {
	me.Port.canvWidth, me.Port.canvHeight = float64(canv.absViewWidth), float64(canv.absViewHeight)
	me.Port.init()
	me.Technique = Core.Render.KnownTechniques[technique](me)
	if me.RenderStates.DepthTest {
		me.RenderStates.Other.ClearBits = gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT
	} else {
		me.RenderStates.Other.ClearBits = gl.COLOR_BUFFER_BIT
	}
}

func (me *RenderView) RenderTechniqueQuad() (tech *RenderTechniqueQuad) {
	tech, _ = me.Technique.(*RenderTechniqueQuad)
	return
}

func (me *RenderView) RenderTechniqueScene() (tech *RenderTechniqueScene) {
	tech, _ = me.Technique.(*RenderTechniqueScene)
	return
}

//	Only used for RenderView.Port
type RenderViewport struct {
	relative, shouldScissor bool
	relX, relY, relW, relH  float64
	absX, absY, absW, absH  int
	aspect                  float64
	glVpX, glVpY            gl.Int
	glVpW, glVpH            gl.Sizei
	canvWidth, canvHeight   float64
}

func (me *RenderViewport) init() {
	me.SetRelative(0, 0, 1, 1)
}

//	Sets the absolute viewport origin and dimensions in pixels.
func (me *RenderViewport) SetAbsolute(x, y, width, height int) {
	me.relative, me.absX, me.absY, me.absW, me.absH = false, x, y, width, height
	me.update()
}

//	Sets the RenderCanvas-relative viewport origin and dimensions, with the value 1.0
//	representing the maximum extent of the viewport on that respective axis.
func (me *RenderViewport) SetRelative(x, y, width, height float64) {
	me.relative, me.relX, me.relY, me.relW, me.relH = true, x, y, width, height
	me.update()
}

func (me *RenderViewport) update() {
	me.shouldScissor = !(me.relative && me.relX == 0 && me.relY == 0 && me.relW == 1 && me.relH == 1)
	if me.relative {
		me.absW, me.absH = int(me.relW*float64(me.canvWidth)), int(me.relH*float64(me.canvHeight))
		me.absX, me.absY = int(me.relX*float64(me.canvWidth)), int(me.relY*float64(me.canvHeight))
	}
	me.glVpX, me.glVpY, me.glVpW, me.glVpH = gl.Int(me.absX), gl.Int(me.absY), gl.Sizei(me.absW), gl.Sizei(me.absH)
	me.aspect = float64(me.absW) / float64(me.absH)
}

//#begin-gt -gen-reflib.gt T:RenderView L:Core.Render.Canvases[id].Views

//	Only used for Core.Render.Canvases[id].Views
type RenderViewLib []*RenderView

func (me *RenderViewLib) add(ref *RenderView) {
	*me = append(*me, ref)
	ref.init()
	return
}

func (me *RenderViewLib) init() {
	*me = make(RenderViewLib, 0, 4)
}

func (me *RenderViewLib) dispose() {
	me.Remove(0, 0)
}

func (me RenderViewLib) IsOk(id int) bool {
	return id > -1 && id < len(me)
}

func (me *RenderViewLib) Remove(fromID, num int) {
	if l := len(*me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		for i := fromID; i < fromID+num; i++ {
			(*me)[fromID].dispose()
		}
		before, after := (*me)[:fromID], (*me)[fromID+num:]
		*me = append(before, after...)
	}
}

func (me RenderViewLib) Walk(on func(ref *RenderView)) {
	for id := 0; id < len(me); id++ {
		on(me[id])
	}
}

//#end-gt
