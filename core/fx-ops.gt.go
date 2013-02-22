package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
	ustr "github.com/metaleap/go-util/str"
)

func fxKnownProcIDs() []string {
	return []string{"Color", "Tex2D", "TexCube", "Grayscale", "Orangify", "Coords", "Gamma"}
}

//	Implemented by specialized types such as FxOpTex2D, FxOpGrayscale etc.
//	Those are created and initializes exclusively through FxOps.Enable("{procID}") or the specialized FxOps.EnableFoo() methods.
type FxOp interface {
	init(string, int)

	fxImageIDs() []string

	qualifiers(inout string) string

	use()

	//	Disables this FxOp.
	//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
	Disable()

	//	Returns whether this FxOp is currently disabled.
	Disabled() bool

	//	Enables this FxOp.
	//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
	Enable()

	//	Returns whether this FxOp is currently enabled.
	Enabled() bool

	//	The procID of this FxOp. This is one of the Core.Rendering.Fx.KnownProcIDs.
	//	For example, "Tex2D" for an FxOpTex2D, "Grayscale" for an FxOpGrayscale etc.
	ProcID() string

	ProcIndex() int

	SetMixWeight(weight float64)

	setProcIndex(index int)

	//	Toggles this FxOp.
	//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
	Toggle()

	unifName(t, n string) string
}

func newFxOp(procID string, procIndex int) (me FxOp) {
	switch procID {
	case "Grayscale":
		me = &FxOpGrayscale{}
	case "Orangify":
		me = &FxOpOrangify{}
	case "Tex2D":
		me = &FxOpTex2D{}
	case "TexCube":
		me = &FxOpTexCube{}
	case "Coords":
		me = &FxOpCoords{}
	case "Color":
		me = &FxOpColor{}
	case "Gamma":
		me = &FxOpGamma{}
	}
	if me != nil {
		me.init(procID, procIndex)
	}
	return
}

type fxOpBase struct {
	disabled  bool
	procID    string
	procIndex int
	weight    gl.Float
	unifNames map[string]string
}

func (me *fxOpBase) init(procID string, procIndex int) {
	me.procID = procID
	me.setProcIndex(procIndex)
	me.weight = 1
}

func (me *fxOpBase) fxImageIDs() (ids []string) {
	return
}

func (me *fxOpBase) qualifiers(inout string) string {
	return ""
}

func (me *fxOpBase) unifName(t, n string) (un string) {
	if un = me.unifNames[t+"_"+n]; len(un) == 0 {
		un = strf("uni_%s_%s%d_%s", t, me.procID, me.procIndex, n)
		me.unifNames[t+"_"+n] = un
	}
	return
}

func (me *fxOpBase) use() {
	thrRend.curProg.Uniform1f(me.unifName("float", "MixWeight"), me.weight)
}

//	Disables this FxOp.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *fxOpBase) Disable() {
	me.disabled = true
}

//	Returns whether this FxOp is currently disabled.
func (me *fxOpBase) Disabled() bool {
	return me.disabled
}

//	Enables this FxOp.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *fxOpBase) Enable() {
	me.disabled = false
}

//	Returns whether this FxOp is currently enabled.
func (me *fxOpBase) Enabled() bool {
	return !me.disabled
}

//	The procID of this FxOp. This is one of the Core.Rendering.Fx.KnownProcIDs.
//	For example, "Tex2D" for an FxOpTex2D, "Grayscale" for an FxOpGrayscale etc.
func (me *fxOpBase) ProcID() string {
	return me.procID
}

func (me *fxOpBase) ProcIndex() int {
	return me.procIndex
}

func (me *fxOpBase) SetMixWeight(weight float64) {
	me.weight = gl.Float(weight)
}

func (me *fxOpBase) setProcIndex(index int) {
	me.unifNames = map[string]string{}
	me.procIndex = index
}

//	Toggles this FxOp.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *fxOpBase) Toggle() {
	me.disabled = !me.disabled
}

//	Yields the grayscale of the current pixel color.
type FxOpGrayscale struct {
	fxOpBase
}

//	pow(color, 1/2.2)
type FxOpGamma struct {
	fxOpBase
}

//	A very simple color "effect". Tints the current color orange-red.
type FxOpOrangify struct {
	fxOpBase
}

//	Colors geometry based on texture coordinates.
type FxOpCoords struct {
	fxOpBase
}

//	Colors geometry based on a specified color.
type FxOpColor struct {
	fxOpBase
	Rgb ugl.GlVec3
}

func (me *FxOpColor) SetRgb(rgb ...gl.Float) *FxOpColor {
	me.Rgb.Set(rgb...)
	return me
}

func (me *FxOpColor) use() {
	me.fxOpBase.use()
	thrRend.curProg.Uniform3fv(me.unifName("vec3", "Rgb"), 1, &me.Rgb[0])
}

type fxOpTexBase struct {
	fxOpBase

	//	The sampler to be used.
	Sampler *ugl.Sampler
}

func (me *fxOpTexBase) init(s string, i int) {
	me.fxOpBase.init(s, i)
	me.Sampler = &Core.Rendering.Fx.Samplers.FullFilteringRepeat
}

func (me *fxOpTexBase) use() {
	me.fxOpBase.use()
	thrRend.tmpSampler = me.Sampler
	Core.useSampler()
}

//	Samples from a 2D texture.
type FxOpTex2D struct {
	fxOpTexBase
	glTex *ugl.Texture2D
}

func (me *FxOpTex2D) init(s string, i int) {
	me.fxOpTexBase.init(s, i)
}

func (me *FxOpTex2D) SetImageID(imageID string) *FxOpTex2D {
	me.glTex = &Core.Libs.Images.Tex2D[imageID].glTex
	return me
}

func (me *FxOpTex2D) use() {
	me.fxOpTexBase.use()
	gl.ActiveTexture(gl.Enum(gl.TEXTURE0 + me.procIndex))
	me.glTex.Bind()
	thrRend.curProg.Uniform1i(me.unifName("sampler2D", "Img"), gl.Int(me.procIndex))
}

type FxOpTexCube struct {
	fxOpTexBase
	glTex *ugl.TextureCube
}

func (me *FxOpTexCube) init(s string, i int) {
	me.fxOpTexBase.init(s, i)
	me.Sampler = &Core.Rendering.Fx.Samplers.FullFilteringClamp
}

func (me *FxOpTexCube) SetImageID(imageID string) {
	me.glTex = &Core.Libs.Images.TexCubes[imageID].glTex
}

func (me *FxOpTexCube) use() {
	me.fxOpTexBase.use()
	me.glTex.Bind()
	thrRend.curProg.Uniform1i(me.unifName("samplerCube", "Img"), 0)
}

//	Used for FxEffect.Ops/OpsX and Camera.Rendering.FxOps.
type FxOps []FxOp

//	Disables the nth (0-based) FxOp with the specified procID,
//	or all FxOps with the specified procID if n < 0.
//	The procID must be one of the Core.Rendering.Fx.KnownProcIDs.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) Disable(procID string, n int) {
	idx, all := -1, n < 0
	for _, op := range me {
		if op.ProcID() == procID {
			if idx++; all || idx == n {
				op.Disable()
			}
			if !all {
				break
			}
		}
	}
}

//	Enables the nth (0-based) FxOp with the specified procID,
//	or all FxOps with the specified procID if n < 0.
//	If me has no FxOp with the specified procID, appends a new one.
//	The procID must be one of the Core.Rendering.Fx.KnownProcIDs.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) Enable(procID string, n int) (fxOp FxOp) {
	var op FxOp
	all, idx := n < 0, -1
	for _, op = range *me {
		if op.ProcID() == procID {
			if idx++; all || idx == n {
				op.Enable()
				fxOp = op
			}
		}
	}
	if fxOp == nil {
		if fxOp = newFxOp(procID, idx+1); fxOp != nil {
			*me = append(*me, fxOp)
		}
	}
	return
}

//	Returns the nth (0-based) FxOp with the specified procID.
//	The procID must be one of the Core.Rendering.Fx.KnownProcIDs.
func (me FxOps) Get(procID string, n int) (op FxOp) {
	idx := -1
	for _, op = range me {
		if op.ProcID() == procID {
			if idx++; idx == n {
				return
			}
		}
	}
	op = nil
	return
}

//	Toggles the nth (0-based) FxOp with the specified procID,
//	or all FxOps with the specified procID if n < 0.
//	If me has no FxOp with the specified procID, appends a new one.
//	The procID must be one of the Core.Rendering.Fx.KnownProcIDs.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) Toggle(procID string, n int) {
	var op FxOp
	matcher := ustr.NewMatcher(procID)
	idx, found, all := -1, false, n < 0
	for _, op = range *me {
		if matcher.IsMatch(op.ProcID()) {
			if idx++; all || idx == n {
				op.Toggle()
			}
			if found = true; !all {
				break
			}
		}
	}
	if (!found) && !matcher.HasWildcardPatterns() {
		if op = newFxOp(procID, idx+1); op != nil {
			*me = append(*me, op)
		}
	}
}

//#begin-gt -gen-fx-ops.gt GT_MULT_SEP:, N:Tex2D,TexCube,Orangify,Grayscale,Coords,Color,Gamma

//	Convenience short-hand for me.Disable("Tex2D", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) DisableTex2D(n int) {
	me.Disable("Tex2D", n)
}

//	Convenience short-hand for me.Enable("Tex2D", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) EnableTex2D(n int) *FxOpTex2D {
	return me.Enable("Tex2D", n).(*FxOpTex2D)
}

//	Convenience short-hand for me.Get("Tex2D", n).
func (me FxOps) GetTex2D(n int) (op *FxOpTex2D) {
	op, _ = me.Get("Tex2D", n).(*FxOpTex2D)
	return
}

//	Convenience short-hand for me.Toggle("Tex2D", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleTex2D(n int) {
	me.Toggle("Tex2D", n)
}

//	Convenience short-hand for me.Disable("TexCube", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) DisableTexCube(n int) {
	me.Disable("TexCube", n)
}

//	Convenience short-hand for me.Enable("TexCube", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) EnableTexCube(n int) *FxOpTexCube {
	return me.Enable("TexCube", n).(*FxOpTexCube)
}

//	Convenience short-hand for me.Get("TexCube", n).
func (me FxOps) GetTexCube(n int) (op *FxOpTexCube) {
	op, _ = me.Get("TexCube", n).(*FxOpTexCube)
	return
}

//	Convenience short-hand for me.Toggle("TexCube", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleTexCube(n int) {
	me.Toggle("TexCube", n)
}

//	Convenience short-hand for me.Disable("Orangify", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) DisableOrangify(n int) {
	me.Disable("Orangify", n)
}

//	Convenience short-hand for me.Enable("Orangify", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) EnableOrangify(n int) *FxOpOrangify {
	return me.Enable("Orangify", n).(*FxOpOrangify)
}

//	Convenience short-hand for me.Get("Orangify", n).
func (me FxOps) GetOrangify(n int) (op *FxOpOrangify) {
	op, _ = me.Get("Orangify", n).(*FxOpOrangify)
	return
}

//	Convenience short-hand for me.Toggle("Orangify", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleOrangify(n int) {
	me.Toggle("Orangify", n)
}

//	Convenience short-hand for me.Disable("Grayscale", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) DisableGrayscale(n int) {
	me.Disable("Grayscale", n)
}

//	Convenience short-hand for me.Enable("Grayscale", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) EnableGrayscale(n int) *FxOpGrayscale {
	return me.Enable("Grayscale", n).(*FxOpGrayscale)
}

//	Convenience short-hand for me.Get("Grayscale", n).
func (me FxOps) GetGrayscale(n int) (op *FxOpGrayscale) {
	op, _ = me.Get("Grayscale", n).(*FxOpGrayscale)
	return
}

//	Convenience short-hand for me.Toggle("Grayscale", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleGrayscale(n int) {
	me.Toggle("Grayscale", n)
}

//	Convenience short-hand for me.Disable("Coords", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) DisableCoords(n int) {
	me.Disable("Coords", n)
}

//	Convenience short-hand for me.Enable("Coords", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) EnableCoords(n int) *FxOpCoords {
	return me.Enable("Coords", n).(*FxOpCoords)
}

//	Convenience short-hand for me.Get("Coords", n).
func (me FxOps) GetCoords(n int) (op *FxOpCoords) {
	op, _ = me.Get("Coords", n).(*FxOpCoords)
	return
}

//	Convenience short-hand for me.Toggle("Coords", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleCoords(n int) {
	me.Toggle("Coords", n)
}

//	Convenience short-hand for me.Disable("Color", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) DisableColor(n int) {
	me.Disable("Color", n)
}

//	Convenience short-hand for me.Enable("Color", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) EnableColor(n int) *FxOpColor {
	return me.Enable("Color", n).(*FxOpColor)
}

//	Convenience short-hand for me.Get("Color", n).
func (me FxOps) GetColor(n int) (op *FxOpColor) {
	op, _ = me.Get("Color", n).(*FxOpColor)
	return
}

//	Convenience short-hand for me.Toggle("Color", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleColor(n int) {
	me.Toggle("Color", n)
}

//	Convenience short-hand for me.Disable("Gamma", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) DisableGamma(n int) {
	me.Disable("Gamma", n)
}

//	Convenience short-hand for me.Enable("Gamma", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) EnableGamma(n int) *FxOpGamma {
	return me.Enable("Gamma", n).(*FxOpGamma)
}

//	Convenience short-hand for me.Get("Gamma", n).
func (me FxOps) GetGamma(n int) (op *FxOpGamma) {
	op, _ = me.Get("Gamma", n).(*FxOpGamma)
	return
}

//	Convenience short-hand for me.Toggle("Gamma", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleGamma(n int) {
	me.Toggle("Gamma", n)
}

//#end-gt
