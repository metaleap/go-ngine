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

type TCanvas struct {
	Cameras  []*TCamera
	Disabled bool

	viewWidth, viewHeight int
}

func NewCanvas(viewWidth, viewHeight int, addCam bool) *TCanvas {
	var canvas = &TCanvas{}
	canvas.viewWidth, canvas.viewHeight = viewWidth, viewHeight
	if addCam {
		canvas.Cameras = []*TCamera{NewCamera(canvas, Core.Options.DefaultRenderTechnique)}
	} else {
		canvas.Cameras = []*TCamera{}
	}
	return canvas
}

func (me *TCanvas) Dispose() {
	for _, cam := range me.Cameras {
		cam.Dispose()
	}
}

func (me *TCanvas) render() {
	for _, curCam = range me.Cameras {
		if !curCam.Disabled {
			Core.CurCamera = curCam
			curScene = Core.Scenes[curCam.RenderSceneKey]
			curCam.render()
		}
	}
}
