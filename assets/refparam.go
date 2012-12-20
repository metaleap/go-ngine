package assets

//	References a previously defined parameter.
type RefParam RefSid

//	Returns its current value.
func (me RefParam) S() string {
	return string(me)
}

//	Returns its current value as a RefSid.
func (me RefParam) Sr() RefSid {
	return RefSid(me)
}

//	Modifies its current value.
func (me *RefParam) Set(v string) {
	*me = RefParam(v)
}
