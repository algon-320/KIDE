package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/online_judge"
	"github.com/algon-320/KIDE/setting"
	"github.com/algon-320/KIDE/util"
	"golang.org/x/crypto/ssh/terminal"
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
	fd := int(os.Stdout.Fd())
	termWidth, _, err := terminal.GetSize(fd)
	if err != nil {
		termWidth = 80
	}

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

			// WA
			if out != c.Output {
				util.PrintTitle(termWidth, 4, "=", "input")
				fmt.Print(c.Input)
				util.PrintTitle(termWidth, 4, "=", "your answer")
				fmt.Print(out)
				util.PrintTitle(termWidth, 4, "=", "correct answer")
				fmt.Print(c.Output)
				fmt.Println(strings.Repeat("=", termWidth))

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

		util.PrintTitle(termWidth, 4, "=", "input")
		fmt.Print(c.Input)
		util.PrintTitle(termWidth, 4, "=", "output")
		out, err := lang.Run(filename, c.Input, true) // 画面出力しながら実行
		fmt.Println(strings.Repeat("=", termWidth))

		if err != nil {
			return err
		}

		if out == c.Output {
			fmt.Println(util.ESCS_COL_GREEN_B + "Passed" + util.ESCS_COL_OFF)
		} else {
			fmt.Println(util.ESCS_COL_RED_B + "Wrong answer" + util.ESCS_COL_OFF)
			util.PrintTitle(termWidth, 4, "=", "your answer")
			fmt.Print(out)
			util.PrintTitle(termWidth, 4, "=", "correct answer")
			fmt.Print(c.Output)
			fmt.Println(strings.Repeat("=", termWidth))
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
	if tmp, ok := setting.Get("General.SaveSourceFileAfterAccepted", ""); ok {
		saveSourceFileAfterAccepted = tmp.(bool)
	} else {
		fmt.Println("Do you want to copy the source file after the solution is accepted (or pretests passed) ?")
		saveSourceFileAfterAccepted = util.AskYesNo()
		setting.Set("General.SaveSourceFileAfterAccepted", saveSourceFileAfterAccepted)
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
	var tmp interface{}
	var ok bool
	if tmp, ok = setting.Get("General.SaveSourceFileDirectory", ""); ok {
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
// 設定で実行コマンドが指定されていた場合に、標準入力にソースコードを投げ、標準出力から読み取ったものを返す
// {EXE_DIR}を実行ファイルのあるディレクトリのパスとして使える
func processSource(sourceCode string) string {
	var cmd *exec.Cmd
	if tmp, exist := setting.Get("General.SourcecodeProcess.Command", ""); exist {
		exeDir, _ := os.Executable()
		exeDir = filepath.Dir(exeDir)
		cmdStr := tmp.(string)
		expanded := strings.Replace(cmdStr, "{EXE_DIR}", exeDir, 1)
		cmd = util.Command(expanded)

		cmd.Stdin = bytes.NewBufferString(sourceCode)
		cmd.Stderr = os.Stderr
		var out bytes.Buffer
		cmd.Stdout = &out

		cmd.Run()
		return out.String()
	}
	return sourceCode
}
