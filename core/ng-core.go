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
	Libs        struct {
		Effects   LibFxEffects
		Materials LibFxMaterials
		Images    struct {
			SplashScreen FxImage2D
			TexCubes     LibFxImageCubes
			Tex2D        LibFxImage2Ds
		}
		Meshes LibMeshes

		Scenes SceneLib
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

			procFuncs map[string]string
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
		&Core.Libs.Materials, &Core.Libs.Effects,
		&Core.Libs.Images.Tex2D, &Core.Libs.Images.TexCubes,
		&Core.Libs.Meshes, Core.MeshBuffers, &Core.Libs.Scenes,
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
	Core.showSplash()
	Core.isInit = true
}

func (_ *EngineCore) initLibs() {
	libs := &Core.Libs
	for _, c := range []ctorable{
		&libs.Images.Tex2D, &libs.Images.TexCubes,
		&libs.Effects, &libs.Materials, &libs.Meshes, &libs.Scenes,
	} {
		c.ctor()
	}
}

func (_ *EngineCore) initRendering() {
	rend := &Core.Rendering
	rend.KnownTechniques = map[string]RenderTechniqueProvider{
		"Quad":  newRenderTechniqueQuad,
		"Scene": newRenderTechniqueScene,
	}
	rend.Fx.KnownProcIDs = fxKnownProcIDs()
	rend.Fx.procFuncs = make(map[string]string, len(rend.Fx.KnownProcIDs))
	for _, shaderFunc := range rend.Fx.KnownProcIDs {
		rend.Fx.procFuncs[shaderFunc] = strf("fx_%s", shaderFunc)
	}

	rend.states.ForceClearColor(Options.Rendering.DefaultClearColor)
	rend.Fx.Samplers.FullFilteringRepeat.Create().EnableFullFiltering(true, 8).SetWrap(gl.REPEAT)
	rend.Fx.Samplers.FullFilteringClamp.Create().EnableFullFiltering(true, 8).SetWrap(gl.CLAMP_TO_EDGE)
	rend.Fx.Samplers.NoFilteringClamp.Create().DisableAllFiltering(false).SetWrap(gl.CLAMP_TO_BORDER)
	rend.Canvases = append(RenderCanvases{}, newRenderCanvas(true, true, 1, 1))
	if quadFx := &rend.Canvases.Final().AddNewCameraQuad().RenderTechniqueQuad().Effect; Options.Initialization.DefaultCanvas.GammaViaShader {
		quadFx.Ops.EnableGamma(-1)
		quadFx.KeepOpsLast = append(quadFx.KeepOpsLast, "Gamma")
		quadFx.UpdateRoutine()
	}
}

func (_ *EngineCore) onResizeWindow(viewWidth, viewHeight int) {
	if Core.isInit {
		for _, canv := range Core.Rendering.Canvases {
			canv.onResize(viewWidth, viewHeight)
		}
		Diag.LogIfGlErr("onResizeWindow")
	}
}

func (_ *EngineCore) GpuSyncImageLibs() (err error) {
	type imgChan chan FxImage
	imgLoad := func(ch imgChan, img FxImage) {
		if !img.Loaded() {
			if err := img.Load(); err != nil {
				Diag.LogErr(err)
			}
		}
		ch <- img
	}

	num := len(Core.Libs.Images.Tex2D) + len(Core.Libs.Images.TexCubes)
	ch, done := make(imgChan, num), 0
	for _, t2d := range Core.Libs.Images.Tex2D {
		go imgLoad(ch, t2d)
	}
	for _, tcm := range Core.Libs.Images.TexCubes {
		go imgLoad(ch, tcm)
	}

	for img := range ch {
		//	As soon as the first image is processed/loaded, it can be uploaded while others are still busy
		if done++; done >= num {
			close(ch)
		}
		if img.Loaded() {
			if err = img.GpuSync(); err != nil {
				Diag.LogErr(err)
			}
		}
	}
	return
}

func (_ *EngineCore) refreshWinSizeRels() {
	Core.onResizeWindow(UserIO.Window.width, UserIO.Window.height)
}

func (_ *EngineCore) showSplash() {
	splash := &Core.Libs.Images.SplashScreen
	splash.init()
	splash.PreProcess.FlipY, splash.PreProcess.ToLinear, splash.PreProcess.ToBgra = false, false, false
	splash.Load()
	splash.GpuSync()
	thrRend.quadTex = &splash.glTex
	Core.onRender()
	splash.Unload()
	Core.Libs.Images.SplashScreen.InitFrom.RawData = Core.Libs.Images.SplashScreen.InitFrom.RawData[:0]
	glfw.SwapBuffers()
}

func (_ *EngineCore) useProg(prog *ugl.Program) {
	if thrRend.curProg != prog {
		thrRend.curProg = prog
		thrRend.curProg.Use()
	}
}

func (_ *EngineCore) useSampler(sampler *ugl.Sampler, unit gl.Uint) {
	if thrRend.curSampler[unit] != sampler {
		thrRend.curSampler[unit] = sampler
		sampler.Bind(unit)
	}
}

func (_ *EngineCore) useTechFx() {
	if thrRend.curTech != thrRend.nextTech || thrRend.curEffect != thrRend.nextEffect {
		thrRend.curMeshBuf = nil
		thrRend.curTech = thrRend.nextTech
		if !(len(thrRend.curCam.Rendering.FxOps) == 0 && len(thrRend.nextEffect.OpsX) == 0) {
			thrRend.nextEffect.OpsX = thrRend.curCam.Rendering.FxOps
			thrRend.nextEffect.UpdateRoutine()
		}
		thrRend.curEffect = thrRend.nextEffect
		prog := glc.shaderMan.ensureProg()
		Core.useProg(prog)
		thrRend.curEffect.use()
	}
	return
}
