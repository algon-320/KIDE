package util

import "fmt"

const (
	// Escape sequences to colorize text
	ESCS_COL_RED_B    = "\033[31;1m"
	ESCS_COL_GREEN_B  = "\033[32;1m"
	ESCS_COL_YELLOW_B = "\033[33;1m"
	ESCS_COL_CYAN_B   = "\033[34;1m"
	ESCS_COL_PURPLE_B = "\033[35;1m"
	ESCS_BOLD         = "\033[1m"
	ESCS_COL_OFF      = "\033[0m"

	// PrefixInfo    = "\033[34;1m[INFO]\033[0m "
	// PrefixDebug   = "\033[32;1m[DEBUG]\033[0m "
	// PrefixError   = "\033[31;1m[ERROR]\033[0m "
	// PrefixCaution = "\033[33;1m[CAUTION]\033[0m "
	PrefixInfo     = ESCS_COL_CYAN_B + "● " + ESCS_COL_OFF   // Prefix of Infomation message.
	PrefixDebug    = ESCS_COL_GREEN_B + "● " + ESCS_COL_OFF  // Prefix of Debug message.
	PrefixError    = ESCS_COL_RED_B + "● " + ESCS_COL_OFF    // Prefix of Error message.
	PrefixCaution  = ESCS_COL_YELLOW_B + "● " + ESCS_COL_OFF // Prefix of Caution message.
	PrefixQuestion = ESCS_COL_PURPLE_B + "● " + ESCS_COL_OFF // Prefix of Question message.
)

// DebugPrint ... デバックメッセージを出力する
func DebugPrint(text string) {
	// fmt.Println(PrefixDebug + text)
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
