package core

import (
	"math"

	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

var (
	//	The heart and brain of go:ngine --- a container for all runtime resources and responsible for rendering.
	Core EngineCore

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
	curProg, tmpProg       *ugl.Program
	curTechnique           renderTechnique
	curScene               *Scene
)

//	Consider EngineCore a "Singleton" type, only valid use is the core.Core global variable.
//	The heart and brain of go:ngine --- a container for all runtime resources and responsible for rendering.
type EngineCore struct {
	MeshBuffers *MeshBuffers
	Options     EngineOptions
	Libs        struct {
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
		Samplers struct {
			NoFilteringClamp    ugl.Sampler
			FullFilteringRepeat ugl.Sampler
		}
		states ugl.RenderStates
	}

	isInit bool
	fileIO fileIO
}

func (me *EngineCore) dispose() {
	me.isInit = false
	for _, disp := range []disposable{
		&me.Rendering.Canvases,
		&me.Libs.Images.I2D, &me.Libs.Effects, &me.Libs.Materials, &me.Libs.Meshes, &me.Libs.Scenes,
		me.MeshBuffers,
	} {
		disp.dispose()
	}
	me.Rendering.PostFx.dispose()
	me.Rendering.Samplers.FullFilteringRepeat.Dispose()
	me.Rendering.Samplers.NoFilteringClamp.Dispose()
	techs = nil
}

func (me *EngineCore) init() {
	initTechniques()
	me.initRenderingStates()
	me.MeshBuffers = newMeshBuffers()
	me.initLibs()
	me.Rendering.Canvases = RenderCanvases{}
	me.Rendering.PostFx.init()
	curCanvas = me.Rendering.Canvases.AddNew(true, true, 1, 1)
	curCam = curCanvas.Cameras.Add(NewCamera3D())
	me.isInit = true
}

func (me *EngineCore) initLibs() {
	libs := &me.Libs
	for _, c := range []ctorable{&libs.Images.I2D, &libs.Effects, &libs.Materials, &libs.Meshes, &libs.Scenes} {
		c.ctor()
	}
}

func (me *EngineCore) initRenderingStates() {
	me.Rendering.states.ForceClearColor(me.Options.Rendering.DefaultClearColor)
	me.Rendering.Samplers.FullFilteringRepeat.Create().EnableFullFiltering(true, 8).SetWrap(gl.REPEAT)
	me.Rendering.Samplers.NoFilteringClamp.Create().DisableAllFiltering(false).SetWrap(gl.CLAMP_TO_BORDER)
}

func (me *EngineCore) onLoop() {
}

func (me *EngineCore) onRender() {
	me.Rendering.Samplers.FullFilteringRepeat.Bind(0)
	for curCanvIndex, curCanvas = range me.Rendering.Canvases {
		if (curCanvas.EveryNthFrame == 1) || ((curCanvas.EveryNthFrame > 1) && (math.Mod(Stats.fpsAll, curCanvas.EveryNthFrame) == 0)) {
			curCanvas.render()
		}
	}
	me.Rendering.Samplers.NoFilteringClamp.Bind(0)
	me.Rendering.PostFx.render()
	ugl.LogLastError("onrender")
}

//	prepares all renderBatches for the next onRender() call...
func (me *EngineCore) onPreRender() {

}

func (me *EngineCore) onResizeWindow(viewWidth, viewHeight int) {
	if me.isInit {
		me.Options.winWidth, me.Options.winHeight = viewWidth, viewHeight
		me.Rendering.PostFx.glWidth, me.Rendering.PostFx.glHeight = gl.Sizei(viewWidth), gl.Sizei(viewHeight)
		for _, canv := range me.Rendering.Canvases {
			canv.onResize(viewWidth, viewHeight)
		}
		ugl.LogLastError("onResizeWindow")
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

func (me *EngineCore) SyncUpdates() {
	var (
		err error
		ok  bool
	)
	ugl.LogLastError("EngineCore.SyncUpdates() -- pre")
	me.onResizeWindow(me.Options.winWidth, me.Options.winHeight)
	ugl.LogLastError("EngineCore.SyncUpdates() -- resizewin")
	for _, img := range me.Libs.Images.I2D {
		if !img.Loaded() {
			if _, ok = asyncResources[img]; !ok {
				img.Load()
			}
		}
	}
	ugl.LogLastError("EngineCore.SyncUpdates() -- imgupload")
	for _, mesh := range me.Libs.Meshes {
		if !mesh.gpuSynced {
			if err = mesh.GpuUpload(); err != nil {
				Diag.LogErr(err)
			}
		}
	}
	ugl.LogLastError("EngineCore.SyncUpdates() -- meshupload")
	return
}

func (me *EngineCore) useProgram(name string) {
	if tmpProg = glc.progMan.Programs[name]; tmpProg != curProg {
		curProg = tmpProg
		curProg.Use()
	}
}

func (me *EngineCore) useTechnique(technique renderTechnique) {
	if technique != curTechnique {
		curMeshBuf = nil
		curTechnique = technique
		me.useProgram(curTechnique.name())
	}
}
