package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

var (
	curMeshBuf *MeshBuffer
	curCanvIndex, lastCanvIndex, curIndex int
	curMatKey, curStr string
	curCam *Camera
	curCanvas *RenderCanvas
	curMat *Material
	curMesh *Mesh
	curModel *Model
	curNode *Node
	curProg, tmpProg *ugl.ShaderProgram
	curTechnique renderTechnique
	curScene *Scene
)

type engineCore struct {
	AssetManager *assetManager
	Canvases renderCanvases
	DefaultCanvasIndex int
	Materials materials
	MeshBuffers *meshBuffers
	Meshes meshes
	Options *engineOptions
	Scenes scenes
	Textures textures
}

func newEngineCore (options *engineOptions) *engineCore {
	initTechniques()
	Core = &engineCore {}
	Core.Options = options
	Core.AssetManager = newAssetManager()
	Core.Options.DefaultTextureParams.setAgain()
	Core.Meshes = meshes {}
	Core.Materials = materials {}
	Core.Textures = textures {}
	Core.Scenes = scenes {}
	Core.Canvases = renderCanvases {}
	curCanvas = Core.Canvases.Add(Core.Canvases.New(options.winWidth, options.winHeight, true))
	curCam = curCanvas.Cameras[0]
	Core.MeshBuffers = newMeshBuffers()
	return Core
}

func (me *engineCore) Dispose () {
	for _, canvas := range me.Canvases { canvas.Dispose() }
	for _, mesh := range me.Meshes { mesh.GpuDelete() }
	for _, tex := range me.Textures { tex.GpuDelete() }
	me.MeshBuffers.dispose()
	techs, Core = nil, nil
}

func (me *engineCore) onLoop () {
}

func (me *engineCore) onRender () {
	lastCanvIndex = len(me.Canvases) - 1
	for curCanvIndex, curCanvas = range me.Canvases {
		if !curCanvas.Disabled {
			curCanvas.render()
		}
	}
}

func (me *engineCore) onSecTick () {
	for tex, _ := range asyncTextures {
		if (tex.img != nil) {
			tex.GpuSync()
			delete(asyncTextures, tex)
			break
		}
	}
}

func (me *engineCore) resizeView (viewWidth, viewHeight int) {
	var defaultCanvas = me.Canvases[me.DefaultCanvasIndex]
	me.Options.winWidth, me.Options.winHeight = viewWidth, viewHeight
	defaultCanvas.viewWidth, defaultCanvas.viewHeight = viewWidth, viewHeight
	for _, cam := range defaultCanvas.Cameras {
		cam.ViewPort.update()
		cam.updatePerspective()
	}
}

func (me *engineCore) SyncUpdates () {
	var err error
	for key, tex := range me.Textures {
		if !tex.gpuSynced {
			tex.GpuSync()
			glLogLastError("engineCore.SyncUpdates(texkey=%s)", key)
		}
	}
	for _, mesh := range me.Meshes {
		if !mesh.gpuSynced {
			if err = mesh.GpuUpload(); err != nil { LogError(err) }
		}
	}
	glLogLastError("engineCore.SyncUpdates()")
	return
}

func (me *engineCore) useProgram (name string) {
	if tmpProg = glShaderMan.AllProgs[name]; tmpProg != curProg {
		curProg = tmpProg
		gl.UseProgram(curProg.Program)
	}
}

func (me *engineCore) useTechnique (technique renderTechnique) {
	if technique != curTechnique {
		curMeshBuf = nil
		curTechnique = technique
		me.useProgram(curTechnique.name())
	}
}
