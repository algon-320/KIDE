package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/online_judge"
	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/util"
)

func downloadSampleCase(problemURL string) error {
	oj, err := online_judge.FromProblemURL(problemURL)
	if err != nil {
		return err
	}
	err = oj.NewProblem(problemURL)
	if err != nil {
		return err
	}
	return nil
}

func run(lang language.Language) error {
	filename, err := language.FindSourceCode(lang)
	if err != nil {
		return err
	}

	util.DebugPrint(fmt.Sprintf("running --> %s", filename))

	_, err = lang.Run(filename, "", true) // 結果を表示しながら実行
	if err != nil {
		return err
	}
	return nil
}

// caseID : 負ならすべてのサンプルケースをテスト
func tester(lang language.Language, problemID string, caseID int) error {
	filename, err := language.FindSourceCode(lang)
	if err != nil {
		return err
	}

	p, err := online_judge.LoadProblem(problemID)
	if err != nil {
		return err
	}

	if caseID < 0 {
		// すべてのサンプルケースをテスト
		samplePassed := true
		for _, c := range p.Cases {
			out, err := lang.Run(filename, c.Input, false) // 画面出力しないで実行
			if err != nil {
				return err
			}

			if out != c.Output {
				util.PrintTitle(30, 4, "=", "input")
				fmt.Print(c.Input)
				util.PrintTitle(30, 4, "=", "your answer")
				fmt.Print(out)
				util.PrintTitle(30, 4, "=", "correct answer")
				fmt.Print(c.Output)
				fmt.Println(strings.Repeat("=", 30))

				fmt.Println(util.ESCS_COL_RED_B + "Wrong answer" + util.ESCS_COL_OFF)
				samplePassed = false
			}
		}
		if samplePassed {
			fmt.Println(util.ESCS_COL_GREEN_B + "Samplecases passed" + util.ESCS_COL_OFF)
			return submit(filename, lang, p) // 確認して提出
		}
	} else if 0 < caseID && caseID <= len(p.Cases) {
		c := p.Cases[caseID-1]

		util.PrintTitle(30, 4, "=", "input")
		fmt.Print(c.Input)
		util.PrintTitle(30, 4, "=", "output")
		out, err := lang.Run(filename, c.Input, true) // 画面出力しながら実行
		fmt.Println(strings.Repeat("=", 30))

		if err != nil {
			return err
		}

		if out == c.Output {
			fmt.Println(util.ESCS_COL_GREEN_B + "Passed" + util.ESCS_COL_OFF)
		} else {
			fmt.Println(util.ESCS_COL_RED_B + "Wrong answer" + util.ESCS_COL_OFF)
			util.PrintTitle(30, 4, "=", "your answer")
			fmt.Print(out)
			util.PrintTitle(30, 4, "=", "correct answer")
			fmt.Print(c.Output)
			fmt.Println(strings.Repeat("=", 30))
		}
	} else {
		return fmt.Errorf(util.PrefixError+"case id should be 1 to %d", len(p.Cases))
	}
	return nil
}

func submit(souceFilename string, lang language.Language, p *online_judge.Problem) error {
	fmt.Printf("Do you really submit the solution `%s` to problem `%s` ?\n", souceFilename, p.Name)
	yes := util.AskYesNo()
	if !yes {
		fmt.Println(util.PrefixInfo + "Submit cancelled.")
		return nil
	}

	// この時点で提出するソースコードが確定
	sourceCodeBytes, err := ioutil.ReadFile(souceFilename)
	if err != nil {
		return err
	}
	sourceCodeStr := string(sourceCodeBytes)

	// process
	sourceCodeStr = processSource(sourceCodeStr)

	if sourceCodeStr == "" {
		fmt.Println(util.PrefixInfo + "Submit cancelled.")
		return nil
	}

	res, err := p.Oj.Submit(p, sourceCodeStr, lang)
	if err != nil {
		return err
	}
	res.Print() // ジャッジ結果を表示

	// 保存
	var saveSourceFileAfterAccepted bool
	tmp, ok := setting.Get("General.SaveSourceFileAfterAccepted", "")
	if !ok {
		fmt.Println("Do you want to copy the source file after the solution is accepted (or pretests passed) ?")
		saveSourceFileAfterAccepted = util.AskYesNo()
		setting.Set("General.SaveSourceFileAfterAccepted", saveSourceFileAfterAccepted)
	} else {
		saveSourceFileAfterAccepted = tmp.(bool)
	}

	if (res.Status == online_judge.JudgeStatusAC || res.Status == online_judge.JudgeStatusPP) &&
		saveSourceFileAfterAccepted {
		// header
		// line1: problem url
		// line2: submission url
		// line3: submitted date
		// line4: judge result
		// line5: (empty line)
		addedSource := lang.CommentOut("problem: "+p.URL) + "\n" +
			lang.CommentOut("submission: "+res.URL) + "\n" +
			lang.CommentOut(res.Date.String()) + "\n" +
			lang.CommentOut(res.Status.ToString()) + "\n\n" +
			sourceCodeStr
		sourceFileName := p.Name + "_" + res.Date.Format("20060102150405") + lang.FileExtension()
		saveSourceFile(sourceFileName, []byte(addedSource), p)
	}
	return nil
}

// ソースコードを General.SaveSourceFileDirectory 以下に保存する
func saveSourceFile(sourceFilename string, sourceCode []byte, p *online_judge.Problem) error {
	exeDir, _ := os.Executable()
	exeDir = filepath.Dir(exeDir)

	var saveSourceFileDir string
	tmp, ok := setting.Get("General.SaveSourceFileDirectory", "")
	if ok {
		saveSourceFileDir = tmp.(string)
		expanded := strings.Replace(saveSourceFileDir, "{EXE_DIR}", exeDir, 1)
		if !util.FileExists(expanded) {
			ok = false
		}
	}

	if !ok {
		for {
			fmt.Println("Put the directory path to save the source file after the solution was accepted.")
			fmt.Println("You can use `{EXE_DIR}` (without quotes) to designate the executable file's directory.")
			saveSourceFileDir = util.AskString("path")
			expanded := strings.Replace(saveSourceFileDir, "{EXE_DIR}", exeDir, 1)

			// 有効なパスか確認
			if err := os.MkdirAll(expanded, 0755); err == nil {
				break
			} else {
				fmt.Println("Invalid path. Try again ...")
			}
		}
		setting.Set("General.SaveSourceFileDirectory", saveSourceFileDir)
	}

	// 現在のディレクトリを保持しておいて、抜けるときに復元する
	prev, _ := os.Getwd()
	defer os.Chdir(prev)

	expanded := strings.Replace(saveSourceFileDir, "{EXE_DIR}", exeDir, 1)
	os.Chdir(expanded)

	if !util.FileExists(p.Oj.Name()) {
		if err := os.Mkdir(p.Oj.Name(), 0775); err != nil {
			return fmt.Errorf(util.PrefixError+"%s", err)
		}
		absPath, _ := filepath.Abs(p.Oj.Name())
		fmt.Println(fmt.Sprintf(util.PrefixInfo+"Created directory `%s`", absPath))
	}
	os.Chdir(p.Oj.Name())

	if err := ioutil.WriteFile(sourceFilename, sourceCode, 0644); err != nil {
		return fmt.Errorf(util.PrefixError+"%s", err)
	}

	fmt.Println(util.PrefixInfo + fmt.Sprintf("Saved the source file as `%s`", sourceFilename))
	return nil
}

// ソースコードを整形する
func processSource(sourceCode string) string {
	// 整形ツール(自分用) : TODO
	mycpp := func(lines []string) string {
		const (
			SKIPBEGIN      = "//SKIPBEGIN"
			SKIPEND        = "//SKIPEND"
			PROBLEM_STRING = "KIDE_PROBLEM_"
		)

		ret := ""
		state := 0
		skip := false
		problemID := ""
		for i := 0; i < len(lines); i++ {
			// SKIPBEGIN ~ SKIPEND の行を無視
			if strings.TrimSpace(lines[i]) == SKIPBEGIN {
				for strings.TrimSpace(lines[i]) != SKIPEND {
					i++
				}
				continue
			}

			if len(lines[i]) >= 1 && lines[i][:1] == "#" {
				// 問題を指定する行
				re := regexp.MustCompile("#define " + PROBLEM_STRING + "(.*)")
				group := re.FindSubmatch([]byte(lines[i]))
				if group != nil {
					problemID = string(group[1])
					skip = true
					ret += "\n"
					continue

				} else if len(lines[i]) >= 3 && lines[i][:3] == "#if" {
					// #if / #ifdef / #ifndef の行
					state++
					re = regexp.MustCompile("#ifdef " + PROBLEM_STRING)
					group = re.FindSubmatch([]byte(lines[i]))
					if group != nil {
						if lines[i] == ("#ifdef " + PROBLEM_STRING + problemID) {
							skip = false
						}
						continue
					}

				} else if len(lines[i]) >= 6 && lines[i][:6] == "#endif" {
					// #endif の行
					if state == 0 {
						fmt.Println(util.PrefixError + "Missing format (processSource)")
						return ""
					}

					state--
					if state == 0 {
						skip = true
						continue
					}
				}
			}

			if skip {
				continue
			}
			ret += lines[i] + "\n"
		}

		return ret
	}
	lines := strings.Split(sourceCode, "\n")
	return mycpp(lines)
}
