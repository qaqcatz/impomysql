package stage1

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/qaqcatz/impomysql/connector"
)

type InitResult struct {
	InitSql string
	Err error
	ExecResult *connector.Result
}

// Init: for the input sql, remove aggregate function(and group by),
// window function, LEFT|RIGHT JOIN, Limit.
//
// Note that:
//
// (1) The transformed sql may fail to execute.
//
// (2) we only Support SELECT statement.
func Init(sql string) *InitResult {
	initResult := &InitResult{
		InitSql: "",
		Err: nil,
		ExecResult: nil,
	}
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		initResult.Err = errors.Wrap(err, "[Init]parse error")
		return initResult
	}
	if stmtNodes == nil || len(stmtNodes) == 0 {
		initResult.Err = errors.New("[Init]stmtNodes == nil || len(stmtNodes) == 0 ")
		return initResult
	}
	rootNode := &stmtNodes[0]

	switch (*rootNode).(type) {
	case *ast.SelectStmt:
	case *ast.SetOprStmt:
	default:
		initResult.Err = errors.New("[Init]*rootNode is not *ast.SelectStmt or *ast.SetOprStmt")
		return initResult
	}

	v := &InitVisitor{}
	(*rootNode).Accept(v)

	buf := new(bytes.Buffer)
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err = (*rootNode).Restore(ctx)
	if err != nil {
		initResult.Err = errors.Wrap(err, "[Init]restore error")
		return initResult
	}
	initResult.InitSql = buf.String()
	return initResult
}

// InitAndExec: Init + exec
func InitAndExec(sql string, conn *connector.Connector) *InitResult {
	initResult := Init(sql)
	if initResult.Err != nil {
		return initResult
	}
	result := conn.ExecSQL(initResult.InitSql)
	initResult.ExecResult = result
	return initResult
}