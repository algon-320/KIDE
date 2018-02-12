package util

import (
	"fmt"
	"os"
	"strings"
)

// AskChoose ... 選択肢を表示してユーザに選ばせる
// choice: 選択肢, message: 表示するメッセージ "PrefixQuestion {message} > {ここにユーザが入力}"の形式で表示される
// return: インデックス (0-indexed)
func AskChoose(choice []string, message string) int {
	for {
		for i, v := range choice {
			fmt.Fprintf(os.Stderr, "%d) %s\n", i+1, v)
		}
		var ans int
		fmt.Fprint(os.Stderr, PrefixQuestion+message+" > ")
		fmt.Scan(&ans)
		if 1 <= ans && ans <= len(choice) {
			return ans - 1
		}

		fmt.Fprintf(os.Stderr, PrefixCaution+"please choose the number from 1 to %d. try again ...\n", len(choice))
	}
}

// AskYesNo ... YesかNoかをユーザに問う
// return: Yesならtrue, Noならfalse
func AskYesNo() bool {
	fmt.Fprint(os.Stderr, PrefixQuestion+"Please respond with 'yes' or 'no' [y/N]: ")
	var resp string
	fmt.Scan(&resp)
	resp = strings.ToLower(resp)
	if resp == "y" || resp == "ye" || resp == "yes" {
		return true
	}
	return false
}

// AskString ... ユーザに文字列を入力させる
// message: 表示されるメッセージ "PrefixQuestion {message} > {ここにユーザが入力}"の形式で表示される
// return: 入力された文字列
func AskString(message string) string {
	fmt.Fprint(os.Stderr, PrefixQuestion+message+" > ")
	var ans string
	fmt.Scanln(&ans)
	return ans
}
