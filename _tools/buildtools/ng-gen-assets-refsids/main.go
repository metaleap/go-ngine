package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	ngr "github.com/go3d/go-ngine/assets/pkgreflect"
	uio "github.com/metaleap/go-util/io"
	ustr "github.com/metaleap/go-util/str"
)

type typeMap map[reflect.Type]bool

var (
	hasSidType   reflect.Type
	allStructs   = map[string]reflect.Type{}
	haveSids     = map[reflect.Type]bool{}
	isResolver   = map[reflect.Type]bool{}
	typesWritten = map[reflect.Type]bool{}
	typeDeps     = map[reflect.Type]typeMap{}
)

func anyOf(t reflect.Type, kinds ...reflect.Kind) bool {
	for _, k := range kinds {
		if t.Kind() == k {
			return true
		}
	}
	return false
}

func elemType(rt reflect.Type) reflect.Type {
	for anyOf(rt, reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice) {
		rt = rt.Elem()
	}
	return rt
}

func isInlineStructField(t reflect.Type) bool {
	return (t.Kind() == reflect.Struct) && (len(t.Name()) == 0)
}

func mapTypeDeps(forType reflect.Type) {
	var (
		dt          reflect.Type
		depTypes    typeMap
		checkFields func(d, c reflect.Type)
	)
	isHasSid := (forType.Name() == "HasSid")
	//	a separate func also lets us walk into embeds and inline structs within the context of depType
	checkFields = func(depType, checkType reflect.Type) {
		var (
			hasSid bool
			ft     reflect.Type
			sf     reflect.StructField
		)
		if _, hasSid = checkType.FieldByName("HasSid"); hasSid {
			haveSids[checkType] = true
		}
		for i := 0; i < checkType.NumField(); i++ {
			if sf = checkType.Field(i); sf.Name != "Def" {
				if sf.Anonymous {
					if isHasSid && (sf.Name == "HasSid") {
						depTypes[depType] = true
					} else {
						checkFields(depType, sf.Type)
					}
				} else if ft = elemType(sf.Type); ft == forType {
					depTypes[depType] = true
				} else if isInlineStructField(ft) {
					checkFields(depType, ft)
				}
			}
		}
	}
	//	type may be encountered many times, only check it once
	if typeDeps[forType] == nil {
		//	check for all known struct types if one of them depends on forType
		depTypes = typeMap{}
		for _, dt = range allStructs {
			checkFields(dt, dt)
		}
		typeDeps[forType] = depTypes
		//	now collect sub-dependencies recursively
		for dt, _ = range depTypes {
			mapTypeDeps(dt)
		}
	}
}

func testResolvers(rt reflect.Type) bool {
	if _, ok := isResolver[rt]; !ok {
		var (
			sf          reflect.StructField
			ft          reflect.Type
			checkFields func(reflect.Type, bool, bool)
		)
		isResolver[rt] = false
		checkFields = func(ct reflect.Type, dbg bool, d2 bool) {
			if (ct != hasSidType) && (ct.Kind() == reflect.Struct) {
				for i := 0; i < ct.NumField(); i++ {
					if sf = ct.Field(i); sf.Name != "Def" {
						if ft = elemType(sf.Type); sf.Anonymous || isInlineStructField(ft) {
							checkFields(ft, false, false)
						} else if testResolvers(ft) || haveSids[ft] {
							isResolver[rt] = true
						}
					}
				}
			}
		}
		checkFields(rt, false, false)
		for dt, _ := range typeDeps[rt] {
			testResolvers(dt)
		}
	}
	return isResolver[rt]
}

func sfmt(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func writeMethod(rt reflect.Type, force bool) (outSrc string) {
	var (
		count      int
		walkFields func(reflect.Type, string)
	)
	walkFields = func(tt reflect.Type, pref string) {
		var (
			isPtr        bool
			amper, lnpre string
			et, ft       reflect.Type
			sf           reflect.StructField
			beginIfNil   func()
			endIfNil     func()
		)
		beginIfNil = func() {
			lnpre = "\t"
			if isPtr = ft.Kind() == reflect.Ptr; isPtr {
				outSrc += sfmt("\tif %s != nil {\n", pref+sf.Name)
				amper = ""
				lnpre += "\t"
			} else {
				amper = "&"
			}
		}
		endIfNil = func() {
			lnpre = "\t"
			if isPtr {
				outSrc += "\t}\n"
			}
		}
		for i := 0; i < tt.NumField(); i++ {
			if sf = tt.Field(i); sf.Name != "Def" {
				ft = sf.Type
				et = elemType(ft)
				if sf.Anonymous {
					walkFields(et, pref)
				} else {
					if haveSids[et] && !sf.Anonymous {
						switch ft.Kind() {
						case reflect.Array, reflect.Map, reflect.Slice:
							outSrc += sfmt("\tfor _, sidItem := range %s {\n\t\tbag.valRaw, bag.valAsRes, bag.sid = sidItem, %s, sidItem.Sid\n\t\tif val = sidResolveCore(path, bag); val != nil {\n\t\t\treturn\n\t\t}\n\t}\n", pref+sf.Name, ustr.Ifs(isResolver[et], "sidItem", "nil"))
						default:
							beginIfNil()
							outSrc += sfmt(lnpre+"bag.valRaw, bag.valAsRes, bag.sid = %s, %s, %s.Sid\n"+lnpre+"if val = sidResolveCore(path, bag); val != nil {\n"+lnpre+"\treturn\n"+lnpre+"}\n", amper+pref+sf.Name, ustr.Ifs(isResolver[et], amper+pref+sf.Name, "nil"), pref+sf.Name)
							endIfNil()
						}
						count++
					} else if isResolver[et] {
						switch ft.Kind() {
						case reflect.Array, reflect.Map, reflect.Slice:
							outSrc += sfmt("\tfor _, subItem := range %s {\n\t\tif val = subItem.resolveSidPath(path, bag); val != nil {\n\t\t\treturn\n\t\t}\n\t}\n", pref+sf.Name)
						default:
							beginIfNil()
							outSrc += sfmt(lnpre+"if val = %s.resolveSidPath(path, bag); val != nil {\n"+lnpre+"\treturn\n"+lnpre+"}\n", pref+sf.Name)
							endIfNil()
						}
						count++
					} else if isInlineStructField(et) {
						walkFields(et, pref+sf.Name+".")
					}
				}
			}
		}
	}
	outSrc += sfmt("func (me *%s) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {\n", rt.Name())
	walkFields(rt, "me.")
	if isResolver[rt] || (count > 0) || force {
		typesWritten[rt] = true
		outSrc += "\treturn\n}\n\n"
		if _, ok := rt.FieldByName("Id"); ok && (isResolver[rt] || (count > 0)) {
			outSrc += sfmt("func (me *%s) resolver(id string) (rsr refSidResolver) {\n\tif (id == me.Id) || (id == \".\") {\n\t\trsr = me\n\t}\n\treturn\n}\n\n", rt.Name())
		}
	} else {
		outSrc = ""
	}
	return
}

func main() {
	var (
		ok bool
		tn string
		rt reflect.Type
	)
	runtime.LockOSThread()
	outFilePath := filepath.Join(os.Args[1], "-gen-refsids.go")
	outSrc := "package assets\n\n"
	//	collect all struct types and catch the HasSid struct type
	for _, rt = range ngr.Types {
		if rt.Kind() == reflect.Struct {
			if !ustr.HasAnyPrefix(rt.Name(), "Lib", "Base", "Has", "Ref") {
				allStructs[rt.Name()] = rt
			}
			if rt.Name() == "HasSid" {
				hasSidType = rt
			}
		}
	}
	//	map all direct and indirect type dependencies on HasSid type
	mapTypeDeps(hasSidType)
	for rt, _ = range typeDeps {
		testResolvers(rt)
	}
	for rt, _ = range typeDeps {
		if rt != hasSidType {
			outSrc += writeMethod(rt, false)
		}
	}
	for tn, rt = range allStructs {
		if _, ok = ngr.Types["Libs"+tn]; ustr.HasAnySuffix(tn, "Def") && (!typesWritten[rt]) && ok {
			outSrc += writeMethod(rt, true)
		}
	}
	uio.WriteTextFile(outFilePath, outSrc[:len(outSrc)-1]+"\n")
}
