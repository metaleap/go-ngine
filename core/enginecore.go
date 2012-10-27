package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

var (
	curCanvIndex, lastCanvIndex int
	curMatKey string
	curCam *TCamera
	curCanvas *TRenderCanvas
	curMat *TMaterial
	curMesh *TMesh
	curNode *TNode
	curProg, tmpProg *glutil.TShaderProgram
	curTechnique iRenderTechnique
	curScene *TScene
)

type tEngineCore struct {
	AssetManager *tAssetManager
	Canvases tRenderCanvases
	CurCamera *TCamera
	CurCanvas *TRenderCanvas
	DefaultCanvasIndex int
	Materials tMaterials
	MeshBuffer *tMeshBuffer
	Meshes tMeshes
	Models tModels
	Options *tOptions
	Scenes tScenes
	Textures tTextures
}

func newEngineCore (options *tOptions) *tEngineCore {
	var defCanvas *TRenderCanvas
	Core = &tEngineCore {}
	Core.Options = options
	Core.AssetManager = newAssetManager()
	Core.Options.DefaultTextureParams.setAgain()
	Core.Meshes = tMeshes {}
	Core.Materials = tMaterials {}
	Core.Textures = tTextures {}
	Core.Scenes = tScenes {}
	Core.Canvases = tRenderCanvases {}
	defCanvas = Core.Canvases.Add(Core.Canvases.New(options.winWidth, options.winHeight, true))
	Core.CurCanvas = defCanvas
	Core.CurCamera = defCanvas.Cameras[0]
	Core.MeshBuffer = newMeshBuffer(1000, 1000)
	return Core
}

func (me *tEngineCore) Dispose () {
	for _, rt := range techs { rt.dispose() }
	techs = map[string]iRenderTechnique {}
	for _, canvas := range me.Canvases { canvas.Dispose() }
	for _, scene := range me.Scenes { scene.Dispose() }
	for _, mesh := range me.Meshes { mesh.GpuDelete() }
	for _, tex := range me.Textures { tex.GpuDelete() }
	me.MeshBuffer.dispose()
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

func (me *tEngineCore) renderLoop () {
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
		glLogLastError("tEngineCore.SyncUpdates(texkey=%s)", key)
	}
	for key, mesh := range me.Meshes {
		if !mesh.gpuSynced { mesh.GpuSync() }
		glLogLastError("tEngineCore.SyncUpdates(meshkey=%s)", key)
	}
}

func (me *tEngineCore) useProgram (name string) {
	if tmpProg = glShaderMan.AllProgs[name]; tmpProg != curProg {
		curProg = tmpProg
		gl.UseProgram(curProg.Program)
	}
}

func (me *tEngineCore) useTechnique (technique iRenderTechnique) {
	if technique != curTechnique {
		curTechnique = technique
		me.useProgram(curTechnique.name())
	}
}
