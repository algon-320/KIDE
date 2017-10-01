package online_judge

import (
	"fmt"
	"testing"
)

func TestFromProblemURL(t *testing.T) {
	fmt.Println("testing : online_judge.go > FromProblemURL")
	yc, err1 := FromProblemURL("https://yukicoder.me/problems/no/273")
	if err1 != nil {
		t.Error(err1)
		return
	}
	if yc.Name() != Yukicoder.Name() {
		t.Error("yukicoderの問題URLのパースに失敗")
		return
	}

	cf, err2 := FromProblemURL("http://codeforces.com/contest/835/problem/E")
	if err2 != nil {
		t.Error(err2)
		return
	}
	if cf.Name() != Codeforces.Name() {
		t.Error("Codeforcesの問題URLのパースに失敗")
		return
	}

	ac, err3 := FromProblemURL("https://beta.atcoder.jp/contests/arc079/tasks/arc079_b")
	if err3 != nil {
		t.Error(err3)
		return
	}
	if ac.Name() != ac.Name() {
		t.Error("AtCoderの問題URLのパースに失敗")
		return
	}
}
