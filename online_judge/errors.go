package online_judge

import (
	"fmt"

	"github.com/algon-320/KIDE/util"
)

type ErrFailedToLogin struct {
	oj_name string
	message string
}

func (e ErrFailedToLogin) Error() string {
	if e.message == "" {
		return util.PrefixError + fmt.Sprintf("Failed to login to `%s`", e.oj_name)
	} else {
		return util.PrefixError + fmt.Sprintf("Failed to login to `%s` : %s", e.oj_name, e.message)
	}
}

//-----------------

type ErrNoSuchOnlineJudge struct {
	oj_name string
}

func (e ErrNoSuchOnlineJudge) Error() string {
	return util.PrefixError + fmt.Sprintf("No such online judge `%s`", e.oj_name)
}

//-----------------

type ErrInvalidProblemURL struct {
	url string
}

func (e ErrInvalidProblemURL) Error() string {
	return util.PrefixError + fmt.Sprintf("Invalid problem url `%s`", e.url)
}

//-----------------

type ErrUnsuportedLanguage struct {
	name string
}

func (e ErrUnsuportedLanguage) Error() string {
	return util.PrefixError + fmt.Sprintf("Unsuported language `%s`", e.name)
}

//-----------------

type ErrFailedToSubmit struct {
	message string
}

func (e ErrFailedToSubmit) Error() string {
	return util.PrefixError + fmt.Sprintf("Failed to Submit the solution : %s", e.message)
}

//-----------------

type ErrFailedToLoadSamplecase struct {
	message string
}

func (e ErrFailedToLoadSamplecase) Error() string {
	return util.PrefixError + fmt.Sprintf("Failed to Load samplecase : %s", e.message)
}
