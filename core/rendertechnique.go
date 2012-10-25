package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

var (
	techs = map[string]iRenderTechnique {
	}
	techMakers = map[string]fTechniqueConstructor {
		"geometry": newTechnique_Geometry,
		"unlit": newTechnique_Unlit,
	}
)

type fTechniqueConstructor func () iRenderTechnique

type iRenderTechnique interface {
	Name () string
	OnPreRender ()
	OnRenderMesh ()
	OnRenderNode ()
	Program () *glutil.TShaderProgram
}

	func getRenderTechnique (name string) iRenderTechnique {
		var t = techs[name]
		var maker = techMakers[name]
		if (t == nil) && (maker != nil) {
			techs[name] = maker()
			t = techs[name]
		}
		return t
	}

type tTechnique_Base struct {
	prog *glutil.TShaderProgram
}

	func (me *tTechnique_Base) setProg (name string, unifs []string, attrs []string) {
		var prog = glShaderMan.AllProgs[name]
		prog.SetUnifLocations("uMatCam", "uMatModelView", "uMatProj")
		if len(unifs) > 0 { prog.SetUnifLocations(unifs...) }
		prog.SetAttrLocations("aPos")
		if len(attrs) > 0 { prog.SetAttrLocations(attrs...) }
		me.prog = prog
	}

type tTechnique_Geometry struct {
	tTechnique_Base
}

	func newTechnique_Geometry () iRenderTechnique {
		var tech = &tTechnique_Geometry {}
		tech.tTechnique_Base.setProg("geometry", nil, nil)
		return tech
	}

	func (me *tTechnique_Geometry) Name () string {
		return me.prog.Name
	}

	func (me *tTechnique_Geometry) OnPreRender () {
	}

	func (me *tTechnique_Geometry) OnRenderMesh () {
	}

	func (me *tTechnique_Geometry) OnRenderNode () {
	}

	func (me *tTechnique_Geometry) Program () *glutil.TShaderProgram {
		return me.prog
	}

type tTechnique_Unlit struct {
	tTechnique_Base
}

	func newTechnique_Unlit () iRenderTechnique {
		var tech = &tTechnique_Unlit {}
		tech.tTechnique_Base.setProg("unlit", []string { "uTex0" }, []string { "aTexCoords" })
		return tech
	}

	func (me *tTechnique_Unlit) Name () string {
		return me.prog.Name
	}

	func (me *tTechnique_Unlit) OnPreRender () {
	}

	func (me *tTechnique_Unlit) OnRenderMesh () {
	}

	func (me *tTechnique_Unlit) OnRenderNode () {
		if curNode.glVertTexCoordsBuf > 0 {
			gl.BindBuffer(gl.ARRAY_BUFFER, curNode.glVertTexCoordsBuf)
			gl.EnableVertexAttribArray(curProg.AttrLocs["aTexCoords"])
			gl.VertexAttribPointer(curProg.AttrLocs["aTexCoords"], 2, gl.FLOAT, gl.FALSE, 0, gl.Pointer(nil))
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, Core.Textures[curNode.mat.texKey].glTex)
			gl.Uniform1i(curProg.UnifLocs["uTex0"], 0)
		}
	}

	func (me *tTechnique_Unlit) Program () *glutil.TShaderProgram {
		return me.prog
	}
