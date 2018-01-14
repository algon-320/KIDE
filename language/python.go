package language

// PYTHON ... Python2
var PYTHON = &languageBase{
	name:           "Python2",
	fileExtension:  ".py",
	compileCommand: "",
	runningCommand: "python {SOURCEFILE_PATH}",
	commentBegin:   "# ",
	commentEnd:     "",
}

type python struct {
	languageBase
}
