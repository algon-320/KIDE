package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

// FileExists ... `filepath`で与えられたパスが存在するかどうかを返す
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

// IsSameFile ...`path1`、`path2`のファイルの内容が同じかどうかを返す (バイト列比較)
func IsSameFile(path1, path2 string) (bool, error) {
	data1, err := ioutil.ReadFile(path1)
	if err != nil {
		return false, fmt.Errorf("cannot open `%s`", path1)
	}
	data2, err := ioutil.ReadFile(path2)
	if err != nil {
		return false, fmt.Errorf("cannot open `%s`", path2)
	}
	return reflect.DeepEqual(data1, data2), nil
}

// FileCopy ... `targetPath`のファイルを`destinationPath`にコピーする
func FileCopy(targetPath, destinationPath string) error {
	cont, err := ioutil.ReadFile(targetPath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(destinationPath, cont, 0644)
}

// Command ... ` `で分割されたコマンド文字列から`*exec.Cmd`を作って返す
func Command(spaceSeparatedCmd string) *exec.Cmd {
	cmdSep := strings.Split(spaceSeparatedCmd, " ")
	var cmd *exec.Cmd
	if len(cmdSep) == 1 {
		cmd = exec.Command(cmdSep[0])
	} else {
		cmd = exec.Command(cmdSep[0], cmdSep[1:]...)
	}
	return cmd
}
