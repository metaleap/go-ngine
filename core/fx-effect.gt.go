package core

import (
	ustr "github.com/metaleap/go-util/str"
)

//	Declares the visual appearance of a surface.
//	An FxEffect can be reused for multiple surfaces, it is bound to geometry via an FxMaterial.
type FxEffect struct {
	ID int

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
}

func (me *FxEffect) dispose() {
}

func (me *FxEffect) init() {
	me.uberPnames = make(map[string]string, len(Core.Rendering.KnownTechniques))
}

func (me *FxEffect) UpdateRoutine() {
	var (
		id, pid string
		idx, i  int
		buf     ustr.Buffer
	)
	if len(me.KeepOpsLast) > 0 {
		idx = len(me.Ops) - 1
		for _, id = range me.KeepOpsLast {
			for i, me.tmpOp = range me.Ops {
				if me.tmpOp.ProcID() == id {
					if i < idx {
						me.Ops[i], me.Ops[idx] = me.Ops[idx], me.Ops[i]
						idx--
					}
				}
			}
		}
	}

	ops, x, counts := me.Ops, len(me.OpsX) > 0, map[string]int{}
doOps:
	for _, me.tmpOp = range ops {
		if me.tmpOp.Enabled() {
			pid = me.tmpOp.ProcID()
			buf.Write("_%s", pid)
			i = counts[pid]
			me.tmpOp.setProcIndex(i)
			counts[pid] = i + 1
		}
	}
	if x {
		x, ops = false, me.OpsX
		goto doOps
	}

	me.uberName = buf.String()
	for id, _ = range Core.Rendering.KnownTechniques {
		me.uberPnames[id] = strf("uber_%s%s", id, me.uberName)
	}
	thrRend.curEffect = nil
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

//#begin-gt -gen-lib.gt T:FxEffect L:Effects

//	Only used for Core.Libs.Effects.
type FxEffectLib []FxEffect

func (me *FxEffectLib) AddNew() (ref *FxEffect) {
	id := -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			id = i
			break
		}
	}
	if id < 0 {
		if id = len(*me); id == cap(*me) {
			nu := make(FxEffectLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, FxEffect{})
	}
	ref = &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *FxEffectLib) Compact() {
	var (
		before, after []FxEffect
		ref           *FxEffect
		oldID, i      int
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			before, after = (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	changed := make(map[int]int, len(*me))
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID != i {
			ref = &(*me)[i]
			oldID, ref.ID = ref.ID, i
			changed[oldID] = i
		}
	}
	if len(changed) > 0 {
		me.onFxEffectIDsChanged(changed)
		Options.Libs.OnIDsChanged.Effects(changed)
	}
}

func (me *FxEffectLib) ctor() {
	*me = make(FxEffectLib, 0, Options.Libs.InitialCap)
}

func (me *FxEffectLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me FxEffectLib) Get(id int) (ref *FxEffect) {
	if id > -1 && id < len(me) {
		if ref = &me[id]; ref.ID != id {
			ref = nil
		}
	}
	return
}

func (me FxEffectLib) Has(id int) (has bool) {
	if id > -1 && id < len(me) {
		has = me[id].ID == id
	}
	return
}

func (me FxEffectLib) Ok(id int) bool {
	return me[id].ID > -1
}

func (me FxEffectLib) Remove(fromID, num int) {
	if l := len(me); fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onFxEffectIDsChanged(changed)
		Options.Libs.OnIDsChanged.Effects(changed)
	}
}

func (me FxEffectLib) Walk(on func(ref *FxEffect)) {
	for id := 0; id < len(me); id++ {
		if me[id].ID > -1 {
			on(&me[id])
		}
	}
}

//#end-gt
