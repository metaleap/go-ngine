package core

type TCanvas struct {
	Cameras []*TCamera
	Disabled bool

	viewWidth, viewHeight int
}

func NewCanvas (viewWidth, viewHeight int, addCam bool) *TCanvas {
	var canvas = &TCanvas {}
	canvas.viewWidth, canvas.viewHeight = viewWidth, viewHeight
	if addCam {
		canvas.Cameras = []*TCamera { NewCamera(canvas, "geometry") }
	} else {
		canvas.Cameras = []*TCamera {}
	}
	return canvas
}

func (me *TCanvas) Dispose () {
	for _, cam := range me.Cameras { cam.Dispose() }
}

func (me *TCanvas) render () {
	for _, curCam = range me.Cameras {
		if !curCam.Disabled {
			Core.CurCamera = curCam
			curScene = Core.Scenes[curCam.RenderSceneKey]
			curCam.render()
		}
	}
}
