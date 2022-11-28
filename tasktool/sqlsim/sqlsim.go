package sqlsim

import (
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
// 1. mkdir sqlsim, read bugs and ddl in task path if exists, create connector
//
// 2. for each bug in bugs, simplify (ddl, bug) and save the result in sqlsim. see SqlSim.
func SqlSimTask(config *task.TaskConfig, publicConn *connector.Connector) error {
	// 1. mkdir sqlsim, read bugs and ddl in task path if exist, create connector
	ddlPath := config.DDLPath
	bugsPath := config.GetTaskBugsPath()
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

	sqlSimPath := path.Join(config.GetTaskPath(), "sqlsim")
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
	err = conn.InitDBWithDDLPath(ddlPath)
	if err != nil {
		return err
	}

	// 2. for each bug in bugs, simplify (ddl, bug) and save the result in sqlsim.
	bugsDir, err := ioutil.ReadDir(bugsPath)
	if err != nil {
		return errors.Wrap(err, "[SqlSimTask]read dir error")
	}
	for _, bugJsonFile := range bugsDir {
		if !strings.HasSuffix(bugJsonFile.Name(), ".json") {
			continue
		}
		bugJsonPath := path.Join(bugsPath, bugJsonFile.Name())
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
	bug.OriginalResult = conn.ExecSQL(bug.OriginalSql)
	if bug.OriginalResult.Err != nil {
		return bug.OriginalResult.Err
	}
	bug.MutatedResult = conn.ExecSQL(bug.MutatedSql)
	if bug.MutatedResult.Err != nil {
		return bug.MutatedResult.Err
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

var SimDMLFuncs = []func(report *task.BugReport, connector2 *connector.Connector) error {
	rmWith,
	rmUnion,
	rmHint,
	rmOrderBy,
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