package core

//	prepares all renderBatches for the next onRender() call...
func (me *EngineCore) onPreRender() {
	for _, curCanvas = range me.Rendering.Canvases {
		if curCanvas.renderThisFrame() {
			curCanvas.preRender()
		}
	}
}

func (me *RenderCanvas) preRender() {
	for _, curCam = range me.Cameras {
		curCam.preRender()
	}
}

func (me *Camera) preRender() {
	if me.Rendering.Enabled {

	}
}

func (me *Node) preRender() {

}
