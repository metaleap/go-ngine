package core

import (
	"log"

	ugo "github.com/metaleap/go-util"
)

type LogCategory int

const (
	LogCatMisc    LogCategory = 1
	LogCatMeshes  LogCategory = 2
	LogCatShaders LogCategory = 4
	LogCatAll     LogCategory = LogCatMeshes | LogCatMisc | LogCatShaders
)

var (
	Diag = EngineDiagnostics{
		LogCategories:       LogCatAll,
		LogErrorsDuringLoop: true,
	}
)

type EngineDiagnostics struct {
	LogCategories       LogCategory
	LogErrorsDuringLoop bool
}

func (me *EngineDiagnostics) Log(cat LogCategory, fmt string, fmtArgs ...interface{}) {
	if (me.LogCategories & cat) == cat {
		log.Printf(fmt, fmtArgs...)
	}
}

func (me *EngineDiagnostics) LogErr(err error) {
	ugo.LogError(err)
}

func (me *EngineDiagnostics) LogMeshes(fmt string, fmtArgs ...interface{}) {
	me.Log(LogCatMeshes, fmt, fmtArgs...)
}

func (me *EngineDiagnostics) LogMisc(fmt string, fmtArgs ...interface{}) {
	me.Log(LogCatMisc, fmt, fmtArgs...)
}

func (me *EngineDiagnostics) LogShaders(fmt string, fmtArgs ...interface{}) {
	me.Log(LogCatShaders, fmt, fmtArgs...)
}
