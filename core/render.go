package core

import (
	gl "github.com/go3d/go-opengl/core"
)

func (_ *EngineCore) onRender() {
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
	for cam := 0; cam < len(me.Cams); cam++ {
		thrRend.curCam = me.Cams[cam]
		thrRend.curCam.render()
	}
	Core.Render.states.SetFramebufferSrgb(false)
	if me.isRtt {
		me.frameBuf.Unbind()
		thrRend.quadTex = me.frameBuf.RenderTextureHandle(0)
	}
}

func (me *Camera) render() {
	if me.Enabled {
		thrRend.curTech, thrRend.curEffect = me.Render.Technique, nil
		Core.Render.states.Apply(&me.thrRend.states)
		if me.Render.Viewport.shouldScissor {
			Core.Render.states.ForceEnableScissorTest()
			gl.Scissor(me.Render.Viewport.glVpX, me.Render.Viewport.glVpY, me.Render.Viewport.glVpW, me.Render.Viewport.glVpH)
		}
		gl.Viewport(me.Render.Viewport.glVpX, me.Render.Viewport.glVpY, me.Render.Viewport.glVpW, me.Render.Viewport.glVpH)
		gl.Clear(me.thrRend.states.Other.ClearBits)
		me.Render.Technique.render()
		if me.Render.Viewport.shouldScissor {
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
	if scene := thrRend.curCam.scene(); scene != nil && scene.RootNode.thrRend.camRender[thrRend.curCam] {
		scene.RootNode.renderChildren()
		scene.RootNode.renderSelf() // might be a skybox so "render" the root last
	}
}

func (me *Node) renderChildren() {
	for _, subNode := range me.ChildNodes.M {
		if subNode.thrRend.camRender[thrRend.curCam] {
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
					thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &me.thrRend.camProjMats[thrRend.curCam][0])
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
				thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &me.thrRend.camProjMats[thrRend.curCam][0])
				gl.DrawElementsBaseVertex(gl.TRIANGLES, mesh.raw.lastNumIndices, gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(mesh.meshBufOffsetIndices)), gl.Int(mesh.meshBufOffsetBaseIndex))
				if me.Rendering.skyMode {
					thrRend.curProg.Uniform1i("uni_int_Sky", 0)
					gl.DepthFunc(gl.LESS)
				}
			}
		}
	}
}
