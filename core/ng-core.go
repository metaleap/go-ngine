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
	MeshBuffers MeshBufferLib
	Libs        EngineLibs
	Render      struct {
		Canvases RenderCanvasLib
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
	Core.Libs.dispose()
	Core.Render.Fx.Samplers.FullFilteringRepeat.Dispose()
	Core.Render.Fx.Samplers.FullFilteringClamp.Dispose()
	Core.Render.Fx.Samplers.NoFilteringClamp.Dispose()
}

func (_ *EngineCore) init() (err error) {
	Core.Libs.init()
	Core.initRendering()
	err = Core.showSplash()
	Core.isInit = true
	return
}

func (_ *EngineCore) initRendering() {
	rend := &Core.Render
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
	if quadFx := &rend.Canvases.AddNew(true, 1, 1).AddNewCameraQuad().RenderTechniqueQuad().Effect; Options.Initialization.DefaultCanvas.GammaViaShader {
		quadFx.Ops.EnableGamma(-1)
		quadFx.KeepOpsLast = append(quadFx.KeepOpsLast, "Gamma")
		quadFx.UpdateRoutine()
	}
}

func (_ *EngineCore) onResizeWindow(viewWidth, viewHeight int) {
	if Core.isInit {
		for i := 0; i < len(Core.Render.Canvases); i++ {
			Core.Render.Canvases[i].onResize(viewWidth, viewHeight)
		}
		// Diag.LogIfGlErr("onResizeWindow")
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

	num := len(Core.Libs.Images.Tex2D) + len(Core.Libs.Images.TexCube)
	ch, done := make(imgChan, num), 0
	Core.Libs.Images.Tex2D.Walk(func(t2d *FxImage2D) {
		go imgLoad(ch, t2d)
	})
	Core.Libs.Images.TexCube.Walk(func(tcm *FxImageCube) {
		go imgLoad(ch, tcm)
	})

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

func (_ *EngineCore) showSplash() (err error) {
	splash := &Core.Libs.Images.SplashScreen
	splash.init()
	splash.Preprocess.FlipY, splash.Preprocess.ToLinear, splash.Preprocess.ToBgra = false, false, false
	splash.Load()
	if err = splash.GpuSync(); err != nil {
		return
	}
	thrRend.quadTex = splash.glTex.GlHandle
	Core.onRender()
	splash.Unload()
	splash.InitFrom.RawData = nil
	glfw.SwapBuffers()
	return
}

func (_ *EngineCore) useTechFx() {
	if thrRend.curTech != thrRend.nextTech || thrRend.curEffect != thrRend.nextEffect {
		thrRend.curTech = thrRend.nextTech
		if !(len(thrRend.curCam.Rendering.FxOps) == 0 && len(thrRend.nextEffect.OpsX) == 0) {
			thrRend.nextEffect.OpsX = thrRend.curCam.Rendering.FxOps
			thrRend.nextEffect.UpdateRoutine()
		}
		thrRend.curEffect = thrRend.nextEffect
		thrRend.curProg = ogl.uber.ensureProg()
		thrRend.curProg.Use()
		thrRend.curEffect.use()
	}
	return
}
