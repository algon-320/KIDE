package online_judge

import (
	"fmt"
	"html"
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

type codeforces struct {
	name        string
	url         string
	loginURL    string
	sessionFile string
}

// Codeforces ... オンラインジャッジ: Codeforces
var Codeforces = &codeforces{
	name:        "Codeforces",
	url:         "http://codeforces.com/",
	loginURL:    "http://codeforces.com/enter",
	sessionFile: "session_codeforces.dat",
}

func (cf *codeforces) getLangID(lang language.Language) (string, error) {
	switch lang {
	case language.CPP:
		return "50", nil // 50 : GNU G++14 6.2.0
	case language.PYTHON:
		return "7", nil // 7 : Python 2.7.12
	default:
		return "", &ErrUnsuportedLanguage{name: lang.Name()}
	}
}

func (cf *codeforces) loadAccount() (string, string) {
	var handle string
	tmp, ok := setting.Get("OnlineJudge.Codeforces.Handle", "CODEFORCES_HANDLE")
	if !ok {
		handle = util.AskString("What is your Codeforces account id ?")
		setting.Set("OnlineJudge.Codeforces.Handle", handle)
	} else {
		handle = tmp.(string)
	}

	var password string
	tmp, ok = setting.Get("OnlineJudge.Codeforces.Password", "CODEFORCES_PASSWORD")
	if !ok {
		password = util.AskString("What is your Codeforces account password ?")
		setting.Set("OnlineJudge.Codeforces.Password", password)
	} else {
		password = tmp.(string)
	}

	return handle, password
}

func (cf *codeforces) login() (*browser.Browser, error) {
	handle, password := cf.loadAccount()
	data := map[string]string{"handle": handle, "password": password}

	br := surf.NewBrowser()

	cjar := util.LoadLoginSession(cf.sessionFile, cf.url)
	if cjar != nil {
		br.SetCookieJar(cjar)
		if cf.checkLoggedin(br) {
			fmt.Println(util.PrefixInfo + "Loaded session of Codeforces.")
			return br, nil
		}
	}

	// 新たにログイン
	fmt.Println(util.PrefixInfo + "Login to Codeforces...")

	if err := br.Open(cf.loginURL); err != nil {
		return nil, &ErrFailedToLogin{oj_name: cf.Name(), message: "Failed to open login page."}
	}

serch_form:
	for _, fm := range br.Forms() {
		if fm == nil {
			continue
		}

		for k, v := range data {
			if err := fm.Input(k, v); err != nil {
				continue serch_form
			}
		}
		if err := fm.Submit(); err != nil {
			return nil, &ErrFailedToLogin{oj_name: cf.Name(), message: "Failed to submit login information."}
		}

		if cf.checkLoggedin(br) {
			cookies := br.SiteCookies()
			util.SaveLoginSession(cf.sessionFile, cookies)
			return br, nil
		}
		return nil, &ErrFailedToLogin{oj_name: cf.Name(), message: "Incorrect username or password."}
	}

	return nil, &ErrFailedToLogin{oj_name: cf.Name(), message: "No form found."}
}

func (cf *codeforces) checkLoggedin(br *browser.Browser) bool {
	prevURL := br.Url()
	if prevURL != nil {
		defer br.Open(prevURL.String())
	}

	// 問題の提出ページを開いて開ければログインできている
	no1URL := "http://codeforces.com/contest/1/submit"
	br.Open(no1URL)
	return br.Url().String() == no1URL
}

func (cf *codeforces) Name() string {
	return cf.name
}

func (cf *codeforces) Submit(p *Problem, sourceCode string, lang language.Language) (*JudgeResult, error) {
	br, err := cf.login()
	if err != nil {
		return nil, err
	}

	submitURL := cf.url + fmt.Sprintf("contest/%s/submit", p.ContestID)
	if err := br.Open(submitURL); err != nil {
		return nil, err
	}

	langID, err := cf.getLangID(lang)
	if err != nil {
		return nil, err
	}

	for _, fm := range br.Forms() {
		if fm == nil {
			continue
		}

		if err := fm.Input("programTypeId", langID); err != nil {
			continue
		}
		if err := fm.Input("submittedProblemIndex", p.ID); err != nil {
			continue
		}
		if err := fm.Input("source", sourceCode); err != nil {
			continue
		}
		if err := fm.Submit(); err != nil {
			return nil, err
		}
		break
	}

	mysubmissionsURL := cf.url + fmt.Sprintf("contest/%s/my", p.ContestID)

	if br.Url().String() != mysubmissionsURL {
		return nil, &ErrFailedToSubmit{message: "might be the same solution."}
	}

	fmt.Println(util.PrefixInfo + "Your solution was successfully submitted.")

	var res JudgeResult
	res.Date = time.Now()
	res.Problem = p
	res.Code = sourceCode
	res.Language = lang
	res.Status = JudgeStatusUNK
	res.URL, _ = br.Dom().Find(".status-frame-datatable").Find("tr:nth-of-type(2) > td:nth-of-type(1) > a").Attr("href")
	res.URL, _ = br.ResolveStringUrl(res.URL)

	// get Judge Status
	var status string
	for {
		br.Open(mysubmissionsURL)
		verdict := br.Dom().Find(".status-frame-datatable").Find("tr:nth-of-type(2) > td:nth-of-type(6)")

		status = verdict.Find("span.submissionVerdictWrapper").First().Text()
		v, _ := verdict.Attr("waiting")
		if v == "false" {
			break
		}
		fmt.Println(status)
		time.Sleep(CheckInterval)
	}

	switch {
	case strings.HasPrefix(status, "Accepted"):
		res.Status = JudgeStatusAC
	case strings.HasPrefix(status, "Pretests passed"):
		res.Status = JudgeStatusPP
	case strings.HasPrefix(status, "Wrong answer"):
		res.Status = JudgeStatusWA
	case strings.HasPrefix(status, "Compilation error"):
		res.Status = JudgeStatusCE
	case strings.HasPrefix(status, "Runtime error"):
		res.Status = JudgeStatusRE
	case strings.HasPrefix(status, "Time limit exceeded"):
		res.Status = JudgeStatusTLE
	case strings.HasPrefix(status, "Memory limit exceeded"):
		res.Status = JudgeStatusMLE
	default:
		res.Status = JudgeStatusUNK
	}

	return &res, nil
}

func (cf *codeforces) NewProblem(url string) error {
	isValid, isSet := cf.IsValidURL(url)
	if !isValid {
		return &ErrInvalidProblemURL{url: url}
	}

	br, err := cf.login()
	if err != nil {
		return err
	}

	downloadProblem := func(problemURL string) error {
		var p Problem
		p.Oj = Codeforces
		p.URL = problemURL

		re1 := regexp.MustCompile(cf.url + "contest/(.+)/problem/(.+)")
		re2 := regexp.MustCompile(cf.url + "problemset/problem/(.+)/(.+)")
		group := re1.FindSubmatch([]byte(problemURL))
		if group == nil {
			group = re2.FindSubmatch([]byte(problemURL))
			if group == nil {
				return &ErrInvalidProblemURL{url: problemURL}
			}
		}

		p.ContestID = string(group[1]) // contest no.
		p.ID = string(group[2])        // A, B, C, and so on.
		p.Name = string(group[1]) + "_" + string(group[2])

		br.Open(problemURL)
		doc := br.Dom()

		var testCase TestCase
		doc.Find("div.sample-test > div").Each(func(_ int, s *goquery.Selection) {
			if s.HasClass("input") {
				pre, _ := goquery.OuterHtml(s.Find("pre"))
				pre = strings.Replace(pre, "<br/>", "\n", -1)
				testCase.Input = pre[5 : len(pre)-6] // <pre>と</pre>を取り除く
				testCase.Input = html.UnescapeString(testCase.Input)
				testCase.Input = util.AddBR(testCase.Input)

			} else if s.HasClass("output") {
				pre, _ := goquery.OuterHtml(s.Find("pre"))
				pre = strings.Replace(pre, "<br/>", "\n", -1)
				testCase.Output = pre[5 : len(pre)-6] // <pre>と</pre>を取り除く
				testCase.Output = html.UnescapeString(testCase.Output)
				testCase.Output = util.AddBR(testCase.Output)
				p.Cases = append(p.Cases, testCase)
			}
		})

		p.Print()
		return p.Save()
	}

	if isSet {
		util.DebugPrint("download [codeforces] problem set")

		br.Open(url)
		doc := br.Dom()
		doc.Find("table.problems").Find("tr").Each(func(i int, tr *goquery.Selection) {
			if i == 0 {
				return
			}
			problemURL, _ := tr.Find("td:first-of-type > a").Attr("href")
			problemURL, _ = br.ResolveStringUrl(problemURL)
			// fmt.Println(problemURL)
			err := downloadProblem(problemURL)
			if err != nil {
				fmt.Println(util.PrefixError + fmt.Sprintf("%s", err))
			}
		})
		return nil
	}
	return downloadProblem(url)
}

func (cf *codeforces) IsValidURL(url string) (bool, bool) {
	urlBytes := []byte(url)
	if regexp.MustCompile(cf.url+"contest/[0-9]+/problem/.+").Match(urlBytes) ||
		regexp.MustCompile(cf.url+"problemset/problem/[0-9]+/.+").Match(urlBytes) ||
		regexp.MustCompile(cf.url+"gym/[0-9]+/problem/.+").Match(urlBytes) ||
		regexp.MustCompile(cf.url+"problemset/problem/[0-9]+/.+").Match(urlBytes) ||
		regexp.MustCompile(cf.url+"group/.+/contest/[0-9]+/problem/.+").Match(urlBytes) {
		return true, false
	} else if regexp.MustCompile(cf.url+"contest/[0-9]+").Match(urlBytes) ||
		regexp.MustCompile(cf.url+"gym/[0-9]+").Match(urlBytes) ||
		regexp.MustCompile(cf.url+"group/.+/contest/[0-9]+").Match(urlBytes) {
		return true, true
	} else {
		return false, false
	}
}

func (cf *codeforces) MarshalJSON() ([]byte, error) {
	return []byte(`"` + cf.Name() + `"`), nil
}
