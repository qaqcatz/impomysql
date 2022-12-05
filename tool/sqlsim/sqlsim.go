package sqlsim

import (
	"fmt"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/pkg/errors"
	"github.com/qaqcatz/impomysql/connector"
	"github.com/qaqcatz/impomysql/task"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// SqlSimTask:
//
// ckstable first!
//
// 1. mkdir sqlsim, read bugs and ddl in task path if exists, create connector
//
// 2. for each bug in bugs, simplify (ddl, bug) and save the result in sqlsim. see SqlSim.
//
// Update: During affversion, we found that some new features cannot run on the old version of DBMS.
// We will try to simplify these new features in sqlsim.
// Actually, we will verify which features in ./resources/impo.yy can not run on mysql 5.0.15 (the oldest version in
// mysql download page: https://downloads.mysql.com/archives/community/), and try to remove them.
//
// There are a lot of functions with prefix rm or frm in sqlsim.
//
// - rm means it is a normal simplified function.
//
// - frm means it is responsible for simplifying new features.
func SqlSimTask(config *task.TaskConfig, publicConn *connector.Connector) error {
	// 1. mkdir sqlsim, read bugs and ddl in task path if exist, create connector
	ddlPath := config.DDLPath
	mayStablePath := path.Join(config.GetTaskPath(), "maystable")
	exists, err := pathExists(ddlPath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	exists, err = pathExists(mayStablePath)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	sqlSimPath := path.Join(config.GetTaskPath(), "sqlsim")
	_ = os.RemoveAll(sqlSimPath)
	_ = os.Mkdir(sqlSimPath, 0777)

	var conn *connector.Connector = nil
	if publicConn != nil {
		conn = publicConn
	} else {
		conn, err = connector.NewConnector(config.Host, config.Port, config.Username, config.Password, config.DbName)
		if err != nil {
			return err
		}
	}

	// 2. for each bug in bugs, simplify (ddl, bug) and save the result in sqlsim.
	err = conn.InitDBWithDDLPath(ddlPath)
	if err != nil {
		return err
	}
	bugsDir, err := ioutil.ReadDir(mayStablePath)
	if err != nil {
		return errors.Wrap(err, "[SqlSimTask]read dir error")
	}
	for _, bugJsonFile := range bugsDir {
		if !strings.HasSuffix(bugJsonFile.Name(), ".json") {
			continue
		}
		bugJsonPath := path.Join(mayStablePath, bugJsonFile.Name())
		err = SqlSim(conn, sqlSimPath, ddlPath, bugJsonPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// SqlSim:
//
// 1. simplify dml: try to remove each node in original/mutated sql,
// simplify if the result does not change or the implication oracle can still detect the bug.
//
// 2. write the simplified ddl and bug(json+log) into sqlsim.
func SqlSim(conn *connector.Connector, outputPath string, ddlPath string, bugJsonPath string) error {
	bug, err := task.NewBugReport(bugJsonPath)
	if err != nil {
		return err
	}
	// originalSql / mutatedSql may error, we should return nil but not error.
	// I already filtered out the error in the task, why? May be a bug of mysql driver...
	bug.OriginalResult = conn.ExecSQL(bug.OriginalSql)
	if bug.OriginalResult.Err != nil {
		//return bug.OriginalResult.Err
		fmt.Println(bugJsonPath + "'s originalSql executes error, but I don't known why...")
		return nil
	}
	bug.MutatedResult = conn.ExecSQL(bug.MutatedSql)
	if bug.MutatedResult.Err != nil {
		//return bug.MutatedResult.Err
		fmt.Println(bugJsonPath + "'s mutatedSql executes error, but I don't known why...")
		return nil
	}

	// 1. simplify dml: try to remove each node in original/mutated sql,
	// simplify if the result does not change or the implication oracle can still detect the bug.
	err = SimDML(bug, conn)
	if err != nil {
		return err
	}

	// 2. write the simplified ddl and bug(json+log) into sqlsim.
	err = bug.SaveBugReport(outputPath)
	if err != nil {
		return err
	}

	return nil
}

// do not adjust the order!
var SimDMLFuncs = []func(report *task.BugReport, connector2 *connector.Connector) error{
	frmWith,
	rmUnion,
	frmHint,
	rmOrderBy,
	rmBinOpTrue,
	rmBinOpFalse,
	frmTimeFunc,
	frmStrFunc,
	frmInfoFunc,
	frmCharset,
}

func SimDML(bug *task.BugReport, conn *connector.Connector) error {
	for _, simDMLFunc := range SimDMLFuncs {
		err := simDMLFunc(bug, conn)
		if err != nil {
			return err
		}
	}
	return nil
}
