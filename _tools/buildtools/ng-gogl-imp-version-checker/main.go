//	http://go-ngine.com/blog/2012/10/19/painless-xml-parsing-in-go-plus-how-to-analyse-your-source-tree-for-package-references
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	util "github.com/metaleap/go-util"
	uio "github.com/metaleap/go-util/io"
)

type glNameInfo struct {
	filePaths []string
	nameKind  string // can be function, enum
	glVersion string
}

var (
	curFilePath = ""
	glNames     = map[string]*glNameInfo{}
	glVersions  = map[string][]string{}
	specDoc     *xmlx.Document
)

func checkGoFile(_ *uio.DirWalker, filePath string, isDir bool) bool {
	if strings.HasSuffix(filePath, ".go") {
		var (
			fset       = token.NewFileSet()
			astFile    *ast.File
			err        error
			hasGoglImp = false
		)
		if astFile, err = parser.ParseFile(fset, filePath, nil, parser.ImportsOnly); err != nil {
			panic(err)
		}
		for _, s := range astFile.Imports {
			if hasGoglImp = (strings.Index(s.Path.Value, "github.com/go3d/go-opengl/api/") >= 0); hasGoglImp {
				break
			}
		}
		if hasGoglImp {
			if astFile, err = parser.ParseFile(fset, filePath, nil, 0); err != nil {
				panic(err)
			}
			curFilePath = filePath
			ast.Inspect(astFile, inspectNode)
		}
	}
	return true
}

func inspectNode(node ast.Node) bool {
	var (
		x, sel *ast.Ident
		gni    *glNameInfo
	)
	switch selExpr := node.(type) {
	case *ast.SelectorExpr:
		switch xExpr := selExpr.X.(type) {
		case *ast.Ident:
			x = xExpr
		}
		sel = selExpr.Sel
	}
	if (x != nil) && (sel != nil) && (x.Name == "gl") {
		if gni = glNames[sel.Name]; gni == nil {
			gni, glNames[sel.Name] = &glNameInfo{[]string{}, "", ""}, gni
		}
		if !inSlice(gni.filePaths, curFilePath) {
			gni.filePaths = append(gni.filePaths, curFilePath)
		}
		return false
	}
	return true
}

func inSlice(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func loadSpecXml() {
	specDoc = xmlx.New()
	if err := specDoc.LoadBytes(uio.ReadBinaryFile(util.GopathSrcGithub("go3d", "go-ngine", "_tools", "buildtools", "ng-gogl-imp-version-checker", "opengl.xml"), true), nil); err != nil {
		panic(err)
	}
}

func main() {
	var (
		ver, kind string
		verList   []string
	)
	loadSpecXml()
	enumNodes, funcNodes := specDoc.SelectNodesRecursive("*", "enum"), specDoc.SelectNodesRecursive("*", "function")
	if errs := uio.NewDirWalker(nil, checkGoFile).Walk(util.GopathSrcGithub("go3d")); len(errs) > 0 {
		panic(errs[0])
	}
	for glName, _ := range glNames {
		kind, ver = "", ""
		for _, enode := range enumNodes {
			if enode.As("*", "name") == glName {
				ver, kind = enode.As("*", "version"), "enum"
			}
		}
		if len(ver) == 0 {
			for _, fnode := range funcNodes {
				if fnode.As("*", "name") == glName {
					ver, kind = fnode.As("*", "version"), "function"
				}
			}
		}
		if len(ver) > 0 {
			glNames[glName].glVersion, glNames[glName].nameKind = ver, kind
			if ver > "3.2" {
				if verList = glVersions[ver]; (verList == nil) || (len(verList) == 0) {
					verList = []string{}
				}
				if !inSlice(verList, glName) {
					verList = append(verList, glName)
				}
				glVersions[ver] = verList
			}
		}
	}
	for ver, verList = range glVersions {
		fmt.Printf("GL v%v used %vx:\n", ver, len(verList))
		for _, glName := range verList {
			if gni := glNames[glName]; gni != nil {
				fmt.Printf("\t%v %v:\n", gni.nameKind, glName)
				for _, filePath := range gni.filePaths {
					fmt.Printf("\t\t%v\n", filePath)
				}
			}
		}
	}
}
