package util

import (
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// ShiftJIS2UTF8 ... Shift_JISでエンコードされた文字列をUTF8に直す
func ShiftJIS2UTF8(str string) (string, error) {
	strReader := strings.NewReader(str)
	decodedReader := transform.NewReader(strReader, japanese.ShiftJIS.NewDecoder())
	decoded, err := ioutil.ReadAll(decodedReader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// AddBR ... 文字列の最後に改行がなければ加える
func AddBR(s string) string {
	s = strings.TrimRight(s, "\n")
	s += "\n"
	return s
}
