package core

import (
	gl "github.com/go3d/go-opengl/core"
)

func (me *EngineCore) onRender() {
	me.Rendering.Samplers.FullFilteringRepeat.Bind(0)
	for _, thrRend.curCanv = range me.Rendering.Canvases {
		if thrRend.curCanv.renderThisFrame() {
			thrRend.curCanv.render()
		}
	}
	me.Rendering.Samplers.NoFilteringClamp.Bind(0)
	me.Rendering.PostFx.render()
}

func (me *RenderCanvas) render() {
	me.frameBuf.Bind()
	for _, thrRend.curCam = range me.Cameras {
		thrRend.curCam.render()
	}
	me.frameBuf.Unbind()
}

func (me *Camera) render() {
	if me.Enabled {
		Core.Rendering.states.Apply(&me.thrRend.states)
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceEnableScissorTest()
		}
		Core.useTechnique(me.thrRend.technique)
		if me.Rendering.Viewport.shouldScissor {
			gl.Scissor(me.Rendering.Viewport.glVpX, me.Rendering.Viewport.glVpY, me.Rendering.Viewport.glVpW, me.Rendering.Viewport.glVpH)
		}
		gl.Viewport(me.Rendering.Viewport.glVpX, me.Rendering.Viewport.glVpY, me.Rendering.Viewport.glVpW, me.Rendering.Viewport.glVpH)
		gl.Clear(me.thrRend.states.Other.ClearBits)
		if thrRend.curScene = me.scene; thrRend.curScene != nil {
			thrRend.curScene.RootNode.render()
		}
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceDisableScissorTest()
		}
	}
}

func (me *Node) render() {
	if me.Enabled {
		if thrRend.curNode = me; me.model != nil {
			thrRend.curTechnique.onRenderNode()
			gl.UniformMatrix4fv(thrRend.curProg.UnifLocs["uMatModelProj"], 1, gl.FALSE, &me.thrRend.matProjs[thrRend.curCam][0])
			me.model.render()
		}
		for me.thrRend.curId, me.thrRend.curSubNode = range me.ChildNodes.M {
			me.thrRend.curSubNode.render()
		}
	}
}

func (me *Model) render() {
	me.mesh.render()
}

func (me *Mesh) render() {
	if thrRend.curMeshBuf != me.meshBuffer {
		me.meshBuffer.use()
	}
	gl.DrawElementsBaseVertex(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(me.meshBufOffsetIndices)), gl.Int(me.meshBufOffsetBaseIndex))
	// gl.DrawElements(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Pointer(nil))
}

func (me *PostFx) render() {
	thrRend.curProg, thrRend.curMat, thrRend.curTechnique, thrRend.curMatId = nil, nil, nil, ""
	Core.Rendering.states.DisableDepthTest()
	Core.Rendering.states.DisableFaceCulling()
	Core.Rendering.Samplers.NoFilteringClamp.Bind(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, me.glWidth, me.glHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	me.glVao.Bind()
	me.prog.Use()
	mainCanvas.frameBuf.BindTexture(0)
	gl.Uniform1i(me.prog.UnifLocs["uTexRendering"], 0)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}
