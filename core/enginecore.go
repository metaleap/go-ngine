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
	AssetManager *tAssetManager
	Canvases []*tRenderCanvas
	CurCamera *tCamera
	CurCanvas *tRenderCanvas
	DefaultCanvasIndex int
	Materials map[string]*TMaterial
	Meshes tMeshMap
	Options *tOptions
	Scenes map[string]*TScene
	Textures map[string]*TTexture
}

func newEngineCore (options *tOptions) *tEngineCore {
	var defCanvas *tRenderCanvas
	Core = &tEngineCore {}
	Core.Options = options
	Core.AssetManager = newAssetManager()
	Core.Options.DefaultTextureParams.setAgain()
	Core.Meshes = tMeshMap {}
	Core.Materials = map [string] *TMaterial {}
	Core.Textures = map[string] *TTexture {}
	Core.Scenes = map [string] *TScene {}
	defCanvas = NewCanvas(options.winWidth, options.winHeight, true)
	Core.Canvases = [] *tRenderCanvas { defCanvas }
	Core.CurCanvas = defCanvas
	Core.CurCamera = defCanvas.Cameras[0]
	return Core
}

func (me *tEngineCore) Dispose () {
	for _, canvas := range me.Canvases { canvas.Dispose() }
	for _, scene := range me.Scenes { scene.Dispose() }
	for _, mesh := range me.Meshes { mesh.Dispose() }
	for _, tex := range me.Textures { tex.GpuDelete() }
	Core = nil
}

func (me *tEngineCore) onSecTick () {
	var allDone = true
	for tex, texDone := range asyncTextures {
		if !texDone {
			if (tex.img != nil) {
				tex.GpuSync()
				asyncTextures[tex] = true
			} else {
				allDone = false
			}
		}
	}
	if allDone { asyncTextures = map[*TTexture]bool {} }
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

func (me *tEngineCore) resizeView (viewWidth, viewHeight int) {
	var defaultCanvas = me.Canvases[me.DefaultCanvasIndex]
	me.Options.winWidth, me.Options.winHeight = viewWidth, viewHeight
	defaultCanvas.viewWidth, defaultCanvas.viewHeight = viewWidth, viewHeight
	for _, cam := range defaultCanvas.Cameras {
		cam.ViewPort.update()
		cam.updatePerspective()
	}
}

func (me *tEngineCore) SyncUpdates () {
	for key, tex := range me.Textures {
		tex.GpuSync()
		nglcore.LogLastError("tEngineCore.SyncUpdates(texkey=%s)", key)
	}
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
