package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	ioutil "github.com/go-ngine/go-util/io"
)

var (
	curFilePath = ""
	glNames = map[string][]string {}
)

func checkGoFile (filePath string) {
	var fset = token.NewFileSet()
	var astFile *ast.File
	var err error
	var hasGoglImp = false
	if strings.Index(filePath, "_trash") >= 0 { return }
	if astFile, err = parser.ParseFile(fset, filePath, nil, parser.ImportsOnly); err != nil {
		panic(err)
	}
	for _, s := range astFile.Imports {
		if hasGoglImp = (strings.Index(s.Path.Value, "github.com/chsc/gogl/") >= 0); hasGoglImp { break }
	}
	if hasGoglImp {
		if astFile, err = parser.ParseFile(fset, filePath, nil, 0); err != nil { panic(err) }
		curFilePath = filePath
		ast.Inspect(astFile, inspectNode)
	}
}

func inspectNode (node ast.Node) bool {
	var x, sel *ast.Ident
	var sl []string
	switch selExpr := node.(type) {
	case *ast.SelectorExpr:
		switch xExpr := selExpr.X.(type) {
		case *ast.Ident:
			x = xExpr
		}
		sel = selExpr.Sel
	}
	if (x != nil) && (sel != nil) && (x.Name == "gl") {
		if sl = glNames[sel.Name]; sl == nil { sl = []string {} }
		if !inSlice(sl, curFilePath) { sl, glNames[sel.Name] = append(sl, curFilePath), sl }
		return false
	}
	return true
}

func inSlice (slice []string, val string) bool {
	for _, v := range(slice) { if v == val { return true } }
	return false
}

func main() {
	ioutil.WalkDirectory(filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-ngine"), ".go", checkGoFile, true)
	for glName, filePaths := range glNames {
		println(glName)
		for _, fp := range filePaths {
			fmt.Printf("\t\t%v\n", fp)
		}
	}
}
