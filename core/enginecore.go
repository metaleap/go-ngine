package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
	nga "github.com/go3d/go-ngine/assets"
)

var (
	curMeshBuf                            *MeshBuffer
	curCanvIndex, lastCanvIndex, curIndex int
	curMatKey, curStr                     string
	curCam                                *Camera
	curCanvas                             *RenderCanvas
	curMat                                *Material
	curMesh                               *Mesh
	curModel                              *Model
	curNode                               *Node
	curProg, tmpProg                      *ugl.ShaderProgram
	curTechnique                          renderTechnique
	curScene                              *Scene
)

//	Consider EngineCore a "Singleton" type, only valid use is the core.Core global variable.
//	The heart and brain of go:ngine --- a container for all runtime resources and responsible for rendering.
type EngineCore struct {
	AssetManager       *assetManager
	Cameras            cameras
	Canvases           renderCanvases
	DefaultCanvasIndex int
	Materials          materials
	MeshBuffers        *meshBuffers
	Meshes             meshes
	Options            *EngineOptions
	Scenes             scenes
	Textures           textures
}

func newEngineCore(options *EngineOptions) {
	initTechniques()
	Core = &EngineCore{}
	Core.Options = options
	Core.AssetManager = newAssetManager()
	Core.Options.DefaultTextureParams.setAgain()
	Core.Materials = materials{}
	Core.Meshes = meshes{}
	Core.Textures = textures{}
	Core.Scenes = scenes{}
	Core.Cameras = cameras{}
	Core.Canvases = renderCanvases{}
	curCanvas = Core.Canvases.Add(Core.Canvases.New(options.winWidth, options.winHeight))
	curCanvas.SetCameraIDs("")
	Core.MeshBuffers = newMeshBuffers()

	nga.OnBeforeSyncAll = func() { Core.onAssetsSyncing() }
	nga.OnAfterSyncAll = func() { Core.onAssetsSynced() }
	nga.CameraDefs.OnSync = func() { Core.Cameras.syncAssetChanges() }
	nga.ImageDefs.OnSync = func() { Core.Textures.syncAssetChanges() }
}

func (me *EngineCore) dispose() {
	for _, cam := range me.Cameras {
		cam.dispose()
	}
	for _, canvas := range me.Canvases {
		canvas.Dispose()
	}
	for _, mesh := range me.Meshes {
		mesh.GpuDelete()
	}
	for _, tex := range me.Textures {
		tex.GpuDelete()
	}
	me.MeshBuffers.dispose()
	techs, Core = nil, nil
}

func (me *EngineCore) onLoop() {
}

func (me *EngineCore) onAssetsSynced() {
	curCanvas.SetCameraIDs("")
}

func (me *EngineCore) onAssetsSyncing() {
}

func (me *EngineCore) onRender() {
	lastCanvIndex = len(me.Canvases) - 1
	for curCanvIndex, curCanvas = range me.Canvases {
		if !curCanvas.Disabled {
			curCanvas.render()
		}
	}
}

func (me *EngineCore) onSec() {
	for tex, done := range asyncTextures {
		if done || (tex.img != nil) {
			delete(asyncTextures, tex)
			tex.GpuSync()
			break
		}
	}
}

func (me *EngineCore) resizeView(viewWidth, viewHeight int) {
	defaultCanvas := me.Canvases[me.DefaultCanvasIndex]
	me.Options.winWidth, me.Options.winHeight = viewWidth, viewHeight
	defaultCanvas.viewWidth, defaultCanvas.viewHeight = viewWidth, viewHeight
	for _, cam := range me.Cameras {
		cam.ViewPort.Update()
		cam.UpdatePerspective()
	}
}

func (me *EngineCore) SyncUpdates() {
	var err error
	for key, tex := range me.Textures {
		if !tex.gpuSynced {
			tex.GpuSync()
			glLogLastError("EngineCore.SyncUpdates(texkey=%s)", key)
		}
	}
	for _, mesh := range me.Meshes {
		if !mesh.gpuSynced {
			if err = mesh.GpuUpload(); err != nil {
				logError(err)
			}
		}
	}
	glLogLastError("EngineCore.SyncUpdates()")
	return
}

func (me *EngineCore) useProgram(name string) {
	if tmpProg = glShaderMan.AllProgs[name]; tmpProg != curProg {
		curProg = tmpProg
		gl.UseProgram(curProg.Program)
	}
}

func (me *EngineCore) useTechnique(technique renderTechnique) {
	if technique != curTechnique {
		curMeshBuf = nil
		curTechnique = technique
		me.useProgram(curTechnique.name())
	}
}
