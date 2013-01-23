package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

//	ONLY used for Core.Rendering.PostFx.
type PostFx struct {
	glVao             gl.Uint
	glWidth, glHeight gl.Sizei
	prog              *ugl.ShaderProgram
}

func (me *PostFx) dispose() {
	gl.DeleteVertexArrays(1, &me.glVao)
}

func (me *PostFx) init() {
	me.prog = glShaderMan.progs["postfx"]
	me.prog.SetUnifLocations("uTexRendering")
	gl.GenVertexArrays(1, &me.glVao)
}

func (me *PostFx) render() {
	curProg, curMat, curTechnique, curMatKey = nil, nil, nil, ""
	gl.UseProgram(me.prog.Program)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	gl.Viewport(0, 0, me.glWidth, me.glHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	mainCanvas.frameBuf.BindTexture(0)
	gl.Uniform1i(me.prog.UnifLocs["uTexRendering"], 0)
	gl.BindVertexArray(me.glVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}
