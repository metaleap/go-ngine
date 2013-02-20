package core

import (
	"sync"

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
			TexCubes     LibFxImageCubes
			Tex2D        LibFxImage2Ds
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
				FullFilteringClamp  ugl.Sampler
			}

			procs map[string]*fxProc
		}
		KnownTechniques map[string]RenderTechniqueProvider

		states ugl.RenderStates
	}

	isInit bool
	fileIO fileIO
}

func (_ *EngineCore) dispose() {
	Core.isInit = false
	for _, disp := range []disposable{
		&Core.Rendering.Canvases,
		&Core.Libs.Images.Tex2D, &Core.Libs.Images.TexCubes, &Core.Libs.Effects, &Core.Libs.Materials, &Core.Libs.Meshes, &Core.Libs.Scenes,
		Core.MeshBuffers,
	} {
		disp.dispose()
	}
	Core.Rendering.Fx.Samplers.FullFilteringRepeat.Dispose()
	Core.Rendering.Fx.Samplers.FullFilteringClamp.Dispose()
	Core.Rendering.Fx.Samplers.NoFilteringClamp.Dispose()
}

func (_ *EngineCore) init() {
	Core.MeshBuffers = newMeshBuffers()
	Core.initLibs()
	Core.initRendering()
	splash := &Core.Libs.Images.SplashScreen
	splash.InitFrom.RawData = embeddedBinaries["splash.png"]
	splash.init()
	splash.PreProcess.FlipY, splash.PreProcess.SrgbToLinear = false, false
	splash.Load()
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
	for _, c := range []ctorable{&libs.Images.Tex2D, &libs.Images.TexCubes, &libs.Effects, &libs.Materials, &libs.Meshes, &libs.Scenes} {
		c.ctor()
	}
}

func (_ *EngineCore) initRendering() {
	rend := &Core.Rendering
	rend.KnownTechniques = map[string]RenderTechniqueProvider{
		"Quad":  newRenderTechniqueQuad,
		"Scene": newRenderTechniqueScene,
	}
	rend.Fx.KnownProcIDs = []string{"Tex2D", "TexCube", "Grayscale", "Orangify", "Colored", "Gamma"}
	rend.Fx.procs = map[string]*fxProc{}
	for _, shaderFunc := range rend.Fx.KnownProcIDs {
		rend.Fx.procs[shaderFunc] = newFxProc(shaderFunc)
	}

	rend.states.ForceClearColor(Core.Options.Rendering.DefaultClearColor)
	rend.Fx.Samplers.FullFilteringRepeat.Create().EnableFullFiltering(true, 8).SetWrap(gl.REPEAT)
	rend.Fx.Samplers.FullFilteringClamp.Create().EnableFullFiltering(true, 8).SetWrap(gl.CLAMP_TO_EDGE)
	rend.Fx.Samplers.NoFilteringClamp.Create().DisableAllFiltering(false).SetWrap(gl.CLAMP_TO_BORDER)
	rend.Canvases = append(RenderCanvases{}, newRenderCanvas(true, true, 1, 1))
	rend.Canvases.Final().AddNewCameraQuad()
}

func (_ *EngineCore) onResizeWindow(viewWidth, viewHeight int) {
	if Core.isInit {
		for _, canv := range Core.Rendering.Canvases {
			canv.onResize(viewWidth, viewHeight)
		}
		Diag.LogIfGlErr("onResizeWindow")
	}
}

func (_ *EngineCore) onSec() {
	if Diag.LogGLErrorsInLoopOnSec {
		Diag.LogIfGlErr("onSec")
	}
}

func (_ *EngineCore) SyncUpdates() (err error) {
	var (
		wait sync.WaitGroup
	)
	imgLoad := func(img FxImage) {
		if !img.Loaded() {
			wait.Add(1)
			go func() {
				if err = img.Load(); err != nil {
					Diag.LogErr(err)
				}
				wait.Done()
			}()
		}
	}
	imgPush := func(img FxImage) {
		if img.Loaded() {
			if err = img.GpuSync(); err != nil {
				Diag.LogErr(err)
			}
		}
	}

	Diag.LogIfGlErr("EngineCore.SyncUpdates() -- pre")
	Core.onResizeWindow(UserIO.Window.width, UserIO.Window.height)
	Diag.LogIfGlErr("EngineCore.SyncUpdates() -- resizewin")

	for _, t2d := range Core.Libs.Images.Tex2D {
		imgLoad(t2d)
	}
	for _, tcm := range Core.Libs.Images.TexCubes {
		imgLoad(tcm)
	}
	wait.Wait()
	for _, t2d := range Core.Libs.Images.Tex2D {
		imgPush(t2d)
	}
	for _, tcm := range Core.Libs.Images.TexCubes {
		imgPush(tcm)
	}
	Diag.LogIfGlErr("EngineCore.SyncUpdates() -- imgupload")

	for _, mesh := range Core.Libs.Meshes {
		if !mesh.gpuSynced {
			if err = mesh.GpuUpload(); err != nil {
				Diag.LogErr(err)
			}
		}
	}
	Diag.LogIfGlErr("EngineCore.SyncUpdates() -- meshupload")
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
