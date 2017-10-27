package online_judge

import (
	"fmt"
	"html"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/util"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

type yukicoder struct {
	name        string
	url         string
	loginURL    string
	sessionFile string
}

// Yukicoder ... オンラインジャッジ: yukicoder
var Yukicoder = &yukicoder{
	name:        "yukicoder",
	url:         "https://yukicoder.me/",
	loginURL:    "https://yukicoder.me/auth/twitter",
	sessionFile: "session_yukicoder.dat",
}

func (yc *yukicoder) getLangID(lang language.Language) (string, error) {
	switch lang {
	case language.CPP:
		return "cpp14", nil // C++14 (gcc 7.1.0)
	case language.PYTHON:
		return "python", nil // Python2 (2.7.13)
	default:
		return "", &ErrUnsuportedLanguage{name: lang.Name()}
	}
}

func (yc *yukicoder) loadAccount() (string, string) {
	var handle string
	tmp, ok := setting.Get("OnlineJudge.yukicoder.Handle", "YUKICODER_HANDLE")
	if !ok {
		handle = util.AskString("What is your yukicoder account id (twitter id) ?")
		setting.Set("OnlineJudge.yukicoder.Handle", handle)
	} else {
		handle = tmp.(string)
	}

	var password string
	tmp, ok = setting.Get("OnlineJudge.yukicoder.Password", "YUKICODER_PASSWORD")
	if !ok {
		password = util.AskString("What is your yukicoder account password (twitter password) ?")
		setting.Set("OnlineJudge.yukicoder.Password", password)
	} else {
		password = tmp.(string)
	}

	return handle, password
}

func (yc *yukicoder) login() (*browser.Browser, error) {
	br := surf.NewBrowser()

	cjar := util.LoadLoginSession(yc.sessionFile, yc.url)
	if cjar != nil {
		br.SetCookieJar(cjar)
		if yc.checkLoggedin(br) {
			fmt.Println(util.PrefixInfo + "Loaded session of yukicoder.")
			return br, nil
		}
	}

	// 新たにログイン
	fmt.Println(util.PrefixInfo + "Login to yukicoder...")

	if err := br.Open(yc.loginURL); err != nil {
		return nil, &ErrFailedToLogin{oj_name: yc.Name(), message: "Failed to open login page."}
	}
	if len(br.Forms()) < 2 {
		return nil, &ErrFailedToLogin{oj_name: yc.Name(), message: "Longin form not found. The login form of twitter might be changed."}
	}

	handle, password := yc.loadAccount()

	fm := br.Forms()[1]
	qs := fm.Dom()
	qs.Find("input[name='cancel']").First().Remove()
	fm = browser.NewForm(br, qs)
	if err := fm.Input("session[username_or_email]", handle); err != nil {
		return nil, &ErrFailedToLogin{oj_name: yc.Name(), message: "No input of username found."}
	}
	if err := fm.Input("session[password]", password); err != nil {
		return nil, &ErrFailedToLogin{oj_name: yc.Name(), message: "No input of password found."}
	}
	if err := fm.Submit(); err != nil {
		return nil, &ErrFailedToLogin{oj_name: yc.Name(), message: "Failed to submit login information."}
	}
	br.Click("a.maintain-context")

	if yc.checkLoggedin(br) {
		cookies := br.SiteCookies()
		util.SaveLoginSession(yc.sessionFile, cookies)
		return br, nil
	}
	return nil, &ErrFailedToLogin{oj_name: yc.Name(), message: "Incorrect username or password"}
}

func (yc *yukicoder) checkLoggedin(br *browser.Browser) bool {
	prevURL := br.Url()
	if prevURL != nil {
		defer br.Open(prevURL.String())
	}

	// 問題の提出ページを開いて開ければログインできている
	no1URL := "https://yukicoder.me/problems/no/1/submit"
	br.Open(no1URL)
	return br.Title() != "yukicoder"
}

func (yc *yukicoder) Name() string {
	return yc.name
}

func (yc *yukicoder) Submit(p *Problem, sourceCode string, lang language.Language) (*JudgeResult, error) {
	br, err := yc.login()
	if err != nil {
		return nil, err
	}

	submitURL := yc.url + fmt.Sprintf("problems/no/%s/submit", p.ID)
	if err := br.Open(submitURL); err != nil {
		return nil, err
	}

	langID, err := yc.getLangID(lang)
	if err != nil {
		return nil, err
	}

	for _, fm := range br.Forms() {
		if fm == nil {
			continue
		}

		if err := fm.Input("lang", langID); err != nil {
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

	if br.Url().String() == submitURL {
		return nil, &ErrFailedToSubmit{message: "might be empty."}
	}

	fmt.Println(util.PrefixInfo + "Your solution was successfully submitted.")

	mysubmissionURL := br.Url().String()

	var res JudgeResult
	res.Date = time.Now()
	res.Problem = p
	res.Code = sourceCode
	res.Language = lang
	res.Status = JudgeStatusUNK
	res.URL = mysubmissionURL

	// get Judge Status
	var status string
waiting:
	for {
		br.Open(mysubmissionURL)
		status = br.Dom().Find("#status").Text()

		switch status {
		case "AC":
			res.Status = JudgeStatusAC
			break waiting
		case "WA":
			res.Status = JudgeStatusWA
			break waiting
		case "CE":
			res.Status = JudgeStatusCE
			break waiting
		case "RE":
			res.Status = JudgeStatusRE
			break waiting
		case "TLE":
			res.Status = JudgeStatusTLE
			break waiting
		case "MLE":
			res.Status = JudgeStatusMLE
			break waiting
		case "OLE":
			res.Status = JudgeStatusOLE
			break waiting
		case "IE":
			res.Status = JudgeStatusIE
			break waiting
		}

		fmt.Println("waiting for judge ...")
		time.Sleep(CheckInterval)
	}

	return &res, nil
}

func (yc *yukicoder) NewProblem(url string) error {
	isValid, isSet := yc.IsValidURL(url)
	if !isValid {
		return &ErrInvalidProblemURL{url: url}
	}

	br, err := yc.login()
	if err != nil {
		return err
	}

	downloadProblem := func(problemURL string) error {
		var p Problem
		p.Oj = Yukicoder
		p.URL = problemURL

		re := regexp.MustCompile(yc.url + "problems/no/(.+)")
		group := re.FindSubmatch([]byte(problemURL))
		if group == nil {
			return &ErrInvalidProblemURL{url: problemURL}
		}

		p.ContestID = ""
		p.ID = string(group[1])
		p.Name = string(group[1])

		br.Open(problemURL)
		doc := br.Dom()

		var testCase TestCase
		doc.Find("div.sample > div").Each(func(_ int, s *goquery.Selection) {
			testCase.Input = s.Find("pre:nth-of-type(1)").Text()
			testCase.Input = html.UnescapeString(testCase.Input)
			testCase.Input = util.AddBR(testCase.Input)
			testCase.Output = s.Find("pre:nth-of-type(2)").Text()
			testCase.Output = html.UnescapeString(testCase.Output)
			testCase.Output = util.AddBR(testCase.Output)
			p.Cases = append(p.Cases, testCase)
		})

		p.Print()
		return p.Save()
	}

	if isSet {
		util.DebugPrint("download [yukicoder] problem set")

		br.Open(url)
		doc := br.Dom()
		doc.Find("tbody").Find("tr").Each(func(i int, tr *goquery.Selection) {
			problemURL, _ := tr.Find("td:nth-of-type(3) > a").Attr("href")
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

func (yc *yukicoder) IsValidURL(url string) (bool, bool) {
	urlBytes := []byte(url)
	if regexp.MustCompile(yc.url + "problems/no/[0-9]+").Match(urlBytes) {
		return true, false
	} else if regexp.MustCompile(yc.url + "contests/[0-9]+").Match(urlBytes) {
		return true, true
	} else {
		return false, false
	}
}

func (yc *yukicoder) MarshalJSON() ([]byte, error) {
	return []byte(`"` + yc.Name() + `"`), nil
}
