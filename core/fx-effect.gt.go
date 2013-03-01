package core

import (
	ustr "github.com/metaleap/go-util/str"
)

//	Declares the visual appearance of a surface.
//	An FxEffect can be reused for multiple surfaces, it is bound to geometry via an FxMaterial.
type FxEffect struct {
	//	An ordered collection of all FxProcs that make up this effect.
	//	When changing the ordering or disabling, enabling or toggling individual FxProcs,
	//	you need to call the FxEffect.UpdateRoutine() method to reflect such changes.
	//	All other dynamic, individual FxProc-specific parameter changes
	//	(colors, image bindings, weights etc.pp.) do not require this.
	FxProcs

	ID int

	Ext FxProcs

	KeepProcIDsLast []string

	uberName   string
	uberPnames map[string]string
}

func (me *FxEffect) dispose() {
}

func (me *FxEffect) init() {
	me.uberPnames = make(map[string]string, len(Core.Render.KnownTechniques))
	me.FxProcs = make(FxProcs, 0, 4)
}

func (me *FxEffect) UpdateRoutine() {
	var (
		i   int
		buf ustr.Buffer
		id  string
	)
	if len(me.KeepProcIDsLast) > 0 {
		me.FxProcs.EnsureLast(me.KeepProcIDsLast...)
		me.Ext.EnsureLast(me.KeepProcIDsLast...)
	}

	ops, ext, counts := me.FxProcs, len(me.Ext) > 0, make(map[string]int, len(me.FxProcs)+len(me.Ext))
doOps:
	for o := 0; o < len(ops); o++ {
		if ops[o].Enabled {
			buf.Write("_%s", ops[o].procID)
			i = counts[ops[o].procID]
			ops[o].setProcIndex(i)
			counts[ops[o].procID] = i + 1
		}
	}
	if ext {
		ext, ops = false, me.Ext
		goto doOps
	}

	me.uberName = buf.String()
	for id, _ = range Core.Render.KnownTechniques {
		me.uberPnames[id] = strf("uber_%s%s", id, me.uberName)
	}
	thrRend.curEffect = nil
}

func (me *FxEffect) use() {
	me.useProcs(me.FxProcs)
	me.useProcs(me.Ext)
}

func (me *FxEffect) useProcs(ops FxProcs) {
	for i := 0; i < len(ops); i++ {
		if ops[i].Enabled {
			ops[i].use()
		}
	}
}

//#begin-gt -gen-lib.gt T:FxEffect L:Core.Libs.Effects

//	Only used for Core.Libs.Effects
type FxEffectLib []FxEffect

func (me *FxEffectLib) AddNew() (id int) {
	id = -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			id = i
			break
		}
	}
	if id == -1 {
		if id = len(*me); id == cap(*me) {
			nu := make(FxEffectLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, FxEffect{})
	}
	ref := &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *FxEffectLib) Compact() {
	var (
		before, after []FxEffect
		ref           *FxEffect
		oldID, i      int
		compact       bool
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			compact, before, after = true, (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	if compact {
		changed := make(map[int]int, len(*me))
		for i = 0; i < len(*me); i++ {
			if ref = &(*me)[i]; ref.ID != i {
				oldID, ref.ID = ref.ID, i
				changed[oldID] = i
			}
		}
		if len(changed) > 0 {
			me.onFxEffectIDsChanged(changed)
		}
	}
}

func (me *FxEffectLib) init() {
	*me = make(FxEffectLib, 0, Options.Libs.InitialCap)
}

func (me *FxEffectLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me FxEffectLib) get(id int) (ref *FxEffect) {
	if me.IsOk(id) {
		ref = &me[id]
	}
	return
}

func (me FxEffectLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me FxEffectLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me FxEffectLib) Remove(fromID, num int) {
	if l := len(me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onFxEffectIDsChanged(changed)
	}
}

func (me FxEffectLib) Walk(on func(ref *FxEffect)) {
	for id := 0; id < len(me); id++ {
		if me.Ok(id) {
			on(&me[id])
		}
	}
}

//#end-gt
