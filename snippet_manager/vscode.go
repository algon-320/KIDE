package snippet_manager

import (
	"fmt"
	"strings"
)

// VScode ... Visual Studio Code
var VScode *vscode

type vscode struct{}

func (editor *vscode) Name() string {
	return "vscode"
}

func (editor *vscode) generateSnippets(snips []Snippet) string {
	temp := []Snippet{}
	for _, snip := range snips {
		snipTemp := Snippet{}
		// エスケープ
		for k := range snip {
			var str string
			for _, c := range snip[k] {
				if c == '"' {
					str += `\"`
				} else if c == '\\' {
					str += `\\`
				} else {
					str += string(c)
				}
			}
			snipTemp[k] = str
		}
		temp = append(temp, snipTemp)
	}

	var ret string
	for i, snip := range temp {
		ret += fmt.Sprintf("\"%s\": {\n", snip["NAME"])
		ret += fmt.Sprintf("\t\"description\": \"%s\",\n", snip["NAME"])
		ret += fmt.Sprintf("\t\"prefix\": \"%s\",\n", snip["TRIGGER"])
		ret += "\t\"body\": [\n"

		body := strings.Split(snip["CODE"], "\n")
		for j, l := range body {
			ret += fmt.Sprintf("\t\t\"%s\"", l)
			if j != len(body)-1 {
				ret += ",\n"
			} else {
				ret += "\n"
			}
		}

		ret += "\t]\n"
		ret += "}"
		if i != len(snips)-1 {
			ret += ",\n\n"
		} else {
			ret += "\n"
		}
	}
	return ret
}
