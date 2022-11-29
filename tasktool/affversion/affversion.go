package affversion

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/task"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// global lock
var affVersionLock sync.Mutex

// AffVersionTask:
// sqlsim first!
//
// Verify whether the bugs detected by tasks can be reproduced on the specified version of DBMS.
// You need to deploy the specified version of DBMS in config.Host:config.Port yourself.
//
// We will create a sqlite database `affversion.db` under the sibling directory of config.GetTaskPath() with a table:
//   CREATE TABLE IF NOT EXISTS `affversion` (`taskId` INT, `bugJsonName` TEXT, `version` TEXT, `status` INT);
//
// - `taskId`: the id of a task, e.g. 0, 1, 2, ...
//
// - `bugJsonName`: the json file name of a bug, e.g. bug-0-21-FixMHaving1U,
// you can use task-`taskId`/sqlsim/`bugJsonName` to read the bug.
//
// - `version`, `status`: whether the bug can be reproduced on the specified version of DBMS.
// `version` can be an arbitrary non-empty string, it is recommended to use tag or commit id.
// `status`: 1-yes; 0-no; -1-error.
//
// If whereVersionEQ == "", we will verify each bug under config.GetTaskPath()/sqlsim,
//
// else we will only verify these bugs:
//   SELECT `bugJsonName` FROM `affversion`
//   WHERE `taskId` = config.TaskId AND `version` = whereVersionEQ AND `status`=1
//
// According to the reproduction status of the bug, we will insert a new record to `affversion`:
//   INSERT INTO `affversion` (`taskId`, `bugJsonName`, `version`, `status`)
//   SELECT taskId, bugJsonName, version, status
//   WHERE NOT EXISTS
//   (SELECT * from `affversion`
//   WHERE `taskId`=taskId AND `bugJsonName`=bugJsonName AND `version`=version AND `status`=status);
func AffVersionTask(config *task.TaskConfig, publicConn *connector.Connector, version string, whereVersionEQ string) error {
	if version == "" {
		return errors.New("[AffVersionTask]version empty")
	}

	// check path
	ddlPath := config.DDLPath
	sqlSimPath := path.Join(config.GetTaskPath(), "sqlsim")
	exists, err := pathExists(ddlPath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	exists, err = pathExists(sqlSimPath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	// open sqlite database affVersion and create table if not exists affVersion
	taskSibPath := filepath.Join(config.GetTaskPath(), "..")
	affVersionDBPath := path.Join(taskSibPath, "affversion.db")
	affVersionDB, err := sql.Open("sqlite3", affVersionDBPath)
	defer affVersionDB.Close()
	if err != nil {
		return errors.Wrap(err, "[AffVersionTask]open database error")
	}
	_, err = affVersionDB.Exec(`CREATE TABLE IF NOT EXISTS affversion (
    taskId INT, bugJsonName TEXT, 
    version TEXT, status INT);`)
	if err != nil {
		return errors.Wrap(err, "[AffVersionTask]create table error")
	}

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

	var bugJsonNames []string
	if whereVersionEQ == "" {
		//  verify each bug under sqlSimPath
		bugJsonNames, err = getBugsFromDir(sqlSimPath)
		if err != nil {
			return err
		}
	} else {
		// only verify these bugs:
		//   SELECT `bugJsonName` FROM `affversion`
		//   WHERE `taskId` = config.TaskId AND `version` = whereVersionEQ AND `status`=1
		bugJsonNames, err = getBugsFromDB(affVersionDB, config.TaskId, whereVersionEQ)
	}

	if len(bugJsonNames) != 0 {
		err = conn.InitDBWithDDLPath(ddlPath)
		if err != nil {
			return err
		}
		for _, bugJsonName := range bugJsonNames {
			bugJsonPath := path.Join(sqlSimPath, bugJsonName)
			err = doVerify(bugJsonPath, config.TaskId, bugJsonName, version, affVersionDB, conn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getBugsFromDir(bugsPath string) ([]string, error) {
	bugJsonNames := make([]string, 0)

	bugsDir, err := ioutil.ReadDir(bugsPath)
	if err != nil {
		return nil, errors.Wrap(err, "[getBugsFromDir]read dir error")
	}
	for _, bugJsonFile := range bugsDir {
		if !strings.HasSuffix(bugJsonFile.Name(), ".json") {
			continue
		}
		bugJsonNames = append(bugJsonNames, bugJsonFile.Name())
	}
	return bugJsonNames, nil
}

func getBugsFromDB(db *sql.DB, taskId int, whereVersionEQ string) ([]string, error) {
	bugJsonNames := make([]string, 0)

	rows, err := db.Query(`SELECT bugJsonName FROM affversion WHERE 
	taskId = `+strconv.Itoa(taskId)+` AND 
    version = '`+whereVersionEQ+`' AND 
    status=1`)
	if err != nil {
		return nil, errors.Wrap(err, "[getBugsFromDB]select bugs error")
	}
	defer rows.Close()
	for rows.Next() {
		var bugJsonName string
		err = rows.Scan(&bugJsonName)
		if err != nil {
			return nil, errors.Wrap(err, "[getBugsFromDB]scan row error")
		}
		bugJsonNames = append(bugJsonNames, bugJsonName)
	}
	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "[getBugsFromDB]rows err")
	}

	return bugJsonNames, nil
}

func doVerify(bugJsonPath string, taskId int, bugJsonName string,
	version string, affVersionDB *sql.DB,
	conn *connector.Connector) error {

	bug, err := task.NewBugReport(bugJsonPath)
	if err != nil {
		return err
	}
	originalResult := conn.ExecSQL(bug.OriginalSql)
	mutatedResult := conn.ExecSQL(bug.MutatedSql)
	check, err := oracle.Check(originalResult, mutatedResult, bug.IsUpper)

	status := -1
	if err != nil {
		status = -1
	} else {
		if check {
			status = 0
		} else {
			status = 1
		}
	}

	affVersionLock.Lock()

	//   INSERT INTO `affversion` (`taskId`, `bugJsonName`, `version`, `status`)
	//   SELECT taskId, bugJsonName, version, status
	//   WHERE NOT EXISTS
	//   (SELECT * from `affversion`
	//   WHERE `taskId`=taskId AND `bugJsonName`=bugJsonName AND `version`=version AND `status`=status);
	_, err = affVersionDB.Exec(`INSERT INTO affversion (taskId, bugJsonName, version, status)
	  SELECT `+strconv.Itoa(taskId)+`, '`+bugJsonName+`', '`+version+`', `+strconv.Itoa(status)+` 
	  WHERE NOT EXISTS
	  (SELECT * from affversion
	  WHERE taskId=`+strconv.Itoa(taskId)+` AND 
	  bugJsonName='`+bugJsonName+`' AND 
	  version='`+version+`' AND 
	  status=`+strconv.Itoa(status)+`);`)
	if err != nil {
		return errors.Wrap(err, "[doVerify]insert bug error")
	}

	affVersionLock.Unlock()

	return nil
}