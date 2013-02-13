package core

import (
	gl "github.com/go3d/go-opengl/core"
)

func (me *EngineCore) onRender() {
	for thrRend.curCanvIndex = len(Core.Rendering.Canvases) - 1; thrRend.curCanvIndex >= 0; thrRend.curCanvIndex-- {
		if thrRend.curCanv = Core.Rendering.Canvases[thrRend.curCanvIndex]; thrRend.curCanv.renderThisFrame() {
			thrRend.curCanv.render()
		}
	}
}

func (me *RenderCanvas) render() {
	if !me.isFinal {
		me.frameBuf.Bind()
	}
	Core.Rendering.states.SetFramebufferSrgb(me.Srgb)
	for _, thrRend.curCam = range me.Cameras {
		thrRend.curCam.render()
	}
	Core.Rendering.states.SetFramebufferSrgb(false)
	if !me.isFinal {
		me.frameBuf.Unbind()
		me.frameBuf.BindTexture(0)
	}
}

func (me *Camera) render() {
	if me.Enabled {
		thrRend.curTechnique = me.Rendering.Technique
		Core.Rendering.states.Apply(&me.thrRend.states)
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceEnableScissorTest()
		}
		if me.Rendering.Viewport.shouldScissor {
			gl.Scissor(me.Rendering.Viewport.glVpX, me.Rendering.Viewport.glVpY, me.Rendering.Viewport.glVpW, me.Rendering.Viewport.glVpH)
		}
		gl.Viewport(me.Rendering.Viewport.glVpX, me.Rendering.Viewport.glVpY, me.Rendering.Viewport.glVpW, me.Rendering.Viewport.glVpH)
		gl.Clear(me.thrRend.states.Other.ClearBits)
		thrRend.curScene = me.scene
		me.Rendering.Technique.render()
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceDisableScissorTest()
		}
	}
}

func (me *RenderTechniqueQuad) render() {
	thrRend.curMat, thrRend.curMatId = nil, ""
	Core.Rendering.Fx.Samplers.NoFilteringClamp.Bind(0)
	me.glVao.Bind()
	Core.useTechFx(me, &me.DefaultEffect)
	gl.Uniform1i(thrRend.curProg.UnifLocs["uni_sampler2D_Tex2D"], 0)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}

func (me *RenderTechniqueScene) render() {
	thrRend.curMat, thrRend.curMatId = nil, ""
	Core.Rendering.Fx.Samplers.FullFilteringRepeat.Bind(0)
	if thrRend.curScene != nil {
		thrRend.curScene.RootNode.render()
	}
}

func (me *Node) render() {
	if me.Enabled {
		if thrRend.curNode = me; me.model != nil {
			thrRend.curTechnique.onRenderNode()
			gl.UniformMatrix4fv(thrRend.curProg.UnifLocs["uni_mat4_VertexMatrix"], 1, gl.FALSE, &me.thrRend.matProjs[thrRend.curCam][0])
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
