package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path"
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

func checkForFileCopy (filePath string) {
/*
	var fileCopyMatch = "//CP "
	var err error
	var file *os.File
	var rawBytes []byte
	var tmp string
	if file, err = os.Open(filePath); err == nil {
		if rawBytes, err = ioutil.ReadAll(file); err == nil {
			if tmp = string(rawBytes); strings.HasPrefix(tmp, fileCopyMatch) {
				tmp = tmp[len(fileCopyMatch) : ]
				tmp = tmp[ : strings.IndexAny(tmp, "\r\n")]
				ioutil.WriteFile(tmp, rawBytes, os.ModePerm)
			}
		}
	}
*/
}

func checkForMainPackage (filePath string) bool {
	var err error
	var file *os.File
	var rawBytes []byte
	var tmp string
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
/*
	for _, shaderSource = range fShaders {
		shaderSource.src = includeShaders(shaderSource.name, shaderSource.src, iShaders)
		// if shaderSource.name == "cast.fs" { ioutil.WriteFile(outFilePath + ".fs", []byte(shaderSource.src), os.ModePerm) }
		glslSrc += fmt.Sprintf("\tFShaders[\"%s\"] = %#v\n", shaderSource.name[: strings.LastIndex(shaderSource.name, ".")], shaderSource.src)
	}
*/
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
	var ngineMatch = "github.com/go-ngine/go-ngine/"
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

func old_main () {
	var startTime = time.Now()
	var flagDebug = flag.Bool("d", false, "debug: true to add, false to strip debug symbols")
	var flagRun = flag.Bool("r", false, "run the built binary?")
	var flagFilePath = flag.String("f", "", "file: current .go source file from which to build")
	var pathSep = string(os.PathSeparator)
	var subPathMatch, pkgPrefix, pkgOutPrefix = pathSep + "src" + pathSep + "terra" + pathSep, "terra" + pathSep, "terra-"
	var origFilePath, gdDirPath, srcFilePath, tmpFileName, buildPkg, buildOut string
	var dir *os.File
	var dirPath, err = os.Getwd()
	var cmd *exec.Cmd
	var cmdArgs, dirFileNames, pkgNameParts []string
	var cmdRawOut []byte
//	var pos int
//	var isMainPkg bool
	if err != nil { panic(err) }
	flag.Parse()
	origFilePath = *flagFilePath
//Restart:
	if (len(origFilePath) > 0) && strings.Contains(origFilePath, subPathMatch) {
		gdDirPath = origFilePath[0 : strings.Index(origFilePath, subPathMatch)]
		if strings.HasSuffix(origFilePath, ".exe.go") { srcFilePath = origFilePath }
		for (len(srcFilePath) == 0) && (strings.Index(dirPath, pkgPrefix) > 0) {
			if dir, err = os.Open(dirPath); err != nil { panic(err) }
			if dirFileNames, err = dir.Readdirnames(0); err != nil{ dir.Close(); panic(err) }
			for _, tmpFileName = range dirFileNames {
				if strings.HasSuffix(tmpFileName, ".exe.go") {
					srcFilePath = path.Join(dirPath, tmpFileName)
					break
				}
			}
			dirPath = path.Dir(dirPath)
			dir.Close()
		}
		if len(srcFilePath) > 0 {
			checkForFileCopy(srcFilePath)
			buildPkg = pkgPrefix + path.Dir(srcFilePath[strings.Index(srcFilePath, subPathMatch) + len(subPathMatch) :])
			pkgNameParts = strings.Split(buildPkg[len(pkgPrefix) :], pathSep)
			buildOut = path.Join(gdDirPath, pkgOutPrefix + pkgNameParts[0], path.Base(srcFilePath[0 : len(srcFilePath) - len(".go")]))
			if strings.Contains(buildOut, pkgOutPrefix + "client" + pathSep) { generateShadersFile(path.Join(path.Dir(srcFilePath), "_glsl"), path.Join(path.Dir(srcFilePath), "engine", "shaders", "-auto-generated-glsl-src.go"), "shaders", !*flagDebug) }
			cmdArgs = []string { "install"/*, "-o", buildOut*/ }
			if !*flagDebug { cmdArgs = append(cmdArgs, "-gcflags=-B", "-ldflags=-s") }
			cmdArgs = append(cmdArgs, buildPkg)
			fmt.Printf("go %v :\n", cmdArgs)
			cmdRawOut, err = exec.Command("go", cmdArgs...).CombinedOutput()
			if len(cmdRawOut) > 0 { fmt.Printf("%v\n", trimLines(string(cmdRawOut), 5)) }
			if err != nil { fmt.Printf("%+v\n", err); os.Exit(1) }
			cmdRawOut, err = exec.Command("cp", "-f", path.Join(os.ExpandEnv("$GOBIN"), pkgNameParts[len(pkgNameParts) - 1]), buildOut).CombinedOutput()
			if len(cmdRawOut) > 0 { fmt.Printf("%v\n", trimLines(string(cmdRawOut), 5)) }
			if err != nil { fmt.Printf("%+v\n", err); os.Exit(1) }
			fmt.Printf("Built %v in %v\n", buildOut, time.Now().Sub(startTime))
			if *flagRun {
				cmd = exec.Command(buildOut)
				cmd.Dir = path.Dir(buildOut)
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				err = cmd.Start()
				if err != nil { fmt.Printf("\nRUN ERROR:\n%+v\n", err); os.Exit(1) }
			}
		} else {
			fmt.Printf("No src file found building from %v, doing a module compile only.\n", origFilePath)
			cmdArgs = []string { "-o", path.Join(gdDirPath, "_tmp", "6g.o"), origFilePath }
			cmdRawOut, err = exec.Command("6g", cmdArgs...).CombinedOutput()
			if len(cmdRawOut) > 0 { fmt.Printf("%v\n", trimLines(string(cmdRawOut), 5)) }
			if err != nil { fmt.Printf("%+v\n", err); os.Exit(1) }
		}
	} else {
		checkForFileCopy(origFilePath)
		subPathMatch, pkgPrefix = pathSep + "src" + pathSep + "github.com" + pathSep + "go-ngine" + pathSep, "github.com" + pathSep + "go-ngine" + pathSep
		if strings.Contains(origFilePath, subPathMatch) {
			gdDirPath = origFilePath[ : strings.Index(origFilePath, subPathMatch)]
			buildPkg = origFilePath[strings.Index(origFilePath, subPathMatch) + len(subPathMatch) : ]
//			for pos = strings.LastIndex(buildPkg, pathSep); pos >= 0; pos = strings.LastIndex(buildPkg, pathSep) {
//				buildPkg = buildPkg[ : pos]
//				if strings.HasPrefix(buildPkg, "go-ngine" + pathSep) || strings.HasPrefix(buildPkg, "go-glutil" + pathSep) { continue }
				if buildPkg == "go-ngine" { generateShadersFile(path.Join(gdDirPath, subPathMatch, "go-ngine", "_glsl"), path.Join(gdDirPath, subPathMatch, "go-ngine", "-auto-generated-glsl-src.go"), "ngine", !*flagDebug) }
				cmdArgs = []string { "install"/*, "-o", buildOut*/ }
				if !*flagDebug { cmdArgs = append(cmdArgs, "-gcflags=-B", "-ldflags=-s") }
				cmdArgs = append(cmdArgs, subPathMatch + buildPkg)
				fmt.Printf("go %v :\n", cmdArgs)
				cmdRawOut, err = exec.Command("go", cmdArgs...).CombinedOutput()
				if len(cmdRawOut) > 0 { fmt.Printf("%v\n", trimLines(string(cmdRawOut), 5)) }
				if err != nil { fmt.Printf("%+v\n", err) }
//				if (isMainPkg) { fmt.Println("DOBREAK"); break }
//			}
/*
			if buildPkg == "go-ngine" {
				subPathMatch, pkgPrefix = pathSep + "src" + pathSep + "terra" + pathSep, "terra" + pathSep
				origFilePath = path.Join(gdDirPath, subPathMatch, "_glrttest", "rttest.exe.go")
				goto Restart
			}
*/
		} else {
			fmt.Printf("Not in known package path: %v (does not match %v)\n", origFilePath, subPathMatch)
		}
	}
}

func trimLines (str string, maxLines int) string {
	return str
	var lines = strings.Split(str, "\n")
	if len(lines) > maxLines { lines = lines[0 : maxLines] }
	return strings.Join(lines, "\n")
}
