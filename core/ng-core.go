package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
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
	MeshBuffers *MeshBuffers
	Options     EngineOptions
	Libs        struct {
		Cameras LibCameras
		Effects LibFxEffects
		Images  struct {
			I2D LibFxImage2Ds
		}
		Materials Materials
		Meshes    Meshes
		Scenes    Scenes
		Textures  Textures
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
	for _, canvas := range me.Rendering.Canvases {
		canvas.Dispose()
	}
	for _, disp := range []disposable{&me.Libs.Cameras, &me.Libs.Images.I2D, &me.Libs.Effects} {
		disp.dispose()
	}
	for _, mesh := range me.Libs.Meshes {
		mesh.GpuDelete()
	}
	for _, tex := range me.Libs.Textures {
		tex.GpuDelete()
	}
	me.MeshBuffers.dispose()
	techs = nil
}

func (me *EngineCore) init(options *EngineOptions) {
	initTechniques()
	me.initRenderingStates()
	me.Options = *options
	me.Options.DefaultTextureParams.setAgain()
	me.MeshBuffers = newMeshBuffers()
	me.initLibs()
	me.Rendering.Canvases = RenderCanvases{}
	curCanvas = me.Rendering.Canvases.Add(me.Rendering.Canvases.New(options.winWidth, options.winHeight))
	curCanvas.SetCameraIDs("")
	me.isInit = true
}

func (me *EngineCore) initLibs() {
	libs := &me.Libs
	for _, c := range []ctorable{&libs.Cameras, &libs.Images.I2D, &libs.Effects} {
		c.ctor()
	}
	libs.Materials = Materials{}
	libs.Meshes = Meshes{}
	libs.Textures = Textures{}
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
	for tex, done := range asyncTextures {
		if done || (tex.img != nil) {
			delete(asyncTextures, tex)
			tex.GpuSync()
			break
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
	var err error
	for key, tex := range me.Libs.Textures {
		if !tex.Loaded() {
			tex.load()
		}
		if !tex.gpuSynced {
			tex.GpuSync()
			glLogLastError("EngineCore.SyncUpdates(texkey=%s)", key)
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
