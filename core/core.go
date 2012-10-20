package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"

	nglcore "github.com/go3d/go-ngine/core/glcore"
)

var (
	curCanvIndex, lastCanvIndex int
	curCam *tCamera
	curCanvas *tRenderCanvas
	curMesh *TMesh
	curNode *TNode
	curProg, tmpProg *glutil.TShaderProgram
	curTechnique iRenderTechnique
	curScene *TScene
)

type tEngineCore struct {
	Canvases []*tRenderCanvas
	CurCamera *tCamera
	CurCanvas *tRenderCanvas
	DefaultCanvasIndex int
	Materials map[string]*TMaterial
	Meshes map[string]*TMesh
	Options *tOptions
	Scenes map[string]*TScene
}

func newEngineCore (viewWidth, viewHeight int) *tEngineCore {
	var defCanvas *tRenderCanvas
	Core = &tEngineCore {}
	Core.Options = NewOptions(true)
	Core.Meshes = map [string] *TMesh {}
	Core.Materials = map [string] *TMaterial {}
	Core.Scenes = map [string] *TScene {}
	defCanvas = NewCanvas(viewWidth, viewHeight, true)
	Core.Canvases = [] *tRenderCanvas { defCanvas }
	Core.CurCanvas = defCanvas
	Core.CurCamera = defCanvas.Cameras[0]
	return Core
}

func (me *tEngineCore) Dispose () {
	for _, canvas := range me.Canvases { canvas.Dispose() }
	for _, scene := range me.Scenes { scene.Dispose() }
	for _, mesh := range me.Meshes { mesh.Dispose() }
	for _, mat := range me.Materials { mat.Dispose() }
	Core = nil
}

func (me *tEngineCore) RenderLoop () {
	lastCanvIndex = len(me.Canvases) - 1
	for curCanvIndex, curCanvas = range me.Canvases {
		if !curCanvas.Disabled {
			me.CurCanvas = curCanvas
			curCanvas.render()
		}
	}
}

func (me *tEngineCore) ResizeView (viewWidth, viewHeight int) {
	var defaultCanvas = me.Canvases[me.DefaultCanvasIndex]
	defaultCanvas.viewWidth, defaultCanvas.viewHeight = viewWidth, viewHeight
	for _, cam := range defaultCanvas.Cameras {
		cam.ViewPort.update()
		cam.updatePerspective()
	}
}

func (me *tEngineCore) SyncUpdates () {
	for key, mesh := range me.Meshes {
		if !mesh.glInit { mesh.initBuffer() }
		if !mesh.glSynced { mesh.updateBuffer() }
		nglcore.LogLastError("tEngineCore.SyncUpdates(meshkey=%s)", key)
	}
}

func (me *tEngineCore) useProgram (name string) {
	if tmpProg = nglcore.ShaderMan.AllProgs[name]; tmpProg != curProg {
		curProg = tmpProg
		gl.UseProgram(curProg.Program)
	}
}

func (me *tEngineCore) useTechnique (technique iRenderTechnique) {
	if technique != curTechnique {
		curTechnique = technique
		me.useProgram(curTechnique.Name())
	}
}
