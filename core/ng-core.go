package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

var (
	asyncResources                        = map[asyncResource]bool{}
	curMeshBuf                            *MeshBuffer
	curCanvIndex, lastCanvIndex, curIndex int
	curMatKey, curStr                     string
	curCam                                *Camera
	curCanvas                             *RenderCanvas
	curMat                                *FxMaterial
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
	MeshBuffers *MeshBuffers
	Options     EngineOptions
	Libs        struct {
		Cameras     LibCameras
		Effects     LibFxEffects
		FxMaterials LibFxMaterials
		Images      struct {
			I2D LibFxImage2Ds
		}
		Meshes Meshes
		Scenes Scenes
	}
	Rendering struct {
		Canvases           RenderCanvases
		DefaultCanvasIndex int
		States             ugl.RenderStates
	}

	isInit bool
	fileIO fileIO
}

func (me *EngineCore) dispose() {
	me.isInit = false
	for _, disp := range []disposable{
		&me.Rendering.Canvases,
		&me.Libs.Cameras, &me.Libs.Images.I2D, &me.Libs.Effects, &me.Libs.FxMaterials, &me.Libs.Meshes,
		me.MeshBuffers,
	} {
		disp.dispose()
	}
	techs = nil
}

func (me *EngineCore) init(options *EngineOptions) {
	initTechniques()
	me.initRenderingStates()
	me.Options = *options
	// me.Options.DefaultTextureParams.setAgain()
	me.MeshBuffers = newMeshBuffers()
	me.initLibs()
	me.Rendering.Canvases = RenderCanvases{}
	curCanvas = me.Rendering.Canvases.Add(me.Rendering.Canvases.New(options.winWidth, options.winHeight))
	curCanvas.SetCameraIDs("")
	me.isInit = true
}

func (me *EngineCore) initLibs() {
	libs := &me.Libs
	for _, c := range []ctorable{&libs.Cameras, &libs.Images.I2D, &libs.Effects, &libs.FxMaterials} {
		c.ctor()
	}
	libs.Meshes = Meshes{}
	libs.Scenes = Scenes{}
}

func (me *EngineCore) initRenderingStates() {
	rs := &me.Rendering.States
	rs.EnableFaceCulling()
	rs.EnableDepthTest()
	rs.SetClearColor(0, 0, 0, 1)
}

func (me *EngineCore) onLoop() {
}

func (me *EngineCore) onRender() {
	lastCanvIndex = len(me.Rendering.Canvases) - 1
	for curCanvIndex, curCanvas = range me.Rendering.Canvases {
		if !curCanvas.Disabled {
			curCanvas.render()
		}
	}
}

func (me *EngineCore) onSec() {
	for r, d := range asyncResources {
		if d {
			delete(asyncResources, r)
			r.onAsyncDone()
		}
	}
}

func (me *EngineCore) resizeView(viewWidth, viewHeight int) {
	if me.isInit {
		defaultCanvas := me.Rendering.Canvases[me.Rendering.DefaultCanvasIndex]
		me.Options.winWidth, me.Options.winHeight = viewWidth, viewHeight
		defaultCanvas.viewWidth, defaultCanvas.viewHeight = viewWidth, viewHeight
		for _, cam := range me.Libs.Cameras {
			cam.ViewPort.update()
			cam.ApplyMatrices()
		}
	}
}

func (me *EngineCore) SyncUpdates() {
	var (
		err error
		ok  bool
	)
	for _, img := range me.Libs.Images.I2D {
		if !img.Loaded() {
			if _, ok = asyncResources[img]; !ok {
				img.Load()
			}
		}
	}
	for _, mesh := range me.Libs.Meshes {
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
	if tmpProg = glShaderMan.progs[name]; tmpProg != curProg {
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
