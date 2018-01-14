package language

import "github.com/algon-320/KIDE/setting"

// JAVA ... Java
var JAVA Language

const (
	defaultCompileCommandJAVA = "javac {SOURCEFILE_PATH}"
	defaultRunningCommandJAVA = "java Main"
)

func init() {
	var compileCmd, runningCmd string
	if v, ok := setting.Get("Language.Java.CompileCommand", ""); ok {
		compileCmd = v.(string)
	} else {
		compileCmd = defaultCompileCommandJAVA
		setting.Set("Language.Java.CompileCommand", compileCmd)
	}
	if v, ok := setting.Get("Language.Java.RunningCommand", ""); ok {
		runningCmd = v.(string)
	} else {
		runningCmd = defaultRunningCommandJAVA
		setting.Set("Language.Java.RunningCommand", runningCmd)
	}

	JAVA = &languageBase{
		name:           "Java",
		fileExtension:  ".java",
		compileCommand: defaultCompileCommandJAVA,
		runningCommand: defaultRunningCommandJAVA,
		commentBegin:   "// ",
		commentEnd:     "",
	}
}

type java struct {
	languageBase
}
