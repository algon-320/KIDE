package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/snippet_manager"

	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/online_judge"
	"github.com/algon-320/KIDE/util"
	"github.com/urfave/cli"
)

func cmdRun(c *cli.Context) error {
	lang := language.GetLanguage(c.String("language"))
	if err := run(lang); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func cmdTester(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError(util.PrefixError+"few args", 1)
	}

	lang := language.GetLanguage(c.String("language"))
	problemID := c.Args().First()
	if err := tester(lang, problemID, c.Int("case")); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func cmdDl(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError(util.PrefixError+"few args", 1)
	}

	url := c.Args().First()
	if err := downloadSampleCase(url); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func cmdSubmit(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError(util.PrefixError+"few args", 1)
	}

	lang := language.GetLanguage(c.String("language"))

	filename, err := language.FindSourceCode(lang)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	p, err := online_judge.LoadProblem(c.Args().First())
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = submit(filename, lang, p)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func cmdView(c *cli.Context) error {
	if c.NArg() < 1 {
		// 引数が無い場合はすべて表示
		list := online_judge.GetAllProblemID()

		// make table
		title := []string{"id", "name", "oj name", "url"}
		data := [][]string{}
		for _, v := range list {
			p, err := online_judge.LoadProblem(v)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			data = append(data, []string{v, p.Name, p.Oj.Name(), p.URL})
		}

		util.PrintTable(title, data, true)
	} else {
		// 引数がある場合は引数で指定された問題を表示する
		problemID := c.Args().First()
		p, err := online_judge.LoadProblem(problemID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		p.Print()
	}
	return nil
}

func cmdProcesser(c *cli.Context) error {
	lang := language.GetLanguage(c.String("language"))

	filename, err := language.FindSourceCode(lang)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	sourceCode, err := ioutil.ReadFile(filename)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	sourceCodeStr := string(sourceCode)

	sourceCodeStr = processSource(sourceCodeStr)
	fmt.Print(sourceCodeStr)
	return nil
}

func cmdAtCoderConv(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError(util.PrefixError+"few args", 1)
	}
	url := c.Args().First()
	re := regexp.MustCompile("https://(.+)\\.contest.atcoder.jp/(.*)")
	group := re.FindSubmatch([]byte(url))
	if group == nil {
		return cli.NewExitError(util.PrefixError+"error", 1)
	}
	contest := string(group[1])
	var suffix string
	if len(group) > 2 {
		suffix = string(group[2])
	}
	newURL := "https://beta.atcoder.jp/contests/" + contest + "/" + suffix
	fmt.Println(newURL)
	return nil
}

func cmdCodeforcesMySubmissionsViewer(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError(util.PrefixError+"few args", 1)
	}

	contestID, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return cli.NewExitError(util.PrefixError+"designate contest id like\"1038\"", 1)
	}
	online_judge.Codeforces.ShowMySubmissions(contestID)
	return nil
}

func cmdSnippetManager(c *cli.Context) error {
	var editorList = []string{"markdown (library output)"}
	for _, e := range snippet_manager.EditorList {
		editorList = append(editorList, e.Name())
	}
	i := util.AskChoose(editorList, "Select format you want to print with")
	if i == 0 {
		md := snippet_manager.GenerateMarkdown()
		fmt.Println(md)
	} else {
		editorName := editorList[i]
		var editor snippet_manager.Editor
		for _, e := range snippet_manager.EditorList {
			if e.Name() == editorName {
				editor = e
				break
			}
		}
		snip := snippet_manager.ExportSnippets(editor)
		fmt.Println(snip)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "KIDE"
	app.Usage = "Kyopro-Iikanjini-Dekiru-Environment"

	var defaultLangName string
	if tmp, exist := setting.Get("Language.DefaultLanguageName", ""); exist {
		defaultLangName = tmp.(string)
	}

	app.Commands = []cli.Command{
		{
			Name:      "run",
			Aliases:   []string{"r"},
			Usage:     "Runs the source code here",
			UsageText: "run [command options]",
			Action:    cmdRun,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "language, l",
					Value: defaultLangName,
					Usage: "runnig as `LANGUAGE`",
				},
			},
		},
		{
			Name:      "tester",
			Aliases:   []string{"t"},
			Usage:     "Tests samplecases",
			UsageText: "tester [problem id] [command options]",
			Action:    cmdTester,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "language, l",
					Value: defaultLangName,
					Usage: "testing in `LANGUAGE`",
				},
				cli.IntFlag{
					Name:  "case, c",
					Value: -1,
					Usage: "testing only one case. `INDEX` is index of samples (1-indexed value)",
				},
			},
		},
		{
			Name:      "dl",
			Aliases:   []string{"d"},
			Usage:     "Downloads samplecases of the problem",
			UsageText: "dl [problem page url]",
			Action:    cmdDl,
		},
		{
			Name:      "submit",
			Aliases:   []string{"s"},
			Usage:     "Submits the solution",
			UsageText: "submit [problem id] [command options]",
			Action:    cmdSubmit,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "language, l",
					Value: defaultLangName,
					Usage: "submit as `LANGUAGE`",
				},
			},
		},
		{
			Name:    "view",
			Aliases: []string{"v"},
			Usage:   "Shows problems",
			Action:  cmdView,
		},
		{
			Name:      "processer",
			Aliases:   []string{"p"},
			Usage:     "Processes source code and output",
			UsageText: "processer [command options]",
			Action:    cmdProcesser,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "language, l",
					Value: defaultLangName,
					Usage: "processing as `LANGUAGE`",
				},
			},
		},
		{
			Name:      "cf-mysubmissions",
			Usage:     "Shows mysubmissions of codeforces contest",
			UsageText: "cf-mysubmissions [contest id]",
			Action:    cmdCodeforcesMySubmissionsViewer,
		},
		{
			Name:   "snippet",
			Action: cmdSnippetManager,
		},
	}

	app.Run(os.Args)
}
