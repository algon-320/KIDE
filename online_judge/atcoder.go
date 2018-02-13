package online_judge

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/util"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
)

type atcoder struct {
	name        string
	url         string
	loginURL    string
	sessionFile string
}

// AtCoder ... オンラインジャッジ: AtCoder
var AtCoder = &atcoder{
	name:        "AtCoder",
	url:         "https://beta.atcoder.jp/",
	loginURL:    "https://beta.atcoder.jp/login",
	sessionFile: "session_atcoder.dat",
}

func (ac *atcoder) getLangID(lang language.Language) (string, error) {
	switch lang {
	case language.CPP:
		return "3003", nil // C++14 (GCC 5.4.1)
	case language.PYTHON2:
		return "3022", nil // Python2 (2.7.6)
	case language.PYTHON3:
		return "3023", nil // Python3 (3.4.3)
	case language.JAVA:
		return "3016", nil // Java8 (OpenJDK 1.8.0)
	default:
		return "", &ErrUnsuportedLanguage{name: lang.Name()}
	}
}

func (ac *atcoder) loadAccount() (string, string) {
	var handle string
	if tmp, ok := setting.Get("OnlineJudge.AtCoder.Handle", "ATCODER_HANDLE"); ok {
		handle = tmp.(string)
	} else {
		handle = util.AskString("What is your AtCoder account id ?")
		setting.Set("OnlineJudge.AtCoder.Handle", handle)
	}

	var password string
	if tmp, ok := setting.Get("OnlineJudge.AtCoder.Password", "ATCODER_PASSWORD"); ok {
		password = tmp.(string)
	} else {
		password = util.AskString("What is your AtCoder account password ?")
		setting.Set("OnlineJudge.AtCoder.Password", password)
	}

	return handle, password
}

func (ac *atcoder) login() (*browser.Browser, error) {
	handle, password := ac.loadAccount()
	data := map[string]string{"username": handle, "password": password}

	br := surf.NewBrowser()

	cjar := util.LoadLoginSession(ac.sessionFile, ac.url)
	if cjar != nil {
		br.SetCookieJar(cjar)
		if ac.checkLoggedin(br) {
			fmt.Fprintln(os.Stderr, util.PrefixInfo+"Loaded session of AtCoder.")
			return br, nil
		}
	}

	// 新たにログイン
	fmt.Fprintln(os.Stderr, util.PrefixInfo+"login to AtCoder ...")

	if err := br.Open(ac.loginURL); err != nil {
		return nil, &ErrFailedToLogin{oj_name: ac.Name(), message: "Failed to open login page."}
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
			return nil, &ErrFailedToLogin{oj_name: ac.Name(), message: "Failed to submit login information."}
		}

		if ac.checkLoggedin(br) {
			cookies := br.SiteCookies()
			util.SaveLoginSession(ac.sessionFile, cookies)
			return br, nil
		}
		return nil, &ErrFailedToLogin{oj_name: ac.Name(), message: "Incorrect username or password."}
	}

	return nil, &ErrFailedToLogin{oj_name: ac.Name(), message: "No form found."}
}

func (ac *atcoder) checkLoggedin(br *browser.Browser) bool {
	prevURL := br.Url()
	if prevURL != nil {
		defer br.Open(prevURL.String())
	}

	// AGC001の提出ページにアクセスして開くことが出来ればログイン出来ている
	// (ログインしていないとログインページに飛ばされる)
	agc001URL := "https://beta.atcoder.jp/contests/agc001/submit"
	br.Open(agc001URL)
	return br.Url().String() == agc001URL
}

func (ac *atcoder) Name() string {
	return ac.name
}

func (ac *atcoder) Submit(p *Problem, sourceCode string, lang language.Language) (*JudgeResult, error) {
	br, err := ac.login()
	if err != nil {
		return nil, err
	}

	submitURL := ac.url + fmt.Sprintf("contests/%s/submit", p.ContestID)
	if err := br.Open(submitURL); err != nil {
		return nil, err
	}

	langID, err := ac.getLangID(lang)
	if err != nil {
		return nil, err
	}

	for _, fm := range br.Forms() {
		if fm == nil {
			continue
		}

		if err := fm.Input("sourceCode", sourceCode); err != nil {
			continue
		}

		qs := fm.Dom()
		qs.Find(fmt.Sprintf("#select-lang-%s > select", p.Name)).SetAttr("name", "data.LanguageId")
		fm = browser.NewForm(br, qs)

		if err := fm.Input("sourceCode", sourceCode); err != nil {
			continue
		}
		if err := fm.SelectByOptionValue("data.LanguageId", langID); err != nil {
			continue
		}
		if err := fm.SelectByOptionValue("data.TaskScreenName", p.Name); err != nil {
			continue
		}

		if err := fm.Submit(); err != nil {
			return nil, err
		}

		break
	}

	mysubmissionsURL := ac.url + fmt.Sprintf("contests/%s/submissions/me", p.ContestID)

	if br.Url().String() != mysubmissionsURL {
		return nil, &ErrFailedToSubmit{message: "might be empty."}
	}

	fmt.Fprintln(os.Stderr, util.PrefixInfo+"Your solution was successfully submitted.")

	var res JudgeResult
	res.Date = time.Now()
	res.Problem = p
	res.Code = sourceCode
	res.Language = lang
	res.Status = JudgeStatusUNK
	res.URL, _ = br.Dom().Find("tbody > tr:nth-of-type(1) > td:last-of-type").Find("a").Attr("href")
	res.URL, _ = br.ResolveStringUrl(res.URL)

	// get Judge Status
	var status string
	watingCnt := 0
waiting:
	for {
		br.Open(mysubmissionsURL)
		status = br.Dom().Find("tbody > tr:nth-of-type(1) > td:nth-of-type(7) > span").Text()

		switch {
		case strings.Contains(status, "AC"):
			res.Status = JudgeStatusAC
			break waiting
		case strings.Contains(status, "WA"):
			res.Status = JudgeStatusWA
			break waiting
		case strings.Contains(status, "CE"):
			res.Status = JudgeStatusCE
			break waiting
		case strings.Contains(status, "RE"):
			res.Status = JudgeStatusRE
			break waiting
		case strings.Contains(status, "TLE"):
			res.Status = JudgeStatusTLE
			break waiting
		case strings.Contains(status, "MLE"):
			res.Status = JudgeStatusMLE
			break waiting
		case strings.Contains(status, "OLE"):
			res.Status = JudgeStatusOLE
			break waiting
		case strings.Contains(status, "IE"):
			res.Status = JudgeStatusIE
			break waiting
		}

		if watingCnt == 0 {
			fmt.Fprint(os.Stderr, util.PrefixInfo+"waiting for judge .")
		} else {
			fmt.Fprint(os.Stderr, ".")
		}
		watingCnt++
		time.Sleep(CheckInterval)
	}
	fmt.Fprint(os.Stderr, "\n")

	if res.Status != JudgeStatusUNK {
		return &res, nil
	}
	return nil, &ErrFailedToSubmit{message: "no submit form found."}
}

func (ac *atcoder) NewProblem(url string) error {
	isValid, isSet := ac.IsValidURL(url)
	if !isValid {
		return &ErrInvalidProblemURL{url: url}
	}

	br, err := ac.login()
	if err != nil {
		return err
	}

	downloadProblem := func(problem_url string) error {
		var p Problem
		p.Oj = AtCoder
		p.URL = problem_url

		re := regexp.MustCompile(ac.url + "contests/(.+)/tasks/(.+)")
		group := re.FindSubmatch([]byte(problem_url))
		if group == nil {
			return &ErrInvalidProblemURL{url: problem_url}
		}
		p.ContestID = string(group[1])
		p.Name = string(group[2])

		br.Open(problem_url)
		doc := br.Dom()
		title := doc.Find("#main-container > div > div:nth-child(2) > span").Text()
		p.ID = title[0:1]

		var testCase TestCase
		japanese := false
		doc.Find("div.part").Each(func(_ int, s *goquery.Selection) {
			h3Text := s.Find("h3").Text()
			switch {
			case strings.HasPrefix(h3Text, "入力例"):
				japanese = true
				testCase.Input = s.Find("pre").Text()
				testCase.Input = html.UnescapeString(testCase.Input)
				testCase.Input = util.AddBR(testCase.Input)
			case strings.HasPrefix(h3Text, "出力例"):
				testCase.Output = s.Find("pre").Text()
				testCase.Output = html.UnescapeString(testCase.Output)
				testCase.Output = util.AddBR(testCase.Output)
				p.Cases = append(p.Cases, testCase)

			case strings.HasPrefix(h3Text, "Sample Input") && !japanese:
				testCase.Input = s.Find("pre").Text()
				testCase.Input = html.UnescapeString(testCase.Input)
				testCase.Input = util.AddBR(testCase.Input)
			case strings.HasPrefix(h3Text, "Sample Output") && !japanese:
				testCase.Output = s.Find("pre").Text()
				testCase.Output = html.UnescapeString(testCase.Output)
				testCase.Output = util.AddBR(testCase.Output)
				p.Cases = append(p.Cases, testCase)
			}
		})

		p.Print()
		return p.Save()
	}

	if isSet {
		util.DebugPrint("download [atcoder] problem set")

		br.Open(url)
		doc := br.Dom()
		doc.Find("tbody > tr").Each(func(_ int, tr *goquery.Selection) {
			problemURL, _ := tr.Find("td:first-of-type > a").Attr("href")
			problemURL, _ = br.ResolveStringUrl(problemURL)
			// fmt.Fprintln(os.Stderr, problemURL)
			err := downloadProblem(problemURL)
			if err != nil {
				fmt.Fprintln(os.Stderr, util.PrefixError+fmt.Sprintf("%s", err))
			}
		})
		return nil
	}
	return downloadProblem(url)
}

func (ac *atcoder) IsValidURL(url string) (bool, bool) {
	urlBytes := []byte(url)
	if regexp.MustCompile(ac.url + "contests/.+/tasks/.+").Match(urlBytes) {
		return true, false
	} else if regexp.MustCompile(ac.url + "contests/.+/tasks$").Match(urlBytes) {
		return true, true
	} else {
		return false, false
	}
}

func (ac *atcoder) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ac.Name() + `"`), nil
}
