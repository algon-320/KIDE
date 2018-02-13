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
	if v, ok := setting.Get("Language.Python2.CompileCommand", ""); ok {
		compileCmd = v.(string)
	} else {
		compileCmd = defaultCompileCommandPYTHON
		setting.Set("Language.Python2.CompileCommand", compileCmd)
	}
	if v, ok := setting.Get("Language.Python2.RunningCommand", ""); ok {
		runningCmd = v.(string)
	} else {
		runningCmd = defaultRunningCommandPYTHON
		setting.Set("Language.Python2.RunningCommand", runningCmd)
	}
	PYTHON2 = &languageBase{
		name:           "Python2",
		fileExtension:  ".py",
		compileCommand: defaultCompileCommandPYTHON,
		runningCommand: defaultRunningCommandPYTHON,
		commentBegin:   "# ",
		commentEnd:     "",
	}

	compileCmd = ""
	runningCmd = ""
	if v, ok := setting.Get("Language.Python3.CompileCommand", ""); ok {
		compileCmd = v.(string)
	} else {
		compileCmd = defaultCompileCommandPYTHON
		setting.Set("Language.Python3.CompileCommand", compileCmd)
	}
	if v, ok := setting.Get("Language.Python3.RunningCommand", ""); ok {
		runningCmd = v.(string)
	} else {
		runningCmd = defaultRunningCommandPYTHON
		setting.Set("Language.Python3.RunningCommand", runningCmd)
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
