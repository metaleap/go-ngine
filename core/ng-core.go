package core

import (
	glfw "github.com/go-gl/glfw"
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
			SplashScreen FxImage2D
			I2D          LibFxImage2Ds
		}
		Meshes LibMeshes
		Scenes LibScenes
	}
	Rendering struct {
		Canvases RenderCanvases
		Fx       struct {
			KnownProcIDs []string
			Samplers     struct {
				NoFilteringClamp    ugl.Sampler
				FullFilteringRepeat ugl.Sampler
			}

			procs map[string]*fxProc
		}
		Techniques RenderTechniques

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
		Core.MeshBuffers, &Core.Rendering.Techniques,
	} {
		disp.dispose()
	}
	Core.Rendering.Fx.Samplers.FullFilteringRepeat.Dispose()
	Core.Rendering.Fx.Samplers.NoFilteringClamp.Dispose()
}

func (_ *EngineCore) init() {
	Core.MeshBuffers = newMeshBuffers()
	Core.initLibs()
	initRenderTechniques()
	Core.initRendering()
	splash := &Core.Libs.Images.SplashScreen
	splash.InitFrom.RawData = embeddedBinaries["splash.png"]
	splash.init()
	splash.Load()
	splash.FlipY, splash.ConvertToLinear = false, false
	splash.GpuSync()
	thrRend.tmpQuadTex = &splash.glTex
	splash.Unload()
	embeddedBinaries = nil
	Core.isInit = true
	Core.onRender()
	glfw.SwapBuffers()
}

func (_ *EngineCore) initLibs() {
	libs := &Core.Libs
	for _, c := range []ctorable{&libs.Images.I2D, &libs.Effects, &libs.Materials, &libs.Meshes, &libs.Scenes} {
		c.ctor()
	}
}

func (_ *EngineCore) initRendering() {
	rend := &Core.Rendering
	rend.states.ForceClearColor(Core.Options.Rendering.DefaultClearColor)
	rend.Fx.Samplers.FullFilteringRepeat.Create().EnableFullFiltering(true, 8).SetWrap(gl.REPEAT)
	rend.Fx.Samplers.NoFilteringClamp.Create().DisableAllFiltering(false).SetWrap(gl.CLAMP_TO_BORDER)
	rend.Canvases = append(RenderCanvases{}, newRenderCanvas(true, true, 1, 1))
	rend.Canvases.Final().AddNewCameraQuad()
}

func initRenderTechniques() {
	type techCtor func(string) RenderTechnique
	rend, knownTechniques := &Core.Rendering, map[string]techCtor{
		"Quad":  newRenderTechniqueQuad,
		"Scene": newRenderTechniqueScene,
	}
	rend.Techniques = RenderTechniques{}
	for name, _ := range knownTechniques {
		//	the fxprocs need the names, but the techs mustn't be created yet as their ctor needs the fxprocs......
		rend.Techniques[name] = nil
	}
	rend.Fx.KnownProcIDs = []string{"Tex2D", "Grayscale", "Orangify", "Colored", "Gamma"}
	rend.Fx.procs = map[string]*fxProc{}
	for _, shaderFunc := range rend.Fx.KnownProcIDs {
		rend.Fx.procs[shaderFunc] = newFxProc(shaderFunc)
	}
	for name, ctor := range knownTechniques {
		rend.Techniques[name] = ctor(name)
	}
}

func (_ *EngineCore) onResizeWindow(viewWidth, viewHeight int) {
	if Core.isInit {
		for _, canv := range Core.Rendering.Canvases {
			canv.onResize(viewWidth, viewHeight)
		}
		ugl.LogLastError("onResizeWindow")
	}
}

func (_ *EngineCore) onSec() {
	if Diag.LogGLErrorsInLoopOnSec {
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
	Core.onResizeWindow(UserIO.Window.width, UserIO.Window.height)
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

func (_ *EngineCore) useProg() {
	if thrRend.curProg != thrRend.tmpProg {
		thrRend.curProg = thrRend.tmpProg
		thrRend.curProg.Use()
	}
}

func (_ *EngineCore) useSampler() {
	if thrRend.curSampler[thrRend.tmpTexUnit] != thrRend.tmpSampler {
		thrRend.curSampler[thrRend.tmpTexUnit] = thrRend.tmpSampler
		thrRend.tmpSampler.Bind(thrRend.tmpTexUnit)
	}
}

func (_ *EngineCore) useTechFx() {
	if thrRend.curTech != thrRend.tmpTech || thrRend.curEffect != thrRend.tmpEffect {
		thrRend.curMeshBuf = nil
		thrRend.curTech = thrRend.tmpTech
		thrRend.curEffect = thrRend.tmpEffect
		glc.shaderMan.ensureProg()
		Core.useProg()
	}
	return
}
