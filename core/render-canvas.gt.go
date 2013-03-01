package core

import (
	"math"

	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

//	Represents a surface (texture framebuffer) that can be rendered to.
type RenderCanvas struct {
	//	This should be an non-negative integer, it's a float64 just to avoid a
	//	type conversion. How often this RenderCanvas is included in rendering:
	//	1 = every frame (this is the default value)
	//	2 = every 2nd frame
	//	3, 5, 8 etc. = every 3rd, 5th, 8th etc. frame
	//	0 = this RenderCanvas is disabled for rendering
	EveryNthFrame float64

	Cams CameraLib

	Srgb bool

	isRtt, viewSizeRelative     bool
	absViewWidth, absViewHeight float64
	relViewWidth, relViewHeight float64

	frameBuf ugl.Framebuffer
}

func (me *RenderCanvas) init() {
	me.EveryNthFrame, me.isRtt = 1, len(Core.Render.Canvases) > 1
	me.Cams.init()
}

func (me *RenderCanvas) CurrentAbsoluteSize() (width, height int) {
	width, height = int(me.absViewWidth), int(me.absViewHeight)
	return
}

func (me *RenderCanvas) AddNewCamera2D(allowOverlaps bool) (cam *Camera) {
	cam = new(Camera)
	me.Cams.add(cam)
	cam.setup2D(me, allowOverlaps)
	Core.refreshWinSizeRels()
	return
}

func (me *RenderCanvas) AddNewCamera3D() (cam *Camera) {
	cam = new(Camera)
	me.Cams.add(cam)
	cam.setup3D(me)
	Core.refreshWinSizeRels()
	return
}

func (me *RenderCanvas) AddNewCameraQuad() (cam *Camera) {
	cam = new(Camera)
	me.Cams.add(cam)
	cam.setupQuad(me)
	Core.refreshWinSizeRels()
	return
}

func (me *RenderCanvas) dispose() {
	me.Cams.dispose()
	if me.isRtt {
		me.frameBuf.Dispose()
	}
}

func (me *RenderCanvas) onResize(viewWidth, viewHeight int) {
	if me.viewSizeRelative {
		me.absViewWidth, me.absViewHeight = me.relViewWidth*float64(viewWidth), me.relViewHeight*float64(viewHeight)
	}
	if me.isRtt {
		me.frameBuf.Resize(gl.Sizei(int(me.absViewWidth)), gl.Sizei(int(me.absViewHeight)))
	}
	for cam := 0; cam < len(me.Cams); cam++ {
		me.Cams[cam].Rendering.Viewport.canvWidth, me.Cams[cam].Rendering.Viewport.canvHeight = me.absViewWidth, me.absViewHeight
		me.Cams[cam].Rendering.Viewport.update()
		me.Cams[cam].ApplyMatrices()
	}
}

func (me *RenderCanvas) renderThisFrame() bool {
	return me.EveryNthFrame == 1 || (me.EveryNthFrame > 1 && math.Mod(Stats.fpsAll, me.EveryNthFrame) == 0)
}

//	Sets the 2 dimensions of this render canvas.
//	If relative is true, width and height are interpreted relative to the resolution of the OpenGL context's default framebuffer, with 1 being 100%.
//	Otherwise, width and height are absolute pixel dimensions.
func (me *RenderCanvas) SetSize(relative bool, width, height float64) {
	if me.viewSizeRelative = relative; me.viewSizeRelative {
		me.relViewWidth, me.relViewHeight = width, height
	} else {
		me.absViewWidth, me.absViewHeight = width, height
	}
	me.onResize(UserIO.Window.width, UserIO.Window.height)
}

func (me *RenderCanvasLib) AddNew(relative bool, width, height float64) (canv *RenderCanvas) {
	canv = new(RenderCanvas)
	me.add(canv)
	canv.SetSize(relative, width, height)
	if !canv.isRtt {
		canv.frameBuf.GlTarget, canv.Srgb = gl.FRAMEBUFFER, !Options.Initialization.DefaultCanvas.GammaViaShader
	} else {
		canv.frameBuf.Create(gl.Sizei(UserIO.Window.width), gl.Sizei(UserIO.Window.height), false)
		canv.frameBuf.AttachRendertexture()
		canv.frameBuf.AttachRenderbuffer()
	}
	Diag.LogIfGlErr("newRenderCanvas()")
	Core.refreshWinSizeRels()
	return
}

//#begin-gt -gen-reflib.gt T:RenderCanvas L:Core.Render.Canvases

//	Only used for Core.Render.Canvases
type RenderCanvasLib []*RenderCanvas

func (me *RenderCanvasLib) add(ref *RenderCanvas) {
	*me = append(*me, ref)
	ref.init()
	return
}

func (me *RenderCanvasLib) init() {
	*me = make(RenderCanvasLib, 0, 4)
}

func (me *RenderCanvasLib) dispose() {
	me.Remove(0, 0)
}

func (me RenderCanvasLib) Get(id int) (ref *RenderCanvas) {
	if me.IsOk(id) {
		ref = me[id]
	}
	return
}

func (me RenderCanvasLib) IsOk(id int) bool {
	return id > -1 && id < len(me)
}

func (me *RenderCanvasLib) Remove(fromID, num int) {
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

func (me RenderCanvasLib) Walk(on func(ref *RenderCanvas)) {
	for id := 0; id < len(me); id++ {
		on(me[id])
	}
}

//#end-gt
