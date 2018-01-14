package language

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/algon-320/KIDE/util"
)

const (
	PREV_SOURCE = "previous.txt" // TODO : previous.txt の名前をsettingで変えられるようにするべき
)

// Language ... 言語定義用インターフェース
type Language interface {
	Name() string
	String() string
	FileExtension() string
	Run(sourcePath string, input string, print bool) (string, error)
	CommentOut(line string) string
	UnComment(line string) string
}

// LanguageBase ... 言語定義用(LanguageインターフェースのRunメソッド以外を実装済み)
type LanguageBase struct {
	name           string
	fileExtension  string
	compileCommand string // {SOURCEFILE_PATH} の部分がすべてソースコードのパスに置換される
	runningCommand string // {SOURCEFILE_PATH} の部分がすべてソースコードのパスに置換される
	commentBegin   string
	commentEnd     string
}

// Name ... 言語の名前を返す
func (l *LanguageBase) Name() string {
	return l.name
}
func (l *LanguageBase) String() string {
	return l.Name()
}

// FileExtension ... ソースファイルの拡張子
func (l *LanguageBase) FileExtension() string {
	return l.fileExtension
}

// CommentOut ... line で与えられた文字列をコメントアウトして返す
func (l *LanguageBase) CommentOut(line string) string {
	return l.commentBegin + line + l.commentEnd
}

// UnComment ... commentedLine で与えられたコメントアウトされた文字列のコメントを外して返す
func (l *LanguageBase) UnComment(commentedLine string) string {
	lenBegin := len(l.commentBegin)
	lenEnd := len(l.commentEnd)
	return commentedLine[lenBegin : len(commentedLine)-lenEnd]
}

// compile ... sourcePath で与えられたパスのソースコードをコンパイルする(変更がない場合は何もしない)
func (l *LanguageBase) compile(sourcePath string) error {
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
func (l *LanguageBase) Run(sourcePath string, input string, print bool) (string, error) {
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

// GetLanguage ... 言語名からLanguageを返す(仮) TODO
func GetLanguage(name string) Language {
	switch name {
	case "c++":
		return CPP
	case "python":
		return PYTHON
	case "java":
		return JAVA
	default:
		fmt.Println(util.PrefixCaution + "unsupported language.")
		return CPP // TODO: settingから読むべき
	}
}

// FindSourceCode ... 実行する対象のソースコードを決めて、ファイル名を返す(複数ある場合はユーザに問う)
func FindSourceCode(lang Language) (string, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return "", err
	}

	// ディレクトリにある lang のソースコードのリストをつくる
	list := []string{}
	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == lang.FileExtension() {
				list = append(list, file.Name())
			}
		}
	}

	if len(list) == 1 {
		util.DebugPrint("unique source file : " + list[0])
		return list[0], nil
	} else if len(list) == 0 {
		util.DebugPrint("no source file")
		return "", &ErrNoSourceCode{name: lang.Name()}
	} else {
		util.DebugPrint("show selection ... choose file")
		sel := util.AskChoose(list, "Choose source file")
		return list[sel], nil
	}
}

// utility ---------------------------------------------------------------------
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
