package language

// CPP ... C++14
var CPP = &LanguageBase{
	name:           "C++",
	fileExtension:  ".cpp",
	compileCommand: "g++ -std=gnu++1y -O2 {SOURCEFILE_PATH} -o a.out",
	runningCommand: "./a.out",
	commentBegin:   "// ",
	commentEnd:     "",
}

type cpp struct {
	LanguageBase
}
