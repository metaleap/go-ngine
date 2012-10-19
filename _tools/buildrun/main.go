package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

type shaderSrc struct {
	name, src string
}

type shaderSrcSortable []shaderSrc

func (p shaderSrcSortable) Swap (i, j int) { p[i], p[j] = p[j], p[i] }
func (p shaderSrcSortable) Len () int { return len(p) }
func (p shaderSrcSortable) Less (i, j int) bool { return p[i].name < p[j].name }

type shaderSrcSortables struct {
	vert, tessCtl, tessEval, geo, frag, comp shaderSrcSortable
}

func (this shaderSrcSortables) MapAll () map[string]shaderSrcSortable {
	return map[string]shaderSrcSortable { "Vertex": this.vert, "TessCtl": this.tessCtl, "TessEval": this.tessEval, "Geometry": this.geo, "Fragment": this.frag, "Compute": this.comp }
}

func checkForMainPackage (filePath string) bool {
	var err error
	var file *os.File
	var rawBytes []byte
	var tmp string
	if strings.Index(filePath, "buildrun") >= 0 { panic("buildrun tool cannot build itself!") }
	if file, err = os.Open(filePath); err == nil {
		defer file.Close()
		if rawBytes, err = ioutil.ReadAll(file); err == nil {
			if tmp = string(rawBytes); strings.HasPrefix(tmp, "package main\n") || strings.HasPrefix(tmp, "package main\r") {
				return true
			}
		}
	}
	return false
}

func collectShaders (srcDirPath string, allShaders *shaderSrcSortables, iShaders map[string]string, stripComments bool) {
	var err error
	var src *os.File
	var fileInfos []os.FileInfo
	var fileInfo os.FileInfo
	var rawSrc []byte
	var fileName, shaderSource string
	var isIncShader, isVertShader, isTessCtlShader, isTessEvalShader, isGeoShader, isFragShader, isCompShader bool
	var pos1, pos2 int
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
						if src, err = os.Open(filepath.Join(srcDirPath, fileName)); err == nil {
							if rawSrc, err = ioutil.ReadAll(src); err == nil {
								shaderSource = string(rawSrc);
								if stripComments {
									for {
										if pos1 = strings.Index(shaderSource, "/*"); pos1 < 0 { break }
										if pos2 = strings.Index(shaderSource, "*/"); pos2 < pos1 { break }
										shaderSource = shaderSource[0 : pos1] + shaderSource[pos2 + 2 :]
									}
								}
								if isIncShader { iShaders[fileName] = shaderSource }
								if isVertShader { allShaders.vert = append(allShaders.vert, shaderSrc { fileName, shaderSource }) }
								if isTessCtlShader { allShaders.tessCtl = append(allShaders.tessCtl, shaderSrc { fileName, shaderSource }) }
								if isTessEvalShader { allShaders.tessEval = append(allShaders.tessEval, shaderSrc { fileName, shaderSource }) }
								if isGeoShader { allShaders.geo = append(allShaders.geo, shaderSrc { fileName, shaderSource }) }
								if isFragShader { allShaders.frag = append(allShaders.frag, shaderSrc { fileName, shaderSource }) }
								if isCompShader { allShaders.comp = append(allShaders.comp, shaderSrc { fileName, shaderSource }) }
							}
							src.Close()
						}
					}
				}
			}
		}
	}
}

func generateShadersFile (srcDirPath, outFilePath, pkgName string, stripComments bool) {
	var err error
	var src *os.File
	var shaderSource shaderSrc
	var allNames = []string {}
	var rawSrc []byte
	var glslSrc = "package " + pkgName + "\n\nfunc init () {\n\tvar rss = &TShaderSources { map[string]string {}, map[string]string {}, map[string]string {}, map[string]string {}, map[string]string {}, map[string]string {} }\n"
	var glslOldSrc = ""
	var allShaders = shaderSrcSortables { shaderSrcSortable {}, shaderSrcSortable {}, shaderSrcSortable {}, shaderSrcSortable {}, shaderSrcSortable {}, shaderSrcSortable {} }
	var iShaders = map[string]string {}
	var shaderName string
	if src, err = os.Open(outFilePath); err == nil {
		if rawSrc, err = ioutil.ReadAll(src); err == nil { glslOldSrc = string(rawSrc); }
		src.Close()
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
			if shaderName = shaderSource.name[: strings.LastIndex(shaderSource.name, ".")]; !inSlice(allNames, shaderName) {
				allNames = append(allNames, shaderName)
			}
			glslSrc += fmt.Sprintf("\trss.%s[\"%s\"] = %#v\n", varName, shaderName, includeShaders(shaderSource.name, shaderSource.src, iShaders))
		}
	}
	if glslSrc += fmt.Sprintf("\tShaderMan.AllSources = rss\n\tShaderMan.AllNames = %#v\n}\n", allNames); glslSrc != glslOldSrc {
		ioutil.WriteFile(outFilePath, []byte(glslSrc), os.ModePerm)
	}
}

func includeShaders (fileName, shaderSource string, iShaders map[string]string) string {
	var lines = strings.Split(shaderSource, "\n")
	var linePrefix = "#pragma incl "
	var str string
	var i int
	var includes []string
	for i, str = range lines {
		if strings.HasPrefix(str, linePrefix) {
			includes = strings.Split(str[len(linePrefix) :], " ")
			break
		}
	}
	if len(includes) > 0 {
		shaderSource = fmt.Sprintf("#line 1 \"%v\"\n", fileName) + strings.Join(lines[: i], "\n")
		for _, str = range includes {
			shaderSource += fmt.Sprintf("\n#line %v \"%v\"\n", 1, str)
			shaderSource += fmt.Sprintf("%v\n", iShaders[str])
		}
		shaderSource += fmt.Sprintf("#line %v \"%v\"\n", i + 1, fileName)
		shaderSource += strings.Join(lines[i + 1 :], "\n")
		return includeShaders(fileName, shaderSource, iShaders)
	}
	return shaderSource
}

func inSlice (slice []string, val string) bool {
	for _, v := range(slice) { if v == val { return true } }
	return false
}

func main () {
	runtime.LockOSThread()
	var startTime = time.Now()
	var pathSep = string(os.PathSeparator)
	var flagFilePath = flag.String("f", "", "file: current .go source file from which to build")
	var goInstPath string
	var goPath = os.ExpandEnv("$GOPATH")
	var isMainPkg bool
	var origFilePath, cmdRunPath string
	var cmdRawOut []byte
	var err error
	var ngineMatch = "github.com/go3d/go-ngine/"
	var nginePath string
	var allowRun = true
	flag.Parse()
	origFilePath = *flagFilePath
	isMainPkg = checkForMainPackage(origFilePath)
	goInstPath = strings.Replace(origFilePath [len(filepath.Join(goPath, "src") + pathSep) : ], pathSep, "/", -1)
	goInstPath = goInstPath [0 : strings.LastIndex(goInstPath, "/")]
	if (strings.HasPrefix(goInstPath, ngineMatch)) {
		nginePath = origFilePath [ : strings.Index(origFilePath, "github.com") + len(ngineMatch)]
		generateShadersFile(filepath.Join(nginePath, "client", "glcore", "_glsl"), filepath.Join(nginePath, "client", "glcore", "-auto-generated-glsl-src.go"), "glcore", true)
	}
	cmdRawOut, err = exec.Command("go", "install", goInstPath).CombinedOutput()
	if len(cmdRawOut) > 0 {
		allowRun = false
		fmt.Printf("%v\n", trimLines(string(cmdRawOut), 5))
	}
	if err != nil {
		allowRun = false
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("TOTAL BUILD TIME: %v\n", time.Now().Sub(startTime))
	if (allowRun && isMainPkg) {
		cmdRunPath = filepath.Join(goPath, "bin", goInstPath [strings.LastIndex(goInstPath, "/") + 1 : ])
		cmdRawOut, err = exec.Command(cmdRunPath).CombinedOutput()
		if len(cmdRawOut) > 0 { fmt.Printf("%v\n", trimLines(string(cmdRawOut), 10)) }
		if err != nil { fmt.Printf("%+v\n", err) }
	}
}

func trimLines (str string, maxLines int) string {
	return str
	var lines = strings.Split(str, "\n")
	if len(lines) > maxLines { lines = lines[0 : maxLines] }
	return strings.Join(lines, "\n")
}
