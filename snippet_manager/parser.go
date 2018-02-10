package snippet_manager

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/algon-320/KIDE/util"
)

func parseSnippet(filePath string) (Snippet, error) {
	parse := func(lines []string) (Snippet, error) {
		snip := Snippet{}
		var multilineTag string
		for _, l := range lines {
			l += string("\n")
			var tagname, content string
			if len(l) > 1 && l[0] == '<' {
				flag := 0
				for _, c := range l {
					if (flag == 0 && c == '<') || (flag == 1 && c == '>') || (flag == 2 && c == ' ') {
						flag++
						continue
					}
					if flag == 1 {
						tagname += string(c)
					}
					if flag == 2 || flag == 3 {
						content += string(c)
					}
				}
				if flag != 0 && flag != 2 && flag != 3 {
					return nil, fmt.Errorf(util.PrefixError + "invalid format: " + l)
				}
				if len(tagname) > 1 {
					if tagname[0] == '*' {
						tagname = tagname[1:]
						multilineTag = tagname
					} else {
						multilineTag = ""
					}
				}
			}
			if multilineTag != "" {
				if tagname == "" {
					snip[multilineTag] += l
				}
			} else if tagname != "" {
				snip[tagname] = content
			}
		}
		// 末尾の改行削除
		for k, v := range snip {
			snip[k] = v[:len(v)-1]
		}
		return snip, nil
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return parse(strings.Split(string(data), "\n"))
}
