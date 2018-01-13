package language

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/algon-320/KIDE/util"
)

// JAVA ... Java
var JAVA = &java{
	name:           "Java",
	fileExtension:  ".java",
	compileCommand: "javac {SOURCEFILE_PATH}",
	runCommand:     "java Main",
}

type java struct {
	name           string
	fileExtension  string
	compileCommand string // `{SOURCEFILE_PATH}`の部分がソースファイルのパスに置き換わる
	runCommand     string
}

func (l *java) Name() string {
	return l.name
}

func (l *java) String() string {
	return l.name
}

func (l *java) FileExtension() string {
	return l.fileExtension
}

// Run ... 実行
// input : 標準入力として与える文字列
// print : 標準出力、標準エラー出力を画面に出力するかどうか
// return : 実行結果の標準出力, この関数のエラー
func (l *java) Run(sourcePath string, input string, print bool) (string, error) {
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

func (l *java) CommentOut(line string) string {
	return "// " + line
}
func (l *java) UnComment(line string) string {
	text := line[3:]
	return text
}

func (l *java) compile(sourcePath string) error {
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