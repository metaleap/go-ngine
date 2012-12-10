package main

import (
	"flag"
	"fmt"
	"strings"

	nga "github.com/go3d/go-ngine/assets"
	ugo "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
)

var (
	outFilePath = flag.String("dst", ugo.GopathSrcGithub("go3d", "go-ngine", "assets", "collada", "imp-1.5", "skeleton.txt"), "out file path")
)

func main() {
	const (
		srcImpCtor = `
func obj_%s(xn *xmlx.Node, n string) (obj *nga.%s) {
	if (xn != nil) && (len(n) > 0) {
		xn = xsn(xn, n)
	}
	if xn != nil {
		obj = nga.New%s()
	}
	return
}
`
		srcImpNew = `
func obj_%s(xn *xmlx.Node, n string) (obj *nga.%s) {
	if (xn != nil) && (len(n) > 0) {
		xn = xsn(xn, n)
	}
	if xn != nil {
		obj = new(nga.%s)
	}
	return
}
`
		srcImpN = `
func objs_%s(xn *xmlx.Node, n string) (objs []*nga.%s) {
	xns := xsns(xn, n)
	objs = make([]*nga.%s, len(xns))
	for i, xn := range xns {
		objs[i] = obj_%s(xn, "")
	}
	return
}
`
	)
	flag.Parse()
	src := ""
	ok := false

	for n, _ := range nga.Types {
		if !(strings.HasPrefix(n, "Base") || strings.HasSuffix(n, "Base") || strings.HasPrefix(n, "Has")) {
			if _, ok = nga.Functions["New"+n]; ok && (n != "FxImageInitFrom") {
				src += fmt.Sprintf(srcImpCtor, n, n, n)
			} else {
				src += fmt.Sprintf(srcImpNew, n, n, n)
			}
			src += fmt.Sprintf(srcImpN, n, n, n, n)
		}
	}
	uio.WriteTextFile(*outFilePath, src)
}
