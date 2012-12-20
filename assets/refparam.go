package assets

//	References a previously defined parameter.
type RefParam RefSid

//	Returns the Sid currently referenced by me.
func (me RefParam) S() string {
	return string(me)
}

//	Modifies the Sid currently referenced by me.
func (me *RefParam) SetParamRef(v string) {
	*me = RefParam(v)
}

//	Returns the Sid currently referenced by me as a RefSid.
func (me RefParam) Sr() RefSid {
	return RefSid(me)
}
