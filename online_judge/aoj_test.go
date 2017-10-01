package online_judge

import (
	"fmt"
	"testing"
)

func TestAojNewProblem(t *testing.T) {
	fmt.Println("testing : aoj.go > AOJ.NewProblem")

	type tmp struct {
		url      string
		id       string
		name     string
		oj       string
		num_case int
	}

	test := func(testcase *tmp) {
		if err := AOJ.NewProblem(testcase.url); err != nil {
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
			url:      "http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ITP1_1_A&lang=jp",
			id:       "ITP1_1_A",
			name:     "ITP1_1_A",
			oj:       AOJ.Name(),
			num_case: 1,
		},
	}

	for _, t := range testcases {
		test(t)
	}
}

func TestAojIsValidURL(t *testing.T) {
	fmt.Println("testing : aoj.go > AOJ.IsValidURL")

	type res struct {
		isVaild      bool
		isProblemSet bool
	}

	testcase := map[string]res{
		"http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ITP1_1_A&lang=jp": res{isVaild: true, isProblemSet: false},
		"http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=0000":             res{isVaild: true, isProblemSet: false},
	}

	for k, v := range testcase {
		p, s := AOJ.IsValidURL(k)
		if p != v.isVaild || s != v.isProblemSet {
			t.Error("AOJ.IsValidURL made incorrect judgement.")
		}
	}
}

// // ！！！！実際に提出を行なうので注意！！！！
// func TestAojSubmit(t *testing.T) {
// 	fmt.Println("testing : aoj.go > AOJ.Submit")
// 	err := AOJ.NewProblem("http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ITP1_1_A&lang=jp")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	p, err := LoadProblem("ITP1_1_A")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	test := func(expect JudgeStatus, lang language.Language, source string) {
// 		res, err := AOJ.Submit(p, source, lang)
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
//   cout << "Hello World" << endl;
//   return 0;
// }
// `)

// 	test(JudgeStatusWA, language.CPP, "// test submission  "+time.Now().String()+`
// #include <bits/stdc++.h>
// using namespace std;
// int main() {
//   cout << "Hello,World" << endl;
//   return 0;
// }
// `)
// }
