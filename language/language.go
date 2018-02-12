package language

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/algon-320/KIDE/util"
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

var languageList = []*Language{
	&CPP,
	&PYTHON,
	&JAVA,
}

// GetLanguage ... 言語名からLanguageを返す(仮) TODO
func GetLanguage(name string) Language {
	for _, lang := range languageList {
		if name == (*lang).Name() {
			return *lang
		}
	}
	fmt.Fprintln(os.Stderr, util.PrefixCaution+"unsupported language.")
	return CPP // TODO: settingから読むべき
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
