package task

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// BugReport: output to TaskConfig.GetBugsPath() / BugId - SqlId - MutationName .log and
// TaskConfig.GetBugsPath() / BugId - SqlId - MutationName .json
type BugReport struct {
	ReportTime     string            `json:"reportTime"`
	BugId          int               `json:"bugId"`
	SqlId          int               `json:"sqlId"`
	MutationName   string            `json:"mutationName"`
	IsUpper        bool              `json:"isUpper"` // true: theoretically, OriginResult < NewResult
	OriginalSql    string            `json:"originalSql"`
	OriginalResult *connector.Result `json:"-"`
	MutatedSql     string            `json:"mutatedSql"`
	MutatedResult  *connector.Result `json:"-"`
}

func (bugReport *BugReport) ToString() string {
	str := "**************************************************\n"
	str += "[MutationName] " + bugReport.MutationName + "\n"
	str += "**************************************************\n"
	str += "[IsUpper] " + strconv.FormatBool(bugReport.IsUpper) + "\n"
	str += "**************************************************\n"
	str += "[OriginalResult]\n"
	str += bugReport.OriginalResult.ToString() + "\n"
	str += "**************************************************\n"
	str += "[MutatedResult]\n"
	str += bugReport.MutatedResult.ToString() + "\n"
	str += "**************************************************\n"
	str += "\n"
	str += "-- OriginalSql\n"
	str += bugReport.OriginalSql + ";\n"
	str += "-- MutatedSql\n"
	str += bugReport.MutatedSql + ";\n"
	return str
}

// BugReport.SaveBugReport: output to TaskConfig.GetTaskBugsPath() / BugId - SqlId - MutationName .log and
// TaskConfig.GetTaskBugsPath() / BugId - SqlId - MutationName .json.
//
// Note that create if not exists TaskConfig.GetTaskBugsPath()
func (bugReport *BugReport) SaveBugReport(taskBugsPath string) error {
	_ = os.Mkdir(taskBugsPath, 0777)
	bugSig := "bug-" + strconv.Itoa(bugReport.BugId) + "-" + strconv.Itoa(bugReport.SqlId) + "-" + bugReport.MutationName
	// log
	bugLogPath := path.Join(taskBugsPath, bugSig+".log")
	err := ioutil.WriteFile(bugLogPath, []byte(bugReport.ToString()), 0777)
	if err != nil {
		return errors.Wrap(err, "[BugReport.SaveBugReport]write log error")
	}
	// json
	bugJsonPath := path.Join(taskBugsPath, bugSig+".json")
	jsonData, err := json.Marshal(bugReport)
	if err != nil {
		return errors.Wrap(err, "[BugReport.SaveBugReport]marshal error")
	}
	err = ioutil.WriteFile(bugJsonPath, jsonData, 0777)
	if err != nil {
		return errors.Wrap(err, "[BugReport.SaveBugReport]write json error")
	}
	return nil
}