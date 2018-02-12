package snippet_manager

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/util"
)

const settingRootDir = "snippet_manager.root_dir"

// Snippet ... スニペットデータ
type Snippet map[string]string

// Editor ... スニペットを出力する対象のテキストエディタ
type Editor interface {
	Name() string
	generateSnippets([]Snippet) string
}

// EditorList ... 利用可能なエディタの一覧
var EditorList = []Editor{
	VScode,
}

var rootSnippetsDir string
var snippets []Snippet

// findSnip ... ディレクトリ`dir`以下から.snipのファイルを再帰的に検索してパスの一覧を返す
func findSnip(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, findSnip(filepath.Join(dir, file.Name()))...)
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext == ".snip" {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}
	return paths
}

func init() {
	tmp, ok := setting.Get(settingRootDir, "")
	if ok {
		rootSnippetsDir = tmp.(string)
	} else {
		for {
			tmp := util.AskString("スニペットファイルのある親ディレクトリを入力してください(絶対パス)")
			if util.FileExists(tmp) {
				rootSnippetsDir = tmp
				setting.Set(settingRootDir, rootSnippetsDir)
				break
			}
		}
	}

	// スニペットデータをロード
	snipList := findSnip(rootSnippetsDir)
	for _, f := range snipList {
		snip, err := parseSnippet(f)
		if err != nil {
			panic(err)
		}
		snippets = append(snippets, snip)
	}
}

// GenerateMarkdown ... Markdownを生成
func GenerateMarkdown() string {
	var md string
	for _, snip := range snippets {
		var str string
		str += fmt.Sprintf("## %s\n", snip["NAME"])
		str += fmt.Sprintf("### Note\n")
		str += fmt.Sprintf("%s\n", snip["NOTE"])
		str += fmt.Sprintf("### Code\n")
		str += fmt.Sprintf("```cpp\n")
		str += fmt.Sprintf("%s\n", snip["CODE"])
		str += fmt.Sprintf("```\n\n")
		md += str
	}
	return md
}

// GenerateHTML ... HTMLページを生成
func GenerateHTML() {
	// md := GenerateMarkdown()
}

// GenerateTex ... texファイルを生成
func GenerateTex() {
	// md := GenerateMarkdown()
}

// ExportSnippets ... スニペットを`editor`用の形式で出力する
func ExportSnippets(editor Editor) string {
	return editor.generateSnippets(snippets)
}
