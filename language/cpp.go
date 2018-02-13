package language

import "github.com/algon-320/KIDE/setting"

// CPP ... C++14
var CPP Language

const (
	defaultCompileCommandCPP = "g++ -std=c++11 -o a.out {SOURCEFILE_PATH}"
	defaultRunningCommandCPP = "./a.out"
)

func init() {
	var compileCmd, runningCmd string
	if v, ok := setting.Get("Language.C++.CompileCommand", ""); ok {
		compileCmd = v.(string)
	} else {
		compileCmd = defaultCompileCommandCPP
		setting.Set("Language.C++.CompileCommand", compileCmd)
	}
	if v, ok := setting.Get("Language.C++.RunningCommand", ""); ok {
		runningCmd = v.(string)
	} else {
		runningCmd = defaultRunningCommandCPP
		setting.Set("Language.C++.RunningCommand", runningCmd)
	}

	CPP = &languageBase{
		name:           "C++",
		fileExtension:  ".cpp",
		compileCommand: compileCmd,
		runningCommand: runningCmd,
		commentBegin:   "// ",
		commentEnd:     "",
	}
}

type cpp struct {
	languageBase
}
