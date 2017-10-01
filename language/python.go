package language

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/algon-320/KIDE/util"
)

// PYTHON ... Python2
var PYTHON = &python{
	name:          "Python2",
	fileExtension: ".py",
	runCommand:    "python {SOURCEFILE_PATH}",
}

type python struct {
	name          string
	fileExtension string
	runCommand    string // `{SOURCEFILE_PATH}`の部分がソースファイルのパスに置き換わる
}

func (l *python) Name() string {
	return l.name
}

func (l *python) FileExtension() string {
	return l.fileExtension
}

// Run ... 実行
// input : 標準入力として与える文字列
// print : 標準出力、標準エラー出力を画面に出力するかどうか
// return : 実行結果の標準出力, エラー
func (l *python) Run(sourcePath string, input string, print bool) (string, error) {
	ret := new(bytes.Buffer)
	l.runCommand = strings.Replace(l.runCommand, "{SOURCEFILE_PATH}", sourcePath, -1)
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

func (l *python) CommentOut(text string) string {
	return "# " + text
}
func (l *python) UnComment(comment string) string {
	text := comment[2:]
	return text
}

func (l *python) String() string {
	return l.name
}
