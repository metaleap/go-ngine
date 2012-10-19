package client

import (
	"log"
	"math"
	"runtime"

	glfw "github.com/jteeuwen/glfw"

	glutil "github.com/go-ngine/go-util/gl"

	ncore "github.com/go-ngine/go-ngine/client/core"
	nglcore "github.com/go-ngine/go-ngine/client/glcore"
)

type TEngineLoop struct {
	IsLooping bool
	SecTickLast, TickNow, TickLast, TickDelta float64
	OnLoopHandlers []func ()
	OnSecTick func ()
}

func (me *TEngineLoop) Loop () {
	var onLoopHandler func ()
	me.SecTickLast, me.TickNow = glfw.Time(), glfw.Time()
	Stats.fps = 0
	if (!me.IsLooping) {
		me.IsLooping = true
		nglcore.LogLastError("ngine.PreLoop")
		log.Printf("Enter loop...")
		for me.IsLooping {
			Core.RenderLoop()
			Stats.fps++
			Stats.fpsAll++
			me.TickLast = me.TickNow
			me.TickNow = glfw.Time()
			me.TickDelta = me.TickNow - me.TickLast
			Core.Timer.LastTick, Core.Timer.NowTick, Core.Timer.TickDelta = me.TickLast, me.TickNow, me.TickDelta
			if math.Floor(me.TickNow) != me.SecTickLast {
				runtime.GC()
				if Stats.TrackGC {
					if Stats.GcTime = glfw.Time() - me.TickNow; Stats.GcTime > Stats.GcMaxTime { Stats.GcMaxTime = Stats.GcTime }
				}
				Stats.fpsSecs++
				Stats.Fps, Stats.fps, me.SecTickLast = Stats.fps, 0, math.Floor(me.TickNow)
				me.OnSecTick()
			}
			for _, onLoopHandler = range me.OnLoopHandlers {
				onLoopHandler()
			}
			Windowing.glfwOnLoop()
		}
		nglcore.LogLastError("ngine.PostLoop")
	}
}

type TEngineStats struct {
	Fps int
	GcTime, GcMaxTime float64
	TrackGC bool
	fps, fpsAll, fpsSecs int
}

func (me *TEngineStats) FpsAvg () int {
	return me.fpsAll / me.fpsSecs
}

var (
	AssetRootDirPath = "."
	Loop *TEngineLoop
	Core *ncore.TEngineCore
	Stats *TEngineStats
	Windowing = newWindowing()
)

func Dispose () {
	nglcore.LogLastError("ngine.Pre_Dispose")
	if Core != nil { Core.Dispose() }
	nglcore.LogLastError("ngine.Core_Dispose")
	nglcore.Dispose()
	Windowing.glfwExit()
	Core, Loop, Stats = nil, nil, nil
}

func Init (winWidth, winHeight int, winFullScreen bool, vsync int, assetRootDirPath, winTitle string, onSecTick func ()) error {
	var err error
	if err = Windowing.glfwInit(winWidth, winHeight, winFullScreen, vsync, winTitle); err == nil {
		if err = nglcore.Init(); err == nil {
			AssetRootDirPath, Loop, Stats = assetRootDirPath, &TEngineLoop {}, &TEngineStats {}
			Core = ncore.NewEngineCore(Windowing.WinWidth, Windowing.WinHeight)
			Loop.OnSecTick = onSecTick
			Loop.OnLoopHandlers = [] func () {}
			log.Println(glutil.GlConnInfo())
		}
	}
	return err
}
