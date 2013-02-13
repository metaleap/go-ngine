package core

import (
	ugl "github.com/go3d/go-opengl/util"
)

type FxOp interface {
	init(string)

	Disable()
	Disabled() bool
	Enable()
	Enabled() bool
	ProcID() string
	Toggle()
}

func NewFxOp(procID string) (me FxOp) {
	switch procID {
	case "Grayscale":
		me = &FxOpGrayscale{}
	case "RedTest":
		me = &FxOpRedTest{}
	case "Tex2D":
		me = &FxOpTex2D{}
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

func (me *fxOpBase) Disable() {
	me.disabled = true
}

func (me *fxOpBase) Disabled() bool {
	return me.disabled
}

func (me *fxOpBase) Enable() {
	me.disabled = false
}

func (me *fxOpBase) Enabled() bool {
	return !me.disabled
}

func (me *fxOpBase) ProcID() string {
	return me.procID
}

func (me *fxOpBase) Toggle() {
	me.disabled = !me.disabled
}

type FxOpGrayscale struct {
	fxOpBase
}

type FxOpRedTest struct {
	fxOpBase
}

type FxOpTex2D struct {
	fxOpBase

	ImageID string
	Sampler *ugl.Sampler
}

type FxOps []FxOp

func (me FxOps) Disable(procID string, index int) {
	idx, all := -1, index < 0
	for _, op := range me {
		if op.ProcID() == procID {
			if idx++; all || idx == index {
				op.Disable()
			}
			if !all {
				break
			}
		}
	}
}

func (me *FxOps) Enable(procID string, index int) (fxOp FxOp) {
	var op FxOp
	found, all, idx := false, index < 0, -1
	for _, op = range *me {
		if op.ProcID() == procID {
			if idx++; all || idx == index {
				op.Enable()
			}
			if fxOp, found = op, true; !all {
				break
			}
		}
	}
	if !found {
		fxOp = NewFxOp(procID)
		*me = append(*me, fxOp)
	}
	return
}

func (me FxOps) Get(procID string, index int) (op FxOp) {
	idx := -1
	for _, op = range me {
		if op.ProcID() == procID {
			if idx++; idx == index {
				return
			}
		}
	}
	op = nil
	return
}

func (me *FxOps) Toggle(procID string, index int) {
	idx, found, all := -1, false, index < 0
	for _, op := range *me {
		if op.ProcID() == procID {
			if idx++; all || idx == index {
				op.Toggle()
			}
			if found = true; !all {
				break
			}
		}
	}
	if !found {
		me.Enable(procID, index)
	}
}

//#begin-gt -gen-fx-ops.gt GT_MULT_SEP:, N:Tex2D,RedTest,Grayscale

//	Convenience short-hand for me.Disable("Tex2D", index)
func (me FxOps) DisableTex2D(index int) {
	me.Disable("Tex2D", index)
}

//	Convenience short-hand for me.Enable("Tex2D", index)
func (me *FxOps) EnableTex2D(index int) *FxOpTex2D {
	return me.Enable("Tex2D", index).(*FxOpTex2D)
}

//	Convenience short-hand for me.Get("Tex2D", index)
func (me FxOps) GetTex2D(index int) *FxOpTex2D {
	return me.Get("Tex2D", index).(*FxOpTex2D)
}

//	Convenience short-hand for me.Toggle("Tex2D", index)
func (me *FxOps) ToggleTex2D(index int) {
	me.Toggle("Tex2D", index)
}



//	Convenience short-hand for me.Disable("RedTest", index)
func (me FxOps) DisableRedTest(index int) {
	me.Disable("RedTest", index)
}

//	Convenience short-hand for me.Enable("RedTest", index)
func (me *FxOps) EnableRedTest(index int) *FxOpRedTest {
	return me.Enable("RedTest", index).(*FxOpRedTest)
}

//	Convenience short-hand for me.Get("RedTest", index)
func (me FxOps) GetRedTest(index int) *FxOpRedTest {
	return me.Get("RedTest", index).(*FxOpRedTest)
}

//	Convenience short-hand for me.Toggle("RedTest", index)
func (me *FxOps) ToggleRedTest(index int) {
	me.Toggle("RedTest", index)
}



//	Convenience short-hand for me.Disable("Grayscale", index)
func (me FxOps) DisableGrayscale(index int) {
	me.Disable("Grayscale", index)
}

//	Convenience short-hand for me.Enable("Grayscale", index)
func (me *FxOps) EnableGrayscale(index int) *FxOpGrayscale {
	return me.Enable("Grayscale", index).(*FxOpGrayscale)
}

//	Convenience short-hand for me.Get("Grayscale", index)
func (me FxOps) GetGrayscale(index int) *FxOpGrayscale {
	return me.Get("Grayscale", index).(*FxOpGrayscale)
}

//	Convenience short-hand for me.Toggle("Grayscale", index)
func (me *FxOps) ToggleGrayscale(index int) {
	me.Toggle("Grayscale", index)
}

//#end-gt
