package core

import (
	"strings"
)

type fxProc struct {
	FuncName string
}

func newFxProc(name string) (me *fxProc) {
	me = &fxProc{FuncName: "fx_" + name}
	return
}

type FxEffect struct {
	Ops FxOps

	uberProcIDs []string
	uberName    string
}

func (me *FxEffect) dispose() {
}

func (me *FxEffect) init() {
}

func (me *FxEffect) UpdateRoutine() {
	me.uberProcIDs = nil
	for _, op := range me.Ops {
		if op.Enabled() {
			me.uberProcIDs = append(me.uberProcIDs, op.ProcID())
		}
	}
	me.uberName = strings.Join(me.uberProcIDs, "_")
	thrRend.curEffect = nil
}

//#begin-gt -gen-lib.gt T:FxEffect

//	Initializes and returns a new FxEffect with default parameters.
func NewFxEffect() (me *FxEffect) {
	me = &FxEffect{}
	me.init()
	return
}

//	A hash-table of FxEffects associated by IDs. Only for use in Core.Libs.
type LibFxEffects map[string]*FxEffect

//	Creates and initializes a new FxEffect with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibFxEffects) AddNew(id string) (obj *FxEffect) {
	obj = NewFxEffect()
	me[id] = obj
	return
}

func (me *LibFxEffects) ctor() {
	*me = LibFxEffects{}
}

func (me *LibFxEffects) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}

func (me LibFxEffects) Remove(id string) {
	if obj := me[id]; obj != nil {
		obj.dispose()
	}
	delete(me, id)
}

//#end-gt
