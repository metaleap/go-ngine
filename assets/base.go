package assets

type base struct {
	ID string
}

	func (me *base) init (id string) {
		me.ID = id
	}

type baseDef struct {
	base
}

type baseInst struct {
	base
}

type baseLib struct {
	base
}
