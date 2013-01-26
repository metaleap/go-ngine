package core

import (
	"log"
	"time"

	ugl "github.com/go3d/go-glutil"
)

type glShaderManager struct {
	names   []string
	sources struct {
		compute, fragment, geometry, tessCtl, tessEval, vertex map[string]string
	}
	progs map[string]*ugl.Program
}

func newShaderManager() (me *glShaderManager) {
	me = &glShaderManager{progs: map[string]*ugl.Program{}}
	src := &me.sources
	src.compute, src.fragment, src.geometry, src.tessCtl, src.tessEval, src.vertex = map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}
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
	timeStart, src, defines := time.Now(), &me.sources, map[string]interface{}{}
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
