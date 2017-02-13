package core

import (
	"log"

	"github.com/metaleap/go-util-misc"
	ugl "github.com/metaleap/go-opengl/util"
)

var (
	Diag = NgDiag{
		LogCategories: DiagLogCatAll,
		LogCategoryNames: map[NgDiagLogCategory]string{
			DiagLogCatMisc:    "[INFO]\t\t",
			DiagLogCatShaders: "[SHADER]\t",
			DiagLogCatImages:  "[IMAGES]\t",
			DiagLogCatMeshes:  "[MESHES]\t",
		},
		LogGLErrorsInLoopOnSec: false,
	}
)

//	NgDiag is a singleton type, only used for the core.Diag package-global exported variable.
//	It is only aware of that instance and does not support any other NgDiag instances.
type NgDiag struct {
	LogCategories          NgDiagLogCategory
	LogCategoryNames       map[NgDiagLogCategory]string
	LogGLErrorsInLoopOnSec bool
}

func (_ *NgDiag) Log(cat NgDiagLogCategory, fmt string, fmtArgs ...interface{}) {
	if (Diag.LogCategories & cat) == cat {
		log.Printf(Diag.LogCategoryNames[cat]+fmt, fmtArgs...)
	}
}

func (_ *NgDiag) LogErr(err error) {
	ugo.LogError(err)
}

func (_ *NgDiag) LogIfGlErr(fmt string, fmtArgs ...interface{}) {
	ugl.LogLastError(fmt, fmtArgs...)
}

func (_ *NgDiag) LogImages(fmt string, fmtArgs ...interface{}) {
	Diag.Log(DiagLogCatImages, fmt, fmtArgs...)
}

func (_ *NgDiag) LogMeshes(fmt string, fmtArgs ...interface{}) {
	Diag.Log(DiagLogCatMeshes, fmt, fmtArgs...)
}

func (_ *NgDiag) LogMisc(fmt string, fmtArgs ...interface{}) {
	Diag.Log(DiagLogCatMisc, fmt, fmtArgs...)
}

func (_ *NgDiag) LogShaders(fmt string, fmtArgs ...interface{}) {
	Diag.Log(DiagLogCatShaders, fmt, fmtArgs...)
}

type NgDiagLogCategory int

const (
	DiagLogCatMisc    NgDiagLogCategory = 1
	DiagLogCatMeshes  NgDiagLogCategory = 2
	DiagLogCatShaders NgDiagLogCategory = 4
	DiagLogCatImages  NgDiagLogCategory = 8
	DiagLogCatAll     NgDiagLogCategory = DiagLogCatMeshes | DiagLogCatMisc | DiagLogCatShaders | DiagLogCatImages
)

func init() {
	ugo.LogErrorFormat = "[ERROR!]\t%v"
}
