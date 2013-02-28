package core

import (
	"strings"
	"time"

	ugl "github.com/go3d/go-opengl/util"
	ustr "github.com/metaleap/go-util/str"
)

type uberShaderFunc struct {
	name, rawSrc         string
	dependsOn, allInputs []string
}

func newUberShaderFunc(name, rawSrc string) (me *uberShaderFunc) {
	me = &uberShaderFunc{name: name, rawSrc: rawSrc}
	return
}

type uberShaderJob struct {
	pname              string
	progAtts, progUnis []string
	quals              map[string]string
}

func (me *uberShaderJob) init(pname string) {
	me.pname = pname
	me.progAtts, me.progUnis = make([]string, 0, 8), make([]string, 0, 8)
	me.quals = make(map[string]string, 8)
}

type uberShaders struct {
	rawSources map[string]string
	allFuncs   struct {
		// compute  map[string]*uberShaderFunc
		// domain   map[string]*uberShaderFunc
		fragment map[string]*uberShaderFunc
		// geometry map[string]*uberShaderFunc
		// hull     map[string]*uberShaderFunc
		shared map[string]*uberShaderFunc
		vertex map[string]*uberShaderFunc
	}
}

func (me *uberShaders) init() {
	me.rawSources = make(map[string]string, 2)
}

func (me *uberShaders) allMaps() map[string]*map[string]*uberShaderFunc {
	return map[string]*map[string]*uberShaderFunc{
		// "cx_": &me.allFuncs.compute,
		// "dx_": &me.allFuncs.domain,
		"fx_": &me.allFuncs.fragment,
		// "gx_": &me.allFuncs.geometry,
		// "hx_": &me.allFuncs.hull,
		"vx_": &me.allFuncs.vertex,
		"":    &me.allFuncs.shared,
	}
}

func (me *uberShaders) inoutTypeSpec(inout string) string {
	inout = inout[strings.Index(inout, "_")+1:]
	return inout[:strings.Index(inout, "_")]
}

func (me *uberShaders) loadFromRawSources() {
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

func (me *uberShaders) processFuncs() {
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
					ustr.AppendUnique(&fn2.dependsOn, fn.name)
				}
			}
			if len(pref) == 0 {
				//	we're in the map of shared-funcs, so need to check all other maps for dep-on-fn too
				for pref2, mp2 = range mps {
					if pref2 != pref {
						for _, fn2 = range *mp2 {
							if fn2 != fn && strings.Index(fn2.rawSrc, fn.name) > 0 {
								ustr.AppendUnique(&fn2.dependsOn, fn.name)
							}
						}
					}
				}
			}
			//	annotate this func for any uniforms, varyings or attributes
			for _, pref = range []string{"att_", "uni_", "var_"} {
				for _, pref2 = range ustr.ExtractAllIdentifiers(fn.rawSrc, pref) {
					ustr.AppendUnique(&fn.allInputs, pref2)
				}
			}
		}
	}
}

func (me *uberShaders) ensureProg() (prog *ugl.Program) {
	pname := thrRend.curEffect.uberPnames[thrRend.curTech.name()]
	if prog = ogl.progs.Get(pname); prog == nil {
		var job uberShaderJob
		job.init(pname)
		prog = ogl.progs.AddNew(pname)
		var err error
		if err = me.setShaderSources(prog, &job, thrRend.curTech.name(), thrRend.curEffect); err == nil {
			var dur time.Duration
			if dur, err = ogl.makeProgs(false, pname); err == nil {
				Diag.LogShaders("Built new GLSL shader program '%s' in %v", pname, dur)
				Stats.addProgCompile(1, dur.Nanoseconds())
				prog.Tag = thrRend.curTech
				if err = prog.SetAttrLocations(job.progAtts...); err == nil {
					if err = prog.SetUnifLocations(job.progUnis...); err == nil {
						for _, meshBuf := range Core.MeshBuffers.bufs {
							if err = meshBuf.setupVao(prog, thrRend.curTech); err != nil {
								break
							}
						}
					}
				}
			}
		}
		if err != nil {
			Diag.LogErr(err)
		}
	}
	return
}

func (me *uberShaders) setShaderSources(prog *ugl.Program, job *uberShaderJob, vertTech string, fragFx *FxEffect) (err error) {
	fragInputs := map[string]bool{}
	if prog.Sources.In.Fragment, err = me.setShaderSourceFrag(job, fragFx, fragInputs); err == nil {
		prog.Sources.In.Vertex, err = me.setShaderSourceVert(job, vertTech, fragInputs)
	}
	return
}

func (me *uberShaders) setShaderSourceEnsureFunc(op FxOp, fn *uberShaderFunc, srcBody *ustr.Buffer, inputs map[string]bool) error {
	var (
		str, _procID_ string
		parts         []string
		df            *uberShaderFunc
	)
	isFxOp := strings.HasPrefix(fn.name, "fx_") && op != nil
	if isFxOp {
		_procID_ = "_" + op.ProcID() + "_"
	}
	rewriteUnis := map[string]string{}
	for _, str = range fn.allInputs {
		if isFxOp && strings.HasPrefix(str, "uni_") && strings.Contains(str, _procID_) {
			parts = strings.Split(str, "_")
			rewriteUnis[str] = op.unifName(parts[1], parts[3])
			str = rewriteUnis[str]
		}
		inputs[str] = true
	}
	for _, str = range fn.dependsOn {
		if isFxOp && strings.HasPrefix(str, "fx_") {
			return errf("%s depends on %s. One fx_ func must not depend directly on another.", fn.name, str)
		}
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
	str = fn.rawSrc
	if isFxOp {
		str = strings.Replace(str, fn.name, strf("%s%d", fn.name, op.ProcIndex()), -1)
		for k, v := range rewriteUnis {
			str = strings.Replace(str, k, v, -1)
		}
	}
	srcBody.Writeln(str + "\n")
	return nil
}

func (me *uberShaders) setShaderSourceFrag(job *uberShaderJob, fx *FxEffect, inputs map[string]bool) (fragSrc string, err error) {
	var (
		srcBody, srcHead ustr.Buffer
		shid             string
		op               FxOp
		ops              FxOps
		shFunc           *uberShaderFunc
	)
	srcHead.Writeln("out vec3 out_Color;")
	allOps := []FxOps{fx.Ops, fx.OpsX}
	for _, ops = range allOps {
		for _, op = range ops {
			if op.Enabled() {
				if shid = Core.Render.Fx.procFuncs[op.ProcID()]; len(shid) == 0 {
					err = errf("uberShader.setShaderSourceFrag('%s') -- unknown fxProc ID '%s'", job.pname, op.ProcID())
					return
				}
				if shFunc = me.allFuncs.fragment[shid]; shFunc == nil {
					err = errf("uberShader.setShaderSourceFrag('%s') -- unknown fragment func '%s'", job.pname, shid)
					return
				}
				me.setShaderSourceEnsureFunc(op, shFunc, &srcBody, inputs)
				inputs[op.unifName("float", "MixWeight")] = true
				for shid, _ = range inputs {
					if job.quals[shid] = op.qualifiers(shid); len(job.quals[shid]) > 0 {
						job.quals[shid] = job.quals[shid] + " "
					}
				}
			}
		}
	}

	for shid, _ = range inputs {
		switch shid[:4] {
		case "uni_":
			ustr.AppendUnique(&job.progUnis, shid)
			srcHead.Writeln("uniform %s %s;", me.inoutTypeSpec(shid), shid)
		case "var_":
			srcHead.Writeln("%sin %s %s;", job.quals[shid], me.inoutTypeSpec(shid), shid)
		}
	}
	srcBody.Writeln("void main () {")
	srcBody.Writeln("\tvec3 vCol = vec3(0);")
	for _, ops = range allOps {
		for _, op = range ops {
			if op.Enabled() {
				shFunc = me.allFuncs.fragment[Core.Render.Fx.procFuncs[op.ProcID()]]
				srcBody.Writeln("\tvCol = mix(vCol, %s%d(vCol), %s);", shFunc.name, op.ProcIndex(), op.unifName("float", "MixWeight"))
			}
		}
	}
	srcBody.Writeln("\tout_Color = vCol;")
	srcBody.Writeln("}")
	fragSrc = srcHead.String() + "\n" + srcBody.String()
	return
}

func (me *uberShaders) setShaderSourceVert(job *uberShaderJob, vertTech string, varyings map[string]bool) (vertSrc string, err error) {
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
			srcHead.Writeln("%sout %s %s;", job.quals[inout], me.inoutTypeSpec(inout), inout)
		}
		if shFunc = me.allFuncs.vertex[fname]; shFunc == nil {
			err = errf("uberShader.setShaderSourceVert('%s') -- unknown vertex func '%s'", job.pname, fname)
			return
		}
		me.setShaderSourceEnsureFunc(nil, shFunc, &srcBody, inputs)
	}
	for inout, _ = range inputs {
		switch inout[:4] {
		case "att_":
			ustr.AppendUnique(&job.progAtts, inout)
			srcHead.Writeln("in %s %s;", me.inoutTypeSpec(inout), inout)
		case "uni_":
			ustr.AppendUnique(&job.progUnis, inout)
			srcHead.Writeln("uniform %s %s;", me.inoutTypeSpec(inout), inout)
		}
	}
	srcBody.Writeln("void main () {")
	for _, inout = range outputs {
		srcBody.Writeln("\t%s = vx_%s_%s();", inout, vertTech, inout)
	}
	srcBody.Writeln("}")
	vertSrc = srcHead.String() + "\n" + srcBody.String()
	return
}
