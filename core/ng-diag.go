package core

import (
	"log"

	ugl "github.com/go3d/go-opengl/util"
	ugo "github.com/metaleap/go-util"
)

type EngineDiagLogCategory int

const (
	DiagLogCatMisc    EngineDiagLogCategory = 1
	DiagLogCatMeshes  EngineDiagLogCategory = 2
	DiagLogCatShaders EngineDiagLogCategory = 4
	DiagLogCatAll     EngineDiagLogCategory = DiagLogCatMeshes | DiagLogCatMisc | DiagLogCatShaders
)

var (
	Diag = EngineDiag{
		LogCategories: DiagLogCatAll,
		LogCategoryNames: map[EngineDiagLogCategory]string{
			DiagLogCatMisc:    "[INFO]\t\t",
			DiagLogCatShaders: "[SHADER]\t",
			DiagLogCatMeshes:  "[MESHES]\t",
		},
		LogGLErrorsInLoopOnSec: true,
	}
)

func init() {
	ugo.LogErrorFormat = "[ERROR!]\t%v"
}

//	EngineDiag is a singleton type, only used for the core.Diag package-global exported variable.
//	It is only aware of that instance and does not support any other EngineDiag instances.
type EngineDiag struct {
	LogCategories          EngineDiagLogCategory
	LogCategoryNames       map[EngineDiagLogCategory]string
	LogGLErrorsInLoopOnSec bool
	WriteTmpFilesTo        struct {
		BaseDirName    string
		ShaderPrograms string
	}
}

func (_ *EngineDiag) Log(cat EngineDiagLogCategory, fmt string, fmtArgs ...interface{}) {
	if (Diag.LogCategories & cat) == cat {
		log.Printf(Diag.LogCategoryNames[cat]+fmt, fmtArgs...)
	}
}

func (_ *EngineDiag) LogErr(err error) {
	ugo.LogError(err)
}

func (_ *EngineDiag) LogIfGlErr(fmt string, fmtArgs ...interface{}) {
	ugl.LogLastError(fmt, fmtArgs...)
}

func (_ *EngineDiag) LogMeshes(fmt string, fmtArgs ...interface{}) {
	Diag.Log(DiagLogCatMeshes, fmt, fmtArgs...)
}

func (_ *EngineDiag) LogMisc(fmt string, fmtArgs ...interface{}) {
	Diag.Log(DiagLogCatMisc, fmt, fmtArgs...)
}

func (_ *EngineDiag) LogShaders(fmt string, fmtArgs ...interface{}) {
	Diag.Log(DiagLogCatShaders, fmt, fmtArgs...)
}
