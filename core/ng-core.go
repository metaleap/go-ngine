package core

import (
	"math"

	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

var (
	asyncResources         = map[asyncResource]bool{}
	curMeshBuf             *MeshBuffer
	curCanvIndex, curIndex int
	curMatKey, curStr      string
	curCam                 *Camera
	curCanvas              *RenderCanvas
	curMat                 *FxMaterial
	curMesh                *Mesh
	curModel               *Model
	curNode                *Node
	curProg, tmpProg       *ugl.ShaderProgram
	curTechnique           renderTechnique
	curScene               *Scene
)

//	Consider EngineCore a "Singleton" type, only valid use is the core.Core global variable.
//	The heart and brain of go:ngine --- a container for all runtime resources and responsible for rendering.
type EngineCore struct {
	MeshBuffers *MeshBuffers
	Options     EngineOptions
	Libs        struct {
		Cameras   LibCameras
		Effects   LibFxEffects
		Materials LibFxMaterials
		Images    struct {
			I2D LibFxImage2Ds
		}
		Meshes LibMeshes
		Scenes LibScenes
	}
	Rendering struct {
		Canvases RenderCanvases
		PostFx   PostFx
		States   ugl.RenderStates
	}

	isInit bool
	fileIO fileIO
}

func (me *EngineCore) dispose() {
	me.isInit = false
	for _, disp := range []disposable{
		&me.Rendering.Canvases,
		&me.Libs.Cameras, &me.Libs.Images.I2D, &me.Libs.Effects, &me.Libs.Materials, &me.Libs.Meshes, &me.Libs.Scenes,
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
	curCanvas = me.Rendering.Canvases.AddNew(true)
	curCanvas.SetCameraIDs("")
	me.isInit = true
}

func (me *EngineCore) initLibs() {
	libs := &me.Libs
	for _, c := range []ctorable{&libs.Cameras, &libs.Images.I2D, &libs.Effects, &libs.Materials, &libs.Meshes, &libs.Scenes} {
		c.ctor()
	}
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
	for curCanvIndex, curCanvas = range me.Rendering.Canvases {
		if (curCanvas.EveryNthFrame == 1) || ((curCanvas.EveryNthFrame > 1) && (math.Mod(Stats.fpsAll, curCanvas.EveryNthFrame) == 0)) {
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

func (me *EngineCore) onResizeWindow(viewWidth, viewHeight int) {
	if me.isInit {
		me.Options.winWidth, me.Options.winHeight = viewWidth, viewHeight
		for _, canv := range me.Rendering.Canvases {
			if canv.viewSizeRelative {
				canv.absViewWidth, canv.absViewHeight = int(canv.relViewWidth*float64(viewWidth)), int(canv.relViewHeight*float64(viewHeight))
			}
		}
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
	me.onResizeWindow(me.Options.winWidth, me.Options.winHeight)
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
