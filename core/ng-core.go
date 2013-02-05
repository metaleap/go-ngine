package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

var (
	//	The heart and brain of go:ngine --- a container for all runtime resources and responsible for rendering.
	Core EngineCore
)

//	EngineCore is a singleton type, only used for the core.Core package-global exported variable.
//	It is only aware of that instance and does not support any other EngineCore instances.
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

func (_ *EngineCore) dispose() {
	Core.isInit = false
	for _, disp := range []disposable{
		&Core.Rendering.Canvases,
		&Core.Libs.Images.I2D, &Core.Libs.Effects, &Core.Libs.Materials, &Core.Libs.Meshes, &Core.Libs.Scenes,
		Core.MeshBuffers,
	} {
		disp.dispose()
	}
	Core.Rendering.PostFx.dispose()
	Core.Rendering.Samplers.FullFilteringRepeat.Dispose()
	Core.Rendering.Samplers.NoFilteringClamp.Dispose()
	techs = nil
}

func (_ *EngineCore) init() {
	initTechniques()
	Core.initRenderingStates()
	Core.MeshBuffers = newMeshBuffers()
	Core.initLibs()
	Core.Rendering.Canvases = RenderCanvases{}
	Core.Rendering.PostFx.init()
	Core.Rendering.Canvases.AddNew(true, true, 1, 1).AddNewCamera3D()
	Core.isInit = true
}

func (_ *EngineCore) initLibs() {
	libs := &Core.Libs
	for _, c := range []ctorable{&libs.Images.I2D, &libs.Effects, &libs.Materials, &libs.Meshes, &libs.Scenes} {
		c.ctor()
	}
}

func (_ *EngineCore) initRenderingStates() {
	Core.Rendering.states.ForceClearColor(Core.Options.Rendering.DefaultClearColor)
	Core.Rendering.Samplers.FullFilteringRepeat.Create().EnableFullFiltering(true, 8).SetWrap(gl.REPEAT)
	Core.Rendering.Samplers.NoFilteringClamp.Create().DisableAllFiltering(false).SetWrap(gl.CLAMP_TO_BORDER)
}

func (_ *EngineCore) onResizeWindow(viewWidth, viewHeight int) {
	if Core.isInit {
		Core.Options.winWidth, Core.Options.winHeight = viewWidth, viewHeight
		Core.Rendering.PostFx.glWidth, Core.Rendering.PostFx.glHeight = gl.Sizei(viewWidth), gl.Sizei(viewHeight)
		for _, canv := range Core.Rendering.Canvases {
			canv.onResize(viewWidth, viewHeight)
		}
		ugl.LogLastError("onResizeWindow")
	}
}

func (_ *EngineCore) onSec() {
	if Diag.LogErrorsDuringLoop {
		ugl.LogLastError("onSec")
	}
	for r, d := range thrRend.asyncResources {
		if d {
			delete(thrRend.asyncResources, r)
			r.onAsyncDone()
		}
	}
}

func (_ *EngineCore) SyncUpdates() {
	var (
		err error
		ok  bool
	)
	ugl.LogLastError("EngineCore.SyncUpdates() -- pre")
	Core.onResizeWindow(Core.Options.winWidth, Core.Options.winHeight)
	ugl.LogLastError("EngineCore.SyncUpdates() -- resizewin")
	for _, img := range Core.Libs.Images.I2D {
		if !img.Loaded() {
			if _, ok = thrRend.asyncResources[img]; !ok {
				img.Load()
			}
		}
	}
	ugl.LogLastError("EngineCore.SyncUpdates() -- imgupload")
	for _, mesh := range Core.Libs.Meshes {
		if !mesh.gpuSynced {
			if err = mesh.GpuUpload(); err != nil {
				Diag.LogErr(err)
			}
		}
	}
	ugl.LogLastError("EngineCore.SyncUpdates() -- meshupload")
	return
}

func (_ *EngineCore) useProgram(name string) {
	if thrRend.tmpProg = glc.progMan.Programs[name]; thrRend.tmpProg != thrRend.curProg {
		thrRend.curProg = thrRend.tmpProg
		thrRend.curProg.Use()
	}
}

func (_ *EngineCore) useTechnique(technique renderTechnique) {
	if technique != thrRend.curTechnique {
		thrRend.curMeshBuf = nil
		thrRend.curTechnique = technique
		Core.useProgram(thrRend.curTechnique.name())
	}
}
