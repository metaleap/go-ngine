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

	OpsX FxOps

	KeepOpsLast []string

	tmpOp      FxOp
	uberName   string
	uberPnames map[string]string

	upd struct {
		i, idx  int
		id, pid string
		counts  map[string]int
		buf     ustr.Buffer
	}
}

func (me *FxEffect) dispose() {
}

func (me *FxEffect) init() {
	me.uberPnames = make(map[string]string, len(Core.Rendering.KnownTechniques))
	me.upd.counts = map[string]int{}
}

func (me *FxEffect) UpdateRoutine() {
	if me.upd.buf.Reset(); len(me.KeepOpsLast) > 0 {
		me.upd.idx = len(me.Ops) - 1
		for _, me.upd.id = range me.KeepOpsLast {
			for me.upd.i, me.tmpOp = range me.Ops {
				if me.tmpOp.ProcID() == me.upd.id {
					if me.upd.i < me.upd.idx {
						me.Ops[me.upd.i], me.Ops[me.upd.idx] = me.Ops[me.upd.idx], me.Ops[me.upd.i]
						me.upd.idx--
					}
				}
			}
		}
	}
	for me.upd.pid, _ = range me.upd.counts {
		me.upd.counts[me.upd.pid] = 0
	}
	for _, me.tmpOp = range me.Ops {
		me.updateRoutineOp()
	}
	for _, me.tmpOp = range me.OpsX {
		me.updateRoutineOp()
	}
	me.uberName = me.upd.buf.String()
	for me.upd.id, _ = range Core.Rendering.KnownTechniques {
		me.uberPnames[me.upd.id] = strf("uber_%s%s", me.upd.id, me.uberName)
	}
	thrRend.curEffect = nil
}

func (me *FxEffect) updateRoutineOp() {
	if me.tmpOp.Enabled() {
		me.upd.pid = me.tmpOp.ProcID()
		me.upd.buf.Write("_%s", me.upd.pid)
		me.upd.i = me.upd.counts[me.upd.pid]
		me.tmpOp.setProcIndex(me.upd.i)
		me.upd.counts[me.upd.pid] = me.upd.i + 1
	}
}

func (me *FxEffect) use() {
	for _, me.tmpOp = range me.Ops {
		if me.tmpOp.Enabled() {
			me.tmpOp.use()
		}
	}
	for _, me.tmpOp = range me.OpsX {
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
