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
	if scene := me.Camera.scene(); scene != nil && me.Camera.thrRend.nodeRender[&scene.RootNode] {
		scene.RootNode.renderChildren()
		scene.RootNode.renderSelf() // might be a skybox so "render" the root last
	}
}

func (me *Node) renderChildren() {
	for _, subNode := range me.ChildNodes.M {
		if thrRend.curCam.thrRend.nodeRender[subNode] {
			subNode.renderSelf()
			subNode.renderChildren()
		}
	}
}

func (me *Node) renderSelf() {
	if mesh := me.mesh(); mesh != nil {
		if mat := me.material(); mat != nil {
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
				if me.Rendering.skyMode {
					Core.Render.states.DisableFaceCulling()
					gl.DepthFunc(gl.LEQUAL)
					thrRend.curProg.Uniform1i("uni_int_Sky", 1)
				}
				thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &thrRend.curCam.thrRend.nodeProjMats[me][0])
				gl.DrawElementsBaseVertex(gl.TRIANGLES, mesh.raw.lastNumIndices, gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(mesh.meshBufOffsetIndices)), gl.Int(mesh.meshBufOffsetBaseIndex))
				if me.Rendering.skyMode {
					thrRend.curProg.Uniform1i("uni_int_Sky", 0)
					gl.DepthFunc(gl.LESS)
				}
			}
		}
	}
}
