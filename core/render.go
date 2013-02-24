package core

import (
	gl "github.com/go3d/go-opengl/core"
)

func (me *EngineCore) onRender() {
	var canv *RenderCanvas
	for canvIndex := len(Core.Rendering.Canvases) - 1; canvIndex >= 0; canvIndex-- {
		if canv = Core.Rendering.Canvases[canvIndex]; canv.renderThisFrame() {
			canv.render()
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
		thrRend.quadTex = &me.frameBufTex.Texture2D
	}
}

func (me *Camera) render() {
	if me.Enabled {
		thrRend.curTech = me.Rendering.Technique
		Core.Rendering.states.Apply(&me.thrRend.states)
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceEnableScissorTest()
		}
		if me.Rendering.Viewport.shouldScissor {
			gl.Scissor(me.Rendering.Viewport.glVpX, me.Rendering.Viewport.glVpY, me.Rendering.Viewport.glVpW, me.Rendering.Viewport.glVpH)
		}
		gl.Viewport(me.Rendering.Viewport.glVpX, me.Rendering.Viewport.glVpY, me.Rendering.Viewport.glVpW, me.Rendering.Viewport.glVpH)
		gl.Clear(me.thrRend.states.Other.ClearBits)
		me.Rendering.Technique.render()
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceDisableScissorTest()
		}
	}
}

func (me *RenderTechniqueQuad) render() {
	thrRend.curMat = nil
	me.glVao.Bind()
	me.fxTex.glTex = thrRend.quadTex
	thrRend.nextTech, thrRend.nextEffect = me, &me.Effect
	Core.useTechFx()
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}

func (me *RenderTechniqueScene) render() {
	thrRend.curMat, thrRend.nextTech = nil, me
	if me.cam.scene != nil {
		me.cam.scene.RootNode.renderChildren()
		me.cam.scene.RootNode.renderSelf() // might be a skybox so "render" the root last
	}
}

func (me *Node) renderChildren() {
	for _, subNode := range me.ChildNodes.M {
		subNode.renderSelf()
		subNode.renderChildren()
	}
}

func (me *Node) renderSelf() {
	if thrRend.curNode = me; me.model != nil {
		if mat := me.EffectiveMaterial(); mat != thrRend.curMat {
			if thrRend.curMat = mat; thrRend.curMat != nil {
				thrRend.nextEffect = Core.Libs.Effects[thrRend.curMat.DefaultEffectID]
				Core.useTechFx()
			}
		}
		if me.Rendering.skyMode {
			Core.Rendering.states.DisableFaceCulling()
			gl.DepthFunc(gl.LEQUAL)
			thrRend.curProg.Uniform1i("uni_int_Sky", 1)
		}
		thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &me.thrRend.camProjMats[thrRend.curCam][0])
		me.model.render(me)
		if me.Rendering.skyMode {
			thrRend.curProg.Uniform1i("uni_int_Sky", 0)
			gl.DepthFunc(gl.LESS)
		}
	}
}

func (me *Model) render(node *Node) {
	me.mesh.render(node)
}

func (me *Mesh) render(node *Node) {
	if thrRend.curMeshBuf != me.meshBuffer {
		me.meshBuffer.use()
	}
	if thrRend.curMat.HasFaceEffects() {
		for fidx, face := range me.raw.faces {
			thrRend.nextEffect = thrRend.curMat.faceEffect(face)
			Core.useTechFx()
			me.meshBuffer.use()
			thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &node.thrRend.camProjMats[thrRend.curCam][0])
			gl.DrawElementsBaseVertex(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(me.meshBufOffsetIndices+(int32(fidx)*3*4))), gl.Int(me.meshBufOffsetBaseIndex))
		}
	} else {
		gl.DrawElementsBaseVertex(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(me.meshBufOffsetIndices)), gl.Int(me.meshBufOffsetBaseIndex))
	}
}
