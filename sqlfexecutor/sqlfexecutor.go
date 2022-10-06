package sqlfexecutor

import (
	"errors"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/qaqcatz/impomysql/connector"
	"io/ioutil"
	"strconv"
	"time"
)

// ParseError: parse error message
type ParseError struct {
	Id int // origin index(split by ';')
	Sql string // error sql statement
	Err error // parse error message
}

func (parseError *ParseError) ToString() string {
	return "[parse error "+strconv.Itoa(parseError.Id)+"] "+parseError.Sql+"\n[error message] " + parseError.Err.Error()
}

// SQLFExecutor: Read .sql file(MySQL) file, parse each sql to ast, execute them, get the results
type SQLFExecutor struct {
	ParseErrs []*ParseError // some sql statements may have syntax error
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
	if sqlFExecutor.ParseErrs != nil && len(sqlFExecutor.ParseErrs) != 0 {
		str += "|Parse error|: " + strconv.Itoa(len(sqlFExecutor.ParseErrs)) + "\n"
		str += "========================================\n"
		for _, parseError := range sqlFExecutor.ParseErrs {
			str += parseError.ToString() + "\n"
		}
		str += "========================================\n"
	}
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

// ExtractSQL: Extract sql statements by ';':
//   - ignore the ';' in ``, '', "";
//   - ignore the escaped characters in ``, '', "";
// Note that: Comments cannot have ';'
func ExtractSQL(s string) []string {
	res := make([]string, 0)
	start := 0
	flag := -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			if flag == -1 {
				flag = '\''
			} else {
				if flag == '\'' {
					flag = -1
				}
			}
		case '"':
			if flag == -1 {
				flag = '"'
			} else {
				if flag == '"' {
					flag = -1
				}
			}
		case '`':
			if flag == -1 {
				flag = '`'
			} else {
				if flag == '`' {
					flag = -1
				}
			}
		case '\\':
			if flag != -1 {
				i++
			}
		case ';':
			if flag == -1 {
				res = append(res, s[start:i+1])
				start = i+1
			}
		default:
			continue
		}
	}
	return res
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

func NewSQLFExecutorB(sqlBytes []byte) (*SQLFExecutor, error) {
	startTime := time.Now()
	sqls := ExtractSQL(string(sqlBytes))
	parseErrs := make([]*ParseError, 0)
	asts := make([]ast.StmtNode, 0)
	for i, sql := range sqls {
		p := parser.New()
		stmtNode, _, err := p.Parse(sql, "", "")
		if err != nil {
			parseErrs = append(parseErrs, &ParseError{
				Id: i,
				Sql: sql,
				Err: err,
			})
			continue
		}
		if stmtNode == nil || len(stmtNode) != 1 {
			parseErrs = append(parseErrs, &ParseError{
				Id: i,
				Sql: sql,
				Err: errors.New("stmtNode == nil || len(stmtNode) != 1"),
			})
			continue
		}
		asts = append(asts, stmtNode[0])
	}

	parseTime := time.Since(startTime)

	return &SQLFExecutor{
		ParseErrs: parseErrs,
		ASTs: asts,
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