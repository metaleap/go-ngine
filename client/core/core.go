package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go-ngine/go-util/gl"

	nglcore "github.com/go-ngine/go-ngine/client/glcore"
)

var (
	Core *TEngineCore

	curCanvIndex, lastCanvIndex int
	curCam *TCamera
	curCanvas *TCanvas
	curMesh *TMesh
	curNode *TNode
	curProg, tmpProg *glutil.TShaderProgram
	curTechnique ITechnique
	curScene *TScene
)

type TEngineCore struct {
	Canvases []*TCanvas
	CurCamera *TCamera
	CurCanvas *TCanvas
	DefaultCanvasIndex int
	Materials map[string]*TMaterial
	Meshes map[string]*TMesh
	Options *TOptions
	Scenes map[string]*TScene
	Timer *TTimer
}

func NewEngineCore (viewWidth, viewHeight int) *TEngineCore {
	var defCanvas *TCanvas
	Core = &TEngineCore {}
	Core.Options = NewOptions(true)
	Core.Meshes = map [string] *TMesh {}
	Core.Materials = map [string] *TMaterial {}
	Core.Scenes = map [string] *TScene {}
	Core.Timer = &TTimer {}
	defCanvas = NewCanvas(viewWidth, viewHeight, true)
	Core.Canvases = [] *TCanvas { defCanvas }
	Core.CurCanvas = defCanvas
	Core.CurCamera = defCanvas.Cameras[0]
	return Core
}

func (me *TEngineCore) Dispose () {
	for _, canvas := range me.Canvases { canvas.Dispose() }
	for _, scene := range me.Scenes { scene.Dispose() }
	for _, mesh := range me.Meshes { mesh.Dispose() }
	for _, mat := range me.Materials { mat.Dispose() }
	Core = nil
}

func (me *TEngineCore) RenderLoop () {
	lastCanvIndex = len(me.Canvases) - 1
	for curCanvIndex, curCanvas = range me.Canvases {
		if !curCanvas.Disabled {
			me.CurCanvas = curCanvas
			curCanvas.render()
		}
	}
}

func (me *TEngineCore) ResizeView (viewWidth, viewHeight int) {
	var defaultCanvas = me.Canvases[me.DefaultCanvasIndex]
	defaultCanvas.viewWidth, defaultCanvas.viewHeight = viewWidth, viewHeight
	for _, cam := range defaultCanvas.Cameras {
		cam.ViewPort.update()
		cam.updatePerspective()
	}
}

func (me *TEngineCore) SyncUpdates () {
	for key, mesh := range me.Meshes {
		if !mesh.glInit { mesh.initBuffer() }
		if !mesh.glSynced { mesh.updateBuffer() }
		nglcore.LogLastError("TEngineCore.SyncUpdates(meshkey=%s)", key)
	}
}

func (me *TEngineCore) useProgram (name string) {
	if tmpProg = nglcore.ShaderMan.AllProgs[name]; tmpProg != curProg {
		curProg = tmpProg
		gl.UseProgram(curProg.Program)
	}
}

func (me *TEngineCore) useTechnique (technique ITechnique) {
	if technique != curTechnique {
		curTechnique = technique
		me.useProgram(curTechnique.Name())
	}
}
