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

type renderCanvases []*RenderCanvas

	func (me renderCanvases) New (viewWidth, viewHeight int) (rc *RenderCanvas) {
		rc = newRenderCanvas(viewWidth, viewHeight)
		return
	}

	func (me *renderCanvases) Add (canvas *RenderCanvas) *RenderCanvas {
		*me = append(*me, canvas)
		return canvas
	}

	func (me renderCanvases) Index (canvas *RenderCanvas) int {
		for curCanvIndex, curCanvas = range me { if curCanvas == canvas { return curCanvIndex } }
		return -1
	}

	func (me *renderCanvases) Remove (canvas *RenderCanvas) {
		if curCanvIndex = me.Index(canvas); curCanvIndex >= 0 {
			*me = append((*me)[ : curCanvIndex], (*me)[curCanvIndex + 1 : ]...)
		}
	}

type RenderCanvas struct {
	Disabled bool

	camIDs []string
	cams []*Camera
	viewWidth, viewHeight int
}

func newRenderCanvas (viewWidth, viewHeight int) (me *RenderCanvas) {
	me = &RenderCanvas {}
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

func (me *RenderCanvas) SetCameraIDs (camIDs ... string) {
	me.camIDs, me.cams = camIDs, make([]*Camera, len(camIDs))
	for i, camID := range me.camIDs {
		me.cams[i] = Core.Cameras[camID]
	}
}
