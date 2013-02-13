package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

//	Implemented by specialized types such as FxOpTex2D, FxOpGrayscale etc.
//	Those are created and initializes exclusively through FxOps.Enable("{procID}") or the specialized FxOps.EnableFoo() methods.
type FxOp interface {
	init(string)

	nifyFunc() bool
	nifyUnif(fullName, typeName string) bool
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

	//	Toggles this FxOp.
	//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
	Toggle()
}

func newFxOp(procID string) (me FxOp) {
	switch procID {
	case "Grayscale":
		me = &FxOpGrayscale{}
	case "Orangify":
		me = &FxOpOrangify{}
	case "Tex2D":
		me = &FxOpTex2D{}
	case "Colored":
		me = &FxOpColored{}
	}
	if me != nil {
		me.init(procID)
	}
	return
}

type fxOpBase struct {
	disabled bool
	procID   string
}

func (me *fxOpBase) init(procID string) {
	me.procID = procID
}

func (me *fxOpBase) nifyFunc() bool {
	return false
}

func (me *fxOpBase) nifyUnif(fullName, typeName string) bool {
	return false
}

func (me *fxOpBase) use() {
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

//	Toggles this FxOp.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *fxOpBase) Toggle() {
	me.disabled = !me.disabled
}

//	Yields the grayscale of the current pixel color.
type FxOpGrayscale struct {
	fxOpBase
}

//	A very simple color "effect". Tints the current color orange-red.
type FxOpOrangify struct {
	fxOpBase
}

//	Colors geometry based on texture coordinates.
type FxOpColored struct {
	fxOpBase
}

//	Samples from a 2D texture.
type FxOpTex2D struct {
	fxOpBase

	bindTex func(string)

	//	The ID (in Core.Libs.Images.I2D) of the FxImage2D to be sampled.
	ImageID string

	//	The sampler to be used.
	//	Defaults to &Core.Rendering.Fx.Samplers.FullFilteringRepeat.
	Sampler *ugl.Sampler
}

func (me *FxOpTex2D) init(s string) {
	me.fxOpBase.init(s)
	me.Sampler = &Core.Rendering.Fx.Samplers.FullFilteringRepeat
	me.bindTex = func(imgID string) {
		Core.Libs.Images.I2D[imgID].glTex.Bind()
	}
}

func (me *FxOpTex2D) nifyFunc() bool {
	return true
}

func (me *FxOpTex2D) nifyUnif(fullName, typeName string) bool {
	return typeName == "sampler2D"
}

func (me *FxOpTex2D) use() {
	me.Sampler.Bind(0)
	me.bindTex(me.ImageID)
	gl.Uniform1i(thrRend.curProg.UnifLocs["uni_sampler2D_Tex2D"], 0)
}

//	Only used for FxEffect.Ops.
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
	found, all, idx := false, n < 0, -1
	for _, op = range *me {
		if op.ProcID() == procID {
			if idx++; all || idx == n {
				op.Enable()
			}
			if fxOp, found = op, true; !all {
				break
			}
		}
	}
	if !found {
		if fxOp = newFxOp(procID); fxOp != nil {
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

func (me *FxOps) SwapAll(procID1, procID2 string) {
	for i, op := range *me {
		switch op.ProcID() {
		case procID1:
			(*me)[i] = newFxOp(procID2)
		case procID2:
			(*me)[i] = newFxOp(procID1)
		}
	}
}

//	Toggles the nth (0-based) FxOp with the specified procID,
//	or all FxOps with the specified procID if n < 0.
//	If me has no FxOp with the specified procID, appends a new one.
//	The procID must be one of the Core.Rendering.Fx.KnownProcIDs.
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) Toggle(procID string, n int) {
	var op FxOp
	idx, found, all := -1, false, n < 0
	for _, op = range *me {
		if op.ProcID() == procID {
			if idx++; all || idx == n {
				op.Toggle()
			}
			if found = true; !all {
				break
			}
		}
	}
	if !found {
		if op = newFxOp(procID); op != nil {
			*me = append(*me, op)
		}
	}
}

//#begin-gt -gen-fx-ops.gt GT_MULT_SEP:, N:Tex2D,Orangify,Grayscale,Colored

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
func (me FxOps) GetTex2D(n int) *FxOpTex2D {
	return me.Get("Tex2D", n).(*FxOpTex2D)
}

//	Convenience short-hand for me.Toggle("Tex2D", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleTex2D(n int) {
	me.Toggle("Tex2D", n)
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
func (me FxOps) GetOrangify(n int) *FxOpOrangify {
	return me.Get("Orangify", n).(*FxOpOrangify)
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
func (me FxOps) GetGrayscale(n int) *FxOpGrayscale {
	return me.Get("Grayscale", n).(*FxOpGrayscale)
}

//	Convenience short-hand for me.Toggle("Grayscale", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleGrayscale(n int) {
	me.Toggle("Grayscale", n)
}



//	Convenience short-hand for me.Disable("Colored", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me FxOps) DisableColored(n int) {
	me.Disable("Colored", n)
}

//	Convenience short-hand for me.Enable("Colored", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) EnableColored(n int) *FxOpColored {
	return me.Enable("Colored", n).(*FxOpColored)
}

//	Convenience short-hand for me.Get("Colored", n).
func (me FxOps) GetColored(n int) *FxOpColored {
	return me.Get("Colored", n).(*FxOpColored)
}

//	Convenience short-hand for me.Toggle("Colored", n).
//	For this change to be applied, call FxEffect.UpdateRoutine() subsequently.
func (me *FxOps) ToggleColored(n int) {
	me.Toggle("Colored", n)
}

//#end-gt
