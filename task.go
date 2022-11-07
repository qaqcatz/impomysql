package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/mutation/stage1"
	"github.com/qaqcatz/impomysql/mutation/stage2"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

type TaskConfig struct {
	OutputPath      string `json:"outputPath"` // default: ./output
	DBMS            string `json:"dbms"`       // default: mysql
	TaskId          int    `json:"taskId"`     // must >= 0, the final output path: OutputPath / DBMS / TaskId
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	DbName          string `json:"dbname"`
	MysqlClientPath string `json:"mysqlClientPath"` // default: /usr/bin/mysql
	DDLPath         string `json:"ddlPath"`         // can not have ';' in comment
	DMLPath         string `json:"dmlPath"`         // can not hava ';' in comment
	Seed            int64  `json:"seed"`            // <= 0: current time
}

func (taskConfig *TaskConfig) GetOutputPath() string {
	return path.Join(taskConfig.OutputPath, taskConfig.DBMS, "task-"+strconv.Itoa(taskConfig.TaskId))
}

// taskInputCheck: check input, assign default value
func TaskInputCheck(config *TaskConfig) error {
	if config.OutputPath == "" {
		config.OutputPath = "./output"
	}
	if config.DBMS == "" {
		config.DBMS = "mysql"
	}
	// Host, Port, Username, Password, DbName, MysqlClientPath can be checked by connector
	if config.TaskId < 0 {
		return errors.New("TaskInputCheck: taskId < 0 ")
	}
	if config.DDLPath == "" {
		return errors.New("TaskInputCheck: empty ddlPath")
	}
	if config.DMLPath == "" {
		return errors.New("TaskInputCheck: empty dmlPath")
	}
	if config.Seed <= 0 {
		config.Seed = time.Now().UnixNano()
	}
	return nil
}

type RdSql struct {
	Id  int    `json:"id"`
	Sql string `json:"sql"`
}

// ExtractSQL: Extract sql statements by ';':
//   - ignore the ';' in ``, '', "";
//   - ignore the escaped characters in ``, '', "";
// Note that: Comments cannot have ';'
func ExtractSQL(s string) []*RdSql {
	res := make([]*RdSql, 0)
	start := 0
	flag := -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			if flag == -1 {
				flag = '\''
			} else {
				if flag == '\'' {
					flag = -1
				}
			}
		case '"':
			if flag == -1 {
				flag = '"'
			} else {
				if flag == '"' {
					flag = -1
				}
			}
		case '`':
			if flag == -1 {
				flag = '`'
			} else {
				if flag == '`' {
					flag = -1
				}
			}
		case '\\':
			if flag != -1 {
				i++
			}
		case ';':
			if flag == -1 {
				res = append(res, &RdSql{
					Id:  len(res),
					Sql: s[start : i+1],
				})
				start = i + 1
			}
		default:
			continue
		}
	}
	return res
}

// SaveSqlsJson: output []*RdSql to jsonPath(outputPath/DBMS/taskId/ ddl | dml .json)
func SaveSqlsJson(sqls []*RdSql, taskOutputPath string, isDDL bool) error {
	ddml := "ddl.json"
	if !isDDL {
		ddml = "dml.json"
	}
	jsonPath := path.Join(taskOutputPath, ddml)
	jsonData, err := json.Marshal(sqls)
	if err != nil {
		return errors.New("SaveSqlsJson: marshal error: " + err.Error())
	}
	err = ioutil.WriteFile(jsonPath, jsonData, 0777)
	if err != nil {
		return errors.New("SaveSqlsJson: write error: " + err.Error())
	}
	return nil
}

// InitDDLSqls: init database and execute ddl sqls
func InitDDLSqls(ddlSqls []*RdSql, conn *connector.Connector) error {
	err := conn.InitDB()
	if err != nil {
		return errors.New("InitDDLSqls: init database error: " + err.Error())
	}
	for i, ddlSql := range ddlSqls {
		result := conn.ExecSQL(ddlSql.Sql)
		if result.Err != nil {
			return errors.New("InitDDLSqls: exec ddl sql " + strconv.Itoa(i) + " error: " + result.Err.Error())
		}
	}
	return nil
}

// TaskResult: result of task, output to outputPath/DBMS/taskId/result.json.
type TaskResult struct {
	StartTime        string `json:"startTime"`
	DDLSqlsNum       int    `json:"ddlSqlsNum"`
	DMLSqlsNum       int    `json:"dmlSqlsNum"`
	EndInitTime      string `json:"endInitTime"`
	Stage1ErrNum     int    `json:"stage1ErrNum"`
	Stage1ExecErrNum int    `json:"stage1ExecErrNum"`
	Stage2ErrNum     int    `json:"stage2ErrNum"`
	Stage2ExecNum    int    `json:"stage2ExecNum"`
	Stage2ExecErrNum int    `json:"stage2ExecErrNum"`
	ImpoBugsNum      int    `json:"impoBugsNum"`
	SaveBugErrNum    int    `json:"saveBugErrNum"`
	EndTime          string `json:"endTime"`
}

// SaveTaskResult: output to outputPath/DBMS/taskId/result.json
func (taskResult *TaskResult) SaveTaskResult(taskOutputPath string) error {
	jsonPath := path.Join(taskOutputPath, "result.json")
	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		return errors.New("SaveTaskResult: marshal error: " + err.Error())
	}
	err = ioutil.WriteFile(jsonPath, jsonData, 0777)
	if err != nil {
		return errors.New("SaveTaskResult: write json error: " + err.Error())
	}
	return nil
}

// BugReport: output to outputPath/DBMS/taskId/bugs/ BugId @ SqlId @ MutationName .json
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

// ToString: output to outputPath/DBMS/taskId/bugs/ BugId @ SqlId @ MutationName .log
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

// SaveLogicalBug: output to outputPath/DBMS/taskId/bugs/ bug - BugId - SqlId - MutationName .json and
// outputPath/DBMS/taskId/bugs/ bug - BugId - SqlId - MutationName .log
func SaveLogicalBug(bugsPath string, bugId int, sqlId int, mutationName string, isUpper bool,
	originalSql string, originalResult *connector.Result, mutatedSql string, mutatedResult *connector.Result) error {
	bugReport := &BugReport{
		ReportTime:     time.Now().String(),
		BugId:          bugId,
		SqlId:          sqlId,
		MutationName:   mutationName,
		IsUpper:        isUpper,
		OriginalSql:    originalSql,
		OriginalResult: originalResult,
		MutatedSql:     mutatedSql,
		MutatedResult:  mutatedResult,
	}
	bugSig := "bug-" + strconv.Itoa(bugId) + "-" + strconv.Itoa(sqlId) + "-" + mutationName
	bugJsonPath := path.Join(bugsPath, bugSig+".json")
	bugLogPath := path.Join(bugsPath, bugSig+".log")
	jsonData, err := json.Marshal(bugReport)
	if err != nil {
		return errors.New("SaveLogicalBug: marshal error: " + err.Error())
	}
	err = ioutil.WriteFile(bugJsonPath, jsonData, 0777)
	if err != nil {
		return errors.New("SaveLogicalBug: write json error: " + err.Error())
	}
	err = ioutil.WriteFile(bugLogPath, []byte(bugReport.ToString()), 0777)
	if err != nil {
		return errors.New("SaveLogicalBug: write log error: " + err.Error())
	}
	return nil
}

// Run:  Basic task for finding logical bugs:
//
// 0. check input
//
// 1. init:
//   1.1 create if not exists outputPath/DBMS/taskId
//   1.2 init outputPath/DBMS/taskId/task.log, create logger, write into outputPath/DBMS/taskId/task.log
//   1.3 create connector
//   1.4 read ddl, write ddl into outputPath/DBMS/taskId/ ddl .json, init database, execute ddl
//   1.5 read dml, write dml into outputPath/DBMS/taskId/ dml .json
//   1.6 init dir outputPath/DBMS/taskId/bugs
//
// 2. run:
// for each dml sql, do:
//   2.1 stage1.InitAndExec
//   2.2 stage2.MutateAllAndExec
//   2.3 use oracle.Check to detect logical bugs, save logical bugs into
//   outputPath/DBMS/taskId/bugs/ BugId @ SqlId @ MutationName .json and
//   outputPath/DBMS/taskId/bugs/ BugId @ SqlId @ MutationName .log, see BugReport.
//   2.4 save the result of stage1+stage2 into outputPath/DBMS/taskId/result.json, see TaskResult
//
// Note that: remove database after task finished
func RunTask(config *TaskConfig) error {
	// 0. check input
	if err := TaskInputCheck(config); err != nil {
		return errors.New("RunTask: input check error: " + err.Error())
	}

	// 1 init
	startTime := time.Now().String()
	// **************************************************
	// 1.1 create if not exists outputPath/DBMS/taskId
	taskOutputPath := config.GetOutputPath()
	_ = os.MkdirAll(taskOutputPath, 0777)
	// 1.2 init outputPath/DBMS/taskId/task.log, create logger, write into outputPath/DBMS/taskId/task.log
	taskLogPath := path.Join(taskOutputPath, "task.log")
	err := ioutil.WriteFile(taskLogPath, []byte(""), 0777)
	if err != nil {
		return errors.New("RunTask: init log file error: " + err.Error())
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	file, err := os.OpenFile(taskLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return errors.New("RunTask: create logger error: " + err.Error())
	}
	writers := []io.Writer{file}
	multiWriter := io.MultiWriter(writers...)
	logger.SetOutput(multiWriter)
	logger.SetLevel(logrus.InfoLevel)
	// 1.3 create connector
	logger.Info("create connector")
	conn, err := connector.NewConnector(config.Host, config.Port, config.Username, config.Password, config.DbName)
	if err != nil {
		logger.Error("create connector error: " + err.Error())
		return errors.New("RunTask: create connector error: " + err.Error())
	}
	// 1.4 read ddl, write ddl into outputPath/DBMS/taskId/ ddl .json, init database, execute ddl
	logger.Info("init ddl")
	ddlData, err := ioutil.ReadFile(config.DDLPath)
	if err != nil {
		logger.Error("read ddl error: " + err.Error())
		return errors.New("RunTask: read ddl error: " + err.Error())
	}
	ddlSqls := ExtractSQL(string(ddlData))
	err = SaveSqlsJson(ddlSqls, taskOutputPath, true)
	if err != nil {
		logger.Error("write ddl.json error: " + err.Error())
		return errors.New("RunTask: write ddl.json error: " + err.Error())
	}
	err = InitDDLSqls(ddlSqls, conn)
	// remove database after task finished
	defer conn.RmDB()
	if err != nil {
		logger.Error("init ddl sqls error: " + err.Error())
		return errors.New("RunTask: init ddl sqls error: " + err.Error())
	}
	// 1.5 read dml, write dml into outputPath/DBMS/taskId/ dml .json
	logger.Info("init dml")
	dmlData, err := ioutil.ReadFile(config.DMLPath)
	if err != nil {
		logger.Error("read dml error: " + err.Error())
		return errors.New("RunTask: read dml error: " + err.Error())
	}
	dmlSqls := ExtractSQL(string(dmlData))
	err = SaveSqlsJson(dmlSqls, taskOutputPath, false)
	if err != nil {
		logger.Error("write dml.json error: " + err.Error())
		return errors.New("RunTask: write dml.json error: " + err.Error())
	}
	// 1.6 init dir outputPath/DBMS/taskId/bugs
	// mkdir bugs
	bugsPath := path.Join(taskOutputPath, "bugs")
	_ = os.RemoveAll(bugsPath)
	_ = os.Mkdir(bugsPath, 0777)
	// **************************************************
	endInitTime := time.Now().String()
	// end 1

	// 2. run
	logger.Info("Running **************************************************")
	taskResult := &TaskResult{
		StartTime:        startTime,
		DDLSqlsNum:       len(ddlSqls),
		DMLSqlsNum:       len(dmlSqls),
		EndInitTime:      endInitTime,
		Stage1ErrNum:     0,
		Stage1ExecErrNum: 0,
		Stage2ErrNum:     0,
		Stage2ExecNum:    0,
		Stage2ExecErrNum: 0,
		ImpoBugsNum:      0,
		SaveBugErrNum:    0,
		EndTime:          "",
	}
	// **************************************************
	// for each sql, do:
	total := len(dmlSqls)
	cur := 0
	for i, dmlSql := range dmlSqls {

		// rate
		if cur > total/20 {
			cur = 0
			logger.Info("[Rate]", i, "/", total)
		} else {
			cur += 1
		}

		// 2.1 stage1.InitAndExec
		stage1Result := stage1.InitAndExec(dmlSql.Sql, conn)
		// handle stage1 error
		if stage1Result.Err != nil {
			taskResult.Stage1ErrNum += 1
			logger.Error("--------------------------------------------------")
			logger.Error("[Stage1 Error]", "(", dmlSql.Id, ")", dmlSql.Sql)
			logger.Error(stage1Result.Err)
			logger.Error("--------------------------------------------------")
			continue
		}
		// handle stage1 execute error
		if stage1Result.ExecResult.Err != nil {
			taskResult.Stage1ExecErrNum += 1
			logger.Error("==================================================")
			logger.Error("[Stage1 Exec Error]", "(", dmlSql.Id, ")", stage1Result.InitSql)
			logger.Error(stage1Result.ExecResult.Err)
			logger.Error("==================================================")
			continue
		}

		originalSql := stage1Result.InitSql
		originalResult := stage1Result.ExecResult

		// 2.2 stage2.MutateAllAndExec
		stage2Result := stage2.MutateAllAndExec(originalSql, config.Seed+int64(i), conn)
		// handle stage2 error
		if stage2Result.Err != nil {
			taskResult.Stage2ErrNum += 1
			logger.Error("--------------------------------------------------")
			logger.Error("[Stage2 Error]", "(", dmlSql.Id, ")", originalSql)
			logger.Error(stage2Result.Err)
			logger.Error("--------------------------------------------------")
			continue
		}
		// for each mutation unit
		taskResult.Stage2ExecNum += len(stage2Result.MutateUnits)
		for _, mutateUnit := range stage2Result.MutateUnits {
			// handle stage2 execute error
			if mutateUnit.Err != nil { // treat as stage2 execute error
				taskResult.Stage2ExecErrNum += 1
				logger.Error("==================================================")
				logger.Error("[Stage2 Exec Error]", "(", dmlSql.Id, "-", mutateUnit.Name, ")", mutateUnit.Sql)
				logger.Error(mutateUnit.Err)
				logger.Error("==================================================")
				continue
			}
			if mutateUnit.ExecResult.Err != nil {
				taskResult.Stage2ExecErrNum += 1
				logger.Error("==================================================")
				logger.Error("[Stage2 Exec Error]", "(", dmlSql.Id, "-", mutateUnit.Name, ")", mutateUnit.Sql)
				logger.Error(mutateUnit.ExecResult.Err)
				logger.Error("==================================================")
				continue
			}

			mutationName := mutateUnit.Name
			isUpper := mutateUnit.IsUpper
			mutatedSql := mutateUnit.Sql
			mutatedResult := mutateUnit.ExecResult

			//   2.3 use oracle.Check to detect logical bugs, save logical bugs into
			//   outputPath/DBMS/taskId/bugs/ BugId @ SqlId @ MutationName .json and
			//   outputPath/DBMS/taskId/bugs/ BugId @ SqlId @ MutationName .log, see BugReport.
			if !oracle.Check(originalResult, mutatedResult, isUpper) {
				// logical bug!!!
				bugId := taskResult.ImpoBugsNum
				taskResult.ImpoBugsNum += 1
				fmt.Println("task", config.TaskId, "detected a logical bug!", bugId, dmlSql.Id, mutationName)
				err := SaveLogicalBug(bugsPath, bugId, dmlSql.Id, mutationName, isUpper,
					originalSql, originalResult, mutatedSql, mutatedResult)
				if err != nil {
					taskResult.SaveBugErrNum += 1
					logger.Error("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
					logger.Error("[Save Bug Error] ", "bug-", bugId, "-", dmlSql.Id, "-", mutationName, ":", err)
					logger.Error("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
					continue
				}
			}
		}
	}
	// 2.4 save the result of stage1+stage2 into outputPath/DBMS/taskId/result.json, see TaskResult
	taskResult.EndTime = time.Now().String()
	err = taskResult.SaveTaskResult(taskOutputPath)
	if err != nil {
		logger.Error("??????????????????????????????????????????????????")
		logger.Error("[Save Result Error] ", err)
		logger.Error("??????????????????????????????????????????????????")
	}
	// **************************************************
	logger.Info("Finished **************************************************")
	// end 2
	return nil
}
