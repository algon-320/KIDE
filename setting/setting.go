package setting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/algon-320/KIDE/util"
)

const (
	// SettingFilename ...設定ファイル(JSON)のファイル名
	SettingFilename = "settings.json"
)

var wrapper map[string]interface{} = make(map[string]interface{})

func init() {
	// 実行ファイルのディレクトリに移動
	prevDir, _ := os.Getwd()
	defer os.Chdir(prevDir)
	exeDir, _ := os.Executable()
	exeDir = filepath.Dir(exeDir)
	os.Chdir(exeDir)

	if !util.FileExists(SettingFilename) {
		fmt.Println(util.PrefixCaution + "No setting file `" + SettingFilename + "`.")
		return
	}

	bytes, err := ioutil.ReadFile(SettingFilename)
	if err != nil {
		fmt.Println(util.PrefixError+"File open error:", err)
		return
	}

	err = json.Unmarshal(bytes, &wrapper)
	if err != nil {
		fmt.Println(util.PrefixError+"JSON Unmarshal error:", err)
		return
	}
}

// Get ... selectorで指定された設定を返す(環境変数がある場合はそれを返す)
// selector: `.`区切りで指定する envVarKey: ここで指定した環境変数の値あるならそれを返す
func Get(selector string, envVarKey string) (interface{}, bool) {
	// 環境変数が設定されている場合はその値を返す
	if envVarKey != "" {
		if v, fnd := os.LookupEnv(envVarKey); fnd {
			util.DebugPrint("setting loaded from environment variable `" + envVarKey + "`")
			return v, true
		}
	}

	sel := strings.Split(selector, ".")
	cur := wrapper
	for i := 0; i < len(sel)-1; i++ {
		v := sel[i]
		m, ok := cur[v]
		if !ok {
			return nil, false
		}
		cur = m.(map[string]interface{})
	}

	v := sel[len(sel)-1]
	ret, ok := cur[v]
	if !ok {
		return nil, false
	}
	return ret, true
}

// Set ... selectorで指定された設定に値を上書きする
// selector: `.`区切りで指定する value:書き込む値
func Set(selector string, value interface{}) error {
	sel := strings.Split(selector, ".")
	cur := wrapper
	for i := 0; i < len(sel)-1; i++ {
		v := sel[i]
		m, ok := cur[v]
		if !ok {
			cur[v] = make(map[string]interface{})
			m = cur[v]
		}
		cur = m.(map[string]interface{})
	}
	v := sel[len(sel)-1]
	cur[v] = value
	save()
	return nil
}

func save() {
	// 実行ファイルのディレクトリに移動
	prevDir, _ := os.Getwd()
	defer os.Chdir(prevDir)
	exeDir, _ := os.Executable()
	exeDir = filepath.Dir(exeDir)
	os.Chdir(exeDir)

	jsonBytes, err := json.Marshal(wrapper)
	if err != nil {
		fmt.Println(util.PrefixError+"JSON Marshal error:", err)
		return
	}

	var buf bytes.Buffer
	json.Indent(&buf, jsonBytes, "", "  ")

	err = ioutil.WriteFile(SettingFilename, buf.Bytes(), 0600)
	if err != nil {
		fmt.Println(util.PrefixError+"File write error:", err)
		return
	}

	util.DebugPrint("updated setteing")
}
