package language

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

// FindSourceCode ... ソースコードのファイル名を返す
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

//--------------------------------------------------------------------

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
