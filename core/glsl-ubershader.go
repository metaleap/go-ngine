package core

import (
	"strings"
	"time"

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
	opIndices        map[FxOp]int
	pname            string
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

func (me *uberShader) ensureProg() {
	me.pname = thrRend.curEffect.uberPnames[thrRend.curTech.name()]
	if thrRend.tmpProg = glc.progMan.Programs[me.pname]; thrRend.tmpProg == nil {
		var err error
		vertTech, fragFx := thrRend.curTech.name(), thrRend.curEffect
		me.tmpAtts, me.tmpUnis = nil, nil
		if err = me.setShaderSources(vertTech, fragFx); err == nil {
			var dur time.Duration
			if dur, err = glcProgsMake(true, me.pname); err == nil {
				Diag.LogShaders("Built new GLSL shader program '%s' in %v", me.pname, dur)
				Stats.addProgCompile(1, dur.Nanoseconds())
				thrRend.tmpProg = glc.progMan.Programs[me.pname]
				if err = thrRend.tmpProg.SetAttrLocations(me.tmpAtts...); err == nil {
					err = thrRend.tmpProg.SetUnifLocations(me.tmpUnis...)
				}
			}
		}
		if err != nil {
			Diag.LogErr(err)
		}
	}
}

func (me *uberShader) setShaderSources(vertTech string, fragFx *FxEffect) (err error) {
	fragInputs := map[string]bool{}
	me.opIndices = map[FxOp]int{}
	if err = me.setShaderSourceFrag(fragFx, fragInputs); err == nil {
		err = me.setShaderSourceVert(vertTech, fragInputs)
	}
	me.opIndices = nil
	return
}

func (me *uberShader) setShaderSourceEnsureFunc(op FxOp, fn *uberShaderFunc, srcBody *ustr.Buffer, inputs map[string]bool) error {
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
			return errf(str)
		}
		me.setShaderSourceEnsureFunc(nil, df, srcBody, inputs)
	}
	srcBody.Writeln(fn.rawSrc + "\n")
	return nil
}

func (me *uberShader) setShaderSourceFrag(fx *FxEffect, inputs map[string]bool) (err error) {
	var (
		srcBody, srcHead ustr.Buffer
		shid             string
		op               FxOp
		shader           *fxProc
		shFunc           *uberShaderFunc
	)
	srcHead.Writeln("out vec3 out_Color;")
	opc, opCounts := 0, map[string]int{}
	for _, op = range fx.Ops {
		if op.Enabled() {
			opc = opCounts[op.ProcID()]
			me.opIndices[op] = opc
			opCounts[op.ProcID()] = opc + 1
		}
	}
	opCounts = nil

	for _, op = range fx.Ops {
		if op.Enabled() {
			if shader = Core.Rendering.Fx.procs[op.ProcID()]; shader == nil {
				err = errf("uberShader.setShaderSourceFrag('%s') -- unknown fxProc ID '%s'", me.pname, op.ProcID())
				return
			}
			if shFunc = me.funcs.fragment[shader.FuncName]; shFunc == nil {
				err = errf("uberShader.setShaderSourceFrag('%s') -- unknown fragment func '%s'", me.pname, shader.FuncName)
				return
			}
			me.setShaderSourceEnsureFunc(op, shFunc, &srcBody, inputs)
		}
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
	for _, op = range fx.Ops {
		if op.Enabled() {
			shFunc = me.funcs.fragment[Core.Rendering.Fx.procs[op.ProcID()].FuncName]
			srcBody.Writeln("\t%s(vCol);", shFunc.name)
		}
	}
	srcBody.Writeln("\tout_Color = vCol;")
	srcBody.Writeln("}")
	glc.progMan.RawSources.Fragment[me.pname] = srcHead.String() + "\n" + srcBody.String()
	return
}

func (me *uberShader) setShaderSourceVert(vertTech string, varyings map[string]bool) (err error) {
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
			err = errf("uberShader.setShaderSourceVert('%s') -- unknown vertex func '%s'", me.pname, fname)
			return
		}
		me.setShaderSourceEnsureFunc(nil, shFunc, &srcBody, inputs)
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
	glc.progMan.RawSources.Vertex[me.pname] = srcHead.String() + "\n" + srcBody.String()
	return
}
