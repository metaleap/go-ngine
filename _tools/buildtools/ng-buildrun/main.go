package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type shaderSrc struct {
	name, src string
}

type shaderSrcSortable []shaderSrc

func (p shaderSrcSortable) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p shaderSrcSortable) Len() int           { return len(p) }
func (p shaderSrcSortable) Less(i, j int) bool { return p[i].name < p[j].name }

type shaderSrcSortables struct {
	vert, tessCtl, tessEval, geo, frag, comp shaderSrcSortable
}

func (this shaderSrcSortables) MapAll() map[string]shaderSrcSortable {
	return map[string]shaderSrcSortable{"Vertex": this.vert, "TessCtl": this.tessCtl, "TessEval": this.tessEval, "Geometry": this.geo, "Fragment": this.frag, "Compute": this.comp}
}

func collectShaders(srcDirPath string, allShaders *shaderSrcSortables, iShaders map[string]string, stripComments bool) {
	var (
		err                                                                                                   error
		src                                                                                                   *os.File
		fileInfos                                                                                             []os.FileInfo
		fileInfo                                                                                              os.FileInfo
		rawSrc                                                                                                []byte
		fileName, shaderSource                                                                                string
		isIncShader, isVertShader, isTessCtlShader, isTessEvalShader, isGeoShader, isFragShader, isCompShader bool
		pos1, pos2                                                                                            int
	)
	if src, err = os.Open(srcDirPath); err == nil {
		fileInfos, err = src.Readdir(0)
		src.Close()
		if err == nil {
			for _, fileInfo = range fileInfos {
				fileName = fileInfo.Name()
				if fileInfo.IsDir() {
					collectShaders(filepath.Join(srcDirPath, fileName), allShaders, iShaders, stripComments)
				} else {
					isIncShader, isVertShader, isTessCtlShader, isTessEvalShader, isGeoShader, isFragShader, isCompShader = strings.HasSuffix(fileName, ".glsl"), strings.HasSuffix(fileName, ".glvs"), strings.HasSuffix(fileName, ".gltc"), strings.HasSuffix(fileName, ".glte"), strings.HasSuffix(fileName, ".glgs"), strings.HasSuffix(fileName, ".glfs"), strings.HasSuffix(fileName, ".glcs")
					if isIncShader || isVertShader || isTessCtlShader || isTessEvalShader || isGeoShader || isFragShader || isCompShader {
						if rawSrc, err = ioutil.ReadFile(filepath.Join(srcDirPath, fileName)); err == nil {
							shaderSource = string(rawSrc)
							if stripComments {
								for {
									if pos1, pos2 = strings.Index(shaderSource, "/*"), strings.Index(shaderSource, "*/"); (pos1 < 0) || (pos2 < pos1) {
										break
									}
									shaderSource = shaderSource[0:pos1] + shaderSource[pos2+2:]
								}
							}
							if isIncShader {
								iShaders[fileName] = shaderSource
							}
							if isVertShader {
								allShaders.vert = append(allShaders.vert, shaderSrc{fileName, shaderSource})
							}
							if isTessCtlShader {
								allShaders.tessCtl = append(allShaders.tessCtl, shaderSrc{fileName, shaderSource})
							}
							if isTessEvalShader {
								allShaders.tessEval = append(allShaders.tessEval, shaderSrc{fileName, shaderSource})
							}
							if isGeoShader {
								allShaders.geo = append(allShaders.geo, shaderSrc{fileName, shaderSource})
							}
							if isFragShader {
								allShaders.frag = append(allShaders.frag, shaderSrc{fileName, shaderSource})
							}
							if isCompShader {
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
		err                    error
		shaderSource           shaderSrc
		allNames               = []string{}
		rawSrc                 []byte
		glslSrc                = "package " + pkgName + "\n\nfunc init() {\n\tvar rss = newGlShaderSources()\n"
		allShaders             = shaderSrcSortables{shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}}
		iShaders               = map[string]string{}
		glslOldSrc, shaderName string
	)
	if rawSrc, err = ioutil.ReadFile(outFilePath); err == nil {
		glslOldSrc = string(rawSrc)
	}
	collectShaders(srcDirPath, &allShaders, iShaders, stripComments)
	sort.Sort(allShaders.comp)
	sort.Sort(allShaders.frag)
	sort.Sort(allShaders.geo)
	sort.Sort(allShaders.tessCtl)
	sort.Sort(allShaders.tessEval)
	sort.Sort(allShaders.vert)
	for varName, shaderSrcItem := range allShaders.MapAll() {
		for _, shaderSource = range shaderSrcItem {
			if shaderName = shaderSource.name[:strings.LastIndex(shaderSource.name, ".")]; !inSlice(allNames, shaderName) {
				allNames = append(allNames, shaderName)
			}
			glslSrc += fmt.Sprintf("\trss.%s[\"%s\"] = %#v\n", varName, shaderName, includeShaders(shaderSource.name, shaderSource.src, iShaders))
		}
	}
	if glslSrc += fmt.Sprintf("\tglShaderMan.AllSources = rss\n\tglShaderMan.AllNames = %#v\n}\n", allNames); glslSrc != glslOldSrc {
		ioutil.WriteFile(outFilePath, []byte(glslSrc), os.ModePerm)
	}
	return true
}

func includeShaders(fileName, shaderSource string, iShaders map[string]string) string {
	var (
		lines      = strings.Split(shaderSource, "\n")
		linePrefix = "#pragma incl "
		str        string
		i          int
		includes   []string
	)
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
	var (
		nginePath   = os.Args[1]
		srcDirPath  = filepath.Join(nginePath, "core", "_glsl")
		outFilePath = filepath.Join(nginePath, "core", "-auto-generated-glsl-src.go")
	)
	runtime.LockOSThread()
	fmt.Printf("Merging shader files inside %v into %v... ", strings.Replace(srcDirPath, nginePath, ".", -1), strings.Replace(outFilePath, nginePath, ".", -1))
	generateShadersFile(srcDirPath, outFilePath, "core", true)
	fmt.Println("DONE.")
}
