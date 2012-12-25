package assets

import (
	"strings"
)

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

//	If me.V is nil or force is true: resolves the Sid path in me.S and sets V to the result.
//	For possible root arguments, see RefSidRoot. If no match is found for the full path, V
//	will become nil (rather than, say, a partial-path-match result-value).
func (me *RefSid) Resolve(root RefSidRoot, force bool) {
	if force || (me.V == nil) {
		parts := strings.Split(me.S, "/")
		if resolver := root.resolver(parts[0]); (resolver != nil) && (len(parts) > 1) {
			bag := &refSidBag{}
			last := parts[len(parts)-1]
			if pos := strings.Index(last, "."); pos > 0 {
				bag.fielder = last[pos+1:]
				parts[len(parts)-1] = last[:pos]
			} else if pos = strings.Index(last, "("); pos > 0 {
				bag.indexer = last[pos+1:]
				parts[len(parts)-1] = last[:pos]
			}
			me.V = resolver.resolveSidPath(parts[1:], bag)
		} else {
			me.V = resolver
		}
	}
}

type refSidBag struct {
	fielder, indexer, sid string
	valRaw                interface{}
	valAsRes              refSidResolver
}

type refSidFielder interface {
	accessField(field string) interface{}
}

type refSidIndexer interface {
	accessIndex(indices string) interface{}
}

type refSidResolver interface {
	resolveSidPath(path []string, bag *refSidBag) (val interface{})
}

//	This interface needs to be passed to the RefSid.Resolve() method to resolve a Sid path:
//	Implemented by all "LibFooDefs" types, plus all types that embed HasId and directly or
//	indirectly lead to types that embed HasSid -- this includes almost all "FooDef" types.
//	The latter ignore the Id part of the Sid path (but it is still required syntactically),
//	whereas the former (all "LibFooDefs" types) require it.
type RefSidRoot interface {
	resolver(part0 string) refSidResolver
}

func sidResolveCore(path []string, bag *refSidBag) interface{} {
	if bag.sid == path[0] {
		if len(path) == 1 {
			return bag.valRaw
		} else if bag.valAsRes != nil {
			return bag.valAsRes.resolveSidPath(path[1:], bag)
		}
	}
	return nil
}
