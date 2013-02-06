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

	Cameras Cameras

	isFinal, viewSizeRelative   bool
	absViewWidth, absViewHeight int
	relViewWidth, relViewHeight float64

	frameBuf ugl.Framebuffer
}

func newRenderCanvas(isFinal, relative bool, width, height float64) (me *RenderCanvas) {
	me = &RenderCanvas{EveryNthFrame: 1, isFinal: isFinal}
	if isFinal {
		relative, width, height = true, 1, 1
	}
	me.SetSize(relative, width, height)
	if isFinal {
		me.frameBuf.GlTarget = gl.FRAMEBUFFER
	} else {
		me.frameBuf.Create(gl.Sizei(Core.Options.winWidth), gl.Sizei(Core.Options.winHeight), false)
		me.frameBuf.AttachRendertexture(ugl.NewFramebufferRendertexture(Core.Options.Rendering.PostFx.TextureRect))
		me.frameBuf.AttachRenderbuffer(ugl.NewFramebufferRenderbuffer())
	}
	ugl.LogLastError("newRenderCanvas(%v x %v)", width, height)
	return
}

func (me *RenderCanvas) CurrentAbsoluteSize() (width, height int) {
	width, height = me.absViewWidth, me.absViewHeight
	return
}

func (me *RenderCanvas) AddNewCamera2D(allowOverlaps bool) (cam *Camera) {
	cam = newCamera2D(me, allowOverlaps)
	me.Cameras = append(me.Cameras, cam)
	return
}

func (me *RenderCanvas) AddNewCamera3D() (cam *Camera) {
	cam = newCamera3D(me)
	me.Cameras = append(me.Cameras, cam)
	return
}

func (me *RenderCanvas) AddNewCameraQuad() (cam *Camera) {
	cam = newCameraQuad(me)
	me.Cameras = append(me.Cameras, cam)
	return
}

func (me *RenderCanvas) dispose() {
	if !me.isFinal {
		me.frameBuf.Dispose()
	}
}

func (me *RenderCanvas) onResize(viewWidth, viewHeight int) {
	if me.viewSizeRelative {
		me.absViewWidth, me.absViewHeight = int(me.relViewWidth*float64(viewWidth)), int(me.relViewHeight*float64(viewHeight))
	}
	if !me.isFinal {
		me.frameBuf.Resize(gl.Sizei(me.absViewWidth), gl.Sizei(me.absViewHeight))
	}
	for _, cam := range me.Cameras {
		cam.Rendering.Viewport.update()
		cam.ApplyMatrices()
	}
}

//	Removes me from Core.Rendering.Canvases and deletes its associated GPU resources.
//	This renders me invalid for further use.
func (me *RenderCanvas) Remove() {
	if !me.isFinal {
		sl := Core.Rendering.Canvases
		for i, c := range sl {
			if c == me {
				Core.Rendering.Canvases = append(sl[:i], sl[i+1:]...)
			}
		}
		me.dispose()
	}
}

func (me *RenderCanvas) renderThisFrame() bool {
	return (me.EveryNthFrame == 1) || ((me.EveryNthFrame > 1) && (math.Mod(Stats.fpsAll, me.EveryNthFrame) == 0))
}

//	Sets the 2 dimensions of this render canvas.
//	If relative is true, width and height are interpreted relative to the resolution of the OpenGL context's default framebuffer, with 1 being 100%.
//	Otherwise, width and height are absolute pixel dimensions.
func (me *RenderCanvas) SetSize(relative bool, width, height float64) {
	if me.viewSizeRelative = relative; me.viewSizeRelative {
		me.relViewWidth, me.relViewHeight = width, height
	} else {
		me.absViewWidth, me.absViewHeight = int(width), int(height)
	}
	me.onResize(Core.Options.winWidth, Core.Options.winHeight)
}

//	Only used for Core.Rendering.Canvases.
type RenderCanvases []*RenderCanvas

func (me *RenderCanvases) dispose() {
	for _, c := range *me {
		c.dispose()
	}
	*me = nil
}

//	Adds a new RenderCanvas and returns it.
//	The relative, width and height values are passed to a call to SetSize().
func (me *RenderCanvases) AddNew(relative bool, width, height float64) (rc *RenderCanvas) {
	rc = newRenderCanvas(false, relative, width, height)
	*me = append(*me, rc)
	return
}

func (me RenderCanvases) Final() *RenderCanvas {
	return me[0]
}

func (me RenderCanvases) Walk(onCanv func(*RenderCanvas), onCam func(*Camera)) {
	for _, canv := range me {
		if onCanv != nil {
			onCanv(canv)
		}
		if onCam != nil {
			for _, cam := range canv.Cameras {
				onCam(cam)
			}
		}
	}
}

/*

Canvas 1: lo-res, 2-fx, LDR render
	-	3D geometry CAM					FB1
	-	postfx (TV, B/W)				FB2

Canvas 2: hi-res, many-fx HDR render
	-	3D geometry CAM					FB3		16
	-	HDR postfx (AO, bloom, tonemap)	FB4
	-	LDR postfx (dof, SS, MB, vig)	FB5

	-	2D HUD / gui CAM				FB5
	-	3D mini-map CAM					FB5

Canvas 3: screen
	- SMAA / gamma postfx pass			FB0

*/
