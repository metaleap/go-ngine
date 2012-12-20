package assets

//	References a resource by its scoped identifier (Sid).
type RefSid string

//	Returns its current value.
func (me RefSid) S() string {
	return string(me)
}

//	Modifies its current value.
func (me *RefSid) Set(v string) {
	*me = RefSid(v)
}
