package core

import (
	"math"

	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

//	Represents a surface (texture framebuffer) that can be rendered to.
type RenderCanvas struct {
	ID int

	//	This should be an non-negative integer, it's a float64 just to avoid a
	//	type conversion. How often this RenderCanvas is included in rendering:
	//	1 = every frame (this is the default value)
	//	2 = every 2nd frame
	//	3, 5, 8 etc. = every 3rd, 5th, 8th etc. frame
	//	0 = this RenderCanvas is disabled for rendering
	EveryNthFrame float64

	Cameras CameraLib

	Srgb bool

	isRtt, viewSizeRelative     bool
	absViewWidth, absViewHeight float64
	relViewWidth, relViewHeight float64

	frameBuf ugl.Framebuffer
}

func (me *RenderCanvas) init() {
	me.EveryNthFrame, me.isRtt, me.Cameras = 1, me.ID > 0, make(CameraLib, 0, 4)
	me.SetSize(true, 1, 1)
	if !me.isRtt {
		me.frameBuf.GlTarget, me.Srgb = gl.FRAMEBUFFER, !Options.Initialization.DefaultCanvas.GammaViaShader
	} else {
		me.frameBuf.Create(gl.Sizei(UserIO.Window.width), gl.Sizei(UserIO.Window.height), false)
		me.frameBuf.AttachRendertexture()
		me.frameBuf.AttachRenderbuffer()
	}
	Diag.LogIfGlErr("newRenderCanvas()")
	Core.refreshWinSizeRels()
}

func (me *RenderCanvas) CurrentAbsoluteSize() (width, height int) {
	width, height = int(me.absViewWidth), int(me.absViewHeight)
	return
}

func (me *RenderCanvas) AddNewCamera2D(allowOverlaps bool) (cam *Camera) {
	cam = me.Cameras.AddNew()
	cam.setup2D(me, allowOverlaps)
	Core.refreshWinSizeRels()
	return
}

func (me *RenderCanvas) AddNewCamera3D() (cam *Camera) {
	cam = me.Cameras.AddNew()
	cam.setup3D(me)
	Core.refreshWinSizeRels()
	return
}

func (me *RenderCanvas) AddNewCameraQuad() (cam *Camera) {
	cam = me.Cameras.AddNew()
	cam.setupQuad(me)
	Core.refreshWinSizeRels()
	return
}

func (me *RenderCanvas) dispose() {
	me.Cameras.dispose()
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
	for cam := 0; cam < len(me.Cameras); cam++ {
		if me.Cameras.Ok(cam) {
			me.Cameras[cam].Rendering.Viewport.canvWidth, me.Cameras[cam].Rendering.Viewport.canvHeight = me.absViewWidth, me.absViewHeight
			me.Cameras[cam].Rendering.Viewport.update()
			me.Cameras[cam].ApplyMatrices()
		}
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

//#begin-gt -gen-lib.gt T:RenderCanvas L:Core.Rendering.Canvases

//	Only used for Core.Rendering.Canvases
type RenderCanvasLib []RenderCanvas

func (me *RenderCanvasLib) AddNew() (ref *RenderCanvas) {
	id := -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			id = i
			break
		}
	}
	if id < 0 {
		if id = len(*me); id == cap(*me) {
			nu := make(RenderCanvasLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, RenderCanvas{})
	}
	ref = &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *RenderCanvasLib) Compact() {
	var (
		before, after []RenderCanvas
		ref           *RenderCanvas
		oldID, i      int
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			before, after = (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	changed := make(map[int]int, len(*me))
	for i = 0; i < len(*me); i++ {
		if ref = &(*me)[i]; ref.ID != i {
			oldID, ref.ID = ref.ID, i
			changed[oldID] = i
		}
	}
	if len(changed) > 0 {
		me.onRenderCanvasIDsChanged(changed)
	}
}

func (me *RenderCanvasLib) ctor() {
	*me = make(RenderCanvasLib, 0, Options.Libs.InitialCap)
}

func (me *RenderCanvasLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me RenderCanvasLib) Get(id int) (ref *RenderCanvas) {
	if id > -1 && id < len(me) {
		if ref = &me[id]; ref.ID != id {
			ref = nil
		}
	}
	return
}

func (me RenderCanvasLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me RenderCanvasLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me RenderCanvasLib) Remove(fromID, num int) {
	if l := len(me); fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onRenderCanvasIDsChanged(changed)
	}
}

func (me RenderCanvasLib) Walk(on func(ref *RenderCanvas)) {
	for id := 0; id < len(me); id++ {
		if me[id].ID > -1 {
			on(&me[id])
		}
	}
}

//#end-gt
