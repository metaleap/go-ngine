package core

import (
	"strings"

	ustr "github.com/metaleap/go-util/str"
)

type uberShaderFunc struct {
	isMain       bool
	name, rawSrc string
	dependsOn    map[string]bool
}

func newUberShaderFunc(name, rawSrc string) (me *uberShaderFunc) {
	me = &uberShaderFunc{name: name, rawSrc: rawSrc, dependsOn: map[string]bool{}, isMain: strings.Index(name, "_Main") == 2}
	return
}

type uberShader struct {
	rawSources map[string]string
	funcs      struct {
		compute  map[string]*uberShaderFunc
		domain   map[string]*uberShaderFunc
		fragment map[string]*uberShaderFunc
		geometry map[string]*uberShaderFunc
		hull     map[string]*uberShaderFunc
		shared   map[string]*uberShaderFunc
		vertex   map[string]*uberShaderFunc
	}
}

func (me *uberShader) init() {
	me.rawSources = map[string]string{}
}

func (me *uberShader) mapOfMaps() map[string]*map[string]*uberShaderFunc {
	return map[string]*map[string]*uberShaderFunc{
		"cx_": &me.funcs.compute,
		"dx_": &me.funcs.domain,
		"fx_": &me.funcs.fragment,
		"gx_": &me.funcs.geometry,
		"hx_": &me.funcs.hull,
		"vx_": &me.funcs.vertex,
		"":    &me.funcs.shared,
	}
}

func (me *uberShader) loadFromRawSources() {
	var (
		lines, funcLines []string
		line, funcName   string
		mp               *map[string]*uberShaderFunc
		i, lastFunc      int
	)
	mps := me.mapOfMaps()
	for _, mp = range mps {
		*mp = map[string]*uberShaderFunc{}
	}
	for fileName, fileSrc := range me.rawSources {
		if strings.HasPrefix(fileName, "fx-") {
			if lines = strings.Split(fileSrc, "\n"); len(lines) > 0 {
				for i, line = range lines {
					if !strings.HasPrefix(line, "\t") {
						if strings.HasPrefix(line, "}") {
							mp, funcLines = nil, lines[lastFunc:i+1]
							funcName = funcLines[0]
							funcName = funcName[strings.Index(funcName, " ")+1:]
							funcName = funcName[:ustr.IndexAny(funcName, " ", "(")]
							if mp = mps[funcName[:3]]; mp == nil {
								mp = mps[""]
							}
							(*mp)[funcName] = newUberShaderFunc(funcName, strings.Join(funcLines, "\n"))
							lastFunc = i + 1
						} else {
							lastFunc = i
						}
					}
				}
			}
		}
	}
}

func (me *uberShader) processFuncs() {
	var (
		pref, pref2 string
		fn, fn2     *uberShaderFunc
		mp, mp2     *map[string]*uberShaderFunc
	)
	mps := me.mapOfMaps()
	for pref, mp = range mps {
		for _, fn = range *mp {
			//	annotate other funcs depending on this func
			for _, fn2 = range *mp {
				if fn2 != fn && strings.Index(fn2.rawSrc, fn.name) > 0 {
					fn2.dependsOn[fn.name] = true
				}
			}
			if len(pref) == 0 {
				//	we're in the map of shared-funcs, so need to check all other maps for dep-on-fn too
				for pref2, mp2 = range mps {
					if pref2 != pref {
						for _, fn2 = range *mp2 {
							if fn2 != fn && strings.Index(fn2.rawSrc, fn.name) > 0 {
								fn2.dependsOn[fn.name] = true
							}
						}
					}
				}
			}
			//	annotate this func for any uniforms, varyings or attributes
		}
	}
	for pref, mp = range mps {
		println(pref + "===>")
		for _, fn = range *mp {
			println(fmtStr("%s deps: %#v", fn.name, fn.dependsOn))
		}
	}
}
