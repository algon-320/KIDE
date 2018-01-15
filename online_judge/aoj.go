package online_judge

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/util"
)

type aoj struct {
	name        string
	url         string
	loginURL    string
	sessionFile string
}

// AOJ ... オンラインジャッジ: Aizu Online Judge
var AOJ = &aoj{
	name:        "Aizu Online Judge",
	url:         "http://judge.u-aizu.ac.jp/onlinejudge/index.jsp",
	loginURL:    "http://judge.u-aizu.ac.jp/onlinejudge/index.jsp",
	sessionFile: "session_aoj.dat",
}

func (a *aoj) getLangID(lang language.Language) (string, error) {
	switch lang {
	case language.CPP:
		return "C++14", nil // C++14
	case language.PYTHON:
		return "Python", nil // Python2
	case language.JAVA:
		return "JAVA", nil // JAVA
	default:
		return "", &ErrUnsuportedLanguage{name: lang.Name()}
	}
}

func (a *aoj) loadAccount() (string, string) {
	var handle string
	if tmp, ok := setting.Get("OnlineJudge.AOJ.Handle", "AOJ_HANDLE"); ok {
		handle = tmp.(string)
	} else {
		handle = util.AskString("What is your AOJ account id ?")
		setting.Set("OnlineJudge.AOJ.Handle", handle)
	}

	var password string
	if tmp, ok := setting.Get("OnlineJudge.AOJ.Password", "AOJ_PASSWORD"); ok {
		password = tmp.(string)
	} else {
		password = util.AskString("What is your AOJ account password ?")
		setting.Set("OnlineJudge.AOJ.Password", password)
	}

	return handle, password
}

func (a *aoj) extractIDs(submitURL string) (string, string) {
	urlObj, _ := url.Parse(submitURL)
	fragment := urlObj.Fragment
	hv := strings.Split(fragment, "/")
	if hv[0] == "submit" {
		if len(hv) == 2 {
			return "", hv[1]
		}
		return hv[1], hv[2]
	}
	return "", ""
}

func (a *aoj) Name() string {
	return a.name
}

func (a *aoj) Submit(p *Problem, sourceCode string, lang language.Language) (*JudgeResult, error) {
	langID, err := a.getLangID(lang)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocument(p.URL)
	if err != nil {
		return nil, err
	}

	submitURL, ok := doc.Find("#pageinfo > div > div > a:last-of-type").Attr("href")
	if !ok {
		return nil, &ErrFailedToSubmit{message: "cannot detect submit link."}
	}
	submitURL = "http://judge.u-aizu.ac.jp/onlinejudge/" + submitURL

	lessonID, problemID := a.extractIDs(submitURL)
	postURL := "http://judge.u-aizu.ac.jp/onlinejudge/webservice/submit"

	handle, password := a.loadAccount()

	values := url.Values{}
	values.Set("userID", handle)
	values.Add("password", password)
	values.Add("problemNO", problemID)
	values.Add("lessonID", lessonID)
	values.Add("language", langID)
	values.Add("sourceCode", sourceCode)

	req, err := http.NewRequest("POST", postURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, &ErrFailedToSubmit{message: "response status :" + fmt.Sprint(resp.StatusCode)}
	}

	bbody, _ := ioutil.ReadAll(resp.Body)
	sbody := string(bbody)
	// 成功した時は"0\n"が返る
	if sbody != "0\n" {
		return nil, &ErrFailedToSubmit{message: "Incorrect Handle or Password."}
	}

	allResultsURL := "http://judge.u-aizu.ac.jp/onlinejudge/status.jsp"
	doc, err = goquery.NewDocument(allResultsURL)
	if err != nil {
		return nil, err
	}

	submissionID := ""
	doc.Find("#tableRanking").Find("tr:nth-of-type(n + 2)").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if s.Find("td:nth-of-type(2) > a").Text() == handle {
			submissionID = s.Find("td:nth-of-type(1) > a").Text()
			return false
		}
		return true
	})

	if submissionID == "" {
		return nil, &ErrFailedToSubmit{message: "no judge result"}
	}

	time.Sleep(time.Second) // 一秒間待つ

	resultURL := "http://judge.u-aizu.ac.jp/onlinejudge/review.jsp?rid=" + submissionID
	var judgeRes JudgeResult
	judgeRes.Date = time.Now()
	judgeRes.Problem = p
	judgeRes.Code = sourceCode
	judgeRes.Language = lang
	judgeRes.Status = JudgeStatusUNK
	judgeRes.URL = resultURL

	waiting := true
	watingCnt := 0
	for waiting {
		doc, err = goquery.NewDocument(resultURL)
		if err != nil {
			return nil, err
		}

		doc.Find("table:nth-of-type(3)").Find("tr:nth-of-type(n + 2)").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			status := s.Find("td:nth-of-type(2)").Text()
			status = strings.TrimRight(status, "\n")

			if strings.HasSuffix(status, "Accepted") {
				judgeRes.Status = JudgeStatusAC
				return true
			} else if strings.HasSuffix(status, "Wrong Answer") {
				judgeRes.Status = JudgeStatusWA
				waiting = false
				return false
			} else if strings.HasSuffix(status, "Compile Error") {
				judgeRes.Status = JudgeStatusCE
				waiting = false
				return false
			} else if strings.HasSuffix(status, "Runtime Error") {
				judgeRes.Status = JudgeStatusRE
				waiting = false
				return false
			} else if strings.HasSuffix(status, "Time Limit Exceeded") {
				judgeRes.Status = JudgeStatusTLE
				waiting = false
				return false
			} else if strings.HasSuffix(status, "Memory Limit Exceeded") {
				judgeRes.Status = JudgeStatusMLE
				waiting = false
				return false
			} else if strings.HasSuffix(status, "Output Limit Exceeded") {
				judgeRes.Status = JudgeStatusOLE
				waiting = false
				return false
			} else if strings.HasSuffix(status, "-") {
				judgeRes.Status = JudgeStatusUNK
				return false
			}
			return true
		})

		if judgeRes.Status == JudgeStatusAC {
			break
		}

		if watingCnt == 0 {
			fmt.Print(util.PrefixInfo + "waiting for judge .")
		} else {
			fmt.Print(".")
		}
		watingCnt++
		time.Sleep(CheckInterval)
	}
	fmt.Print("\n")

	return &judgeRes, nil
}

func (a *aoj) NewProblem(url string) error {

	// ---------------- メモ ----------------
	// 入力例
	// 出力例
	// サンプル入力
	// サンプル出力
	// Sample Input
	// Sample Output
	// Output for the Sample Input
	// Output for Sample Input
	// --------------------------------------

	isValid, _ := a.IsValidURL(url)
	if !isValid {
		return &ErrInvalidProblemURL{url: url}
	}

	downloadProblem := func(problemURL string) error {
		var p Problem
		p.Oj = AOJ
		p.URL = problemURL

		re := regexp.MustCompile("http://judge.u-aizu.ac.jp/onlinejudge/description.jsp\\?id=(.+?)(?:&.*)?$")
		group := re.FindSubmatch([]byte(problemURL))
		if group == nil {
			fmt.Println(group)
			return &ErrInvalidProblemURL{url: problemURL}
		}
		p.ID = string(group[1])
		p.Name = string(group[1])
		p.ContestID = ""

		doc, err := goquery.NewDocument(problemURL)
		if err != nil {
			return err
		}

		var testCase TestCase
		doc.Find("h2,h3").Each(func(_ int, s *goquery.Selection) {
			utfText, _ := util.ShiftJIS2UTF8(s.Text())

			if strings.HasPrefix(utfText, "入力例") ||
				strings.HasPrefix(utfText, "サンプル入力") ||
				strings.HasPrefix(utfText, "Sample Input") {
				testCase.Input = s.Next().Text()
				testCase.Input = html.UnescapeString(testCase.Input)
				testCase.Input = util.AddBR(testCase.Input)

			} else if strings.HasPrefix(utfText, "出力例") ||
				strings.HasPrefix(utfText, "サンプル出力") ||
				strings.HasPrefix(utfText, "Sample Output") ||
				strings.HasPrefix(utfText, "Output for") {
				testCase.Output = s.Next().Text()
				testCase.Output = html.UnescapeString(testCase.Output)
				testCase.Output = util.AddBR(testCase.Output)
				p.Cases = append(p.Cases, testCase)
			}
		})

		p.Print()
		return p.Save()
	}

	if err := downloadProblem(url); err != nil {
		return err
	}
	return nil
}

func (a *aoj) IsValidURL(url string) (bool, bool) {
	if strings.HasPrefix(url, "http://judge.u-aizu.ac.jp/onlinejudge/description.jsp") {
		return true, false
	}
	return false, false
}

func (a *aoj) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.Name() + `"`), nil
}
