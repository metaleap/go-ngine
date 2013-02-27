package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	ugfx "github.com/metaleap/go-util/gfx"
	uio "github.com/metaleap/go-util/io"
	ustr "github.com/metaleap/go-util/str"
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

func collectShaders(srcDirPath string, allShaders *shaderSrcSortables, incShaders map[string]string, stripComments bool) {
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
					collectShaders(filepath.Join(srcDirPath, fileName), allShaders, incShaders, stripComments)
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
								incShaders[fileName] = shaderSource
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

func generateShadersSource(srcDirPath string, stripComments bool) (err error, newSrc string) {
	var (
		shaderSource       shaderSrc
		shaderName, tmpSrc string
		srcBuf             ustr.Buffer
	)
	srcBuf.Writeln("\tglc.progMan.Init()\n\tglc.uberShader.init()")
	allShaders := shaderSrcSortables{shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}, shaderSrcSortable{}}
	incShaders := map[string]string{}
	collectShaders(srcDirPath, &allShaders, incShaders, stripComments)
	sort.Sort(allShaders.comp)
	sort.Sort(allShaders.frag)
	sort.Sort(allShaders.geo)
	sort.Sort(allShaders.tessCtl)
	sort.Sort(allShaders.tessEval)
	sort.Sort(allShaders.vert)
	allNames := map[string]bool{}
	for varName, shaderSrcItem := range allShaders.mapAll() {
		for _, shaderSource = range shaderSrcItem {
			if shaderName = shaderSource.name[:strings.LastIndex(shaderSource.name, ".")]; !allNames[shaderName] {
				srcBuf.Writeln("\tglc.progMan.AddNew(%#v)", shaderName)
				allNames[shaderName] = true
			}
			srcBuf.Writeln("\tglc.progMan.Get(%#v).Sources.In.%s = %#v", shaderName, varName, includeShaders(shaderSource.name, shaderSource.src, incShaders))
		}
	}
	for shaderName, tmpSrc = range incShaders {
		srcBuf.Writeln("\tglc.uberShader.rawSources[%#v] = %#v", shaderName[:strings.Index(shaderName, ".")], tmpSrc)
	}
	newSrc = srcBuf.String()
	return
}

func includeShaders(fileName, shaderSource string, incShaders map[string]string) string {
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
		shaderSource = fmt.Sprintf("#line 1\n") + strings.Join(lines[:i], "\n")
		for _, str = range includes {
			shaderSource += fmt.Sprintf("\n#line %v\n", 1)
			shaderSource += fmt.Sprintf("%v\n", incShaders[str])
		}
		shaderSource += fmt.Sprintf("#line %v\n", i+1)
		shaderSource += strings.Join(lines[i+1:], "\n")
		return includeShaders(fileName, shaderSource, incShaders)
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

var (
	wait                   sync.WaitGroup
	outFileTime            time.Time
	outFilePath, nginePath string
	newSrc                 struct {
		shaders, embeds string
	}
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var srcTimeGlsl, srcTimeEmbeds time.Time
	force := false
	nginePath = os.Args[1]
	outFilePath = filepath.Join(nginePath, "core", "-gen-embed.go")
	if fileInfo, err := os.Stat(outFilePath); err == nil {
		outFileTime = fileInfo.ModTime()
	} else {
		force = true
	}
	if outFileTime.IsZero() {
		force = true
	}

	srcDirPathEmbeds := filepath.Join(nginePath, "core", "_embed")
	if !force {
		if errs := uio.NewDirWalker(nil, uio.NewWalkerVisitor_IsNewerThan(outFileTime, &srcTimeGlsl)).Walk(srcDirPathEmbeds); len(errs) > 0 {
			panic(errs[0])
		}
	}

	if force || srcTimeGlsl.UnixNano() > outFileTime.UnixNano() || srcTimeEmbeds.UnixNano() > outFileTime.UnixNano() {
		fmt.Printf("Re-generating %s...\n", outFilePath)
		wait.Add(1)
		go makeShaders(srcDirPathEmbeds)
		wait.Add(1)
		go makeEmbeds(srcDirPathEmbeds)
		wait.Wait()
		uio.WriteTextFile(outFilePath, fmt.Sprintf("package core\n\n//\tGenerated by ng-buildrun\nfunc init() {\n%s\n%s\n}", newSrc.shaders, newSrc.embeds))
	}
}

func makeEmbeds(srcDirPath string) {
	defer wait.Done()
	filePath := filepath.Join(srcDirPath, "splash.png")
	var buf ustr.Buffer
	buf.Writeln("\t//\tEmbedded binary from %s", filePath)
	if raw := uio.ReadBinaryFile(filePath, true); len(raw) > 0 {
		if strings.HasSuffix(filePath, ".png") {
			if src, _, err := image.Decode(bytes.NewReader(raw)); err == nil {
				dst, _ := ugfx.CloneImage(src, false)
				ugfx.PreprocessImage(src, dst, true, true, true)
				w := new(bytes.Buffer)
				png.Encode(w, dst)
				raw = w.Bytes()
			} else {
				panic(err)
			}
		}
		if len(raw) > 0 {
			buf.Writeln("\tCore.Libs.Images.SplashScreen.InitFrom.RawData = %#v", raw)
		}
	}
	newSrc.embeds = buf.String()
}

func makeShaders(srcDirPath string) {
	defer wait.Done()
	if err, nsrc := generateShadersSource(srcDirPath, true); err != nil {
		panic(err)
	} else {
		newSrc.shaders = fmt.Sprintf("\t//\tGLSL shader sources from %s\n%s", srcDirPath, nsrc)
	}
}
