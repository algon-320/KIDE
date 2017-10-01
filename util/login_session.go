package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"

	"github.com/headzoo/surf/util"
)

// LoadLoginSession ... セッションファイルを読み込む
func LoadLoginSession(filename string, cookieURL string) *cookiejar.Jar {
	// 実行ファイルのディレクトリに移動
	prevDir, _ := os.Getwd()
	defer os.Chdir(prevDir)
	exeDir, _ := os.Executable()
	exeDir = filepath.Dir(exeDir)
	os.Chdir(exeDir)

	if !util.FileExists(filename) {
		return nil
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(PrefixError+"File open error:", err)
		return nil
	}

	var cookies []*http.Cookie
	err = json.Unmarshal(bytes, &cookies)
	if err != nil {
		fmt.Println(PrefixError+"JSON Unmarshal error:", err)
		return nil
	}
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(cookieURL)
	jar.SetCookies(u, cookies)
	return jar
}

// SaveLoginSession .. セッションファイルを保存する
func SaveLoginSession(filename string, cookies []*http.Cookie) {
	// 実行ファイルのディレクトリに移動
	prev, _ := os.Getwd()
	defer os.Chdir(prev)
	exeDir, _ := os.Executable()
	exeDir = filepath.Dir(exeDir)
	os.Chdir(exeDir)

	bytes, err := json.Marshal(cookies)
	if err != nil {
		fmt.Println(PrefixError+"JSON Marshal error:", err)
		return
	}

	err = ioutil.WriteFile(filename, bytes, 0600)
	if err != nil {
		fmt.Println(PrefixError+"File write error:", err)
		return
	}
	fmt.Println(PrefixInfo + fmt.Sprintf("Saved session as `%s`", filename))
}
