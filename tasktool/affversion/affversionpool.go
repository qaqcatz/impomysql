package affversion

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

// AffVersionTaskPool: like task and task pool, see AffVersionTask
func AffVersionTaskPool(config *task.TaskPoolConfig, version string, whereVersionEQ string) error {
	// check task pool path
	taskPoolPath := config.GetTaskPoolPath()
	exists, err := pathExists(taskPoolPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("[AffVersionTaskPool]task pool path does not exist")
	}

	// create connectors pool
	connPool, err := connector.NewConnectorPool(config.Host, config.Port, config.Username, config.Password,
		config.DbPrefix, config.ThreadNum)
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
			fmt.Println("[Rate]", i, "/", total)
		} else {
			cur += 1
		}

		// wait for a free connector
		conn := connPool.WaitForFree()
		waitGroup.Add(1)
		go PrepareAndRunAffVersionTask(taskConfigJsonPath, &waitGroup, conn, connPool,
			version, whereVersionEQ)
	}
	waitGroup.Wait()
	fmt.Println("Finished!")

	return nil
}

func PrepareAndRunAffVersionTask(taskConfigJsonPath string,
	waitGroup *sync.WaitGroup,
	conn *connector.Connector, connPool *connector.ConnectorPool,
	version string, whereVersionEQ string) {

	defer func() {
		connPool.BackToPool(conn)
		waitGroup.Done()
	}()

	taskConfig, err := task.NewTaskConfig(taskConfigJsonPath)
	if err != nil {
		panic(fmt.Sprintf("[PrepareAndRunAffVersionTask]new task config error: %+v\n", err))
	}
	err = AffVersionTask(taskConfig, conn, version, whereVersionEQ)
	if err != nil {
		panic(fmt.Sprintf("[PrepareAndRunAffVersionTask]sqlsim task error: %+v\n", err))
	}
}