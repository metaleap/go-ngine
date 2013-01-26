package core

import (
	"log"
	"time"

	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

type glShaderManager struct {
	names   []string
	sources glShaderSources
	progs   map[string]*ugl.Program
}

func newShaderManager() (me *glShaderManager) {
	me = &glShaderManager{progs: map[string]*ugl.Program{}}
	return
}

func (me *glShaderManager) dispose() {
	for _, prog := range me.progs {
		prog.Dispose()
	}
	me.progs = map[string]*ugl.Program{}
}

func (me *glShaderManager) compileAll() (err error) {
	var prog *ugl.Program
	src, timeStart, defines := me.sources, time.Now(), map[string]interface{}{}
	for _, name := range me.names {
		prog = ugl.NewProgram(name)
		prog.Create()
		if err = prog.CompileAndLinkShaders(src.compute[name], src.fragment[name], src.geometry[name], src.tessCtl[name], src.tessEval[name], src.vertex[name], defines); err != nil {
			prog.Dispose()
			return
		} else {
			me.progs[name] = prog
		}
	}
	if err == nil {
		log.Printf("Shader compilation time: %v\n", time.Now().Sub(timeStart))
	}
	return
}

type glShaderSources struct {
	compute, fragment, geometry, tessCtl, tessEval, vertex map[string]string
}

func (me *glShaderSources) enumerate() map[gl.Enum]int {
	return map[gl.Enum]int{
		0 /*gl.COMPUTE_SHADER*/ :  0,
		gl.FRAGMENT_SHADER:        1,
		gl.GEOMETRY_SHADER:        2,
		gl.TESS_CONTROL_SHADER:    3,
		gl.TESS_EVALUATION_SHADER: 4,
		gl.VERTEX_SHADER:          5,
	}
}

func (me *glShaderSources) init() {
	me.compute, me.fragment, me.geometry, me.tessCtl, me.tessEval, me.vertex = map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}
}

func (me *glShaderSources) source(glShaderType gl.Enum, shaderName string) string {
	switch glShaderType {
	case gl.FRAGMENT_SHADER:
		return me.fragment[shaderName]
	case gl.GEOMETRY_SHADER:
		return me.geometry[shaderName]
	case gl.TESS_CONTROL_SHADER:
		return me.tessCtl[shaderName]
	case gl.TESS_EVALUATION_SHADER:
		return me.tessEval[shaderName]
	case gl.VERTEX_SHADER:
		return me.vertex[shaderName]
	default:
		return me.compute[shaderName]
	}
	return ""
}
