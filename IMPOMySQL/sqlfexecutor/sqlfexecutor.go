package sqlfexecutor

import (
	"errors"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/qaqcatz/IMPOMySQL/IMPOMySQL/connector"
	"io/ioutil"
	"strconv"
	"time"
)

// SQLFExecutor: read .sql file(MySQL) file, parse each sql to ast, execute them, get the results
type SQLFExecutor struct {
	ASTs []ast.StmtNode
	Results []*connector.Result
	ReadTime time.Duration // total time of reading .sql file
	ParseTime time.Duration // total time of parsing sqls to ASTs
	ExecuteTime time.Duration // total execute time
	PassedSQLNum int // the number of passed sql
	FailedSQLNum int // the number of failed sql
}

func (sqlFExecutor *SQLFExecutor) ToString() string {
	str := ""
	if sqlFExecutor.ASTs == nil || len(sqlFExecutor.ASTs) == 0 {
		return str + "|ASTs| = 0"
	}
	str += "|ASTs|: " + strconv.Itoa(len(sqlFExecutor.ASTs)) + "\n"
	str += "Read Time: " + sqlFExecutor.ReadTime.String() + "\n"
	str += "Parse Time: " +  sqlFExecutor.ParseTime.String()
	if sqlFExecutor.Results == nil || len(sqlFExecutor.Results) == 0 {
		return str + "\n|Results| = 0"
	}
	str += "\nExec Time: " + sqlFExecutor.ExecuteTime.String() + "\n"
	str += "Passed SQL Num: " + strconv.Itoa(sqlFExecutor.PassedSQLNum) + "\n"
	str += "Failed SQL Num: " + strconv.Itoa(sqlFExecutor.FailedSQLNum)
	for i, result := range sqlFExecutor.Results {
		str += "\n==================================================\n"
		str += "[sql "+strconv.Itoa(i)+"]: " + sqlFExecutor.ASTs[i].Text() + "\n"
		str += "[result "+strconv.Itoa(i)+"]: " + result.ToString() + "\n"
		str += "=================================================="
	}
	return str
}

func NewSQLFExecutor(filePath string) (*SQLFExecutor, error) {
	startTime := time.Now()
	sqls, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("NewSQLFExecutor: read sqls file error: " + err.Error())
	}
	readTime := time.Since(startTime)
	sqlFExecutor, err := NewSQLFExecutorB(sqls)
	if err != nil {
		return nil, errors.New("NewSQLFExecutor: " + err.Error())
	}
	sqlFExecutor.ReadTime = readTime
	return sqlFExecutor, nil
}

func NewSQLFExecutorB(sqls []byte) (*SQLFExecutor, error) {
	startTime := time.Now()
	p := parser.New()
	stmtNodes, _, err := p.Parse(string(sqls), "", "")
	if err != nil {
		return nil, errors.New("NewSQLFExecutorB: parse sqls error: " + err.Error())
	}
	parseTime := time.Since(startTime)

	return &SQLFExecutor{
		ASTs: stmtNodes,
		ParseTime: parseTime,
	}, nil
}

func (sqlFExecutor *SQLFExecutor) Exec(conn *connector.Connector) {
	startTime := time.Now()
	sqlFExecutor.PassedSQLNum = 0
	sqlFExecutor.FailedSQLNum = 0
	for _, AST := range sqlFExecutor.ASTs {
		result := conn.ExecSQL(AST.Text())
		if result.Err != nil {
			sqlFExecutor.FailedSQLNum += 1
		} else {
			sqlFExecutor.PassedSQLNum += 1
		}
		sqlFExecutor.Results = append(sqlFExecutor.Results, result)
	}
	sqlFExecutor.ExecuteTime = time.Since(startTime)
}