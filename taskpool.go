package main

import (
	"encoding/json"
	"errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/nanoshlib"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"
)

type TaskPoolConfig struct {
	OutputPath      string `json:"outputPath"` // default: ./output
	DBMS            string `json:"dbms"`       // default: mysql
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	DbPrefix        string `json:"dbPrefix"`
	GoRandGenPath   string `json:"goRandGenPath"`
	ZZPath          string `json:"zzPath"`
	YYPath          string `json:"yyPath"`
	QueriesNum      int    `json:"queriesNum"`
	Seed            int64  `json:"seed"` // <= 0: current time
	ThreadNum       int    `json:"threadNum"` // default: 16
	MaxTasks        int    `json:"maxTasks"`  // <= 0: no limit
	MaxTimeS        int    `json:"maxTimeS"`  // <= 0: no limit
}

// TaskPoolInputCheck: check input, assign default value
func TaskPoolInputCheck(config *TaskPoolConfig) error {
	if config.OutputPath == "" {
		config.OutputPath = "./output"
	}
	if config.DBMS == "" {
		config.DBMS = "mysql"
	}
	if config.GoRandGenPath == "" {
		return errors.New("TaskPoolInputCheck: empty goRandGenPath")
	}
	if config.ZZPath == "" {
		return errors.New("TaskPoolInputCheck: empty zzPath")
	}
	if config.YYPath == "" {
		return errors.New("TaskPoolInputCheck: empty yyPath")
	}
	if config.QueriesNum <= 0 {
		return errors.New("TaskPoolInputCheck: queriesNum <= 0")
	}
	if config.Seed <= 0 {
		config.Seed = time.Now().UnixNano()
	}
	if config.ThreadNum <= 0 {
		return errors.New("TaskPoolInputCheck: threadNum <= 0")
	}
	return nil
}

func (taskPoolConfig *TaskPoolConfig) GetOutputPath() string {
	return path.Join(taskPoolConfig.OutputPath, taskPoolConfig.DBMS)
}

// RunTaskPool:
//
// 0. check input
//
// 1. init
//   1.1 init OutputPath/DBMS
//   1.2 init logger, write to OutputPath/DBMS/taskpool.log and os.Stdout
//   1.3 create thread pool with size ThreadNum, fill with *connector.Connector,
//   the database name of each connector is config.DbPrefix + thread id
// 2. run, use thread pool(with size ThreadNum) to continuously execute tasks.
// Each thread can only perform one task at the same time. see PrepareAndRunTask
func RunTaskPool(config *TaskPoolConfig) error {
	// 0. check input
	err := TaskPoolInputCheck(config)
	if err != nil {
		return errors.New("RunTaskPool: check input error: " + err.Error())
	}

	// 1. init
	// **************************************************
	// 1.1 init OutputPath/DBMS
	outputPath := config.GetOutputPath()
	_ = os.RemoveAll(outputPath)
	_ = os.MkdirAll(outputPath, 0777)
	// 1.2 init logger, write to OutputPath/DBMS/taskpool.log and os.Stdout
	logPath := path.Join(outputPath, "taskpool.log")
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return errors.New("RunTaskPool: create logger error: " + err.Error())
	}
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
			config.DbPrefix + strconv.Itoa(i))
		if err != nil {
			logger.Error("create connector error: " + err.Error())
			return errors.New("RunTaskPool: create connector error: " + err.Error())
		}
		threadPool <- conn
	}
	// **************************************************
	// end 1

	// 2. run, use thread pool(with size ThreadNum) to continuously execute tasks.
	// Each thread can only perform one task at the same time. see PrepareAndRunTask
	startTime := time.Now()
	totalTaskNum := 0
	var finTaskNum int32 = 0
	var errTaskNum int32 = 0
	logger.Info("Running **************************************************")
	// **************************************************
	for {

		// wait for a free connector
		conn := <- threadPool

		// max time limit
		if config.MaxTimeS > 0 && time.Since(startTime) >= time.Duration(config.MaxTimeS)*time.Second {
			logger.Info("max time!")
			break
		}
		// max task limit
		if config.MaxTimeS > 0 && atomic.LoadInt32(&finTaskNum) >= int32(config.MaxTasks) {
			logger.Info("max tasks!")
			break
		}

		// execute a new task
		taskId := totalTaskNum
		totalTaskNum += 1
		go PrepareAndRunTask(conn, threadPool, &finTaskNum, &errTaskNum, taskId,
			config, logger)
	}
	// **************************************************
	logger.Info("[total time] ", time.Since(startTime).String())
	logger.Info("[total number of tasks] ", totalTaskNum)
	logger.Info("[total number of finished tasks] ", finTaskNum)
	logger.Info("[total number of error tasks] ", errTaskNum)
	logger.Info("Finished **************************************************")
	// end 2
	return nil
}

// PrepareAndRunTask:
//   1. create task config, create task dir, write task config into task dir
//   2. go-randgen write output.data.sql + output.rand.sql into task dir
//   3. run task
func PrepareAndRunTask(conn *connector.Connector, threadPool chan *connector.Connector, finTaskNum *int32, errTaskNum *int32, taskId int,
	config *TaskPoolConfig, logger *logrus.Logger) {

	defer func () {
		atomic.AddInt32(finTaskNum, 1)
		threadPool <- conn
	} ()

	logger.Info("Run task", taskId)

	// 1. create task config, create task dir, write task config into task dir
	taskConfig := &TaskConfig{
		OutputPath: config.OutputPath,
		DBMS: config.DBMS,
		TaskId: taskId,
		Host: config.Host,
		Port: config.Port,
		Username: config.Username,
		Password: config.Password,
		DbName: conn.DbName,
		DDLPath: "",
		DMLPath: "",
		Seed: config.Seed + int64(taskId),
	}
	taskDir := taskConfig.GetOutputPath()
	taskConfig.DDLPath = path.Join(taskDir, "output.data.sql")
	taskConfig.DMLPath = path.Join(taskDir, "output.rand.sql")

	_ = os.MkdirAll(taskDir, 0777)

	taskConfigJsonPath := path.Join(taskDir, "config.json")
	taskConfigJsonData, err := json.Marshal(taskConfig)
	if err != nil {
		atomic.AddInt32(errTaskNum, 1)
		logger.Error("task", taskId, " json marshal error: ", err)
		return
	}
	err = ioutil.WriteFile(taskConfigJsonPath, taskConfigJsonData, 0777)
	if err != nil {
		atomic.AddInt32(errTaskNum, 1)
		logger.Error("task", taskId, " write config json error: ", err)
		return
	}
	// 2. go-randgen output.data.sql + output.rand.sql
	// cd taskDir && ./go-randgen gentest -Z zzPath -Y yyPath -Q queriesNum -B --seed seed
	goRandGenAbsPath, err := filepath.Abs(config.GoRandGenPath)
	if err != nil {
		atomic.AddInt32(errTaskNum, 1)
		logger.Error("task", taskId, " get go-randgen abs path error: ", err)
		return
	}
	zzAbsPath, err := filepath.Abs(config.ZZPath)
	if err != nil {
		atomic.AddInt32(errTaskNum, 1)
		logger.Error("task", taskId, " get zz abs path error: ", err)
		return
	}
	yyAbsPath, err := filepath.Abs(config.YYPath)
	if err != nil {
		atomic.AddInt32(errTaskNum, 1)
		logger.Error("task", taskId, " get yy abs path error: ", err)
		return
	}
	randGenCmd := "cd "+taskDir+" && "+goRandGenAbsPath+" gentest "+
		" -Z "+zzAbsPath+" -Y "+yyAbsPath+
		" -Q "+strconv.Itoa(config.QueriesNum)+" --seed "+strconv.FormatInt(taskConfig.Seed, 10)+
		" -B "
	logger.Info("task", taskId, " randgan cmd: ", randGenCmd)
	_, errBuf, err := nanoshlib.Exec(randGenCmd, -1)
	if err != nil {
		atomic.AddInt32(errTaskNum, 1)
		errStr := ""
		if errBuf != nil {
			errStr = string(errBuf)
		}
		logger.Error("task", taskId, " randgen error: ", err, ": ", errStr)
		return
	}
	// 3. run task
	err = RunTask(taskConfig, conn)
	if err != nil {
		atomic.AddInt32(errTaskNum, 1)
		logger.Error("task", taskId, " run task error: ", err)
		return
	}

	logger.Info("task", taskId, " Finished")
}