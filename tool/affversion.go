package tool

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/mutation/oracle"
	"github.com/qaqcatz/impomysql/task"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// AffVersion:
// Verify whether the bugs detected by tasks can be reproduced on the specified version of DBMS.
//
// - dbmsOutputPath: the OutputPath of your tasks + '/' + the DBMS of your tasks, for example, ./output/mysql
//
// - version: the specified version of DBMS, needs to be a unique string, it is recommended to use tag or commit id.
//
// - dsn, threadNum: You need to deploy the specified version of DBMS in advance and provide your dsn, format:
//   username^password^host^port^dbPrefix
//   you cannot use '^' in any of username, password, host, port, dbPrefix
// for each thread i, we will create a connector with dsn "username:password@tcp(host:port)/dbPrefix+i"
//
// Before introducing whereVersionEQ, you need to know how AffVersion works:
//
// (1) init affversion.db:
// We will create a sqlite database `affversion.db` in dbmsOutputPath with a table:
//   CREATE TABLE `affversion` (`taskPath` TEXT, `bugJsonName` TEXT, `version` TEXT);
//   CREATE INDEX `versionidx` ON `affversion` (`version`);
// If `affversion.db` does not exist, we will create database `affversion.db` and table `affversion`,
// then traverse each task in dbmsOutputPath, traverse each bug in taskPath/bugs(if exists) and update table `affversion`:
//   INSERT INTO `affversion` VALUES (taskPath, bugJsonName, "");
//
// (2) load bugs group by taskPath:
//   SELECT `taskPath`, `bugJsonName` FROM `affversion` WHERE `version` = whereVersionEQ
// We will save these bugs in a map group by taskPath, so that each group only needs to execute ddl once.
//
// Obviously, If whereVersionEQ="", you will get all bugs.
//
// (3) verify each group in parallel:
// Each group will be assigned a thread.
// We will first init database with ddl.
// Then, for each bug in this group, we will verify whether the bug can be reproduced on the specified version of DBMS.
// If it can be reproduced, we will:
//   INSERT INTO `affversion` (`taskPath`, `bugJsonName`, `version`) SELECT taskPath, bugJsonName, version
//   WHERE NOT EXISTS
//   (SELECT * from `affversion` WHERE `taskPath`=taskPath AND `bugJsonName`=bugJsonName AND `version`=version);
// This is done to ensure that each row is unique. (We will also ensure thread safety)
//
// Now you understand how AffVersion works, you can query the table `affversion` to get the information you want.
func AffVersion(dbmsOutputPath string, version string, dsn string, threadNum int, whereVersionEQ string) error {
	// get abs path
	dbmsOutputPath, err := filepath.Abs(dbmsOutputPath)
	if err != nil {
		return errors.Wrap(err, "[AffVersion]path abs error")
	}
	// create connectors
	dsnUnits := strings.Split(dsn, "^")
	if len(dsnUnits) != 5 {
		return errors.New("[AffVersion]len(dsnUnits) != 5: "+dsn)
	}
	username := dsnUnits[0]
	password := dsnUnits[1]
	host := dsnUnits[2]
	portStr := dsnUnits[3]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.Wrap(err, "[AffVersion]parse port error")
	}
	dbPrefix := dsnUnits[4]
	threadPool := make(chan *connector.Connector, threadNum)
	for i := 0; i < threadNum; i++ {
		conn, err := connector.NewConnector(host, port, username, password, dbPrefix+strconv.Itoa(i))
		if err != nil {
			return err
		}
		threadPool <- conn
	}
	// (1) init affversion.db:
	affVersionDBPath := path.Join(dbmsOutputPath, "affversion.db")
	affVersionDBPathExists, err := pathExists(affVersionDBPath)
	if err != nil {
		return err
	}
	// sql.Open will create database if not exists.
	affVersionDB, err := sql.Open("sqlite3", affVersionDBPath)
	defer affVersionDB.Close()
	if err != nil {
		return errors.Wrap(err, "[AffVersion]open database error")
	}
	// if it is the first time to open db, create table and insert data
	if !affVersionDBPathExists {
		_, err = affVersionDB.Exec(`CREATE TABLE affversion (taskPath TEXT, bugJsonName TEXT, version TEXT);`)
		if err != nil {
			return errors.Wrap(err, "[AffVersion]create table error")
		}
		_, err = affVersionDB.Exec(`CREATE INDEX versionidx ON affversion (version);`)
		if err != nil {
			return errors.Wrap(err, "[AffVersion]create index error")
		}
		bugJsonPaths, err := getAllBugsFromDir(dbmsOutputPath);
		if err != nil {
			return err
		}
		for _, bugJsonPath := range bugJsonPaths {
			taskPath := filepath.Join(bugJsonPath, "../", "../")
			bugJsonName := filepath.Base(bugJsonPath)
			_, err = affVersionDB.Exec(`INSERT INTO affversion VALUES ('`+taskPath+`', '`+bugJsonName+`', '');`)
			if err != nil {
				return errors.Wrap(err, "[AffVersion]insert bug error")
			}
		}
	}
	// (2) load bugs group by taskPath:
	bugGroups, err := getBugGroupsFromDB(affVersionDB, whereVersionEQ)
	if err != nil {
		return err
	}
	// (3) verify each group in parallel
	var waitgroup sync.WaitGroup
	var mutex sync.Mutex
	for taskPath, bugGroup := range bugGroups {
		// wait for a free connector
		conn := <- threadPool
		waitgroup.Add(1)
		go doVerify(version, affVersionDB, &waitgroup, &mutex, conn, threadPool, taskPath, bugGroup)
	}
	waitgroup.Wait()
	return nil
}

// getAllBugsFromDir: recursively traverse each bug in dbmsOutputPath, get bugJsonPaths
func getAllBugsFromDir(dbmsOutputPath string) ([]string, error) {
	bugJsonPaths := make([]string, 0)

	dbmsOutputDir, err := ioutil.ReadDir(dbmsOutputPath)
	if err != nil {
		return nil, errors.Wrap(err, "[getAllBugsFromDir]read dir error")
	}
	for _, taskDir := range dbmsOutputDir {
		if !taskDir.IsDir() {
			continue
		}
		bugsPath := path.Join(dbmsOutputPath, taskDir.Name(), "bugs")
		bugsPathExists, err := pathExists(bugsPath)
		if err != nil {
			return nil, err
		}
		if !bugsPathExists {
			continue
		}

		bugsDir, err := ioutil.ReadDir(bugsPath)
		if err != nil {
			return nil, errors.Wrap(err, "[getAllBugsFromDir]read dir error")
		}
		for _, bugJsonFile := range bugsDir {
			if !strings.HasSuffix(bugJsonFile.Name(), ".json") {
				continue
			}
			bugJsonPaths = append(bugJsonPaths, path.Join(bugsPath, bugJsonFile.Name()))
		}
	}
	return bugJsonPaths, nil
}

// getBugGroupsFromDB: SELECT `taskPath`, `bugJsonName` FROM `affversion` WHERE `version` = whereVersionEQ,
// group by taskPath
func getBugGroupsFromDB(db *sql.DB, whereVersionEQ string) (map[string][]string, error) {
	bugGroups := make(map[string][]string)

	rows, err := db.Query("SELECT taskPath, bugJsonName FROM affversion WHERE version='"+whereVersionEQ+"'")
	if err != nil {
		return nil, errors.Wrap(err, "[getBugGroupsFromDB]select bug error")
	}
	defer rows.Close()
	for rows.Next() {
		var taskPath string
		var bugJsonName string
		err = rows.Scan(&taskPath, &bugJsonName)
		if err != nil {
			return nil, errors.Wrap(err, "[getBugGroupsFromDB]scan row error")
		}
		if bugGroup, ok := bugGroups[taskPath]; ok {
			bugGroups[taskPath] = append(bugGroup, bugJsonName)
		} else {
			bugGroups[taskPath] = []string{bugJsonName}
		}
	}
	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "[getBugGroupsFromDB]rows err")
	}

	return bugGroups, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, errors.Wrap(err, "[PathExists]file stat error")
}

func doVerify(version string, affVersionDB *sql.DB,
	waitGroup *sync.WaitGroup, mutex *sync.Mutex,
	conn *connector.Connector, threadPool chan *connector.Connector,
	taskPath string, bugGroup []string) {

	defer func() {
		threadPool <- conn
		waitGroup.Done()
	}()

	// init database with ddl
	err := conn.InitDBWithDDLPath(path.Join(taskPath, "output.data.sql"))
	if err != nil {
		panic(err)
	}
	for _, bugJsonName := range bugGroup {
		bug, err := task.NewBugReport(path.Join(taskPath, "bugs", bugJsonName))
		if err != nil {
			panic(err)
		}
		originalResult := conn.ExecSQL(bug.OriginalSql)
		mutatedResult := conn.ExecSQL(bug.MutatedSql)
		if oracle.Check(originalResult, mutatedResult, bug.IsUpper) {
			continue
		}
		mutex.Lock()
		// INSERT INTO `affversion` (`taskPath`, `bugJsonName`, `version`) SELECT taskPath, bugJsonName, version
		// WHERE NOT EXISTS
		// (SELECT * from `affversion` WHERE `taskPath`=taskPath AND `bugJsonName`=bugJsonName AND `version`=version);
		_, err = affVersionDB.Exec(`INSERT INTO affversion (taskPath, bugJsonName, version) `+
			`SELECT '`+taskPath+`', '`+bugJsonName+`', '`+version+`' WHERE NOT EXISTS `+
			`(SELECT * from affversion WHERE taskPath='`+taskPath+`' AND bugJsonName='`+bugJsonName+`' AND version='`+version+`');`)
		mutex.Unlock()
		if err != nil {
			panic("[doVerify]insert bug error: " + err.Error())
		}
	}
}