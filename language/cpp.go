package language

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/algon-320/KIDE/util"
)

// CPP ... C++14
var CPP = &cpp{
	name:           "C++",
	fileExtension:  ".cpp",
	compileCommand: "g++ -std=gnu++1y -O2 {SOURCEFILE_PATH} -o a.out",
	runCommand:     "./a.out",
}

type cpp struct {
	name           string
	fileExtension  string
	compileCommand string // `{SOURCEFILE_PATH}`の部分がソースファイルのパスに置き換わる
	runCommand     string
}

func (l *cpp) Name() string {
	return l.name
}

func (l *cpp) String() string {
	return l.name
}

func (l *cpp) FileExtension() string {
	return l.fileExtension
}

// Run ... 実行
// input : 標準入力として与える文字列
// print : 標準出力、標準エラー出力を画面に出力するかどうか
// return : 実行結果の標準出力, この関数のエラー
func (l *cpp) Run(sourcePath string, input string, print bool) (string, error) {
	if err := l.compile(sourcePath); err != nil {
		return "", err
	}

	ret := new(bytes.Buffer)
	cmd := util.Command(l.runCommand)

	var stdin io.Reader
	var stdout io.Writer
	var stderr io.Writer
	if input == "" {
		stdin = os.Stdin
	} else {
		stdin = bytes.NewBufferString(input)
	}

	if print {
		stdout = io.MultiWriter(os.Stdout, ret)
		stderr = os.Stderr
	} else {
		stdout = ret
		stderr = nil
	}

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return "", &ErrRuntimeError{}
	}

	return string(ret.Bytes()), nil
}

func (l *cpp) CommentOut(line string) string {
	return "// " + line
}
func (l *cpp) UnComment(line string) string {
	text := line[3:]
	return text
}

func (l *cpp) compile(sourcePath string) error {
	skip, err := checkSkipCompile(sourcePath)
	if err != nil {
		return err
	}
	if skip {
		util.DebugPrint("Source file isn't changed. Skip compiling.")
		return nil
	}

	util.DebugPrint("Source file is changed. Compiling ...")

	cmd := util.Command(strings.Replace(l.compileCommand, "{SOURCEFILE_PATH}", sourcePath, -1))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return &ErrCompileError{}
	}
	util.DebugPrint("Successfully compiled!")
	return copySourceFile(sourcePath)
}
