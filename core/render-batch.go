package core

/*

render sorting per pass:

- by VAO					(gl.BindVertexArray)
- by material shader		(gl.UseProgram, currently not happening)
- by texture				(gl.BindTexture, for multi-tex also gl.ActiveTexture and gl.Uniform1i)
- "by other gl.UniformXyz"	(undefined currently)
- by node					(gl.UniformMatrix)
- by mesh					(gl.DrawElements)

*/

type renderBatch struct {
	cam *Camera
}

func newRenderBatch(cam *Camera) (me *renderBatch) {
	me = &renderBatch{cam: cam}
	return
}

func (me *renderBatch) render() {
	//	[GPU] per mat: diffusetex.bind(), uniform(diffuse)
	//	[CPU] per node: mult node-matrix with cam-matrix
	//	[GPU] per node: gl.UniformMatrix4fv()
	//	[GPU] meshBuf.Use()  -> gl.BindVertexArray()
	//	[GPU] gl.DrawElementsBaseVertex()
}
