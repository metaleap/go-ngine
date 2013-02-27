package core

//	A Model is a parameterized instantiation of its parent Mesh geometry
//	with unique appearance, material or other properties.
//	
//	Each Mesh provides at least one Model, the "default model" (with ID ""),
//	accessible via someMesh.Models.Default(). To create new models for a Mesh,
//	call someMesh.Models["sourceModelID"].Clone("newModelID").
type Model struct {
	ID    int
	MatID int
	Name  string
}

func (me *Model) dispose() {
}

func (me *Model) init() {
	me.MatID = -1
	return
}

//	Creates a copy of me and adds it to the parent Mesh's Models
//	hash-table under the specified newModelID.
func (me *Model) Clone() (clonedModel *Model) {
	clonedModel = Core.Libs.Models.AddNew()
	clonedModel.MatID = me.MatID
	return
}

//#begin-gt -gen-lib.gt T:Model L:Core.Libs.Models

//	Only used for Core.Libs.Models
type ModelLib []Model

func (me *ModelLib) AddNew() (ref *Model) {
	id := -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			id = i
			break
		}
	}
	if id < 0 {
		if id = len(*me); id == cap(*me) {
			nu := make(ModelLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, Model{})
	}
	ref = &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *ModelLib) Compact() {
	var (
		before, after []Model
		ref           *Model
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
		if ref = &(*me)[i]; ref.ID != i {
			oldID, ref.ID = ref.ID, i
			changed[oldID] = i
		}
	}
	if len(changed) > 0 {
		me.onModelIDsChanged(changed)
	}
}

func (me *ModelLib) ctor() {
	*me = make(ModelLib, 0, Options.Libs.InitialCap)
}

func (me *ModelLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me ModelLib) Get(id int) (ref *Model) {
	if id > -1 && id < len(me) {
		if ref = &me[id]; ref.ID != id {
			ref = nil
		}
	}
	return
}

func (me ModelLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me ModelLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me ModelLib) Remove(fromID, num int) {
	if l := len(me); fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onModelIDsChanged(changed)
	}
}

func (me ModelLib) Walk(on func(ref *Model)) {
	for id := 0; id < len(me); id++ {
		if me[id].ID > -1 {
			on(&me[id])
		}
	}
}

//#end-gt
