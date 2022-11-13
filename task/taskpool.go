package task

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type TaskPoolConfig struct {
	OutputPath  string `json:"outputPath"` // default: ./output
	DBMS        string `json:"dbms"`       // default: mysql
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DbPrefix    string `json:"dbPrefix"`
	Seed        int64  `json:"seed"` // <= 0: current time
	RandGenPath string `json:"randGenPath"`
	ZZPath      string `json:"zzPath"`
	YYPath      string `json:"yyPath"`
	QueriesNum  int    `json:"queriesNum"`
	ThreadNum   int    `json:"threadNum"`
	MaxTasks    int    `json:"maxTasks"` // <= 0: no limit
	MaxTimeS    int    `json:"maxTimeS"` // <= 0: no limit
}

// TaskPoolConfig.GetTaskPoolPath:
//   path.Join(taskPoolConfig.OutputPath, taskPoolConfig.DBMS)
func (taskPoolConfig *TaskPoolConfig) GetTaskPoolPath() string {
	return path.Join(taskPoolConfig.OutputPath, taskPoolConfig.DBMS)
}

func NewTaskPoolConfig(configJsonPath string) (*TaskPoolConfig, error) {
	configData, err := ioutil.ReadFile(configJsonPath)
	if err != nil {
		return nil, errors.Wrap(err, "[NewTaskPoolConfig]read task config error")
	}
	var configT TaskPoolConfig
	err = json.Unmarshal(configData, &configT)
	if err != nil {
		return nil, errors.Wrap(err, "[NewTaskPoolConfig]unmarshal task config error")
	}
	config := &configT
	return InitTaskPoolConfig(config)
}

// InitTaskPoolConfig: check input, assign default value, convert path to abs, create config
func InitTaskPoolConfig(config *TaskPoolConfig) (*TaskPoolConfig, error) {
	if config.OutputPath == "" {
		p, err := filepath.Abs("./output")
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskPoolConfig]path abs error")
		}
		config.OutputPath = p
	} else {
		p, err := filepath.Abs(config.OutputPath)
		if err != nil {
			return nil, errors.Wrap(err, "[InitTaskPoolConfig]path abs error")
		}
		config.OutputPath = p
	}
	if config.DBMS == "" {
		config.DBMS = "mysql"
	}
	if config.Seed <= 0 {
		config.Seed = time.Now().UnixNano()
	}
	if config.RandGenPath == "" {
		return nil, errors.New("[InitTaskPoolConfig]empty randGenPath")
	}
	if config.ZZPath == "" {
		return nil, errors.New("[InitTaskPoolConfig]empty zzPath")
	}
	if config.YYPath == "" {
		return nil, errors.New("[InitTaskPoolConfig]empty yyPath")
	}
	if config.QueriesNum <= 0 {
		return nil, errors.New("[InitTaskPoolConfig]queriesNum <= 0")
	}
	if config.ThreadNum <= 0 {
		return nil, errors.New("[InitTaskPoolConfig]threadNum <= 0")
	}
	return config, nil
}

type TaskPoolResult struct {
	lock              sync.Mutex
	StartTime         string `json:"startTime"`
	TotalTaskNum      int    `json:"totalTaskNum"`
	FinishedTaskNum   int    `json:"finishedTaskNum"`
	ErrorTaskNum      int    `json:"errorTaskNum"`
	ErrorTaskIds      []int  `json:"errorTaskIds"`
	Stage1WarnNum     int    `json:"stage1WarnNum"`
	Stage1WarnTaskIds []int  `json:"stage1WarnTaskIds"`
	Stage2WarnNum     int    `json:"stage2WarnNum"`
	Stage2WarnTaskIds []int  `json:"stage2WarnTaskIds"`
	BugsNum           int    `json:"bugsNum"`
	BugTaskIds        []int  `json:"bugTaskIds"`
	EndTime           string `json:"endTime"`
}

// TaskPoolResult.SaveTaskPoolResult: output to taskPoolPath/result.json
func (taskPoolResult *TaskPoolResult) SaveTaskPoolResult(taskPoolPath string) error {
	taskPoolResult.lock.Lock()
	defer taskPoolResult.lock.Unlock()
	jsonPath := path.Join(taskPoolPath, "result.json")
	jsonData, err := json.Marshal(taskPoolResult)
	if err != nil {
		return errors.Wrap(err, "[TaskPoolResult.SaveTaskPoolResult]marshal error")
	}
	err = ioutil.WriteFile(jsonPath, jsonData, 0777)
	if err != nil {
		return errors.Wrap(err, "[TaskPoolResult.SaveTaskPoolResult]write json error")
	}
	return nil
}

// RunTaskPool:
//
// 1. init
//   1.1 init TaskPoolConfig.GetTaskPoolPath()
//   1.2 create logger, write to TaskPoolConfig.GetTaskPoolPath()/taskpool.log and os.Stdout
//   1.3 create thread pool with size ThreadNum, fill with *connector.Connector,
//   the database name of each connector is config.DbPrefix + thread id
// 2. run, use thread pool(with size ThreadNum) to continuously execute tasks.
// Each thread can only perform one task at the same time, see PrepareAndRunTask().
// Save taskpool result.
func RunTaskPool(config *TaskPoolConfig) (*TaskPoolResult, error) {
	// 1. init
	startTime := time.Now()
	// **************************************************
	// 1.1 init TaskPoolConfig.GetTaskPoolPath()
	_ = os.RemoveAll(config.GetTaskPoolPath())
	_ = os.MkdirAll(config.GetTaskPoolPath(), 0777)
	// 1.2 create logger, write to TaskPoolConfig.GetTaskPoolPath()/taskpool.log and os.Stdout
	loggerPath := path.Join(config.GetTaskPoolPath(), "taskpool.log")
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	file, err := os.OpenFile(loggerPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return nil, errors.Wrap(err, "[RunTaskPool]create logger error")
	}
	defer file.Close()
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	multiWriter := io.MultiWriter(writers...)
	logger.SetOutput(multiWriter)
	logger.SetLevel(logrus.InfoLevel)
	// 1.3 create thread pool with size ThreadNum, fill with *connector.Connector,
	// the database name of each connector is config.DbPrefix + thread id
	threadPool := make(chan *connector.Connector, config.ThreadNum)
	for i := 0; i < config.ThreadNum; i++ {
		conn, err := connector.NewConnector(config.Host, config.Port, config.Username, config.Password,
			config.DbPrefix+strconv.Itoa(i))
		if err != nil {
			logger.Error("create connector error: " + err.Error())
			return nil, err
		}
		threadPool <- conn
	}
	// **************************************************
	// end 1

	// 2. run, use thread pool(with size ThreadNum) to continuously execute tasks.
	// Each thread can only perform one task at the same time, see PrepareAndRunTask().
	// Save taskpool result.
	logger.Info("Running **************************************************")
	taskPoolResult := &TaskPoolResult{
		StartTime:         startTime.String(),
		TotalTaskNum:      0,
		FinishedTaskNum:   0,
		ErrorTaskNum:      0,
		ErrorTaskIds:      make([]int, 0),
		Stage1WarnNum:     0,
		Stage1WarnTaskIds: make([]int, 0),
		Stage2WarnNum:     0,
		Stage2WarnTaskIds: make([]int, 0),
		BugsNum:           0,
		BugTaskIds:        make([]int, 0),
		EndTime:           "",
	}
	// **************************************************
	for {

		// wait for a free connector
		conn := <-threadPool

		// max time limit
		if config.MaxTimeS > 0 && time.Since(startTime) >= time.Duration(config.MaxTimeS)*time.Second {
			logger.Info("max time!")
			break
		}
		// max task limit
		finishedTaskNum := 0
		taskPoolResult.lock.Lock()
		finishedTaskNum = taskPoolResult.FinishedTaskNum
		taskPoolResult.lock.Unlock()
		if config.MaxTasks > 0 && finishedTaskNum >= config.MaxTasks {
			logger.Info("max tasks!")
			break
		}

		// execute a new task
		taskId := taskPoolResult.TotalTaskNum
		taskPoolResult.TotalTaskNum += 1
		go PrepareAndRunTask(config, logger, threadPool, conn, taskPoolResult, taskId)
	}
	// save taskpool result.
	taskPoolResult.EndTime = time.Now().String()
	err = taskPoolResult.SaveTaskPoolResult(config.GetTaskPoolPath())
	if err != nil {
		logger.Error("??????????????????????????????????????????????????")
		logger.Error("[Save Result Error] ", err)
		logger.Error("??????????????????????????????????????????????????")
	}
	// **************************************************
	logger.Info("Finished **************************************************")
	// end 2
	return taskPoolResult, nil
}

// PrepareAndRunTask:
//   1. create and save task config
//   2. run task
func PrepareAndRunTask(config *TaskPoolConfig, logger *logrus.Logger, threadPool chan *connector.Connector,
	conn *connector.Connector, taskPoolResult *TaskPoolResult, taskId int) {

	defer func() {
		taskPoolResult.lock.Lock()
		taskPoolResult.FinishedTaskNum += 1
		taskPoolResult.lock.Unlock()
		threadPool <- conn
	}()

	logger.Info("Run task", taskId)

	// 1. create and save task config
	taskConfig := &TaskConfig{
		OutputPath: config.OutputPath,
		DBMS:       config.DBMS,
		TaskId:     taskId,
		Host:       conn.Host,
		Port:       conn.Port,
		Username:   conn.Username,
		Password:   conn.Password,
		DbName:     conn.DbName,
		Seed:       config.Seed + int64(taskId),
		RdGenPath:  config.RandGenPath,
		ZZPath:     config.ZZPath,
		YYPath:     config.YYPath,
		QueriesNum: config.QueriesNum,
		NeedDML:    false,
	}
	taskConfig, err := InitTaskConfig(taskConfig)
	if err != nil {
		logger.Error("task", taskId, " init task config error: ", err)
		return
	}
	err = taskConfig.SaveConfig(config.GetTaskPoolPath())
	if err != nil {
		logger.Error("task", taskId, " save task config error: ", err)
		return
	}

	// 2. run task
	taskResult, err := RunTask(taskConfig, conn, logger)
	if err != nil {
		taskPoolResult.lock.Lock()
		taskPoolResult.ErrorTaskNum += 1
		taskPoolResult.ErrorTaskIds = append(taskPoolResult.ErrorTaskIds, taskId)
		taskPoolResult.lock.Unlock()
		logger.Error("task", taskId, " run task error: ", err)
		return
	}
	taskPoolResult.lock.Lock()
	stage1WarnNum := taskResult.Stage1ExecErrNum-taskResult.Stage1IgExecErrNum
	if stage1WarnNum > 0 {
		taskPoolResult.Stage1WarnNum += stage1WarnNum
		taskPoolResult.Stage1WarnTaskIds = append(taskPoolResult.Stage1WarnTaskIds, taskId)
	}
	stage2WarnNum := taskResult.Stage2UnitExecErrNum-taskResult.Stage2IgUnitExecErrNum
	if stage2WarnNum > 0 {
		taskPoolResult.Stage2WarnNum += stage2WarnNum
		taskPoolResult.Stage2WarnTaskIds = append(taskPoolResult.Stage2WarnTaskIds, taskId)
	}
	bugsNum := taskResult.ImpoBugsNum
	if bugsNum > 0 {
		taskPoolResult.BugsNum += bugsNum
		taskPoolResult.BugTaskIds = append(taskPoolResult.BugTaskIds, taskId)
	}
	taskPoolResult.lock.Unlock()

	logger.Info("task", taskId, " Finished")
}
