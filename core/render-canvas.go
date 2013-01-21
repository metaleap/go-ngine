package core

//	Represents a surface (texture framebuffer) that can be rendered to.
type RenderCanvas struct {
	//	This MUST be an non-negative integer, it's a float64 just to avoid
	//	type conversion. How often this RenderCanvas is included in rendering:
	//	1 = every frame (this is the default value)
	//	2 = every 2nd frame
	//	8 = every 8th frame
	//	0 = RenderCanvas is disabled for rendering
	EveryNthFrame float64

	isMain, viewSizeRelative    bool
	absViewWidth, absViewHeight int
	relViewWidth, relViewHeight float64
	camIDs                      []string
	cams                        []*Camera
}

func newRenderCanvas() (me *RenderCanvas) {
	me = &RenderCanvas{EveryNthFrame: 1}
	me.SetSize(true, 1, 1)
	return
}

func (me *RenderCanvas) dispose() {
}

func (me *RenderCanvas) Main() bool {
	return me.isMain
}

func (me *RenderCanvas) Remove() {
	sl := Core.Rendering.Canvases
	for i, c := range sl {
		if c == me {
			Core.Rendering.Canvases = append(sl[:i], sl[i+1:]...)
		}
	}
	me.dispose()
}

func (me *RenderCanvas) render() {
	for _, curCam = range me.cams {
		curCam.render()
	}
}

func (me *RenderCanvas) SetCameraIDs(camIDs ...string) {
	me.camIDs, me.cams = camIDs, make([]*Camera, len(camIDs))
	for i, camID := range me.camIDs {
		me.cams[i] = Core.Libs.Cameras[camID]
	}
}

func (me *RenderCanvas) SetMain() {
	for _, canv := range Core.Rendering.Canvases {
		canv.isMain = (canv == me)
	}
}

func (me *RenderCanvas) SetSize(relative bool, width, height float64) {
	if me.viewSizeRelative = relative; me.viewSizeRelative {
		me.relViewWidth, me.relViewHeight = width, height
	} else {
		me.absViewWidth, me.absViewHeight = int(width), int(height)
	}
}

type RenderCanvases []*RenderCanvas

func (me *RenderCanvases) dispose() {
	for _, c := range *me {
		c.dispose()
	}
	*me = RenderCanvases{}
}

func (me *RenderCanvases) AddNew(isMain bool) (rc *RenderCanvas) {
	rc = newRenderCanvas()
	*me = append(*me, rc)
	if (!isMain) && ((len(*me) == 0) || (me.Main() == nil)) {
		isMain = true
	}
	if isMain {
		rc.SetMain()
	}
	return
}

func (me RenderCanvases) Main() (main *RenderCanvas) {
	for _, main = range me {
		if main.isMain {
			return
		}
	}
	main = nil
	return
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
