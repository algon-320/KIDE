package language

import "github.com/algon-320/KIDE/setting"

// PYTHON2 ... Python
var PYTHON2, PYTHON3 Language

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

	PYTHON2 = &languageBase{
		name:           "Python2",
		fileExtension:  ".py",
		compileCommand: defaultCompileCommandPYTHON,
		runningCommand: defaultRunningCommandPYTHON,
		commentBegin:   "# ",
		commentEnd:     "",
	}
	PYTHON3 = &languageBase{
		name:           "Python3",
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
