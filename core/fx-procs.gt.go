package core

import (
	"github.com/go-utils/ustr"
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

func fxKnownProcIDs() []string {
	return []string{"Color", "Tex2D", "TexCube", "Grayscale", "Orangify", "Coords", "Gamma"}
}

type FxProc struct {
	Enabled bool

	weight    gl.Float
	procIndex int
	procID    string
	unifNames map[[2]string]string

	Color struct {
		Rgb ugl.GlVec3
	}

	Tex struct {
		ImageID int
		glUnitU gl.Uint
		glUnitI gl.Int
		glUnitE gl.Enum

		Sampler ugl.Sampler
	}
}

func (me *FxProc) init(procID string, procIndex int) {
	me.procID, me.procIndex, me.Enabled, me.weight = procID, -1, true, 1
	me.setProcIndex(procIndex)
	me.Tex.ImageID = -1
	me.Tex.Sampler = Core.Render.Fx.Samplers.FullFilteringRepeat
}

func (me *FxProc) qualifiers(inout string) (q string) {
	return
}

func (me *FxProc) setProcIndex(index int) {
	if index != me.procIndex {
		me.procIndex = index
		me.unifNames = map[[2]string]string{}
		if me.IsTex() {
			me.Tex.glUnitI = gl.Int(index)
			me.Tex.glUnitU = gl.Uint(index)
			me.Tex.glUnitE = gl.Enum(gl.TEXTURE0 + index)
		}
	}
}

func (me *FxProc) unifName(t, n string) (un string) {
	k := [2]string{t, n}
	if un = me.unifNames[k]; len(un) == 0 {
		un = strf("uni_%s_%s%d_%s", t, me.procID, me.procIndex, n)
		me.unifNames[k] = un
	}
	return
}

func (me *FxProc) use() {
	thrRend.curProg.Uniform1f(me.unifName("float", "MixWeight"), me.weight)
	if me.IsColor() {
		thrRend.curProg.UniformVec3(me.unifName("vec3", "Rgb"), &me.Color.Rgb)
	}
	if me.IsTex() {
		me.Tex.Sampler.Bind(me.Tex.glUnitU)
		if me.IsTex2D() {
			thrRend.curProg.Uniform1i(me.unifName("sampler2D", "Img"), me.Tex.glUnitI)
			if Core.Libs.Images.Tex2D.IsOk(me.Tex.ImageID) {
				ugl.Cache.BindTextureTo(me.Tex.glUnitU, Core.Libs.Images.Tex2D[me.Tex.ImageID].glTex.GlHandle, Core.Libs.Images.Tex2D[me.Tex.ImageID].glTex.GlTarget)
			}
		} else if me.IsTexCube() {
			thrRend.curProg.Uniform1i(me.unifName("samplerCube", "Img"), me.Tex.glUnitI)
			if Core.Libs.Images.TexCube.IsOk(me.Tex.ImageID) {
				ugl.Cache.BindTextureTo(me.Tex.glUnitU, Core.Libs.Images.TexCube[me.Tex.ImageID].glTex.GlHandle, Core.Libs.Images.TexCube[me.Tex.ImageID].glTex.GlTarget)
			}
		}
	}
}

func (me *FxProc) SetMixWeight(weight float64) {
	me.weight = gl.Float(weight)
}

func (me *FxProc) Toggle() {
	me.Enabled = !me.Enabled
}

func (me *FxProc) Color_SetRgb(rgb ...gl.Float) *FxProc {
	me.Color.Rgb.Set(rgb...)
	return me
}

func (me *FxProc) Tex_SetImageID(imageID int) *FxProc {
	me.Tex.ImageID = imageID
	return me
}

func (me *FxProc) IsTex() bool {
	return me.IsTex2D() || me.IsTexCube()
}

//	Used for FxEffect.Procs and Camera.Rendering.FxProcs.
type FxProcs []FxProc

//	Disables the nth (0-based) FxProc with the specified procID,
//	or all FxProcs with the specified procID if n < 0.
//	The procID must be one of the Core.Render.Fx.KnownProcIDs.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxProcs) Disable(procID string, n int) {
	idx, all := -1, n < 0
	for i := 0; i < len(me); i++ {
		if me[i].procID == procID {
			if idx++; all || idx == n {
				me[i].Enabled = false
			}
			if !all {
				break
			}
		}
	}
}

//	Enables the nth (0-based) FxProc with the specified procID,
//	or all FxProcs with the specified procID if n < 0.
//	If me has no FxProc with the specified procID, appends a new one.
//	The procID must be one of the Core.Render.Fx.KnownProcIDs.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) Enable(procID string, n int) (proc *FxProc) {
	all, idx := n < 0, -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].procID == procID {
			if idx++; all || idx == n {
				(*me)[i].Enabled = true
				proc = &(*me)[i]
			}
		}
	}
	if proc == nil {
		*me = append(*me, FxProc{})
		proc = &(*me)[len(*me)-1]
		proc.init(procID, idx+1)
	}
	return
}

func (me FxProcs) EnsureLast(lastProcIDs ...string) {
	idx := len(me) - 1
	for _, id := range lastProcIDs {
		for i := 0; i < len(me); i++ {
			if me[i].procID == id {
				if i < idx {
					me[i], me[idx] = me[idx], me[i]
					idx--
				}
			}
		}
	}
}

//	Returns the nth (0-based) FxProc with the specified procID.
//	The procID must be one of the Core.Render.Fx.KnownProcIDs.
func (me FxProcs) Get(procID string, n int) (proc *FxProc) {
	idx := -1
	for i := 0; i < len(me); i++ {
		if me[i].procID == procID {
			if idx++; idx == n {
				proc = &me[i]
				return
			}
		}
	}
	return
}

//	Toggles the nth (0-based) FxProc with the specified procID,
//	or all FxProcs with the specified procID if n < 0.
//	If me has no FxProc with the specified procID, appends a new one.
//	The procID must be one of the Core.Render.Fx.KnownProcIDs.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) Toggle(procID string, n int) {
	var matcher ustr.Matcher
	matcher.AddPatterns(procID)
	idx, found, all := -1, false, n < 0
	for i := 0; i < len(*me); i++ {
		if matcher.IsMatch((*me)[i].procID) {
			if idx++; all || idx == n {
				(*me)[i].Toggle()
			}
			if found = true; !all {
				break
			}
		}
	}
	if (!found) && !matcher.HasWildcardPatterns() {
		*me = append(*me, FxProc{})
		(*me)[len(*me)-1].init(procID, idx+1)
	}
}











//#begin-gt -gen-fx-procs.gt GT_MULT_SEP:, N:Tex2D,TexCube,Orangify,Grayscale,Coords,Color,Gamma

func (me *FxProc) IsTex2D() bool {
	return me.procID == "Tex2D"
}

//	Convenience short-hand for me.Disable("Tex2D", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxProcs) DisableTex2D(n int) {
	me.Disable("Tex2D", n)
}

//	Convenience short-hand for me.Enable("Tex2D", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) EnableTex2D(n int) *FxProc {
	return me.Enable("Tex2D", n)
}

//	Convenience short-hand for me.Get("Tex2D", n).
func (me FxProcs) GetTex2D(n int) *FxProc {
	return me.Get("Tex2D", n)
}

//	Convenience short-hand for me.Toggle("Tex2D", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) ToggleTex2D(n int) {
	me.Toggle("Tex2D", n)
}

func (me *FxProc) IsTexCube() bool {
	return me.procID == "TexCube"
}

//	Convenience short-hand for me.Disable("TexCube", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxProcs) DisableTexCube(n int) {
	me.Disable("TexCube", n)
}

//	Convenience short-hand for me.Enable("TexCube", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) EnableTexCube(n int) *FxProc {
	return me.Enable("TexCube", n)
}

//	Convenience short-hand for me.Get("TexCube", n).
func (me FxProcs) GetTexCube(n int) *FxProc {
	return me.Get("TexCube", n)
}

//	Convenience short-hand for me.Toggle("TexCube", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) ToggleTexCube(n int) {
	me.Toggle("TexCube", n)
}

func (me *FxProc) IsOrangify() bool {
	return me.procID == "Orangify"
}

//	Convenience short-hand for me.Disable("Orangify", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxProcs) DisableOrangify(n int) {
	me.Disable("Orangify", n)
}

//	Convenience short-hand for me.Enable("Orangify", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) EnableOrangify(n int) *FxProc {
	return me.Enable("Orangify", n)
}

//	Convenience short-hand for me.Get("Orangify", n).
func (me FxProcs) GetOrangify(n int) *FxProc {
	return me.Get("Orangify", n)
}

//	Convenience short-hand for me.Toggle("Orangify", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) ToggleOrangify(n int) {
	me.Toggle("Orangify", n)
}

func (me *FxProc) IsGrayscale() bool {
	return me.procID == "Grayscale"
}

//	Convenience short-hand for me.Disable("Grayscale", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxProcs) DisableGrayscale(n int) {
	me.Disable("Grayscale", n)
}

//	Convenience short-hand for me.Enable("Grayscale", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) EnableGrayscale(n int) *FxProc {
	return me.Enable("Grayscale", n)
}

//	Convenience short-hand for me.Get("Grayscale", n).
func (me FxProcs) GetGrayscale(n int) *FxProc {
	return me.Get("Grayscale", n)
}

//	Convenience short-hand for me.Toggle("Grayscale", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) ToggleGrayscale(n int) {
	me.Toggle("Grayscale", n)
}

func (me *FxProc) IsCoords() bool {
	return me.procID == "Coords"
}

//	Convenience short-hand for me.Disable("Coords", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxProcs) DisableCoords(n int) {
	me.Disable("Coords", n)
}

//	Convenience short-hand for me.Enable("Coords", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) EnableCoords(n int) *FxProc {
	return me.Enable("Coords", n)
}

//	Convenience short-hand for me.Get("Coords", n).
func (me FxProcs) GetCoords(n int) *FxProc {
	return me.Get("Coords", n)
}

//	Convenience short-hand for me.Toggle("Coords", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) ToggleCoords(n int) {
	me.Toggle("Coords", n)
}

func (me *FxProc) IsColor() bool {
	return me.procID == "Color"
}

//	Convenience short-hand for me.Disable("Color", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxProcs) DisableColor(n int) {
	me.Disable("Color", n)
}

//	Convenience short-hand for me.Enable("Color", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) EnableColor(n int) *FxProc {
	return me.Enable("Color", n)
}

//	Convenience short-hand for me.Get("Color", n).
func (me FxProcs) GetColor(n int) *FxProc {
	return me.Get("Color", n)
}

//	Convenience short-hand for me.Toggle("Color", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) ToggleColor(n int) {
	me.Toggle("Color", n)
}

func (me *FxProc) IsGamma() bool {
	return me.procID == "Gamma"
}

//	Convenience short-hand for me.Disable("Gamma", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxProcs) DisableGamma(n int) {
	me.Disable("Gamma", n)
}

//	Convenience short-hand for me.Enable("Gamma", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) EnableGamma(n int) *FxProc {
	return me.Enable("Gamma", n)
}

//	Convenience short-hand for me.Get("Gamma", n).
func (me FxProcs) GetGamma(n int) *FxProc {
	return me.Get("Gamma", n)
}

//	Convenience short-hand for me.Toggle("Gamma", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxProcs) ToggleGamma(n int) {
	me.Toggle("Gamma", n)
}

//#end-gt
