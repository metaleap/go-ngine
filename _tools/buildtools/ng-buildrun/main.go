package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	uio "github.com/metaleap/go-util/io"
)

type shaderSrc struct {
	name, src string
}

type shaderSrcSortable []shaderSrc

func (me shaderSrcSortable) Swap(i, j int)      { me[i], me[j] = me[j], me[i] }
func (me shaderSrcSortable) Len() int           { return len(me) }
func (me shaderSrcSortable) Less(i, j int) bool { return me[i].name < me[j].name }

type shaderSrcSortables struct {
	vert, tessCtl, tessEval, geo, frag, comp shaderSrcSortable
}

func (me shaderSrcSortables) mapAll() map[string]shaderSrcSortable {
	return map[string]shaderSrcSortable{"Vertex": me.vert, "TessCtl": me.tessCtl, "TessEval": me.tessEval, "Geometry": me.geo, "Fragment": me.frag, "Compute": me.comp}
}

func collectShaders(srcDirPath string, allShaders *shaderSrcSortables, iShaders map[string]string, stripComments bool) {
	var (
		fileInfos                                                   []os.FileInfo
		fileName, shaderSource                                      string
		isInc, isVert, isTessCtl, isTessEval, isGeo, isFrag, isComp bool
		pos1, pos2                                                  int
	)
	if src, err := os.Open(srcDirPath); err == nil {
		fileInfos, err = src.Readdir(0)
		src.Close()
		if err == nil {
			for _, fileInfo := range fileInfos {
				fileName = fileInfo.Name()
				if fileInfo.IsDir() {
					collectShaders(filepath.Join(srcDirPath, fileName), allShaders, iShaders, stripComments)
				} else {
					isInc, isVert, isTessCtl, isTessEval, isGeo, isFrag, isComp = strings.HasSuffix(fileName, ".glsl"), strings.HasSuffix(fileName, ".glvs"), strings.HasSuffix(fileName, ".gltc"), strings.HasSuffix(fileName, ".glte"), strings.HasSuffix(fileName, ".glgs"), strings.HasSuffix(fileName, ".glfs"), strings.HasSuffix(fileName, ".glcs")
					if isInc || isVert || isTessCtl || isTessEval || isGeo || isFrag || isComp {
						if shaderSource = uio.ReadTextFile(filepath.Join(srcDirPath, fileName), false, ""); len(shaderSource) > 0 {
							if stripComments {
								for {
									if pos1, pos2 = strings.Index(shaderSource, "/*"), strings.Index(shaderSource, "*/"); (pos1 < 0) || (pos2 < pos1) {
										break
									}
									shaderSource = shaderSource[0:pos1] + shaderSource[pos2+2:]
								}
							}
							if isInc {
								iShaders[fileName] = shaderSource
							}
							if isVert {
								allShaders.vert = append(allShaders.vert, shaderSrc{fileName, shaderSource})
							}
							if isTessCtl {
								allShaders.tessCtl = append(allShaders.tessCtl, shaderSrc{fileName, shaderSource})
							}
							if isTessEval {
								allShaders.tessEval = append(allShaders.tessEval, shaderSrc{fileName, shaderSource})
							}
							if isGeo {
								allShaders.geo = append(allShaders.geo, shaderSrc{fileName, shaderSource})
							}
							if isFrag {
								allShaders.frag = append(allShaders.frag, shaderSrc{fileName, shaderSource})
							}
							if isComp {
								allShaders.comp = append(allShaders.comp, shaderSrc{fileName, shaderSource})
							}
						}
					}
				}
			}
		}
	}
}

func generateShadersFile(srcDirPath, outFilePath, pkgName string, stripComments bool) bool {
	var (
		shaderSource           shaderSrc
		allNames               []string
		glslOldSrc, shaderName string
	)
	glslSrc := "package " + pkgName + "\n\nfunc init() {\n\tvar rss = newGlShaderSources()\n"
	allShaders := shaderSrcSortables{shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}}
	iShaders := map[string]string{}
	glslOldSrc = uio.ReadTextFile(outFilePath, false, "")
	collectShaders(srcDirPath, &allShaders, iShaders, stripComments)
	sort.Sort(allShaders.comp)
	sort.Sort(allShaders.frag)
	sort.Sort(allShaders.geo)
	sort.Sort(allShaders.tessCtl)
	sort.Sort(allShaders.tessEval)
	sort.Sort(allShaders.vert)
	for varName, shaderSrcItem := range allShaders.mapAll() {
		for _, shaderSource = range shaderSrcItem {
			if shaderName = shaderSource.name[:strings.LastIndex(shaderSource.name, ".")]; !inSlice(allNames, shaderName) {
				allNames = append(allNames, shaderName)
			}
			glslSrc += fmt.Sprintf("\trss.%s[\"%s\"] = %#v\n", varName, shaderName, includeShaders(shaderSource.name, shaderSource.src, iShaders))
		}
	}
	if glslSrc += fmt.Sprintf("\tglShaderMan.AllSources = rss\n\tglShaderMan.AllNames = %#v\n}\n", allNames); glslSrc != glslOldSrc {
		uio.WriteTextFile(outFilePath, glslSrc)
	}
	return true
}

func includeShaders(fileName, shaderSource string, iShaders map[string]string) string {
	const linePrefix = "#pragma incl "
	var (
		str      string
		i        int
		includes []string
	)
	lines := strings.Split(shaderSource, "\n")
	for i, str = range lines {
		if strings.HasPrefix(str, linePrefix) {
			includes = strings.Split(str[len(linePrefix):], " ")
			break
		}
	}
	if len(includes) > 0 {
		shaderSource = fmt.Sprintf("#line 1 \"%v\"\n", fileName) + strings.Join(lines[:i], "\n")
		for _, str = range includes {
			shaderSource += fmt.Sprintf("\n#line %v \"%v\"\n", 1, str)
			shaderSource += fmt.Sprintf("%v\n", iShaders[str])
		}
		shaderSource += fmt.Sprintf("#line %v \"%v\"\n", i+1, fileName)
		shaderSource += strings.Join(lines[i+1:], "\n")
		return includeShaders(fileName, shaderSource, iShaders)
	}
	return shaderSource
}

func inSlice(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func main() {
	var outTime, srcTime, tmpTime int64
	runtime.LockOSThread()
	nginePath := os.Args[1]
	srcDirPath, outFilePath := filepath.Join(nginePath, "core", "_glsl"), filepath.Join(nginePath, "core", "-auto-generated-glsl-src.go")
	if fileInfo, err := os.Stat(outFilePath); err == nil {
		outTime = fileInfo.ModTime().UnixNano()
		ff := func(filePath string, rec bool) bool {
			if fileInfo, err = os.Stat(filePath); (err == nil) && !fileInfo.IsDir() {
				if tmpTime = fileInfo.ModTime().UnixNano(); tmpTime > srcTime {
					srcTime = tmpTime
				}
			}
			return srcTime <= outTime
		}
		uio.WalkDirectory(srcDirPath, "", ff, true)
	}
	if srcTime > outTime {
		fmt.Printf("Re-merging changed shader files inside %v into %v... ", strings.Replace(srcDirPath, nginePath, ".", -1), strings.Replace(outFilePath, nginePath, ".", -1))
		generateShadersFile(srcDirPath, outFilePath, "core", true)
		fmt.Println("DONE.")
	}
}
