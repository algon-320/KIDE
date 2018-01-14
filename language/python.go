package language

import "github.com/algon-320/KIDE/setting"

// PYTHON ... Python
var PYTHON Language

const (
	defaultCompileCommandPYTHON = ""
	defaultRunningCommandPYTHON = "python {SOURCEFILE_PATH}"
)

func init() {
	var compileCmd, runningCmd string
	if v, ok := setting.Get("Language.Python.CompileCommand", ""); ok {
		compileCmd = v.(string)
	} else {
		compileCmd = defaultCompileCommandPYTHON
		setting.Set("Language.Python.CompileCommand", compileCmd)
	}
	if v, ok := setting.Get("Language.Python.RunningCommand", ""); ok {
		runningCmd = v.(string)
	} else {
		runningCmd = defaultRunningCommandPYTHON
		setting.Set("Language.Python.RunningCommand", runningCmd)
	}

	PYTHON = &languageBase{
		name:           "Python",
		fileExtension:  ".py",
		compileCommand: defaultCompileCommandPYTHON,
		runningCommand: defaultRunningCommandPYTHON,
		commentBegin:   "# ",
		commentEnd:     "",
	}

}

type python struct {
	languageBase
}
