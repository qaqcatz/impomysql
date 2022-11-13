package task

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/mutation/stage1"
	"github.com/qaqcatz/impomysql/mutation/stage2"
	"github.com/qaqcatz/nanoshlib"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type TaskConfig struct {
	OutputPath string `json:"outputPath"` // default: ./output
	DBMS       string `json:"dbms"`       // default: mysql
	TaskId     int    `json:"taskId"`     // >= 0
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	DbName     string `json:"dbname"`
	Seed       int64  `json:"seed"`    // <= 0: current time
	DDLPath    string `json:"ddlPath"` // ddl.sql, can not have ';' in comment
	DMLPath    string `json:"dmlPath"` // dml.sql, can not hava ';' in comment
	RdGenPath  string `json:"rdGenPath"` // use go-randgen -Z ZZPath -Y YYPath -Q QueriesNum -B --seed Seed if RdGenPath != ""
	ZZPath     string `json:"zzPath"`
	YYPath     string `json:"yyPath"`
	QueriesNum int    `json:"queriesNum"`
	NeedDML    bool   `json:"needDML"` // output dml.sql or not
}

// TaskConfig.GetTaskPath:
//   path.Join(taskConfig.OutputPath, taskConfig.DBMS, "task-"+strconv.Itoa(taskConfig.TaskId))
func (taskConfig *TaskConfig) GetTaskPath() string {
	return path.Join(taskConfig.OutputPath, taskConfig.DBMS, "task-"+strconv.Itoa(taskConfig.TaskId))
}

// TaskConfig.GetTaskBugsPath:
//   path.Join(taskConfig.OutputPath, taskConfig.DBMS, "task-"+strconv.Itoa(taskConfig.TaskId), "bugs")
func (taskConfig *TaskConfig) GetTaskBugsPath() string {
	return path.Join(taskConfig.OutputPath, taskConfig.DBMS, "task-"+strconv.Itoa(taskConfig.TaskId), "bugs")
}

// TaskConfig.SaveConfig: save config into:
//   path.Join(taskInputPath, "task-"+strconv.Itoa(taskConfig.TaskId)+"-config.json")
// You should init taskInputPath yourself.
func (taskConfig *TaskConfig) SaveConfig(taskInputPath string) error {
	jsonPath := path.Join(taskInputPath, "task-"+strconv.Itoa(taskConfig.TaskId)+"-config.json")
	jsonData, err := json.Marshal(taskConfig)
	if err != nil {
		return errors.Wrap(err, "[TaskConfig.SaveConfig]marshal error")
	}
	err = ioutil.WriteFile(jsonPath, jsonData, 0777)
	if err != nil {
		return errors.Wrap(err, "[TaskConfig.SaveConfig]write json error")
	}
	return nil
}

func NewTaskConfig(configJsonPath string) (*TaskConfig, error) {
	configData, err := ioutil.ReadFile(configJsonPath)
	if err != nil {
		return nil, errors.Wrap(err, "[NewTaskConfig]read task config error")
	}
	var configT TaskConfig
	err = json.Unmarshal(configData, &configT)
	if err != nil {
		return nil, errors.Wrap(err, "[NewTaskConfig]unmarshal task config error")
	}
	config := &configT
	return InitTaskConfig(config)
}

// InitTaskConfig: check input, assign default value, convert path to abs, create config
func InitTaskConfig(config *TaskConfig) (*TaskConfig, error) {
	if config.OutputPath == "" {
		p, err := filepath.Abs("./output")
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskConfig]path abs error")
		}
		config.OutputPath = p
	} else {
		p, err := filepath.Abs(config.OutputPath)
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskConfig]path abs error")
		}
		config.OutputPath = p
	}
	if config.DBMS == "" {
		config.DBMS = "mysql"
	}
	// Host, Port, Username, Password, DbName, MysqlClientPath can be checked by connector
	if config.TaskId < 0 {
		return nil, errors.New("[InitTaskConfig]taskId < 0 ")
	}
	if config.Seed <= 0 {
		config.Seed = time.Now().UnixNano()
	}
	if config.RdGenPath == "" {
		if config.DDLPath == "" {
			return nil, errors.New("[InitTaskConfig]empty ddlPath")
		}

		p, err := filepath.Abs(config.DDLPath)
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskConfig]path abs error")
		}
		config.DDLPath = p

		if config.DMLPath == "" {
			return nil, errors.New("[InitTaskConfig]empty dmlPath")
		}

		p, err = filepath.Abs(config.DMLPath)
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskConfig]path abs error")
		}
		config.DMLPath = p
	} else {
		p, err := filepath.Abs(config.RdGenPath)
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskConfig]path abs error")
		}
		config.RdGenPath = p

		if config.ZZPath == "" {
			return nil, errors.New("[InitTaskConfig]empty ZZPath")
		}

		p, err = filepath.Abs(config.ZZPath)
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskConfig]path abs error")
		}
		config.ZZPath = p

		if config.YYPath == "" {
			return nil, errors.New("[InitTaskConfig]empty YYPath")
		}

		p, err = filepath.Abs(config.YYPath)
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskConfig]path abs error")
		}
		config.YYPath = p

		if config.QueriesNum <= 0 {
			return nil, errors.New("[InitTaskConfig]queriesNum <= 0 ")
		}

		// set TaskConfig.DDLPath = TaskConfig.GetTaskPath()/output.data.sql
		config.DDLPath = path.Join(config.GetTaskPath(), "output.data.sql")
		// set TaskConfig.DMLPath = TaskConfig.GetTaskPath()/output.rand.sql.
		config.DMLPath = path.Join(config.GetTaskPath(), "output.rand.sql")
	}
	return config, nil
}

type TaskResult struct {
	StartTime              string `json:"startTime"`
	DDLSqlsNum             int    `json:"ddlSqlsNum"`
	DMLSqlsNum             int    `json:"dmlSqlsNum"`
	EndInitTime            string `json:"endInitTime"`
	Stage1ErrNum           int    `json:"stage1ErrNum"`
	Stage1ExecErrNum       int    `json:"stage1ExecErrNum"`
	Stage1IgExecErrNum     int    `json:"stage1IgExecErrNum"`
	Stage2ErrNum           int    `json:"stage2ErrNum"`
	Stage2UnitNum          int    `json:"stage2UnitNum"`
	Stage2UnitErrNum       int    `json:"stage2UnitErrNum"`
	Stage2UnitExecErrNum   int    `json:"stage2UnitExecErrNum"`
	Stage2IgUnitExecErrNum int    `json:"stage2IgUnitExecErrNum"`
	ImpoBugsNum            int    `json:"impoBugsNum"`
	SaveBugErrNum          int    `json:"saveBugErrNum"`
	EndTime                string `json:"endTime"`
}

// TaskResult.SaveTaskResult: output to taskPath/result.json
func (taskResult *TaskResult) SaveTaskResult(taskPath string) error {
	jsonPath := path.Join(taskPath, "result.json")
	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		return errors.Wrap(err, "[TaskResult.SaveTaskResult]marshal error")
	}
	err = ioutil.WriteFile(jsonPath, jsonData, 0777)
	if err != nil {
		return errors.Wrap(err, "[TaskResult.SaveTaskResult]write json error")
	}
	return nil
}

// Run: Basic task for finding logical bugs:
//
// 1. init:
//   1.1 init TaskConfig.GetTaskPath()
//   1.2 create logger
//   1.3 create connector(if publicConn == nil, otherwise just use public Conn)
//   1.4 go-randgen if TaskConfig.RdGenPath != ""
//   1.5 read ddl, init database, execute ddl
//   1.6 read dml
//
// 2. run:
// for each dml sql, do:
//   2.1 stage1.InitAndExec
//   2.2 stage2.MutateAllAndExec
//   2.3 use oracle.Check to detect logical bugs
//   2.4 save task result
func RunTask(config *TaskConfig, publicConn *connector.Connector, publicLogger *logrus.Logger) (*TaskResult, error) {
	// 1 init
	startTime := time.Now().String()
	// **************************************************
	// 1.1 init TaskConfig.GetTaskPath()
	_ = os.RemoveAll(config.GetTaskPath())
	_ = os.MkdirAll(config.GetTaskPath(), 0777)
	// 1.2 create logger
	loggerPath := path.Join(config.GetTaskPath(), "task.log")
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	file, err := os.OpenFile(loggerPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return nil, errors.Wrap(err, "[RunTask]create logger error")
	}
	defer file.Close()
	writers := []io.Writer{file}
	multiWriter := io.MultiWriter(writers...)
	logger.SetOutput(multiWriter)
	logger.SetLevel(logrus.InfoLevel)
	// 1.3 create connector(if publicConn == nil, otherwise just use public Conn)
	var conn *connector.Connector = nil
	if publicConn == nil {
		logger.Info("create connector")
		conn, err = connector.NewConnector(config.Host, config.Port, config.Username, config.Password, config.DbName)
		if err != nil {
			logger.Error("create connector error: " + err.Error())
			return nil, err
		}
	} else {
		logger.Info("use public connector")
		conn = publicConn
	}

	// 1.4 go-randgen
	// Note that we already changed config.DDLPath and config.DMLPath in NewTaskConfig().
	if config.RdGenPath != "" {
		// cd TaskConfig.GetTaskPath() &&
		// TaskConfig.RdGenPath gentest
		// -Z TaskConfig.ZZPath
		// -Y TaskConfig.YYPath
		// -Q TaskConfig.QueriesNum
		// --seed TaskConfig.Seed
		// -B
		// output: output.data.sql, output.rand.sql
		randGenCmd := "cd "+config.GetTaskPath()+" && "+config.RdGenPath+" gentest "+
			" -Z "+config.ZZPath+
			" -Y "+config.YYPath+
			" -Q "+strconv.Itoa(config.QueriesNum)+
			" --seed "+strconv.FormatInt(config.Seed, 10)+
			" -B "
		logger.Info(" randgan cmd: ", randGenCmd)
		_, errBuf, err := nanoshlib.Exec(randGenCmd, -1)
		if err != nil {
			errStr := ""
			if errBuf != nil {
				errStr = string(errBuf)
			}
			logger.Error("randgen error: ", err, ": ", errStr)
			return nil, errors.Wrap(err, "randgen error")
		}

		if !config.NeedDML {
			// defer remove TaskConfig.DMLPath
			defer os.Remove(config.DMLPath)
		}
	}

	// 1.5 read ddl, init database, execute ddl
	logger.Info("init ddl")
	ddlData, err := ioutil.ReadFile(config.DDLPath)
	if err != nil {
		logger.Error("read ddl error: " + err.Error())
		return nil, errors.Wrap(err, "[RunTask]read ddl error")
	}
	ddlSqls := ExtractSQL(string(ddlData))
	err = InitDDLSqls(ddlSqls, conn)
	if err != nil {
		logger.Error("init ddl sqls error: " + err.Error())
		return nil, err
	}
	// 1.6 read dml
	logger.Info("init dml")
	dmlData, err := ioutil.ReadFile(config.DMLPath)
	if err != nil {
		logger.Error("read dml error: " + err.Error())
		return nil, errors.Wrap(err, "[RunTask]read dml error")
	}
	dmlSqls := ExtractSQL(string(dmlData))
	// **************************************************
	endInitTime := time.Now().String()
	// end 1

	// 2. run
	logger.Info("Running **************************************************")
	taskResult := &TaskResult{
		StartTime:              startTime,
		DDLSqlsNum:             len(ddlSqls),
		DMLSqlsNum:             len(dmlSqls),
		EndInitTime:            endInitTime,
		Stage1ErrNum:           0,
		Stage1ExecErrNum:       0,
		Stage1IgExecErrNum:     0,
		Stage2ErrNum:           0,
		Stage2UnitNum:          0,
		Stage2UnitErrNum:       0,
		Stage2UnitExecErrNum:   0,
		Stage2IgUnitExecErrNum: 0,
		ImpoBugsNum:            0,
		SaveBugErrNum:          0,
		EndTime:                "",
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
			//logger.Error("--------------------------------------------------")
			//logger.Error("[Stage1 Error]", "(", dmlSql.Id, ")", dmlSql.Sql)
			//logger.Error(stage1Result.Err)
			//logger.Error("--------------------------------------------------")
			continue
		}
		// handle stage1 execute error
		if stage1Result.ExecResult.Err != nil {
			taskResult.Stage1ExecErrNum += 1
			taskResult.Stage1IgExecErrNum += 1 // ignore all errors of stage1
			//logger.Error("==================================================")
			//logger.Error("[Stage1 Exec Error]", "(", dmlSql.Id, ")", stage1Result.InitSql)
			//logger.Error(stage1Result.ExecResult.Err)
			//logger.Error("==================================================")
			continue
		}

		originalSql := stage1Result.InitSql
		originalResult := stage1Result.ExecResult

		// 2.2 stage2.MutateAllAndExec
		stage2Result := stage2.MutateAllAndExec(originalSql, config.Seed+int64(i), conn)
		// handle stage2 error
		if stage2Result.Err != nil {
			taskResult.Stage2ErrNum += 1
			//logger.Error("--------------------------------------------------")
			//logger.Error("[Stage2 Error]", "(", dmlSql.Id, ")", originalSql)
			//logger.Error(stage2Result.Err)
			//logger.Error("--------------------------------------------------")
			continue
		}
		// for each mutation unit
		taskResult.Stage2UnitNum += len(stage2Result.MutateUnits)
		for _, mutateUnit := range stage2Result.MutateUnits {
			// handle stage2 unit error
			if mutateUnit.Err != nil {
				taskResult.Stage2UnitErrNum += 1
				//logger.Error("==================================================")
				//logger.Error("[Stage2 Unit Error]", "(", dmlSql.Id, "-", mutateUnit.Name, ")", mutateUnit.Sql)
				//logger.Error(mutateUnit.Err)
				//logger.Error("==================================================")
				continue
			}
			// handle stage2 unit exec error
			if mutateUnit.ExecResult.Err != nil {
				taskResult.Stage2UnitExecErrNum += 1
				if IgnoreError(mutateUnit.Name, mutateUnit.ExecResult) {
					taskResult.Stage2IgUnitExecErrNum += 1
					continue
				}
				logger.Error("==================================================")
				logger.Error("[Stage2 Unit Exec Error]", "(", dmlSql.Id, "-", mutateUnit.Name, ")", mutateUnit.Sql)
				logger.Error(mutateUnit.ExecResult.Err)
				logger.Error("==================================================")
				continue
			}

			mutationName := mutateUnit.Name
			isUpper := mutateUnit.IsUpper
			mutatedSql := mutateUnit.Sql
			mutatedResult := mutateUnit.ExecResult

			//   2.3 use oracle.Check to detect logical bugs
			if !oracle.Check(originalResult, mutatedResult, isUpper) {
				// logical bug!!!
				bugId := taskResult.ImpoBugsNum
				taskResult.ImpoBugsNum += 1

				logger.Info("logical bug!!! bugId = ", bugId, " sqlId = ", dmlSql.Id, " mutationName = ", mutationName)
				if publicLogger != nil {
					publicLogger.Info("task-", strconv.Itoa(config.TaskId), " detected a logical bug!!! bugId = ",
						bugId, " sqlId = ", dmlSql.Id, " mutationName = ", mutationName)
				}

				bugReport := &BugReport{
					ReportTime:     time.Now().String(),
					BugId:          bugId,
					SqlId:          dmlSql.Id,
					MutationName:   mutationName,
					IsUpper:        isUpper,
					OriginalSql:    originalSql,
					OriginalResult: originalResult,
					MutatedSql:     mutatedSql,
					MutatedResult:  mutatedResult,
				}
				err := bugReport.SaveBugReport(config.GetTaskBugsPath())
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
	// 2.4 save task result
	taskResult.EndTime = time.Now().String()
	err = taskResult.SaveTaskResult(config.GetTaskPath())
	if err != nil {
		logger.Error("??????????????????????????????????????????????????")
		logger.Error("[Save Result Error] ", err)
		logger.Error("??????????????????????????????????????????????????")
	}
	// **************************************************
	logger.Info("Finished **************************************************")
	// end 2
	return taskResult, nil
}
