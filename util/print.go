package util

import (
	"fmt"
	"strings"
)

const (
	// Escape sequences to colorize text
	ESCS_COL_RED_B    = "\033[31;1m"
	ESCS_COL_GREEN_B  = "\033[32;1m"
	ESCS_COL_YELLOW_B = "\033[33;1m"
	ESCS_COL_CYAN_B   = "\033[34;1m"
	ESCS_COL_PURPLE_B = "\033[35;1m"
	ESCS_BOLD         = "\033[1m"
	ESCS_COL_OFF      = "\033[0m"

	PrefixInfo     = ESCS_COL_CYAN_B + "● " + ESCS_COL_OFF   // 水色
	PrefixDebug    = ESCS_COL_GREEN_B + "● " + ESCS_COL_OFF  // 緑
	PrefixError    = ESCS_COL_RED_B + "● " + ESCS_COL_OFF    // 赤
	PrefixCaution  = ESCS_COL_YELLOW_B + "● " + ESCS_COL_OFF // 黄色
	PrefixQuestion = ESCS_COL_PURPLE_B + "● " + ESCS_COL_OFF // 紫
)

// DebugPrint ... デバックメッセージを出力する
func DebugPrint(text string) {
	// fmt.Fprintln(os.Stderr, PrefixDebug + text)
}

// PrintTable ... テーブルを書く
// border: trueにすると境界線付きで出力する
func PrintTable(title []string, data [][]string, border bool) {
	width := []int{}

	n := len(data)
	m := len(title)
	for _, v := range title {
		width = append(width, len(v))
	}

	for i := 0; i < n; i++ {
		for j := 0; j < len(data[i]); j++ {
			if width[j] < len(data[i][j]) {
				width[j] = len(data[i][j])
			}
		}
	}

	for j := 0; j < m; j++ {
		fmt.Print(ESCS_BOLD + title[j] + ESCS_COL_OFF)
		amount := width[j] - len(title[j]) + 2
		for ; amount > 0; amount-- {
			fmt.Print(" ")
		}
	}
	fmt.Print("\n")

	if border {
		for j := 0; j < m; j++ {
			for k := 0; k < width[j]+2; k++ {
				fmt.Print("-")
			}
		}
		fmt.Print("\n")
	}

	for i := 0; i < n; i++ {
		for j := 0; j < len(data[i]); j++ {
			fmt.Print(data[i][j])
			amount := width[j] - len(data[i][j]) + 2
			for ; amount > 0; amount-- {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

// SprintTitle ... タイトルテキストを装飾付きで表示して改行した文字列を返す
//     例:) SprintTitle(20, 5, "-+", "Hello,World") --> "-+-+- Hello,World -+-+-+-+\n"
// width: 全体の幅
// prefixWidth: textの前の装飾の幅
// ornament: 装飾として用いる文字列
// text: 表示する文字列
func SprintTitle(width int, prefixWidth int, ornament string, text string) string {
	if width < len(text) {
		return ""
	}
	front := strings.Repeat(ornament, prefixWidth/len(ornament)+1)[:prefixWidth]
	back := strings.Repeat(ornament, (width-prefixWidth-len(text)-2)/len(ornament)+1)[:width-prefixWidth-len(text)-2]
	return fmt.Sprintf("%s %s %s\n", front, text, back)
}

// SprintTitlef ... SprintTitleのフォーマット文字列版
func SprintTitlef(width int, prefixWidth int, ornament string, ftext string, args ...interface{}) string {
	text := fmt.Sprintf(ftext, args...)
	return SprintTitle(width, prefixWidth, ornament, text)
}

// PrintTitle ... タイトルテキストを装飾付きで表示して改行する
//     例:) PrintTitle(20, 5, "-+", "Hello,World") --> "-+-+- Hello,World -+-+-+-+"
// width: 全体の幅
// prefixWidth: textの前の装飾の幅
// ornament: 装飾として用いる文字列
// text: 表示する文字列
func PrintTitle(width int, prefixWidth int, ornament string, text string) {
	fmt.Print(SprintTitle(width, prefixWidth, ornament, text))
}

// PrintTitlef ... PrintTitleのフォーマット文字列版
func PrintTitlef(width int, prefixWidth int, ornament string, ftext string, args ...interface{}) {
	fmt.Print(SprintTitlef(width, prefixWidth, ornament, ftext, args...))
}
