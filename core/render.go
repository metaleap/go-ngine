package core

import (
	"math"

	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

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

func (me *RenderCanvas) render() {
	me.frameBuf.Bind()
	for _, curCam = range me.Cameras {
		curCam.render()
	}
	me.frameBuf.Unbind()
}

func (me *Camera) render() {
	if me.Rendering.Enabled {
		Core.Rendering.states.Apply(&me.Rendering.States.RenderStatesBag)
		Core.Rendering.states.ForceEnableScissorTest()
		curScene = Core.Libs.Scenes[me.Rendering.SceneID]
		Core.useTechnique(me.technique)
		me.matCamProj.SetFromMult4(&me.matProj, &me.Controller.mat)
		//me.glMatCamProj.Load(&me.matCamProj)
		//gl.UniformMatrix4fv(curProg.UnifLocs["uMatCamProj"], 1, gl.FALSE, &me.glMatCamProj[0])
		// gl.UniformMatrix4fv(curProg.UnifLocs["uMatProj"], 1, gl.FALSE, &me.glmatProj[0])
		me.technique.onPreRender()
		gl.Scissor(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		gl.Viewport(me.Rendering.ViewPort.glVpX, me.Rendering.ViewPort.glVpY, me.Rendering.ViewPort.glVpW, me.Rendering.ViewPort.glVpH)
		if me.Rendering.States.ClearColor[3] > 0 {
			gl.Clear(me.Rendering.States.ClearFlags)
		}
		curScene.RootNode.render()
		Core.Rendering.states.ForceDisableScissorTest()
	}
}

func (me *Node) render() {
	if !me.Disabled {
		if curNode, curMesh, curModel = me, me.mesh, me.model; curModel != nil {
			curTechnique.onRenderNode()
			if curCam.Perspective.Use {
				me.matModelProj.SetFromMult4(&curCam.matCamProj, &me.Transform.matModelView)
			} else {
				me.matModelProj = me.Transform.matModelView
			}
			me.glMatModelProj.Load(&me.matModelProj)
			gl.UniformMatrix4fv(curProg.UnifLocs["uMatModelProj"], 1, gl.FALSE, &me.glMatModelProj[0])
			me.model.render()
		}
		for me.curKey, me.curSubNode = range me.ChildNodes.M {
			me.curSubNode.render()
		}
	}
}
