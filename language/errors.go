package language

import (
	"fmt"

	"github.com/algon-320/KIDE/util"
)

type ErrNoSuchLanguage struct {
	ext_name string
}

func (e ErrNoSuchLanguage) Error() string {
	return util.PrefixError + fmt.Sprintf("Not suported source file `%s`.", e.ext_name)
}

type ErrCompileError struct {
}

func (e ErrCompileError) Error() string {
	return util.PrefixError + fmt.Sprintf("Compile Error.")
}

type ErrRuntimeError struct {
}

func (e ErrRuntimeError) Error() string {
	return util.PrefixError + fmt.Sprintf("Runtime Error.")
}

type ErrNoSourceCode struct {
	name string
}

func (e ErrNoSourceCode) Error() string {
	return util.PrefixError + fmt.Sprintf("No %s source file found.", e.name)
}
