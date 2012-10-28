package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

var (
	techs map[string]iRenderTechnique
)

type fTechniqueConstructor func (string) iRenderTechnique

func initTechniques () {
	techs = map[string]iRenderTechnique {}
	for techName, techMaker := range map[string]fTechniqueConstructor { "rt_unlit_colored": newTechnique_UnlitColored, "rt_unlit_textured": newTechnique_UnlitTextured } {
		techs[techName] = techMaker(techName)
	}
}

type iRenderTechnique interface {
	initMeshBuffer (meshBuffer *TMeshBuffer)
	name () string
	onPreRender ()
	onRenderMesh ()
	onRenderNode ()
}

type tTechnique_Base struct {
	prog *glutil.TShaderProgram
}

	func (me *tTechnique_Base) initMeshBuffer (meshBuffer *TMeshBuffer) {
		gl.EnableVertexAttribArray(me.prog.AttrLocs["aPos"])
		gl.VertexAttribPointer(me.prog.AttrLocs["aPos"], 3, gl.FLOAT, gl.FALSE, 8 * 4, gl.Pointer(nil))
	}

	func (me *tTechnique_Base) name () string {
		return me.prog.Name
	}

	func (me *tTechnique_Base) onPreRender () {
	}

	func (me *tTechnique_Base) onRenderMesh () {
	}

	func (me *tTechnique_Base) onRenderNode () {
	}

	func (me *tTechnique_Base) setProg (name string, unifs []string, attrs []string) {
		var prog = glShaderMan.AllProgs[name]
		prog.SetUnifLocations("uMatCam", "uMatModelView", "uMatProj")
		if len(unifs) > 0 { prog.SetUnifLocations(unifs...) }
		prog.SetAttrLocations("aPos")
		if len(attrs) > 0 { prog.SetAttrLocations(attrs...) }
		me.prog = prog
	}

type tTechnique_UnlitColored struct {
	tTechnique_Base
}

	func newTechnique_UnlitColored (progName string) iRenderTechnique {
		var tech = &tTechnique_UnlitColored {}
		tech.tTechnique_Base.setProg(progName, nil, nil)
		return tech
	}

type tTechnique_UnlitTextured struct {
	tTechnique_Base
}

	func newTechnique_UnlitTextured (progName string) iRenderTechnique {
		var tech = &tTechnique_UnlitTextured {}
		tech.tTechnique_Base.setProg(progName, []string { "uTex0" }, []string { "aTexCoords" })
		return tech
	}

	func (me *tTechnique_UnlitTextured) initMeshBuffer (meshBuffer *TMeshBuffer) {
		me.tTechnique_Base.initMeshBuffer(meshBuffer)
		gl.EnableVertexAttribArray(me.prog.AttrLocs["aTexCoords"])
		gl.VertexAttribPointer(me.prog.AttrLocs["aTexCoords"], 2, gl.FLOAT, gl.FALSE, 8 * 4, gl.Offset(nil, 3 * 4))
	}

	func (me *tTechnique_UnlitTextured) onRenderNode () {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, Core.Textures[curNode.mat.texKey].glTex)
		gl.Uniform1i(curProg.UnifLocs["uTex0"], 0)
	}
