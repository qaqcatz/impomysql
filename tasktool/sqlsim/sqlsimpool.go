package sqlsim

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

// SqlSimTaskPool: like task and task pool, see SqlSimTask
func SqlSimTaskPool(config *task.TaskPoolConfig) error {
	// check task pool path
	taskPoolPath := config.GetTaskPoolPath()
	exists, err := pathExists(taskPoolPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("[SqlSimTaskPool]task pool path does not exist")
	}

	// create logger
	// 1.2 create logger, write to TaskPoolConfig.GetTaskPoolPath()/taskpool.log and os.Stdout
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	writers := []io.Writer{
		os.Stdout,
	}
	multiWriter := io.MultiWriter(writers...)
	logger.SetOutput(multiWriter)
	logger.SetLevel(logrus.InfoLevel)

	// create connectors pool
	connPool, err := connector.NewConnectorPool(config.Host, config.Port, config.Username, config.Password,
		config.DbPrefix, config.ThreadNum)
	if err != nil {
		return err
	}

	// for each task config json, call SqlSimTask
	taskPoolDir, err := ioutil.ReadDir(taskPoolPath)
	if err != nil {
		return errors.Wrap(err, "[SqlSimTaskPool]read dir error")
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
			logger.Info("[Rate]", i, "/", total)
		} else {
			cur += 1
		}

		// wait for a free connector
		conn := connPool.WaitForFree()
		waitGroup.Add(1)
		go PrepareAndRunSqlSimTask(taskConfigJsonPath, &waitGroup, conn, connPool)
	}
	waitGroup.Wait()
	logger.Info("Finished!")
	return nil
}

func PrepareAndRunSqlSimTask(taskConfigJsonPath string,
	waitGroup *sync.WaitGroup,
	conn *connector.Connector, connPool *connector.ConnectorPool) {

	defer func() {
		connPool.BackToPool(conn)
		waitGroup.Done()
	}()

	taskConfig, err := task.NewTaskConfig(taskConfigJsonPath)
	if err != nil {
		panic(fmt.Sprintf("[PrepareAndRunSqlSimTask]new task config error: %+v\n", err))
	}
	err = SqlSimTask(taskConfig, conn)
	if err != nil {
		panic(fmt.Sprintf("[PrepareAndRunSqlSimTask]sqlsim task error: %+v\n", err))
	}
}