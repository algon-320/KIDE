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
	labelMap := map[JudgeStatus]string{
		JudgeStatusAC:  "Accepted",
		JudgeStatusPP:  "Pretest passed",
		JudgeStatusWA:  "Wrong answer",
		JudgeStatusCE:  "Compilation error",
		JudgeStatusRE:  "Runtime error",
		JudgeStatusTLE: "Time limit exceeded",
		JudgeStatusMLE: "Memory limit exceeded",
		JudgeStatusOLE: "Output limit exceeded",
		JudgeStatusIE:  "Internal error",
		JudgeStatusUNK: "Unknown",
	}
	label, ok := labelMap[js]
	if !ok {
		return "Unknown"
	}
	return label
}

// GetColorESCS ... ジャッジ結果に対応する色のエスケープシーケンスを返す
func (js JudgeStatus) GetColorESCS() string {
	colorMap := map[JudgeStatus]string{
		JudgeStatusAC:  util.ESCS_COL_GREEN_B,
		JudgeStatusPP:  util.ESCS_COL_GREEN_B,
		JudgeStatusWA:  util.ESCS_COL_RED_B,
		JudgeStatusCE:  util.ESCS_COL_PURPLE_B,
		JudgeStatusRE:  util.ESCS_COL_YELLOW_B,
		JudgeStatusTLE: util.ESCS_COL_YELLOW_B,
		JudgeStatusMLE: util.ESCS_COL_YELLOW_B,
		JudgeStatusOLE: util.ESCS_COL_YELLOW_B,
		JudgeStatusIE:  util.ESCS_COL_PURPLE_B,
	}
	col, ok := colorMap[js]
	if !ok {
		return util.ESCS_COL_PURPLE_B
	}
	return col
}

// String ... エスケープシーケンス付きで文字列化(出力用)
func (js JudgeStatus) String() string {
	return js.GetColorESCS() + js.ToString() + util.ESCS_COL_OFF
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
