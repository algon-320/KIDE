package online_judge

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/util"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
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
	loginURL:    "http://judge.u-aizu.ac.jp/onlinejudge/signin.jsp",
	sessionFile: "session_aoj.dat",
}

func (a *aoj) getLangID(lang language.Language) (string, error) {
	switch lang {
	case language.CPP:
		return "C++14", nil // C++14
	case language.PYTHON2:
		return "Python", nil // Python2
	case language.PYTHON3:
		return "Python3", nil // Python3
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

func (a *aoj) login() (*browser.Browser, error) {
	handle, password := a.loadAccount()
	data := map[string]string{"userID": handle, "password": password}

	br := surf.NewBrowser()

	cjar := util.LoadLoginSession(a.sessionFile, a.url)
	if cjar != nil {
		br.SetCookieJar(cjar)
		if a.checkLoggedin(br) {
			fmt.Fprintln(os.Stderr, util.PrefixInfo+"Loaded session of AOJ.")
			return br, nil
		}
	}

	// 新たにログイン
	fmt.Fprintln(os.Stderr, util.PrefixInfo+"login to AOJ ...")

	type Payload struct {
		UserID   string `json:"id"`
		Password string `json:"password"`
	}
	payload := Payload{
		UserID:   data["userID"],
		Password: data["password"],
	}
	bytes, _ := json.Marshal(&payload)

	req, err := http.NewRequest("POST", "https://judgeapi.u-aizu.ac.jp/session", strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/JSON")

	client := &http.Client{Jar: br.CookieJar()}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if a.checkLoggedin(br) {
		cookieURL, _ := url.Parse("judge.u-aizu.ac.jp")
		cookies := br.CookieJar().Cookies(cookieURL)
		util.SaveLoginSession(a.sessionFile, cookies)
		return br, nil
	}

	return nil, &ErrFailedToLogin{oj_name: a.Name(), message: "Failed to login."}
}

func (a *aoj) checkLoggedin(br *browser.Browser) bool {
	br.Open("https://judgeapi.u-aizu.ac.jp/self")
	if br.StatusCode() != 200 {
		return false
	}
	return true
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
	postURL := "https://judgeapi.u-aizu.ac.jp/submissions"

	problemPath := problemID
	if len(lessonID) > 0 {
		problemPath = lessonID + "_" + problemID
	}

	type Payload struct {
		ProblemID  string `json:"problemId"`
		Language   string `json:"language"`
		SourceCode string `json:"sourceCode"`
	}
	payload := Payload{
		ProblemID:  problemPath,
		Language:   langID,
		SourceCode: sourceCode,
	}
	bytes, _ := json.Marshal(&payload)

	br, err := a.login()
	if err != nil {
		return nil, &ErrFailedToSubmit{message: err.(*ErrFailedToLogin).message}
	}

	// submit
	{
		req, err := http.NewRequest("POST", postURL, strings.NewReader(string(bytes)))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/JSON")
		client := &http.Client{Jar: br.CookieJar()}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, &ErrFailedToSubmit{message: "response status :" + fmt.Sprint(resp.StatusCode)}
		}
	}

	result := struct {
		Status []struct {
			RunID          string `xml:"run_id"`
			UserID         string `xml:"user_id"`
			ProblemID      string `xml:"problem_id"`
			SubmissionDate string `xml:"submission_date"`
			Status         string `xml:"status"`
			Language       string `xml:"language"`
			Cputime        string `xml:"cputime"`
			Memory         string `xml:"memory"`
			CodeSize       string `xml:"code_size"`
		} `xml:"status"`
	}{}
	handle, _ := a.loadAccount()
	// 結果を取る
	{
		time.Sleep(1 * time.Second)
		resp, err := http.Get("http://judge.u-aizu.ac.jp/onlinejudge/webservice/status_log?user_id=" + handle + "&limit=1")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		xmldata, _ := ioutil.ReadAll(resp.Body)
		if err := xml.Unmarshal(xmldata, &result); err != nil {
			return nil, err
		}
	}

	submissionID := strings.Trim(result.Status[0].RunID, "\n")
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
	for {
		status := strings.Trim(result.Status[0].Status, "\n")
		if strings.HasSuffix(status, "Accepted") {
			judgeRes.Status = JudgeStatusAC
			waiting = false
		} else if strings.HasSuffix(status, "Wrong Answer") {
			judgeRes.Status = JudgeStatusWA
			waiting = false
		} else if strings.HasSuffix(status, "Compile Error") {
			judgeRes.Status = JudgeStatusCE
			waiting = false
		} else if strings.HasSuffix(status, "Runtime Error") {
			judgeRes.Status = JudgeStatusRE
			waiting = false
		} else if strings.HasSuffix(status, "Time Limit Exceeded") {
			judgeRes.Status = JudgeStatusTLE
			waiting = false
		} else if strings.HasSuffix(status, "Memory Limit Exceeded") {
			judgeRes.Status = JudgeStatusMLE
			waiting = false
		} else if strings.HasSuffix(status, "Output Limit Exceeded") {
			judgeRes.Status = JudgeStatusOLE
			waiting = false
		} else if strings.HasSuffix(status, "-") {
			judgeRes.Status = JudgeStatusUNK
		}

		if !waiting {
			break
		}

		util.SaveCursorPos()
		{
			fmt.Fprint(os.Stderr, util.ESCS_COL_REVERSE+
				"waiting for judge "+strings.Repeat(".", watingCnt)+
				util.ESCS_COL_OFF)
			watingCnt++
			time.Sleep(CheckInterval)

			util.ClearCurrentLine()
		}
		util.RestoreCursorPos()
	}
	fmt.Fprint(os.Stderr, "\n")

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
			fmt.Fprintln(os.Stderr, group)
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

	return downloadProblem(url)
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
