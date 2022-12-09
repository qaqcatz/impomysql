package affversion

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/nanoshlib"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

// AffVersionTaskPool: like task and task pool, see AffVersionTask.
// Old versions may crash or exception, we need to save logs for debugging.
//   logPath: taskPoolPath/affversion-version.log
func AffVersionTaskPool(config *task.TaskPoolConfig, threadNum int, port int, version string, whereVersionEQ string) error {
	// check task pool path
	taskPoolPath := config.GetTaskPoolPath()
	exists, err := pathExists(taskPoolPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("[AffVersionTaskPool]task pool path does not exist")
	}

	// create logger
	loggerPath := path.Join(taskPoolPath, "affversion-"+version+".log")
	_ = os.Remove(loggerPath)
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	file, err := os.OpenFile(loggerPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return errors.Wrap(err, "[AffVersionTaskPool]create logger error")
	}
	defer file.Close()
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	multiWriter := io.MultiWriter(writers...)
	logger.SetOutput(multiWriter)
	logger.SetLevel(logrus.InfoLevel)

	// create connectors pool
	connPool, err := connector.NewConnectorPool(config.Host, port, config.Username, config.Password,
		config.DbPrefix, threadNum)
	if err != nil {
		return err
	}

	// for each task config json, call SqlSimTask
	taskPoolDir, err := ioutil.ReadDir(taskPoolPath)
	if err != nil {
		return errors.Wrap(err, "[AffVersionTaskPool]read dir error")
	}
	taskConfigJsonPaths := make([]string, 0)
	for _, taskConfigJsonFile := range taskPoolDir {
		if !strings.HasSuffix(taskConfigJsonFile.Name(), ".json") {
			continue
		}
		if !strings.HasPrefix(taskConfigJsonFile.Name(), "task-") {
			continue
		}

		taskConfigJsonPath := path.Join(taskPoolPath, taskConfigJsonFile.Name())
		taskConfigJsonPaths = append(taskConfigJsonPaths, taskConfigJsonPath)
	}

	var waitGroup sync.WaitGroup
	total := len(taskConfigJsonPaths)
	cur := 0

	for i, taskConfigJsonPath := range taskConfigJsonPaths {

		// rate
		if cur > total/20 {
			cur = 0
			fds, err := monitorFds()
			fdNumStr := ""
			if err != nil {
				fdNumStr = err.Error()
			} else {
				fdNumStr = strconv.Itoa(len(fds))
			}
			logger.Info("[Rate]", i, "/", total, "(fd num: ", fdNumStr, ")")
		} else {
			cur += 1
		}

		// wait for a free connector
		conn := connPool.WaitForFree()
		waitGroup.Add(1)
		go PrepareAndRunAffVersionTask(logger, taskConfigJsonPath, &waitGroup, conn, connPool,
			port, version, whereVersionEQ)
	}
	waitGroup.Wait()
	logger.Info("Finished!")
	//logger.Info("debug fds:")
	//fds, err := monitorFds()
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	for _, fd := range fds {
	//		fmt.Println(fd)
	//	}
	//}
	return nil
}

func PrepareAndRunAffVersionTask(logger *logrus.Logger, taskConfigJsonPath string,
	waitGroup *sync.WaitGroup,
	conn *connector.Connector, connPool *connector.ConnectorPool,
	port int, version string, whereVersionEQ string) {

	defer func() {
		connPool.BackToPool(conn)
		waitGroup.Done()
	}()

	// task may fail due to dbms crash or exception, do not use panic here! just log the error
	taskConfig, err := task.NewTaskConfig(taskConfigJsonPath)
	if err != nil {
		logger.Error("[PrepareAndRunAffVersionTask]new task config error: ", err)
		return
	}
	err = AffVersionTask(taskConfig, conn, port, version, whereVersionEQ)
	if err != nil {
		logger.Error("[PrepareAndRunAffVersionTask]affversion task "+strconv.Itoa(taskConfig.TaskId)+" error: ", err)
		return
	}
}

func monitorFds() ([]string, error) {
	outStream, errStream, err := nanoshlib.Exec(fmt.Sprintf("ls -l /proc/%v/fd", os.Getpid()), -1)
	if err != nil {
		return nil, errors.New("[monitorFd]count fd error: " + err.Error() + ": " + errStream)
	}
	lines := strings.Split(strings.TrimSpace(string(outStream)), "\n")
	return lines, nil
}