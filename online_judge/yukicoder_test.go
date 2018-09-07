package online_judge

import (
	"fmt"
	"testing"
)

func TestYukicoderLogin(t *testing.T) {
	fmt.Println("testing : yukicoder.go > Yukicoder.login")
	br, err := Yukicoder.login()
	if err != nil {
		t.Error("yukicoderのログイン処理でエラーが発生しました。")
		t.Error(err)
		return
	}
	br.Open("https://yukicoder.me/problems/no/1/submit")
	if br.Title() != "提出 - yukicoder" {
		t.Errorf("yukicoderにログイン出来ていません！ : %s", br.Title())
		return
	}
}

func TestYukicoderNewProblem(t *testing.T) {
	fmt.Println("testing : yukicoder.go > Yukicoder.NewProblem")
	if err := Yukicoder.NewProblem("https://yukicoder.me/problems/no/543"); err != nil {
		t.Error(err)
	}

	p, err := LoadProblem("543")
	if err != nil {
		t.Error(err)
		return
	}

	if p.ID != "543" {
		t.Error("id検出エラー")
	}
	if p.Name != "543" {
		t.Error("問題名検出エラー")
	}
	if p.URL != "https://yukicoder.me/problems/no/543" {
		t.Error("URLエラー")
	}
	if p.Oj.Name() != Yukicoder.Name() {
		t.Error("OJエラー")
	}
	if len(p.Cases) != 2 {
		t.Error("サンプルケースを抽出出来ていません！")
	}
}

func TestYukicoderIsValidURL(t *testing.T) {
	fmt.Println("testing : yukicoder.go > Yukicoder.IsValidURL")

	type res struct {
		isVaild      bool
		isProblemSet bool
	}

	testcase := map[string]res{
		"https://yukicoder.me/":                res{isVaild: false, isProblemSet: false},
		"https://yukicoder.me/problems/no/557": res{isVaild: true, isProblemSet: false},
		"https://yukicoder.me/contests/173":    res{isVaild: true, isProblemSet: true},
	}

	for k, v := range testcase {
		p, s := Yukicoder.IsValidURL(k)
		if p != v.isVaild || s != v.isProblemSet {
			t.Error("Yukicoder.IsValidURL made incorrect judgement.")
		}
	}
}

// ！！！！実際に提出を行なうので注意！！！！
// func TestYukicoderSubmit(t *testing.T) {
// 	fmt.Println("testing : yukicoder.go > Yukicoder.Submit")

// 	err := Yukicoder.NewProblem("https://yukicoder.me/problems/no/543")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	p, err := LoadProblem("543")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	test := func(expect JudgeStatus, source string) {
// 		res, err := Yukicoder.Submit(p, source, language.CPP)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if res.Status != expect {
// 			t.Error("ジャッジ結果の取得に失敗しています。")
// 		}
// 	}

// 	test(JudgeStatusAC, "// test submission  "+time.Now().String()+`
// #include <bits/stdc++.h>
// using namespace std;
// int main() {
//   string a,b;
//   cin>>a>>b;
//   cout<<b<<" "<<a<<endl;
//   return 0;
// }
// `)

// 	test(JudgeStatusWA, "// test submission  "+time.Now().String()+`
// #include <bits/stdc++.h>
// using namespace std;
// int main() {
//   string a,b;
//   cin>>a>>b;
//   cout<<a<<" "<<b<<endl;
//   return 0;
// }
// `)
// }
