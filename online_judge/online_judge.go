package online_judge

import (
	"net/url"

	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/util"
)

// OnlineJudge ... オンラインジャッジ定義用のインターフェース
type OnlineJudge interface {
	Name() string
	Submit(*Problem, string, language.Language) (*JudgeResult, error)
	NewProblem(string) error
	IsValidURL(string) (bool, bool) // isValid, isProblemSet
	MarshalJSON() ([]byte, error)
}

func FromName(ojName string) (OnlineJudge, error) {
	switch ojName {
	case AtCoder.Name():
		return AtCoder, nil
	case Codeforces.Name():
		return Codeforces, nil
	case Yukicoder.Name():
		return Yukicoder, nil
	case AOJ.Name():
		return AOJ, nil
	// ここに追加
	default:
		return nil, &ErrNoSuchOnlineJudge{oj_name: ojName}
	}
}

func FromProblemURL(problemURL string) (OnlineJudge, error) {
	if p, _ := AtCoder.IsValidURL(problemURL); p {
		util.DebugPrint(problemURL + " is a url of " + AtCoder.Name())
		return AtCoder, nil
	}
	if p, _ := Codeforces.IsValidURL(problemURL); p {
		util.DebugPrint(problemURL + " is a url of " + Codeforces.Name())
		return Codeforces, nil
	}
	if p, _ := Yukicoder.IsValidURL(problemURL); p {
		util.DebugPrint(problemURL + " is a url of " + Yukicoder.Name())
		return Yukicoder, nil
	}
	if p, _ := AOJ.IsValidURL(problemURL); p {
		util.DebugPrint(problemURL + " is a url of " + AOJ.Name())
		return AOJ, nil
	}
	// ここに追加

	urlObj, _ := url.Parse(problemURL)
	return nil, &ErrNoSuchOnlineJudge{oj_name: urlObj.Hostname()}
}
