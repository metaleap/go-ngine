package assets

import (
	"strings"
)

//	Returns a RefSidResolver based on the specified arg (typically, an Id).
type GetRefSidResolver func(arg string) RefSidResolver

//	References a previously defined parameter.
type RefParam struct {
	//	A parameter reference technically always refers to a Sid.
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

//	References a resource by its scoped identifier (Sid).
type RefSid struct {
	//	The Sid path currently referenced.
	//	To be set ONLY through the NewRefSid() constructor or SetSidRef() method!
	S string
	//	The resolved value referenced by this Sid path.
	//	This is always a pointer: so V may be a *SidFloat but it will never be a SidFloat.
	//	To be set ONLY through the Resolve() method! Reset to nil by the SetSidRef() method.
	V interface{}
}

//	Creates and returns a new RefSid, its S initialized with the specified sidRef.
func NewRefSid(sidRef string) (rs *RefSid) {
	rs = &RefSid{}
	rs.SetSidRef(sidRef)
	return
}

//	Sets S to sidRef and resets V to nil.
func (me *RefSid) SetSidRef(sidRef string) {
	me.S, me.V = sidRef, nil
}

//	Resolves this Sid reference (if V is nil or force is true), sets and returns V.
//	If no match is found for the full path, V will become nil
//	(rather than, say, a partial-path-match result-value).
func (me *RefSid) Resolve(root RefSidResolverRoot, force bool) interface{} {
	if force || (me.V == nil) {
		parts := strings.Split(me.S, "/")
		if resolver := root.resolver(parts[0]); (resolver != nil) && (len(parts) > 1) {
			me.V = resolver.resolveSidPath(parts[1:])
		} else {
			me.V = resolver
		}
	}
	return me.V
}

//	Resolves a Sid path. Though its resolveSidPath() method is unexported,
//	all "applicable" struct types in this package implement this interface.
type RefSidResolver interface {
	resolveSidPath(path []string) (val interface{})
}

type RefSidResolverRoot interface {
	resolver(part0 string) RefSidResolver
}

func sidResolveCore(path []string, val interface{}, res RefSidResolver, sid string) interface{} {
	if sid == path[0] {
		if len(path) == 1 {
			return val
		} else if res != nil {
			return res.resolveSidPath(path[1:])
		}
	}
	return nil
}
