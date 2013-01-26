package core

import (
	gl "github.com/chsc/gogl/gl42"
	ugl "github.com/go3d/go-glutil"
)

//	ONLY used for Core.Rendering.PostFx.
type PostFx struct {
	glVao             ugl.VertexArray
	glWidth, glHeight gl.Sizei
	prog              *ugl.Program
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
	//Core.Rendering.Samplers.NoFilteringClamp.Bind(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, me.glWidth, me.glHeight)
	ugl.LogLastError("pre-clrscr")
	gl.Clear(gl.COLOR_BUFFER_BIT)
	ugl.LogLastError("post-clrscr")
	me.glVao.Bind()
	me.prog.Use()
	mainCanvas.frameBuf.BindTexture(0)
	gl.Uniform1i(me.prog.UnifLocs["uTexRendering"], 0)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	me.glVao.Unbind()
}
