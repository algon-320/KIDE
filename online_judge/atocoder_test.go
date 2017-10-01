package online_judge

import (
	"fmt"
	"testing"
)

func TestAtCoderLogin(t *testing.T) {
	fmt.Println("testing : atcoder.go > AtCoder.login")
	br, err := AtCoder.login()
	if err != nil {
		t.Fatal("AtCoderのログイン処理でエラーが発生しました。", err)
		return
	}
	br.Open("https://beta.atcoder.jp/settings")
	if br.Title() != "General Settings - AtCoder" {
		t.Fatalf("AtCoderにログイン出来ていません！ : %s", br.Title())
		return
	}
}

func TestAtCoderNewProblem(t *testing.T) {
	fmt.Println("testing : atcoder.go > AtCoder.NewProblem")

	type tmp struct {
		url      string
		id       string
		name     string
		contest  string
		oj       string
		num_case int
	}

	test := func(testcase *tmp) {
		if err := AtCoder.NewProblem(testcase.url); err != nil {
			t.Error(err)
		}

		p, err := LoadProblem(testcase.id)
		if err != nil {
			t.Error(err)
			return
		}

		if p.ID != testcase.id {
			t.Error("id検出エラー")
		}
		if p.Name != testcase.name {
			t.Error("問題名検出エラー")
		}
		if p.ContestID != testcase.contest {
			t.Error("コンテストid検出エラー")
		}
		if p.URL != testcase.url {
			t.Error("URLエラー")
		}
		if p.Oj.Name() != testcase.oj {
			t.Error("OJエラー")
		}
		if len(p.Cases) != testcase.num_case {
			t.Error("サンプルケースを抽出出来ていません！")
		}
	}

	testcases := []*tmp{
		&tmp{
			url:      "https://beta.atcoder.jp/contests/chokudai_S001/tasks/chokudai_S001_a",
			id:       "A",
			name:     "chokudai_S001_a",
			contest:  "chokudai_S001",
			oj:       AtCoder.Name(),
			num_case: 4,
		},
		&tmp{
			url:      "https://beta.atcoder.jp/contests/arc002/tasks/arc002_1",
			id:       "A",
			name:     "arc002_1",
			contest:  "arc002",
			oj:       AtCoder.Name(),
			num_case: 4,
		},
		// TODO: ARC001では失敗する
		// &tmp{
		//   url:"https://beta.atcoder.jp/contests/arc001/tasks/arc001_1",
		//   id:"A",
		//   name:"arc001_1",
		//   contest:"arc001",
		//   oj:AtCoder.Name(),
		//   num_case:3,
		// },
	}

	for _, t := range testcases {
		test(t)
	}
}

func TestAtCoderIsValidURL(t *testing.T) {
	fmt.Println("testing : atcoder.go > AtCoder.IsValidURL")

	type res struct {
		isVaild      bool
		isProblemSet bool
	}

	testcase := map[string]res{
		"https://atcoder.jp/post/37":                             res{isVaild: false, isProblemSet: false},
		"https://beta.atcoder.jp/contests/abc070":                res{isVaild: false, isProblemSet: false},
		"https://beta.atcoder.jp/contests/abc070/tasks":          res{isVaild: true, isProblemSet: true},
		"https://beta.atcoder.jp/contests/abc070/tasks/abc070_a": res{isVaild: true, isProblemSet: false},
		"https://beta.atcoder.jp/contests/abc070/clarifications": res{isVaild: false, isProblemSet: false},
	}

	for k, v := range testcase {
		p, s := AtCoder.IsValidURL(k)
		if p != v.isVaild || s != v.isProblemSet {
			t.Error("AtCoder.IsValidURL made incorrect judgement.")
		}
	}
}

// ！！！！実際に提出を行なうので注意！！！！
// func TestAtCoderSubmit(t *testing.T) {
// 	fmt.Println("testing : atcoder.go > AtCoder.Submit")
// 	err := AtCoder.NewProblem("https://beta.atcoder.jp/contests/abc068/tasks/abc068_a")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	p, err := LoadProblem("a")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	test := func(expect JudgeStatus, lang language.Language, source string) {
// 		res, err := AtCoder.Submit(p, source, lang)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if res.Status != expect {
// 			t.Error("ジャッジ結果の取得に失敗しています。")
// 		}
// 	}

// 	test(JudgeStatusAC, language.CPP, "// test submission  "+time.Now().String()+`
// #include <bits/stdc++.h>
// using namespace std;
// int main() {
//   string s;
//   cin >> s;
//   cout << "ABC"+s << endl;
//   return 0;
// }
// `)

// 	test(JudgeStatusWA, language.CPP, "// test submission  "+time.Now().String()+`
// #include <bits/stdc++.h>
// using namespace std;
// int main() {
//   int s;
//   cin >> s;
//   cout << "ABC"+s << endl;
//   return 0;
// }
// `)
// }
