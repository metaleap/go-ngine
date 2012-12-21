package main

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	nga "github.com/go3d/go-ngine/assets"
	uio "github.com/metaleap/go-util/io"
)

type typeMap map[reflect.Type]bool

var (
	allStructs []reflect.Type
	typeDeps   = map[reflect.Type]typeMap{}
)

func findTypeDeps(hasFieldOfType reflect.Type) {
	if typeDeps[hasFieldOfType] == nil {
		var (
			i      int
			ft, rt reflect.Type
			types  = typeMap{}
			pt     = reflect.PtrTo(hasFieldOfType)
		)
		for _, rt = range allStructs {
			for i = 0; i < rt.NumField(); i++ {
				if ft = rt.Field(i).Type; (ft == hasFieldOfType) || (ft == pt) {
					types[rt] = true
					break
				} else {
				}
			}
		}
		typeDeps[hasFieldOfType] = types
		for rt, _ = range types {
			findTypeDeps(rt)
		}
	}
}

func main() {
	var (
		hasSidType, reflType reflect.Type
	)
	runtime.LockOSThread()
	outFilePath := filepath.Join(os.Args[1], "-gen-refsids.go")
	outSrc := "package assets\n\n"
	for _, reflType = range nga.Types {
		if reflType.Kind() == reflect.Struct {
			allStructs = append(allStructs, reflType)
			if reflType.Name() == "HasSid" {
				hasSidType = reflType
			}
		}
	}
	findTypeDeps(hasSidType)
	println(len(typeDeps))
	uio.WriteTextFile(outFilePath, outSrc)
}
