package assets

import (
	"strings"
)

//	References a previously defined parameter.
type RefParam struct {
	RefSid
}

//	Convenience short-hand for me.RefSid.SetSidRef(sidRef)
func (me *RefParam) SetParamRef(sidRef string) {
	me.RefSid.SetSidRef(sidRef)
}

//	Creates and returns a new RefParam initialized with the specified paramRef.
func NewRefParam(paramRef string) (rs *RefParam) {
	rs = &RefParam{}
	rs.SetParamRef(paramRef)
	return
}

//	Resolves a Sid path.
type RefSidResolver interface {
	//	The returned val is always a pointer: so val may be a *SidFloat but it will never be a SidFloat.
	//	If no match is found for the full path, should always return nil, instead of a partial-path match.
	ResolveSidPath(path []string) (val interface{})
}

//	References a resource by its scoped identifier (Sid).
type RefSid struct {
	s     string
	parts []string
	val   interface{}
}

//	Creates and returns a new RefSid initialized with the specified sidRef.
func NewRefSid(sidRef string) (rs *RefSid) {
	rs = &RefSid{}
	rs.SetSidRef(sidRef)
	return
}

//	Returns the Sid currently referenced by me.
func (me *RefSid) S() string {
	return me.s
}

//	Modifies the Sid currently referenced by me.
func (me *RefSid) SetSidRef(sidRef string) {
	me.s = sidRef
	me.parts = nil
	me.val = nil
}

//	Resolves this Sid reference and returns the value, field, attribute, resource it targets.
//	The returned val is always a pointer: so val may be a *SidFloat but it will never be a SidFloat.
func (me *RefSid) V(rsr func(id string) RefSidResolver) (val interface{}) {
	if me.val == nil {
		if len(me.parts) == 0 {
			me.parts = strings.Split(me.S(), "/")
		}
		resolver := rsr(me.parts[0])
		if (resolver != nil) && (len(me.parts) > 1) {
			me.val = resolver.ResolveSidPath(me.parts[1:])
		} else {
			me.val = resolver
		}
	}
	val = me.val
	return
}
