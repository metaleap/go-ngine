package core

import (
	"strings"
	"time"

	ugl "github.com/go3d/go-opengl/util"
	ustr "github.com/metaleap/go-util/str"
)

type uberShaderFunc struct {
	name, rawSrc      string
	dependsOn, inputs map[string]bool
}

func newUberShaderFunc(name, rawSrc string) (me *uberShaderFunc) {
	me = &uberShaderFunc{
		name: name, rawSrc: rawSrc,
		dependsOn: map[string]bool{}, inputs: map[string]bool{},
	}
	return
}

type uberShader struct {
	rawSources       map[string]string
	tmpAtts, tmpUnis []string
	funcs            struct {
		// compute  map[string]*uberShaderFunc
		// domain   map[string]*uberShaderFunc
		fragment map[string]*uberShaderFunc
		// geometry map[string]*uberShaderFunc
		// hull     map[string]*uberShaderFunc
		shared map[string]*uberShaderFunc
		vertex map[string]*uberShaderFunc
	}
}

func (me *uberShader) init() {
	me.rawSources = map[string]string{}
}

func (me *uberShader) allMaps() map[string]*map[string]*uberShaderFunc {
	return map[string]*map[string]*uberShaderFunc{
		// "cx_": &me.funcs.compute,
		// "dx_": &me.funcs.domain,
		"fx_": &me.funcs.fragment,
		// "gx_": &me.funcs.geometry,
		// "hx_": &me.funcs.hull,
		"vx_": &me.funcs.vertex,
		"":    &me.funcs.shared,
	}
}

func (me *uberShader) inoutTypeSpec(inout string) string {
	inout = inout[strings.Index(inout, "_")+1:]
	return inout[:strings.Index(inout, "_")]
}

func (me *uberShader) loadFromRawSources() {
	var (
		lines, funcLines []string
		line, funcName   string
		mp               *map[string]*uberShaderFunc
		i, lastFunc      int
	)
	mps := me.allMaps()
	for _, mp = range mps {
		*mp = map[string]*uberShaderFunc{}
	}
	for fileName, fileSrc := range me.rawSources {
		if strings.HasPrefix(fileName, "uber-") {
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
	mps := me.allMaps()
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
			for _, pref = range []string{"att_", "uni_", "var_"} {
				for _, pref2 = range ustr.ExtractAllIdentifiers(fn.rawSrc, pref) {
					fn.inputs[pref2] = true
				}
			}
		}
	}
	// for pref, mp = range mps {
	// 	println(pref + "===>")
	// 	for _, fn = range *mp {
	// 		println(fmtStr("%s deps: %#v", fn.name, fn.dependsOn))
	// 		println(fmtStr("%s inputs: %#v", fn.name, fn.inputs))
	// 	}
	// }
}

func (me *uberShader) program(vertTech string, fragFx *FxEffect) (prog *ugl.Program) {
	pname := "uber_" + vertTech + "_" + fragFx.uberName
	if prog = glc.progMan.Programs[pname]; prog == nil {
		var err error
		me.tmpAtts, me.tmpUnis = nil, nil
		if err = me.setShaderSources(pname, vertTech, fragFx); err == nil {
			var dur time.Duration
			if dur, err = glcProgsMake(true, pname); err == nil {
				Diag.LogShaders("Built new shader program '%s' in %v", pname, dur)
				Stats.addProgCompile(1, dur.Nanoseconds())
				prog = glc.progMan.Programs[pname]
				if err = prog.SetAttrLocations(me.tmpAtts...); err == nil {
					err = prog.SetUnifLocations(me.tmpUnis...)
				}
			}
		}
		if err != nil {
			Diag.LogErr(err)
		}
	}
	return
}

func (me *uberShader) setShaderSources(pname, vertTech string, fragFx *FxEffect) (err error) {
	fragInputs := map[string]bool{}
	if err = me.setShaderSourceFrag(pname, fragFx, fragInputs); err == nil {
		err = me.setShaderSourceVert(pname, vertTech, fragInputs)
	}
	return
}

func (me *uberShader) setShaderSourceEnsureFunc(fn *uberShaderFunc, srcBody *ustr.Buffer, inputs map[string]bool) error {
	var (
		str string
		df  *uberShaderFunc
	)
	for str, _ = range fn.inputs {
		inputs[str] = true
	}
	for str, _ = range fn.dependsOn {
		for _, m := range me.allMaps() {
			if df = (*m)[str]; df != nil {
				break
			}
		}
		if df == nil {
			return fmtErr(str)
		}
		me.setShaderSourceEnsureFunc(df, srcBody, inputs)
	}
	srcBody.Writeln(fn.rawSrc + "\n")
	return nil
}

func (me *uberShader) setShaderSourceFrag(pname string, fx *FxEffect, inputs map[string]bool) (err error) {
	var (
		srcBody, srcHead ustr.Buffer
		shid             string
		procID           string
		shader           *fxProc
		shFunc           *uberShaderFunc
	)
	srcHead.Writeln("out vec3 out_Color;")
	for _, procID = range fx.uberProcIDs {
		if shader = Core.Rendering.Fx.procs[procID]; shader == nil {
			err = fmtErr("uberShader.setShaderSourceFrag('%s') -- unknown fxProc ID '%s'", pname, procID)
			return
		}
		if shFunc = me.funcs.fragment[shader.FuncName]; shFunc == nil {
			err = fmtErr("uberShader.setShaderSourceFrag('%s') -- unknown fragment func '%s'", pname, shader.FuncName)
			return
		}
		me.setShaderSourceEnsureFunc(shFunc, &srcBody, inputs)
	}
	for shid, _ = range inputs {
		switch shid[:4] {
		case "uni_":
			ustr.AppendUnique(&me.tmpUnis, shid)
			srcHead.Writeln("uniform %s %s;", me.inoutTypeSpec(shid), shid)
		case "var_":
			srcHead.Writeln("in %s %s;", me.inoutTypeSpec(shid), shid)
		}
	}
	srcBody.Writeln("void main () {")
	srcBody.Writeln("\tvec3 vCol = vec3(0);")
	for _, procID = range fx.uberProcIDs {
		shFunc = me.funcs.fragment[Core.Rendering.Fx.procs[procID].FuncName]
		srcBody.Writeln("\t%s(vCol);", shFunc.name)
	}
	srcBody.Writeln("\tout_Color = vCol;")
	srcBody.Writeln("}")
	glc.progMan.RawSources.Fragment[pname] = srcHead.String() + "\n" + srcBody.String()
	return
}

func (me *uberShader) setShaderSourceVert(pname, vertTech string, varyings map[string]bool) (err error) {
	var (
		i                int
		srcBody, srcHead ustr.Buffer
		shFunc           *uberShaderFunc
		inout, fname     string
		outputs          = []string{"gl_Position"}
		inputs           = map[string]bool{}
	)
	for inout, _ = range varyings {
		if inout[:4] == "var_" {
			outputs = append(outputs, inout)
		}
	}
	for i, inout = range outputs {
		if fname = "vx_" + vertTech + "_" + inout; i > 0 {
			srcHead.Writeln("out %s %s;", me.inoutTypeSpec(inout), inout)
		}
		if shFunc = me.funcs.vertex[fname]; shFunc == nil {
			err = fmtErr("uberShader.setShaderSourceVert('%s') -- unknown vertex func '%s'", pname, fname)
			return
		}
		me.setShaderSourceEnsureFunc(shFunc, &srcBody, inputs)
	}
	for inout, _ = range inputs {
		switch inout[:4] {
		case "att_":
			ustr.AppendUnique(&me.tmpAtts, inout)
			srcHead.Writeln("in %s %s;", me.inoutTypeSpec(inout), inout)
		case "uni_":
			ustr.AppendUnique(&me.tmpUnis, inout)
			srcHead.Writeln("uniform %s %s;", me.inoutTypeSpec(inout), inout)
		}
	}
	srcBody.Writeln("void main () {")
	for _, inout = range outputs {
		srcBody.Writeln("\t%s = vx_%s_%s();", inout, vertTech, inout)
	}
	srcBody.Writeln("}")
	glc.progMan.RawSources.Vertex[pname] = srcHead.String() + "\n" + srcBody.String()
	return
}
