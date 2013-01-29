package core

import (
	gl "github.com/chsc/gogl/gl42"
)

func (me *EngineCore) onRender() {
	me.Rendering.Samplers.FullFilteringRepeat.Bind(0)
	for _, curCanvas = range me.Rendering.Canvases {
		if curCanvas.renderThisFrame() {
			curCanvas.render()
		}
	}
	me.Rendering.Samplers.NoFilteringClamp.Bind(0)
	me.Rendering.PostFx.render()
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

func (me *Model) render() {
	me.mesh.render()
}

func (me *Mesh) render() {
	if curMeshBuf != me.meshBuffer {
		me.meshBuffer.use()
	}
	gl.DrawElementsBaseVertex(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Offset(nil, uintptr(me.meshBufOffsetIndices)), gl.Int(me.meshBufOffsetBaseIndex))
	// gl.DrawElements(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Pointer(nil))
}

func (me *PostFx) render() {
	curProg, curMat, curTechnique, curMatKey = nil, nil, nil, ""
	Core.Rendering.states.DisableDepthTest()
	Core.Rendering.states.DisableFaceCulling()
	//Core.Rendering.Samplers.NoFilteringClamp.Bind(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, me.glWidth, me.glHeight)
	// ugl.LogLastError("pre-clrscr")
	gl.Clear(gl.COLOR_BUFFER_BIT)
	// ugl.LogLastError("post-clrscr")
	me.glVao.Bind()
	me.prog.Use()
	mainCanvas.frameBuf.BindTexture(0)
	gl.Uniform1i(me.prog.UnifLocs["uTexRendering"], 0)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}
