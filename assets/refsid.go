package assets

import (
	"strconv"
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
//	
//	Sid path examples:
//		-	foo/bar/doodad
//			either: root is a lib that finds object with Id "foo", which resolves path "bar/doodad"
//			or: root is a non-lib object with Id "foo" and resolves path "bar/doodad"
//		-	./bar/doodad
//			root is a non-lib object with its own arbitrary Id and resolves path "bar/doodad"
//		-	bar
//			gets rewritten to "./bar", then: see above
//		-	foo/bar/doodad.Hollow
//			root resolves foo/bar/doodad, then returns pointer to its Hollow field
//			(if doodad supports named-field access by implementing RefSidFielder)
//		-	foo/bar/doodad(2)
//			root resolves foo/bar/doodad, then returns pointer to a "slot" at index 2
//			(if doodad supports indexed-slot access by implementing RefSidIndexer)
func (me *RefSid) Resolve(root RefSidRoot, force bool) {
	if force || (me.V == nil) {
		me.V = nil
		parts := strings.Split(me.S, "/")
		if len(parts) == 1 {
			parts = append([]string{"."}, parts...)
		}
		if resolver := root.sidResolver(parts[0]); (resolver != nil) && (len(parts) > 1) {
			bag := &refSidBag{indexers: [...]int{-1, -1}}
			last := parts[len(parts)-1]
			if pos := strings.Index(last, "."); pos > 0 {
				parts[len(parts)-1] = last[:pos]
				bag.fielder = last[pos+1:]
			} else if pos = strings.Index(last, "("); pos > 0 {
				parts[len(parts)-1] = last[:pos]
				i := 0
				for _, s := range strings.Split(last[pos+1:], ")(") {
					if iv, err := strconv.Atoi(strings.Trim(s, " )(")); err != nil {
						bag.indexers[i] = iv
						if i++; i > 1 {
							break
						}
					}
				}
			}
			if me.V = resolver.sidResolve(parts[1:], bag); me.V != nil {
				if len(bag.fielder) > 0 {
					if fielder, _ := me.V.(RefSidFielder); fielder != nil {
						me.V = fielder.AccessField(bag.fielder)
					} else {
						me.V = nil
					}
				} else if bag.indexers[0] > -1 {
					if indexer, _ := me.V.(RefSidIndexer); indexer != nil {
						me.V = indexer.AccessIndex(bag.indexers[0], bag.indexers[1])
					} else {
						me.V = nil
					}
				}
			}
		}
	}
}

type refSidBag struct {
	fielder, sid string
	indexers     [2]int
	valRaw       interface{}
	valAsRes     refSidResolver
}

func (me *refSidBag) sidResolve(path []string) interface{} {
	if me.sid == path[0] {
		if len(path) == 1 {
			return me.valRaw
		} else if me.valAsRes != nil {
			return me.valAsRes.sidResolve(path[1:], me)
		}
	}
	return nil
}

//	Implemented by select types that embed HasSid to aid resolving Sid paths
//	with a tailing named-field accessor, as in "some/sid/path.fieldName".
type RefSidFielder interface {
	AccessField(fieldName string) interface{}
}

//	Implemented by select types that embed HasSid to aid resolving Sid paths
//	with a tailing indexed-slot accessor, as in "some/sid/path(2)".
type RefSidIndexer interface {
	AccessIndex(i, j int) interface{}
}

type refSidResolver interface {
	sidResolve(path []string, bag *refSidBag) (val interface{})
}

//	This interface needs to be passed to the RefSid.Resolve() method to resolve a Sid path:
//	Implemented by almost all "LibFooDefs" types, plus all types that embed HasId and directly or
//	indirectly lead to fields of types that embed HasSid -- this includes almost all "FooDef" types.
type RefSidRoot interface {
	sidResolver(part0 string) refSidResolver
}
