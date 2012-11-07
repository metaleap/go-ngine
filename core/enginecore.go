package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

var (
	curMeshBuf *TMeshBuffer
	curCanvIndex, lastCanvIndex, curIndex int
	curMatKey, curStr string
	curCam *TCamera
	curCanvas *TRenderCanvas
	curMat *TMaterial
	curMesh *TMesh
	curModel *TModel
	curNode *TNode
	curProg, tmpProg *ugl.ShaderProgram
	curTechnique iRenderTechnique
	curScene *TScene
)

type tEngineCore struct {
	AssetManager *tAssetManager
	Canvases tRenderCanvases
	DefaultCanvasIndex int
	Materials tMaterials
	MeshBuffers *tMeshBuffers
	Meshes tMeshes
	Options *tOptions
	Scenes tScenes
	Textures tTextures
}

func newEngineCore (options *tOptions) *tEngineCore {
	initTechniques()
	Core = &tEngineCore {}
	Core.Options = options
	Core.AssetManager = newAssetManager()
	Core.Options.DefaultTextureParams.setAgain()
	Core.Meshes = tMeshes {}
	Core.Materials = tMaterials {}
	Core.Textures = tTextures {}
	Core.Scenes = tScenes {}
	Core.Canvases = tRenderCanvases {}
	curCanvas = Core.Canvases.Add(Core.Canvases.New(options.winWidth, options.winHeight, true))
	curCam = curCanvas.Cameras[0]
	Core.MeshBuffers = newMeshBuffers()
	return Core
}

func (me *tEngineCore) Dispose () {
	for _, canvas := range me.Canvases { canvas.Dispose() }
	for _, mesh := range me.Meshes { mesh.GpuDelete() }
	for _, tex := range me.Textures { tex.GpuDelete() }
	me.MeshBuffers.dispose()
	techs, Core = nil, nil
}

func (me *tEngineCore) onLoop () {
}

func (me *tEngineCore) onRender () {
	lastCanvIndex = len(me.Canvases) - 1
	for curCanvIndex, curCanvas = range me.Canvases {
		if !curCanvas.Disabled {
			curCanvas.render()
		}
	}
}

func (me *tEngineCore) onSecTick () {
	for tex, _ := range asyncTextures {
		if (tex.img != nil) {
			tex.GpuSync()
			delete(asyncTextures, tex)
			break
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
	var err error
	for key, tex := range me.Textures {
		if !tex.gpuSynced {
			tex.GpuSync()
			glLogLastError("tEngineCore.SyncUpdates(texkey=%s)", key)
		}
	}
	for _, mesh := range me.Meshes {
		if !mesh.gpuSynced {
			if err = mesh.GpuUpload(); err != nil { LogError(err) }
		}
	}
	glLogLastError("tEngineCore.SyncUpdates()")
	return
}

func (me *tEngineCore) useProgram (name string) {
	if tmpProg = glShaderMan.AllProgs[name]; tmpProg != curProg {
		curProg = tmpProg
		gl.UseProgram(curProg.Program)
	}
}

func (me *tEngineCore) useTechnique (technique iRenderTechnique) {
	if technique != curTechnique {
		curMeshBuf = nil
		curTechnique = technique
		me.useProgram(curTechnique.name())
	}
}
