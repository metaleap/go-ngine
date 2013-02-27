package core

import (
	gl "github.com/go3d/go-opengl/core"
)

func (_ *EngineCore) onRender() {
	for cid := len(Core.Rendering.Canvases) - 1; cid >= 0; cid-- {
		if Core.Rendering.Canvases.Ok(cid) && Core.Rendering.Canvases[cid].renderThisFrame() {
			Core.Rendering.Canvases[cid].render()
		}
	}
}

func (me *RenderCanvas) render() {
	if me.isRtt {
		me.frameBuf.Bind()
	}
	Core.Rendering.states.SetFramebufferSrgb(me.Srgb)
	for cam := 0; cam < len(me.Cameras); cam++ {
		if me.Cameras.Ok(cam) {
			thrRend.curCam = &me.Cameras[cam]
			thrRend.curCam.render()
		}
	}
	Core.Rendering.states.SetFramebufferSrgb(false)
	if me.isRtt {
		me.frameBuf.Unbind()
		thrRend.quadTex = me.frameBuf.RenderTextureHandle(0)
	}
}

func (me *Camera) render() {
	if me.Enabled {
		thrRend.curTech, thrRend.curEffect = me.Rendering.Technique, nil
		Core.Rendering.states.Apply(&me.thrRend.states)
		if me.Rendering.Viewport.shouldScissor {
			Core.Rendering.states.ForceEnableScissorTest()
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
	me.glVao.Bind()
	thrRend.nextTech, thrRend.nextEffect = me, &me.Effect
	Core.useTechFx()
	gl.BindTexture(gl.TEXTURE_2D, thrRend.quadTex)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}

func (me *RenderTechniqueScene) render() {
	thrRend.nextTech = me
	if scene := thrRend.curCam.scene(); scene != nil && scene.RootNode.Rendering.Enabled {
		scene.RootNode.renderChildren()
		scene.RootNode.renderSelf() // might be a skybox so "render" the root last
	}
}

func (me *Node) renderChildren() {
	for _, subNode := range me.ChildNodes.M {
		if subNode.Rendering.Enabled {
			subNode.renderSelf()
			subNode.renderChildren()
		}
	}
}

func (me *Node) renderSelf() {
	if mesh, mat := me.mesh(), me.material(); mesh != nil && mat != nil {
		if mat.HasFaceEffects() {
			for fidx, face := range mesh.raw.faces {
				thrRend.nextEffect = mat.faceEffect(face)
				Core.useTechFx()
				mesh.meshBuffer.use()
				thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &me.thrRend.camProjMats[thrRend.curCam][0])
				gl.DrawElementsBaseVertex(gl.TRIANGLES, 3, gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(mesh.meshBufOffsetIndices+(int32(fidx)*3*4))), gl.Int(mesh.meshBufOffsetBaseIndex))
			}
		} else {
			thrRend.nextEffect = Core.Libs.Effects.Get(mat.DefaultEffectID)
			Core.useTechFx()
			mesh.meshBuffer.use()
			if me.Rendering.skyMode {
				Core.Rendering.states.DisableFaceCulling()
				gl.DepthFunc(gl.LEQUAL)
				thrRend.curProg.Uniform1i("uni_int_Sky", 1)
			}
			thrRend.curProg.UniformMatrix4fv("uni_mat4_VertexMatrix", 1, gl.FALSE, &me.thrRend.camProjMats[thrRend.curCam][0])
			gl.DrawElementsBaseVertex(gl.TRIANGLES, gl.Sizei(len(mesh.raw.indices)), gl.UNSIGNED_INT, gl.Util.PtrOffset(nil, uintptr(mesh.meshBufOffsetIndices)), gl.Int(mesh.meshBufOffsetBaseIndex))
			if me.Rendering.skyMode {
				thrRend.curProg.Uniform1i("uni_int_Sky", 0)
				gl.DepthFunc(gl.LESS)
			}
		}
	}
}
