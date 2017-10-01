package online_judge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/algon-320/KIDE/util"
)

const (
	// SamplecaseDir ... 問題のJSONが保存されるディレクトリ
	SamplecaseDir = "samplecases"
	// SamplecaseFilename ... 問題のJSONのファイル名
	SamplecaseFilename = `problem_%s.json` // %sにIDを埋め込む
)

// TestCase ... サンプルケースの入出力
type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Problem struct {
	ID        string      `json:"id"`
	ContestID string      `json:"contest_id"`
	Name      string      `json:"name"`
	URL       string      `json:"url"`
	Oj        OnlineJudge `json:"oj"`
	Cases     []TestCase  `json:"cases"`
}

// TODO : String() にするべき
func (p *Problem) Print() {
	fmt.Println("id:", p.ID)
	fmt.Println("name:", p.Name)
	fmt.Println("contest_id:", p.ContestID)
	fmt.Println("url:", p.URL)
	fmt.Println("oj:", p.Oj.Name())
	for i, tc := range p.Cases {
		fmt.Printf("==== sample case %d =========\n", i)
		fmt.Println("-------- Input ---------")
		fmt.Print(tc.Input)
		fmt.Println("-------- Output --------")
		fmt.Print(tc.Output)
	}
	fmt.Println("============================")
}

// Save ... ファイルに保存する
func (p *Problem) Save() error {
	// 実行ファイルのディレクトリに移動
	prevDir, _ := os.Getwd()
	defer os.Chdir(prevDir)
	exeDir, _ := os.Executable()
	exeDir = filepath.Dir(exeDir)
	os.Chdir(exeDir)

	p.ID = strings.ToUpper(p.ID)
	filename := fmt.Sprintf(SamplecaseFilename, p.ID)

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// サンプルケースのフォルダがなければ作る
	if !util.FileExists(SamplecaseDir) {
		if err = os.Mkdir(SamplecaseDir, 0775); err != nil {
			return err
		}
		absPath, _ := filepath.Abs(SamplecaseDir)
		fmt.Println(fmt.Sprintf(util.PrefixInfo+"Created a directory `%s`", absPath))
	}
	os.Chdir(SamplecaseDir)

	var buf bytes.Buffer
	json.Indent(&buf, jsonBytes, "", "  ")
	err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	fmt.Println(util.PrefixInfo + "Save problem : " + p.ID)
	return nil
}

// LoadProblem ... id で指定された問題を読み込む
func LoadProblem(id string) (*Problem, error) {
	id = strings.ToUpper(id)

	// 実行ファイルのディレクトリ/{SamplecaseDir}に移動
	prevDir, _ := os.Getwd()
	defer os.Chdir(prevDir)
	exeDir, _ := os.Executable()
	exeDir = filepath.Dir(exeDir)
	os.Chdir(exeDir)
	os.Chdir(SamplecaseDir)

	filename := fmt.Sprintf(SamplecaseFilename, id)

	if !util.FileExists(filename) {
		return nil, &ErrFailedToLoadSamplecase{message: filename + " dosen't exist."}
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, &ErrFailedToLoadSamplecase{message: "failed to open " + filename}
	}

	var tmp struct {
		ID        string     `json:"id"`
		ContestID string     `json:"contest_id"`
		Name      string     `json:"name"`
		URL       string     `json:"url"`
		Oj        string     `json:"oj"`
		Cases     []TestCase `json:"cases"`
	}
	err = json.Unmarshal(bytes, &tmp)
	if err != nil {
		return nil, err
	}
	oj, err := FromName(tmp.Oj)
	if err != nil {
		return nil, err
	}
	p := &Problem{
		ID:        tmp.ID,
		ContestID: tmp.ContestID,
		Name:      tmp.Name,
		URL:       tmp.URL,
		Oj:        oj,
		Cases:     tmp.Cases,
	}

	util.DebugPrint("Load problem : " + id)
	return p, nil
}

// GetAllProblemID ... 保存済みの問題の一覧を返す
func GetAllProblemID() []string {
	// 実行ファイルのディレクトリ/{SamplecaseDir}に移動
	prevDir, _ := os.Getwd()
	defer os.Chdir(prevDir)
	exeDir, _ := os.Executable()
	exeDir = filepath.Dir(exeDir)
	os.Chdir(exeDir)
	os.Chdir(SamplecaseDir)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	ret := []string{}

	pat := strings.Replace(SamplecaseFilename, `.`, `\.`, -1)
	pat = strings.Replace(pat, `%s`, `(.*)`, -1)

	re := regexp.MustCompile(pat)
	for _, file := range files {
		if !file.IsDir() {
			group := re.FindSubmatch([]byte(file.Name()))
			if group != nil {
				ret = append(ret, string(group[1]))
			}
		}
	}

	return ret
}
