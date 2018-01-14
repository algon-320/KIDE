package language

// JAVA ... Java
var JAVA = &languageBase{
	name:           "Java",
	fileExtension:  ".java",
	compileCommand: "javac {SOURCEFILE_PATH}",
	runningCommand: "java Main",
	commentBegin:   "// ",
	commentEnd:     "",
}

type java struct {
	languageBase
}
