package sqlsexecutor

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

// SQLSExecutor: Read .sql file(MySQL) file or sqls []string, parse each sql to ast, execute them, get the results
type SQLSExecutor struct {
	ParseErrs []*ParseError // some sql statements may have syntax error
	ASTs []ast.StmtNode
	Results []*connector.Result
	ReadTime time.Duration // total time of reading .sql file
	ParseTime time.Duration // total time of parsing sqls to ASTs
	ExecuteTime time.Duration // total execute time
	PassedSQLNum int // the number of passed sql
	FailedSQLNum int // the number of failed sql
	EmptyRowNumOfPRes int // the number of passed sql witch return empty result
	SumRowNumOfPRes int // The total number of rows of passed Results
}

func (sqlsExecutor *SQLSExecutor) ToShortString() string {
	str := ""
	if sqlsExecutor.ParseErrs != nil && len(sqlsExecutor.ParseErrs) != 0 {
		str += "|Parse error|: " + strconv.Itoa(len(sqlsExecutor.ParseErrs)) + "\n"
		str += "========================================\n"
		for _, parseError := range sqlsExecutor.ParseErrs {
			str += parseError.ToString() + "\n"
		}
		str += "========================================\n"
	}
	if sqlsExecutor.ASTs == nil || len(sqlsExecutor.ASTs) == 0 {
		return str + "|ASTs| = 0"
	}
	str += "|ASTs|: " + strconv.Itoa(len(sqlsExecutor.ASTs)) + "\n"
	str += "Read Time: " + sqlsExecutor.ReadTime.String() + "\n"
	str += "Parse Time: " +  sqlsExecutor.ParseTime.String()
	if sqlsExecutor.Results == nil || len(sqlsExecutor.Results) == 0 {
		return str + "\n|Results| = 0"
	}
	str += "\nExec Time: " + sqlsExecutor.ExecuteTime.String() + "\n"
	str += "Passed SQL Num: " + strconv.Itoa(sqlsExecutor.PassedSQLNum) + "\n"
	str += "Failed SQL Num: " + strconv.Itoa(sqlsExecutor.FailedSQLNum) + "\n"
	str += "Empty RowNum of Passed Results: " + strconv.Itoa(sqlsExecutor.EmptyRowNumOfPRes) + "\n"
	str += "Sum RowNum of Passed Results: " + strconv.Itoa(sqlsExecutor.SumRowNumOfPRes)
	str += "\nFailed SQLs:"
	for i, result := range sqlsExecutor.Results {
		if result.Err != nil {
			str += "\n==================================================\n"
			str += "[sql "+strconv.Itoa(i)+"]: " + sqlsExecutor.ASTs[i].Text() + "\n"
			str += "[result "+strconv.Itoa(i)+"]: " + result.ToString() + "\n"
			str += "=================================================="
		}
	}
	return str
}

func (sqlsExecutor *SQLSExecutor) ToString() string {
	str := sqlsExecutor.ToShortString()
	str += "\nPassed SQLs(empty):"
	for i, result := range sqlsExecutor.Results {
		if result.Err == nil && len(result.Rows) == 0 {
			str += "\n==================================================\n"
			str += "[sql "+strconv.Itoa(i)+"]: " + sqlsExecutor.ASTs[i].Text() + "\n"
			str += "[result "+strconv.Itoa(i)+"]: " + result.ToString() + "\n"
			str += "=================================================="
		}
	}
	str += "\nPassed SQLs:"
	for i, result := range sqlsExecutor.Results {
		if result.Err == nil && len(result.Rows) != 0 {
			str += "\n==================================================\n"
			str += "[sql "+strconv.Itoa(i)+"]: " + sqlsExecutor.ASTs[i].Text() + "\n"
			str += "[result "+strconv.Itoa(i)+"]: " + result.ToString() + "\n"
			str += "=================================================="
		}
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

// NewSQLSExecutor: create SQLSExecutor from .sql file
func NewSQLSExecutor(filePath string) (*SQLSExecutor, error) {
	startTime := time.Now()
	sqls, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("NewSQLSExecutor: read sqls file error: " + err.Error())
	}
	readTime := time.Since(startTime)
	sqlSExecutor, err := NewSQLSExecutorB(sqls)
	if err != nil {
		return nil, errors.New("NewSQLSExecutor: " + err.Error())
	}
	sqlSExecutor.ReadTime = readTime
	return sqlSExecutor, nil
}

// NewSQLSExecutorB: create SQLSExecutor from bytes
func NewSQLSExecutorB(sqlBytes []byte) (*SQLSExecutor, error) {
	sqls := ExtractSQL(string(sqlBytes))
	return NewSQLSExecutorS(sqls)
}

// NewSQLSExecutorS: create SQLSExecutor from []string
func NewSQLSExecutorS(sqls []string) (*SQLSExecutor, error) {
	startTime := time.Now()
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

	return &SQLSExecutor{
		ParseErrs: parseErrs,
		ASTs: asts,
		ParseTime: parseTime,
	}, nil
}

func (sqlsExecutor *SQLSExecutor) Exec(conn *connector.Connector) {
	startTime := time.Now()
	sqlsExecutor.PassedSQLNum = 0
	sqlsExecutor.FailedSQLNum = 0
	sqlsExecutor.EmptyRowNumOfPRes = 0
	sqlsExecutor.SumRowNumOfPRes = 0
	for _, AST := range sqlsExecutor.ASTs {
		result := conn.ExecSQL(AST.Text())
		if result.Err != nil {
			sqlsExecutor.FailedSQLNum += 1
		} else {
			sqlsExecutor.PassedSQLNum += 1
			if len(result.Rows) == 0 {
				sqlsExecutor.EmptyRowNumOfPRes += 1
			}
			sqlsExecutor.SumRowNumOfPRes += len(result.Rows)
		}
		sqlsExecutor.Results = append(sqlsExecutor.Results, result)
	}
	sqlsExecutor.ExecuteTime = time.Since(startTime)
}