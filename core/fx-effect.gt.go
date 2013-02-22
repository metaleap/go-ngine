package core

import (
	ustr "github.com/metaleap/go-util/str"
)

type fxProc struct {
	FuncName string
}

func newFxProc(name string) (me *fxProc) {
	me = &fxProc{FuncName: "fx_" + name}
	return
}

//	Declares the visual appearance of a surface.
//	An FxEffect can be reused for multiple surfaces, it is bound to geometry via an FxMaterial.
type FxEffect struct {
	//	An ordered collection of all FxOps that make up this effect.
	//	When changing the ordering or disabling, enabling or toggling individual FxOps,
	//	you need to call the FxEffect.UpdateRoutine() method to reflect such changes.
	//	Other dynamic, individual FxOp-specific parameter changes (colors, image bindings, weights etc.pp.)
	//	do not require this.
	Ops FxOps

	KeepOpsLast []string

	tmpOp      FxOp
	uberName   string
	uberPnames map[string]string
}

func (me *FxEffect) dispose() {
}

func (me *FxEffect) init() {
	me.uberPnames = make(map[string]string, len(Core.Rendering.KnownTechniques))
	me.KeepOpsLast = []string{"Gamma"}
}

func (me *FxEffect) UpdateRoutine() {
	var (
		buf    ustr.Buffer
		id     string
		i, idx int
		op     FxOp
	)

	idx = len(me.Ops) - 1
	for _, id = range me.KeepOpsLast {
		for i, op = range me.Ops {
			if op.ProcID() == id {
				if i < idx {
					me.Ops[i], me.Ops[idx] = me.Ops[idx], me.Ops[i]
					idx--
				}
			}
		}
	}

	counts := map[string]int{}
	for _, op = range me.Ops {
		if op.Enabled() {
			buf.Write("_%s", op.ProcID())
			i = counts[op.ProcID()]
			op.setProcIndex(i)
			counts[op.ProcID()] = i + 1
		}
	}
	me.uberName = buf.String()
	for id, _ = range Core.Rendering.KnownTechniques {
		me.uberPnames[id] = strf("uber_%s%s", id, me.uberName)
	}
	thrRend.curEffect, thrRend.tmpEffect = nil, nil
}

func (me *FxEffect) use() {
	for _, me.tmpOp = range me.Ops {
		if me.tmpOp.Enabled() {
			me.tmpOp.use()
		}
	}
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
