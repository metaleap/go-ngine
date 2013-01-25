package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

//	ONLY used for Core.Rendering.PostFx.
type PostFx struct {
	glVao             ugl.VertexArray
	glWidth, glHeight gl.Sizei
	prog              *ugl.ShaderProgram
}

func (me *PostFx) dispose() {
	me.glVao.Dispose()
}

func (me *PostFx) init() {
	me.prog = glShaderMan.progs["postfx"]
	me.prog.SetUnifLocations("uTexRendering")
	me.glVao.Create()
}

func (me *PostFx) render() {
	curProg, curMat, curTechnique, curMatKey = nil, nil, nil, ""
	Core.Rendering.states.DisableDepthTest()
	Core.Rendering.states.DisableFaceCulling()
	gl.UseProgram(me.prog.Program)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.Viewport(0, 0, me.glWidth, me.glHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	mainCanvas.frameBuf.BindTexture(0)
	gl.Uniform1i(me.prog.UnifLocs["uTexRendering"], 0)
	me.glVao.Bind()
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}
