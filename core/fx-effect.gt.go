package core

type fxProc struct {
	FuncName string
}

func newFxProc(name string) (me *fxProc) {
	me = &fxProc{FuncName: "fx_" + name}
	return
}

type FxProc struct {
	Enabled bool
	ProcID  string

	index int
}

type FxRoutine struct {
	Procs []*FxProc

	pname string
}

func NewFxRoutine(procIDs ...string) (me *FxRoutine) {
	me = &FxRoutine{}
	for _, procID := range procIDs {
		me.Procs = append(me.Procs, &FxProc{Enabled: true, ProcID: procID})
	}
	me.Update()
	return
}

func (me *FxRoutine) Update() {
	counts := map[string]int{}
	for _, proc := range me.Procs {
		if proc.Enabled {
			proc.index = counts[proc.ProcID]
			counts[proc.ProcID] = proc.index + 1
			me.pname += ("_" + proc.ProcID)
		}
	}
}

type FxEffect struct {
	OldDiffuse *FxColorOrTexture

	Routine *FxRoutine
}

func (me *FxEffect) dispose() {
}

func (me *FxEffect) init() {
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
