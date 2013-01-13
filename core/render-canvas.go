package core

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

type RenderCanvases []*RenderCanvas

func (me RenderCanvases) New(viewWidth, viewHeight int) (rc *RenderCanvas) {
	rc = newRenderCanvas(viewWidth, viewHeight)
	return
}

func (me *RenderCanvases) dispose() {
	for _, c := range *me {
		c.Dispose()
	}
	*me = RenderCanvases{}
}

func (me *RenderCanvases) Add(canvas *RenderCanvas) *RenderCanvas {
	*me = append(*me, canvas)
	return canvas
}

func (me RenderCanvases) Index(canvas *RenderCanvas) int {
	for curCanvIndex, curCanvas = range me {
		if curCanvas == canvas {
			return curCanvIndex
		}
	}
	return -1
}

func (me *RenderCanvases) Remove(canvas *RenderCanvas) {
	if curCanvIndex = me.Index(canvas); curCanvIndex >= 0 {
		*me = append((*me)[:curCanvIndex], (*me)[curCanvIndex+1:]...)
	}
}

//	Represents a surface that can be rendered to, ie. either the
//	screen or a render texture. Equivalent to a "framebuffer object"
//	in OpenGL.
type RenderCanvas struct {
	//	If true, this RenderCanvas is ignored by the rendering runtime.
	Disabled bool

	camIDs                []string
	cams                  []*Camera
	viewWidth, viewHeight int
}

func newRenderCanvas(viewWidth, viewHeight int) (me *RenderCanvas) {
	me = &RenderCanvas{}
	me.viewWidth, me.viewHeight = viewWidth, viewHeight
	return
}

func (me *RenderCanvas) Dispose() {
}

func (me *RenderCanvas) render() {
	for _, curCam = range me.cams {
		if !curCam.Disabled {
			curCam.render()
		}
	}
}

func (me *RenderCanvas) SetCameraIDs(camIDs ...string) {
	me.camIDs, me.cams = camIDs, make([]*Camera, len(camIDs))
	for i, camID := range me.camIDs {
		me.cams[i] = Core.Libs.Cameras[camID]
	}
}
