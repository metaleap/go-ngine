package core

import (
	gl "github.com/go3d/go-opengl/core"
)

func (_ *NgCore) onRender() {
	for cid := len(Core.Render.Canvases) - 1; cid >= 0; cid-- {
		if Core.Render.Canvases[cid].renderThisFrame() {
			Core.Render.Canvases[cid].render()
		}
	}
}

func (me *RenderCanvas) render() {
	if me.isRtt {
		me.frameBuf.Bind()
	}
	Core.Render.states.SetFramebufferSrgb(me.Srgb)
	for view := 0; view < len(me.Views); view++ {
		me.Views[view].render()
	}
	Core.Render.states.SetFramebufferSrgb(false)
	if me.isRtt {
		me.frameBuf.Unbind()
		thrRend.quadTex = me.frameBuf.RenderTextureHandle(0)
	}
}

func (me *RenderView) render() {
	if me.Enabled {
		thrRend.curView, thrRend.curTech, thrRend.curEffect = me, me.Technique, nil
		Core.Render.states.Apply(&me.RenderStates)
		if me.Port.shouldScissor {
			Core.Render.states.ForceEnableScissorTest()
			gl.Scissor(me.Port.glVpX, me.Port.glVpY, me.Port.glVpW, me.Port.glVpH)
		}
		gl.Viewport(me.Port.glVpX, me.Port.glVpY, me.Port.glVpW, me.Port.glVpH)
		gl.Clear(me.RenderStates.Other.ClearBits)
		me.Technique.render()
		if me.Port.shouldScissor {
			Core.Render.states.ForceDisableScissorTest()
		}
	}
}

func (me *RenderTechniqueQuad) render() {
	me.glVao.Bind()
	thrRend.nextTech, thrRend.nextEffect = me, &me.Effect
	Core.useTechFx()
	gl.BindTexture(gl.TEXTURE_2D, thrRend.quadTex)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}

func (me *RenderTechniqueScene) render() {
	thrRend.nextTech = me
	thrRend.curCam = &me.Camera
	if scene := me.Camera.Scene(); scene != nil {
		scene.render()
	}
}

func (me *Scene) render() {
	for id := len(me.allNodes) - 1; id > -1; id-- {
		if me.allNodes.Ok(id) && thrRend.curCam.thrRend.nodeRender[&me.allNodes[id]] {
			me.allNodes[id].render()
		}
	}
}

func (me *SceneNode) render() {
	if mesh, mat := me.meshMat(); mesh != nil && mat != nil {
		if mat.HasFaceEffects() {
			for i, l := int32(0), int32(len(mesh.raw.faces)); i < l; i++ {
				thrRend.nextEffect = mat.faceEffect(&mesh.raw.faces[i])
				Core.useTechFx()
				mesh.meshBuffer.use()
				thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &thrRend.curCam.thrRend.nodeProjMats[me][0])
				gl.DrawElementsBaseVertex(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(mesh.meshBufOffsetIndices+(i*3*4))), gl.Int(mesh.meshBufOffsetBaseIndex))
			}
		} else {
			thrRend.nextEffect = Core.Libs.Effects.get(mat.DefaultEffectID)
			Core.useTechFx()
			mesh.meshBuffer.use()
			if me.Render.skyMode {
				Core.Render.states.DisableFaceCulling()
				gl.DepthFunc(gl.LEQUAL)
				thrRend.curProg.Uniform1i("uni_int_Sky", 1)
			}
			thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &thrRend.curCam.thrRend.nodeProjMats[me][0])
			gl.DrawElementsBaseVertex(gl.TRIANGLES, mesh.raw.lastNumIndices, gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(mesh.meshBufOffsetIndices)), gl.Int(mesh.meshBufOffsetBaseIndex))
			if me.Render.skyMode {
				thrRend.curProg.Uniform1i("uni_int_Sky", 0)
				gl.DepthFunc(gl.LESS)
			}
		}
	}
}
