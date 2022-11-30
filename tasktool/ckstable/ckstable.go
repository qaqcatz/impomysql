package ckstable

import (
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
	"github.com/qaqcatz/nanoshlib"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// CheckStableTask: Some bugs are unstable.
// We will repeat the originalSql/MutatedSql of each bug execNum(recommended 10) times,
// save the stable bugs into directory maystable,
// save the unstable bugs into directory unstable.
func CheckStableTask(config *task.TaskConfig, publicConn *connector.Connector, execNum int) error {
	if execNum <= 0 {
		return errors.New("[CheckStableTask]execNum <= 0")
	}

	// check path
	ddlPath := config.DDLPath
	bugsPath := path.Join(config.GetTaskBugsPath())
	exists, err := pathExists(ddlPath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	exists, err = pathExists(bugsPath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	mayStablePath := path.Join(config.GetTaskPath(), "maystable")
	_ = os.RemoveAll(mayStablePath)
	_ = os.Mkdir(mayStablePath, 0777)
	unStablePath := path.Join(config.GetTaskPath(), "unstable")
	_ = os.RemoveAll(unStablePath)
	_ = os.Mkdir(unStablePath, 0777)

	// create mysql connector
	var conn *connector.Connector = nil
	if publicConn != nil {
		conn = publicConn
	} else {
		conn, err = connector.NewConnector(config.Host, config.Port, config.Username, config.Password, config.DbName)
		if err != nil {
			return err
		}
	}

	// for each bug in bugs, check whether the bug is stable.
	err = conn.InitDBWithDDLPath(ddlPath)
	if err != nil {
		return err
	}
	bugsDir, err := ioutil.ReadDir(bugsPath)
	if err != nil {
		return errors.Wrap(err, "[CheckStableTask]read dir error")
	}
	for _, bugJsonFile := range bugsDir {
		if !strings.HasSuffix(bugJsonFile.Name(), ".json") {
			continue
		}
		bugJsonPath := path.Join(bugsPath, bugJsonFile.Name())
		bugLogPath := bugJsonPath[0:len(bugJsonPath)-5]+".log"
		exists, err := pathExists(bugLogPath)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("[CheckStableTask]miss log: " + bugLogPath)
		}

		ck, err := CheckStable(bugJsonPath, execNum, conn)
		if err != nil {
			return err
		}
		if ck {
			_, errStream, err := nanoshlib.Exec("cp " + bugJsonPath + " " + mayStablePath, -1)
			if err != nil {
				return errors.New("[CheckStableTask]cp stable bug json error: " + errStream)
			}
			_, errStream, err = nanoshlib.Exec("cp " + bugLogPath + " " + mayStablePath, -1)
			if err != nil {
				return errors.New("[CheckStableTask]cp stable bug log error: " + errStream)
			}
		} else {
			_, errStream, err := nanoshlib.Exec("cp " + bugJsonPath + " " + unStablePath, -1)
			if err != nil {
				return errors.New("[CheckStableTask]cp unstable bug json error: " + errStream)
			}
			_, errStream, err = nanoshlib.Exec("cp " + bugLogPath + " " + unStablePath, -1)
			if err != nil {
				return errors.New("[CheckStableTask]cp unstable bug log error: " + errStream)
			}
		}
	}
	return nil
}

func CheckStable(bugJsonPath string, execNum int, conn *connector.Connector) (bool, error) {
	bug, err := task.NewBugReport(bugJsonPath)
	if err != nil {
		return false, err
	}

	var result *connector.Result = nil
	for i := 0; i < execNum; i++ {
		tempResult := conn.ExecSQL(bug.OriginalSql)
		if tempResult.Err != nil {
			return false, nil
		}
		if result == nil {
			result = tempResult
		} else {
			cmp, err := result.CMP(tempResult)
			if err != nil || cmp != 0 {
				return false, nil
			}
		}
	}

	result = nil
	for i := 0; i < execNum; i++ {
		tempResult := conn.ExecSQL(bug.MutatedSql)
		if tempResult.Err != nil {
			return false, nil
		}
		if result == nil {
			result = tempResult
		} else {
			cmp, err := result.CMP(tempResult)
			if err != nil || cmp != 0 {
				return false, nil
			}
		}
	}

	return true, nil
}