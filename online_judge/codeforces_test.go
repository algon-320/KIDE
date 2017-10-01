package online_judge

import (
	"fmt"
	"testing"
)

func TestCodeforcesLogin(t *testing.T) {
	fmt.Println("testing : codeforces.go > Codeforces.login")
	br, err := Codeforces.login()
	if err != nil {
		t.Error("Codeforcesのログイン処理でエラーが発生しました。")
		t.Error(err)
		return
	}
	br.Open("http://codeforces.com/settings/general")
	if br.Title() != "Settings - Codeforces" {
		t.Error("Codeforcesにログイン出来ていません！ : %s", br.Title())
		return
	}
}

func TestCodeforcesNewProblem(t *testing.T) {
	fmt.Println("testing : codeforces.go > Codeforces.NewProblem")
	test := func(task, url string) {
		if err := Codeforces.NewProblem(url); err != nil {
			t.Error(err)
		}

		p, err := LoadProblem(task)
		if err != nil {
			t.Error(err)
			return
		}

		if p.ID != "A" {
			t.Error("id検出エラー")
		}
		if p.Name != "837_A" {
			t.Error("問題名検出エラー")
		}
		if p.ContestID != "837" {
			t.Error("コンテストid検出エラー")
		}
		if p.URL != url {
			t.Error("URLエラー")
		}
		if p.Oj.Name() != Codeforces.Name() {
			t.Error("OJエラー")
		}
		if len(p.Cases) != 3 {
			t.Error("サンプルケースを抽出出来ていません！")
		}
	}
	test("A", "http://codeforces.com/problemset/problem/837/A")
	test("A", "http://codeforces.com/contest/837/problem/A")
}

func TestCodeforcesIsValidURL(t *testing.T) {
	fmt.Println("testing : codeforces.go > Codeforces.IsValidURL")

	type res struct {
		isVaild      bool
		isProblemSet bool
	}

	testcase := map[string]res{
		"http://codeforces.com/contests":                        res{isVaild: false, isProblemSet: false},
		"http://codeforces.com/":                                res{isVaild: false, isProblemSet: false},
		"http://codeforces.com/problemset":                      res{isVaild: false, isProblemSet: false},
		"http://codeforces.com/problemset/problem/839/E":        res{isVaild: true, isProblemSet: false},
		"http://codeforces.com/contest/839":                     res{isVaild: true, isProblemSet: true},
		"http://codeforces.com/contest/839/problem/A":           res{isVaild: true, isProblemSet: false},
		"http://codeforces.com/gym/211512":                      res{isVaild: true, isProblemSet: true},
		"http://codeforces.com/gym/211512/problem/A":            res{isVaild: true, isProblemSet: false},
		"http://codeforces.com/group/ilLpsu4YlI/contest/214226": res{isVaild: true, isProblemSet: true},
	}

	for k, v := range testcase {
		p, s := Codeforces.IsValidURL(k)
		if p != v.isVaild || s != v.isProblemSet {
			t.Error(fmt.Sprintf("Codeforces.IsValidURL made incorrect judgement. %s", k))
		}
	}
}

// ！！！！実際に提出を行なうので注意！！！！
// func TestCodeforcecsSubmit(t *testing.T) {
// 	fmt.Println("testing : codeforces.go > Codeforces.Submit")

// 	err := Codeforces.NewProblem("http://codeforces.com/contest/837/problem/A")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	p, err := LoadProblem("A")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	test := func(expect JudgeStatus, lang language.Language, source string) {
// 		res, err := Codeforces.Submit(p, source, lang)
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
//     int ans = 0;
//     string s;
//     while (true) {
//         cin >> s;
//         if (cin.eof()) break;
//         int tmp = 0;
//         for (auto x: s) {
//             if('A' <= x && x <= 'Z') tmp++;
//         }
//         if (ans < tmp) ans = tmp;
//     }
//     cout << ans << endl;
//     return 0;
// }
// `)
// 	test(JudgeStatusWA, language.CPP, "// test submission  "+time.Now().String()+`
// #include <bits/stdc++.h>
// using namespace std;
// int main() {
//     int ans = 0;
//     string s;
//     while (true) {
//         cin >> s;
//         if (cin.eof()) break;
//         int tmp = 0;
//         for (auto x: s) {
//             if ('A' <= x && x <= 'Z') tmp++;
//         }
//         if (ans < tmp) ans = tmp;
//     }
//     cout << ans+1 << endl;
//     return 0;
// }
// `)
// }
