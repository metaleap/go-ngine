package core

import (
	"log"

	ugl "github.com/go3d/go-opengl/util"
	ugo "github.com/metaleap/go-util"
)

type EngineDiagLogCategory int

const (
	DiagLogCatMisc    EngineDiagLogCategory = 1
	DiagLogMeshes     EngineDiagLogCategory = 2
	DiagLogCatShaders EngineDiagLogCategory = 4
	DiagLogCatAll     EngineDiagLogCategory = DiagLogMeshes | DiagLogCatMisc | DiagLogCatShaders
)

var (
	Diag = EngineDiag{
		LogCategories:       DiagLogCatAll,
		LogErrorsDuringLoop: true,
	}
)

type EngineDiag struct {
	LogCategories       EngineDiagLogCategory
	LogErrorsDuringLoop bool
}

func (me *EngineDiag) Log(cat EngineDiagLogCategory, fmt string, fmtArgs ...interface{}) {
	if (me.LogCategories & cat) == cat {
		log.Printf(fmt, fmtArgs...)
	}
}

func (me *EngineDiag) LogErr(err error) {
	ugo.LogError(err)
}

func (me *EngineDiag) LogIfGlErr(fmt string, fmtArgs ...interface{}) {
	ugl.LogLastError(fmt, fmtArgs...)
}

func (me *EngineDiag) LogMeshes(fmt string, fmtArgs ...interface{}) {
	me.Log(DiagLogMeshes, fmt, fmtArgs...)
}

func (me *EngineDiag) LogMisc(fmt string, fmtArgs ...interface{}) {
	me.Log(DiagLogCatMisc, fmt, fmtArgs...)
}

func (me *EngineDiag) LogShaders(fmt string, fmtArgs ...interface{}) {
	me.Log(DiagLogCatShaders, fmt, fmtArgs...)
}
