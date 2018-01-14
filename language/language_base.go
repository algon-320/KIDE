package language

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/algon-320/KIDE/util"
)

// languageBase ... 言語定義用(LanguageインターフェースのRunメソッド以外を実装済み)
type languageBase struct {
	name           string
	fileExtension  string
	compileCommand string // {SOURCEFILE_PATH} の部分がすべてソースコードのパスに置換される
	runningCommand string // {SOURCEFILE_PATH} の部分がすべてソースコードのパスに置換される
	commentBegin   string
	commentEnd     string
}

// Name ... 言語の名前を返す
func (l *languageBase) Name() string {
	return l.name
}
func (l *languageBase) String() string {
	return l.Name()
}

// FileExtension ... ソースファイルの拡張子
func (l *languageBase) FileExtension() string {
	return l.fileExtension
}

// CommentOut ... line で与えられた文字列をコメントアウトして返す
func (l *languageBase) CommentOut(line string) string {
	return l.commentBegin + line + l.commentEnd
}

// UnComment ... commentedLine で与えられたコメントアウトされた文字列のコメントを外して返す
func (l *languageBase) UnComment(commentedLine string) string {
	lenBegin := len(l.commentBegin)
	lenEnd := len(l.commentEnd)
	return commentedLine[lenBegin : len(commentedLine)-lenEnd]
}

// compile ... sourcePath で与えられたパスのソースコードをコンパイルする(変更がない場合は何もしない)
func (l *languageBase) compile(sourcePath string) error {
	skip, err := checkSkipCompile(sourcePath) // 変更があるか確認
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

// Run ... 実行
// input : 標準入力として与える文字列
// print : 標準出力、標準エラー出力を画面に出力するかどうか
// return : 実行結果の標準出力, この関数のエラー
func (l *languageBase) Run(sourcePath string, input string, print bool) (string, error) {
	if l.compileCommand != "" {
		if err := l.compile(sourcePath); err != nil {
			return "", err
		}
	}

	ret := new(bytes.Buffer)
	cmd := util.Command(strings.Replace(l.runningCommand, "{SOURCEFILE_PATH}", sourcePath, -1))

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

// utility ---------------------------------------------------------------------
const PREV_SOURCE = "previous.txt" // TODO : previous.txt の名前をsettingで変えられるようにするべき
func checkSkipCompile(sourcePath string) (bool, error) {
	sourcePathAbs, _ := filepath.Abs(sourcePath)
	edir, _ := os.Executable()
	dir := filepath.Dir(edir)
	prevSourcePath := filepath.Join(dir, PREV_SOURCE)

	if !util.FileExists(sourcePathAbs) {
		return false, fmt.Errorf(util.PrefixError + "No such source file.")
	}

	if util.FileExists(prevSourcePath) {
		res, err := util.IsSameFile(prevSourcePath, sourcePathAbs)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}
func copySourceFile(sourcePath string) error {
	p, _ := os.Executable()
	dir := filepath.Dir(p)

	sourcePathAbs, _ := filepath.Abs(sourcePath)
	if err := util.FileCopy(sourcePathAbs, filepath.Join(dir, PREV_SOURCE)); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf(util.PrefixInfo+"Copied the souce file to `%s`", PREV_SOURCE))
	return nil
}
