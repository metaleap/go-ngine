package assets

//	References a resource by its scoped identifier (Sid).
type RefSid string

//	Returns the Sid currently referenced by me.
func (me RefSid) S() string {
	return string(me)
}

//	Modifies the Sid currently referenced by me.
func (me *RefSid) SetSidRef(v string) {
	*me = RefSid(v)
}
