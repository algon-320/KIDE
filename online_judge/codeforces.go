package online_judge

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/util"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
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
	case language.PYTHON2:
		return "7", nil // 7 : Python 2.7.12
	case language.PYTHON3:
		return "31", nil // 31 : Python 3.6
	case language.JAVA:
		return "36", nil // Java 1.8.0_131
	default:
		return "", &ErrUnsuportedLanguage{name: lang.Name()}
	}
}

func (cf *codeforces) loadAccount() (string, string) {
	var handle string
	if tmp, ok := setting.Get("OnlineJudge.Codeforces.Handle", "CODEFORCES_HANDLE"); ok {
		handle = tmp.(string)
	} else {
		handle = util.AskString("What is your Codeforces account id ?")
		setting.Set("OnlineJudge.Codeforces.Handle", handle)
	}

	var password string
	if tmp, ok := setting.Get("OnlineJudge.Codeforces.Password", "CODEFORCES_PASSWORD"); ok {
		password = tmp.(string)
	} else {
		password = util.AskString("What is your Codeforces account password ?")
		setting.Set("OnlineJudge.Codeforces.Password", password)
	}

	return handle, password
}

func (cf *codeforces) login() (*browser.Browser, error) {
	handle, password := cf.loadAccount()
	data := map[string]string{"handleOrEmail": handle, "password": password}

	br := surf.NewBrowser()

	cjar := util.LoadLoginSession(cf.sessionFile, cf.url)
	if cjar != nil {
		br.SetCookieJar(cjar)
		if cf.checkLoggedin(br) {
			fmt.Fprintln(os.Stderr, util.PrefixInfo+"Loaded session of Codeforces.")
			return br, nil
		}
	}

	// 新たにログイン
	fmt.Fprintln(os.Stderr, util.PrefixInfo+"Login to Codeforces...")

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

		if err := fm.SelectByOptionValue("programTypeId", langID); err != nil {
			continue
		}
		if err := fm.SelectByOptionValue("submittedProblemIndex", p.ID); err != nil {
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

	fmt.Fprintln(os.Stderr, util.PrefixInfo+"Your solution was successfully submitted.")

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

		util.SaveCursorPos()
		{
			fmt.Fprintln(os.Stderr, util.ESCS_COL_REVERSE+status+util.ESCS_COL_OFF)

			time.Sleep(CheckInterval)

			util.ClearCurrentLine()
		}
		util.RestoreCursorPos()
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
				fmt.Fprintln(os.Stderr, util.PrefixError+fmt.Sprintf("%s", err))
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

func (cf *codeforces) ShowMySubmissions(contestID int) {
	br, err := cf.login()
	if err != nil {
		return
	}

	mysubmissionsURL := cf.url + fmt.Sprintf("contest/%d/my", contestID)
	if err := br.Open(mysubmissionsURL); err != nil {
		return
	}

	judgeFinished := make(map[int]struct{})
	waitingTotal := -1

	messagePrinted := false
	clearLine := func() {
		if messagePrinted {
			// waiting for judge の行を消して、カーソル位置を戻す
			util.ClearCurrentLine()
			util.RestoreCursorPos()
			messagePrinted = false
		}
	}

	for {
		br.Open(mysubmissionsURL)
		trs := br.Dom().Find(".status-frame-datatable").Find("tr")

		submissionIDs := []int{}
		statuses := map[int]JudgeStatus{}
		problems := map[int]string{}
		results := map[int]string{}

		currentWaitingCount := 0
		trs.Each(func(i int, tr *goquery.Selection) {
			if cl, ok := tr.Attr("class"); ok && cl == "first-row" {
				return
			}

			var submissionID int
			if tmp, ok := tr.Attr("data-submission-id"); ok {
				submissionID, _ = strconv.Atoi(tmp)
			} else {
				return
			}
			problems[submissionID] = strings.TrimSpace(tr.Find("td:nth-of-type(4)").Text())

			isWaiting, _ := tr.Find("td.status-verdict-cell").First().Attr("waiting")
			if isWaiting == "true" {
				currentWaitingCount++
				statuses[submissionID] = JudgeStatusUNK
				return
			}

			submissionIDs = append(submissionIDs, submissionID)

			var status JudgeStatus
			span := tr.Find("span.submissionVerdictWrapper").First()
			result := strings.TrimSpace(span.Text())
			switch {
			case strings.HasPrefix(result, "Accepted"):
				status = JudgeStatusAC
			case strings.HasPrefix(result, "Pretests passed"):
				status = JudgeStatusPP
			case strings.HasPrefix(result, "Wrong answer"):
				status = JudgeStatusWA
			case strings.HasPrefix(result, "Compilation error"):
				status = JudgeStatusCE
			case strings.HasPrefix(result, "Runtime error"):
				status = JudgeStatusRE
			case strings.HasPrefix(result, "Time limit exceeded"):
				status = JudgeStatusTLE
			case strings.HasPrefix(result, "Memory limit exceeded"):
				status = JudgeStatusMLE
			default:
				status = JudgeStatusUNK
			}
			statuses[submissionID] = status
			results[submissionID] = result
		})

		sort.Slice(submissionIDs, func(i, j int) bool { return submissionIDs[i] < submissionIDs[j] })

		for _, k := range submissionIDs {
			v := statuses[k]

			// ジャッジ待ちのものはスキップ
			if _, ok := statuses[k]; !ok {
				continue
			}

			if _, ok := judgeFinished[k]; !ok {
				clearLine()
				fmt.Println("url: " + cf.url + fmt.Sprintf("contest/%d/submission/%d", contestID, k))
				fmt.Println("\tproblem: " + problems[k])
				fmt.Println("\tverdict: " + v.GetColorESCS() + results[k] + util.ESCS_COL_OFF)
				fmt.Println()
				judgeFinished[k] = struct{}{} // ジャッジ済み
			}
		}

		if currentWaitingCount == 0 {
			break
		}

		if waitingTotal == -1 {
			waitingTotal = currentWaitingCount
		}

		clearLine()
		util.SaveCursorPos() // カーソル位置を保存
		fmt.Print(util.ESCS_COL_REVERSE + "waiting for judge (" + strconv.Itoa(waitingTotal-currentWaitingCount) + "/" + strconv.Itoa(waitingTotal) + ")" + util.ESCS_COL_OFF)
		messagePrinted = true

		time.Sleep(1 * time.Minute) // 30秒ごとに確認
	}
	clearLine()
}
