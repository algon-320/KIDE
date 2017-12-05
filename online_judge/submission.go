package online_judge

import (
	"fmt"
	"time"

	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/util"
)

// JudgeStatus ... ジャッジの結果
type JudgeStatus int

const (
	// JudgeStatusUNK ... Unknown
	JudgeStatusUNK JudgeStatus = -1
	// JudgeStatusAC ... Accepted
	JudgeStatusAC JudgeStatus = iota
	// JudgeStatusPP ... Pretests passed
	JudgeStatusPP
	// JudgeStatusWA ... Wrong answer
	JudgeStatusWA
	// JudgeStatusCE ... Compilation error
	JudgeStatusCE
	// JudgeStatusRE ... Runtime error
	JudgeStatusRE
	// JudgeStatusTLE ... Time limit exceeded
	JudgeStatusTLE
	// JudgeStatusMLE ... Memory limit exceeded
	JudgeStatusMLE
	// JudgeStatusOLE ... Output limit exceeded
	JudgeStatusOLE
	// JudgeStatusIE ... Internal error
	JudgeStatusIE
)

// ToString ... ジャッジ結果を文字列化します
func (js JudgeStatus) ToString() string {
	switch js {
	case JudgeStatusAC:
		return "Accepted"
	case JudgeStatusPP:
		return "Pretest passed"
	case JudgeStatusWA:
		return "Wrong answer"
	case JudgeStatusCE:
		return "Compilation error"
	case JudgeStatusRE:
		return "Runtime error"
	case JudgeStatusTLE:
		return "Time limit exceeded"
	case JudgeStatusMLE:
		return "Memory limit exceeded"
	case JudgeStatusOLE:
		return "Output limit exceeded"
	case JudgeStatusIE:
		return "Internal error"

	case JudgeStatusUNK:
		fallthrough
	default:
		return util.PrefixCaution + "Unknown"
	}
}

// String ... エスケープシーケンス付きで文字列化(出力用)
func (js JudgeStatus) String() string {
	switch js {
	case JudgeStatusAC:
		return util.ESCS_COL_GREEN_B + js.ToString() + util.ESCS_COL_OFF
	case JudgeStatusPP:
		return util.ESCS_COL_GREEN_B + "Pretest passed" + util.ESCS_COL_OFF
	case JudgeStatusWA:
		return util.ESCS_COL_RED_B + "Wrong answer" + util.ESCS_COL_OFF
	case JudgeStatusCE:
		return util.ESCS_COL_PURPLE_B + "Compilation error" + util.ESCS_COL_OFF
	case JudgeStatusRE:
		return util.ESCS_COL_YELLOW_B + "Runtime error" + util.ESCS_COL_OFF
	case JudgeStatusTLE:
		return util.ESCS_COL_YELLOW_B + "Time limit exceeded" + util.ESCS_COL_OFF
	case JudgeStatusMLE:
		return util.ESCS_COL_YELLOW_B + "Memory limit exceeded" + util.ESCS_COL_OFF
	case JudgeStatusOLE:
		return util.ESCS_COL_YELLOW_B + "Output limit exceeded" + util.ESCS_COL_OFF
	case JudgeStatusIE:
		return util.ESCS_COL_PURPLE_B + "Internal error" + util.ESCS_COL_OFF

	case JudgeStatusUNK:
		fallthrough
	default:
		return "\033[31;1;7m" + util.PrefixCaution + "Unknown" + util.ESCS_COL_OFF
	}
}

// CheckInterval ... ジャッジ結果を確認する間隔
const CheckInterval = 5 * time.Second

// JudgeResult ... ジャッジの詳細データ
type JudgeResult struct {
	Problem  *Problem
	Code     string
	Language language.Language
	Date     time.Time
	URL      string
	Status   JudgeStatus
}

// Print ... ジャッジの詳細を出力する TODO: 文字列で返すようにするべき(またはString()を実装する)
func (res *JudgeResult) Print() {
	util.PrintTitle(80, 4, "#", "JudgeResult")
	// util.PrintTitle(30, 4, "=", "Problem")
	// res.Problem.Print()
	util.PrintTitle(80, 4, "=", "Language")
	fmt.Println(res.Language)
	util.PrintTitle(80, 4, "=", "SourceCode")
	fmt.Println(res.Code)
	util.PrintTitle(80, 4, "=", "Date")
	fmt.Println(res.Date)
	util.PrintTitle(80, 4, "=", "URL")
	fmt.Println(res.URL)
	util.PrintTitle(80, 4, "=", "Status")
	fmt.Println(res.Status)
}
