package core

import (
	"fmt"
	"log"
	"time"

	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

type tGlShaderManager struct {
	AllNames []string
	AllProgs map[string]*glutil.TShaderProgram
	AllSources *tGlShaderSources
}

	func newShaderManager () *tGlShaderManager {
		var sm = &tGlShaderManager {}
		sm.AllProgs = map[string]*glutil.TShaderProgram {}
		return sm
	}

	func (me *tGlShaderManager) dispose () {
		var doClean = func (sprog **glutil.TShaderProgram) {
			var sp *glutil.TShaderProgram = *sprog
			if sp != nil { sp.CleanUp(); *sprog = nil }
		}
		for _, prog := range me.AllProgs {
			doClean(&prog)
		}
	}

	func (me *tGlShaderManager) compileAll () error {
		var err error
		var shaderName, shaderSrc string
		var shaderTypeIndex int
		var glShaderType gl.Enum
		var timeStart = time.Now()
		var glStatus gl.Int
		var glShaders = []gl.Uint { 0, 0, 0, 0, 0, 0 }
		var shaderProg *glutil.TShaderProgram
		var defines = map[string]interface{} { }
		for _, shaderName = range me.AllNames {
			for glShaderType, shaderTypeIndex = range me.AllSources.enumerate() {
				if shaderSrc = me.AllSources.source(glShaderType, shaderName); len(shaderSrc) > 0 {
					glShaders[shaderTypeIndex] = gl.CreateShader(glShaderType)
					glutil.ShaderSource(shaderName, glShaders[shaderTypeIndex], shaderSrc, defines, false, "150")
					gl.CompileShader(glShaders[shaderTypeIndex])
					if gl.GetShaderiv(glShaders[shaderTypeIndex], gl.COMPILE_STATUS, &glStatus); glStatus == 0 { err = fmt.Errorf("SHADER %s: %s\n", shaderName, glutil.ShaderInfoLog(glShaders[shaderTypeIndex], true)) }
				} else {
					glShaders[shaderTypeIndex] = 0
				}
				if err != nil { break }
			}
			if err == nil {
				if shaderProg, err = glutil.NewShaderProgram(shaderName, glShaders[0], glShaders[1], glShaders[2], glShaders[3], glShaders[4], glShaders[5]); err == nil {
					me.AllProgs[shaderName] = shaderProg
					/*
					if shaderName == "postfx" {
						me.Prog_PostFx = shaderProg
					} else if shaderName == "unlit" {
						me.Prog_Unlit = shaderProg
					} else if shaderName == "pvlit" {
						me.Prog_PvLit = shaderProg
					} else if shaderName == "pplit" {
						me.Prog_PpLit = shaderProg
					}
					*/
				}
			}
			if err != nil { break }
		}
		if err == nil { log.Printf("Shader compilation time: %v\n", time.Now().Sub(timeStart)) }
		return err
	}

type tGlShaderSources struct {
	Compute, Fragment, Geometry, TessCtl, TessEval, Vertex map[string]string
}

	func newGlShaderSources () *tGlShaderSources {
		return &tGlShaderSources { map[string]string {}, map[string]string {}, map[string]string {}, map[string]string {}, map[string]string {}, map[string]string {} }
	}

	func (me *tGlShaderSources) enumerate () map[gl.Enum]int {
		return map[gl.Enum]int { 0: 0, gl.FRAGMENT_SHADER: 1, gl.GEOMETRY_SHADER: 2, gl.TESS_CONTROL_SHADER: 3, gl.TESS_EVALUATION_SHADER: 4, gl.VERTEX_SHADER: 5 }
	}

	func (me *tGlShaderSources) source (glShaderType gl.Enum, shaderName string) string {
		switch glShaderType {
		case gl.FRAGMENT_SHADER:
			return me.Fragment[shaderName]
		case gl.GEOMETRY_SHADER:
			return me.Geometry[shaderName]
		case gl.TESS_CONTROL_SHADER:
			return me.TessCtl[shaderName]
		case gl.TESS_EVALUATION_SHADER:
			return me.TessEval[shaderName]
		case gl.VERTEX_SHADER:
			return me.Vertex[shaderName]
		}
		return ""
	}
